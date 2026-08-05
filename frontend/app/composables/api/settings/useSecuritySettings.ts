import { reactive, ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { useApi } from '~/composables/api/useApi'

export function useSecuritySettings() {
  const { user } = useAuthStore()
  const { notify } = useNotificationStore()
  const isLoading = ref(false)
  const isSuccess = ref(false)
  const resetSuccessMessage = ref('')

  const form = reactive({
    email: '',
    currentPassword: '',
    newPassword: '',
    repeatPassword: '',
  })

  const errors = ref<Record<string, string>>({})

  const validate = (): boolean => {
    errors.value = {}
    let isValid = true

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

    if (form.email && form.email !== user?.email) {
      if (!emailRegex.test(form.email)) {
        errors.value.email = 'Введите корректный email.'
        isValid = false
      }
      if (!form.currentPassword) {
        errors.value.confirmEmailPassword = 'Введите пароль для подтверждения смены email.'
        isValid = false
      }
    }

    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])[a-zA-Z0-9]{6,50}$/

    if (form.newPassword || form.repeatPassword) {
      if (!form.currentPassword) {
        errors.value.currentPassword = 'Введите текущий пароль.'
        isValid = false
      }

      if (!passwordRegex.test(form.newPassword)) {
        errors.value.newPassword = 'Пароль: 6-50 симв., латиница, 1 заглавная, 1 строчная.'
        isValid = false
      }

      if (form.newPassword !== form.repeatPassword) {
        errors.value.repeatPassword = 'Пароли не совпадают.'
        isValid = false
      }
    }

    if (!form.email && !form.newPassword) {
      errors.value.global = 'Нет изменений для сохранения.'
      isValid = false
    }

    return isValid
  }

  const handleSendResetLink = async () => {
    if (!user?.email) {
      errors.value.global = 'Email не найден.'
      return
    }

    isLoading.value = true
    resetSuccessMessage.value = ''
    errors.value = {}

    try {
      const { error } = await useApi('/api/auth/forgot-password', {
        method: 'POST',
        body: { email: user.email },
      })

      if (error.value) {
        throw new Error('Не удалось отправить письмо.')
      }

      resetSuccessMessage.value = 'Ссылка для сброса пароля отправлена на вашу почту.'

      setTimeout(() => {
        resetSuccessMessage.value = ''
      }, 5000)
    } catch (e: any) {
      errors.value.global = e.message || 'Ошибка отправки запроса.'
    } finally {
      isLoading.value = false
    }
  }

  const saveChanges = async () => {
    if (!validate()) return

    isLoading.value = true
    isSuccess.value = false
    errors.value = {}
    resetSuccessMessage.value = ''

    try {
      if (form.currentPassword && form.newPassword) {
        const { error } = await useApi<any>('/api/profile/password', {
          method: 'POST',
          body: {
            old_password: form.currentPassword,
            new_password: form.newPassword,
          },
        })

        if (error.value) {
          notify({
            title: 'Неудача',
            content: 'Не удалось обновить пароль',
            type: 'error',
          })
          throw new Error('Не удалось обновить пароль')
        } else {
          notify({
            title: 'Успех',
            content: 'Успешная смена пароля',
            type: 'success',
          })
        }
      }

      const tempCurrentPassword = form.currentPassword
      form.currentPassword = ''
      form.newPassword = ''
      form.repeatPassword = ''

      if (form.email) {
        const { error } = await useApi<any>('/api/profile/email', {
          method: 'POST',
          body: { old_password: tempCurrentPassword, email: form.email },
        })

        if (error.value) {
          notify({
            title: 'Неудача',
            content: 'Не удалось обновить email',
            type: 'error',
          })
          throw new Error('Не удалось обновить пароль')
        } else {
          notify({
            title: 'Успех',
            content: 'Успешная смена email',
            type: 'success',
          })
        }
        form.email = ''
      }

      isSuccess.value = true
    } catch (e: any) {
      errors.value.global = e.message || 'Ошибка при сохранении настроек.'
    } finally {
      isLoading.value = false
    }
  }

  return {
    form,
    errors,
    isLoading,
    isSuccess,
    resetSuccessMessage,
    saveChanges,
    handleSendResetLink,
  }
}
