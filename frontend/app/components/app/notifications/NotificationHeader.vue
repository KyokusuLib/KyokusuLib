<script setup lang="ts">
defineProps<{
  unreadCount: number
}>()
</script>

<template>
  <div class="flex items-center justify-between">
    <div>
      <div class="flex items-center justify-center gap-2">
        <Icon name="ph:bell-bold" size="18" class="mt-1" />
        <h1 class="text-2xl font-black text-zinc-900 dark:text-white">Уведомления</h1>
      </div>
      <div
        v-if="unreadCount"
        class="flex items-center gap-2 text-sm text-zinc-500 dark:text-zinc-400 mt-1"
      >
        <Icon name="ph:envelope-simple-bold" size="18" class="mt-1" />
        {{ unreadCount }} непрочитанных
      </div>
    </div>
  </div>
</template>
