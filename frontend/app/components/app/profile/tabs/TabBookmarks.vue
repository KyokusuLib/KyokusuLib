<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import type { BookmarkCategory } from '@/types/frontend/bookmarks'
import NovelaCard from '@/components/app/novela/NovelaCard.vue'
import { useBookmark } from '~/composables/api/novela/useBookmark'

const props = defineProps<{ userId: number }>()

const { user } = useAuthStore()
const activeCategory = ref<BookmarkCategory>('reading')

const {
  fetchUserBookmarkNovels,
  fetchBookmarkCategories,
  userBookmarkNovels,
  bookmarkCategories,
  loading,
} = useBookmark()

const load = () => fetchUserBookmarkNovels(props.userId, activeCategory.value)

onMounted(async () => {
  await fetchBookmarkCategories(props.userId)
  const readingCat = bookmarkCategories.value.find((c: any) => c.name === 'reading')
  if (readingCat) activeCategory.value = readingCat.id
})

watch(activeCategory, load, { immediate: true })

const scrollContainer = ref<HTMLElement | null>(null)
let isDown = false
let startX = 0
let scrollLeft = 0
const isDragging = ref(false)

const onMouseDown = (e: MouseEvent) => {
  isDown = true
  isDragging.value = false
  if (!scrollContainer.value) return
  startX = e.pageX - scrollContainer.value.offsetLeft
  scrollLeft = scrollContainer.value.scrollLeft
}

const onMouseLeave = () => {
  isDown = false
  isDragging.value = false
}

const onMouseUp = () => {
  isDown = false
  setTimeout(() => {
    isDragging.value = false
  }, 0)
}

const onMouseMove = (e: MouseEvent) => {
  if (!isDown || !scrollContainer.value) return
  e.preventDefault()
  const x = e.pageX - scrollContainer.value.offsetLeft
  const walk = (x - startX) * 2
  if (Math.abs(walk) > 5) {
    isDragging.value = true
  }
  scrollContainer.value.scrollLeft = scrollLeft - walk
}

const handleTabClick = (e: Event, catId: BookmarkCategory) => {
  if (isDragging.value) {
    e.preventDefault()
    e.stopPropagation()
    return
  }
  activeCategory.value = catId
}

const tabVisible = (cat: any) => {
  if (user && user.id === props.userId) return true

  return cat.visible && (cat.user_id === props.userId || cat.user_id === null)
}
</script>

<template>
  <div class="space-y-6">
    <div
      ref="scrollContainer"
      class="flex gap-2 overflow-x-auto pb-2 no-scrollbar cursor-grab active:cursor-grabbing select-none [scrollbar-width:none] [-ms-overflow-style:none] [&::-webkit-scrollbar]:hidden"
      @mousedown="onMouseDown"
      @mouseleave="onMouseLeave"
      @mouseup="onMouseUp"
      @mousemove="onMouseMove"
    >
      <div v-for="cat in bookmarkCategories" :key="cat.id">
        <button
          v-if="tabVisible(cat)"
          @click="(e) => handleTabClick(e, cat.id as BookmarkCategory)"
          class="px-4 py-2 rounded-xl text-xs font-black uppercase tracking-widest transition-all whitespace-nowrap border cursor-pointer"
          :class="[
            activeCategory === cat.id
              ? 'bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 border-transparent shadow-lg'
              : 'bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 text-zinc-500 hover:border-zinc-400',
          ]"
        >
          <div class="flex justify-center items-center gap-3">
            <Icon :name="cat.icon" size="16" />
            <div class="flex items-center gap-2">
              <span>{{ cat.label }}</span>
              <span
                v-if="cat.count !== undefined"
                class="text-[10px] border border-zinc-500 bg-zinc-100 dark:bg-zinc-800 text-zinc-500 h-4 w-4 px-1.5 py-0.5 rounded-full"
                :class="
                  activeCategory === cat.id
                    ? 'bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-zinc-200'
                    : ''
                "
              >
                <div class="relative right-0.5 -top-0.5">
                  {{ cat.count }}
                </div>
              </span>
            </div>
          </div>
        </button>
      </div>
    </div>

    <div class="min-h-100 relative">
      <div
        v-if="loading"
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4"
      >
        <div
          v-for="i in 5"
          :key="i"
          class="aspect-2/3 rounded-2xl bg-zinc-200 dark:bg-zinc-800 animate-pulse"
        ></div>
      </div>

      <div
        v-else-if="!!userBookmarkNovels.length"
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-6"
      >
        <NovelaCard v-for="novela in userBookmarkNovels" :key="novela.id" v-bind="novela" />
      </div>

      <div
        v-else
        class="flex flex-col items-center justify-center py-20 bg-zinc-50 dark:bg-zinc-900/30 border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-3xl"
      >
        <Icon
          name="ph:bookmarks-simple-light"
          size="48"
          class="text-zinc-300 dark:text-zinc-700 mb-4"
        />
        <p class="text-zinc-500 font-medium">В этой категории пока пусто</p>
      </div>
    </div>
  </div>
</template>
