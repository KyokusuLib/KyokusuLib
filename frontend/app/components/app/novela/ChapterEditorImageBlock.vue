<script setup lang="ts">
import { Tooltip } from '@kyokusu-ui/vue'

defineProps<{
  canRemove: boolean
  previewUrl?: string
  caption: string
  position: number
}>()

const emit = defineEmits<{
  'update:caption': [value: string]
  fileSelected: [file: File]
  remove: []
  addTextAfter: []
  addImageAfter: []
  dragStart: []
  dragEnd: []
  previewImage: [url: string]
}>()

const fileInput = ref<HTMLInputElement | null>(null)
const isDragOver = ref(false)

function onFileChange(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (file) emit('fileSelected', file)
}

function onDrop(e: DragEvent) {
  e.preventDefault()
  isDragOver.value = false
  const file = e.dataTransfer?.files?.[0]
  if (file && file.type.startsWith('image/')) {
    emit('fileSelected', file)
  }
}

function onDragOver(e: DragEvent) {
  e.preventDefault()
  isDragOver.value = true
}

function onDragLeave() {
  isDragOver.value = false
}
</script>

<template>
  <div
    class="group relative rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900/50 hover:border-zinc-400 dark:hover:border-zinc-500 transition-colors"
    :class="
      isDragOver
        ? 'border-yellow-400 dark:border-yellow-500 bg-yellow-50/50 dark:bg-yellow-900/20'
        : ''
    "
    @drop="onDrop"
    @dragover="onDragOver"
    @dragleave="onDragLeave"
  >
    <div class="flex items-center gap-2 px-3 py-2 border-b border-zinc-200 dark:border-zinc-700/50">
      <Icon
        name="ph:dots-six-vertical-bold"
        size="14"
        class="text-zinc-400 cursor-grab active:cursor-grabbing"
        draggable="true"
        @dragstart="emit('dragStart')"
        @dragend="emit('dragEnd')"
      />
      <span class="text-xs font-medium text-zinc-500 uppercase tracking-wider">
        Изображение #{{ position }}
      </span>
      <div class="ml-auto flex gap-4 opacity-0 group-hover:opacity-100 transition-opacity">
        <Tooltip text="Добавить текстовый блок ниже">
          <button type="button" class="rounded-md cursor-pointer" @click="emit('addTextAfter')">
            <Icon name="ph:text-t-bold" size="14" class="text-zinc-500" />
          </button>
        </Tooltip>
        <Tooltip text="Добавить изображение ниже">
          <button type="button" class="rounded-md cursor-pointer" @click="emit('addImageAfter')">
            <Icon name="ph:image-bold" size="14" class="text-zinc-500" />
          </button>
        </Tooltip>
        <Tooltip v-if="canRemove" text="Удалить блок">
          <button
            type="button"
            class="rounded-md transition-colors cursor-pointer"
            @click="emit('remove')"
          >
            <Icon name="ph:trash-bold" size="14" class="text-red-500" />
          </button>
        </Tooltip>
      </div>
    </div>

    <div class="p-3">
      <div
        v-if="!previewUrl"
        class="flex flex-col items-center justify-center gap-3 rounded-lg border-2 border-dashed border-zinc-200 dark:border-zinc-700 p-8 cursor-pointer hover:border-yellow-500 dark:hover:border-yellow-500 hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors"
        @click="fileInput?.click()"
      >
        <Icon name="ph:image-bold" size="32" class="text-zinc-400" />
        <p class="text-sm text-zinc-500 text-center">
          Перетащите изображение сюда<br />или нажмите для выбора
        </p>
        <input ref="fileInput" type="file" accept="image/*" class="hidden" @change="onFileChange" />
      </div>

      <template v-else>
        <div
          class="rounded-lg overflow-hidden border border-zinc-200 dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-900/50 mb-3 cursor-zoom-in hover:opacity-90 transition-opacity"
          @click="emit('previewImage', previewUrl)"
        >
          <img :src="previewUrl" class="max-h-75 w-auto mx-auto object-contain" />
        </div>
        <input
          type="text"
          placeholder="Подпись к изображению"
          class="w-full px-3 py-2.5 rounded-lg border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-200 text-sm focus:ring-2 focus:ring-zinc-500 outline-none transition-all"
          @input="emit('update:caption', ($event.target as HTMLInputElement).value)"
        />
      </template>
    </div>
  </div>
</template>
