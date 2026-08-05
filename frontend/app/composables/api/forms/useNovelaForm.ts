import { reactive, ref } from 'vue'
import { useApi } from '~/composables/api/useApi'
import { useNotificationStore } from '@/stores/notification'

export interface NovelaFormState {
  title: string
  alternativeTitles: string
  description: string
  type: string
  ageRating: string
  releaseYear: string
  status: string
  translationStatus: string
  country: string
  categories: string[]
  genres: string[]
}

export function useNovelaForm() {
  const { notify } = useNotificationStore()

  const isLoading = ref(false)
  const isSuccess = ref(false)
  const errors = ref<Record<string, string>>({})

  const form = reactive<NovelaFormState>({
    title: '',
    alternativeTitles: '',
    description: '',
    type: '',
    ageRating: '',
    releaseYear: '',
    status: '',
    translationStatus: '',
    country: '',
    categories: [],
    genres: [],
  })

  const fileInput = ref<HTMLInputElement | null>(null)
  const selectedFile = ref<File | null>(null)
  const previewUrl = ref<string | null>(null)

  const validate = (): boolean => {
    errors.value = {}
    let isValid = true

    if (!form.title.trim()) {
      errors.value.title = 'Название обязательно.'
      isValid = false
    }

    if (!form.type) {
      errors.value.type = 'Выберите тип.'
      isValid = false
    }

    if (!form.country) {
      form.country = 'Неизвестно'
    }

    if (!form.ageRating) {
      errors.value.ageRating = 'Выберите возрастной рейтинг.'
      isValid = false
    }

    if (!form.releaseYear || !/^\d{4}$/.test(form.releaseYear)) {
      errors.value.releaseYear = 'Укажите корректный год (YYYY).'
      isValid = false
    }

    if (form.description.length > 5000) {
      errors.value.description = 'Описание слишком длинное (макс. 5000 символов).'
      isValid = false
    }

    return isValid
  }

  const handleImageClick = () => fileInput.value?.click()

  const handleImageChange = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files[0]) {
      const file = target.files[0]
      if (file.size > 5 * 1024 * 1024) {
        notify({ title: 'Ошибка', content: 'Файл > 5МБ', type: 'error' })
        return
      }
      selectedFile.value = file
      const reader = new FileReader()
      reader.onload = (e) => {
        if (e.target?.result) previewUrl.value = e.target.result as string
      }
      reader.readAsDataURL(file)
    }
  }

  const removeImage = () => {
    selectedFile.value = null
    previewUrl.value = null
    if (fileInput.value) fileInput.value.value = ''
  }

  const submit = async () => {
    if (!validate()) return

    isLoading.value = true
    isSuccess.value = false

    try {
      const formData = new FormData()
      formData.append('title', form.title)
      formData.append('description', form.description.replace(/<\/?(div|br)\s*\/?>/g, '').trim())
      formData.append('type', form.type)
      formData.append('age_rating', form.ageRating)
      formData.append('release_year', form.releaseYear)
      formData.append('status', form.status)
      formData.append('translation_status', form.translationStatus)
      formData.append('country', form.country)

      formData.append(
        'alternative_titles',
        JSON.stringify(form.alternativeTitles.split('/').map((t) => t.trim())),
      )
      formData.append('categories', JSON.stringify(form.categories))
      formData.append('genres', JSON.stringify(form.genres))

      if (selectedFile.value) {
        formData.append('poster', selectedFile.value)
      }

      const { data, error } = await useApi<any>('/api/novela', {
        method: 'POST',
        body: formData,
      })

      if (error.value) {
        const serverMessage = data.value?.error
        throw new Error(serverMessage || 'Ошибка сервера')
      }

      isSuccess.value = true
      notify({
        title: 'Успешно',
        content: 'Новела добавлена',
        type: 'success',
      })

      Object.assign(form, {
        title: '',
        alternativeTitles: '',
        description: '',
        type: '',
        ageRating: '',
        releaseYear: '',
        status: '',
        translationStatus: '',
        categories: [],
        genres: [],
      })
      removeImage()
    } catch (e: any) {
      const msg = e.error || 'Ошибка при создании'
      errors.value.global = msg
      notify({ title: 'Ошибка', content: msg, type: 'error' })
    } finally {
      isLoading.value = false
    }
  }

  return {
    form,
    errors,
    isLoading,
    isSuccess,
    fileInput,
    previewUrl,
    handleImageClick,
    handleImageChange,
    removeImage,
    submit,
  }
}
