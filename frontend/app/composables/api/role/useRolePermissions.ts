import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { KyokusuAppRole, ROLE_WEIGHTS } from '@/types/enums/role-enum'

export function useRolePermissions() {
  const authStore = useAuthStore()

  const currentRole = computed((): KyokusuAppRole => {
    return (authStore.user?.role as KyokusuAppRole) || KyokusuAppRole.USER
  })

  const currentWeight = computed(() => ROLE_WEIGHTS[currentRole.value] || 0)

  const getWeight = (role: string): number => {
    return ROLE_WEIGHTS[role as KyokusuAppRole] || 0
  }

  const hasPermission = (requiredRole: KyokusuAppRole): boolean => {
    const requiredWeight = ROLE_WEIGHTS[requiredRole]
    return currentWeight.value >= requiredWeight
  }

  const isExactRole = (role: KyokusuAppRole): boolean => {
    return currentRole.value === role
  }

  const canManageUser = (targetUserRole: string): boolean => {
    const targetWeight = getWeight(targetUserRole)
    return currentWeight.value > targetWeight
  }

  return {
    currentRole,
    currentWeight,
    hasPermission,
    isExactRole,
    canManageUser,
  }
}
