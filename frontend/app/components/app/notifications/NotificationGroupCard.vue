<script setup lang="ts">
import { computed } from 'vue'
import { Tooltip } from '@kyokusu-ui/vue'
import type { BackendNotification } from '@/types/backend/notification'
import { unreadCountIn, pluralize } from '~/utils/notification'
import NotificationItem from '~/components/app/notifications/NotificationItem.vue'

const props = defineProps<{
  title: string
  notifications: BackendNotification[]
  selectedIds: Set<number>
  expanded: boolean
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

const groupAllSelected = computed(
  () =>
    props.notifications.length > 0 && props.notifications.every((n) => props.selectedIds.has(n.id)),
)

const groupHasUnread = computed(() => unreadCountIn(props.notifications) > 0)
</script>

<template>
  <div
    class="group bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 rounded-3xl transition-all hover:border-yellow-500/30"
  >
    <div class="flex items-center justify-between px-5 py-4">
      <div class="flex items-center gap-3 min-w-0">
        <label @click.stop="emit('toggleSelectGroup', title)" class="shrink-0 cursor-pointer">
          <div
            class="w-5 h-5 rounded-md border-2 flex items-center justify-center transition-colors"
            :class="
              groupAllSelected
                ? 'bg-yellow-600 border-yellow-600'
                : 'border-zinc-400 dark:border-zinc-600'
            "
          >
            <Icon v-if="groupAllSelected" name="ph:check-bold" size="12" class="text-white" />
          </div>
        </label>

        <button
          @click="emit('toggleExpand', title)"
          class="flex items-center gap-3 min-w-0 cursor-pointer"
        >
          <div
            class="w-10 h-10 rounded-xl bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center shrink-0"
          >
            <Icon name="ph:bell-bold" size="20" class="text-zinc-500 dark:text-zinc-400" />
          </div>
          <div class="min-w-0 text-left">
            <p class="text-sm font-bold text-zinc-900 dark:text-white truncate">
              {{ title }}
            </p>
            <p class="text-[10px] font-semibold text-zinc-400 uppercase tracking-widest">
              {{ notifications.length }}
              {{ pluralize(notifications.length, ['уведомление', 'уведомлений']) }}
              <template v-if="unreadCountIn(notifications)">
                · {{ unreadCountIn(notifications) }} непрочитано
              </template>
            </p>
          </div>
        </button>
      </div>

      <div class="flex items-center gap-1 shrink-0">
        <div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
          <Tooltip v-if="groupHasUnread" text="Отметить всё прочитанным">
            <button
              @click.stop="emit('groupMarkRead', title)"
              class="p-1.5 rounded-lg cursor-pointer text-zinc-400 hover:text-green-500 hover:bg-green-50 dark:hover:bg-green-500/10 transition-colors"
            >
              <Icon name="ph:check-circle-bold" size="16" />
            </button>
          </Tooltip>
          <Tooltip text="Удалить все">
            <button
              @click.stop="emit('groupDelete', title)"
              class="p-1.5 rounded-lg cursor-pointer text-zinc-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-500/10 transition-colors"
            >
              <Icon name="ph:trash-bold" size="16" />
            </button>
          </Tooltip>
        </div>

        <button @click="emit('toggleExpand', title)" class="cursor-pointer flex items-center p-1">
          <Icon
            name="ph:caret-down-bold"
            size="12"
            class="text-zinc-400 transition-transform duration-300"
            :class="{ 'rotate-180': expanded }"
          />
        </button>
      </div>
    </div>

    <Transition name="collapse">
      <div v-if="expanded" class="px-5 pb-4 space-y-2">
        <NotificationItem
          v-for="n in [...notifications].sort(
            (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
          )"
          :key="n.id"
          :notification="n"
          :selected="selectedIds.has(n.id)"
          @update:selected="emit('update:selected', $event)"
          @mark-read="emit('markRead', $event)"
          @delete="emit('delete', $event)"
        />
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.collapse-enter-active,
.collapse-leave-active {
  transition: all 0.3s cubic-bezier(0.23, 1, 0.32, 1);
  overflow: hidden;
}
.collapse-enter-from,
.collapse-leave-to {
  opacity: 0;
  max-height: 0;
  padding-top: 0;
  padding-bottom: 0;
}
.collapse-enter-to,
.collapse-leave-from {
  max-height: 1000px;
}
</style>
