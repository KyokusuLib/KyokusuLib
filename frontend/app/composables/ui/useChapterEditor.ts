import type { ChapterBlock } from '~/stores/chapterDraft'
import { uid } from '~/utils/chapter'

export function useChapterEditor(blocks: Ref<ChapterBlock[]>) {
  function addTextBlock(afterIndex: number, content = '') {
    blocks.value.splice(afterIndex + 1, 0, {
      id: uid('block'),
      type: 'text',
      content,
    })
  }

  function addImageBlock(afterIndex: number) {
    blocks.value.splice(afterIndex + 1, 0, {
      id: uid('block'),
      type: 'image',
      content: '',
    })
  }

  function removeBlock(index: number) {
    if (blocks.value.length <= 1) return
    blocks.value.splice(index, 1)
  }

  function reorderBlock(fromIndex: number, toIndex: number) {
    if (fromIndex === toIndex) return
    const [item] = blocks.value.splice(fromIndex, 1)
    blocks.value.splice(toIndex, 0, item)
  }

  return {
    addTextBlock,
    addImageBlock,
    removeBlock,
    reorderBlock,
  }
}
