import { ref, type Ref } from 'vue'

export function useImageEdit(avatarModel: Ref<string>, bannerModel: Ref<string>) {
  const avatarInput = ref<HTMLInputElement | null>(null)
  const bannerInput = ref<HTMLInputElement | null>(null)
  const isUploading = ref(false)

  function trigger(el: HTMLInputElement | null) {
    el?.click()
  }

  function onFileSelect(event: Event): Promise<string | null> {
    return new Promise((resolve) => {
      const target = event.target as HTMLInputElement
      const file = target.files?.[0]
      if (!file) {
        resolve(null)
        return
      }

      const reader = new FileReader()
      reader.onload = (e) => resolve((e.target?.result as string) ?? null)
      reader.readAsDataURL(file)
    })
  }

  function clear(el: HTMLInputElement | null) {
    if (el) el.value = ''
  }

  async function onAvatarChange(event: Event) {
    const dataUrl = await onFileSelect(event)
    if (dataUrl) avatarModel.value = dataUrl
  }

  async function onBannerChange(event: Event) {
    const dataUrl = await onFileSelect(event)
    if (dataUrl) bannerModel.value = dataUrl
  }

  function removeAvatar() {
    avatarModel.value = ''
    clear(avatarInput.value)
  }

  function removeBanner() {
    bannerModel.value = ''
    clear(bannerInput.value)
  }

  return {
    avatarInput,
    bannerInput,
    isUploading,
    trigger,
    onAvatarChange,
    onBannerChange,
    removeAvatar,
    removeBanner,
  }
}
