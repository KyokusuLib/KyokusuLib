export interface UploadedImage {
  id: string
  sourceUrl: string
  previewUrl: string
  caption: string
  isGallery?: boolean
}

let imgIdCounter = 0

export function useImageUpload() {
  const images = ref<UploadedImage[]>([])

  function fileToDataUrl(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onload = () => resolve(String(reader.result || ''))
      reader.onerror = () => reject(new Error('Failed to read file'))
      reader.readAsDataURL(file)
    })
  }

  async function addFile(file: File) {
    const id = `img-${Date.now()}-${++imgIdCounter}`
    const sourceUrl = await fileToDataUrl(file)
    const item: UploadedImage = {
      id,
      sourceUrl,
      previewUrl: URL.createObjectURL(file),
      caption: '',
    }
    images.value.push(item)
    return item
  }

  function removeImageById(id: string) {
    const index = images.value.findIndex((img) => img.id === id)
    if (index === -1) return
    const img = images.value[index]
    if (img) URL.revokeObjectURL(img.previewUrl)
    images.value.splice(index, 1)
  }

  function getImageById(id: string) {
    return images.value.find((img) => img.id === id) || null
  }

  function updateCaption(id: string, caption: string) {
    const img = images.value.find((i) => i.id === id)
    if (!img) return
    img.caption = caption
  }

  function clearAll() {
    for (const img of images.value) {
      URL.revokeObjectURL(img.previewUrl)
    }
    images.value = []
  }

  onBeforeUnmount(() => {
    clearAll()
  })

  return {
    images,
    addFile,
    removeImageById,
    getImageById,
    updateCaption,
    clearAll,
  }
}
