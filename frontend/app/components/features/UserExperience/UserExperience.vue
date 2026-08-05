<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  currentExp: number
  expToNextLevel: number
  level: number
  levelTitle: string
}

const props = defineProps<Props>()

const progress = computed(() => {
  if (props.expToNextLevel === 0) return 0
  const LEVEL_MAX_EXP = props.expToNextLevel + props.currentExp

  return LEVEL_MAX_EXP > 0 ? (props.currentExp / LEVEL_MAX_EXP) * 100 : 0
})
</script>

<template>
  <div
    class="group flex items-center mt-1 gap-3 p-1 pl-1.5 pr-4 rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900/50 shadow-sm transition-all duration-300 hover:border-yellow-500/30 hover:shadow-lg hover:shadow-yellow-500/5 cursor-default h-10 w-full sm:w-auto min-w-[280px]"
  >
    <div
      class="relative shrink-0 flex items-center justify-center w-6 h-6 rounded-xl bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 shadow-md"
    >
      <span class="text-[15px] font-black italic">{{ level }}</span>
      <div
        class="absolute inset-0 bg-linear-to-tr from-white/10 to-transparent rounded-xl pointer-events-none"
      ></div>
    </div>

    <div class="flex-1 flex flex-col gap-1 min-w-0">
      <div class="flex justify-between items-end leading-none">
        <span
          class="text-[9px] font-black uppercase tracking-[0.2em] ml-px text-zinc-400 group-hover:text-yellow-600 transition-colors"
        >
          Уровень: {{ levelTitle }}
        </span>
        <span class="text-[9px] font-bold text-zinc-500 dark:text-zinc-400 tabular-nums">
          {{ Math.floor(progress) }}%
        </span>
      </div>

      <div class="h-1.5 bg-zinc-100 dark:bg-zinc-800 rounded-full overflow-hidden relative">
        <div
          class="absolute inset-y-0 left-0 bg-linear-to-r from-yellow-400 via-yellow-500 to-yellow-600 rounded-full transition-all duration-1000 ease-out"
          :style="{ width: `${progress}%` }"
        >
          <div class="absolute right-0 top-0 bottom-0 w-4 bg-white/30 blur-sm"></div>
        </div>
      </div>
    </div>

    <div
      class="hidden lg:flex flex-col items-end leading-none ml-2"
      :title="`До следующего уровня: ${expToNextLevel} XP`"
    >
      <span class="text-[10px] font-black text-zinc-900 dark:text-zinc-200">
        {{ currentExp }}
      </span>
      <span class="text-[8px] font-bold text-zinc-400 uppercase tracking-tighter"> XP </span>
    </div>
  </div>
</template>

<style scoped>
.group:hover .bg-white\/30 {
  animation: shine 1.5s infinite;
}

@keyframes shine {
  0% {
    opacity: 0.2;
    transform: translateX(-100%);
  }
  50% {
    opacity: 0.5;
  }
  100% {
    opacity: 0.2;
    transform: translateX(100%);
  }
}
</style>
