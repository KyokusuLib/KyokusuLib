import { $api } from '~/composables/api/useApi'
import type { BackendTgPost } from '@/types/backend/tg-post'

export function useTgPosts() {
  const posts = useState<BackendTgPost[]>('tg-posts', () => [])
  const stream = useState<EventSource | null>('tg-posts-stream', () => null)
  const nextCursor = useState<string | null>('tg-posts-next-cursor', () => null)
  const hasMore = useState<boolean>('tg-posts-has-more', () => false)
  const isLoadingMore = useState<boolean>('tg-posts-loading-more', () => false)

  let reconnectAttempts = 0
  const MAX_RECONNECT_DELAY = 30000

  function connectStream() {
    if (import.meta.server) return

    disconnectStream()
    reconnectAttempts = 0

    const es = new EventSource('/api/tg/stream')

    es.addEventListener('tg_post', (event: MessageEvent) => {
      try {
        const post: BackendTgPost = JSON.parse(event.data)
        posts.value = [post, ...posts.value.filter((p) => p.id !== post.id)].slice(0, 100)
      } catch {}
    })

    es.addEventListener('tg_post_deleted', (event: MessageEvent) => {
      try {
        const { id } = JSON.parse(event.data) as { id: number }
        posts.value = posts.value.filter((p) => p.id !== id)
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
    setTimeout(() => connectStream(), delay)
  }

  function disconnectStream() {
    if (stream.value) {
      stream.value.close()
      stream.value = null
    }
  }

  async function fetchTgPosts(limit = 50, cursor?: string) {
    isLoadingMore.value = Boolean(cursor)
    try {
      const response = await $fetch.raw<BackendTgPost[]>('/api/tg/posts', {
        query: { limit, cursor },
        headers: useRequestHeaders(['cookie']),
      })
      const data = response._data ?? []

      if (cursor) {
        const known = new Set(posts.value.map((post) => post.id))
        posts.value = [...posts.value, ...data.filter((post) => !known.has(post.id))]
      } else {
        posts.value = data
      }

      nextCursor.value = response.headers.get('x-next-cursor')
      hasMore.value = response.headers.get('x-has-more') === 'true'
      return data
    } catch {
      if (!cursor) {
        posts.value = []
        nextCursor.value = null
        hasMore.value = false
      }
      return []
    } finally {
      isLoadingMore.value = false
    }
  }

  async function loadMoreTgPosts() {
    if (!nextCursor.value || isLoadingMore.value) return []
    return fetchTgPosts(50, nextCursor.value)
  }

  async function deleteTgPost(id: number) {
    await $api(`/tg/posts/${id}`, { method: 'DELETE' })
    posts.value = posts.value.filter((p) => p.id !== id)
  }

  return {
    posts,
    hasMore,
    isLoadingMore,
    fetchTgPosts,
    loadMoreTgPosts,
    connectStream,
    disconnectStream,
    deleteTgPost,
  }
}
