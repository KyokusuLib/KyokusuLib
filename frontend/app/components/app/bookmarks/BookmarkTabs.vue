<script setup lang="ts">
import { ref } from 'vue'
import type { BookmarkCategory } from '@/types/frontend/bookmarks'

const props = defineProps<{
  categories: { id: string; label: string; icon: string; count?: number }[]
  activeCategory: BookmarkCategory
}>()

const emit = defineEmits<{
  (e: 'update:activeCategory', val: BookmarkCategory): void
}>()

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
  emit('update:activeCategory', catId)
}
</script>

<template>
  <div
    ref="scrollContainer"
    class="flex gap-2 overflow-x-auto pb-4 no-scrollbar cursor-grab active:cursor-grabbing select-none"
    @mousedown="onMouseDown"
    @mouseleave="onMouseLeave"
    @mouseup="onMouseUp"
    @mousemove="onMouseMove"
  >
    <button
      v-for="cat in categories"
      :key="cat.id"
      @click="(e) => handleTabClick(e, cat.id as BookmarkCategory)"
      class="px-5 py-1.5 rounded-2xl text-xs font-black uppercase tracking-widest transition-all whitespace-nowrap border cursor-pointer"
      :class="[
        activeCategory === cat.id
          ? 'bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 border-transparent shadow-xl shadow-zinc-900/10 dark:shadow-white/10'
          : 'bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 text-zinc-500 hover:border-zinc-400',
      ]"
    >
      <div class="flex justify-center items-center gap-3">
        <Icon :name="cat.icon" size="18" />
        <div class="flex items-center justify-center gap-2">
          <span>{{ cat.label }}</span>
          <span
            v-if="cat.count !== undefined"
            class="flex items-center justify-center border border-zinc-500 bg-zinc-200 dark:bg-zinc-800 text-zinc-500 px-1.5 py-0.5 h-4 w-4 text-[10px] rounded-full"
            :class="
              activeCategory === cat.id
                ? 'bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-zinc-200'
                : ''
            "
          >
            <div class="ml-[1px]">
              {{ cat.count }}
            </div>
          </span>
        </div>
      </div>
    </button>
  </div>
</template>
