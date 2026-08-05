<script setup lang="ts">
import { watch } from 'vue'
import { PageLayout, Separator } from '@kyokusu-ui/vue'

import { useRoute } from 'vue-router'
import { useCatalog } from '@/composables/api/novela/useCatalog'

import CatalogNovelaList from '@/components/app/catalog/CatalogNovelaList.vue'
import CatalogFilterSidebar from '@/components/app/catalog/CatalogFilterSidebar.vue'
import CatalogMobileFilters from '@/components/app/catalog/CatalogMobileFilters.vue'
import { useNovelaFilters } from '~/composables/api/novela/useNovelaFilters'

const route = useRoute()

// Build initial filter params from query BEFORE catalog fetch
const initialFilterParams: Record<string, string> = {}
if (typeof route.query.sort === 'string' && route.query.sort)
  initialFilterParams.sort = route.query.sort
if (typeof route.query.type === 'string' && route.query.type)
  initialFilterParams.type = route.query.type
if (typeof route.query.status === 'string' && route.query.status)
  initialFilterParams.status = route.query.status
if (typeof route.query.translation_status === 'string' && route.query.translation_status)
  initialFilterParams.translation_status = route.query.translation_status
if (typeof route.query.ageRating === 'string' && route.query.ageRating)
  initialFilterParams.ageRating = route.query.ageRating
if (typeof route.query.search === 'string' && route.query.search)
  initialFilterParams.search = route.query.search

const {
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
} = await useCatalog(initialFilterParams)

const { filters } = useNovelaFilters()

// Sync module-level filters with initial query params (for Select display in sidebar)
if (initialFilterParams.sort) filters.sort = initialFilterParams.sort
if (initialFilterParams.type) filters.type = initialFilterParams.type
if (initialFilterParams.status) filters.status = initialFilterParams.status
if (initialFilterParams.translation_status)
  filters.translationStatus = initialFilterParams.translation_status
if (initialFilterParams.ageRating) filters.ageRating = initialFilterParams.ageRating
if (initialFilterParams.search) filters.search = initialFilterParams.search

// Watch for subsequent query changes (NOT immediate — avoids hydration clash)
watch(
  () => route.query,
  async (query) => {
    filters.sort = ''
    filters.type = ''
    filters.status = ''
    filters.translationStatus = ''
    filters.ageRating = ''
    filters.search = ''

    let hasQuery = false

    if (typeof query.sort === 'string' && query.sort) {
      filters.sort = query.sort
      hasQuery = true
    }
    if (typeof query.type === 'string' && query.type) {
      filters.type = query.type
      hasQuery = true
    }
    if (typeof query.status === 'string' && query.status) {
      filters.status = query.status
      hasQuery = true
    }
    if (typeof query.translation_status === 'string' && query.translation_status) {
      filters.translationStatus = query.translation_status
      hasQuery = true
    }
    if (typeof query.ageRating === 'string' && query.ageRating) {
      filters.ageRating = query.ageRating
      hasQuery = true
    }
    if (typeof query.search === 'string' && query.search) {
      filters.search = query.search
      hasQuery = true
    }

    if (hasQuery) {
      await onApplyFilters()
    }
  },
)
</script>

<template>
  <PageLayout>
    <div
      class="min-h-full lg:w-350 bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 transition-colors duration-300 border border-zinc-200 dark:border-zinc-800 rounded-3xl sm:rounded-[2.5rem] shadow-sm"
    >
      <div class="w-full max-w-8xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
        <div class="mb-10">
          <h1
            class="text-3xl sm:text-4xl font-black uppercase tracking-tight text-zinc-900 dark:text-white mb-4 flex items-center gap-4"
          >
            <Icon
              name="ph:books-bold"
              size="36"
              class="text-yellow-400/30 dark:text-yellow-500/20 -rotate-6"
            />
            <span class="relative">
              Каталог
              <Icon
                name="ph:book-open-bold"
                size="44"
                class="absolute -top-2 left-1/2 -translate-x-1/2 text-yellow-400/10 dark:text-yellow-500/10 -z-10 rotate-12"
              />
            </span>
          </h1>
          <Separator />
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-[7fr_2fr] gap-12">
          <CatalogNovelaList
            :novels="allNovels"
            :is-initial-loading="isInitialLoading"
            :is-fetching-more="isFetchingMore"
            :has-more="hasMore"
            @toggle-filters="toggleFilters"
            @intersect="loadNextPage"
          />

          <CatalogFilterSidebar
            @apply="onApplyFilters"
            @search-update="onSearchUpdate"
            @reset="onResetFilters"
          />
        </div>
      </div>
    </div>

    <CatalogMobileFilters
      :is-open="isFilterOpen"
      @close="closeFilters"
      @apply="onApplyFilters"
      @search-update="onSearchUpdate"
      @reset="onResetFilters"
    />
  </PageLayout>
</template>
