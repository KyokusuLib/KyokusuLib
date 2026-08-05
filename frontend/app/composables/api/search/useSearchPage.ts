import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { $api } from '@/composables/api/useApi'
import { useDebounceFn, useIntersectionObserver } from '@vueuse/core'
import type { SearchCategory, SearchResultItem } from '@/composables/api/search/useSearch'
import { useSearch } from '@/composables/api/search/useSearch'

export function useSearchPage() {
  const route = useRoute()
  const router = useRouter()
  const { addRecentSearch } = useSearch()

  const query = ref((route.query.q as string) || '')
  const activeCategory = ref<SearchCategory>((route.query.type as SearchCategory) || 'ranobe')
  const isLoading = ref(false)
  const isLoadingMore = ref(false)
  const results = ref<SearchResultItem[]>([])
  const limit = 50
  const offset = ref(0)
  const hasMore = ref(true)

  const loadData = async (isLoadMore = false) => {
    if (!query.value || query.value.trim().length < 2) {
      results.value = []
      hasMore.value = false
      return
    }

    if (isLoadMore) {
      isLoadingMore.value = true
    } else {
      isLoading.value = true
      offset.value = 0
      hasMore.value = true
      results.value = []
    }

    try {
      const searchQuery = query.value.trim()
      let newResults: SearchResultItem[] = []

      if (activeCategory.value === 'ranobe') {
        const data = await $api<SearchResultItem[]>('/api/novela', {
          query: { search: searchQuery, limit, offset: offset.value },
        })
        newResults = data || []
      } else if (activeCategory.value === 'authors') {
        const data = await $api<SearchResultItem[]>('/api/author', {
          query: { search: searchQuery, limit, offset: offset.value },
        })
        newResults = data || []
      } else if (activeCategory.value === 'users') {
        const data = await $api<SearchResultItem[]>('/api/user', {
          query: { search: searchQuery, limit, offset: offset.value },
        })
        newResults = data || []
      } else if (activeCategory.value === 'teams') {
        const data = await $api<SearchResultItem[]>('/api/teams', {
          query: { search: searchQuery, limit, offset: offset.value },
        })
        newResults = data || []
      }

      if (newResults.length < limit) {
        hasMore.value = false
      } else {
        hasMore.value = true
      }

      if (isLoadMore) {
        results.value.push(...newResults)
      } else {
        results.value = newResults
      }
    } catch (error) {
      console.error('Full search failed:', error)
      if (!isLoadMore) results.value = []
    } finally {
      isLoading.value = false
      isLoadingMore.value = false
    }
  }

  const debouncedSearch = useDebounceFn(() => {
    loadData(false)
    // Update URL without reloading
    router.replace({
      query: {
        ...route.query,
        q: query.value || undefined,
        type: activeCategory.value,
      },
    })
  }, 500)

  watch([query, activeCategory], () => {
    if (query.value.trim().length >= 2) {
      isLoading.value = true
      debouncedSearch()
    } else {
      results.value = []
      hasMore.value = false
      router.replace({
        query: {
          ...route.query,
          q: undefined,
          type: activeCategory.value,
        },
      })
    }
  })

  onMounted(() => {
    if (query.value.trim().length >= 2) {
      loadData(false)
    }
  })

  const loadMoreSentinel = ref<HTMLElement | null>(null)
  useIntersectionObserver(loadMoreSentinel, ([{ isIntersecting }]) => {
    if (
      isIntersecting &&
      hasMore.value &&
      !isLoading.value &&
      !isLoadingMore.value &&
      query.value.trim().length >= 2
    ) {
      offset.value += limit
      loadData(true)
    }
  })

  const handleSelectResult = (item: SearchResultItem) => {
    addRecentSearch(query.value)

    if (activeCategory.value === 'ranobe') {
      navigateTo(`/novela/${item.id}`)
    } else if (activeCategory.value === 'authors') {
      navigateTo(`/author/${item.id}`)
    } else if (activeCategory.value === 'users') {
      navigateTo(`/profile/${item.id}`)
    } else if (activeCategory.value === 'teams') {
      navigateTo(`/team/${item.slug || item.id}`)
    }
  }

  return {
    query,
    activeCategory,
    isLoading,
    isLoadingMore,
    results,
    hasMore,
    loadMoreSentinel,
    handleSelectResult,
  }
}
