import { useApi } from '@/composables/api/useApi'
import type { UserDefinitions } from '@/types/backend/user'

export function useUserExperiance() {
  const getDefinitions = async (): Promise<UserDefinitions | null> => {
    try {
      const { data, error } = await useApi<UserDefinitions>(`/api/user/experiance/definitions`)

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
