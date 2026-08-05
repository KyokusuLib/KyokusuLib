<script setup lang="ts">
import { staticImage } from '@/utils/str'
import { computed } from 'vue'
import type { TotalNovelaStatistics } from '@/types/backend/statistics'
import { roundTo } from '@/utils/num'

const props = defineProps<{
  items: TotalNovelaStatistics[]
}>()

const topThree = computed(() => props.items.slice(0, 3))

const rankConfig = {
  0: {
    label: '1',
    barHeight: 'h-84 sm:h-120',
    cardClass:
      'bg-gradient-to-b from-amber-400/20 to-amber-600/10 dark:from-amber-500/15 dark:to-amber-700/5 border-amber-400/40 dark:border-amber-500/30 ring-amber-400/30 dark:ring-amber-500/20',
    badgeClass: 'bg-amber-500 text-white shadow-lg shadow-amber-500/30',
    icon: 'ph:crown-bold',
    iconColor: 'text-amber-400',
    order: 'order-2',
    scale: 'scale-100',
  },
  1: {
    label: '2',
    barHeight: 'h-66 sm:h-84',
    cardClass:
      'bg-gradient-to-b from-zinc-300/20 to-zinc-400/10 dark:from-zinc-400/15 dark:to-zinc-600/5 border-zinc-300/40 dark:border-zinc-500/30 ring-zinc-300/30 dark:ring-zinc-500/20',
    badgeClass: 'bg-zinc-400 text-white shadow-lg shadow-zinc-400/30',
    icon: 'ph:medal-bold',
    iconColor: 'text-zinc-400 dark:text-zinc-300',
    order: 'order-1',
    scale: 'scale-100',
  },
  2: {
    label: '3',
    barHeight: 'h-54 sm:h-66',
    cardClass:
      'bg-gradient-to-b from-orange-400/20 to-orange-600/10 dark:from-orange-500/15 dark:to-orange-700/5 border-orange-400/40 dark:border-orange-500/30 ring-orange-400/30 dark:ring-orange-500/20',
    badgeClass: 'bg-orange-500 text-white shadow-lg shadow-orange-500/30',
    icon: 'ph:medal-bold',
    iconColor: 'text-orange-400',
    order: 'order-3',
    scale: 'scale-100',
  },
}
</script>

