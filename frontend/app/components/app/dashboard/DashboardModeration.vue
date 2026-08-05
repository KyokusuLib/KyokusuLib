<script setup lang="ts">
import { useModeration } from '@/composables/api/moderation/useModeration'
import { onMounted } from 'vue'

const {
  pendingContent,
  isLoading,
  isActionLoading,
  fetchPending,
  approveVolume,
  rejectVolume,
  approveChapter,
  rejectChapter,
} = useModeration()

onMounted(() => {
  fetchPending()
})

const formatChapterNumber = (num: number) => {
  return Number.isInteger(num) ? num : num.toFixed(1)
}
</script>

<template>
  <div class="space-y-8">
    <div v-if="isLoading" class="py-24 flex flex-col items-center justify-center text-zinc-400">
      <Icon name="ph:spinner-bold" size="32" class="animate-spin mb-4" />
      <p class="text-xs font-bold uppercase tracking-[0.2em]">Загрузка...</p>
    </div>

    <div v-else-if="pendingContent">
      <div
        v-if="
          (!pendingContent.volumes || pendingContent.volumes.length === 0) &&
          (!pendingContent.chapters || pendingContent.chapters.length === 0)
        "
        class="py-24 flex flex-col items-center justify-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-[3rem] text-zinc-400 bg-zinc-50/50 dark:bg-zinc-900/10"
      >
        <Icon name="ph:check-circle-bold" size="64" class="opacity-20 mb-4 text-green-500" />
        <p class="font-bold uppercase tracking-[0.2em] text-[10px]">Очередь чиста</p>
      </div>

      <div v-else class="grid grid-cols-1 xl:grid-cols-2 gap-8">
        <div v-if="pendingContent.volumes && pendingContent.volumes.length > 0" class="space-y-4">
          <div class="flex items-center gap-2 mb-6">
            <div
              class="w-8 h-8 rounded-xl bg-yellow-500/10 text-yellow-500 flex items-center justify-center"
            >
              <Icon name="ph:books-bold" size="18" />
            </div>
            <h2 class="text-xl font-bold text-zinc-900 dark:text-white">
              Тома
              <span class="text-sm text-zinc-500 font-medium"
                >({{ pendingContent.volumes.length }})</span
              >
            </h2>
          </div>

          <div
            v-for="volume in pendingContent.volumes"
            :key="volume.id"
            class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 shadow-sm transition-all hover:border-yellow-500/30 group"
          >
            <div class="flex justify-between items-start mb-4">
              <div>
                <h3 class="font-bold text-zinc-900 dark:text-white">
                  Том {{ volume.number }}
                  <span v-if="volume.title" class="text-zinc-500 ml-1 font-medium"
                    >— {{ volume.title }}</span
                  >
                </h3>
                <p class="text-xs text-zinc-400 mt-1">
                  Добавил пользователь: #{{ volume.created_by }}
                </p>
              </div>
              <div
                class="bg-yellow-500/10 text-yellow-600 dark:text-yellow-400 px-2 py-1 rounded text-[10px] font-bold uppercase"
              >
                Ожидает
              </div>
            </div>

            <div class="flex gap-2">
              <button
                @click="approveVolume(volume.id)"
                :disabled="isActionLoading"
                class="flex-1 py-2 bg-green-50 hover:bg-green-100 dark:bg-green-500/10 dark:hover:bg-green-500/20 text-green-600 dark:text-green-400 font-semibold text-sm rounded-xl transition-colors flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
              >
                <Icon name="ph:check-bold" />
                Одобрить
              </button>
              <button
                @click="rejectVolume(volume.id)"
                :disabled="isActionLoading"
                class="flex-1 py-2 bg-red-50 hover:bg-red-100 dark:bg-red-500/10 dark:hover:bg-red-500/20 text-red-600 dark:text-red-400 font-semibold text-sm rounded-xl transition-colors flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
              >
                <Icon name="ph:x-bold" />
                Отклонить
              </button>
            </div>
          </div>
        </div>

        <!-- Chapters Queue -->
        <div v-if="pendingContent.chapters && pendingContent.chapters.length > 0" class="space-y-4">
          <div class="flex items-center gap-2 mb-6">
            <div
              class="w-8 h-8 rounded-xl bg-blue-500/10 text-blue-500 flex items-center justify-center"
            >
              <Icon name="ph:file-text-bold" size="18" />
            </div>
            <h2 class="text-xl font-bold text-zinc-900 dark:text-white">
              Главы
              <span class="text-sm text-zinc-500 font-medium"
                >({{ pendingContent.chapters.length }})</span
              >
            </h2>
          </div>

          <div
            v-for="chapter in pendingContent.chapters"
            :key="chapter.id"
            class="bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-5 shadow-sm transition-all hover:border-blue-500/30 group"
          >
            <div class="flex justify-between items-start mb-4">
              <div>
                <h3 class="font-bold text-zinc-900 dark:text-white">
                  Глава {{ formatChapterNumber(chapter.number) }}
                  <span v-if="chapter.title" class="text-zinc-500 ml-1 font-medium"
                    >— {{ chapter.title }}</span
                  >
                </h3>
                <p class="text-xs text-zinc-400 mt-1">
                  Добавил пользователь: #{{ chapter.created_by }}
                </p>
              </div>
              <div
                class="bg-yellow-500/10 text-yellow-600 dark:text-yellow-400 px-2 py-1 rounded text-[10px] font-bold uppercase"
              >
                Ожидает
              </div>
            </div>

            <div class="mb-4">
              <p
                class="text-xs text-zinc-500 dark:text-zinc-400 bg-zinc-50 dark:bg-zinc-800/50 p-3 rounded-xl max-h-24 overflow-y-auto line-clamp-3"
              >
                {{ chapter.content || 'Без текста (могут быть изображения)' }}
              </p>
            </div>

            <div class="flex gap-2">
              <button
                @click="approveChapter(chapter.id)"
                :disabled="isActionLoading"
                class="flex-1 py-2 bg-green-50 hover:bg-green-100 dark:bg-green-500/10 dark:hover:bg-green-500/20 text-green-600 dark:text-green-400 font-semibold text-sm rounded-xl transition-colors flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
              >
                <Icon name="ph:check-bold" />
                Одобрить
              </button>
              <button
                @click="rejectChapter(chapter.id)"
                :disabled="isActionLoading"
                class="flex-1 py-2 bg-red-50 hover:bg-red-100 dark:bg-red-500/10 dark:hover:bg-red-500/20 text-red-600 dark:text-red-400 font-semibold text-sm rounded-xl transition-colors flex items-center justify-center gap-2 cursor-pointer disabled:opacity-50"
              >
                <Icon name="ph:x-bold" />
                Отклонить
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
