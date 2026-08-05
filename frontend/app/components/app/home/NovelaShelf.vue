<script setup lang="ts">
import NovelaCard from './NovelaCard.vue'
import type { NovelaDetails } from '~/types/backend/novela'

interface Props {
  title: string
  titleIcon: string
  items: NovelaDetails[]
  viewAllHref?: string
}

defineProps<Props>()
</script>

<template>
  <section class="flex flex-col gap-6 w-full py-2">
    <div class="flex items-end justify-between px-2">
      <div class="space-y-1">
        <h2
          class="flex items-start gap-3 text-xl md:text-2xl font-black tracking-tight text-zinc-900 dark:text-white uppercase italic"
        >
          <div class="flex flex-col items-center gap-3">
            {{ title }}
            <div class="h-1 w-full bg-yellow-500 rounded-full"></div>
          </div>
          <Icon :name="titleIcon" size="20" class="text-zinc-700 dark:text-zinc-200 mt-2" />
        </h2>
      </div>

      <NuxtLink
        v-if="viewAllHref"
        :to="viewAllHref"
        class="group flex items-center gap-2 text-[10px] font-black text-zinc-400 hover:text-yellow-500 transition-colors uppercase tracking-[0.2em]"
      >
        Все сразу
        <svg
          class="w-3 h-3 transition-transform group-hover:translate-x-1"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M9 5l7 7-7 7" />
        </svg>
      </NuxtLink>
    </div>

    <div class="relative group/shelf">
      <div
        class="flex max-h-[400px] gap-5 overflow-x-auto pb-6 px-2 snap-x scroll-smooth no-scrollbar"
      >
        <NovelaCard v-for="novela in items" :key="novela.id" :novela="novela" />
      </div>

      <div
        class="absolute top-0 right-0 h-full w-20 bg-linear-to-l from-white dark:from-[#0f0f0f] to-transparent pointer-events-none opacity-0 group-hover/shelf:opacity-100 transition-opacity duration-500"
      ></div>
    </div>
  </section>
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
