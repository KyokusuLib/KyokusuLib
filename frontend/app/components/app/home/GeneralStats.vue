<script setup lang="ts">
import { computed } from 'vue'
import type { GeneralStatistics } from '@/types/backend/statistics'

const props = defineProps({
  stats: {
    type: Object as () => GeneralStatistics,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
})

const blocks = computed(() => [
  {
    label: 'Закладок',
    value: props.stats?.bookmarkCount ?? 0,
    icon: 'ph:bookmark-simple-bold',
  },
  {
    label: 'Прочитанно',
    value: props.stats?.readCount ?? 0,
    icon: 'ph:eye-bold',
  },
  {
    label: 'Оценок',
    value: props.stats?.ratingCount ?? 0,
    icon: 'ph:star-bold',
  },
  {
    label: 'Комментариев',
    value: props.stats?.commentCount ?? 0,
    icon: 'ph:chat-circle-bold',
  },
  {
    label: 'Томов',
    value: props.stats?.volumeCount ?? 0,
    icon: 'ph:books-bold',
  },
  {
    label: 'Глав',
    value: props.stats?.chapterCount ?? 0,
    icon: 'ph:file-text-bold',
  },
])

function formatNumber(n: number): string {
  if (n >= 1_000_000) return (n / 1_000_000).toFixed(1).replace(/\.0$/, '') + 'M'
  if (n >= 1_000) return (n / 1_000).toFixed(1).replace(/\.0$/, '') + 'K'
  return n.toLocaleString()
}
</script>

<template>
  <div
    class="bg-white dark:bg-zinc-900/40 rounded-[2.5rem] border border-zinc-200 dark:border-zinc-800 p-6 shadow-sm"
  >
    <div class="border-t border-zinc-200 dark:border-zinc-700 mb-6"></div>

    <div v-if="isLoading" class="flex justify-center py-12">
      <Icon name="ph:spinner-gap-bold" size="28" class="animate-spin text-zinc-400" />
    </div>

    <div v-else class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3 sm:gap-4">
      <div
        v-for="block in blocks"
        :key="block.label"
        class="flex flex-col items-center justify-center gap-3 rounded-2xl px-4 py-6 ring-1 ring-zinc-200/50 dark:ring-zinc-700/50 bg-zinc-50 dark:bg-zinc-800/30 transition-colors"
      >
        <div
          class="flex items-center justify-center w-11 h-11 rounded-xl bg-zinc-200 dark:bg-zinc-700 text-zinc-500 dark:text-zinc-400"
        >
          <Icon :name="block.icon" size="22" />
        </div>

        <span class="text-2xl sm:text-3xl font-black tracking-tight text-zinc-800 dark:text-white">
          {{ formatNumber(block.value) }}
        </span>

        <span
          class="text-xs sm:text-sm font-medium text-zinc-500 dark:text-zinc-400 text-center leading-tight"
        >
          {{ block.label }}
        </span>
      </div>
    </div>
  </div>
</template>
