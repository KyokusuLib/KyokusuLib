import { useApi, $api } from '~/composables/api/useApi'
import { useRolePermissions } from '../role/useRolePermissions'

import type { GetUserDto } from '@/types/backend/user'
import type { DashboardRowUserStatus } from '~/types/enums/dashboard-table'
import { KyokusuAppRole } from '~/types/enums/role-enum'

export function useUserApi() {
  const { hasPermission } = useRolePermissions()
  const { notify } = useNotificationStore()

  const getUser = async (userId: number): Promise<GetUserDto | null> => {
    if (!userId) return null

    try {
      const { data, error } = await useApi<GetUserDto>(`/api/user/${userId}`)

      if (error.value) {
        console.error('USER FETCHER ERROR:', error.value)
        return null
      }
      return data.value!
    } catch (e) {
      console.error('USER FETCHER EXCEPTION:', e)
      return null
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
    deleteUser,
    updateUserStatus,
  }
}
