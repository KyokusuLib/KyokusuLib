<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useIntersectionObserver } from '@vueuse/core'
import { useTeam } from '~/composables/api/teams/useTeams'
import { staticImage } from '@/utils/str'
import type { CustomRoleNames, TeamMember } from '~/types/frontend/teams'
import TeamMemberCard from '~/components/app/teams/TeamMemberCard.vue'

const route = useRoute()
const router = useRouter()
const slug = computed(() => route.params.id as string)

const { getTeamBySlug, getTeamMembers } = useTeam()

const customRoleNames = computed<CustomRoleNames[]>(() => {
  if (!team.value) return []
  return [
    { value: 'member', label: team.value.role_names.member },
    { value: 'moderator', label: team.value.role_names.moderator },
    { value: 'owner', label: team.value.role_names.owner },
  ]
})

const { data: team, status: teamStatus } = await useAsyncData(`team-${slug.value}`, () =>
  getTeamBySlug(slug.value),
)

const members = ref<TeamMember[]>([])
const isLoading = ref(true)
const isLoadingMore = ref(false)
const hasMore = ref(true)
const offset = ref(0)
const limit = 20

const loadMembers = async (isLoadMore = false) => {
  if (isLoadMore) {
    isLoadingMore.value = true
  } else {
    isLoading.value = true
    offset.value = 0
    members.value = []
  }

  try {
    const newMembers = await getTeamMembers(slug.value, limit, offset.value)
    if (newMembers.length < limit) {
      hasMore.value = false
    } else {
      hasMore.value = true
    }

    if (isLoadMore) {
      members.value.push(...newMembers)
    } else {
      members.value = newMembers
    }
  } catch (e) {
    console.error('Error loading members:', e)
    hasMore.value = false
  } finally {
    isLoading.value = false
    isLoadingMore.value = false
  }
}

const loadMoreSentinel = ref<HTMLElement | null>(null)

useIntersectionObserver(
  loadMoreSentinel,
  ([{ isIntersecting }]) => {
    if (isIntersecting && hasMore.value && !isLoading.value && !isLoadingMore.value) {
      offset.value += limit
      loadMembers(true)
    }
  },
  { threshold: 0.1 },
)

const handleMemberUpdated = (updatedMember: TeamMember) => {
  const idx = members.value.findIndex((m) => m.user.id === updatedMember.user.id)
  if (idx !== -1) {
    members.value = members.value.map((m, i) => (i === idx ? updatedMember : m))
  }
}

onMounted(() => {
  loadMembers(false)
})

useSeoMeta({
  title: () =>
    team.value?.name ? `Участники: ${team.value.name} - KyokusuLib` : 'Участники команды',
})
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 flex flex-col font-sans transition-colors duration-300"
  >
    <div class="w-full max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
      <div class="mb-8 flex items-center gap-4">
        <button
          @click="router.back()"
          class="p-2 sm:px-4 sm:py-2 bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-full hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors flex items-center gap-2 cursor-pointer shadow-sm"
        >
          <Icon name="ph:arrow-left-bold" />
          <span class="hidden sm:inline font-bold text-sm">Назад</span>
        </button>

        <div v-if="team" class="flex items-center gap-3">
          <div
            class="w-10 h-10 rounded-full overflow-hidden border border-zinc-200 dark:border-zinc-700 bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center"
          >
            <img
              v-if="team.avatar_url"
              :src="staticImage(team.avatar_url)"
              class="w-full h-full object-cover"
            />
            <Icon v-else name="ph:users-three-bold" class="text-zinc-400" />
          </div>
          <div>
            <h1
              class="text-xl sm:text-2xl font-black uppercase tracking-tight text-zinc-900 dark:text-white leading-none"
            >
              Участники команды
            </h1>
            <p class="text-sm text-zinc-500 dark:text-zinc-400 font-medium">
              {{ team.name }}
            </p>
          </div>
        </div>
      </div>

      <div v-if="isLoading && members.length === 0" class="flex flex-col gap-4">
        <div
          v-for="i in 10"
          :key="i"
          class="h-20 rounded-2xl bg-zinc-100 dark:bg-zinc-800/50 animate-pulse w-full"
        ></div>
      </div>

      <div v-else-if="members.length > 0" class="flex flex-col gap-4">
        <TeamMemberCard
          v-for="member in members"
          :key="member.user.id"
          :member="member"
          :customRoleNames="customRoleNames"
          :slug="slug"
          @updated="handleMemberUpdated"
        />
      </div>

      <div
        v-else-if="!isLoading"
        class="flex flex-col items-center justify-center py-32 text-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-3xl bg-white dark:bg-[#18181b]"
      >
        <Icon name="ph:users-slash-bold" size="64" class="text-zinc-300 dark:text-zinc-700 mb-4" />
        <h2 class="text-xl font-bold text-zinc-900 dark:text-zinc-100 mb-2">Нет участников</h2>
        <p class="text-sm font-medium text-zinc-500">В этой команде пока нет участников.</p>
      </div>

      <!-- Intersection Observer Sentinel -->
      <div ref="loadMoreSentinel" class="h-20 flex items-center justify-center mt-4">
        <div
          v-if="isLoadingMore"
          class="flex flex-col items-center text-zinc-400 dark:text-zinc-500"
        >
          <Icon name="ph:spinner-gap-bold" size="32" class="animate-spin text-yellow-500 mb-2" />
          <span class="text-xs font-bold uppercase tracking-wider">Загрузка...</span>
        </div>
      </div>
    </div>
  </div>
</template>
