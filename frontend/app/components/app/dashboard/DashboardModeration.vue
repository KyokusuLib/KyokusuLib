<script setup lang="ts">
import { ref, watch } from 'vue'
import { Button, ModalWindow, Spinner, RichText } from '@kyokusu-ui/vue'
import ModerationCard from '@/components/app/dashboard/ModerationCard.vue'
import {
  useModeration,
  type PendingVolume,
  type PendingChapter,
} from '@/composables/api/moderation/useModeration'

const {
  pendingContent,
  isLoading,
  isActionLoading,
  fetchPending,
  approveVolume,
  rejectVolume,
  approveChapter,
  rejectChapter,
  updateVolume,
  updateChapter,
} = useModeration()

onMounted(() => {
  fetchPending()
})

// --- Edit modal state ---
const isEditModalOpen = ref(false)
const editingVolume = ref<PendingVolume | null>(null)
const editingChapter = ref<PendingChapter | null>(null)

const volumeNumber = ref(1)
const volumeTitle = ref('')
const chapterNumber = ref(1)
const chapterTitle = ref('')
const chapterContent = ref('')

const isEditingVolume = computed(() => editingVolume.value !== null)
const modalTitle = computed(() =>
  isEditingVolume.value ? 'Редактирование тома' : 'Редактирование главы',
)

const openVolumeEdit = (volume: PendingVolume) => {
  editingVolume.value = volume
  editingChapter.value = null
  volumeNumber.value = volume.number
  volumeTitle.value = volume.title || ''
  isEditModalOpen.value = true
}

const openChapterEdit = (chapter: PendingChapter) => {
  editingChapter.value = chapter
  editingVolume.value = null
  chapterNumber.value = chapter.number
  chapterTitle.value = chapter.title || ''
  chapterContent.value = chapter.content || ''
  isEditModalOpen.value = true
}

const closeModal = () => {
  isEditModalOpen.value = false
  editingVolume.value = null
  editingChapter.value = null
}

watch(isEditModalOpen, (open) => {
  if (!open) {
    editingVolume.value = null
    editingChapter.value = null
  }
})

const saveVolume = async () => {
  if (!editingVolume.value) return
  const ok = await updateVolume(editingVolume.value.id, {
    volume_number: volumeNumber.value,
    title: volumeTitle.value,
  })
  if (ok) closeModal()
}

const saveChapter = async () => {
  if (!editingChapter.value) return
  const ok = await updateChapter(editingChapter.value.id, {
    chapter_number: chapterNumber.value,
    title: chapterTitle.value,
    content: chapterContent.value,
  })
  if (ok) closeModal()
}

const handleSave = () => {
  if (isEditingVolume.value) {
    saveVolume()
  } else {
    saveChapter()
  }
}
</script>

<template>
  <div class="space-y-10">
    <!-- Loading -->
    <div v-if="isLoading" class="py-24 flex flex-col items-center justify-center text-zinc-400">
      <Spinner size="lg" variant="primary" class="mb-4" />
      <p class="text-xs font-bold uppercase tracking-[0.2em]">Загрузка...</p>
    </div>

    <!-- Empty state -->
    <div
      v-else-if="
        pendingContent &&
        (!pendingContent.volumes || pendingContent.volumes.length === 0) &&
        (!pendingContent.chapters || pendingContent.chapters.length === 0)
      "
      class="py-24 flex flex-col items-center justify-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-[3rem] text-zinc-400 bg-zinc-50/50 dark:bg-zinc-900/10"
    >
      <Icon name="ph:check-circle-bold" size="64" class="opacity-20 mb-4 text-green-500" />
      <p class="font-bold uppercase tracking-[0.2em] text-[10px]">Очередь чиста</p>
    </div>

    <template v-else-if="pendingContent">
      <!-- Volumes queue -->
      <section v-if="pendingContent.volumes && pendingContent.volumes.length > 0">
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

        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
          <ModerationCard
            v-for="volume in pendingContent.volumes"
            :key="volume.id"
            kind="volume"
            :item="volume"
            :loading="isActionLoading"
            @approve="approveVolume(volume.id)"
            @reject="rejectVolume(volume.id)"
            @edit="openVolumeEdit(volume)"
          />
        </div>
      </section>

      <!-- Chapters queue -->
      <section v-if="pendingContent.chapters && pendingContent.chapters.length > 0">
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

        <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 gap-4">
          <ModerationCard
            v-for="chapter in pendingContent.chapters"
            :key="chapter.id"
            kind="chapter"
            :item="chapter"
            :loading="isActionLoading"
            @approve="approveChapter(chapter.id)"
            @reject="rejectChapter(chapter.id)"
            @edit="openChapterEdit(chapter)"
          />
        </div>
      </section>
    </template>

    <!-- Edit modal -->
    <ModalWindow
      v-model="isEditModalOpen"
      :title="modalTitle"
      center-title
      width="w-full max-w-2xl"
    >
      <div v-if="editingVolume" class="flex flex-col gap-4 py-4">
        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-semibold text-zinc-700 dark:text-zinc-300 ml-1">
            Номер тома
          </label>
          <input
            v-model.number="volumeNumber"
            type="number"
            min="1"
            class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700 rounded-xl text-zinc-900 dark:text-zinc-100 outline-none focus:ring-2 focus:ring-yellow-500 transition-all"
            placeholder="Например, 1"
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label class="text-sm font-semibold text-zinc-700 dark:text-zinc-300 ml-1">
            Название (опционально)
          </label>
          <input
            v-model="volumeTitle"
            type="text"
            class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700 rounded-xl text-zinc-900 dark:text-zinc-100 outline-none focus:ring-2 focus:ring-yellow-500 transition-all"
            placeholder="Например, Начало пути"
          />
        </div>
      </div>

      <div v-else-if="editingChapter" class="flex flex-col gap-4 py-4">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="flex flex-col gap-1.5">
            <label class="text-sm font-semibold text-zinc-700 dark:text-zinc-300 ml-1">
              Номер главы
            </label>
            <input
              v-model.number="chapterNumber"
              type="number"
              min="1"
              step="0.1"
              class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700 rounded-xl text-zinc-900 dark:text-zinc-100 outline-none focus:ring-2 focus:ring-yellow-500 transition-all"
              placeholder="Например, 1.5"
            />
          </div>

          <div class="flex flex-col gap-1.5">
            <label class="text-sm font-semibold text-zinc-700 dark:text-zinc-300 ml-1">
              Название (опционально)
            </label>
            <input
              v-model="chapterTitle"
              type="text"
              class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700 rounded-xl text-zinc-900 dark:text-zinc-100 outline-none focus:ring-2 focus:ring-yellow-500 transition-all"
              placeholder="Например, Неожиданная встреча"
            />
          </div>
        </div>

        <div class="flex flex-col gap-1.5">
          <RichText
            id="chapter-content"
            v-model="chapterContent"
            label="Содержание главы"
            placeholder="Введите текст главы..."
          />
        </div>
      </div>

      <template #footer>
        <div class="flex justify-end gap-3">
          <Button variant="outline" @click="closeModal">Отмена</Button>
          <Button
            variant="primary"
            :loading="isActionLoading"
            :disabled="isEditingVolume ? volumeNumber < 1 : chapterNumber < 1"
            @click="handleSave"
          >
            Сохранить
          </Button>
        </div>
      </template>
    </ModalWindow>
  </div>
</template>
