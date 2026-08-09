import { $api } from '@/composables/api/useApi'
import { useNotificationStore } from '@/stores/notification'

export interface PendingVolume {
  id: string
  novela_id: number
  novela_title: string
  number: number
  title: string
  status: string
  created_by: number
  created_by_name: string
}

export interface PendingChapter {
  id: string
  volume_id: string
  volume_number: number
  novela_id: number
  novela_title: string
  number: number
  title: string
  content: string
  status: string
  created_by: number
  created_by_name: string
}

export interface PendingContent {
  volumes: PendingVolume[]
  chapters: PendingChapter[]
}

export function useModeration() {
  const pendingContent = useState<PendingContent | null>('pending-content', () => null)
  const isLoading = useState('moderation-loading', () => false)
  const isActionLoading = useState('moderation-action-loading', () => false)
  const { notify } = useNotificationStore()

  const fetchPending = async () => {
    isLoading.value = true
    try {
      const data = await $api<PendingContent>('/api/moderation/pending')
      pendingContent.value = data
    } catch (e: any) {
      console.error(e)
      notify({
        type: 'error',
        title: 'Ошибка',
        content: e.message || 'Не удалось загрузить данные для модерации',
      })
    } finally {
      isLoading.value = false
    }
  }

  const approveVolume = async (id: string) => {
    isActionLoading.value = true
    try {
      await $api(`/api/moderation/volumes/${id}/approve`, { method: 'POST' })
      notify({ type: 'success', title: 'Успех', content: 'Том одобрен' })
      if (pendingContent.value) {
        pendingContent.value.volumes = pendingContent.value.volumes.filter((v) => v.id !== id)
      }
    } catch (e: any) {
      console.error(e)
      notify({ type: 'error', title: 'Ошибка', content: e.message || 'Не удалось одобрить том' })
    } finally {
      isActionLoading.value = false
    }
  }

  const rejectVolume = async (id: string) => {
    isActionLoading.value = true
    try {
      await $api(`/api/moderation/volumes/${id}/reject`, { method: 'POST' })
      notify({ type: 'success', title: 'Успех', content: 'Том отклонен' })
      if (pendingContent.value) {
        pendingContent.value.volumes = pendingContent.value.volumes.filter((v) => v.id !== id)
      }
    } catch (e: any) {
      console.error(e)
      notify({ type: 'error', title: 'Ошибка', content: e.message || 'Не удалось отклонить том' })
    } finally {
      isActionLoading.value = false
    }
  }

  const approveChapter = async (id: string) => {
    isActionLoading.value = true
    try {
      await $api(`/api/moderation/chapters/${id}/approve`, { method: 'POST' })
      notify({ type: 'success', title: 'Успех', content: 'Глава одобрена' })
      if (pendingContent.value) {
        pendingContent.value.chapters = pendingContent.value.chapters.filter((c) => c.id !== id)
      }
    } catch (e: any) {
      console.error(e)
      notify({ type: 'error', title: 'Ошибка', content: e.message || 'Не удалось одобрить главу' })
    } finally {
      isActionLoading.value = false
    }
  }

  const rejectChapter = async (id: string) => {
    isActionLoading.value = true
    try {
      await $api(`/api/moderation/chapters/${id}/reject`, { method: 'POST' })
      notify({ type: 'success', title: 'Успех', content: 'Глава отклонена' })
      if (pendingContent.value) {
        pendingContent.value.chapters = pendingContent.value.chapters.filter((c) => c.id !== id)
      }
    } catch (e: any) {
      console.error(e)
      notify({ type: 'error', title: 'Ошибка', content: e.message || 'Не удалось отклонить главу' })
    } finally {
      isActionLoading.value = false
    }
  }

  const updateVolume = async (id: string, payload: { volume_number: number; title: string }) => {
    isActionLoading.value = true
    try {
      await $api(`/api/moderation/volumes/${id}`, { method: 'PUT', body: payload })
      notify({ type: 'success', title: 'Успех', content: 'Том обновлен' })
      if (pendingContent.value) {
        pendingContent.value.volumes = pendingContent.value.volumes.map((v) =>
          v.id === id ? { ...v, number: payload.volume_number, title: payload.title } : v,
        )
      }
      return true
    } catch (e: any) {
      console.error(e)
      notify({ type: 'error', title: 'Ошибка', content: e.message || 'Не удалось обновить том' })
      return false
    } finally {
      isActionLoading.value = false
    }
  }

  const updateChapter = async (
    id: string,
    payload: { chapter_number: number; title: string; content: string },
  ) => {
    isActionLoading.value = true
    try {
      await $api(`/api/moderation/chapters/${id}`, { method: 'PUT', body: payload })
      notify({ type: 'success', title: 'Успех', content: 'Глава обновлена' })
      if (pendingContent.value) {
        pendingContent.value.chapters = pendingContent.value.chapters.map((c) =>
          c.id === id
            ? {
                ...c,
                number: payload.chapter_number,
                title: payload.title,
                content: payload.content,
              }
            : c,
        )
      }
      return true
    } catch (e: any) {
      console.error(e)
      notify({ type: 'error', title: 'Ошибка', content: e.message || 'Не удалось обновить главу' })
      return false
    } finally {
      isActionLoading.value = false
    }
  }

  return {
    pendingContent,
    isLoading,
    isActionLoading,
    fetchPending,
    approveVolume,
    rejectVolume,
    approveChapter,
    rejectChapter,
    updateVolume,
    updateChapter,
  }
}
