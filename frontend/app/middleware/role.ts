import { KyokusuAppRole, ROLE_WEIGHTS } from '~/types/enums/role-enum'

export default defineNuxtRouteMiddleware((to) => {
  const authStore = useAuthStore()

  const minRole = to.meta.minRole as KyokusuAppRole | undefined

  if (!minRole) return

  const userRole = (authStore.user?.role as KyokusuAppRole) || KyokusuAppRole.USER
  const userWeight = ROLE_WEIGHTS[userRole] || 0
  const requiredWeight = ROLE_WEIGHTS[minRole] || 0

  if (userWeight < requiredWeight) {
    return navigateTo('/')
  }
})
