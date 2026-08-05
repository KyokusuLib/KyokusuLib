import { reactive, ref } from 'vue'
import { useNotificationStore } from '@/stores/notification'
import type { Team } from '@/types/frontend/teams'
import { staticImage } from '@/utils/str'
import { useTeam } from './useTeams'

export interface TeamSettingsState {
  description: string
  owner_role_name: string
  moderator_role_name: string
  member_role_name: string
  team_type: string
}

export function useTeamSettings(initialTeam: Team) {
  const { notify } = useNotificationStore()
  const { updateTeam, isLoading: isUpdating } = useTeam()

  const errors = ref<Record<string, string>>({})

  const form = reactive<TeamSettingsState>({
    description: initialTeam.description || '',
    owner_role_name: initialTeam.role_names?.owner || 'Владелец',
    moderator_role_name: initialTeam.role_names?.moderator || 'Модератор',
    member_role_name: initialTeam.role_names?.member || 'Участник',
    team_type: initialTeam.team_type || 'open',
  })

  const fileInput = ref<HTMLInputElement | null>(null)
  const selectedFile = ref<File | null>(null)
  const previewUrl = ref<string | null>(staticImage(initialTeam.avatar_url || '') || null)

  const bannerInput = ref<HTMLInputElement | null>(null)
  const selectedBanner = ref<File | null>(null)
  const previewBannerUrl = ref<string | null>(staticImage(initialTeam.banner_url || '') || null)

  const validate = (): boolean => {
    errors.value = {}
    let isValid = true

    if (form.description.length > 2000) {
      errors.value.description = 'Описание слишком длинное (макс. 2000 символов).'
      isValid = false
    }

    if (form.owner_role_name.length > 100) {
      errors.value.owner_role_name = 'Название слишком длинное (макс. 100 символов).'
      isValid = false
    }
    if (form.moderator_role_name.length > 100) {
      errors.value.moderator_role_name = 'Название слишком длинное (макс. 100 символов).'
      isValid = false
    }
    if (form.member_role_name.length > 100) {
      errors.value.member_role_name = 'Название слишком длинное (макс. 100 символов).'
      isValid = false
    }

    return isValid
  }

  const handleImageClick = () => {
    fileInput.value?.click()
  }

  const handleImageChange = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files[0]) {
      const file = target.files[0]

      if (file.size > 5 * 1024 * 1024) {
        notify({
          title: 'Ошибка',
          content: 'Размер файла не должен превышать 5МБ',
          type: 'error',
        })
        return
      }

      selectedFile.value = file

      const reader = new FileReader()
      reader.onload = (e) => {
        if (e.target?.result) {
          previewUrl.value = e.target.result as string
        }
      }
      reader.readAsDataURL(file)
    }
  }

  const removeImage = () => {
    selectedFile.value = null
    previewUrl.value = null
    if (fileInput.value) fileInput.value.value = ''
  }

  const handleBannerClick = () => {
    bannerInput.value?.click()
  }

  const handleBannerChange = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files[0]) {
      const file = target.files[0]

      if (file.size > 5 * 1024 * 1024) {
        notify({
          title: 'Ошибка',
          content: 'Размер файла не должен превышать 5МБ',
          type: 'error',
        })
        return
      }

      selectedBanner.value = file

      const reader = new FileReader()
      reader.onload = (e) => {
        if (e.target?.result) {
          previewBannerUrl.value = e.target.result as string
        }
      }
      reader.readAsDataURL(file)
    }
  }

  const removeBanner = () => {
    selectedBanner.value = null
    previewBannerUrl.value = null
    if (bannerInput.value) bannerInput.value.value = ''
  }

  const submitUpdate = async (slug: string) => {
    if (!validate()) return null

    const formData = new FormData()
    formData.append('description', replaceTags(form.description))
    formData.append('owner_role_name', form.owner_role_name)
    formData.append('moderator_role_name', form.moderator_role_name)
    formData.append('member_role_name', form.member_role_name)
    formData.append('team_type', form.team_type)

    if (selectedFile.value) {
      formData.append('avatar', selectedFile.value)
    }

    if (selectedBanner.value) {
      formData.append('banner', selectedBanner.value)
    }

    const updatedTeam = await updateTeam(slug, formData)

    if (updatedTeam) {
      return updatedTeam
    }
    return null
  }

  return {
    form,
    errors,
    isUpdating,
    fileInput,
    previewUrl,
    bannerInput,
    previewBannerUrl,
    handleImageClick,
    handleImageChange,
    removeImage,
    handleBannerClick,
    handleBannerChange,
    removeBanner,
    submitUpdate,
  }
}
