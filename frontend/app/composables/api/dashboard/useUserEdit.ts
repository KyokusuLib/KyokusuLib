import { reactive, computed } from 'vue'
import { useUserApi } from '@/composables/api/user/userApi'
import { useNotificationStore } from '@/stores/notification'
import { DashboardRowUserStatus, DashboardRowUserGender } from '@/types/enums/dashboard-table'
import { KyokusuAppRole } from '@/types/enums/role-enum'
import { formatRole, formatStatus, formatGender } from '@/utils/dashboard'
import type { DashboardUser } from '@/types/frontend/dashboard/users'

export interface UserEditForm {
  name: string
  about: string
  gender: string
  birthday: string
  isPublic: boolean
  role: string
  status: string
  picture: string
  banner: string
  isShowTag: boolean
  isShowBookmark: boolean
}

export const GENDER_OPTIONS = Object.values(DashboardRowUserGender).map((value) => ({
  value,
  label: formatGender(value),
}))

export const STATUS_OPTIONS = Object.values(DashboardRowUserStatus).map((value) => ({
  value,
  label: formatStatus(value),
}))

export const ROLE_OPTIONS = Object.values(KyokusuAppRole).map((value) => ({
  value,
  label: formatRole(value),
}))

export function useUserEdit() {
  const { updateUser } = useUserApi()
  const { notify } = useNotificationStore()

  const form = reactive<UserEditForm>({
    name: '',
    about: '',
    gender: 'hidden',
    birthday: '',
    isPublic: true,
    role: 'user',
    status: 'offline',
    picture: '',
    banner: '',
    isShowTag: false,
    isShowBookmark: false,
  })

  const isDirty = computed(() => {
    return form.name !== '' || form.about !== ''
  })

  function loadUser(user: DashboardUser) {
    form.name = user.name
    form.about = user.about ?? ''
    form.gender = user.gender ?? 'hidden'
    form.isPublic = user.is_public
    form.role = user.role
    form.status = user.status
    form.picture = user.picture
    form.banner = user.banner
    form.isShowTag = user.settings?.is_show_tag ?? false
    form.isShowBookmark = user.settings?.is_show_bookmark ?? false

    // clean birthday for DatePicker — same logic as AccountSettings
    let bd = user.birthday ?? ''
    if (
      !bd ||
      bd === '01.01.1' ||
      bd === '0001-01-01T00:00:00Z' ||
      bd === '0001-01-01' ||
      bd.startsWith('0001-01-01')
    ) {
      bd = ''
    } else {
      if (bd.includes('T')) bd = bd.split('T')[0]!
      if (bd.includes(' ')) bd = bd.split(' ')[0]!
    }
    form.birthday = bd
  }

  function reset() {
    form.name = ''
    form.about = ''
    form.gender = 'hidden'
    form.birthday = ''
    form.isPublic = true
    form.role = 'user'
    form.status = 'offline'
    form.picture = ''
    form.banner = ''
    form.isShowTag = false
    form.isShowBookmark = false
  }

  async function save(userId: number): Promise<boolean> {
    try {
      const ok = await updateUser(userId, {
        name: form.name,
        about: form.about,
        gender: form.gender,
        birthday: form.birthday,
        is_public: form.isPublic,
        is_show_tag: form.isShowTag,
        is_show_bookmark: form.isShowBookmark,
        role: form.role,
        status: form.status,
      })

      if (!ok) {
        notify({
          title: 'Ошибка',
          content: 'Не удалось обновить пользователя',
          type: 'error',
        })
        return false
      }

      return true
    } catch (e: any) {
      notify({
        title: 'Ошибка',
        content: e?.message ?? 'Не удалось обновить пользователя',
        type: 'error',
      })
      return false
    }
  }

  return {
    form,
    isDirty,
    GENDER_OPTIONS,
    STATUS_OPTIONS,
    ROLE_OPTIONS,
    loadUser,
    reset,
    save,
  }
}
