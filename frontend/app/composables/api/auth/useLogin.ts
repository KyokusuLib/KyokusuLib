import { reactive, ref } from 'vue'
import { useApi } from '~/composables/api/useApi'

export function useLogin() {
  const authStore = useAuthStore()
  const { notify } = useNotificationStore()
  const config = useRuntimeConfig()

  const isLoading = ref(false)
  const errorMessage = ref('')
  const successMessage = ref('')

  const isForgotPasswordMode = ref(false)
  const isResetPasswordMode = ref(false)
  const resetToken = ref('')

  const form = reactive({
    email: '',
    password: '',
  })

  const toggleForgotPasswordMode = () => {
    isForgotPasswordMode.value = !isForgotPasswordMode.value
    isResetPasswordMode.value = false
    errorMessage.value = ''
    successMessage.value = ''
    form.password = ''
  }

  const handleManualLogin = async () => {
    errorMessage.value = ''
    isLoading.value = true

    try {
      const { data, error } = await useApi<any>('/api/auth/login', {
        method: 'POST',
        body: form,
      })

      if (error.value) {
        throw new Error(error.value.data?.error || 'Неверный email или пароль')
      }

      if (data.value) {
        authStore.user = data.value
        await navigateTo(`/profile/${data.value.id}`)
      }

      notify({
        title: 'Успех',
        content: 'Успешный вход в аккаунт',
        type: 'success',
      })
    } catch (err: any) {
      errorMessage.value = err.message || 'Произошла ошибка при входе'
      notify({
        title: 'Неудача',
        content: errorMessage.value,
        type: 'error',
      })
    } finally {
      isLoading.value = false
    }
  }

  const handleSendResetLink = async () => {
    errorMessage.value = ''
    successMessage.value = ''
    isLoading.value = true

    try {
      const { error } = await useApi('/api/auth/forgot-password', {
        method: 'POST',
        body: { email: form.email },
      })

      if (error.value) {
        throw new Error('Не удалось отправить письмо. Проверьте email.')
      }

      successMessage.value = 'Ссылка для сброса пароля отправлена на вашу почту.'
    } catch (err: any) {
      errorMessage.value = err.message || 'Ошибка отправки запроса'
    } finally {
      isLoading.value = false
    }
  }

  const handleResetPassword = async () => {
    errorMessage.value = ''
    successMessage.value = ''
    isLoading.value = true

    try {
      const { error } = await useApi('/api/auth/reset-password', {
        method: 'POST',
        body: {
          token: resetToken.value,
          password: form.password,
        },
      })

      if (error.value) {
        throw new Error('Ошибка сброса пароля. Возможно, ссылка устарела.')
      }

      successMessage.value = 'Пароль успешно изменен! Теперь вы можете войти.'

      setTimeout(() => {
        isResetPasswordMode.value = false
        resetToken.value = ''
        form.password = ''
        successMessage.value = ''
      }, 3000)
    } catch (err: any) {
      errorMessage.value = err.message || 'Ошибка сброса пароля'
    } finally {
      isLoading.value = false
    }
  }

  const apiBase = config.public.apiBase || ''

  const handleGoogleLogin = () => {
    window.location.href = `${apiBase}/api/auth/google/login`
  }

  const handleDiscordLogin = () => {
    window.location.href = `${apiBase}/api/auth/discord/login`
  }

  return {
    form,
    isLoading,
    errorMessage,
    successMessage,
    isForgotPasswordMode,
    isResetPasswordMode,
    resetToken,
    toggleForgotPasswordMode,
    handleManualLogin,
    handleSendResetLink,
    handleResetPassword,
    handleGoogleLogin,
    handleDiscordLogin,
  }
}
