import { useApi, $api } from '~/composables/api/useApi'
import type {
  CreateUserActivity,
  GetActivityResponse,
  UserActivity,
} from '@/types/backend/user_activity'
import type { ResponseMessage } from '~/types/backend/response_message'

export function useUserActivity() {
  const ACTIVITIES_LIMIT_QUERY = 20
  const activities = useState<UserActivity[]>('user-activities', () => [])
  const isLoadingActivities = useState('is-loading-activities', () => false)

  const fetchActivities = async (userId?: string | number) => {
    isLoadingActivities.value = true
    const url = userId ? `/api/user/activities/${userId}` : '/api/user/activities'

    try {
      const { data, error } = await useApi<GetActivityResponse>(url, {
        query: { limit: ACTIVITIES_LIMIT_QUERY },
      })

      if (error.value) throw error.value

      const result = data.value?.message || []
      activities.value = result

      return result
    } catch (e) {
      console.error('Failed to fetch activities:', e)
      activities.value = []
      return []
    } finally {
      isLoadingActivities.value = false
    }
  }

  const createUserActivity = async (payload: CreateUserActivity) => {
    try {
      const data = await $api<ResponseMessage>('/api/user/activities', {
        method: 'POST',
        body: payload,
      })

      return data
    } catch (e) {
      console.error('Failed to create activity:', e)
      throw e
    }
  }

  return {
    activities,
    isLoadingActivities,
    fetchActivities,
    createUserActivity,
  }
}
