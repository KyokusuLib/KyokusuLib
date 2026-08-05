<script setup lang="ts">
import type { ChapterBlock } from '~/stores/chapterDraft'
import type { UploadedImage } from '~/composables/ui/useImageUpload'
import { useChapterForm } from '~/composables/ui/useChapterForm'
import { staticImage } from '~/utils/str'
import { Select as BaseSelect } from '@kyokusu-ui/vue'
import ChapterEditorBlock from '~/components/app/novela/ChapterEditorBlock.vue'
import ChapterEditorImageBlock from '~/components/app/novela/ChapterEditorImageBlock.vue'
import ImagePreviewModal from '~/components/app/novela/ImagePreviewModal.vue'

const props = defineProps<{
  volumes: { value: string; label: string }[]
  selectedVolumeId: string
  chapterNumber: number
  chapterTitle: string
  posterUrl?: string
  title: string
  subtitle: string
  submitLabel: string
  submittingLabel: string
  isSubmitting: boolean
  // Pre‑populate form for edit mode
  initialBlocks?: ChapterBlock[]
  initialImages?: UploadedImage[]
}>()

const emit = defineEmits<{
  'update:selectedVolumeId': [v: string]
  'update:chapterNumber': [v: number]
  'update:chapterTitle': [v: string]
  submit: [content: string]
}>()

const form = reactive(useChapterForm({ blocks: props.initialBlocks, images: props.initialImages }))
const { previewImageUrl } = form

const chapterNumberError = ref('')

defineExpose({ form })

function onSplitText(
  index: number,
  {
    beforeHTML,
    selectedHTML,
    afterHTML,
  }: { beforeHTML: string; selectedHTML: string; afterHTML: string },
) {
  form.blocks[index].content = beforeHTML + afterHTML
  form.addTextAt(index)
  form.blocks[index + 1].content = selectedHTML
}

