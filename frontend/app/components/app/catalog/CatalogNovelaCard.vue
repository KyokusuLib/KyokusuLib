<script setup lang="ts">
import type { NovelaDetails } from '@/types/backend/novela'
import { staticImage } from '@/utils/str'
import { roundTo } from '@/utils/num'

import {
  NOVELA_CATALOG_CARD_TITLE_MAX_LENGTH,
  NOVELA_CATALOG_CARD_TITLE_DEVIDE,
} from '@/constants/data'

import { Card } from '@kyokusu-ui/vue'

const props = defineProps<{
  novela: NovelaDetails
}>()

const novelaImg = computed(() => {
  return staticImage(props.novela.poster_url)
})

const novelaTitle = computed(() => {
  return props.novela.title.length > NOVELA_CATALOG_CARD_TITLE_MAX_LENGTH &&
    props.novela.title.length > NOVELA_CATALOG_CARD_TITLE_DEVIDE
    ? props.novela.title.slice(0, NOVELA_CATALOG_CARD_TITLE_DEVIDE) + '...'
    : props.novela.title
})

const novelaReleaseYear = computed(() => {
  return props.novela.release_date ? new Date(props.novela.release_date).getFullYear() : 'N/A'
})
</script>

<template>
  <Card
    class="cursor-pointer group h-full w-fit md:w-full"
    variant="outline"
    padding="md"
    shadow
    v-on:click="() => navigateTo(`/novela/${props.novela.id}`)"
  >
    <div class="flex group/card transition-opacity duration-300">
      <div
        class="relative h-52.5 md:h-64 w-44 md:w-full rounded-xl overflow-hidden shadow-lg dark:shadow-2xl border border-zinc-200 dark:border-zinc-800 transition-all duration-300 group-hover/card:scale-[1.02] group-hover/card:shadow-xl dark:group-hover/card:shadow-[0_0_30px_rgba(255,255,255,0.1)] bg-zinc-200 dark:bg-zinc-900 cursor-pointer"
      >
        <img
          :src="novelaImg"
          class="h-64 md:h-64 w-full rounded-xl object-cover pointer-events-none"
          alt="cover"
        />
        <div
          class="absolute top-2 right-2 opacity-80 bg-zinc-900 text-amber-500 font-bold text-[10px] md:text-xs px-2 py-1 rounded-2xl backdrop-blur-sm shadow-sm"
        >
          ★ {{ roundTo(props.novela.rating_details?.total_rating, 1) }}
        </div>
      </div>
    </div>
    <div class="mt-1 ml-1 max-w-44 md:max-w-full">
      <h1
        class="font-bold text-zinc-900 dark:text-white text-[12px] group-hover:text-yellow-500 transition-colors duration-300 truncate md:text-wrap md:wrap-break-words md:hyphens-auto leading-snug"
        lang="ru"
      >
        {{ novelaTitle }}
      </h1>
    </div>

    <div
      class="grid grid-cols-[auto_auto_1fr] gap-x-2 mt-1 ml-1 text-[10px] text-zinc-500 dark:text-zinc-500 font-bold uppercase tracking-tighter"
    >
      <span>{{ novela.type }}</span>
      <span class="w-1 h-1 mt-1.5 rounded-full bg-zinc-300 dark:bg-zinc-700"></span>
      <span class="truncate">{{ novela.status }}</span>

      <span>{{ novelaReleaseYear }}</span>
      <span class="w-1 h-1 mt-1.5 rounded-full bg-zinc-300 dark:bg-zinc-700"></span>
      <span class="truncate">{{ novela.country }}</span>
    </div>
  </Card>
</template>
