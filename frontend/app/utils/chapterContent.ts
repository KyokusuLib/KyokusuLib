import { uid } from '~/utils/chapter'
import type { ChapterBlock } from '~/stores/chapterDraft'
import type { UploadedImage } from '~/composables/ui/useImageUpload'

export interface ParseResult {
  blocks: ChapterBlock[]
  images: UploadedImage[]
}

export function escapeHtml(str: string): string {
  return str
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
    .replaceAll("'", '&#39;')
}

function decodeHtmlEntities(text: string): string {
  return text
    .replaceAll('&amp;', '&')
    .replaceAll('&lt;', '<')
    .replaceAll('&gt;', '>')
    .replaceAll('&quot;', '"')
    .replaceAll('&#39;', "'")
}

export function parseChapterContent(content: string): ParseResult {
  const blocks: ChapterBlock[] = []
  const images: UploadedImage[] = []

  if (!content.trim()) {
    return { blocks: [{ id: uid('block'), type: 'text', content: '' }], images: [] }
  }

  let imgCounter = 0
  const pendingText: string[] = []

  function flushText() {
    if (pendingText.length > 0) {
      blocks.push({ id: uid('block'), type: 'text', content: pendingText.join('\n\n') })
      pendingText.length = 0
    }
  }

  for (const part of content.split('\n\n')) {
    const trimmed = part.trim()
    if (!trimmed) continue

    const figureMatch = trimmed.match(/<figure[^>]*>(.*?)<\/figure>/s)
    if (figureMatch) {
      flushText()

      const imgTagMatch = figureMatch[1].match(/<img[^>]*>/)
      const captionMatch = figureMatch[1].match(/<figcaption>(.*?)<\/figcaption>/)
      if (imgTagMatch) {
        const srcMatch = imgTagMatch[0].match(/src="(.*?)"/)
        const altMatch = imgTagMatch[0].match(/alt="(.*?)"/)
        if (srcMatch) {
          const src = srcMatch[1] || ''
          const caption = captionMatch ? decodeHtmlEntities(captionMatch[1]) : altMatch?.[1] || ''
          const imgId = `img-edit-${Date.now()}-${++imgCounter}`
          images.push({ id: imgId, sourceUrl: src, previewUrl: src, caption })
          blocks.push({ id: uid('block'), type: 'image', content: imgId })
        }
      }
    } else {
      pendingText.push(trimmed)
    }
  }

  flushText()

  if (blocks.length === 0) {
    blocks.push({ id: uid('block'), type: 'text', content: '' })
  }

  return { blocks, images }
}

export function buildChapterContent(blocks: ChapterBlock[], images: UploadedImage[]): string {
  return blocks
    .map((block) => {
      if (block.type === 'text') return block.content
      const img = images.find((i) => i.id === block.content)
      if (!img) return ''
      const caption = img.caption?.trim()
        ? `<figcaption>${escapeHtml(img.caption)}</figcaption>`
        : ''
      return `<figure style="text-align: center;"><img src="${img.sourceUrl}" alt="${escapeHtml(img.caption || 'chapter image')}" style="display:block;width:100%;max-width:28rem;height:auto;margin:0 auto;object-fit:cover;border-radius:0.75rem;box-shadow:0 10px 15px -3px rgba(0,0,0,0.1),0 4px 6px -4px rgba(0,0,0,0.1);cursor:pointer"/>${caption}</figure>`
    })
    .filter(Boolean)
    .join('\n\n')
}
