<script setup lang="ts">
import { computed } from 'vue'

import { Label } from '@kyokusu-ui/vue'

const yearFrom = defineModel<number | null>('yearFrom', { default: null })
const yearTo = defineModel<number | null>('yearTo', { default: null })

const MIN_YEAR = 1990
const MAX_YEAR = new Date().getFullYear()

const fromVal = computed({
  get: () => yearFrom.value ?? MIN_YEAR,
  set: (val: number) => {
    const clamped = Math.min(Math.max(val, MIN_YEAR), yearTo.value ?? MAX_YEAR)
    yearFrom.value = clamped === MIN_YEAR ? null : clamped
  },
})

const toVal = computed({
  get: () => yearTo.value ?? MAX_YEAR,
  set: (val: number) => {
    const clamped = Math.max(Math.min(val, MAX_YEAR), yearFrom.value ?? MIN_YEAR)
    yearTo.value = clamped === MAX_YEAR ? null : clamped
  },
})

const rangeStyle = computed(() => {
  const left = ((fromVal.value - MIN_YEAR) / (MAX_YEAR - MIN_YEAR)) * 100
  const right = 100 - ((toVal.value - MIN_YEAR) / (MAX_YEAR - MIN_YEAR)) * 100
  return { left: `${left}%`, right: `${right}%` }
})
</script>

<template>
  <div class="space-y-2">
    <Label label="Год выпуска" />

    <div
      class="flex items-center justify-between text-xs text-zinc-500 dark:text-zinc-400 font-medium"
    >
      <span>{{ fromVal }}</span>
      <span class="text-zinc-400 dark:text-zinc-500">—</span>
      <span>{{ toVal }}</span>
    </div>

    <div class="relative h-6 flex items-center">
      <div class="absolute inset-x-0 h-1.5 rounded-full bg-zinc-200 dark:bg-zinc-700" />
      <div class="absolute h-1.5 rounded-full bg-yellow-500" :style="rangeStyle" />

      <!-- Min range input -->
      <input
        type="range"
        :min="MIN_YEAR"
        :max="MAX_YEAR"
        :value="fromVal"
        step="1"
        class="absolute inset-x-0 w-full h-1.5 appearance-none bg-transparent pointer-events-none [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-yellow-500 [&::-webkit-slider-thumb]:shadow-md [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-white dark:[&::-webkit-slider-thumb]:border-zinc-800 [&::-webkit-slider-thumb]:cursor-pointer [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:appearance-none [&::-moz-range-thumb]:w-4 [&::-moz-range-thumb]:h-4 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-yellow-500 [&::-moz-range-thumb]:shadow-md [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-white dark:[&::-moz-range-thumb]:border-zinc-800 [&::-moz-range-thumb]:cursor-pointer z-10"
        @input="fromVal = Number(($event.target as HTMLInputElement).value)"
      />

      <!-- Max range input -->
      <input
        type="range"
        :min="MIN_YEAR"
        :max="MAX_YEAR"
        :value="toVal"
        step="1"
        class="absolute inset-x-0 w-full h-1.5 appearance-none bg-transparent pointer-events-none [&::-webkit-slider-thumb]:pointer-events-auto [&::-webkit-slider-thumb]:appearance-none [&::-webkit-slider-thumb]:w-4 [&::-webkit-slider-thumb]:h-4 [&::-webkit-slider-thumb]:rounded-full [&::-webkit-slider-thumb]:bg-yellow-500 [&::-webkit-slider-thumb]:shadow-md [&::-webkit-slider-thumb]:border-2 [&::-webkit-slider-thumb]:border-white dark:[&::-webkit-slider-thumb]:border-zinc-800 [&::-webkit-slider-thumb]:cursor-pointer [&::-moz-range-thumb]:pointer-events-auto [&::-moz-range-thumb]:appearance-none [&::-moz-range-thumb]:w-4 [&::-moz-range-thumb]:h-4 [&::-moz-range-thumb]:rounded-full [&::-moz-range-thumb]:bg-yellow-500 [&::-moz-range-thumb]:shadow-md [&::-moz-range-thumb]:border-2 [&::-moz-range-thumb]:border-white dark:[&::-moz-range-thumb]:border-zinc-800 [&::-moz-range-thumb]:cursor-pointer z-20"
        @input="toVal = Number(($event.target as HTMLInputElement).value)"
      />
    </div>

    <div class="flex items-center justify-between text-[10px] text-zinc-400 dark:text-zinc-500">
      <span>{{ MIN_YEAR }}</span>
      <span>{{ MAX_YEAR }}</span>
    </div>
  </div>
</template>
