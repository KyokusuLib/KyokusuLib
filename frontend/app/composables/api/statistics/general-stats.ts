import { ref } from 'vue'
import { $api } from '@/composables/api/useApi'
import type { GeneralStatistics } from '@/types/backend/statistics'

export const useGeneralStatistics = () => {
  const stats = ref<GeneralStatistics | null>(null)
  const isLoading = ref(false)

  const fetchGeneralStatistics = async (): Promise<GeneralStatistics | null> => {
    isLoading.value = true
    try {
      const response = await $api<GeneralStatistics>(`/api/novelas/statistics/general`)
      stats.value = response
      return response
    } catch (error) {
      console.error(error)
      return null
    } finally {
      isLoading.value = false
    }
  }

  return {
    stats,
    isLoading,
    fetchGeneralStatistics,
  }
}
