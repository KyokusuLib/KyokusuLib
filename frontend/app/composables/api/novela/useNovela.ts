import type { NovelaDetails } from '@/types/backend/novela'
import { $api } from '@/composables/api/useApi'
import { useNotificationStore } from '@/stores/notification'
import type { NovelsQueryParams } from '@/types/frontend/query/novela-query'
import type { NovelaFilters } from '@/types/frontend/novela/novela-filters'

export function useNovela() {
  const novela = useState<NovelaDetails | null>('novela-data', () => null)

  const { notify } = useNotificationStore()

  const isLoading = ref(false)
  const isUpdating = ref(false)

  const handleError = (e: any, defaultMessage: string) => {
    console.error(e)
    notify({
      type: 'error',
      title: 'Ошибка',
      content: e?.data?.message || e?.message || defaultMessage,
    })
  }

  const fetchNovels = async (params: NovelsQueryParams | NovelaFilters) => {
    isLoading.value = true
    try {
      const data = await $api<NovelaDetails[]>('/api/novela', {
        params,
      })
      return data
    } catch (e) {
      handleError(e, 'Не удалось загрузить каталог новелл')
      return []
    } finally {
      isLoading.value = false
    }
  }

  const fetchNovela = async (id: string | number) => {
    isLoading.value = true
    try {
      const data = await $api<NovelaDetails>(`/api/novela/${id}`)
      novela.value = data
      return data
    } catch (e) {
      handleError(e, 'Не удалось загрузить данные новеллы')
    } finally {
      isLoading.value = false
    }
  }

  const updateNovela = async (id: number, payload: any, file?: File | null) => {
    isUpdating.value = true
    const formData = new FormData()

    formData.append('data', JSON.stringify(payload))
    if (file) formData.append('poster', file)

    try {
      await $api(`/api/novela/${id}`, {
        method: 'PUT',
        body: formData,
      })

      notify({
        type: 'success',
        title: 'Успех',
        content: 'Новелла успешно обновлена',
      })

      await fetchNovela(id)
    } catch (e) {
      handleError(e, 'Ошибка при обновлении новеллы')
    } finally {
      isUpdating.value = false
    }
  }

  const addVolume = async (novelaId: number, volumeNumber: number, title: string) => {
    try {
      const res = await $api<{ id: string; message: string; status: string }>(
        `/api/novela/${novelaId}/volumes`,
        {
          method: 'POST',
          body: { volume_number: volumeNumber, title },
        },
      )
      notify({
        type: 'success',
        title: 'Успех',
        content: res.message || 'Том успешно добавлен',
      })
      return { id: res.id, title, number: volumeNumber, status: res.status }
    } catch (e) {
      handleError(e, 'Не удалось добавить том')
    }
  }

  const addChapter = async (
    novelaId: number,
    volumeId: string | number,
    chapterNumber: number,
    title: string,
    content: string,
  ) => {
    try {
      const res = await $api<{ id: string; message: string; status: string }>(
        `/api/novela/volumes/${volumeId}/chapters`,
        {
          method: 'POST',
          body: { chapter_number: chapterNumber, title, content },
        },
      )
      return res
    } catch (e) {
      handleError(e, 'Не удалось добавить главу')
    }
  }

  const addChapterImage = async (
    chapterId: string,
    imageUrl: string,
    caption: string,
    position: number,
  ) => {
    try {
      await $api(`/api/novela/chapters/${chapterId}/images`, {
        method: 'POST',
        body: { image_url: imageUrl, caption, position },
      })
    } catch (e) {
      handleError(e, 'Не удалось добавить изображение')
    }
  }

  const deleteChapterImages = async (chapterId: string) => {
    try {
      await $api(`/api/novela/chapters/${chapterId}/images`, {
        method: 'DELETE',
      })
    } catch (e) {
      handleError(e, 'Не удалось удалить изображения')
    }
  }

  const deleteVolume = async (novelaId: number, volumeId: string) => {
    try {
      await $api(`/api/novela/volumes/${volumeId}`, {
        method: 'DELETE',
      })

      if (novela.value) {
        novela.value = {
          ...novela.value,
          volumes: novela.value.volumes.filter((v) => v.id !== volumeId),
        }
      }
      notify({
        type: 'success',
        title: 'Успех',
        content: 'Том удален',
      })
    } catch (e) {
      handleError(e, 'Не удалось удалить том')
    }
  }

  const deleteChapter = async (novelaId: number, chapterId: string) => {
    try {
      await $api(`/api/novela/chapters/${chapterId}`, {
        method: 'DELETE',
      })

      if (novela.value) {
        novela.value = {
          ...novela.value,
          volumes: novela.value.volumes.map((volume) => ({
            ...volume,
            chapters: volume.chapters.filter((c) => c.id !== chapterId),
          })),
        }
      }
      notify({
        type: 'success',
        title: 'Успех',
        content: 'Глава удалена',
      })
    } catch (e) {
      handleError(e, 'Не удалось удалить главу')
    }
  }

  const updateChapter = async (
    novelaId: number,
    chapterId: string,
    chapterNumber: number,
    title: string,
    content: string,
  ) => {
    try {
      await $api(`/api/novela/chapters/${chapterId}`, {
        method: 'PUT',
        body: { chapter_number: chapterNumber, title, content },
      })

      if (novela.value) {
        novela.value = {
          ...novela.value,
          volumes: novela.value.volumes.map((volume) => ({
            ...volume,
            chapters: volume.chapters.map((c) =>
              c.id === chapterId ? { ...c, title, number: chapterNumber, content } : c,
            ),
          })),
        }
      }
      notify({
        type: 'success',
        title: 'Успех',
        content: 'Глава обновлена',
      })
    } catch (e) {
      handleError(e, 'Не удалось обновить главу')
    }
  }

  return {
    novela,
    isLoading,
    isUpdating,
    fetchNovela,
    fetchNovels,
    updateNovela,
    addVolume,
    deleteVolume,
    addChapter,
    addChapterImage,
    deleteChapterImages,
    deleteChapter,
    updateChapter,
  }
}
