<script setup lang="ts">
import { computed } from 'vue'
import type { NovelaRatingCategory } from '@/types/backend/novela'

const props = defineProps<{
  ratingDetails: NovelaRatingCategory
  isExpanded: boolean
}>()

const emit = defineEmits(['toggle'])

const maxCount = computed(() => {
  if (!props.ratingDetails.nc_items?.length) return 0
  return Math.max(...props.ratingDetails.nc_items.map((item) => item.count))
})

const getPercentage = (count: number) => {
  if (maxCount.value === 0) return 0
  return (count / maxCount.value) * 100
}
</script>

<template>
  <div
    class="w-full bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 rounded-2xl overflow-hidden transition-colors duration-300"
  >
    <button
      @click="emit('toggle')"
      class="w-full flex items-center justify-between p-4 cursor-pointer hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors outline-none group"
    >
      <div class="flex items-center gap-3 text-left">
        <div class="flex flex-col">
          <span class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-500"
            >Рейтинг</span
          >
          <div class="flex items-baseline gap-2">
            <span class="text-2xl font-black text-zinc-900 dark:text-white">
              {{ ratingDetails.total_rating.toFixed(1) }}
            </span>
            <span class="text-sm font-medium text-zinc-400"
              >/ {{ ratingDetails.total }} оценок</span
            >
          </div>
        </div>
      </div>

      <div
        class="w-8 h-8 rounded-full flex items-center justify-center bg-zinc-100 dark:bg-zinc-800 transition-transform duration-300"
        :class="{ 'rotate-180': isExpanded }"
      >
        <Icon name="lucide:chevron-down" class="text-zinc-500 text-lg" />
      </div>
    </button>

    <div
      class="grid transition-[grid-template-rows] duration-300 ease-in-out"
      :class="isExpanded ? 'grid-template-rows-[1fr]' : 'grid-template-rows-[0fr]'"
    >
      <div class="overflow-hidden">
        <div class="px-4 pb-4">
          <div class="pt-4 border-t border-zinc-100 dark:border-zinc-800 space-y-3">
            <div
              v-for="item in ratingDetails.nc_items"
              :key="item.value"
              class="flex items-center gap-1 group/row"
            >
              <div class="flex items-center gap-1.5 min-w-11">
                <Icon name="ph:star-fill" class="text-yellow-500" size="12" />
                <span class="text-xs font-bold text-zinc-500 dark:text-zinc-400">{{
                  item.value
                }}</span>
              </div>

              <div class="flex-1 h-1.5 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden">
                <div
                  class="h-full bg-yellow-500 transition-[width] duration-700 ease-out"
                  :style="{ width: isExpanded ? `${getPercentage(item.count)}%` : '0%' }"
                ></div>
              </div>

              <div
                class="flex items-center justify-end gap-1.5 min-w-12 group-hover/row:translate-x-0.5 transition-transform duration-300"
              >
                <span
                  class="text-[11px] font-black text-zinc-400 group-hover/row:text-zinc-900 dark:group-hover/row:text-zinc-200 transition-colors"
                >
                  {{ item.count }}
                </span>
                <Icon name="lucide:users" class="text-zinc-300 dark:text-zinc-400 text-sm" />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.grid-template-rows-\[0fr\] {
  grid-template-rows: 0fr;
}
.grid-template-rows-\[1fr\] {
  grid-template-rows: 1fr;
}
</style>
