import { ref } from 'vue'
import type { NovelaDetails } from '@/types/backend/novela'
import type { NovelaFilters } from '@/types/frontend/novela/novela-filters'
import { useNovela } from '@/composables/api/novela/useNovela'
import { NOVELA_CATALOG_PAGE_SIZE } from '@/constants/data'

export async function useCatalog(initialFilters: NovelaFilters = {}) {
  const { fetchNovels } = useNovela()

  const appliedFilterParams = ref<NovelaFilters>({ ...initialFilters })

  const isFilterOpen = ref(false)
  const toggleFilters = () => {
    isFilterOpen.value = !isFilterOpen.value
  }
  const closeFilters = () => {
    isFilterOpen.value = false
  }

  const page = ref(1)
  const isFetchingMore = ref(false)

  const {
    pending: isInitialLoading,
    refresh,
    data: initialData,
  } = await useAsyncData<NovelaDetails[]>(
    'catalog-novels',
    async () => {
      const params = {
        limit: NOVELA_CATALOG_PAGE_SIZE,
        offset: 0,
        ...appliedFilterParams.value,
      }
      const result = await fetchNovels(params)
      return result ?? []
    },
    { default: () => [] },
  )

  const allNovels = ref<NovelaDetails[]>([...initialData.value])
  const hasMore = ref(initialData.value.length === NOVELA_CATALOG_PAGE_SIZE)

  async function loadNextPage() {
    if (isFetchingMore.value || !hasMore.value) return
    isFetchingMore.value = true
    page.value++
    const params = {
      limit: NOVELA_CATALOG_PAGE_SIZE,
      offset: (page.value - 1) * NOVELA_CATALOG_PAGE_SIZE,
      ...appliedFilterParams.value,
    }
    const result = await fetchNovels(params)
    const fetched = result ?? []
    hasMore.value = fetched.length === NOVELA_CATALOG_PAGE_SIZE
    allNovels.value = [...allNovels.value, ...fetched]
    isFetchingMore.value = false
  }

  async function onApplyFilters(filterParams?: Record<string, any>) {
    if (filterParams) {
      appliedFilterParams.value = filterParams
    }
    closeFilters()
    page.value = 1
    hasMore.value = true
    await refresh()
    allNovels.value = [...initialData.value]
    hasMore.value = initialData.value.length === NOVELA_CATALOG_PAGE_SIZE
  }

  async function onSearchUpdate(filterParams: Record<string, any>) {
    appliedFilterParams.value = filterParams
    page.value = 1
    hasMore.value = true
    await refresh()
    allNovels.value = [...initialData.value]
    hasMore.value = initialData.value.length === NOVELA_CATALOG_PAGE_SIZE
  }

  async function onResetFilters() {
    appliedFilterParams.value = {}
    closeFilters()
    page.value = 1
    hasMore.value = true
    await refresh()
    allNovels.value = [...initialData.value]
    hasMore.value = initialData.value.length === NOVELA_CATALOG_PAGE_SIZE
  }

  return {
    isFilterOpen,
    toggleFilters,
    closeFilters,
    hasMore,
    allNovels,
    isFetchingMore,
    isInitialLoading,
    loadNextPage,
    onApplyFilters,
    onSearchUpdate,
    onResetFilters,
  }
}
