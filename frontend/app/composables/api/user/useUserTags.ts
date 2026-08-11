import { $api } from '@/composables/api/useApi'
import type { UserTagDefinitions } from '@/types/backend/user'

export function useUserTags() {
  const getDefinitions = async (): Promise<UserTagDefinitions | null> => {
    try {
      return await $api<UserTagDefinitions>(`/user/tags`)
    } catch {
      return null
    }
  }

  return {
    getDefinitions,
  }
}
