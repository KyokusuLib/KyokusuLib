import { useApi } from '~/composables/api/useApi'
import { useNotificationStore } from '@/stores/notification'
import type {
  MessageNetworkSettings,
  SocialNewtworkSettings,
} from '@/types/frontend/social-settings'
import { ref } from 'vue'

export function useSocialNetworkSettings() {
  const { notify } = useNotificationStore()
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase || ''

  const socialNetworkData = ref<SocialNewtworkSettings>({
    is_discord_connected: false,
    is_google_connected: false,
  })

  const syncSocialNetworkConnects = async () => {
    try {
      const { data } = await useApi<SocialNewtworkSettings>('/api/socials')

      if (data.value) {
        socialNetworkData.value = data.value
      }
    } catch (e) {
      console.log(e)
    }
  }

  const handleDiscordLogout = async () => {
    try {
      const { data } = await useApi<MessageNetworkSettings>('/api/socials/discord/logout')

      if (data.value) {
        notify({
          title: 'Успех',
          content: 'Произошел выход из discord аккаунта',
          type: 'info',
        })
        socialNetworkData.value.is_discord_connected = false
      }
    } catch (e) {
      console.log(e)
    }
  }

  const handleGoogleLogout = async () => {
    try {
      const { data } = await useApi<MessageNetworkSettings>('/api/socials/google/logout')

      if (data.value) {
        notify({
          title: 'Успех',
          content: 'Произошел выход из google аккаунта',
          type: 'info',
        })
        socialNetworkData.value.is_google_connected = false
      }
    } catch (e) {
      console.log(e)
    }
  }

  const handleGoogleLogin = () => {
    window.location.href = `${apiBase}/api/socials/google/login`
  }

  const handleDiscordLogin = () => {
    window.location.href = `${apiBase}/api/socials/discord/login`
  }

  const handleDiscord = async () => {
    if (socialNetworkData.value.is_discord_connected) {
      await handleDiscordLogout()
    } else {
      handleDiscordLogin()
    }
  }

  const handleGoogle = async () => {
    if (socialNetworkData.value.is_google_connected) {
      await handleGoogleLogout()
    } else {
      handleGoogleLogin()
    }
  }

  if (import.meta.client) {
    syncSocialNetworkConnects()
  }

  return {
    socialNetworkData,
    handleDiscordLogin,
    handleGoogleLogin,
    handleDiscordLogout,
    handleGoogleLogout,
    handleDiscord,
    handleGoogle,
  }
}
