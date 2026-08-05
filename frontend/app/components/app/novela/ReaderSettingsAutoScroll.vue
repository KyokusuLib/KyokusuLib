<script setup lang="ts">
import { useReaderSettings } from '~/composables/ui/useReaderSettings'
import { Toggle } from '@kyokusu-ui/vue'

const { isAutoScrollEnabled } = useReaderSettings()
</script>

<template>
  <div
    class="rounded-2xl bg-zinc-50/80 dark:bg-zinc-800/40 border border-zinc-100 dark:border-zinc-800 p-5 space-y-4"
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2.5">
        <div class="w-8 h-8 rounded-lg bg-yellow-500/10 flex items-center justify-center">
          <Icon name="ph:infinity-bold" size="16" class="text-yellow-600 dark:text-yellow-400" />
        </div>
        <label class="text-sm font-bold text-zinc-500 dark:text-zinc-400 uppercase tracking-wider"
          >Автопрокрутка</label
        >
      </div>

      <button
        @click="isAutoScrollEnabled = !isAutoScrollEnabled"
        class="relative w-11 h-6 rounded-full transition-colors duration-300 shrink-0 cursor-pointer"
        role="switch"
        :class="isAutoScrollEnabled ? 'bg-yellow-500' : 'bg-zinc-200 dark:bg-zinc-700'"
        :aria-checked="isAutoScrollEnabled"
      >
        <div
          class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow-md transition-all duration-300"
          :class="isAutoScrollEnabled ? 'translate-x-5' : 'translate-x-0'"
        ></div>
      </button>
    </div>

    <p class="text-xs text-zinc-500 dark:text-zinc-400 leading-relaxed">
      Следующая глава будет загружаться автоматически при достижении конца текущей.
    </p>
  </div>
</template>
