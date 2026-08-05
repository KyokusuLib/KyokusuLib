import { $api } from '@/composables/api/useApi'
import type { NovelaDetails } from '~/types/backend/novela'
import type { BookmarkCategory } from '~/types/frontend/bookmarks'

export function useBookmark() {
  const userBookmarkNovels = useState<NovelaDetails[]>('user-category-bookmarks', () => [])
  const loading = ref(false)

  const bookmarkCategories = useState<any[]>('user-bookmark-categories', () => [
    { id: 'planned', label: 'В планах', icon: 'ph:book-open-bold' },
    { id: 'reading', label: 'Читаю', icon: 'ph:calendar-blank-bold' },
    { id: 'completed', label: 'Прочитано', icon: 'ph:check-circle-bold' },
    { id: 'on_hold', label: 'Отложено', icon: 'ph:pause-circle-bold' },
    { id: 'dropped', label: 'Брошено', icon: 'ph:prohibit-bold' },
  ])

  const fetchBookmarkCategories = async (user_id: number) => {
    try {
      const data = await $api<any[]>(`/api/bookmarks/categories?user_id=${user_id}`, {
        method: 'GET',
      })
      if (data) {
        bookmarkCategories.value = data.map((cat: any) => {
          let icon = 'ph:bookmark-simple-bold'
          let label = cat.name

          if (cat.name === 'planned') {
            icon = 'ph:book-open-bold'
            label = 'В планах'
          } else if (cat.name === 'reading') {
            icon = 'ph:calendar-blank-bold'
            label = 'Читаю'
          } else if (cat.name === 'completed') {
            icon = 'ph:check-circle-bold'
            label = 'Прочитано'
          } else if (cat.name === 'on_hold') {
            icon = 'ph:pause-circle-bold'
            label = 'Отложено'
          } else if (cat.name === 'dropped') {
            icon = 'ph:prohibit-bold'
            label = 'Брошено'
          }

          return {
            id: cat.id || cat.name,
            name: cat.name,
            label: label,
            icon: icon,
            user_id: cat.user_id,
            count: cat.count || 0,
            visible: cat.visible,
          }
        })
      }
    } catch (e) {
      console.error('Failed to fetch bookmark categories:', e)
    }
  }

  const setBookmark = async (novelaId: number, category: BookmarkCategory) => {
    loading.value = true
    try {
      const body: Record<string, any> = { novela_id: novelaId }
      if (typeof category === 'number') {
        body.category_id = category
      } else {
        body.category = category
      }

      const data = await $api('/api/novela/bookmark', {
        method: 'POST',
        body,
      })

      return data
    } catch (e) {
      console.error('Failed to set bookmark:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  const removeBookmark = async (novelaId: number) => {
    loading.value = true
    try {
      const data = await $api(`/api/novela/bookmark/${novelaId}`, {
        method: 'DELETE',
      })

      return data
    } catch (e) {
      console.error('Failed to remove bookmark:', e)
      throw e
    } finally {
      loading.value = false
    }
  }

  const fetchUserBookmarkNovels = async (userID: number, category: BookmarkCategory) => {
    loading.value = true
    try {
      const data = await $api<NovelaDetails[]>('/api/novela/bookmarks', {
        query: { user_id: userID, category: category },
      })

      userBookmarkNovels.value = data || []
    } catch {
      userBookmarkNovels.value = []
    } finally {
      loading.value = false
    }
  }

  const createBookmarkCategory = async (user_id: number, name: string) => {
    try {
      const data = await $api<{ id: number; message: string }>('/api/bookmarks/categories', {
        method: 'POST',
        body: { name },
      })

      bookmarkCategories.value.push({
        id: data.id,
        user_id: user_id,
        label: name,
        icon: 'ph:bookmark-simple-bold',
        count: 0,
        visible: true,
      })
    } catch (e) {
      console.error('Failed to create bookmark category:', e)
      throw e
    }
  }

  const updateBookmarkCategory = async (id: number, name: string, visible: boolean) => {
    try {
      await $api(`/api/bookmarks/categories/${id}`, {
        method: 'PUT',
        body: { name, visible },
      })
      bookmarkCategories.value = bookmarkCategories.value.map((cat) =>
        cat.id === id ? { ...cat, label: name, visible } : cat,
      )
    } catch (e) {
      console.error('Failed to update bookmark category:', e)
      throw e
    }
  }

  const deleteBookmarkCategory = async (id: number) => {
    try {
      await $api(`/api/bookmarks/categories/${id}`, {
        method: 'DELETE',
      })
      bookmarkCategories.value = bookmarkCategories.value.filter((cat) => cat.id !== id)
    } catch (e) {
      console.error('Failed to delete bookmark category:', e)
      throw e
    }
  }

  return {
    loading,
    bookmarkCategories,
    fetchBookmarkCategories,
    createBookmarkCategory,
    updateBookmarkCategory,
    deleteBookmarkCategory,
    setBookmark,
    removeBookmark,
    fetchUserBookmarkNovels,
    userBookmarkNovels,
  }
}
