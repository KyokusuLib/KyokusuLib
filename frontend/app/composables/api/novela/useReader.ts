import { useApi } from '~/composables/api/useApi'
import type { ChapterReaderResponse } from '@/types/backend/novela'

export function useReader() {
  const chapters = useState<ChapterReaderResponse[]>('reader-chapters', () => [])
  const isLoading = useState('reader-loading', () => false)
  const isAppending = useState('reader-appending', () => false)

  const fetchChapter = async (id: string, append = false) => {
    if (append) isAppending.value = true
    else isLoading.value = true

    try {
      const { data, error } = await useApi<ChapterReaderResponse>(`/api/novela/chapters/${id}`)
      if (error.value) throw error.value

      if (append) {
        const isAlreadyLoaded = chapters.value.some((c) => c.id === data.value?.id)
        if (!isAlreadyLoaded) {
          chapters.value.push(data.value!)
        }
      } else {
        chapters.value = [data.value!]
      }

      return data.value
    } finally {
      isLoading.value = false
      isAppending.value = false
    }
  }

  return {
    chapters,
    isLoading,
    isAppending,
    fetchChapter,
  }
}
