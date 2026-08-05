import { computed } from 'vue'
import { useAuthStore } from '~/stores/auth'
import { GenderSetting } from '~/types/enums/gender-enum'
import { KyokusuAppRole } from '~/types/enums/role-enum'
import { TabProfile } from '@/types/enums/tab-profile'
import type { GetUserDto } from '~/types/backend/user'
import { useApi } from '../useApi'

export function useProfile() {
  const authStore = useAuthStore()
  const fetchedUser = useState<GetUserDto | null>('profile-fetched-user', () => null)
  const profileLoading = useState('profile-loading', () => false)
  const profileError = useState('profile-error', () => false)

  const isSelfProfile = computed(() => {
    if (!authStore.user || !fetchedUser.value) return false
    return authStore.user.id === fetchedUser.value.id
  })

  const profileData = computed<GetUserDto | null>(() => {
    return isSelfProfile.value ? authStore.user : fetchedUser.value
  })

  const isPublicAccount = computed(() => !!profileData.value?.is_public)
  const accountCreated = computed(() =>
    profileData.value?.create_at
      ? new Date(profileData.value.create_at).toLocaleDateString('ru-RU')
      : 'Неизвестно',
  )

  const profileTabs = computed(() => {
    const tabs: { label: string; id: TabProfile }[] = [{ label: 'Обзор', id: TabProfile.Overview }]

    if (
      (authStore.isAuthenticated && isSelfProfile.value) ||
      profileData.value?.settings.is_show_bookmark
    ) {
      tabs.push({ label: 'Закладки', id: TabProfile.Bookmarks })
    }

    tabs.push({ label: 'Команды', id: TabProfile.Teams })
    tabs.push({ label: 'Комментарии', id: TabProfile.Comments })

    return tabs
  })

  const getRoleColor = (role?: string) => {
    const roles: Record<string, string> = {
      [KyokusuAppRole.ADMIN]: 'text-red-400 bg-red-400/10 border-red-400/20',
      [KyokusuAppRole.MODERATOR]: 'text-blue-400 bg-blue-400/10 border-blue-400/20',
    }
    return roles[role || ''] || 'text-zinc-400 bg-zinc-400/10 border-zinc-400/20'
  }

  const getGenderText = (gender?: string) => {
    if (gender === GenderSetting.MALE) return 'Мужской'
    if (gender === GenderSetting.FEMALE) return 'Женский'
    return 'Не указан'
  }

  const checkIsLogin = (last_login?: string | null) => {
    if (!last_login) return false
    return Date.now() - new Date(last_login).getTime() < 300000
  }

  const formatLastLogin = (last_login?: string | null) => {
    if (!last_login) return 'Неизвестно'
    const date = new Date(last_login)
    const isToday = date.toDateString() === new Date().toDateString()
    return isToday
      ? date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' })
      : date.toLocaleDateString('ru-RU')
  }

  const init = async (id: number) => {
    if (!id) return null

    if (authStore.user?.id === id) {
      fetchedUser.value = authStore.user
      return authStore.user
    }

    profileLoading.value = true
    try {
      const { data, error } = await useApi<GetUserDto>(`/api/user/${id}`)
      if (error.value || !data.value) {
        fetchedUser.value = null
        profileError.value = true
        return null
      }
      fetchedUser.value = data.value
      return data.value
    } catch {
      fetchedUser.value = null
      return null
    } finally {
      profileLoading.value = false
    }
  }

  return {
    profileData,
    isSelfProfile,
    isPublicAccount,
    isLoading: profileLoading,
    profileTabs,
    accountCreated: accountCreated,
    userRoleColor: computed(() => getRoleColor(profileData.value?.role)),
    userGender: computed(() => getGenderText(profileData.value?.gender)),
    lastLogin: computed(() => formatLastLogin(profileData.value?.last_login)),
    isLogin: computed(() => profileData.value?.status === 'online'),

    getRoleColor,
    getGenderText,
    checkIsLogin,
    formatLastLogin,
    profileError,
    init,
  }
}
