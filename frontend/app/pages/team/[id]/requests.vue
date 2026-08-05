<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useIntersectionObserver } from '@vueuse/core'
import { useTeam } from '~/composables/api/teams/useTeams'
import { staticImage } from '@/utils/str'
import TeamJoinRequestCard from '~/components/app/teams/TeamJoinRequestCard.vue'
import { useAuthStore } from '@/stores/auth'
import { useRolePermissions } from '@/composables/api/role/useRolePermissions'
import { KyokusuAppRole } from '@/types/enums/role-enum'

const route = useRoute()
const router = useRouter()
const slug = computed(() => route.params.id as string)

const { getTeamBySlug, getTeamJoinRequests, acceptJoinRequest } = useTeam()
const authStore = useAuthStore()
const { hasPermission } = useRolePermissions()

const { data: team, status: teamStatus } = await useAsyncData(`team-${slug.value}`, () =>
  getTeamBySlug(slug.value),
)

const canEdit = computed(() => {
  if (!authStore.isAuthenticated || !authStore.user || !team.value) return false
  return authStore.user.id === team.value.owner_id || hasPermission(KyokusuAppRole.MODERATOR)
})

watchEffect(() => {
  if (teamStatus.value === 'success' && !canEdit.value) {
    router.push(`/team/${slug.value}`)
  }
})

const requests = ref<any[]>([])
const isLoading = ref(true)
const isLoadingMore = ref(false)
const hasMore = ref(true)
const offset = ref(0)
const limit = 20

const loadRequests = async (isLoadMore = false) => {
  if (isLoadMore) {
    isLoadingMore.value = true
  } else {
    isLoading.value = true
    offset.value = 0
    requests.value = []
  }

  try {
    const newRequests = await getTeamJoinRequests(slug.value, limit, offset.value)
    if (newRequests.length < limit) {
      hasMore.value = false
    } else {
      hasMore.value = true
    }

    if (isLoadMore) {
      requests.value.push(...newRequests)
    } else {
      requests.value = newRequests
    }
  } catch (e) {
    console.error('Error loading requests:', e)
    hasMore.value = false
  } finally {
    isLoading.value = false
    isLoadingMore.value = false
  }
}

const handleAccepted = (userId: number) => {
  requests.value = requests.value.filter((req) => req.user.id !== userId)
}

const handleRejected = (userId: number) => {
  requests.value = requests.value.filter((req) => req.user.id !== userId)
}

const isAcceptingAll = ref(false)
const acceptAll = async () => {
  isAcceptingAll.value = true
  try {
    const promises = requests.value.map((req) => acceptJoinRequest(slug.value, req.user.id))
    await Promise.all(promises)

    requests.value = []
  } catch (e) {
    console.error('Failed to accept all:', e)
  } finally {
    isAcceptingAll.value = false
  }
}

const loadMoreSentinel = ref<HTMLElement | null>(null)

useIntersectionObserver(
  loadMoreSentinel,
  ([{ isIntersecting }]) => {
    if (isIntersecting && hasMore.value && !isLoading.value && !isLoadingMore.value) {
      offset.value += limit
      loadRequests(true)
    }
  },
  { threshold: 0.1 },
)

onMounted(() => {
  if (canEdit.value) {
    loadRequests(false)
  }
})

useSeoMeta({
  title: () => (team.value?.name ? `Заявки: ${team.value.name} - KyokusuLib` : 'Заявки в команду'),
})
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 flex flex-col font-sans transition-colors duration-300"
  >
    <div class="w-full max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
      <div class="mb-8 flex items-center justify-between">
        <div class="flex items-center gap-4">
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
                Заявки на вступление
              </h1>
              <p class="text-sm text-zinc-500 dark:text-zinc-400 font-medium">
                {{ team.name }}
              </p>
            </div>
          </div>
        </div>

        <button
          v-if="requests.length > 0"
          @click="acceptAll"
          :disabled="isAcceptingAll"
          class="px-4 py-2 bg-yellow-600 hover:bg-yellow-500 text-white font-bold rounded-full transition-colors flex items-center gap-2 cursor-pointer disabled:opacity-50 text-sm"
        >
          <Icon v-if="isAcceptingAll" name="ph:spinner-gap-bold" class="animate-spin" />
          <Icon v-else name="ph:check-double-bold" />
          <span class="hidden sm:inline">Принять всех</span>
        </button>
      </div>

      <div v-if="isLoading && requests.length === 0" class="flex flex-col gap-4">
        <div
          v-for="i in 10"
          :key="i"
          class="h-20 rounded-2xl bg-zinc-100 dark:bg-zinc-800/50 animate-pulse w-full"
        ></div>
      </div>

      <div v-else-if="requests.length > 0" class="flex flex-col gap-4">
        <TeamJoinRequestCard
          v-for="req in requests"
          :key="req.user.id"
          :request="req"
          :slug="slug"
          @accepted="handleAccepted"
          @rejected="handleRejected"
        />
      </div>

      <div
        v-else-if="!isLoading"
        class="flex flex-col items-center justify-center py-32 text-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-3xl bg-white dark:bg-[#18181b]"
      >
        <Icon
          name="ph:envelope-simple-open-bold"
          size="64"
          class="text-zinc-300 dark:text-zinc-700 mb-4"
        />
        <h2 class="text-xl font-bold text-zinc-900 dark:text-zinc-100 mb-2">Нет заявок</h2>
        <p class="text-sm font-medium text-zinc-500">
          В данный момент нет новых заявок на вступление.
        </p>
      </div>

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
