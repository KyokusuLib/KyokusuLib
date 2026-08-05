<script lang="ts" setup>
import { useInterfaceSettings } from '@/composables/api/settings/useInterfaceSettings'

import TabToggleSettings from '@/components/app/settings/TabToggleSettings.vue'

const { isDarkTheme, isShowTag, isShowBookmark, syncSettingWithBackend } = useInterfaceSettings()

await useAsyncData('settings-sync', () => syncSettingWithBackend())
</script>

<template>
  <div class="space-y-8">
    <div class="flex flex-col items-center">
      <h2 class="text-3xl font-bold mb-1">Интерфейс</h2>
      <p class="text-sm text-zinc-500 dark:text-stone-400">Настройте внешний вид приложения.</p>
    </div>

    <div class="space-y-6 text-left">
      <div class="pb-6 border-b space-y-4 border-zinc-200 dark:border-zinc-800">
        <h3 class="px-2 text-base font-semibold">Оформление</h3>

        <TabToggleSettings
          v-model="isDarkTheme"
          id="theme_toggle"
          title="Тёмная тема"
          imgSrc="ph:palette"
          description="Переключите для смены темы."
        />
      </div>

      <div class="pb-6 border-b space-y-4 border-zinc-200 dark:border-zinc-800">
        <h3 class="px-2 text-base font-semibold">Профиль</h3>

        <TabToggleSettings
          v-model="isShowTag"
          id="show_tag_toggle"
          title="Показывать тэг пользователя"
          imgSrc="ph:tag"
          description="Переключите для показа пользовательского тэга."
        />

        <TabToggleSettings
          v-model="isShowBookmark"
          id="show_bookmark_toggle"
          title="Показывать закладки пользователя"
          imgSrc="ph:bookmark"
          description="Переключите для отображения закладок пользователя в профиле."
        />
      </div>
    </div>
  </div>
</template>
