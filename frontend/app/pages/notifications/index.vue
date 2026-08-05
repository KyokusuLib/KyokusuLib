<script setup lang="ts">
import { watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import { useNotificationManager } from '~/composables/notifications/useNotificationManager'

definePageMeta({ middleware: ['auth'] })

const { isAuthenticated } = storeToRefs(useAuthStore())

const {
  isLoading,
  selectedIds,
  expandedGroups,
  sortField,
  sortOrder,
  list,
  unreadCount,
  groupEntries,
  allSelected,
  selectedCount,
  load,
  toggleSelectAll,
  toggleSelect,
  toggleExpand,
  handleMarkRead,
  handleDelete,
  handleMarkAllRead,
  handleDeleteSelected,
  handleGroupMarkRead,
  handleGroupDelete,
  toggleSelectGroup,
  toggleSortField,
  toggleSortOrder,
} = useNotificationManager()

watch(
  isAuthenticated,
  (auth) => {
    if (auth) load()
  },
  { immediate: true },
)
</script>

<template>
  <div class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] transition-colors duration-300">
    <div class="max-w-3xl mx-auto px-4 py-8 space-y-6">
      <ClientOnly>
        <AppNotificationsNotificationHeader :unread-count="unreadCount" />

        <AppNotificationsNotificationToolbar
          v-if="list.length"
          :total-count="list.length"
          :selected-count="selectedCount"
          :all-selected="allSelected"
          :sort-field="sortField"
          :sort-order="sortOrder"
          @update:all-selected="toggleSelectAll"
          @mark-all-read="handleMarkAllRead"
          @delete-selected="handleDeleteSelected"
          @update:sort-field="toggleSortField($event)"
          @update:sort-order="toggleSortOrder"
        />

        <AppNotificationsNotificationList
          :loading="isLoading"
          :items="list"
          :group-entries="groupEntries"
          :selected-ids="selectedIds"
          :expanded-groups="expandedGroups"
          @toggle-expand="toggleExpand"
          @update:selected="toggleSelect"
          @mark-read="handleMarkRead"
          @delete="handleDelete"
          @toggle-select-group="toggleSelectGroup"
          @group-mark-read="handleGroupMarkRead"
          @group-delete="handleGroupDelete"
        />
      </ClientOnly>
    </div>
  </div>
</template>
