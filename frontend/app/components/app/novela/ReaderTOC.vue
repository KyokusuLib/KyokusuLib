<script setup lang="ts">
import { useNovela } from '~/composables/api/novela/useNovela'

const props = defineProps<{
  modelValue: boolean
  novelaId: number
  currentChapterId: string
  novelaTitle?: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const { novela, fetchNovela, isLoading } = useNovela()

const close = () => emit('update:modelValue', false)

watch(
  () => props.modelValue,
  (isOpen) => {
    if (isOpen && !novela.value) {
      fetchNovela(props.novelaId)
    }

    if (isOpen) {
      nextTick(() => {
        const el = document.querySelector(`[data-chapter-id="${props.currentChapterId}"]`)
        if (el) {
          el.scrollIntoView({ block: 'center', behavior: 'smooth' })
        }
      })
    }
  },
)

const sortedVolumes = computed(() => {
  if (!novela.value) return []
  return [...novela.value.volumes].sort((a, b) => a.number - b.number)
})
</script>

<template>
  <div>
    <!-- Backdrop -->
    <Transition name="fade">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-100 bg-black/20 dark:bg-black/40 backdrop-blur-[2px]"
        @click="close"
      ></div>
    </Transition>

    <!-- TOC Panel -->
    <Transition name="slide-left">
      <div
        v-if="modelValue"
        class="fixed left-0 inset-y-0 z-110 w-80 bg-white dark:bg-zinc-900 border-r border-zinc-200 dark:border-zinc-800 shadow-2xl flex flex-col"
      >
        <div
          class="p-6 border-b border-zinc-100 dark:border-zinc-800 flex items-center justify-between bg-zinc-50/50 dark:bg-zinc-900/50 backdrop-blur-sm"
        >
          <div class="flex flex-col gap-0.5 overflow-hidden">
            <h2 class="text-lg font-black tracking-tight truncate">
              {{ novela?.title || novelaTitle || 'Загрузка...' }}
            </h2>
            <span class="text-[10px] text-zinc-500 font-bold uppercase tracking-widest"
              >Оглавление</span
            >
          </div>
          <button
            @click="close"
            class="p-2 hover:bg-zinc-200 dark:hover:bg-zinc-800 rounded-xl transition-colors shrink-0"
          >
            <Icon name="ph:x-bold" size="20" />
          </button>
        </div>

        <div class="flex-1 overflow-y-auto custom-scrollbar p-4">
          <div v-if="isLoading" class="flex flex-col items-center justify-center py-20 gap-4">
            <div
              class="w-8 h-8 border-3 border-zinc-200 border-t-yellow-500 rounded-full animate-spin"
            ></div>
            <span class="text-xs text-zinc-500 font-medium">Синхронизация глав...</span>
          </div>

          <div v-else class="space-y-6">
            <div v-for="vol in sortedVolumes" :key="vol.id" class="space-y-2">
              <div class="flex items-center gap-3 px-2">
                <div class="h-px grow bg-zinc-100 dark:bg-zinc-800"></div>
                <span
                  class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-400 whitespace-nowrap"
                  >Том {{ vol.number }}</span
                >
                <div class="h-px grow bg-zinc-100 dark:bg-zinc-800"></div>
              </div>

              <div class="grid grid-cols-1 gap-1">
                <NuxtLink
                  v-for="ch in vol.chapters"
                  :key="ch.id"
                  :to="`/novela/reader/${props.novelaId}/${ch.id}`"
                  :data-chapter-id="ch.id"
                  class="flex items-center gap-3 px-2.5 py-1.5 rounded-xl transition-all group"
                  :class="
                    currentChapterId === ch.id
                      ? 'bg-yellow-500 text-white shadow-lg shadow-yellow-500/20'
                      : 'hover:bg-zinc-100 dark:hover:bg-zinc-800 text-zinc-600 dark:text-zinc-400'
                  "
                  @click="close"
                >
                  <span class="text-xs font-mono mt-1.5 opacity-50 shrink-0">#{{ ch.number }}</span>
                  <span class="text-sm font-bold truncate grow">{{ ch.title }}</span>
                  <Icon
                    v-if="currentChapterId === ch.id"
                    name="ph:play-fill"
                    size="14"
                    class="shrink-0"
                  />
                </NuxtLink>
              </div>
            </div>
          </div>
        </div>

        <div
          class="p-4 border-t border-zinc-100 dark:border-zinc-800 bg-zinc-50/50 dark:bg-zinc-900/50"
        >
          <NuxtLink
            v-if="novela"
            :to="`/novela/${novela.id}`"
            class="flex items-center justify-center gap-2 w-full py-2 bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 rounded-xl text-xs font-black uppercase tracking-widest hover:border-yellow-500 hover:text-yellow-500 transition-all active:scale-[0.98]"
          >
            <Icon name="ph:info-bold" size="16" />
            К описанию
          </NuxtLink>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.slide-left-enter-active,
.slide-left-leave-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-left-enter-from,
.slide-left-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #e4e4e7;
  border-radius: 20px;
}
:deep(.dark) .custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #27272a;
}
</style>
