export type NotificationType = 'success' | 'error' | 'info' | 'warning'

export interface NotificationPayload {
  title: string
  content?: string
  type?: NotificationType
  duration?: number
}

export interface NotificationItem extends NotificationPayload {
  id: string
}
