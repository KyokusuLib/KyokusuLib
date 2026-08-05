<script setup lang="ts">
import type { NovelaDetails } from '@/types/backend/novela'
import { staticImage } from '@/utils/str'
import { roundTo } from '@/utils/num'

interface Props {
  novela: NovelaDetails
}

defineProps<Props>()
</script>

<template>
  <NuxtLink
    :to="`/novela/${novela.id}`"
    class="group flex flex-col gap-3 min-w-35 md:min-w-55 snap-start"
  >
    <div
      class="relative aspect-2/3 rounded-2xl overflow-hidden bg-zinc-100 dark:bg-zinc-800 shadow-sm transition-all duration-500 group-hover:shadow-2xl group-hover:shadow-yellow-500/10 group-hover:-translate-y-1.5 border border-transparent dark:group-hover:border-zinc-700"
    >
      <img
        :src="staticImage(novela.poster_url || '')"
        class="w-full h-full transition-transform duration-700 group-hover:scale-110"
        :alt="novela.title"
      />

      <div class="absolute top-2 left-2 flex flex-col gap-1">
        <div
          class="bg-black/60 backdrop-blur-md text-white text-[10px] font-black px-2 py-0.5 rounded-lg border border-white/10 shadow-sm"
        >
          ★ {{ roundTo(novela.rating_details.total_rating, 1) || '0.0' }}
        </div>
      </div>

      <div
        class="absolute inset-0 bg-linear-to-t from-black/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-end p-3"
      >
        <span class="text-[10px] text-white font-bold uppercase tracking-widest">Читать далее</span>
      </div>
    </div>

    <div class="space-y-1 px-1">
      <h3
        class="text-sm font-bold text-zinc-900 dark:text-zinc-100 truncate group-hover:text-yellow-500 transition-colors duration-300"
      >
        {{ novela.title }}
      </h3>
      <div
        class="flex items-center gap-2 text-[10px] text-zinc-500 dark:text-zinc-500 font-bold uppercase tracking-tighter"
      >
        <span>{{ novela.type }}</span>
        <span class="w-1 h-1 rounded-full bg-zinc-300 dark:bg-zinc-700"></span>
        <span class="truncate">{{ novela.status }}</span>
      </div>
    </div>
  </NuxtLink>
</template>
