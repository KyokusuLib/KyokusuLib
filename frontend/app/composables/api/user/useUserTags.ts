import { useApi } from '@/composables/api/useApi'
import type { UserTagDefinitions } from '@/types/backend/user'

export function useUserTags() {
  const getDefinitions = async (): Promise<UserTagDefinitions | null> => {
    try {
      const { data, error } = await useApi<UserTagDefinitions>(`/api/user/tags`)

      if (error.value) {
        return null
      }

      return data.value!
    } catch (_e) {
      return null
    }
  }

  return {
    getDefinitions,
  }
}
