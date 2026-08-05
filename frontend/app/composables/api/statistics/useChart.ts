import { ref, computed } from 'vue'
import { $api } from '@/composables/api/useApi'
import type { NovelaMonthlySeries } from '@/types/backend/statistics'

const MONTH_NAMES = [
  'Январь',
  'Февраль',
  'Март',
  'Апрель',
  'Май',
  'Июнь',
  'Июль',
  'Август',
  'Сентябрь',
  'Октябрь',
  'Ноябрь',
  'Декабрь',
]

const COLORS = [
  '#3b82f6',
  '#10b981',
  '#f59e0b',
  '#ef4444',
  '#8b5cf6',
  '#ec4899',
  '#14b8a6',
  '#f97316',
  '#6366f1',
  '#84cc16',
]

export const useChart = () => {
  const series = ref<NovelaMonthlySeries[]>([])
  const chartLoading = ref(false)

  const fetchMonthlyStats = async (limit = 5) => {
    chartLoading.value = true
    try {
      const data = await $api<NovelaMonthlySeries[]>(
        `/api/novelas/statistics/monthly?limit=${limit}`,
      )
      series.value = data ?? []
    } catch {
      series.value = []
    } finally {
      chartLoading.value = false
    }
  }

  const chartDatasets = computed(() =>
    series.value.map((s, i) => ({
      label: s.novela.title,
      color: COLORS[i % COLORS.length] ?? '#3b82f6',
    })),
  )

  const chartData = computed(() =>
    MONTH_NAMES.map((label, monthIdx) => ({
      label,
      values: series.value.map((s) => s.monthlyReads[monthIdx] ?? 0),
    })),
  )

  return {
    series,
    chartLoading,
    chartDatasets,
    chartData,
    fetchMonthlyStats,
  }
}
