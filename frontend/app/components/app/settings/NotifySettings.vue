<script lang="ts" setup>
import { useNotifySettings } from '@/composables/api/settings/useNotifySettings'

import TabToggleSettings from './TabToggleSettings.vue'

const { profileSettings } = useNotifySettings()
</script>

<template>
  <div class="space-y-8 cursor-default transition-colors duration-300">
    <div class="flex flex-col items-center">
      <h2 class="text-3xl font-bold mb-1 text-zinc-900 dark:text-white">Уведомления</h2>
      <p class="text-sm text-zinc-500 dark:text-stone-400">Настройте уведомления приложения.</p>
    </div>

    <div class="space-y-6 text-left">
      <div class="pb-6 border-b space-y-4 border-zinc-200 dark:border-zinc-800">
        <h3 class="px-2 text-base leading-xl font-semibold text-zinc-800 dark:text-zinc-200">
          Уведомления на сайте
        </h3>

        <TabToggleSettings
          v-model="profileSettings.is_app_notify"
          id="app_notify_toggle"
          title="Оповещения"
          imgSrc="ph:envelope"
          description="Переключите для отключения оповещений на сайте."
        />

        <TabToggleSettings
          v-model="profileSettings.is_new_published_notify"
          id="show_tag_toggle"
          title="Выход новых работ"
          imgSrc="ph:chat-text"
          description="Переключите для отключения оповещений на сайте."
        />
      </div>
    </div>
  </div>
</template>
