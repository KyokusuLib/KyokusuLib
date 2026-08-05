import { reactive, ref } from 'vue'
import { useApi } from '~/composables/api/useApi'
import { useNotificationStore } from '@/stores/notification'
import type { NovelaAuthorDetails } from '@/types/backend/novela'
import { staticImage } from '@/utils/str'

export interface AuthorSettingsState {
  name: string
  about: string
  country: string
  metier: string
}

export function useAuthorSettings(initialAuthor: NovelaAuthorDetails) {
  const { notify } = useNotificationStore()

  const isUpdating = ref(false)
  const errors = ref<Record<string, string>>({})

  const form = reactive<AuthorSettingsState>({
    name: initialAuthor.name || '',
    about: initialAuthor.bio || '',
    country: initialAuthor.country || '',
    metier: initialAuthor.metier || '',
  })

  const fileInput = ref<HTMLInputElement | null>(null)
  const selectedFile = ref<File | null>(null)
  const previewUrl = ref<string | null>(staticImage(initialAuthor.picture || '') || null)

  const validate = (): boolean => {
    errors.value = {}
    let isValid = true

    if (!form.name.trim()) {
      errors.value.name = 'Имя автора обязательно.'
      isValid = false
    }

    if (!form.country.trim()) {
      errors.value.country = 'Страна обязательна.'
      isValid = false
    }

    if (!form.metier.trim()) {
      errors.value.metier = 'Род деятельности обязателен.'
      isValid = false
    }

    if (form.about.length > 2000) {
      errors.value.about = 'Описание слишком длинное (макс. 2000 символов).'
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

  const updateAuthor = async (id: number) => {
    if (!validate()) return null

    isUpdating.value = true

    try {
      const formData = new FormData()
      formData.append('name', form.name)
      formData.append('bio', form.about)
      formData.append('country', form.country)
      formData.append('metier', form.metier)

      if (selectedFile.value) {
        formData.append('picture', selectedFile.value)
      }

      const { error } = await useApi<any>(`/api/author/${id}`, {
        method: 'PUT',
        body: formData,
      })

      if (error.value) {
        const errorMsg = error.value.data?.error || 'Неизвестная ошибка'
        throw new Error(errorMsg)
      }

      notify({
        title: 'Успешно',
        content: 'Автор успешно обновлен',
        type: 'success',
      })

      return {
        id,
        name: form.name,
        bio: form.about,
        country: form.country,
        metier: form.metier,
        picture: previewUrl.value,
      }
    } catch (e: any) {
      const msg = e.message || 'Неизвестная ошибка при обновлении автора'
      errors.value.global = msg
      notify({
        title: 'Ошибка',
        content: msg,
        type: 'error',
      })
      return null
    } finally {
      isUpdating.value = false
    }
  }

  return {
    form,
    errors,
    isUpdating,
    fileInput,
    previewUrl,
    handleImageClick,
    handleImageChange,
    removeImage,
    updateAuthor,
  }
}
