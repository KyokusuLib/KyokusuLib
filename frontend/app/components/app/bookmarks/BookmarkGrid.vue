<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import type { BookmarkCategory } from '@/types/frontend/bookmarks'
import NovelaCard from '@/components/app/novela/NovelaCard.vue'

const props = defineProps<{
  loading: boolean
  novels: any[]
  categories: { id: string; label: string }[]
}>()

const emit = defineEmits<{
  (e: 'change-status', id: number, category: BookmarkCategory | 'remove'): void
}>()

const openMenuId = ref<number | null>(null)

const closeMenu = () => {
  openMenuId.value = null
}

onMounted(() => {
  document.addEventListener('click', closeMenu)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', closeMenu)
})

const toggleMenu = (event: Event, id: number) => {
  event.stopPropagation()
  openMenuId.value = openMenuId.value === id ? null : id
}

const handleSelect = (id: number, category: any) => {
  emit('change-status', id, category)
  closeMenu()
}
</script>

<template>
  <div class="min-h-100 relative">
    <div
      v-if="loading"
      class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 sm:gap-6 pb-6"
    >
      <div
        v-for="i in 10"
        :key="i"
        class="w-full aspect-[2/3] rounded-3xl bg-zinc-200 dark:bg-zinc-800 animate-pulse"
      ></div>
    </div>

    <div
      v-else-if="novels.length"
      class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4 sm:gap-6 pb-6"
    >
      <div v-for="novela in novels" :key="novela.id" class="relative group">
        <div class="w-full h-full">
          <NovelaCard v-bind="novela" class="h-full" />
        </div>

        <div
          class="absolute top-2 right-2 transition-opacity z-10"
          :class="
            openMenuId === novela.id
              ? 'opacity-100'
              : 'opacity-100 sm:opacity-0 group-hover:opacity-100'
          "
        >
          <div class="relative">
            <button
              @click.prevent="(e) => toggleMenu(e, novela.id)"
              class="w-8 h-8 rounded-xl bg-zinc-900/80 backdrop-blur border border-white/10 flex items-center justify-center text-white shadow-lg transition-transform hover:scale-110 cursor-pointer"
            >
              <Icon name="ph:dots-three-bold" size="20" />
            </button>

            <Transition name="popover">
              <div
                v-if="openMenuId === novela.id"
                @click.stop
                class="absolute top-full -right-40 mt-1.4 w-48 p-1.5 z-50 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-xl shadow-xl overflow-hidden origin-top-right"
              >
                <div class="flex flex-col gap-0.5">
                  <button
                    v-for="cat in categories"
                    :key="cat.id"
                    @click="handleSelect(novela.id, cat.id)"
                    class="w-full py-2 rounded-lg text-sm font-medium flex items-center justify-center gap-2 cursor-pointer transition-colors"
                    :class="
                      novela.bookmark === cat.id
                        ? 'bg-yellow-50 dark:bg-yellow-500/10 text-yellow-600 dark:text-yellow-500'
                        : 'hover:bg-yellow-50 dark:hover:bg-yellow-500/10 text-zinc-600 dark:text-zinc-400 hover:text-yellow-600 dark:hover:text-yellow-500'
                    "
                  >
                    {{ cat.label }}
                    <div
                      v-if="novela.bookmark === cat.id"
                      class="w-1.5 h-1.5 rounded-full bg-yellow-500"
                    ></div>
                  </button>

                  <div class="h-px bg-zinc-100 dark:bg-zinc-800 my-1"></div>

                  <button
                    @click="handleSelect(novela.id, 'remove')"
                    class="w-full py-2 flex justify-center rounded-lg text-sm font-medium text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 transition-colors cursor-pointer"
                  >
                    Удалить из закладок
                  </button>
                </div>
              </div>
            </Transition>
          </div>
        </div>
      </div>
    </div>

    <div
      v-else
      class="flex flex-col items-center justify-center p-12 bg-white dark:bg-zinc-900/30 border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-[3rem]"
    >
      <Icon
        name="ph:bookmarks-simple-light"
        size="64"
        class="text-zinc-300 dark:text-zinc-700 mb-6"
      />
      <p class="text-zinc-400 font-bold uppercase tracking-[0.2em] text-xs text-center px-4">
        В этой категории пока пусто
      </p>
    </div>
  </div>
</template>

<style scoped>
.popover-enter-active,
.popover-leave-active {
  transition: all 0.2s ease;
}
.popover-enter-from,
.popover-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
