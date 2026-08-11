import { $api } from '~/composables/api/useApi'
import { useRolePermissions } from '../role/useRolePermissions'

import type { GetUserDto } from '@/types/backend/user'
import type { DashboardRowUserStatus } from '~/types/enums/dashboard-table'
import { KyokusuAppRole } from '~/types/enums/role-enum'

export interface UpdateUserPayload {
  name: string
  about: string
  gender: string
  birthday: string
  is_public: boolean
  is_show_tag: boolean
  is_show_bookmark: boolean
  role: string
  status: string
}

export function useUserApi() {
  const { hasPermission } = useRolePermissions()
  const { notify } = useNotificationStore()

  const getUser = async (userId: number): Promise<GetUserDto | null> => {
    if (!userId) return null

    try {
      return await $api<GetUserDto>(`/user/${userId}`)
    } catch (e) {
      console.error('USER FETCHER ERROR:', e)
      return null
    }
  }

  const updateUser = async (userId: number, payload: UpdateUserPayload): Promise<boolean> => {
    try {
      await $api(`/user/${userId}`, {
        method: 'PUT',
        body: payload,
      })
      return true
    } catch (e) {
      console.error('USER UPDATE ERROR:', e)
      return false
    }
  }

  const deleteUser = async (userId: number) => {
    if (!hasPermission(KyokusuAppRole.ADMIN)) {
      notify({
        title: 'Неудача',
        content: 'У вас нет доступа для использования функции',
        type: 'info',
      })
      return
    }

    try {
      await $api(`/user/${userId}`, {
        method: 'DELETE',
      })

      notify({
        title: 'Успех',
        content: 'Пользователь удалён',
        type: 'success',
      })
    } catch (e: any) {
      notify({
        title: 'Ошибка',
        content: e.message ?? 'Не удалось удалить пользователя',
        type: 'error',
      })
    }
  }

  const updateUserStatus = async (userId: number, status: DashboardRowUserStatus) => {
    if (!hasPermission(KyokusuAppRole.MODERATOR)) {
      notify({
        title: 'Неудача',
        content: 'У вас нет доступа для использования функции',
        type: 'info',
      })
      return
    }

    try {
      await $api(`/user/${userId}/status`, {
        method: 'PUT',
        body: {
          status,
          lastActive: Date.now(),
        },
      })

      notify({
        title: 'Успех',
        content: 'Статус пользователя обновлён',
        type: 'success',
      })
    } catch (e: any) {
      notify({
        title: 'Ошибка',
        content: e.message ?? 'Не удалось обновить статус',
        type: 'error',
      })
    }
  }

  return {
    getUser,
    updateUser,
    deleteUser,
    updateUserStatus,
  }
}
