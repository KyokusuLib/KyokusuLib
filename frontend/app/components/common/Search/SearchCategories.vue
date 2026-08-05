<script setup lang="ts">
import type { SearchCategory } from '@/composables/api/search/useSearch'

defineProps<{
  modelValue: SearchCategory
}>()

const emit = defineEmits(['update:modelValue'])

const categories: { id: SearchCategory; label: string; icon: string }[] = [
  { id: 'ranobe', label: 'Ранобе', icon: 'ph:book-open-bold' },
  { id: 'users', label: 'Пользователи', icon: 'ph:users-bold' },
  { id: 'teams', label: 'Команды', icon: 'ph:users-three-bold' },
  { id: 'authors', label: 'Авторы', icon: 'ph:pen-nib-bold' },
]
</script>

<template>
  <div
    class="flex sm:justify-center justify-start overflow-x-auto gap-2 px-2 sm:px-4 py-3 no-scrollbar"
  >
    <button
      v-for="cat in categories"
      :key="cat.id"
      @click="emit('update:modelValue', cat.id)"
      class="flex items-center gap-2 px-4 py-2 rounded-xl text-[11px] font-black uppercase tracking-widest whitespace-nowrap transition-all duration-300 cursor-pointer"
      :class="
        modelValue === cat.id
          ? 'bg-yellow-50 dark:bg-yellow-500/10 text-yellow-600 dark:text-yellow-500 ring-1 ring-yellow-500/30 shadow-sm'
          : 'bg-zinc-100 text-zinc-500 hover:bg-zinc-200 dark:bg-zinc-800/80 dark:text-zinc-400 dark:hover:bg-zinc-700'
      "
    >
      <Icon :name="cat.icon" class="mb-0.5" size="16" />
      {{ cat.label }}
    </button>
  </div>
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
