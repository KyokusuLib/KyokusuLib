<script setup lang="ts">
import type { BackendNotification } from '@/types/backend/notification'
import { formatDate } from '~/utils/notification'
import { Tooltip } from '@kyokusu-ui/vue'

const props = defineProps<{
  notification: BackendNotification
  selected: boolean
}>()

const emit = defineEmits<{
  (e: 'update:selected', id: number): void
  (e: 'markRead', id: number): void
  (e: 'delete', id: number): void
}>()
</script>

<template>
  <div
    class="flex items-center gap-3 px-3 py-2.5 rounded-2xl transition-colors hover:bg-zinc-50 dark:hover:bg-zinc-800/40"
    :class="{ 'bg-zinc-50 dark:bg-zinc-800/30': selected }"
  >
    <label @click.stop="emit('update:selected', notification.id)" class="shrink-0 cursor-pointer">
      <div
        class="w-5 h-5 rounded-md border-2 flex items-center justify-center transition-colors"
        :class="
          selected ? 'bg-yellow-600 border-yellow-600' : 'border-zinc-400 dark:border-zinc-600'
        "
      >
        <Icon v-if="selected" name="ph:check-bold" size="12" class="text-white" />
      </div>
    </label>

    <div @click="emit('markRead', notification.id)" class="flex-1 min-w-0 cursor-pointer">
      <p
        class="text-sm leading-relaxed"
        :class="
          notification.isRead
            ? 'text-zinc-500 dark:text-zinc-400'
            : 'text-zinc-900 dark:text-white font-semibold'
        "
      >
        {{ notification.message || notification.title }}
      </p>
      <p class="text-[10px] text-zinc-400 font-semibold uppercase tracking-widest mt-0.5">
        {{ formatDate(notification.createdAt) }}
      </p>
    </div>

    <div
      class="flex items-center gap-1 shrink-0 opacity-0 group-hover:opacity-100 transition-opacity"
    >
      <Tooltip :text="notification.isRead ? 'Уже прочитано' : 'Отметить как прочитано'">
        <button
          v-if="!notification.isRead"
          @click.stop="emit('markRead', notification.id)"
          class="p-1.5 rounded-lg cursor-pointer text-zinc-400 hover:text-green-500 transition-colors"
        >
          <Icon name="ph:check-circle-bold" size="16" />
        </button>
      </Tooltip>
      <Tooltip text="Удалить">
        <button
          @click.stop="emit('delete', notification.id)"
          class="p-1.5 rounded-lg cursor-pointer text-zinc-400 hover:text-red-500 transition-colors"
        >
          <Icon name="ph:trash-bold" size="16" />
        </button>
      </Tooltip>
    </div>
  </div>
</template>
