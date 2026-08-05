<script setup lang="ts">
import type { NovelaChapter } from '~/types/backend/novela'
import ChapterListItem from './ChapterListItem.vue'

defineProps<{
  chapters: NovelaChapter[]
  novelaId: number
}>()
</script>

<template>
  <div
    class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl overflow-hidden max-h-150 overflow-y-auto custom-scrollbar"
  >
    <ChapterListItem
      v-for="chapter in chapters"
      :key="chapter.id"
      :chapter="chapter"
      :novela-id="novelaId"
    />

    <div v-if="!chapters.length" class="p-8 text-center text-zinc-500">Главы не найдены</div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #d4d4d8;
  border-radius: 20px;
}
:deep(.dark) .custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #3f3f46;
}
</style>
