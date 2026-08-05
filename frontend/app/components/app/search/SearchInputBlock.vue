<script setup lang="ts">
import { ref, onMounted } from 'vue'
import SearchCategories from '@/components/common/Search/SearchCategories.vue'
import type { SearchCategory } from '@/composables/api/search/useSearch'

const props = defineProps<{
  query: string
  activeCategory: SearchCategory
}>()

const emit = defineEmits(['update:query', 'update:activeCategory'])

const inputRef = ref<HTMLInputElement | null>(null)

onMounted(() => {
  inputRef.value?.focus()
})
</script>

<template>
  <div
    class="bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-3xl overflow-hidden shadow-sm mb-8"
  >
    <div class="group relative flex items-center px-4 sm:px-6 py-4 sm:py-6">
      <div
        class="flex items-center justify-center w-12 h-12 rounded-2xl bg-zinc-100 dark:bg-zinc-800/80 shadow-sm border border-zinc-200 dark:border-zinc-700 shrink-0 relative overflow-hidden group-focus-within:border-yellow-500/50 group-focus-within:ring-2 group-focus-within:ring-yellow-500/20 transition-all duration-300"
      >
        <div
          class="absolute inset-0 bg-yellow-500/10 opacity-0 group-focus-within:opacity-100 transition-opacity duration-300"
        ></div>
        <Icon
          name="ph:magnifying-glass-bold"
          size="24"
          class="text-zinc-400 dark:text-zinc-500 group-focus-within:text-yellow-500 transition-colors duration-300 relative z-10"
        />
      </div>

      <input
        ref="inputRef"
        :value="query"
        @input="emit('update:query', ($event.target as HTMLInputElement).value)"
        type="text"
        class="flex-1 w-full pl-5 pr-14 bg-transparent border-none text-2xl sm:text-3xl font-bold tracking-tight text-zinc-900 dark:text-zinc-100 placeholder-zinc-300 dark:placeholder-zinc-700 focus:outline-none focus:ring-0 selection:bg-yellow-500/30"
        placeholder="Что будем искать?"
        autofocus
      />

      <Transition name="fade-fast">
        <button
          v-if="query"
          @click="emit('update:query', '')"
          class="absolute right-4 sm:right-8 p-2 text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-zinc-100 bg-zinc-200/50 hover:bg-zinc-200 dark:bg-zinc-800/50 dark:hover:bg-zinc-700 rounded-xl transition-all cursor-pointer flex items-center justify-center gap-1 active:scale-95"
        >
          <span class="hidden sm:block text-[10px] font-bold uppercase tracking-widest ml-1"
            >Очистить</span
          >
          <Icon name="ph:x-bold" size="18" />
        </button>
      </Transition>
    </div>

    <div class="px-4 sm:px-6 pb-4 border-t border-zinc-100 dark:border-zinc-800/50 pt-4">
      <SearchCategories
        :model-value="activeCategory"
        @update:model-value="emit('update:activeCategory', $event)"
      />
    </div>
  </div>
</template>

<style scoped>
.fade-fast-enter-active,
.fade-fast-leave-active {
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}
.fade-fast-enter-from,
.fade-fast-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
