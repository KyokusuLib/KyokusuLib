import { computed } from 'vue'
import { storeToRefs } from 'pinia'
import { NEWS_SOURCES } from '@/constants/news'
import { useNewsStore } from '@/stores/news'
import type { NewsSource } from '@/types/frontend/news'

export function useNewsSources() {
  const newsStore = useNewsStore()
  const { selectedSources } = storeToRefs(newsStore)
  const isLoading = useState('news-sources-loading', () => true)

  const loadSources = async () => {
    isLoading.value = true
    await new Promise((resolve) => setTimeout(resolve, 600))
    isLoading.value = false
  }

  const selectedSourceDetails = computed<NewsSource[]>(() =>
    NEWS_SOURCES.filter((source) => selectedSources.value.includes(source.id)),
  )

  const hasSelectedSources = computed(() => selectedSources.value.length > 0)

  return {
    sources: NEWS_SOURCES,
    selectedSources,
    selectedSourceDetails,
    hasSelectedSources,
    isLoading,
    loadSources,
    isSelected: newsStore.isSelected,
    toggleSource: newsStore.toggleSource,
    clearSources: newsStore.clearSources,
  }
}
