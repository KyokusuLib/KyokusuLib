import { useAuthStore } from '~/stores/auth'

export default defineNuxtRouteMiddleware(async (_to, _from) => {
  const authStore = useAuthStore()

  const token = useCookie('KYOKUSU_API_TOKEN')

  if (token.value && !authStore.isAuthenticated) {
    try {
      await authStore.initAuth()
    } catch {
      token.value = null
    }
  }
})
