<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { ModalWindow, Button as BaseButton } from '@kyokusu-ui/vue'
import { useSearch, type SearchResultItem } from '@/composables/api/search/useSearch'
import SearchCategories from './Search/SearchCategories.vue'
import SearchHistory from './Search/SearchHistory.vue'
import SearchResults from './Search/SearchResults.vue'

import type { MostSearched, RecentSearch, RecentSearchType } from '@/types/frontend/search/searches'

const {
  isOpen,
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
  setQueryParam,
  removeQueryParam,
  clearFilters,
  clearRecentSearches,
  fetchMostSearched,
} = useSearch({ immediateWatch: true })

const searchInputRef = ref<HTMLInputElement | null>(null)

const handleSelectResult = (item: SearchResultItem) => {
  const typeMap: Record<string, RecentSearchType> = {
    ranobe: 'query',
    authors: 'author',
    users: 'user',
    teams: 'team',
  }

  addRecentSearch({
    text: searchQuery.value || item.title || item.name || item.username,
    type: typeMap[activeCategory.value] || 'query',
    category: activeCategory.value as any,
  })

  closeSearch()

  const routes: Record<string, string> = {
    ranobe: `/novela/${item.id}`,
    authors: `/author/${item.id}`,
    users: `/profile/${item.id}`,
    teams: `/team/${item.slug || item.id}`,
  }

  navigateTo(routes[activeCategory.value])
}

const handleSelectRecent = (item: RecentSearch | MostSearched) => {
  if ('text' in item) {
    if (item.type === 'genre' || item.type === 'category') {
      setQueryParam({
        id: item.id!,
        label: item.text,
        type: item.type === 'genre' ? 'genres' : 'categories',
      })
    } else {
      searchQuery.value = item.text
      activeCategory.value =
        (item.category as any) || (item.type === 'query' ? 'ranobe' : item.type + 's')

      if (item.type === 'author') activeCategory.value = 'authors'
      if (item.type === 'user') activeCategory.value = 'users'
      if (item.type === 'team') activeCategory.value = 'teams'
    }
  } else {
    // From Popular Searches
    setQueryParam(item)

    // Add popular search to recent history when selected
    addRecentSearch({
      text: item.label,
      type: item.type === 'genres' ? 'genre' : 'category',
      id: item.id,
    })
  }
}

const navigateToSearch = () => {
  const query = searchQuery.value.trim()
  if (query.length >= 2 || hasActiveFilters.value) {
    addRecentSearch({
      text: query,
      type: 'query',
      category: activeCategory.value as any,
    })

    closeSearch()
    navigateTo({
      path: '/search',
      query: {
        q: query || undefined,
        type: activeCategory.value,
        genres: genres.value.length ? genres.value.map((g) => g.label).join(',') : undefined,
        categories: categories.value.length
          ? categories.value.map((c) => c.label).join(',')
          : undefined,
      },
    })
  }
}

const handleKeydown = (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isOpen.value) {
    closeSearch()
  }
}

watch(isOpen, async (val) => {
  if (val) {
    setTimeout(async () => {
      await fetchMostSearched()
      if (searchInputRef.value) {
        searchInputRef.value.focus()
      }
    }, 150)
  }
})

