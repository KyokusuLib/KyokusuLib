import { $api } from '@/composables/api/useApi'
import type { NovelaLikeRequest } from '~/types/frontend/novela/novela-like'

export function useNovelaLike() {
  const loading = ref(false)

  const setNovelaLike = async (payload: NovelaLikeRequest) => {
    loading.value = true
    try {
      const data = $api('/api/novela/like', {
        method: 'POST',
        body: payload,
      })

      return data
    } catch (e) {
      console.error('Failed to set like:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    setNovelaLike,
  }
}
