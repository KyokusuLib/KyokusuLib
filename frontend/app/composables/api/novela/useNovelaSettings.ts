import { reactive, ref } from 'vue'
import { YEAR_DEFAULT } from '~/constants/data'
import type { NovelaDetails } from '~/types/backend/novela'

export function useNovelaSettings(novela: NovelaDetails) {
  const fileInput = ref<HTMLInputElement | null>(null)
  const previewUrl = ref<string | null>(staticImage(novela.poster_url) || null)
  const selectedFile = ref<File | null>(null)

  const form = reactive({
    title: novela.title,
    alternativeTitles: novela.alternative_titles?.join(' / ') || '',
    description: novela.description,
    type: novela.type,
    ageRating: novela.age_rating,
    releaseYear: novela.release_date !== '' ? novela.release_date.split('-')[0]! : YEAR_DEFAULT,
    status: novela.status,
    translationStatus: novela.translation_status,
    country: novela.country,
    genres: novela.genres || [],
    categories: novela.categories || [],
    authors: novela.authors?.map((a) => ({ id: a.id, label: a.name })) || [],
  })

  const handleImageClick = () => fileInput.value?.click()

  const handleImageChange = (event: Event) => {
    const target = event.target as HTMLInputElement
    if (target.files && target.files[0]) {
      const file = target.files[0]
      selectedFile.value = file
      previewUrl.value = URL.createObjectURL(file)
    }
  }

  const removeImage = () => {
    selectedFile.value = null
    previewUrl.value = null
  }

  return {
    form,
    fileInput,
    previewUrl,
    selectedFile,
    handleImageClick,
    handleImageChange,
    removeImage,
  }
}
