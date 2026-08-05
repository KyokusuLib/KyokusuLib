<script setup lang="ts">
import { ref } from 'vue'
import type { NovelaDetails } from '@/types/backend/novela'
import CatalogNovelaCard from '@/components/app/catalog/CatalogNovelaCard.vue'
import CatalogNovelaPopover from '@/components/app/catalog/CatalogNovelaPopover.vue'

defineProps<{
  novels: NovelaDetails[]
  isInitialLoading: boolean
  isFetchingMore: boolean
  hasMore: boolean
}>()

const emit = defineEmits<{
  toggleFilters: []
  intersect: []
}>()

const loadMoreTrigger = ref<HTMLElement | null>(null)

useIntersectionObserver(
  loadMoreTrigger,
  ([{ isIntersecting }]) => {
    if (isIntersecting) {
      emit('intersect')
    }
  },
  { rootMargin: '300px' },
)
</script>

<template>
  <div
    class="min-h-150 bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 rounded-xl p-6 shadow-sm flex flex-col"
  >
    <div class="flex items-center justify-between mb-6">
      <button
        type="button"
        class="lg:hidden flex items-center gap-2 px-4 py-1.5 bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 rounded-xl text-sm font-bold hover:bg-zinc-800 dark:hover:bg-zinc-200 transition-all active:scale-95 cursor-pointer"
        @click="emit('toggleFilters')"
      >
        <Icon name="ph:funnel-bold" size="16" />
        <span>Фильтры</span>
      </button>
    </div>

    <div v-if="isInitialLoading" class="grow flex items-center justify-center py-20">
      <div
        class="w-10 h-10 border-3 border-zinc-300 border-t-yellow-500 rounded-full animate-spin"
      />
    </div>

    <div
      v-else-if="novels.length === 0"
      class="grow flex flex-col items-center justify-center py-20 text-zinc-400"
    >
      <Icon name="ph:books-bold" size="48" class="mb-4 opacity-50" />
      <p class="text-lg font-bold text-zinc-500 dark:text-zinc-400 mb-1">Новеллы не найдены</p>
      <p class="text-sm">Попробуйте изменить параметры фильтрации</p>
    </div>

    <div v-else class="flex flex-col gap-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
        <div
          v-for="novela in novels"
          :key="novela.id"
          class="relative group/popover h-full justify-self-center md:justify-self-stretch"
        >
          <CatalogNovelaCard :novela="novela" />
          <CatalogNovelaPopover :novela="novela" />
        </div>
      </div>

      <div ref="loadMoreTrigger" class="w-full h-10 flex items-center justify-center mt-4">
        <div
          v-if="isFetchingMore"
          class="w-6 h-6 border-2 border-zinc-300 border-t-yellow-500 rounded-full animate-spin"
        />
        <div
          v-else-if="!hasMore && novels.length > 0"
          class="text-sm text-zinc-500 dark:text-zinc-600"
        >
          Вы просмотрели весь каталог
        </div>
      </div>
    </div>
  </div>
</template>
