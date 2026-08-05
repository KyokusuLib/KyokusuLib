<script setup lang="ts">
import { useNovela } from '~/composables/api/novela/useNovela'
import { useChapterDraftStore } from '~/stores/chapterDraft'
import { staticImage } from '~/utils/str'
import { useChapterEditor } from '~/composables/ui/useChapterEditor'
import { useImageUpload } from '~/composables/ui/useImageUpload'
import { useChapterSubmit } from '~/composables/ui/useChapterSubmit'
import { getLastChapterNumberForVolume } from '~/utils/chapter'
import type { NovelaChapter } from '~/types/backend/novela'
import { Select as BaseSelect } from '@kyokusu-ui/vue'
import ChapterEditorBlock from '~/components/app/novela/ChapterEditorBlock.vue'
import ImagePreviewModal from '~/components/app/novela/ImagePreviewModal.vue'
import ChapterEditorImageBlock from '~/components/app/novela/ChapterEditorImageBlock.vue'

const route = useRoute()
const router = useRouter()
const novelaId = route.params.id as string

const { novela, fetchNovela } = useNovela()
const draftStore = useChapterDraftStore()
const existingDraft = draftStore.getDraft(novelaId)

await useAsyncData(`novela-add-chapter-${novelaId}`, () => fetchNovela(novelaId))

const volumes = computed(() => novela.value?.volumes ?? [])

const volumeOptions = computed(() =>
  volumes.value.map((v) => ({
    value: v.id,
    label: `Том ${v.number}${v.title ? ` — ${v.title}` : ''}`,
  })),
)

const selectedVolumeId = ref(existingDraft?.selectedVolumeId ?? '')
const chapterNumber = ref(existingDraft === null ? 1 : existingDraft.chapterNumber + 1)
const chapterTitle = ref(existingDraft?.chapterTitle ?? '')

const blocks = ref(
  existingDraft?.blocks ?? [{ id: `${Date.now()}-block-0`, type: 'text' as const, content: '' }],
)

const { addTextBlock, addImageBlock, removeBlock, reorderBlock } = useChapterEditor(blocks)

const { images, addFile, removeImageById, getImageById, updateCaption } = useImageUpload()
const { isSubmitting, submit } = useChapterSubmit()

const imageBlockOrder = computed(() => {
  let counter = 0
  return blocks.value.map((block) => {
    if (block.type === 'image') {
      counter++
      return counter
    }
    return 0
  })
})

const dragIndex = ref<number | null>(null)
const previewImageUrl = ref<string | null>(null)

function onDragStart(index: number) {
  dragIndex.value = index
}

function onDragOver(e: DragEvent, index: number) {
  e.preventDefault()
  if (dragIndex.value === null || dragIndex.value === index) return
  reorderBlock(dragIndex.value, index)
  dragIndex.value = index
}

function onDragEnd() {
  dragIndex.value = null
}

function addTextAt(index: number) {
  addTextBlock(index)
}

function addImageAt(index: number) {
  addImageBlock(index)
}

async function onImageSelected(index: number, file: File) {
  const block = blocks.value[index]
  if (!block || block.type !== 'image') return

  if (block.content) {
    removeImageById(block.content)
  }

  const image = await addFile(file)
  block.content = image.id
}

function onImageCaptionUpdate(index: number, caption: string) {
  const block = blocks.value[index]
  if (!block || block.type !== 'image' || !block.content) return
  updateCaption(block.content, caption)
}

function onRemoveBlock(index: number) {
  const block = blocks.value[index]
  if (block?.type === 'image' && block.content) {
    removeImageById(block.content)
  }
  removeBlock(index)
}

function previewImage(url: string) {
  previewImageUrl.value = url
}

function closePreview() {
  previewImageUrl.value = null
}

const chapterNumberError = ref('')

function onChapterNumberInput(e: Event) {
  const input = e.target as HTMLInputElement
  const val = input.value
  const num = Number(val)
  if (val === '' || (Number.isInteger(num) && num >= 1)) {
    chapterNumber.value = num || 1
    chapterNumberError.value = ''
  } else {
    chapterNumberError.value = 'Номер главы должен быть целым числом'
  }
}

