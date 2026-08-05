export interface BackendNotification {
  id: number
  userId: number
  title: string
  message: string
  isRead: boolean
  createdAt: string
}

export interface NotificationStats {
  unreadCount: number
}
