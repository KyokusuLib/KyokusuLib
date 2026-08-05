import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import type { GetUserDto } from '@/types/backend/user'
import { KyokusuAppRole } from '@/types/enums/role-enum'
import { useApi } from '~/composables/api/useApi'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<GetUserDto | null>(null)
  const isAuthChecking = ref(true)

  const isAuthenticated = computed(() => !!user.value)
  const isAdmin = computed(() => user.value?.role === KyokusuAppRole.ADMIN)

  async function initAuth() {
    isAuthChecking.value = true
    try {
      const { data, error } = await useApi<GetUserDto>('/api/auth/me')

      if (!error.value && data.value) {
        user.value = data.value
      } else {
        user.value = null
      }
    } catch {
      user.value = null
    } finally {
      isAuthChecking.value = false
    }
  }

  async function fetchUser() {
    return initAuth()
  }

  async function logout() {
    try {
      await useApi('/api/auth/logout', { method: 'POST' })
      user.value = null

      await navigateTo('/login')
    } catch (e) {
      console.error('Logout error:', e)
    }
  }

  return {
    user,
    isAuthChecking,
    isAuthenticated,
    isAdmin,
    initAuth,
    fetchUser,
    logout,
  }
})
