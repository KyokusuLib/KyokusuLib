<script setup lang="ts">
import { Separator } from '@kyokusu-ui/vue'

defineProps<{
  currentPage: number
  hasMore: boolean
  loading: boolean
  total: number
}>()

const emit = defineEmits<{
  (e: 'prev'): void
  (e: 'next'): void
}>()
</script>

<template>
  <div class="flex flex-col-reverse sm:flex-row sm:items-center sm:justify-between gap-3">
    <div class="flex gap-2">
      <button
        :disabled="currentPage <= 1 || loading"
        class="flex-1 sm:flex-none flex items-center justify-center gap-2 px-4 py-2.5 sm:py-2 rounded-xl cursor-pointer text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-900 dark:text-zinc-100"
        @click="emit('prev')"
      >
        <Icon name="ph:arrow-left" class="w-3.5 h-3.5" />
        <span>Назад</span>
      </button>

      <button
        :disabled="!hasMore || loading"
        class="flex-1 sm:flex-none flex items-center justify-center gap-2 px-4 py-2.5 sm:py-2 rounded-xl cursor-pointer text-sm font-medium transition-colors disabled:opacity-50 disabled:cursor-not-allowed bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-900 dark:text-zinc-100"
        @click="emit('next')"
      >
        <span>Вперёд</span>
        <Icon name="ph:arrow-right" class="w-3.5 h-3.5" />
      </button>
    </div>

    <div class="flex items-center justify-center gap-4">
      <span class="text-sm font-medium text-zinc-500 dark:text-zinc-400">
        Страница: {{ currentPage }}
      </span>
      <Separator orientation="vertical" />
      <span class="text-sm font-medium text-zinc-500 dark:text-zinc-400"> Всего: {{ total }} </span>
    </div>
  </div>
</template>