onMounted(async () => {
  window.addEventListener('keydown', handleKeydown)
  await fetchMostSearched()
  if (isOpen.value && searchInputRef.value) {
    searchInputRef.value.focus()
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<template>
  <ModalWindow v-if="isOpen" v-model="isOpen" title="Поиск" width="w-full max-w-2xl">
    <div
      class="flex flex-col max-h-[80vh] -m-6 sm:-m-8 bg-white dark:bg-[#18181b] rounded-2xl overflow-hidden"
    >
      <div
        class="flex flex-col bg-white dark:bg-[#18181b] border-b border-zinc-200 dark:border-zinc-800"
      >
        <div class="px-4 sm:px-10 pt-4 sm:pt-6 pb-2">
          <div
            class="group relative flex items-center bg-zinc-100 dark:bg-zinc-900/50 hover:bg-zinc-200/50 dark:hover:bg-zinc-800 focus-within:bg-white dark:focus-within:bg-[#121214] border border-zinc-200 dark:border-zinc-700 focus-within:border-yellow-500 dark:focus-within:border-yellow-500/50 rounded-2xl px-4 py-2 sm:py-2 transition-all duration-300 shadow-sm focus-within:shadow-md focus-within:ring-4 focus-within:ring-yellow-500/10"
          >
            <Icon
              name="ph:magnifying-glass-bold"
              size="20"
              class="text-zinc-400 dark:text-zinc-500 group-focus-within:text-yellow-500 transition-colors duration-300 shrink-0"
            />

            <input
              ref="searchInputRef"
              v-model="searchQuery"
              @keydown.enter="navigateToSearch"
              type="text"
              class="flex-1 w-full pl-4 pr-10 sm:pr-28 bg-transparent border-none text-sm sm:text-sm font-bold tracking-tight text-zinc-900 dark:text-zinc-100 placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-0 selection:bg-yellow-500/30"
              placeholder="Что будем искать?"
              autofocus
            />

            <Transition name="fade-fast">
              <button
                v-if="searchQuery"
                @click="searchQuery = ''"
                class="absolute right-3 sm:right-4 p-1.5 text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100 bg-zinc-200/50 hover:bg-zinc-200 dark:bg-zinc-800/80 dark:hover:bg-zinc-700 rounded-xl transition-all cursor-pointer flex items-center justify-center gap-1 active:scale-95"
              >
                <span class="hidden sm:block text-[10px] font-bold uppercase tracking-widest ml-1"
                  >Очистить</span
                >
                <Icon name="ph:x-bold" class="mb-0.5" size="12" />
              </button>
            </Transition>
          </div>
        </div>

        <!-- Categories Component -->
        <div class="px-4 sm:px-6 pb-2">
          <SearchCategories v-model="activeCategory" />
        </div>

        <div v-if="hasActiveFilters" class="px-4 sm:px-10 pb-4 flex justify-center flex-wrap gap-2">
          <div
            v-for="genre in genres"
            :key="genre.id"
            class="px-2 py-1 bg-yellow-500/10 border border-yellow-500/20 rounded-full flex items-center gap-1.5 animate-in fade-in zoom-in duration-200"
          >
            <button
              @click="removeQueryParam(genre)"
              class="flex items-center mb-0.5 cursor-pointer text-yellow-600 dark:text-yellow-500 hover:text-yellow-700 dark:hover:text-yellow-400 transition-colors"
            >
              <span
                class="text-[10px] font-bold text-yellow-600 dark:text-yellow-500 uppercase tracking-wider px-1"
                >{{ genre.label }}</span
              >
              <Icon name="ph:x-bold" size="10" />
            </button>
          </div>
          <div
            v-for="category in categories"
            :key="category.id"
            class="px-2 py-1 bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-lg flex items-center gap-1.5 animate-in fade-in zoom-in duration-200"
          >
            <span
              class="text-[10px] font-bold text-zinc-600 dark:text-zinc-400 uppercase tracking-wider"
              >{{ category.label }}</span
            >
            <button
              @click="removeQueryParam(category)"
              class="flex cursor-pointer mb-0.5 text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200 transition-colors"
            >
              <Icon name="ph:x-bold" size="10" />
            </button>
          </div>
          <button
            @click="clearFilters"
            class="cursor-pointer rounded-2xl text-[10px] font-bold text-zinc-400 bg-zinc-100 hover:bg-red-50 dark:bg-zinc-800/50 dark:hover:bg-zinc-500/10 px-2 py-1 uppercase tracking-wider transition-colors"
          >
            Очистить всё
          </button>
        </div>
      </div>

      <!-- Main Content Area -->
      <div class="flex-1 overflow-y-auto min-h-87.5 p-6 custom-scrollbar relative">
        <!-- History & Recommendations (shown only when text is empty) -->
        <SearchHistory
          v-if="!searchQuery"
          class="mb-8"
          :recentSearches="recentSearches"
          :popularSearches="popularSearches"
          @select="handleSelectRecent"
          @remove="removeRecentSearch"
          @clear="clearRecentSearches"
        />

        <!-- Loading State -->
        <div
          v-if="isSearching"
          class="absolute inset-0 flex flex-col items-center justify-center bg-white/80 dark:bg-[#18181b]/80 backdrop-blur-sm z-10"
        >
          <Icon name="ph:spinner-gap-bold" size="40" class="animate-spin text-yellow-500" />
          <p class="mt-4 text-sm font-bold text-zinc-500 uppercase tracking-widest">Поиск...</p>
        </div>

        <!-- Results State -->
        <SearchResults
          v-if="searchResults.length > 0"
          :results="searchResults"
          :activeCategory="activeCategory"
          @select="handleSelectResult"
        />

        <!-- Empty State (shown only if there is a query or filters, and no results) -->
        <div
          v-else-if="(searchQuery.length >= 2 || hasActiveFilters) && !isSearching"
          class="flex flex-col items-center justify-center py-16 text-center h-full"
        >
          <div
            class="w-20 h-20 rounded-full bg-zinc-50 dark:bg-zinc-800/50 flex items-center justify-center mb-4 border border-zinc-100 dark:border-zinc-800"
          >
            <Icon
              name="ph:magnifying-glass-bold"
              size="32"
              class="text-zinc-300 dark:text-zinc-600"
            />
          </div>
          <p class="text-base font-bold text-zinc-900 dark:text-zinc-100">Ничего не найдено</p>
          <p class="text-sm text-zinc-500 mt-2 max-w-xs">
            По вашему запросу нет результатов. Попробуйте изменить фильтры или текст.
          </p>
        </div>
      </div>

      <!-- Footer -->
      <div
        class="px-6 py-4 mb-4 bg-zinc-50 dark:bg-zinc-900/80 border-t border-zinc-200 dark:border-zinc-800 flex justify-between items-center mt-auto rounded-b-2xl shadow-[0_-4px_10px_rgba(0,0,0,0.02)] dark:shadow-[0_-4px_10px_rgba(0,0,0,0.2)]"
      >
        <BaseButton
          variant="outline"
          @click="closeSearch"
          class="!py-1.5 !px-3 !text-xs text-zinc-500 hover:text-zinc-900 dark:hover:text-zinc-100"
        >
          Закрыть (Esc)
        </BaseButton>
        <div class="flex items-center gap-1.5 text-xs font-bold uppercase tracking-widest">
          <span class="text-zinc-800 dark:text-zinc-200">KyokusuLib</span>
          <Icon name="ph:lightning-fill" class="text-yellow-500" size="14" />
        </div>
      </div>
    </div>
  </ModalWindow>
</template>

<style scoped>
.fade-fast-enter-active,
.fade-fast-leave-active {
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}
.fade-fast-enter-from,
.fade-fast-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #3f3f46;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #fbbf24;
}
</style>
