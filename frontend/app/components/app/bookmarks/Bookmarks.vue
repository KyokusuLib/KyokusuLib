<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useBookmark } from '~/composables/api/novela/useBookmark'
import type { BookmarkCategory } from '@/types/frontend/bookmarks'

import BookmarkHeader from './BookmarkHeader.vue'
import BookmarkTabs from './BookmarkTabs.vue'
import BookmarkGrid from './BookmarkGrid.vue'
import BookmarkCategoriesEdit from './BookmarkCategoriesEdit.vue'

const authStore = useAuthStore()
const activeCategory = ref<BookmarkCategory>('reading')

const {
  fetchUserBookmarkNovels,
  fetchBookmarkCategories,
  userBookmarkNovels,
  bookmarkCategories,
  setBookmark,
  removeBookmark,
  loading,
} = useBookmark()

const load = () => {
  if (authStore.user) fetchUserBookmarkNovels(authStore.user.id, activeCategory.value)
}

onMounted(async () => {
  if (authStore.user) {
    await fetchBookmarkCategories(authStore.user.id)
    const readingCat = bookmarkCategories.value.find((c: any) => c.name === 'reading')
    if (readingCat) activeCategory.value = readingCat.id
    load()
  }
})
watch(activeCategory, load)

const handleStatusChange = async (novelaId: number, newCategory: BookmarkCategory | 'remove') => {
  try {
    if (newCategory === 'remove') {
      await removeBookmark(novelaId)
    } else {
      await setBookmark(novelaId, newCategory)
    }

    userBookmarkNovels.value = userBookmarkNovels.value.filter((n) => n.id !== novelaId)

    const newCategorie = bookmarkCategories.value.find((c) => c.id === newCategory)
    newCategorie.count++

    const oldCategorie = bookmarkCategories.value.find((c) => c.id === activeCategory.value)
    oldCategorie.count--
  } catch (error) {
    console.error('Error changing bookmark status', error)
  }
}

const totalBookmarkCounts = computed(() => {
  return bookmarkCategories.value.reduce((acc, c) => acc + c.count, 0)
})
</script>

<template>
  <div
    class="min-h-screen flex flex-col bg-zinc-50 dark:bg-radial-[at_center] dark:from-zinc-800 dark:to-zinc-950 dark:to-90% transition-colors duration-300"
  >
    <div class="w-full max-w-6xl mx-auto px-4 sm:px-6 py-8">
      <BookmarkHeader :totalBookmarkCounts="totalBookmarkCounts" />
      <div class="space-y-8">
        <div class="flex items-start gap-4">
          <BookmarkTabs
            :categories="bookmarkCategories"
            v-model:activeCategory="activeCategory"
            class="flex-1 overflow-x-auto [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden"
          />
          <BookmarkCategoriesEdit />
        </div>
        <BookmarkGrid
          :loading="loading"
          :novels="userBookmarkNovels"
          :categories="bookmarkCategories"
          @change-status="handleStatusChange"
        />
      </div>
    </div>
  </div>
</template>
