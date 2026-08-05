<script setup lang="ts">
import { RichText as BaseRichTextEditor, Tooltip } from '@kyokusu-ui/vue'
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps<{
  modelValue: string
  canRemove: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  remove: []
  addTextAfter: []
  addImageAfter: []
  dragStart: []
  dragEnd: []
  splitText: [payload: { beforeHTML: string; selectedHTML: string; afterHTML: string }]
}>()

const contextMenu = ref({ visible: false, x: 0, y: 0 })

function onContextMenu(e: MouseEvent) {
  const sel = window.getSelection()
  if (!sel || sel.isCollapsed || !sel.rangeCount) return

  const target = e.target as HTMLElement | null
  if (!target?.closest('[contenteditable]')) return

  e.preventDefault()
  contextMenu.value = { visible: true, x: e.clientX, y: e.clientY }
}

function splitSelectedText() {
  const sel = window.getSelection()
  if (!sel || sel.isCollapsed || !sel.rangeCount) {
    closeContextMenu()
    return
  }

  const range = sel.getRangeAt(0)
  const container = range.commonAncestorContainer
  const editorEl = (
    container instanceof HTMLElement
      ? container.closest('[contenteditable]')
      : container.parentElement?.closest('[contenteditable]')
  ) as HTMLElement | null

  if (!editorEl) {
    closeContextMenu()
    return
  }

  // Selected HTML
  const selectedDiv = document.createElement('div')
  selectedDiv.appendChild(range.cloneContents())
  const selectedHTML = selectedDiv.innerHTML

  // HTML before selection
  const beforeRange = document.createRange()
  beforeRange.selectNodeContents(editorEl)
  beforeRange.setEnd(range.startContainer, range.startOffset)
  const beforeDiv = document.createElement('div')
  beforeDiv.appendChild(beforeRange.cloneContents())
  const beforeHTML = beforeDiv.innerHTML

  // HTML after selection
  const afterRange = document.createRange()
  afterRange.selectNodeContents(editorEl)
  afterRange.setStart(range.endContainer, range.endOffset)
  const afterDiv = document.createElement('div')
  afterDiv.appendChild(afterRange.cloneContents())
  const afterHTML = afterDiv.innerHTML

  closeContextMenu()
  emit('splitText', { beforeHTML, selectedHTML, afterHTML })
}

function closeContextMenu() {
  contextMenu.value.visible = false
}

function onKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') closeContextMenu()
}

onMounted(() => {
  document.addEventListener('click', closeContextMenu)
  document.addEventListener('keydown', onKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('click', closeContextMenu)
  document.removeEventListener('keydown', onKeyDown)
})
</script>

<template>
  <div
    class="group relative rounded-xl border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900/50 hover:border-zinc-400 dark:hover:border-zinc-500 transition-colors"
    @contextmenu="onContextMenu"
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
      <span class="text-xs font-medium text-zinc-500 uppercase tracking-wider">Текст</span>
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
          <button type="button" class="rounded-md cursor-pointer" @click="emit('remove')">
            <Icon name="ph:trash-bold" size="14" class="text-red-500" />
          </button>
        </Tooltip>
      </div>
    </div>

    <div class="p-3">
      <BaseRichTextEditor
        :model-value="modelValue"
        @update:model-value="emit('update:modelValue', $event)"
        placeholder="Введите текст главы..."
      />
    </div>

    <Teleport to="body">
      <div
        v-if="contextMenu.visible"
        class="fixed z-99 px-2 py-1 min-w-50 bg-white dark:bg-zinc-800 hover:bg-zinc-100 dark:hover:bg-zinc-700 rounded-lg shadow-xl border border-zinc-200 dark:border-zinc-700 overflow-hidden"
        :style="{ left: `${contextMenu.x}px`, top: `${contextMenu.y}px` }"
      >
        <button
          type="button"
          class="w-full text-left px-3 py-2 text-sm text-zinc-700 dark:text-zinc-300 flex items-center gap-2 cursor-pointer"
          @click.stop="splitSelectedText"
        >
          <Icon name="ph:square-split-horizontal-bold" size="16" />
          Выделить в отдельный блок
        </button>
      </div>
    </Teleport>
  </div>
</template>
