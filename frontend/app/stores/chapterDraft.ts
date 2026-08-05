import { defineStore } from 'pinia'
import { uid } from '~/utils/chapter'

export type BlockType = 'text' | 'image'

export interface ChapterBlock {
  id: string
  type: BlockType
  content: string
}

export interface ChapterImage {
  id: string
  previewUrl: string
  caption: string
  position: number
  file: null
}

export interface ChapterDraft {
  novelaId: string
  selectedVolumeId: string
  chapterNumber: number
  chapterTitle: string
  blocks: ChapterBlock[]
  images: ChapterImage[]
}

const STORAGE_KEY = 'kyokusu-chapter-drafts'

function load(): Record<string, ChapterDraft> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : {}
  } catch {
    return {}
  }
}

function save(drafts: Record<string, ChapterDraft>) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(drafts))
  } catch {
    // quota exceeded — silent
  }
}

export const useChapterDraftStore = defineStore('chapterDraft', {
  state: () => ({
    drafts: load() as Record<string, ChapterDraft>,
  }),

  actions: {
    getDraft(novelaId: string): ChapterDraft | null {
      return this.drafts[novelaId] ?? null
    },

    createDraft(novelaId: string, volumeId: string, chapterNumber: number): ChapterDraft {
      const draft: ChapterDraft = {
        novelaId,
        selectedVolumeId: volumeId,
        chapterNumber,
        chapterTitle: '',
        blocks: [{ id: uid('block'), type: 'text', content: '' }],
        images: [],
      }
      this.drafts[novelaId] = draft
      save(this.drafts)
      return draft
    },

    updateDraft(novelaId: string, updates: Partial<ChapterDraft>) {
      const draft = this.drafts[novelaId]
      if (draft) {
        Object.assign(draft, updates)
        save(this.drafts)
      }
    },

    clearDraft(novelaId: string) {
      delete this.drafts[novelaId]
      save(this.drafts)
    },
  },
})
