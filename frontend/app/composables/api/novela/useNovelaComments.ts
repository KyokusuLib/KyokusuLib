import { $api } from '../useApi'
import type { NovelaCommentResponse } from '~/types/backend/novela'

export function useNovelaComments() {
  const comments = useState<NovelaCommentResponse[]>('novela-comments', () => [])
  const isLoading = ref(false)

  const totalComments = computed(() => {
    const countReplies = (comment: NovelaCommentResponse): number =>
      1 + (comment.replies?.reduce((acc, r) => acc + countReplies(r), 0) ?? 0)
    return comments.value.reduce((acc, c) => acc + countReplies(c), 0)
  })

  const fetchComments = async (novelaId: number) => {
    isLoading.value = true
    try {
      comments.value = await $api<NovelaCommentResponse[]>(`/api/novela/comments/${novelaId}`)
      return comments.value
    } catch {
      return null
    } finally {
      isLoading.value = false
    }
  }

  const addComment = async (payload: {
    novela_id: number
    content: string
    parent_id?: number | null
  }) => {
    return await $api('/api/novela/comments', {
      method: 'POST',
      body: payload,
    })
  }

  const deleteComment = async (commentId: number) => {
    return await $api(`/api/novela/comments/${commentId}`, {
      method: 'DELETE',
    })
  }

  const updateComment = async (commentId: number, content: string) => {
    return await $api(`/api/novela/comments/${commentId}`, {
      method: 'PUT',
      body: { content, updated_at: new Date().toISOString() },
    })
  }

  const setCommentLike = async (commentId: number) => {
    return await $api(`/api/novela/comments/like/${commentId}`, {
      method: 'POST',
    })
  }

  const unsetCommentLike = async (commentId: number) => {
    return await $api(`/api/novela/comments/unlike/${commentId}`, {
      method: 'POST',
    })
  }

  const reportComment = async (commentId: number, reason: string) => {
    return await $api(`/api/novela/comments/${commentId}/report?reason=${reason}`, {
      method: 'POST',
    })
  }

  const cancelCommentReport = async (commentId: number) => {
    return await $api(`/api/novela/comments/${commentId}/report`, {
      method: 'DELETE',
    })
  }

  return {
    comments,
    totalComments,
    isLoading,
    fetchComments,
    addComment,
    deleteComment,
    updateComment,
    setCommentLike,
    unsetCommentLike,
    reportComment,
    cancelCommentReport,
  }
}
