<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    modelValue: number
    min?: number
    max?: number
    step?: number
    label?: string
    suffix?: string
    icon?: string
    format?: (v: number) => string
  }>(),
  {
    min: 0,
    max: 100,
    step: 1,
    suffix: '',
    icon: '',
    format: (v: number) => String(v),
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: number]
}>()

const value = computed({
  get: () => props.modelValue,
  set: (v) => emit('update:modelValue', v),
})

const fillPercent = computed(() => {
  const pct = ((value.value - props.min) / (props.max - props.min)) * 100
  return Math.max(0, Math.min(100, pct))
})
</script>

<template>
  <div
    class="rounded-2xl bg-zinc-50/80 dark:bg-zinc-800/40 border border-zinc-100 dark:border-zinc-800 p-5 space-y-4"
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2.5">
        <div
          v-if="icon"
          class="w-8 h-8 rounded-lg bg-yellow-500/10 flex items-center justify-center"
        >
          <Icon :name="icon" size="16" class="text-yellow-600 dark:text-yellow-400" />
        </div>
        <label
          v-if="label"
          class="text-sm font-bold text-zinc-500 dark:text-zinc-400 uppercase tracking-wider"
        >
          {{ label }}
        </label>
      </div>
      <div class="flex items-baseline gap-0.5">
        <span class="text-lg font-black font-mono text-yellow-500 tabular-nums">{{
          format(value)
        }}</span>
        <span v-if="suffix" class="text-xs font-bold text-zinc-400">{{ suffix }}</span>
      </div>
    </div>

    <div class="px-0.5">
      <div class="relative">
        <input
          v-model.number="value"
          type="range"
          :min="min"
          :max="max"
          :step="step"
          class="range-slider w-full"
          :style="{ '--fill': `${fillPercent}%` }"
        />
      </div>
      <div
        class="flex justify-between mt-2 text-[10px] text-zinc-400 font-bold uppercase tracking-tight select-none"
      >
        <span>{{ min }}</span>
        <span>{{ max }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.range-slider {
  -webkit-appearance: none;
  appearance: none;
  height: 6px;
  cursor: pointer;
  border-radius: 999px;
  outline: none;
  background: transparent;
}

.range-slider::-webkit-slider-runnable-track {
  height: 6px;
  border-radius: 999px;
  background: linear-gradient(to right, #eab308 var(--fill, 0%), #e4e4e7 var(--fill, 0%));
}

:root.dark .range-slider::-webkit-slider-runnable-track {
  background: linear-gradient(to right, #eab308 var(--fill, 0%), #27272a var(--fill, 0%));
}

.range-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  appearance: none;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #eab308;
  cursor: pointer;
  margin-top: -6px;
  box-shadow: 0 2px 8px rgba(234, 179, 8, 0.45);
  transition:
    transform 0.15s ease,
    box-shadow 0.15s ease;
}

.range-slider::-webkit-slider-thumb:hover {
  transform: scale(1.2);
  box-shadow: 0 3px 12px rgba(234, 179, 8, 0.55);
}

.range-slider::-moz-range-track {
  height: 6px;
  border-radius: 999px;
  background: #e4e4e7;
}

:root.dark .range-slider::-moz-range-track {
  background: #27272a;
}

.range-slider::-moz-range-thumb {
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #eab308;
  cursor: pointer;
  border: none;
  box-shadow: 0 2px 8px rgba(234, 179, 8, 0.45);
}

.range-slider::-moz-range-progress {
  height: 6px;
  border-radius: 999px;
  background: #eab308;
}
</style>
