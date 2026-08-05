import { ref } from 'vue'
import { $api } from '@/composables/api/useApi'
import type { DashboardUser } from '@/types/frontend/dashboard/users'

interface FetchUsersParams {
  search?: string
  limit?: number
  offset?: number
}

export function useUsers() {
  const searchQuery = ref('')
  const currentOffset = ref(0)
  const PAGE_SIZE = 20

  const {
    pending: isLoading,
    data: users,
    refresh,
  } = useAsyncData<DashboardUser[]>(
    'dashboard-users',
    async () => {
      const params: FetchUsersParams = {
        limit: PAGE_SIZE,
        offset: currentOffset.value,
      }

      if (searchQuery.value.trim()) {
        params.search = searchQuery.value.trim()
      }

      try {
        return await $api<DashboardUser[]>('/api/user', { params })
      } catch {
        return []
      }
    },
    { default: () => [], watch: [searchQuery] },
  )

  const hasMore = computed(() => (users.value?.length ?? 0) >= PAGE_SIZE)

  async function loadMore() {
    if (!hasMore.value || isLoading.value) return
    currentOffset.value += PAGE_SIZE
    await refresh()
  }

  async function onSearch(query: string) {
    searchQuery.value = query
    currentOffset.value = 0
    await refresh()
  }

  async function refreshUsers() {
    currentOffset.value = 0
    await refresh()
  }

  return {
    users,
    isLoading,
    searchQuery,
    hasMore,
    loadMore,
    onSearch,
    refreshUsers,
  }
}
