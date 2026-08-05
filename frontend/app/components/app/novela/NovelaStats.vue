<script setup lang="ts">
import { computed } from 'vue'
import type { NovelaDetails } from '@/types/backend/novela'
import { getStingYear } from '@/utils/date'

const props = defineProps<{
  novela: NovelaDetails
}>()

const formatNumber = (num: number) => {
  if (!num) return 0
  return new Intl.NumberFormat('ru-RU', {
    notation: 'compact',
    compactDisplay: 'short',
  }).format(num)
}

const stats = computed(() => [
  {
    label: 'Просмотры',
    value: formatNumber(props.novela.views),
    icon: 'ph:eye-bold',
  },
  {
    label: 'В закладках',
    value: formatNumber(props.novela.bookmark_details.total),
    icon: 'ph:bookmark-simple-bold',
  },
  {
    label: 'Год',
    value: getStingYear(props.novela.release_date),
    icon: 'ph:calendar-blank-bold',
  },
  {
    label: 'Лайки',
    value: formatNumber(props.novela.like_count),
    icon: 'ph:heart-bold',
  },
])
</script>

<template>
  <div class="flex flex-wrap gap-5 text-sm cursor-default">
    <div
      v-for="stat in stats"
      :key="stat.label"
      class="flex items-center gap-2 text-white dark:text-zinc-100 group transition-colors hover:text-zinc-900 dark:hover:text-zinc-100"
    >
      <Icon
        :name="stat.icon"
        size="18"
        class="opacity-80 group-hover:opacity-100 transition-opacity"
      />

      <div class="flex items-baseline gap-1.5">
        <span class="font-black text-white dark:text-zinc-200">
          {{ stat.value }}
        </span>
        <span class="hidden sm:inline text-[10px] font-bold uppercase tracking-wider opacity-50">
          {{ stat.label }}
        </span>
      </div>
    </div>
  </div>
</template>
