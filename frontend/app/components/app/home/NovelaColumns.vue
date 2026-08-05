<script setup lang="ts">
import NovelaRowItem from './NovelaRowItem.vue'
import type { NovelaDetails } from '~/types/backend/novela'

interface Props {
  newItems: NovelaDetails[]
  popularItems: NovelaDetails[]
}
defineProps<Props>()
</script>

<template>
  <div
    class="w-full flex flex-col lg:flex-row items-stretch gap-10 lg:gap-8 px-3 py-6 sm:px-6 sm:py-8 bg-white dark:bg-[#0f0f0f] border border-zinc-200 dark:border-zinc-800 rounded-3xl sm:rounded-[2.5rem] shadow-sm overflow-hidden"
  >
    <div class="flex flex-col w-full lg:w-1/2 min-w-0">
      <div class="flex items-center justify-between mb-5 px-1">
        <h3
          class="text-sm sm:text-base md:text-lg font-black uppercase tracking-widest italic border-l-4 border-yellow-500 pl-3 text-zinc-900 dark:text-white truncate"
        >
          Новинки
        </h3>
        <NuxtLink
          to="/catalog?sort=new"
          class="shrink-0 text-[10px] font-bold uppercase text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200 transition-colors ml-4 whitespace-nowrap"
        >
          Все сразу
        </NuxtLink>
      </div>

      <div class="flex flex-col gap-2 w-full">
        <NovelaRowItem
          v-for="(item, index) in newItems"
          :key="item.id"
          :novela="item"
          :index="index"
        />
      </div>
    </div>

    <div
      class="shrink-0 h-px lg:h-auto lg:w-px bg-linear-to-r lg:bg-linear-to-b from-transparent via-zinc-200 dark:via-zinc-800 to-transparent my-2 lg:my-0"
    ></div>

    <div class="flex flex-col w-full lg:w-1/2 min-w-0">
      <div class="flex items-center justify-between mb-5 px-1">
        <h3
          class="text-sm sm:text-base md:text-lg font-black uppercase tracking-widest italic border-l-4 border-yellow-500 pl-3 text-zinc-900 dark:text-white truncate"
        >
          Популярное
        </h3>
        <NuxtLink
          to="/catalog?sort=popular"
          class="shrink-0 text-[10px] font-bold uppercase text-zinc-400 hover:text-zinc-900 dark:hover:text-zinc-200 transition-colors ml-4 whitespace-nowrap"
        >
          Все сразу
        </NuxtLink>
      </div>

      <div class="flex flex-col gap-2 w-full">
        <NovelaRowItem
          v-for="(item, index) in popularItems"
          :key="item.id"
          :novela="item"
          :index="index"
        />
      </div>
    </div>
  </div>
</template>