function onChapterNumberInput(e: Event) {
  const input = e.target as HTMLInputElement
  const val = input.value
  const num = Number(val)
  if (val === '' || (Number.isInteger(num) && num >= 1)) {
    emit('update:chapterNumber', num || 1)
    chapterNumberError.value = ''
  } else {
    chapterNumberError.value = 'Номер главы должен быть целым числом'
  }
}
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 relative overflow-hidden"
  >
    <div class="absolute inset-0 z-0 pointer-events-none">
      <div v-if="posterUrl" class="absolute inset-0">
        <img
          :src="staticImage(posterUrl)"
          class="w-full h-full object-cover blur-[15px] scale-110 opacity-[0.1] dark:opacity-[0.4]"
          alt=""
        />
        <div class="absolute inset-0 bg-zinc-50/40 dark:bg-[#0f0f0f]/60"></div>
      </div>
    </div>

    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8 relative z-10">
      <div class="flex items-center gap-4 mb-8">
        <slot name="back" />
        <div>
          <h1 class="flex items-center gap-3 text-2xl font-black tracking-tight">
            <Icon name="ph:pencil-simple-line-bold" size="28" class="text-zinc-400" />
            {{ title }}
          </h1>
          <p class="text-sm text-zinc-500 mt-1">{{ subtitle }}</p>
        </div>
      </div>

      <form class="space-y-6" @submit.prevent="emit('submit', form.getContent())">
        <div
          class="bg-white/70 dark:bg-zinc-900/70 backdrop-blur-xl border border-zinc-200 dark:border-zinc-800 rounded-2xl p-6 space-y-5 shadow-xl shadow-zinc-200/50 dark:shadow-none"
        >
          <div class="flex items-center gap-2 mb-2">
            <Icon name="ph:info-bold" size="18" class="text-zinc-400" />
            <h2 class="text-lg font-bold">Основная информация</h2>
          </div>

          <div>
            <label
              class="flex items-center gap-2 text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2"
            >
              <Icon name="ph:books-bold" size="16" class="text-zinc-400" />
              Том
            </label>
            <BaseSelect
              :model-value="selectedVolumeId"
              :options="volumes"
              placeholder="Выберите том"
              @update:model-value="emit('update:selectedVolumeId', $event)"
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
                :value="chapterTitle"
                type="text"
                placeholder="Название главы"
                class="w-full px-3 py-2.5 rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-200 text-sm focus:ring-2 focus:ring-zinc-500 outline-none transition-all"
                @input="emit('update:chapterTitle', ($event.target as HTMLInputElement).value)"
              />
            </div>
          </div>
        </div>

        <div
          class="bg-white/70 dark:bg-zinc-900/70 backdrop-blur-xl border border-zinc-200 dark:border-zinc-800 rounded-2xl p-6 shadow-xl shadow-zinc-200/50 dark:shadow-none"
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
                @click="form.addTextAt(form.blocks.length - 1)"
              >
                <Icon name="ph:text-t-bold" size="14" />
                Текст
              </button>
              <button
                type="button"
                class="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-yellow-300 dark:border-yellow-700 text-yellow-700 dark:text-yellow-400 hover:bg-yellow-50 dark:hover:bg-yellow-900/20 transition-colors cursor-pointer"
                @click="form.addImageAt(form.blocks.length - 1)"
              >
                <Icon name="ph:image-bold" size="14" />
                Изображение
              </button>
            </div>
          </div>

          <div class="space-y-3">
            <template v-for="(block, index) in form.blocks" :key="block.id">
              <div
                @dragover="form.onDragOver($event, index)"
                @drop.prevent="form.onDrop(index)"
                :class="
                  form.hoverIndex === index && form.hoverIndex !== form.dragIndex
                    ? 'ring-2 ring-yellow-400/60 rounded-xl'
                    : ''
                "
              >
                <ChapterEditorBlock
                  v-if="block.type === 'text'"
                  :model-value="block.content"
                  :can-remove="form.blocks.length > 1"
                  @update:model-value="block.content = $event"
                  @remove="form.onRemoveBlock(index)"
                  @add-text-after="form.addTextAt(index)"
                  @add-image-after="form.addImageAt(index)"
                  @drag-start="form.onDragStart(index)"
                  @drag-end="form.onDragEnd"
                  @split-text="onSplitText(index, $event)"
                />
                <ChapterEditorImageBlock
                  v-else
                  :can-remove="form.blocks.length > 1"
                  :position="form.imageBlockOrder[index]!"
                  :preview-url="form.getImageById(block.content)?.previewUrl || undefined"
                  :caption="form.getImageById(block.content)?.caption || ''"
                  @update:caption="form.onImageCaptionUpdate(index, $event)"
                  @file-selected="form.onImageSelected(index, $event)"
                  @remove="form.onRemoveBlock(index)"
                  @add-text-after="form.addTextAt(index)"
                  @add-image-after="form.addImageAt(index)"
                  @drag-start="form.onDragStart(index)"
                  @drag-end="form.onDragEnd"
                  @preview-image="form.previewImage"
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
          <slot name="actions" />
          <button
            type="submit"
            :disabled="isSubmitting"
            class="flex items-center gap-2 px-6 py-2 rounded-xl text-sm font-bold bg-zinc-900 text-white dark:bg-zinc-100 dark:text-zinc-900 hover:bg-zinc-800 dark:hover:bg-zinc-200 disabled:opacity-50 transition-all cursor-pointer"
          >
            <Icon v-if="!isSubmitting" name="ph:floppy-disk-bold" size="18" />
            <Icon v-else name="ph:circle-notch-bold" size="18" class="animate-spin" />
            {{ isSubmitting ? submittingLabel : submitLabel }}
          </button>
        </div>
      </form>
    </div>

    <ImagePreviewModal
      v-if="form.previewImageUrl"
      :image-url="form.previewImageUrl"
      @close="form.closePreview"
    />
  </div>
</template>