<template>
  <div
    v-if="topThree.length > 0"
    class="flex items-end justify-center gap-4 sm:gap-6 pt-8 pb-4 px-2"
  >
    <div
      v-for="(item, index) in topThree"
      :key="item.novela.novelaId"
      :class="[
        rankConfig[index as keyof typeof rankConfig]?.order,
        rankConfig[index as keyof typeof rankConfig]?.scale,
      ]"
      class="flex flex-col items-center transition-all duration-300"
    >
      <div
        :class="rankConfig[index as keyof typeof rankConfig]?.badgeClass"
        class="relative -mb-6 z-10 w-11 h-11 sm:w-12 sm:h-12 rounded-full flex flex-col items-center justify-center text-xs sm:text-sm font-bold"
      >
        <Icon
          v-if="index === 0"
          :name="rankConfig[index as keyof typeof rankConfig]?.icon"
          class="w-2.5 h-2.5 sm:w-3.5 sm:h-3.5"
        />
        <span class="leading-none">{{ rankConfig[index as keyof typeof rankConfig]?.label }}</span>
      </div>

      <NuxtLink
        :to="`/novela/${item.novela.novelaId}`"
        :class="[
          rankConfig[index as keyof typeof rankConfig]?.cardClass,
          rankConfig[index as keyof typeof rankConfig]?.barHeight,
        ]"
        class="relative flex flex-col items-center w-24 sm:w-32 border ring-1 rounded-2xl sm:rounded-3xl overflow-hidden group hover:ring-2 transition-all duration-300 cursor-pointer pt-8 sm:pt-10 pb-3 sm:pb-4"
      >
        <div
          class="shrink-0 w-14 h-20 sm:w-16 sm:h-24 rounded-lg sm:rounded-xl overflow-hidden shadow-md group-hover:shadow-xl transition-shadow duration-300"
        >
          <img
            :src="staticImage(item.novela.posterUrl)"
            :alt="item.novela.title"
            class="w-full h-full object-cover"
          />
        </div>

        <div class="mt-2 sm:mt-3 flex items-center gap-1">
          <Icon
            name="ph:star-fill"
            :class="rankConfig[index as keyof typeof rankConfig]?.iconColor"
            size="14"
          />
          <span class="text-xs sm:text-sm font-bold text-zinc-800 dark:text-zinc-200">
            {{ roundTo(item.rating, 1) }}
          </span>
        </div>
        <p
          class="mt-1.5 text-[10px] sm:text-xs text-zinc-600 dark:text-zinc-400 font-bold text-center wrap-break-words px-2 max-w-24 sm:max-w-32 leading-tight"
        >
          {{ item.novela.title }}
        </p>
      </NuxtLink>
    </div>
  </div>

  <div v-else class="flex items-end justify-center gap-4 sm:gap-6 pt-8 pb-4 px-2">
    <div class="flex flex-col items-center order-1 scale-100">
      <div
        class="relative -mb-6 z-10 w-11 h-11 sm:w-12 sm:h-12 rounded-full flex flex-col items-center justify-center text-xs sm:text-sm font-bold bg-zinc-400 text-white shadow-lg shadow-zinc-400/30"
      >
        <span class="leading-none">2</span>
      </div>
      <div
        class="relative flex flex-col items-center w-24 sm:w-32 h-52 sm:h-64 bg-linear-to-b from-zinc-300/20 to-zinc-400/10 dark:from-zinc-400/15 dark:to-zinc-600/5 border border-zinc-300/40 dark:border-zinc-500/30 ring-1 ring-zinc-300/30 dark:ring-zinc-500/20 rounded-2xl sm:rounded-3xl pt-8 sm:pt-10 pb-3 sm:pb-4"
      >
        <div
          class="w-14 h-20 sm:w-16 sm:h-24 rounded-lg sm:rounded-xl bg-zinc-200 dark:bg-zinc-700 animate-pulse"
        />
      </div>
      <p
        class="mt-2 sm:mt-3 text-xs sm:text-sm font-semibold text-zinc-500 dark:text-zinc-500 text-center"
      >
        &mdash;
      </p>
    </div>

    <div class="flex flex-col items-center order-2 scale-105">
      <div
        class="relative -mb-6 z-10 w-11 h-11 sm:w-12 sm:h-12 rounded-full flex flex-col items-center justify-center text-xs sm:text-sm font-bold bg-amber-500 text-white shadow-lg shadow-amber-500/30"
      >
        <Icon name="ph:crown-bold" size="14" class="leading-none" />
        <span class="leading-none">1</span>
      </div>
      <div
        class="relative flex flex-col items-center w-24 sm:w-32 h-52 sm:h-64 bg-linear-to-b from-amber-400/20 to-amber-600/10 dark:from-amber-500/15 dark:to-amber-700/5 border border-amber-400/40 dark:border-amber-500/30 ring-1 ring-amber-400/30 dark:ring-amber-500/20 rounded-2xl sm:rounded-3xl pt-8 sm:pt-10 pb-3 sm:pb-4"
      >
        <div
          class="w-14 h-20 sm:w-16 sm:h-24 rounded-lg sm:rounded-xl bg-amber-200 dark:bg-amber-800/40 animate-pulse"
        />
      </div>
      <p
        class="mt-2 sm:mt-3 text-xs sm:text-sm font-semibold text-zinc-500 dark:text-zinc-500 text-center"
      >
        &mdash;
      </p>
    </div>

    <div class="flex flex-col items-center order-3 scale-100">
      <div
        class="relative -mb-6 z-10 w-11 h-11 sm:w-12 sm:h-12 rounded-full flex flex-col items-center justify-center text-xs sm:text-sm font-bold bg-orange-500 text-white shadow-lg shadow-orange-500/30"
      >
        <span class="leading-none">3</span>
      </div>
      <div
        class="relative flex flex-col items-center w-24 sm:w-32 h-52 sm:h-64 bg-linear-to-b from-orange-400/20 to-orange-600/10 dark:from-orange-500/15 dark:to-orange-700/5 border border-orange-400/40 dark:border-orange-500/30 ring-1 ring-orange-400/30 dark:ring-orange-500/20 rounded-2xl sm:rounded-3xl pt-8 sm:pt-10 pb-3 sm:pb-4"
      >
        <div
          class="w-14 h-20 sm:w-16 sm:h-24 rounded-lg sm:rounded-xl bg-orange-200 dark:bg-orange-800/40 animate-pulse"
        />
      </div>
      <p
        class="mt-2 sm:mt-3 text-xs sm:text-sm font-semibold text-zinc-500 dark:text-zinc-500 text-center"
      >
        &mdash;
      </p>
    </div>
  </div>
</template>
