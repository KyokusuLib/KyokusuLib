<script setup lang="ts">
import type { DashboardUserSettings } from '@/types/frontend/dashboard/users'

defineProps<{ settings: DashboardUserSettings }>()
</script>

<template>
  <div class="flex items-center gap-4 text-xs text-zinc-400 border border-zinc-700 p-2 rounded-xl">
    <span class="inline-flex items-center gap-1">
      <Icon :name="settings.is_show_tag ? 'ph:eye-bold' : 'ph:eye-slash-bold'" size="14" />
      Теги {{ settings.is_show_tag ? 'видны' : 'скрыты' }}
    </span>
    <span class="inline-flex items-center gap-1">
      <Icon :name="settings.is_show_bookmark ? 'ph:eye-bold' : 'ph:eye-slash-bold'" size="14" />
      Закладки {{ settings.is_show_bookmark ? 'видны' : 'скрыты' }}
    </span>
  </div>
</template>
