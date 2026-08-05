<script setup lang="ts">
withDefaults(
  defineProps<{
    placeholder?: string
  }>(),
  { placeholder: 'Поиск...' },
)

const model = defineModel<string>({ required: true })
</script>

<template>
  <div class="relative">
    <Icon
      name="ph:magnifying-glass-bold"
      size="16"
      class="absolute left-4 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none"
    />
    <input
      v-model="model"
      type="text"
      :placeholder="placeholder"
      class="w-full bg-zinc-100 dark:bg-zinc-800 border-none rounded-2xl pl-11 pr-4 py-2.5 text-sm focus:ring-2 focus:ring-yellow-500/50 outline-none transition-all placeholder:text-zinc-400"
    />
  </div>
</template>
