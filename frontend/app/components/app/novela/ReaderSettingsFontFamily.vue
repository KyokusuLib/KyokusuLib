<script setup lang="ts">
import { READER_FONTS, useReaderSettings } from '~/composables/ui/useReaderSettings'

const { fontFamily } = useReaderSettings()
</script>

<template>
  <div
    class="rounded-2xl bg-zinc-50/80 dark:bg-zinc-800/40 border border-zinc-100 dark:border-zinc-800 p-5 space-y-4"
  >
    <div class="flex items-center gap-2.5">
      <div class="w-8 h-8 rounded-lg bg-yellow-500/10 flex items-center justify-center">
        <Icon name="ph:text-aa-bold" size="16" class="text-yellow-600 dark:text-yellow-400" />
      </div>
      <label class="text-sm font-bold text-zinc-500 dark:text-zinc-400 uppercase tracking-wider"
        >Шрифт</label
      >
    </div>

    <div class="grid grid-cols-1 gap-2">
      <button
        v-for="font in READER_FONTS"
        :key="font.value"
        @click="fontFamily = font.value"
        class="group flex items-center justify-between px-4 py-3 rounded-xl border transition-all duration-200 cursor-pointer"
        :class="
          fontFamily === font.value
            ? 'border-yellow-500 bg-yellow-500/5 text-yellow-600 dark:text-yellow-400 shadow-sm shadow-yellow-500/10'
            : 'border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 hover:border-zinc-300 dark:hover:border-zinc-600 hover:bg-zinc-50 dark:hover:bg-zinc-700/50'
        "
      >
        <span :class="[font.value, 'text-sm font-bold']">{{ font.name }}</span>
        <div
          class="w-5 h-5 rounded-full flex items-center justify-center transition-all duration-200"
          :class="
            fontFamily === font.value
              ? 'bg-yellow-500 text-white scale-100'
              : 'bg-zinc-200 dark:bg-zinc-700 text-transparent scale-90 opacity-0 group-hover:opacity-50'
          "
        >
          <Icon v-if="fontFamily === font.value" name="ph:check-bold" size="12" />
        </div>
      </button>
    </div>
  </div>
</template>
