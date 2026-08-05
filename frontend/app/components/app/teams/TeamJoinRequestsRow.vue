<script setup lang="ts">
import { useTeam } from '~/composables/api/teams/useTeams'
import TeamJoinRequestCard from './TeamJoinRequestCard.vue'
import { useAuthStore } from '@/stores/auth'
import { useRolePermissions } from '@/composables/api/role/useRolePermissions'
import { KyokusuAppRole } from '@/types/enums/role-enum'
import { computed, watch } from 'vue'

const props = defineProps<{
  slug: string
  ownerId: number
}>()

const emit = defineEmits(['update:members'])

const { getTeamJoinRequests } = useTeam()
const authStore = useAuthStore()
const { hasPermission } = useRolePermissions()

const canEdit = computed(() => {
  if (!authStore.isAuthenticated || !authStore.user) return false
  return authStore.user.id === props.ownerId || hasPermission(KyokusuAppRole.MODERATOR)
})

const {
  data: requests,
  pending: isLoading,
  refresh,
} = useAsyncData<any[]>(
  `team-join-requests-${props.slug}`,
  async () => {
    if (!canEdit.value) return []
    try {
      return await getTeamJoinRequests(props.slug, 10, 0)
    } catch (e) {
      console.error('Failed to fetch join requests', e)
      return []
    }
  },
  { default: () => [] },
)

watch(canEdit, (newVal) => {
  if (newVal) refresh()
})

const handleAccepted = (userId: number) => {
  if (requests.value) {
    requests.value = requests.value.filter((req) => req.user.id !== userId)
  }
  emit('update:members')
}

const handleRejected = (userId: number) => {
  if (requests.value) {
    requests.value = requests.value.filter((req) => req.user.id !== userId)
  }
}
</script>

<template>
  <div
    v-if="canEdit && (isLoading || requests.length > 0)"
    class="bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 sm:p-8 shadow-sm w-full mt-6"
  >
    <div class="flex items-end justify-between mb-6">
      <h2
        class="flex items-center gap-3 text-xl md:text-2xl font-black tracking-tight text-zinc-900 dark:text-white uppercase italic"
      >
        <div class="flex flex-col items-start gap-1">
          Заявки на вступление
          <div class="h-1 w-1/2 bg-yellow-500 rounded-full"></div>
        </div>
      </h2>
      <div class="flex justify-center sm:justify-end">
        <NuxtLink
          :to="`/team/${slug}/requests`"
          class="text-center px-4 py-2 bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-300 font-bold rounded-full transition-colors text-[12px] duration-300 flex items-center gap-2"
        >
          Все заявки
          <Icon name="ph:arrow-right-bold" />
        </NuxtLink>
      </div>
    </div>

    <div class="mt-2">
      <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="i in 3"
          :key="i"
          class="h-20 rounded-2xl bg-zinc-100 dark:bg-zinc-800/50 animate-pulse"
        ></div>
      </div>
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4">
        <TeamJoinRequestCard
          v-for="req in requests"
          :key="req.user.id"
          :request="req"
          :slug="slug"
          @accepted="handleAccepted"
          @rejected="handleRejected"
        />
      </div>
    </div>
  </div>
</template>
