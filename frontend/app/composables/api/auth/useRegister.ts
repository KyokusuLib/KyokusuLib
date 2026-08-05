import { reactive, ref } from 'vue'
import { useApi } from '~/composables/api/useApi'

export function useRegister() {
  const isLoading = ref(false)
  const errorMessage = ref('')
  const isSuccess = ref(false)
  const registeredEmail = ref('')

  const form = reactive({
    email: '',
    username: '',
    password: '',
    passwordConfirm: '',
  })

  const handleManualRegister = async () => {
    errorMessage.value = ''
    isSuccess.value = false

    if (!form.email || !form.username || !form.password || !form.passwordConfirm) {
      errorMessage.value = 'Заполните все поля'
      return
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(form.email)) {
      errorMessage.value = 'Некорректный формат email'
      return
    }

    const usernameRegex = /^(?=.*[a-z])(?=.*[A-Z])[a-zA-Z]{6,50}$/
    if (!usernameRegex.test(form.username)) {
      errorMessage.value = 'Никнейм: 6-50 символов, латиница, мин. 1 большая и 1 маленькая буква'
      return
    }

    const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{6,50}$/
    if (!passwordRegex.test(form.password)) {
      errorMessage.value =
        'Пароль должен быть от 6 символов, содержать латиницу, цифру, заглавную и строчную букву'
      return
    }

    if (form.password !== form.passwordConfirm) {
      errorMessage.value = 'Пароли не совпадают'
      return
    }

    isLoading.value = true

    try {
      const payload = {
        email: form.email,
        username: form.username,
        password: form.password,
      }

      const { data, error } = await useApi<any>('/api/auth/register', {
        method: 'POST',
        body: payload,
      })

      if (error.value) {
        const errorMsg =
          error.value.data?.message || error.value.data?.error || 'Ошибка при регистрации'
        throw new Error(errorMsg)
      }

      isSuccess.value = true
      registeredEmail.value = data.value?.user?.email || form.email
    } catch (err: any) {
      errorMessage.value = err.message || 'Произошла ошибка при регистрации'
    } finally {
      isLoading.value = false
    }
  }

  return {
    form,
    isLoading,
    errorMessage,
    isSuccess,
    registeredEmail,
    handleManualRegister,
  }
}
