import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async () => {
  const authStore = useAuthStore()
  const token = useCookie('KYOKUSU_API_TOKEN')

  // HttpOnly cookies are readable only on the server; on the client the
  // session is restored from the SSR payload.
  if (import.meta.server) {
    if (token.value && !authStore.isAuthenticated) {
      await authStore.initAuth()
      if (!authStore.isAuthenticated) {
        token.value = null
      }
    }
    return
  }

  const initialized = useState('auth-initialized', () => false)
  if (initialized.value) return
  initialized.value = true

  if (!authStore.isAuthenticated) {
    await authStore.initAuth()
  }
})
