import { ThemeSetting } from '@/types/enums/theme-enum'
import { useAuthStore } from '@/stores/auth'
import { useReaderSettings } from '@/composables/ui/useReaderSettings'

import { useApi, $api } from '../useApi'

export function useInterfaceSettings() {
  const colorMode = useColorMode()
  const readerSettings = useReaderSettings()
  const authStore = useAuthStore()

  const isShowTag = useState<boolean>('settings-is-show-tag', () => false)
  const isShowBookmark = useState<boolean>('settings-is-show-bookmark', () => true)
  const isInitialized = useState<boolean>('settings-is-initialized', () => false)

  const isDarkTheme = computed({
    get: () => colorMode.preference === ThemeSetting.DARK_THEME,
    set: (val: boolean) => {
      colorMode.preference = val ? ThemeSetting.DARK_THEME : ThemeSetting.LIGHT_THEME
    },
  })

  const syncSettingWithBackend = async () => {
    if (isInitialized.value || !authStore.isAuthenticated) return

    try {
      const { data } = await useApi<any>('/api/profile/settings')
      if (data.value) {
        if (data.value.theme) colorMode.preference = data.value.theme
        if (data.value.is_show_tag !== undefined) isShowTag.value = data.value.is_show_tag
        if (data.value.is_show_bookmark !== undefined)
          isShowBookmark.value = data.value.is_show_bookmark
        if (data.value.font_size !== undefined) readerSettings.fontSize.value = data.value.font_size
        if (data.value.line_weight !== undefined)
          readerSettings.lineWeight.value = data.value.line_weight
        if (data.value.font_family !== undefined)
          readerSettings.fontFamily.value = data.value.font_family
        if (data.value.auto_scroll !== undefined)
          readerSettings.isAutoScrollEnabled.value = data.value.auto_scroll

        isInitialized.value = true
      }
    } catch (e) {
      console.error('Sync error', e)
    }
  }

  watch(isShowTag, async (newVal) => {
    if (import.meta.client && isInitialized.value) {
      try {
        await $api('/api/profile/settings/interface', {
          method: 'PATCH',
          body: { is_show_tag: newVal },
        })
      } catch (e) {
        console.error('Failed to save is_show_tag', e)
      }
    }
  })

  watch(isShowBookmark, async (newVal) => {
    if (import.meta.client && isInitialized.value) {
      try {
        await $api('/api/profile/settings/interface', {
          method: 'PATCH',
          body: { is_show_bookmark: newVal },
        })
      } catch (e) {
        console.error('Failed to save is_show_bookmark', e)
      }
    }
  })

  watch(readerSettings.fontSize, async (newVal) => {
    if (import.meta.client && isInitialized.value) {
      try {
        await $api('/api/profile/settings/reader', {
          method: 'PATCH',
          body: { font_size: newVal },
        })
      } catch (e) {
        console.error('Failed to save font_size', e)
      }
    }
  })

  watch(readerSettings.lineWeight, async (newVal) => {
    if (import.meta.client && isInitialized.value) {
      try {
        await $api('/api/profile/settings/reader', {
          method: 'PATCH',
          body: { line_weight: newVal },
        })
      } catch (e) {
        console.error('Failed to save line_weight', e)
      }
    }
  })

  watch(readerSettings.isAutoScrollEnabled, async (newVal) => {
    if (import.meta.client && isInitialized.value) {
      try {
        await $api('/api/profile/settings/reader', {
          method: 'PATCH',
          body: { auto_scroll: newVal },
        })
      } catch (e) {
        console.error('Failed to save auto_scroll', e)
      }
    }
  })

  watch(readerSettings.fontFamily, async (newVal) => {
    if (import.meta.client && isInitialized.value) {
      try {
        await $api('/api/profile/settings/reader', {
          method: 'PATCH',
          body: { font_family: newVal },
        })
      } catch (e) {
        console.error('Failed to save font_family', e)
      }
    }
  })

  return {
    isDarkTheme,
    isShowTag,
    isShowBookmark,
    syncSettingWithBackend,
    colorMode,
  }
}
