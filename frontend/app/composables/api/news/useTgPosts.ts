import { $api } from '~/composables/api/useApi'
import type { BackendTgPost } from '@/types/backend/tg-post'

export function useTgPosts() {
  const posts = useState<BackendTgPost[]>('tg-posts', () => [])
  const stream = useState<EventSource | null>('tg-posts-stream', () => null)

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

  async function fetchTgPosts(limit = 50) {
    try {
      const data = await $api<BackendTgPost[]>('/tg/posts', { query: { limit } })
      posts.value = data
      return data
    } catch {
      return []
    }
  }

  async function deleteTgPost(id: number) {
    await $api(`/tg/posts/${id}`, { method: 'DELETE' })
    posts.value = posts.value.filter((p) => p.id !== id)
  }

  return {
    posts,
    fetchTgPosts,
    connectStream,
    disconnectStream,
    deleteTgPost,
  }
}
