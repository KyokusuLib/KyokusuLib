import { $api } from '@/composables/api/useApi'
import type { NovelaRatingRequest } from '~/types/frontend/novela/novela-rating'

export function useNovelaRating() {
  const loading = ref(false)

  const setNovelaRating = async (payload: NovelaRatingRequest) => {
    loading.value = true
    try {
      const data = await $api('/api/novela/rating', {
        method: 'POST',
        body: payload,
      })

      return data
    } catch (e) {
      console.error('Failed to set rating:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    loading,
    setNovelaRating,
  }
}
