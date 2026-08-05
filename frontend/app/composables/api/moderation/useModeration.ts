import { $api } from '@/composables/api/useApi'
import type { NovelaVolume, NovelaChapter } from '@/types/backend/novela'
import { useNotificationStore } from '@/stores/notification'

export interface PendingContent {
  volumes: NovelaVolume[]
  chapters: NovelaChapter[]
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

  const approveVolume = async (id: number) => {
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

  const rejectVolume = async (id: number) => {
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

  const approveChapter = async (id: number) => {
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

  const rejectChapter = async (id: number) => {
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

  return {
    pendingContent,
    isLoading,
    isActionLoading,
    fetchPending,
    approveVolume,
    rejectVolume,
    approveChapter,
    rejectChapter,
  }
}
