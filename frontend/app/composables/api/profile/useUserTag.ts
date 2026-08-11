import { useNotificationStore } from '@/stores/notification'
import { $api } from '@/composables/api/useApi'
import type { UserTagDTO } from '@/types/backend/user'
import type { ResponseMessage } from '@/types/backend/response_message'

export function useUserTag() {
  const { notify } = useNotificationStore()

  const updateUserTags = async (userId: number, tagIds: number[]): Promise<boolean> => {
    try {
      await $api<ResponseMessage>('/user/tags', {
        method: 'PUT',
        body: { user_id: userId, tag_ids: tagIds },
      })

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
      const data = await $api<ResponseMessage>('/user/tag', {
        method: 'PUT',
        body: dto,
      })

      if (data?.message) {
        notify({
          type: 'success',
          title: 'User tag updated successfully',
          content: data.message,
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
