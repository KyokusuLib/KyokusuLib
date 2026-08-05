import { ref } from 'vue'
import { $api } from '@/composables/api/useApi'
import type { NovelaAuthorDetails } from '~/types/backend/novela'

export function useAuthors() {
  const isSearching = ref(false)
  const foundAuthors = ref<NovelaAuthorDetails[]>([])

  const searchAuthors = async (query: string) => {
    if (!query || query.length < 2) {
      foundAuthors.value = []
      return
    }

    isSearching.value = true
    try {
      const data = await $api<NovelaAuthorDetails[]>('/api/author', {
        query: { search: query },
      })

      foundAuthors.value = data.map((a) => ({
        id: a.id,
        name: a.name,
        label: a.name,
        country: a.country,
        picture: a.picture,
        metier: a.metier,
        bio: a.bio,
      }))
    } catch {
      foundAuthors.value = []
    } finally {
      isSearching.value = false
    }
  }

  const getAuthorById = async (authorId: string) => {
    isSearching.value = true

    try {
      const data = await $api<NovelaAuthorDetails>(`/api/author/${authorId}`)
      return data
    } catch {
      return null
    } finally {
      isSearching.value = false
    }
  }

  return {
    isSearching,
    foundAuthors,
    getAuthorById,
    searchAuthors,
  }
}
