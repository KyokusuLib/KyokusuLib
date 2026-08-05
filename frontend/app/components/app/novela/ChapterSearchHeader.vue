<script setup lang="ts">
import { Tooltip } from '@kyokusu-ui/vue'

defineProps<{
  modelValue: string
  canManage?: boolean
}>()

const route = useRoute()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'add-volume': []
}>()
</script>

<template>
  <div class="flex items-center gap-2 w-full sm:w-auto">
    <Tooltip text="Добавить том">
      <button
        v-if="canManage"
        @click="emit('add-volume')"
        class="flex items-center justify-center rounded-2xl h-8 w-8 p-2 text-sm font-medium whitespace-nowrap transition-colors border bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 hover:border-zinc-400 hover:text-yellow-500 cursor-pointer"
      >
        <Icon name="ph:plus-bold" size="16" />
      </button>
    </Tooltip>
    <Tooltip text="Добавить главу">
      <NuxtLink
        v-if="canManage"
        :to="`/novela/${route.params.id}/add-chapter`"
        class="flex items-center justify-center rounded-2xl h-8 w-8 p-2 text-sm font-medium whitespace-nowrap transition-colors border bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 hover:border-zinc-400 hover:text-yellow-500 cursor-pointer"
      >
        <Icon name="ph:file-plus-bold" size="16" />
      </NuxtLink>
    </Tooltip>
    <div class="relative w-full sm:w-64">
      <input
        :value="modelValue"
        type="text"
        placeholder="Поиск главы..."
        class="w-full pl-9 pr-4 py-2 bg-zinc-100 dark:bg-zinc-800 border-none rounded-lg text-sm text-zinc-900 dark:text-zinc-200 placeholder-zinc-400 focus:ring-2 focus:ring-zinc-500 outline-none transition-all"
        @input="emit('update:modelValue', ($event.target as HTMLInputElement).value)"
      />
      <Icon
        name="ph:magnifying-glass-bold"
        size="16"
        class="absolute left-3 top-2.5 text-zinc-400"
      />
    </div>
  </div>
</template>