async function onSubmit() {
  if (!Number.isInteger(chapterNumber.value) || chapterNumber.value < 1) {
    chapterNumberError.value = 'Номер главы должен быть целым числом'
    return
  }
  const res = await submit(
    novelaId,
    selectedVolumeId.value,
    chapterNumber.value,
    chapterTitle.value,
    blocks.value,
    images.value,
  )
  if (res) {
    draftStore.clearDraft(novelaId)
    if (novela.value) {
      const volume = novela.value.volumes.find((v) => v.id === selectedVolumeId.value)
      if (volume) {
        const newChapter: NovelaChapter = {
          id: res.id,
          title: chapterTitle.value || `Глава ${chapterNumber.value}`,
          number: chapterNumber.value,
          status: res.status,
          content: '',
          images: [],
        }
        volume.chapters.push(newChapter)
      }
    }
    router.push(`/novela/${novelaId}`)
  }
}

function persistDraft() {
  draftStore.updateDraft(novelaId, {
    selectedVolumeId: selectedVolumeId.value,
    chapterNumber: chapterNumber.value,
    chapterTitle: chapterTitle.value,
    blocks: blocks.value,
  })
}

watch([selectedVolumeId, chapterNumber, chapterTitle, blocks], persistDraft, { deep: true })

watch(
  volumes,
  (list) => {
    if (!list.length) return

    if (!selectedVolumeId.value) {
      selectedVolumeId.value = list[0]!.id
    }

    if (!existingDraft) {
      chapterNumber.value = getLastChapterNumberForVolume(list, selectedVolumeId.value)
    }
  },
  { immediate: true },
)

watch(selectedVolumeId, () => {
  chapterNumber.value = getLastChapterNumberForVolume(volumes.value, selectedVolumeId.value)
})

