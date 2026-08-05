<script setup lang="ts">
import type { BackendNotification } from '@/types/backend/notification'
import NotificationGroupCard from '~/components/app/notifications/NotificationGroupCard.vue'
import NotificationEmpty from '~/components/app/notifications/NotificationEmpty.vue'

defineProps<{
  loading: boolean
  items: BackendNotification[]
  groupEntries: [string, BackendNotification[]][]
  selectedIds: Set<number>
  expandedGroups: Set<string>
}>()

const emit = defineEmits<{
  (e: 'toggleExpand', title: string): void
  (e: 'update:selected', id: number): void
  (e: 'markRead', id: number): void
  (e: 'delete', id: number): void
  (e: 'toggleSelectGroup', title: string): void
  (e: 'groupMarkRead', title: string): void
  (e: 'groupDelete', title: string): void
}>()
</script>

<template>
  <div v-if="loading" class="py-24 flex flex-col items-center justify-center text-zinc-400">
    <Icon name="ph:spinner-bold" size="32" class="animate-spin mb-4" />
    <p class="text-xs font-bold uppercase tracking-[0.2em]">Загрузка...</p>
  </div>

  <NotificationEmpty v-else-if="!items.length" />

  <div v-else class="space-y-4">
    <NotificationGroupCard
      v-for="[title, notifications] in groupEntries"
      :key="title"
      :title="title"
      :notifications="notifications"
      :selected-ids="selectedIds"
      :expanded="expandedGroups.has(title)"
      @toggle-expand="emit('toggleExpand', $event)"
      @update:selected="emit('update:selected', $event)"
      @mark-read="emit('markRead', $event)"
      @delete="emit('delete', $event)"
      @toggle-select-group="emit('toggleSelectGroup', $event)"
      @group-mark-read="emit('groupMarkRead', $event)"
      @group-delete="emit('groupDelete', $event)"
    />
  </div>
</template>
