import { useImageUpload } from '~/composables/ui/useImageUpload'
import { useChapterEditor } from '~/composables/ui/useChapterEditor'
import { buildChapterContent } from '~/utils/chapterContent'
import type { ChapterBlock } from '~/stores/chapterDraft'
import type { UploadedImage } from '~/composables/ui/useImageUpload'

export function useChapterForm(initial?: { blocks?: ChapterBlock[]; images?: UploadedImage[] }) {
  const blocks = ref<ChapterBlock[]>(
    initial?.blocks ?? [{ id: `${Date.now()}-block-0`, type: 'text' as const, content: '' }],
  )

  const { images, addFile, removeImageById, getImageById, updateCaption, clearAll } =
    useImageUpload()

  if (initial?.images) {
    for (const img of initial.images) {
      images.value.push(img)
    }
  }

  const { addTextBlock, addImageBlock, removeBlock, reorderBlock } = useChapterEditor(blocks)

  const imageBlockOrder = computed(() => {
    let counter = 0
    return blocks.value.map((b) => (b.type === 'image' ? ++counter : 0))
  })

  const dragIndex = ref<number | null>(null)
  const hoverIndex = ref<number | null>(null)

  function onDragStart(index: number) {
    dragIndex.value = index
    hoverIndex.value = index
  }

  function onDragOver(e: DragEvent, index: number) {
    e.preventDefault()
    if (dragIndex.value === null) return
    hoverIndex.value = index
  }

  function onDrop(index: number) {
    if (dragIndex.value === null) return
    if (dragIndex.value !== index) {
      reorderBlock(dragIndex.value, index)
    }
    hoverIndex.value = null
    dragIndex.value = null
  }

  function onDragEnd() {
    hoverIndex.value = null
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
    if (block.content) removeImageById(block.content)
    const item = await addFile(file)
    block.content = item.id
  }

  function onImageCaptionUpdate(index: number, caption: string) {
    const block = blocks.value[index]
    if (block?.type === 'image' && block.content) updateCaption(block.content, caption)
  }

  function onRemoveBlock(index: number) {
    const block = blocks.value[index]
    if (block?.type === 'image' && block.content) removeImageById(block.content)
    removeBlock(index)
  }

  const previewImageUrl = ref<string | null>(null)
  function previewImage(url: string) {
    previewImageUrl.value = url
  }
  function closePreview() {
    previewImageUrl.value = null
  }

  function getContent(): string {
    return buildChapterContent(blocks.value, images.value)
  }

  onBeforeUnmount(() => clearAll())

  return {
    blocks,
    images,
    imageBlockOrder,
    dragIndex,
    hoverIndex,
    previewImageUrl,
    addTextAt,
    addImageAt,
    onImageSelected,
    onImageCaptionUpdate,
    onRemoveBlock,
    onDragStart,
    onDragOver,
    onDragEnd,
    onDrop,
    previewImage,
    closePreview,
    getImageById,
    getContent,
  }
}