onMounted(() => {
  if (!existingDraft) {
    draftStore.createDraft(novelaId, selectedVolumeId.value, chapterNumber.value)
  }
})
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 relative overflow-hidden"
  >
    <!-- Background Layers -->
    <div class="absolute inset-0 z-0 pointer-events-none">
      <!-- Blurred Image -->
      <div v-if="novela?.poster_url" class="absolute inset-0">
        <img
          :src="staticImage(novela.poster_url)"
          class="w-full h-full object-cover blur-[15px] scale-110 opacity-[0.1] dark:opacity-[0.4]"
          alt=""
        />
        <!-- Overlay for better contrast -->
        <div class="absolute inset-0 bg-zinc-50/40 dark:bg-[#0f0f0f]/60"></div>
      </div>
    </div>

    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 relative z-10">
      <div class="flex items-center gap-4 mb-8">
        <button
          class="flex items-center p-2 rounded-xl bg-white/50 dark:bg-zinc-900/50 backdrop-blur-md border border-zinc-200/50 dark:border-zinc-800/50 hover:bg-zinc-200 dark:hover:bg-zinc-800 transition-colors cursor-pointer"
          @click="router.push(`/novela/${novelaId}`)"
        >
          <Icon name="ph:arrow-left-bold" size="20" />
        </button>
        <div>
          <h1 class="flex items-center gap-3 text-2xl font-black tracking-tight">
            <Icon name="ph:pencil-simple-line-bold" size="28" class="text-zinc-400" />
            Добавить главу
          </h1>
          <p v-if="novela" class="text-sm text-zinc-500 mt-1">{{ novela.title }}</p>
        </div>
      </div>

      <form class="space-y-6" @submit.prevent="onSubmit">
        <div
          class="bg-white/70 dark:bg-zinc-900/70 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-6 space-y-5 shadow-xl shadow-zinc-200/50 dark:shadow-none"
        >
          <div class="flex items-center gap-2 mb-2">
            <Icon name="ph:info-bold" size="18" class="text-zinc-400" />
            <h2 class="text-lg font-bold">Основная информация</h2>
          </div>

          <div class="relative z-20">
            <label
              class="flex items-center gap-2 text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2"
            >
              <Icon name="ph:books-bold" size="16" class="text-zinc-400" />
              Том
            </label>
            <BaseSelect
              v-model="selectedVolumeId"
              :options="volumeOptions"
              placeholder="Выберите том"
            />
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label
                class="flex items-center gap-2 text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2"
              >
                <Icon name="ph:hash-bold" size="16" class="text-zinc-400" />
                Номер главы
              </label>
              <input
                :value="chapterNumber"
                type="number"
                inputmode="numeric"
                min="1"
                step="1"
                class="w-full px-3 py-2.5 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-200 text-sm focus:ring-2 focus:ring-zinc-500 outline-none transition-all"
                :class="chapterNumberError ? 'border-red-400 dark:border-red-600' : ''"
                @input="onChapterNumberInput"
              />
              <p v-if="chapterNumberError" class="text-xs text-red-500 mt-1">
                {{ chapterNumberError }}
              </p>
            </div>
            <div>
              <label
                class="flex items-center gap-2 text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2"
              >
                <Icon name="ph:text-t-bold" size="16" class="text-zinc-400" />
                Название
              </label>
              <input
                v-model="chapterTitle"
                type="text"
                placeholder="Название главы"
                class="w-full px-3 py-2.5 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-200 text-sm focus:ring-2 focus:ring-zinc-500 outline-none transition-all"
              />
            </div>
          </div>
        </div>

        <div
          class="bg-white/70 dark:bg-zinc-900/70 border border-zinc-200 dark:border-zinc-800 rounded-2xl p-6 shadow-xl shadow-zinc-200/50 dark:shadow-none"
        >
          <div class="flex items-center justify-between mb-4">
            <h2 class="flex items-center gap-2 text-lg font-bold">
              <Icon name="ph:list-bullets-bold" size="20" class="text-zinc-400" />
              Содержание
            </h2>
            <div class="flex gap-2">
              <button
                type="button"
                class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors cursor-pointer"
                @click="addTextBlock(blocks.length - 1)"
              >
                <Icon name="ph:text-t-bold" size="14" />
                Текст
              </button>
              <button
                type="button"
                class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-yellow-300 dark:border-yellow-700 text-yellow-700 dark:text-yellow-400 hover:bg-yellow-50 dark:hover:bg-yellow-900/20 transition-colors cursor-pointer"
                @click="addImageBlock(blocks.length - 1)"
              >
                <Icon name="ph:image-bold" size="14" />
                Изображение
              </button>
            </div>
          </div>

          <div class="space-y-3">
            <template v-for="(block, index) in blocks" :key="block.id">
              <div @dragover="onDragOver($event, index)">
                <ChapterEditorBlock
                  v-if="block.type === 'text'"
                  :model-value="block.content"
                  :can-remove="blocks.length > 1"
                  @update:model-value="block.content = $event"
                  @remove="onRemoveBlock(index)"
                  @add-text-after="addTextAt(index)"
                  @add-image-after="addImageAt(index)"
                  @drag-start="onDragStart(index)"
                  @drag-end="onDragEnd"
                />
                <ChapterEditorImageBlock
                  v-else
                  :can-remove="blocks.length > 1"
                  :position="imageBlockOrder[index]!"
                  :preview-url="getImageById(block.content)?.previewUrl || undefined"
                  :caption="getImageById(block.content)?.caption || ''"
                  @update:caption="onImageCaptionUpdate(index, $event)"
                  @file-selected="onImageSelected(index, $event)"
                  @remove="onRemoveBlock(index)"
                  @add-text-after="addTextAt(index)"
                  @add-image-after="addImageAt(index)"
                  @drag-start="onDragStart(index)"
                  @drag-end="onDragEnd"
                  @preview-image="previewImage"
                />
              </div>
            </template>
          </div>

          <p class="flex items-center gap-2 text-xs text-zinc-400 mt-3">
            <Icon name="ph:info-bold" size="14" />
            Перетаскивайте блоки за иконку слева для изменения порядка. Позиция изображения
            определяется порядком блока.
          </p>
        </div>

        <div class="flex items-center justify-end gap-3">
          <button
            type="button"
            class="flex items-center gap-2 px-4 py-2 rounded-xl text-sm font-medium border border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors cursor-pointer"
            @click="router.push(`/novela/${novelaId}`)"
          >
            <Icon name="ph:x-bold" size="16" />
            Отмена
          </button>
          <button
            type="submit"
            :disabled="isSubmitting"
            class="flex items-center gap-2 px-6 py-2 rounded-xl text-sm font-bold bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900 hover:bg-zinc-800 dark:hover:bg-zinc-200 disabled:opacity-50 transition-all cursor-pointer"
          >
            <Icon v-if="!isSubmitting" name="ph:floppy-disk-bold" size="18" />
            <Icon v-else name="ph:circle-notch-bold" size="18" class="animate-spin" />
            {{ isSubmitting ? 'Добавление...' : 'Добавить главу' }}
          </button>
        </div>
      </form>
    </div>

    <ImagePreviewModal v-if="previewImageUrl" :image-url="previewImageUrl" @close="closePreview" />
  </div>
</template>
