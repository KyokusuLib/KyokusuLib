<script setup lang="ts">
import type { NovelaDetails } from '~/types/backend/novela'
import { staticImage } from '@/utils/str'
import { roundTo } from '@/utils/num'

interface Props {
  novela: NovelaDetails
  index: number
}
defineProps<Props>()
</script>

<template>
  <NuxtLink
    :to="`/novela/${novela.id}`"
    class="group flex items-center gap-4 p-2 rounded-2xl transition-all duration-300 hover:bg-zinc-100 dark:hover:bg-zinc-800/40"
  >
    <div
      class="relative shrink-0 w-12 h-16 md:w-14 md:h-20 rounded-xl overflow-hidden bg-zinc-200 dark:bg-zinc-800 shadow-sm"
    >
      <img
        :src="staticImage(novela.poster_url || '')"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
      />
      <div class="absolute inset-0 bg-black/10 group-hover:bg-transparent transition-colors"></div>
    </div>

    <div class="flex flex-col gap-1 min-w-0 overflow-hidden">
      <h4
        class="text-sm font-black text-zinc-900 dark:text-zinc-100 truncate group-hover:text-yellow-500 transition-colors"
      >
        {{ novela.title }}
      </h4>
      <div class="flex items-center gap-2 text-[10px] font-bold uppercase tracking-tighter italic">
        <span class="text-yellow-500 flex items-center gap-0.5">
          ★ {{ roundTo(novela.rating_details.total_rating, 1) || '0.0' }}
        </span>
        <span class="text-zinc-400 dark:text-zinc-500">{{ novela.type }}</span>
      </div>
    </div>

    <div
      class="ml-auto text-xl font-black italic text-zinc-100 dark:text-zinc-800/50 group-hover:text-yellow-500/20 transition-colors"
    >
      #{{ index + 1 }}
    </div>
  </NuxtLink>
</template>
