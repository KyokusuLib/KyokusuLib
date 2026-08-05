<script setup lang="ts">
defineProps<{
  query: string
}>()
</script>

<template>
  <div
    class="flex flex-col items-center justify-center py-32 text-center bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-3xl"
  >
    <div
      class="w-24 h-24 rounded-full bg-zinc-50 dark:bg-zinc-800/50 flex items-center justify-center mb-6 border border-zinc-100 dark:border-zinc-800"
    >
      <Icon name="ph:magnifying-glass-bold" size="40" class="text-zinc-300 dark:text-zinc-600" />
    </div>
    <h3 class="text-xl font-bold text-zinc-900 dark:text-zinc-100 mb-2">Ничего не найдено</h3>
    <p class="text-zinc-500 max-w-sm mx-auto">
      По запросу «<span class="font-semibold text-zinc-700 dark:text-zinc-300">{{ query }}</span
      >» нет результатов. Попробуйте изменить запрос или выбрать другую категорию.
    </p>
  </div>
</template>
