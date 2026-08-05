import { useNovela } from '~/composables/api/novela/useNovela'
import { useNotificationStore } from '~/stores/notification'
import type { UploadedImage } from '~/composables/ui/useImageUpload'
import type { ChapterBlock } from '~/stores/chapterDraft'

export function useChapterSubmit() {
  const { addChapter } = useNovela()
  const { notify } = useNotificationStore()
  const isSubmitting = ref(false)

  function escapeHtml(str: string) {
    return str
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll("'", '&#39;')
  }

  async function submit(
    novelaId: string,
    volumeId: string,
    chapterNumber: number,
    chapterTitle: string,
    blocks: ChapterBlock[],
    images: UploadedImage[],
  ) {
    const hasText = blocks.some((b) => b.type === 'text' && b.content.trim())
    if (!hasText) {
      notify({ type: 'error', title: 'Ошибка', content: 'Добавьте хотя бы один текстовый блок' })
      return null
    }

    if (!chapterTitle) {
      chapterTitle = `Глава ${chapterNumber}`
    }

    const content = blocks
      .map((block) => {
        if (block.type === 'text') {
          return block.content
        }

        const img = images.find((i) => i.id === block.content)
        if (!img) return ''

        const caption = img.caption?.trim()
          ? `<figcaption>${escapeHtml(img.caption)}</figcaption>`
          : ''
        return `<figure><img src="${img.sourceUrl}" alt="${escapeHtml(img.caption || 'chapter image')}"/>${caption}</figure>`
      })
      .filter(Boolean)
      .join('\n\n')

    isSubmitting.value = true
    try {
      const res = await addChapter(Number(novelaId), volumeId, chapterNumber, chapterTitle, content)

      notify({ type: 'success', title: 'Глава добавлена', content: chapterTitle })
      return res
    } catch (e: any) {
      notify({ type: 'error', title: 'Ошибка', content: e?.message || 'Не удалось добавить главу' })
      return null
    } finally {
      isSubmitting.value = false
    }
  }

  return {
    isSubmitting,
    submit,
  }
}
