<script setup lang="ts">
import { StatisticsPeriodEnum } from '@/types/enums/statistics'

defineProps<{
  modelValue: StatisticsPeriodEnum | undefined
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: StatisticsPeriodEnum | undefined): void
}>()

const periods = [
  { label: 'За всё время', value: undefined, icon: 'ph:clock-bold' },
  { label: 'За день', value: StatisticsPeriodEnum.Day, icon: 'ph:sun-bold' },
  { label: 'За неделю', value: StatisticsPeriodEnum.Week, icon: 'ph:calendar-dots-bold' },
  { label: 'За месяц', value: StatisticsPeriodEnum.Month, icon: 'ph:moon-bold' },
]
</script>

<template>
  <div class="flex flex-col sm:flex-row gap-2">
    <button
      v-for="period in periods"
      :key="String(period.value)"
      type="button"
      class="flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl text-sm font-medium transition-colors cursor-pointer"
      :class="[
        modelValue === period.value
          ? 'bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-white ring-1 ring-zinc-300 dark:ring-zinc-600'
          : 'text-zinc-500 hover:bg-zinc-100 dark:text-zinc-400 dark:hover:bg-zinc-800',
      ]"
      @click="emit('update:modelValue', period.value)"
    >
      <Icon :name="period.icon" size="18" />
      {{ period.label }}
    </button>
  </div>
</template>
