<script setup lang="ts">
import { Tooltip, Separator } from '@kyokusu-ui/vue'

defineProps<{
  totalCount: number
  selectedCount: number
  allSelected: boolean
  sortField: string
  sortOrder: 'asc' | 'desc'
}>()

const emit = defineEmits<{
  (e: 'update:allSelected', val: boolean): void
  (e: 'markAllRead'): void
  (e: 'deleteSelected'): void
  (e: 'update:sortField', val: string): void
  (e: 'update:sortOrder', val: 'asc' | 'desc'): void
}>()

const sortOptions = [
  { field: 'created_at', label: 'По дате', icon: 'ph:calendar-bold' },
  { field: 'title', label: 'По группе', icon: 'ph:folder-bold' },
]
</script>

<template>
  <div class="flex items-center justify-between gap-4 flex-wrap">
    <div class="flex items-center gap-2">
      <label
        class="flex items-center gap-2 cursor-pointer select-none"
        @click="emit('update:allSelected', !allSelected)"
      >
        <div
          class="w-5 h-5 rounded-md border-2 flex items-center justify-center transition-colors"
          :class="
            allSelected ? 'bg-yellow-600 border-yellow-600' : 'border-zinc-400 dark:border-zinc-600'
          "
        >
          <Icon v-if="allSelected" name="ph:check-bold" size="12" class="text-white" />
        </div>
        <span class="text-sm font-medium text-zinc-600 dark:text-zinc-400">
          {{ allSelected ? 'Снять выделение' : 'Выбрать всё' }}
        </span>
      </label>

      <span v-if="selectedCount" class="text-xs text-zinc-400 font-medium">
        ({{ selectedCount }})
      </span>
    </div>

    <div class="flex items-center gap-2">
      <Tooltip text="Отметить всё прочитанным">
        <button
          @click="emit('markAllRead')"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-semibold cursor-pointer border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-600 dark:text-zinc-400 hover:border-green-400/50 hover:text-green-600 dark:hover:text-green-400 transition-colors"
        >
          <Icon name="ph:check-circle-bold" size="16" />
          Прочитано
        </button>
      </Tooltip>

      <Tooltip v-if="selectedCount" text="Удалить выбранные">
        <button
          @click="emit('deleteSelected')"
          class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl text-xs font-semibold cursor-pointer border border-zinc-200 dark:border-zinc-700 bg-white dark:bg-zinc-900 text-zinc-600 dark:text-zinc-400 hover:border-red-400/50 hover:text-red-600 dark:hover:text-red-400 transition-colors"
        >
          <Icon name="ph:trash-bold" size="16" />
          Удалить
        </button>
      </Tooltip>

      <Separator orientation="vertical" style="height: 30px" />

      <div class="flex items-center gap-1">
        <button
          v-for="opt in sortOptions"
          :key="opt.field"
          @click="emit('update:sortField', opt.field)"
          class="flex items-center gap-1.5 px-2.5 py-1.5 rounded-xl text-xs font-semibold cursor-pointer transition-colors"
          :class="
            sortField === opt.field
              ? 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-500'
              : 'text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300'
          "
        >
          <Icon :name="opt.icon" size="14" />
          {{ opt.label }}
        </button>

        <button
          @click="emit('update:sortOrder', sortOrder === 'asc' ? 'desc' : 'asc')"
          class="p-1.5 rounded-xl cursor-pointer text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 transition-colors"
        >
          <Icon
            :name="sortOrder === 'asc' ? 'ph:sort-ascending' : 'ph:sort-descending'"
            size="16"
          />
        </button>
      </div>
    </div>
  </div>
</template>
