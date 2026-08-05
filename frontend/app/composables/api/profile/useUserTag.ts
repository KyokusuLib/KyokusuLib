import { useNotificationStore } from '@/stores/notification'
import { useApi } from '@/composables/api/useApi'
import type { UserTagDTO } from '@/types/backend/user'
import type { ResponseMessage } from '@/types/backend/response_message'

export function useUserTag() {
  const { notify } = useNotificationStore()

  const updateUserTags = async (userId: number, tagIds: number[]): Promise<boolean> => {
    try {
      const { error } = await useApi<ResponseMessage>('/api/user/tags', {
        method: 'PUT',
        body: { user_id: userId, tag_ids: tagIds },
      })

      if (error.value) {
        throw new Error(error.value.message)
      }

      notify({
        title: 'Успех!',
        content: 'Профиль был обновлен',
        type: 'success',
      })

      return true
    } catch (e: any) {
      console.error(e)
      return false
    }
  }

  const updateUserTag = async (dto: UserTagDTO) => {
    try {
      const { data, error } = await useApi<ResponseMessage>('/api/user/tag', {
        method: 'PUT',
        body: dto,
      })

      if (error.value) {
        throw new Error(error.value.message)
      }

      if (data.value?.message) {
        notify({
          type: 'success',
          title: 'User tag updated successfully',
          content: data.value.message,
        })
      }
    } catch (e: any) {
      console.error(e)
    }
  }

  return {
    updateUserTag,
    updateUserTags,
  }
}
