import { $api } from '~/composables/api/useApi'
import type { BackendNotification, NotificationStats } from '@/types/backend/notification'

export function useNotifications() {
  const list = useState<BackendNotification[]>('notifications-list', () => [])
  const unreadCount = useState<number>('notifications-unread', () => 0)
  const stream = useState<EventSource | null>('notifications-stream', () => null)
  const { notify } = useNotificationStore()

  let reconnectAttempts = 0
  const MAX_RECONNECT_DELAY = 30000

  function connect() {
    if (import.meta.server) return

    disconnect()
    reconnectAttempts = 0

    const url = '/api/notifications/stream'
    const es = new EventSource(url, { withCredentials: true })

    es.addEventListener('notification', (event: MessageEvent) => {
      try {
        const notification: BackendNotification = JSON.parse(event.data)
        list.value = [notification, ...list.value]

        if (!notification.isRead) {
          unreadCount.value++
          notify({
            title: notification.title,
            content: notification.message,
            type: 'info',
          })
        }
      } catch {}
    })

    es.onerror = () => {
      es.close()
      scheduleReconnect()
    }

    stream.value = es
  }

  function scheduleReconnect() {
    const delay = Math.min(1000 * 2 ** reconnectAttempts, MAX_RECONNECT_DELAY)
    reconnectAttempts++
    setTimeout(() => connect(), delay)
  }

  function disconnect() {
    if (stream.value) {
      stream.value.close()
      stream.value = null
    }
  }

  async function fetchList(limit?: number, offset?: number) {
    try {
      const data = await $api<BackendNotification[]>('/notifications', {
        query: { limit, offset },
      })
      list.value = data
      return data
    } catch {
      return []
    }
  }

  async function fetchStats() {
    try {
      const data = await $api<NotificationStats>('/notifications/stats')
      unreadCount.value = data.unreadCount ?? 0
      return data
    } catch {
      return { unreadCount: 0 }
    }
  }

  async function markRead(id: number) {
    try {
      await $api(`/notifications/${id}/read`, { method: 'PATCH' })
      const item = list.value.find((n) => n.id === id)
      if (item && !item.isRead) {
        item.isRead = true
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
    } catch {}
  }

  async function markAllRead() {
    try {
      await $api('/notifications/read-all', { method: 'PATCH' })
      list.value.forEach((n) => {
        n.isRead = true
      })
      unreadCount.value = 0
    } catch {}
  }

  async function remove(id: number) {
    try {
      await $api(`/notifications/${id}`, { method: 'DELETE' })
      const item = list.value.find((n) => n.id === id)
      if (item && !item.isRead) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
      list.value = list.value.filter((n) => n.id !== id)
    } catch {}
  }

  return {
    list,
    unreadCount,
    connect,
    disconnect,
    fetchList,
    fetchStats,
    markRead,
    markAllRead,
    remove,
  }
}
