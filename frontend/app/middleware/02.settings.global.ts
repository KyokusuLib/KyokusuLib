import { useInterfaceSettings } from '~/composables/api/settings/useInterfaceSettings'

export default defineNuxtRouteMiddleware(async (_to) => {
  const { syncSettingWithBackend } = useInterfaceSettings()
  const authStore = useAuthStore()

  if (import.meta.server && authStore.isAuthenticated) {
    await syncSettingWithBackend()
  }
})
