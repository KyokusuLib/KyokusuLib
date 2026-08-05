import { ref } from 'vue'
import { parseDateToISO } from '@/utils/date'
import PersonIcon from '@/assets/images/special/user.png'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import type { GetUserDto } from '@/types/backend/user'
import { useApi } from '~/composables/api/useApi'

export function useProfileSettings() {
  const authStore = useAuthStore()

  const { user } = authStore
  const { notify } = useNotificationStore()

  const errors = ref<Record<string, string>>({})
  const isLoading = ref(false)
  const isSuccess = ref(false)
  const isBirthDateDisable = ref(
    !!(user?.birthday && !user.birthday.startsWith('0001-01-01') && user.birthday !== '01.01.1'),
  )

  const fileInput = ref<HTMLInputElement | null>(null)
  const selectedFile = ref<File | null>(null)

  const bannerInput = ref<HTMLInputElement | null>(null)
  const selectedBannerFile = ref<File | null>(null)

  const profile = ref<GetUserDto>({ ...user! })

  if (
    !profile.value.birthday ||
    profile.value.birthday === '01.01.1' ||
    profile.value.birthday === '0001-01-01T00:00:00Z' ||
    profile.value.birthday === '0001-01-01' ||
    profile.value.birthday?.startsWith('0001-01-01')
  ) {
    profile.value.birthday = ''
  } else {
    // Clean up time components
    let cleanDate = profile.value.birthday
    if (cleanDate.includes('T')) cleanDate = cleanDate.split('T')[0]!
    if (cleanDate.includes(' ')) cleanDate = cleanDate.split(' ')[0]!

    // Validate that the date can be safely parsed by kyokusu-ui-vue
    if (cleanDate.includes('-')) {
      const parts = cleanDate.split('-')
      if (parts.length !== 3 || parts.some((p) => isNaN(Number(p)))) {
        cleanDate = ''
      }
    } else if (cleanDate.includes('.')) {
      const parts = cleanDate.split('.')
      if (parts.length !== 3 || parts.some((p) => isNaN(Number(p)))) {
        cleanDate = ''
      }
    }
    profile.value.birthday = cleanDate
  }

  const validate = (): boolean => {
    errors.value = {}
    let isValid = true

    const nick = profile.value.name?.trim() || ''
    const allowedCharsRegex = /^[A-Za-z\s]{6,50}$/
    const hasLowercase = /[a-z]/
    const hasUppercase = /[A-Z]/

    if (!nick) {
      errors.value.name = 'Никнейм обязателен.'
      isValid = false
    } else if (!allowedCharsRegex.test(nick)) {
      if (nick.length < 6 || nick.length > 50) {
        errors.value.name = 'Длина никнейма от 6 до 50 символов.'
      } else {
        errors.value.name = 'Никнейм может содержать только латинские буквы (A–Z, a–z) и пробелы.'
      }
      isValid = false
    } else if (!hasLowercase.test(nick) || !hasUppercase.test(nick)) {
      errors.value.name =
        'Никнейм должен содержать минимум 1 строчную (a–z) и 1 заглавную (A–Z) букву.'
      isValid = false
    }

    const aboutText = profile.value.about || ''
    if (aboutText.length > 500) {
      errors.value.about = `Превышен лимит символов (${aboutText.length}/500).`
      isValid = false
    }

    const birthdayText = profile.value.birthday
    if (
      birthdayText === '01.01.1' ||
      birthdayText === '0001-01-01T00:00:00Z' ||
      birthdayText === '0001-01-01'
    ) {
      profile.value.birthday = ''
    }

    return isValid
  }

  const handleAvatarClick = () => {
    fileInput.value?.click()
  }

  const handleFileChange = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files[0]) {
      const file = target.files[0]
      selectedFile.value = file

      const reader = new FileReader()
      reader.onload = (e) => {
        if (e.target?.result) {
          profile.value.picture = e.target.result as string
        }
      }
      reader.readAsDataURL(file)
    }
  }

  const removeAvatar = () => {
    profile.value.picture = PersonIcon
    selectedFile.value = null
    if (fileInput.value) fileInput.value.value = ''
  }

  const handleBannerClick = () => {
    bannerInput.value?.click()
  }

  const handleBannerChange = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files[0]) {
      const file = target.files[0]
      selectedBannerFile.value = file

      const reader = new FileReader()
      reader.onload = (e) => {
        if (e.target?.result) {
          profile.value.banner = e.target.result as string
        }
      }
      reader.readAsDataURL(file)
    }
  }

  const removeBanner = () => {
    profile.value.banner = ''
    selectedBannerFile.value = null
    if (bannerInput.value) bannerInput.value.value = ''
  }

  const saveChanges = async () => {
    if (!validate()) {
      return
    }

    isLoading.value = true
    isSuccess.value = false

    try {
      const updatePayload = {
        name: profile.value.name,
        about: profile.value.about,
        gender: profile.value.gender,
        birthday: parseDateToISO(profile.value.birthday),
        is_public: profile.value.is_public,
      }

      const { error: profileError } = await useApi<any>('/api/profile/update', {
        method: 'PUT',
        body: updatePayload,
      })

      if (profileError.value) {
        notify({
          title: 'Неудача',
          content: 'Что-то пошло не так :(',
          type: 'error',
        })
        throw new Error(profileError.value.data?.error || 'Unknown error')
      }

      if (selectedFile.value) {
        const formData = new FormData()
        formData.append('avatar', selectedFile.value)

        const { data: avatarData, error: avatarError } = await useApi<any>('/api/profile/avatar', {
          method: 'POST',
          body: formData,
        })

        if (avatarError.value) {
          notify({
            title: 'Неудача',
            content: 'Неудалось сохранить аватар',
            type: 'error',
          })
          throw new Error('Failed to upload avatar')
        }

        if (avatarData.value && avatarData.value.picture) {
          profile.value.picture = avatarData.value.picture
        }
      } else if (selectedFile.value === null && profile.value.picture === PersonIcon) {
        const { error: avatarError } = await useApi('/api/profile/avatar', {
          method: 'DELETE',
        })
        if (avatarError.value) {
          notify({
            title: 'Неудача',
            content: 'Неудалось удалить аватар',
            type: 'error',
          })
          throw new Error('Failed to remove avatar')
        }
      }

      if (selectedBannerFile.value) {
        const formData = new FormData()
        formData.append('banner', selectedBannerFile.value)

        const { data: bannerData, error: bannerError } = await useApi<any>('/api/profile/banner', {
          method: 'POST',
          body: formData,
        })

        if (bannerError.value) {
          notify({
            title: 'Неудача',
            content: 'Неудалось сохранить баннер',
            type: 'error',
          })
          throw new Error('Failed to upload banner')
        }

        if (bannerData.value && bannerData.value.banner) {
          profile.value.banner = bannerData.value.banner
        }
      } else if (selectedBannerFile.value === null && !profile.value.banner) {
        const { error: bannerError } = await useApi('/api/profile/banner', {
          method: 'DELETE',
        })

        if (bannerError.value) {
          notify({
            title: 'Неудача',
            content: 'Неудалось удалить баннер',
            type: 'error',
          })
          throw new Error('Failed to remove banner')
        }
      }

      authStore.user = { ...authStore.user, ...profile.value }
      isSuccess.value = true
      notify({
        title: 'Успешно',
        content: 'Настройки профиля сохранены',
        type: 'success',
      })
    } catch (_e: any) {
      errors.value.global = 'Ошибка при сохранении'
      notify({
        title: 'Неудача',
        content: 'Настройки профиля не сохранены',
        type: 'error',
      })
    } finally {
      isLoading.value = false
    }
  }

  const handleDeleteAccount = async () => {
    const { error } = await useApi<{ message: string }>('/api/profile/delete', {
      method: 'POST',
    })

    if (error.value) {
      notify({
        title: 'Неудача',
        content: `Аккаунт не был удалён: ${error.value.data?.error || error.value.message}`,
        type: 'error',
      })
      return
    }

    notify({
      title: 'Успех',
      content: 'Ваш аккаунт был успешно удалён!',
      type: 'success',
    })

    notify({
      title: 'Перенаправление',
      content: 'Вы будете перенаправлены на главную страницу',
      type: 'info',
    })

    navigateTo('/')
  }

  return {
    profile,
    fileInput,
    bannerInput,
    isBirthDateDisable,
    errors,
    isLoading,
    isSuccess,
    handleAvatarClick,
    handleFileChange,
    removeAvatar,
    handleBannerClick,
    handleBannerChange,
    handleDeleteAccount,
    removeBanner,
    saveChanges,
  }
}
