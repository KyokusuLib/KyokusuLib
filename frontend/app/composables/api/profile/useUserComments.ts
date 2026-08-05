import { $api } from '../useApi'
import type { UserCommentResponse } from '~/types/backend/comments'

export function useUserComments() {
  const comments = useState<UserCommentResponse[]>('user-comments', () => [])
  const isLoading = ref(false)

  const fetchUserComments = async (userID: number) => {
    isLoading.value = true
    try {
      const data = await $api<UserCommentResponse[]>(`/api/user/${userID}/comments`)
      comments.value = data || []
      return data
    } catch {
      comments.value = []
    } finally {
      isLoading.value = false
    }
  }

  return {
    comments,
    fetchUserComments,
    isLoading,
  }
}
