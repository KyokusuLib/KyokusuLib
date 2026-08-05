<script setup lang="ts">
import { useNovela } from '@/composables/api/novela/useNovela'
import type { NovelaDetails } from '@/types/backend/novela'
import NovelaCard from '@/components/app/home/NovelaCard.vue'

const props = defineProps<{
  authorId: number
  authorName: string
}>()

const { fetchNovels } = useNovela()

const { data: novels, pending: isLoading } = useAsyncData<NovelaDetails[]>(
  `author-novels-${props.authorId}`,
  async () => {
    try {
      const fetchedNovels = await fetchNovels({ author_id: props.authorId, limit: 12 })
      return fetchedNovels || []
    } catch (e) {
      console.error('Failed to fetch author novels', e)
      return []
    }
  },
  { default: () => [] },
)
</script>

<template>
  <div
    v-if="isLoading || novels.length > 0"
    class="bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 sm:p-8 shadow-sm w-full"
  >
    <div class="flex items-end justify-between mb-6">
      <h2
        class="flex items-center gap-3 text-xl md:text-2xl font-black tracking-tight text-zinc-900 dark:text-white uppercase italic"
      >
        <div class="flex flex-col items-start gap-1">
          Работы автора
          <div class="h-1 w-1/2 bg-yellow-500 rounded-full"></div>
        </div>
      </h2>
    </div>

    <div class="mt-2">
      <div
        v-if="isLoading"
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 sm:gap-6 pb-4"
      >
        <div v-for="i in 5" :key="i" class="w-full">
          <div
            class="aspect-2/3 rounded-2xl bg-zinc-100 dark:bg-zinc-800/50 animate-pulse mb-3"
          ></div>
          <div class="h-4 bg-zinc-100 dark:bg-zinc-800/50 rounded w-3/4 animate-pulse"></div>
        </div>
      </div>

      <div
        v-else
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 sm:gap-6 pb-4"
      >
        <div v-for="novela in novels" :key="novela.id" class="w-full h-full">
          <NovelaCard :novela="novela" />
        </div>
      </div>
    </div>
  </div>
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
