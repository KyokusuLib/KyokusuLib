<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { Separator, Button } from '@kyokusu-ui/vue'
import { useNovelaFilters } from '@/composables/api/novela/useNovelaFilters'
import FilterSearch from './FilterSearch.vue'
import FilterGenre from './FilterGenre.vue'
import FilterCategory from './FilterCategory.vue'
import FilterStatus from './FilterStatus.vue'
import FilterChapters from './FilterChapters.vue'
import FilterYear from './FilterYear.vue'
import FilterCountry from './FilterCountry.vue'
import {
  GENRE_OPTIONS,
  CATEGORY_OPTIONS,
  COUNTRY_OPTIONS,
} from '@/types/frontend/query/novela-filters'
import type { NovelaFilterOption } from '@/types/frontend/query/novela-filters'

interface Props {
  genreOptions?: NovelaFilterOption[]
  categoryOptions?: NovelaFilterOption[]
  countryOptions?: NovelaFilterOption[]
}

const props = withDefaults(defineProps<Props>(), {
  genreOptions: () => GENRE_OPTIONS,
  categoryOptions: () => CATEGORY_OPTIONS,
  countryOptions: () => COUNTRY_OPTIONS,
})

const emit = defineEmits<{
  apply: [params: Record<string, any>]
  reset: []
  'search-update': [params: Record<string, any>]
}>()

const {
  filters,
  hasActiveFilters,
  activeFilterCount,
  resetFilters,
  buildQueryParams,
  debouncedSearch,
} = useNovelaFilters()

const localSearch = ref('')

function onSearchUpdate(value: string) {
  localSearch.value = value
  debouncedSearch(value)
}

watch(
  () => filters.search,
  () => {
    emit('search-update', buildQueryParams())
  },
)

function handleReset() {
  localSearch.value = ''
  resetFilters()
  emit('reset')
}

function handleApply() {
  emit('apply', buildQueryParams())
}

const handleKeydownApply = async (e: KeyboardEvent) => {
  if (e.key === 'Enter') {
    handleApply()
  }
}

const handleKeydownReset = async (e: KeyboardEvent) => {
  if (e.altKey && e.key.toLowerCase() === 'r') {
    handleReset()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydownApply)
  window.addEventListener('keydown', handleKeydownReset)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydownApply)
  window.removeEventListener('keydown', handleKeydownReset)
})
</script>

<template>
  <div class="space-y-6">
    <FilterSearch v-model="localSearch" @update="onSearchUpdate" />

    <Separator />

    <FilterStatus
      v-model:status="filters.status"
      v-model:translationStatus="filters.translationStatus"
      v-model:type="filters.type"
      v-model:ageRating="filters.ageRating"
      v-model:sort="filters.sort"
    />

    <Separator />

    <FilterGenre
      v-if="genreOptions.length"
      :options="genreOptions"
      :model-value="filters.genres"
      @update:model-value="
        (val) => {
          filters.genres = val
        }
      "
    />

    <template v-if="genreOptions.length && categoryOptions.length">
      <Separator />
    </template>

    <FilterCategory
      v-if="categoryOptions.length"
      :options="categoryOptions"
      :model-value="filters.categories"
      @update:model-value="
        (val) => {
          filters.categories = val
        }
      "
    />

    <template v-if="categoryOptions.length">
      <Separator />
    </template>

    <FilterCountry
      v-if="countryOptions.length"
      :options="countryOptions"
      :modelValue="filters.country"
      @update:modelValue="
        (val) => {
          filters.country = val
        }
      "
    />

    <template v-if="countryOptions.length">
      <Separator />
    </template>

    <FilterChapters
      v-model:chapters-from="filters.chaptersFrom"
      v-model:chapters-to="filters.chaptersTo"
    />

    <Separator />

    <FilterYear v-model:year-from="filters.yearFrom" v-model:year-to="filters.yearTo" />

    <Separator />

    <div class="space-y-2.5">
      <Button
        variant="primary"
        class="w-full flex items-center justify-center gap-4"
        @click="handleApply"
      >
        <Icon name="ph:magnifying-glass" size="18" />
        Применить
      </Button>

      <button
        v-if="hasActiveFilters"
        @click="handleReset"
        class="relative w-full py-2 text-sm font-medium text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300 border border-zinc-850 dark:border-zinc-750 transition-colors cursor-pointer rounded-xl"
      >
        Сбросить фильтры
        <ClientOnly>
          <span
            class="absolute -top-0.5 -right-0.5 flex items-center justify-center min-w-4.5 h-4.5 px-1 rounded-full bg-red-600 text-[10px] font-bold text-white leading-none"
            >{{ activeFilterCount }}</span
          >
        </ClientOnly>
      </button>
    </div>
  </div>
</template>
