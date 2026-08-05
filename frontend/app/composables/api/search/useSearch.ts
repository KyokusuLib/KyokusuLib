import { ref, watch, computed } from 'vue'
import { useStorage, useDebounceFn } from '@vueuse/core'
import { $api } from '@/composables/api/useApi'
import type { MostSearched, RecentSearch, RecentSearchType } from '@/types/frontend/search/searches'
import { MOST_SEARCHED_DEFAULT } from '~/constants/data'
import type { NovelaAuthorDetails, NovelaDetails } from '~/types/backend/novela'
import type { GetUserDto } from '~/types/backend/user'
import type { Team } from '~/types/frontend/teams'

export type SearchCategory = 'ranobe' | 'users' | 'teams' | 'authors'
export type SearchResultItem = NovelaDetails | GetUserDto | Team | NovelaAuthorDetails | any

export interface UseSearchOptions {
  immediateWatch?: boolean
}

export function useSearch(options: UseSearchOptions = { immediateWatch: false }) {
  const { user } = useAuthStore()
  const isOpen = useState('search-is-open', () => false)
  const searchQuery = useState('search-query', () => '')

  const genres = useState<MostSearched[]>('search-genres', () => [])
  const categories = useState<MostSearched[]>('search-categories', () => [])

  const activeCategory = useState<SearchCategory>('search-category', () => 'ranobe')
  const isSearching = useState('search-is-searching', () => false)
  const searchResults = useState<SearchResultItem[]>('search-results', () => [])

  const recentSearches = useStorage<RecentSearch[]>('recent-searches', [])

  if (recentSearches.value.length > 0 && typeof recentSearches.value[0] === 'string') {
    recentSearches.value = (recentSearches.value as any).map((text: string) => ({
      text,
      type: 'query',
      category: 'ranobe',
    }))
  }

  const popularSearches = ref<MostSearched[]>(MOST_SEARCHED_DEFAULT)

  const hasActiveFilters = computed(() => genres.value.length > 0 || categories.value.length > 0)
  const isValidQuery = computed(() => searchQuery.value.trim().length >= 2)
  const canSearch = computed(() => isValidQuery.value || hasActiveFilters.value)

  const setQueryParam = (item: MostSearched) => {
    const target = item.type === 'genres' ? genres : categories
    if (!target.value.some((v) => v.id === item.id)) {
      target.value.push(item)
    }
  }

  const removeQueryParam = (item: MostSearched) => {
    const target = item.type === 'genres' ? genres : categories
    target.value = target.value.filter((v) => v.id !== item.id)
  }

  const clearFilters = () => {
    genres.value = []
    categories.value = []
  }

  const openSearch = () => {
    isOpen.value = true
  }

  const closeSearch = () => {
    isOpen.value = false
    searchQuery.value = ''
    clearFilters()
    searchResults.value = []
  }

  const addRecentSearch = (item: {
    text: string
    type: RecentSearchType
    id?: number
    category?: SearchCategory
  }) => {
    if (!item.text.trim()) return

    const current = [...recentSearches.value]
    const index = current.findIndex(
      (existing) => existing.text === item.text && existing.type === item.type,
    )

    if (index > -1) current.splice(index, 1)

    current.unshift({
      text: item.text,
      type: item.type,
      id: item.id,
      category: item.category,
    })

    recentSearches.value = current.slice(0, 5)
  }

  const removeRecentSearch = (text: string) => {
    recentSearches.value = recentSearches.value.filter((q) => q.text !== text)
  }

  const clearRecentSearches = () => {
    recentSearches.value = []
  }

  const fetchMostSearched = async () => {
    try {
      const response = await $fetch<{ genres: Array<string>; categories: Array<string> }>(
        '/api/novela/most-searched',
      )
      popularSearches.value = [
        ...response.genres.map((g, index) => ({ id: index, type: 'genres', label: g })),
        ...response.categories.map((c, index) => ({ id: index, type: 'categories', label: c })),
      ]
    } catch (error) {
      console.error(error)
    }
  }

  const performSearch = async () => {
    if (!canSearch.value) {
      searchResults.value = []
      return
    }

    const query = searchQuery.value.trim()
    isSearching.value = true

    try {
      if (activeCategory.value === 'ranobe') {
        const apiQuery: Record<string, any> = {
          search: query,
          limit: 10,
        }

        if (genres.value.length > 0) {
          apiQuery.genres = genres.value.map((v) => v.label).join(',')
        }

        if (categories.value.length > 0) {
          apiQuery.categories = categories.value.map((v) => v.label).join(',')
        }

        const data = await $api<NovelaDetails[]>('/api/novela', {
          query: apiQuery,
        })
        searchResults.value = data || []
      } else {
        if (searchQuery.value.trim().length < 2) {
          searchResults.value = []
          return
        }

        const endpoint = {
          authors: '/api/author',
          users: '/api/user',
          teams: '/api/teams',
        }[activeCategory.value]

        const data = await $api<any[]>(endpoint, {
          query: {
            search: query,
            limit: 10,
            ...(activeCategory.value === 'teams' && user?.id ? { user_id: user.id } : {}),
          },
        })
        searchResults.value = data || []
      }
    } catch (error) {
      console.error('Search failed:', error)
      searchResults.value = []
    } finally {
      isSearching.value = false
    }
  }

  const debouncedSearch = useDebounceFn(performSearch, 500)

  if (options.immediateWatch) {
    watch(
      [searchQuery, activeCategory, genres, categories],
      () => {
        if (canSearch.value) {
          debouncedSearch()
        } else {
          searchResults.value = []
        }
      },
      { deep: true },
    )
  }

  return {
    isOpen,
    openSearch,
    closeSearch,
    searchQuery,
    activeCategory,
    isSearching,
    searchResults,
    recentSearches,
    popularSearches,
    genres,
    categories,
    hasActiveFilters,
    addRecentSearch,
    removeRecentSearch,
    clearRecentSearches,
    setQueryParam,
    removeQueryParam,
    clearFilters,
    performSearch,
    fetchMostSearched,
  }
}
