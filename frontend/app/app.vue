<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useActivityStore } from '@/stores/activity'
import { useNotificationStore } from '@/stores/notification'
import { useInterfaceSettings } from '~/composables/api/settings/useInterfaceSettings'
import { useNotifications } from '~/composables/api/notifications/useNotifications'
import { NotificationContainer as TheNotifications } from '@kyokusu-ui/vue'

const authStore = useAuthStore()
const activityStore = useActivityStore()
const notificationStore = useNotificationStore()

const { syncSettingWithBackend } = useInterfaceSettings()
const { connect: connectNotifications, disconnect, fetchStats } = useNotifications()

watch(
  () => authStore.isAuthenticated,
  async (auth) => {
    if (auth) {
      connectNotifications()
      await Promise.all([syncSettingWithBackend(), fetchStats()])
    } else {
      disconnect()
    }
  },
  { immediate: true },
)

onMounted(async () => {
  setTimeout(() => {
    const loader = document.getElementById('app-loading-style')
    if (loader) loader.remove()
  }, 100)

  activityStore.initActivityTracking()
})
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <ClientOnly>
      <TheNotifications
        :notifications="notificationStore.notifications"
        @remove="notificationStore.remove"
      />
    </ClientOnly>
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </div>
</template>
