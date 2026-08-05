<script setup lang="ts">
import SearchResults from '@/components/common/Search/SearchResults.vue'
import SearchInputBlock from '@/components/app/search/SearchInputBlock.vue'
import SearchStateLoading from '@/components/app/search/SearchStateLoading.vue'
import SearchStateEmpty from '@/components/app/search/SearchStateEmpty.vue'
import SearchStateInitial from '@/components/app/search/SearchStateInitial.vue'
import { useSearchPage } from '@/composables/api/search/useSearchPage'

const {
  query,
  activeCategory,
  isLoading,
  isLoadingMore,
  results,
  hasMore,
  loadMoreSentinel,
  handleSelectResult,
} = useSearchPage()

useSeoMeta({
  title: () => (query.value ? `Поиск: ${query.value} - KyokusuLib` : 'Поиск - KyokusuLib'),
})
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 flex flex-col font-sans transition-colors duration-300"
  >
    <div class="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
      <div class="mb-8 space-y-2">
        <h1
          class="text-3xl sm:text-4xl font-black uppercase tracking-tight text-zinc-900 dark:text-white"
        >
          Расширенный поиск
        </h1>
        <p class="text-sm text-zinc-500 dark:text-zinc-400">
          Найдите то, что вам нужно, используя глубокий поиск по всем разделам сайта.
        </p>
      </div>

      <SearchInputBlock v-model:query="query" v-model:activeCategory="activeCategory" />

      <div class="min-h-100">
        <SearchStateLoading v-if="isLoading" />

        <div v-else-if="results.length > 0" class="space-y-4">
          <div
            class="flex items-center justify-between pb-4 border-b border-zinc-200 dark:border-zinc-800"
          >
            <h2 class="text-lg font-bold text-zinc-900 dark:text-white flex items-center gap-2">
              <Icon name="ph:list-dashes-bold" class="text-yellow-500" />
              Результаты поиска
            </h2>
            <span
              class="px-3 py-1 rounded-full bg-zinc-100 dark:bg-zinc-800 text-xs font-bold text-zinc-500"
            >
              Найдено: {{ results.length }}
            </span>
          </div>

          <div class="grid grid-cols-1 gap-4 w-full">
            <SearchResults
              :results="results"
              :activeCategory="activeCategory"
              @select="handleSelectResult"
              class="contents"
            />
          </div>

          <div ref="loadMoreSentinel" class="h-10 w-full mt-4 flex items-center justify-center">
            <Icon
              v-if="isLoadingMore"
              name="ph:spinner-gap-bold"
              size="24"
              class="animate-spin text-yellow-500"
            />
          </div>
        </div>

        <SearchStateEmpty v-else-if="query.length >= 2" :query="query" />

        <SearchStateInitial v-else />
      </div>
    </div>
  </div>
</template>
