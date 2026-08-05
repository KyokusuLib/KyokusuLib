<script setup lang="ts">
import { Separator as Separtor } from '@kyokusu-ui/vue'

import AccountSettings from '@/components/app/settings/AccountSettings.vue'
import SafetySettings from '@/components/app/settings/SafetySettings.vue'
import { useSettingsStore } from '@/stores/settings'
import { TabSettings } from '@/types/enums/tab-enum'
import InterfaceSettings from './InterfaceSettings.vue'
import SocialNetworkSettings from './SocialNetworkSettings.vue'
import NotifySettings from './NotifySettings.vue'

const settingsStore = useSettingsStore()

const categories = [
  { id: TabSettings.Account, label: 'Аккаунт', icon: 'ph:user' },
  { id: TabSettings.Security, label: 'Безопасность', icon: 'ph:lock' },
  { id: TabSettings.Notifications, label: 'Уведомления', icon: 'ph:bell' },
  { id: TabSettings.Interface, label: 'Интерфейс', icon: 'ph:palette' },
  { id: TabSettings.SocialNetwork, label: 'Социальные сети', icon: 'ph:link-simple' },
]
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-radial-[at_center] dark:from-zinc-800 dark:to-zinc-950 dark:to-90% text-zinc-900 dark:text-zinc-200 transition-colors duration-300"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="mb-8">
        <div class="flex items-start sm:justify-start justify-center gap-2">
          <h1
            class="cs-text text-left text-4xl leading-xl font-semibold text-zinc-900 dark:text-white mb-4"
          >
            Настройки
          </h1>
          <Icon name="ph:gear-bold" class="mt-2" size="20" />
        </div>
        <Separtor />
      </div>

      <div class="flex flex-col md:flex-row gap-8">
        <aside class="w-full md:w-64 shrink-0">
          <nav class="flex flex-col space-y-1">
            <button
              v-for="category in categories"
              :key="category.id"
              @click="settingsStore.setActiveTab(category.id)"
              :class="[
                settingsStore.activeTab === category.id
                  ? 'bg-zinc-200 text-zinc-900 dark:bg-zinc-700/50 dark:text-white'
                  : 'text-zinc-500 hover:bg-zinc-100 hover:text-zinc-900 dark:text-stone-400 dark:hover:bg-zinc-800 dark:hover:text-zinc-200',
                'group grid grid-cols-[1fr_auto] items-center px-3 py-2.5 text-sm font-medium rounded-md transition-all duration-200 ease-in-out w-full text-left cursor-pointer gap-2',
              ]"
            >
              <span>{{ category.label }}</span>
              <Icon
                :name="category.icon"
                size="18"
                class="text-zinc-700 dark:text-zinc-400 justify-self-end"
              />
            </button>
          </nav>
        </aside>

        <div
          class="flex-1 bg-white border-zinc-200 dark:bg-zinc-900/30 rounded-lg border dark:border-zinc-800 p-6 min-h-100 transition-colors duration-300"
        >
          <div v-if="settingsStore.activeTab === TabSettings.Account" class="h-full">
            <AccountSettings />
          </div>

          <div v-else-if="settingsStore.activeTab === TabSettings.Security" class="space-y-6">
            <div class="h-full">
              <SafetySettings />
            </div>
          </div>

          <div v-else-if="settingsStore.activeTab === TabSettings.Notifications" class="space-y-6">
            <NotifySettings />
          </div>

          <div v-else-if="settingsStore.activeTab === TabSettings.Interface" class="space-y-6">
            <InterfaceSettings />
          </div>

          <div v-else-if="settingsStore.activeTab === TabSettings.SocialNetwork" class="space-y-6">
            <SocialNetworkSettings />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
