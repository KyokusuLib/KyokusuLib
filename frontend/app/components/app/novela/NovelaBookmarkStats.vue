<script setup lang="ts">
import { computed } from 'vue'
import { useBreakpoints, breakpointsTailwind } from '@vueuse/core'
import type { NovelaBookmarkCategory } from '@/types/backend/novela'

const props = defineProps<{
  bookmarkDetails: NovelaBookmarkCategory
  isExpanded: boolean
}>()

const emit = defineEmits(['toggle'])

const breakpoints = useBreakpoints(breakpointsTailwind)
const isMobile = breakpoints.smaller('sm')

const statusConfig: Record<string, { label: string; bg: string; text: string; icon: string }> = {
  reading: {
    label: 'Читают',
    bg: 'bg-emerald-500',
    text: 'text-emerald-500',
    icon: 'ph:book-open-bold',
  },
  completed: {
    label: 'Прочитано',
    bg: 'bg-blue-500',
    text: 'text-blue-500',
    icon: 'ph:check-circle-bold',
  },
  planned: {
    label: 'В планах',
    bg: 'bg-zinc-400',
    text: 'text-zinc-400',
    icon: 'ph:clock-bold',
  },
  on_hold: {
    label: 'Отложено',
    bg: 'bg-amber-500',
    text: 'text-amber-500',
    icon: 'ph:pause-circle-bold',
  },
  dropped: {
    label: 'Брошено',
    bg: 'bg-red-500',
    text: 'text-red-500',
    icon: 'ph:prohibit-bold',
  },
}

const maxCount = computed(() => {
  if (!props.bookmarkDetails.nc_items?.length) return 0
  return Math.max(...props.bookmarkDetails.nc_items.map((item) => item.count))
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
            >В библиотеках</span
          >
          <div class="flex items-baseline gap-2">
            <span class="text-2xl font-black text-zinc-900 dark:text-white">
              {{ bookmarkDetails.total }}
            </span>
            <span class="xs:inline text-sm font-medium text-zinc-400">человек</span>
          </div>
        </div>
      </div>

      <div
        class="w-8 h-8 rounded-full flex items-center justify-center bg-zinc-100 dark:bg-zinc-800 transition-transform duration-300"
        :class="{ 'rotate-180': isExpanded }"
      >
        <Icon name="ph:caret-down-bold" size="18" class="text-zinc-500" />
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
              v-for="item in bookmarkDetails.nc_items"
              :key="item.value"
              class="flex items-center gap-3 group/row"
            >
              <div class="flex items-center gap-2 min-w-fit sm:min-w-[110px]">
                <Icon
                  :name="statusConfig[item.value]?.icon || 'ph:bookmark-simple-bold'"
                  size="16"
                  :class="statusConfig[item.value]?.text"
                />
                <span
                  class="hidden sm:inline text-xs font-bold text-zinc-500 dark:text-zinc-400 whitespace-nowrap"
                >
                  {{ statusConfig[item.value]?.label || item.value }}
                </span>
              </div>

              <div class="flex-1 h-1.5 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden">
                <div
                  class="h-full transition-[width] duration-700 ease-out"
                  :class="statusConfig[item.value]?.bg"
                  :style="{ width: isExpanded ? `${getPercentage(item.count)}%` : '0%' }"
                ></div>
              </div>

              <div
                class="flex items-center justify-end gap-1.5 min-w-[42px] group-hover/row:translate-x-0.5 transition-transform duration-300"
              >
                <span
                  class="text-[11px] font-black text-zinc-400 group-hover/row:text-zinc-900 dark:group-hover/row:text-zinc-200 transition-colors"
                >
                  {{ item.count }}
                </span>
                <Icon name="ph:users-fill" size="14" class="text-zinc-300 dark:text-zinc-600" />
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
