<script setup lang="ts">
import { staticImage } from '@/utils/str'

interface Props {
  id: number
  title: string
  poster_url?: string
  rating?: number
  type?: string
}

defineProps<Props>()
</script>

<template>
  <NuxtLink :to="`/novela/${id}`" class="group relative flex flex-col gap-2 w-full">
    <div
      class="relative aspect-[2/3] rounded-2xl overflow-hidden bg-zinc-100 dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-800 transition-all duration-500 group-hover:shadow-2xl group-hover:shadow-yellow-500/10 group-hover:-translate-y-1"
    >
      <img
        :src="staticImage(poster_url || '')"
        class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
        loading="lazy"
      />

      <div
        v-if="rating"
        class="absolute top-2 left-2 bg-black/60 backdrop-blur-md text-white text-[10px] font-black px-2 py-0.5 rounded-lg border border-white/10 shadow-sm flex items-center gap-1"
      >
        <Icon name="ph:star-fill" class="text-yellow-500" />
        {{ roundTo(rating, 1) }}
      </div>

      <div
        class="absolute inset-0 rounded-2xl bg-linear-to-t from-black/80 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-end p-3"
      >
        <span
          v-if="type"
          class="text-[10px] text-white font-bold uppercase tracking-widest bg-yellow-500/20 px-2 py-1 rounded-md backdrop-blur-sm"
        >
          {{ type }}
        </span>
      </div>
    </div>

    <div class="px-1">
      <h3
        class="text-sm font-bold text-zinc-900 dark:text-zinc-100 truncate group-hover:text-yellow-500 transition-colors duration-300"
      >
        {{ title }}
      </h3>
    </div>
  </NuxtLink>
</template>
