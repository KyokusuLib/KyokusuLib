import { $api } from '@/composables/api/useApi'
import type {
  TotalNovelaStatistics,
  TotalNovelaStatisticsResponse,
} from '@/types/backend/statistics'
import type { StatisticsPeriodEnum } from '@/types/enums/statistics'

export const useTotalStatistics = () => {
  const { notify } = useNotificationStore()
  const novelaTotalStats = ref<TotalNovelaStatistics[]>([])
  const isLoading = ref<boolean>(false)
  const total = ref<number>(0)

  const fetchTotalStatistics = async (
    period?: StatisticsPeriodEnum,
    limit?: number,
    offset?: number,
  ) => {
    isLoading.value = true
    try {
      const response = await $api<TotalNovelaStatisticsResponse>(`/api/novelas/statistics/total`, {
        query: { period, limit, offset },
      })

      if (!response.data || response.data.length === 0) {
        notify({
          title: 'Нет данных',
          content: 'Нет данных для отображения.',
          type: 'warning',
        })
        novelaTotalStats.value = []
        total.value = 0
      } else {
        novelaTotalStats.value = response.data
        total.value = response.total
      }
    } catch (error) {
      console.error(error)
      notify({
        title: 'Ошибка при загрузке данных',
        content: 'Не удалось загрузить данные. Пожалуйста, попробуйте позже.',
        type: 'error',
      })
    } finally {
      isLoading.value = false
    }
  }

  return {
    novelaTotalStats,
    total,
    isLoading,
    fetchTotalStatistics,
  }
}
