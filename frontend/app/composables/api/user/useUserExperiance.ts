import { $api } from '@/composables/api/useApi'
import type { UserDefinitions } from '@/types/backend/user'

export function useUserExperiance() {
  const getDefinitions = async (): Promise<UserDefinitions | null> => {
    try {
      return await $api<UserDefinitions>(`/user/experiance/definitions`)
    } catch {
      return null
    }
  }

  const updateLevel = async (userId: number, level: number, experience: number): Promise<boolean> => {
    try {
      await $api(`/user/experiance`, {
        method: 'PUT',
        body: { user_id: userId, level, experience },
      })
      return true
    } catch (e) {
      console.error(e)
      return false
    }
  }

  return {
    getDefinitions,
    updateLevel,
  }
}
