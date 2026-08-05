<script setup lang="ts">
import HomeHero from './HomeHero/HomeHero.vue'
import { Carousel } from '@kyokusu-ui/vue'
import NovelaShelf from './NovelaShelf.vue'
import NovelaColumns from './NovelaColumns.vue'
import GeneralStats from './GeneralStats.vue'
import { useNovela } from '@/composables/api/novela/useNovela'
import {
  NOVELA_FETCH_LIMIT,
  NOVELA_FETCH_LIMIT_LAST_UPDATED,
  NOVELA_FETCH_LIMIT_NEW,
  NOVELA_FETCH_LIMIT_POPULAR,
  NOVELA_FETCH_LIMIT_TRENDING,
} from '@/constants/data'
import { NovelaSort } from '@/types/frontend/query/novela-query'
import { useGeneralStatistics } from '@/composables/api/statistics/general-stats'

const { fetchNovels } = useNovela()
const { fetchGeneralStatistics } = useGeneralStatistics()

const { data: novels } = await useAsyncData('novelas-home', () =>
  fetchNovels({ limit: NOVELA_FETCH_LIMIT }),
)

const { data: newNovels } = await useAsyncData('novels-new', () =>
  fetchNovels({ sort: NovelaSort.New, limit: NOVELA_FETCH_LIMIT_NEW }),
)

const { data: popularNovels } = await useAsyncData('novels-popular', () =>
  fetchNovels({ sort: NovelaSort.Popular, limit: NOVELA_FETCH_LIMIT_POPULAR }),
)

const { data: updatedNovels } = await useAsyncData('novels-updated', () =>
  fetchNovels({
    sort: NovelaSort.Updated,
    limit: NOVELA_FETCH_LIMIT_LAST_UPDATED,
  }),
)

const { data: trendingNovels } = await useAsyncData('novels-trending', () =>
  fetchNovels({
    sort: NovelaSort.Trending,
    limit: NOVELA_FETCH_LIMIT_TRENDING,
  }),
)

const { data: generalStats, pending: loadingStats } = await useAsyncData('general-stats', () =>
  fetchGeneralStatistics(),
)
</script>

<template>
  <div
    class="min-h-screen flex flex-col bg-zinc-50 dark:bg-radial-[at_center] dark:from-zinc-800 dark:to-zinc-950 dark:to-90% transition-colors duration-300"
  >
    <HomeHero />
    <div class="w-full max-w-6xl mx-auto px-2 sm:px-4">
      <Carousel v-if="!!novels" :items="novels">
        <template #card="{ item }">
          <div
            class="flex flex-col gap-2 md:gap-3 w-35 md:w-55 group/card transition-opacity duration-300"
          >
            <div
              class="relative h-52.5 md:h-80 w-full rounded-xl overflow-hidden shadow-lg dark:shadow-2xl border border-zinc-200 dark:border-zinc-800 transition-transform duration-300 group-hover/card:scale-[1.02] bg-zinc-200 dark:bg-zinc-900 cursor-pointer"
              @click="() => navigateTo(`/novela/${item.id}`)"
            >
              <img
                :src="staticImage(item.poster_url)"
                class="h-full w-full object-cover pointer-events-none"
                alt="cover"
              />
              <div
                class="absolute top-2 left-2 bg-white/90 text-black font-bold text-[10px] md:text-xs px-2 py-1 rounded backdrop-blur-sm shadow-sm"
              >
                ★ {{ roundTo(item.rating_details.total_rating, 1) }}
              </div>
            </div>

            <div
              class="flex flex-col items-center px-1"
              @click="() => navigateTo(`/novela/${item.id}`)"
            >
              <div class="flex flex-col gap-2 items-center h-17.5">
                <h3
                  class="cs-text text-center leading-md line-clamp-2 text-black dark:text-white text-[15px] leading-[1.15] font-semibold text-balance dark:font-semibold group-hover/card:text-yellow-600 dark:group-hover/card:text-yellow-400 transition-colors"
                >
                  {{ item.title }}
                </h3>
                <span class="text-[10px] md:text-sm text-zinc-500 dark:text-zinc-400">
                  {{ item.type }}
                </span>
              </div>

              <button
                class="mt-2 w-full py-1.5 md:py-2 cursor-pointer rounded-lg bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-900 dark:text-white text-xs md:text-sm font-medium transition-colors border border-zinc-200 dark:border-zinc-700"
              >
                Читать
              </button>
            </div>
          </div>
        </template>
      </Carousel>
    </div>

    <div class="w-full max-w-6xl mx-auto px-2 sm:px-4 py-6">
      <NovelaColumns :new-items="newNovels || []" :popular-items="popularNovels || []" />
    </div>

    <div class="w-full max-w-6xl mx-auto px-2 sm:px-4 py-6">
      <NovelaShelf title="Трендовое" titleIcon="ph:books-bold" :items="trendingNovels || []" />

      <NovelaShelf
        title="Свежие главы"
        titleIcon="ph:book-open-bold"
        :items="updatedNovels || []"
      />
    </div>

    <div class="w-full max-w-6xl mx-auto px-2 sm:px-4 py-6">
      <GeneralStats v-if="generalStats != null" :stats="generalStats" :is-loading="loadingStats" />
    </div>
  </div>
</template>
