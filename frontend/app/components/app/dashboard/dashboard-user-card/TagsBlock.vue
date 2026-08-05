<script setup lang="ts">
import type { DashboardUserTag } from '@/types/frontend/dashboard/users'

defineProps<{
  tags: DashboardUserTag[]
  activeTag: string
}>()
</script>

<template>
  <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
    <div class="flex gap-2">
      <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 mb-2">Теги</p>
      <Icon name="ph:tag-bold" size="16" class="text-zinc-400 shrink-0" />
    </div>
    <div class="flex flex-wrap gap-1.5">
      <span
        v-for="tag in tags"
        :key="tag.tag_id"
        class="inline-flex px-2.5 py-1 rounded-lg text-[11px] font-medium bg-zinc-200/70 dark:bg-zinc-700/70 text-zinc-600 dark:text-zinc-300"
        :class="{ 'ring-2 ring-yellow-500/50': tag.tag === activeTag }"
      >
        {{ tag.tag }}
      </span>
    </div>
  </div>
</template>
