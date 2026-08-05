import { reactive, computed } from 'vue'
import { useDebounceFn } from '@vueuse/core'

import { $api } from '@/composables/api/useApi'
import { useNotificationStore } from '@/stores/notification'

import { hasSpecialSymbols } from '@/utils/str'

import type { NovelaFilterState } from '@/types/frontend/query/novela-filters'
import { DEFAULT_FILTER_STATE } from '@/types/frontend/query/novela-filters'
import type { CatalogFilterPreset } from '@/types/backend/catalog-filters'

const filters = reactive<NovelaFilterState>({
  ...DEFAULT_FILTER_STATE,
})

export function useNovelaFilters() {
  const { isAuthenticated } = useAuthStore()
  const savedFilters = ref<CatalogFilterPreset[]>([])

  const { notify } = useNotificationStore()

  const hasSavedFilters = computed(() => savedFilters.value.length > 0)

  const hasActiveFilters = computed(() => {
    return (
      filters.search !== '' ||
      filters.genres.length > 0 ||
      filters.categories.length > 0 ||
      filters.country.length > 0 ||
      filters.status !== '' ||
      filters.translationStatus !== '' ||
      filters.type !== '' ||
      filters.chaptersFrom !== null ||
      filters.chaptersTo !== null ||
      filters.yearFrom !== null ||
      filters.yearTo !== null ||
      filters.ageRating !== '' ||
      filters.sort !== ''
    )
  })

  const activeFilterCount = computed(() => {
    let count = 0
    if (filters.search !== '') count++
    count += filters.genres.length
    count += filters.categories.length
    count += filters.country.length
    if (filters.status !== '') count++
    if (filters.translationStatus !== '') count++
    if (filters.type !== '') count++
    if (filters.ageRating !== '') count++
    if (filters.sort !== '') count++
    if (filters.chaptersFrom !== null || filters.chaptersTo !== null) count++
    if (filters.yearFrom !== null || filters.yearTo !== null) count++
    return count
  })

  function buildQueryParams() {
    const params: Record<string, any> = {}

    if (filters.search) params.search = filters.search
    if (filters.genres.length > 0) params.genres = filters.genres
    if (filters.categories.length > 0) params.categories = filters.categories
    if (filters.country.length > 0) params.country = filters.country
    if (filters.status) params.status = filters.status
    if (filters.translationStatus) params.translation_status = filters.translationStatus
    if (filters.type) params.type = filters.type
    if (filters.ageRating) params.ageRating = filters.ageRating
    if (filters.sort) params.sort = filters.sort
    if (filters.chaptersFrom) params.chapters_from = filters.chaptersFrom
    if (filters.chaptersTo) params.chapters_to = filters.chaptersTo
    if (filters.yearFrom) params.year_from = filters.yearFrom
    if (filters.yearTo) params.year_to = filters.yearTo

    return params
  }

  function resetFilters() {
    Object.assign(filters, DEFAULT_FILTER_STATE)
  }

  function setGenre(genre: string, selected: boolean) {
    if (selected) {
      if (!filters.genres.includes(genre)) {
        filters.genres.push(genre)
      }
    } else {
      filters.genres = filters.genres.filter((g) => g !== genre)
    }
  }

  function setCategory(category: string, selected: boolean) {
    if (selected) {
      if (!filters.categories.includes(category)) {
        filters.categories.push(category)
      }
    } else {
      filters.categories = filters.categories.filter((c) => c !== category)
    }
  }

  async function saveFilters(name: string) {
    name = name.trim()
    if (name === '' || name === null || name === undefined) {
      notify({
        title: 'Внимание!',
        content: 'Название должно содержать мин. один символ.',
        type: 'warning',
      })
      return
    }

    if (hasSpecialSymbols(name)) {
      notify({
        title: 'Внимание!',
        content: 'В название не должно быть спец. символов',
        type: 'warning',
      })
      return
    }

    if (!hasActiveFilters) {
      notify({
        title: 'Внимание!',
        content: 'Вы не выбрали неодного фильтра, что бы сохранить шаблон выберите хотя бы один.',
        type: 'warning',
      })
      return
    }

    try {
      const created = await $api<CatalogFilterPreset>('/api/catalog/filters', {
        method: 'POST',
        body: {
          name: name,
          filters: filters,
        },
      })

      savedFilters.value.push(created)

      notify({
        title: 'Успех',
        content: 'Шаблон поиска сохранён',
        type: 'success',
      })
    } catch (e) {
      console.error('Failed to save filter preset:', e)
      notify({
        title: 'Неудача',
        content: `Произошла ошибка: ${e.message}`,
        type: 'error',
      })
      throw e
    }
  }

  async function getFilters() {
    if (!isAuthenticated) return

    try {
      const presets = await $api<CatalogFilterPreset[]>('/api/catalog/filters')
      savedFilters.value = presets ?? []
    } catch (e) {
      console.error('Failed to fetch filter presets:', e)
      notify({
        title: 'Неудача',
        content: 'Не удалось получить шаблоны поиска',
        type: 'error',
      })
    }
  }

  async function deleteFilterPreset(id: number) {
    try {
      await $api(`/api/catalog/filters/${id}`, { method: 'DELETE' })
      savedFilters.value = savedFilters.value.filter((p) => p.id !== id)
      notify({
        title: 'Успех',
        content: 'Шаблон удалён',
        type: 'success',
      })
    } catch (e) {
      console.error('Failed to delete filter preset:', e)
      notify({
        title: 'Неудача',
        content: `Не удалось удалить шаблон: ${e.message}`,
        type: 'error',
      })
    }
  }

  function loadPreset(preset: CatalogFilterPreset) {
    Object.assign(filters, {
      ...DEFAULT_FILTER_STATE,
      ...preset.filters,
    })
  }

  const debouncedSearch = useDebounceFn((value: string) => {
    filters.search = value
  }, 400)

  return {
    filters,
    savedFilters,
    hasSavedFilters,
    hasActiveFilters,
    activeFilterCount,
    buildQueryParams,
    resetFilters,
    setGenre,
    setCategory,
    saveFilters,
    getFilters,
    deleteFilterPreset,
    loadPreset,
    debouncedSearch,
  }
}
