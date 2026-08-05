import { $api } from '@/composables/api/useApi'
import type { NovelaChapter, NovelaDetails } from '@/types/backend/novela'
import { useAuthStore } from '@/stores/auth'
import { ref, watch, nextTick } from 'vue'
import { useThrottleFn, useEventListener } from '@vueuse/core'
import type { UserLevel } from '~/types/backend/user'

export function useReadProgress() {
  const { isAuthenticated, user } = useAuthStore()
  const { notify } = useNotificationStore()

  function isChapterRead(chapter: NovelaChapter): boolean {
    return chapter.is_read ?? false
  }

  function getLastReadChapterId(novela: NovelaDetails): string | null {
    return novela.last_readed?.chapter_id ?? null
  }

  function getLastReadChapterNumber(novela: NovelaDetails): number | null {
    return novela.last_readed?.chapter_number ?? null
  }

  function getLastReadChapterScroll(novela: NovelaDetails): number | null {
    return novela.last_readed?.chapter_scroll ?? null
  }

  function getFirstChapterId(novela: NovelaDetails): string | null {
    const firstVolume = novela.volumes?.[0]

    if (!firstVolume?.chapters?.length) {
      return null
    }

    return firstVolume.chapters[0]!.id
  }

  function getContinueReadingUrl(novela: NovelaDetails): string | null {
    if (!isAuthenticated) {
      const firstId = getFirstChapterId(novela)

      return firstId ? `/novela/reader/${novela.id}/${firstId}` : null
    }

    const lastId = getLastReadChapterId(novela)

    if (lastId) {
      return `/novela/reader/${novela.id}/${lastId}`
    }

    const firstId = getFirstChapterId(novela)

    return firstId ? `/novela/reader/${novela.id}/${firstId}` : null
  }

  async function saveReadPosition(chapterId: string, scrollPosition: number) {
    if (!isAuthenticated) {
      return
    }

    try {
      await $api('/api/novela/chapters/progress', {
        method: 'POST',
        body: {
          chapter_id: chapterId,
          scroll_position: Math.round(scrollPosition),
        },
      })
    } catch (e) {
      console.error('Failed to save read position:', e)
    }
  }

  async function markChapterAsRead(chapterId: string): Promise<boolean> {
    let isRead = false

    if (!isAuthenticated) {
      return isRead
    }

    try {
      const data = await $api<UserLevel>(`/api/novela/chapter/${chapterId}/mark-as-read`, {
        method: 'POST',
      })

      if (data !== null || data !== undefined) {
        user!.user_level = data
        isRead = true
      }
    } catch (e: any) {
      if (e instanceof Error && e.message === 'Ошибка: Chapter already read') {
        isRead = true
        notify({
          title: 'Внимание',
          type: 'warning',
          content: 'Награда за прочтение этой главы уже получена',
        })
      } else {
        notify({
          title: 'Ошибка',
          type: 'error',
          content: e instanceof Error ? e.message : String(e),
        })
      }
    }

    return isRead
  }

  return {
    isChapterRead,
    getLastReadChapterId,
    getLastReadChapterNumber,
    getLastReadChapterScroll,
    getFirstChapterId,
    getContinueReadingUrl,
    saveReadPosition,
    markChapterAsRead,
  }
}

export function useReaderScrollPosition(chapterId: string, initialScroll = 0) {
  const { saveReadPosition } = useReadProgress()
  const { isAuthenticated } = useAuthStore()

  const currentChapterId = ref(chapterId)
  const lastSavedOffset = ref(0)

  function getChapterElement(id: string) {
    return document.querySelector<HTMLElement>(`[data-chapter-id="${id}"]`)
  }

  function getCurrentOffset() {
    const chapterEl = getChapterElement(currentChapterId.value)

    if (!chapterEl) {
      return 0
    }

    const chapterTop = chapterEl.getBoundingClientRect().top + window.scrollY

    return Math.max(0, Math.round(window.scrollY - chapterTop))
  }

  const throttledSave = useThrottleFn(async () => {
    if (!isAuthenticated) return

    const offset = getCurrentOffset()

    lastSavedOffset.value = offset

    await saveReadPosition(currentChapterId.value, offset)
  }, 500)

  function setupScrollTracking() {
    if (!isAuthenticated) {
      return
    }

    useEventListener(window, 'scroll', throttledSave, { passive: true })
  }

  async function restoreScroll() {
    if (Number.isNaN(initialScroll) || initialScroll < 0) {
      return
    }

    await nextTick()

    let attempts = 0

    const restore = () => {
      const chapterEl = getChapterElement(currentChapterId.value)

      if (!chapterEl) {
        attempts++

        if (attempts < 30) {
          setTimeout(restore, 100)
        }

        return
      }

      const chapterTop = chapterEl.getBoundingClientRect().top + window.scrollY

      window.scrollTo({
        top: chapterTop + initialScroll,
        behavior: 'auto',
      })
    }

    restore()

    const images = Array.from(document.querySelectorAll('img'))

    await Promise.all(
      images.map((img) => {
        if (img.complete) {
          return Promise.resolve()
        }

        return new Promise<void>((resolve) => {
          img.addEventListener('load', () => resolve(), { once: true })

          img.addEventListener('error', () => resolve(), { once: true })
        })
      }),
    )

    await nextTick()

    restore()

    setTimeout(restore, 100)
    setTimeout(restore, 300)
    setTimeout(restore, 700)
  }

  async function saveOnLeave() {
    if (!isAuthenticated) {
      return
    }

    await saveReadPosition(currentChapterId.value, getCurrentOffset())
  }

  watch(currentChapterId, async (newId, oldId) => {
    if (!isAuthenticated || !oldId || oldId === newId) {
      return
    }

    const oldChapterEl = getChapterElement(oldId)

    if (oldChapterEl) {
      const chapterTop = oldChapterEl.getBoundingClientRect().top + window.scrollY

      const offset = Math.max(0, Math.round(window.scrollY - chapterTop))

      await saveReadPosition(oldId, offset)
    }

    lastSavedOffset.value = 0
  })

  return {
    currentChapterId,
    setupScrollTracking,
    restoreScroll,
    saveOnLeave,
  }
}
