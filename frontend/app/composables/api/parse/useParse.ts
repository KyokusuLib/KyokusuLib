import { ref } from 'vue'
import { $api } from '@/composables/api/useApi'
import { useNotificationStore } from '@/stores/notification'

export function useParse() {
  const isLoading = ref(false)
  const notificationStore = useNotificationStore()

  async function parseRanobeHub(id: string) {
    isLoading.value = true
    try {
      const data = await $api<{ message: string }>(`/parse/ranobehub/${id}`, {
        method: 'POST',
      })
      notificationStore.notify({
        title: 'Успех',
        content: data.message || 'Новелла успешно спаршена',
        type: 'success',
      })
      return true
    } catch (err: any) {
      notificationStore.notify({
        title: 'Ошибка',
        content: err.message || 'Не удалось спарсить новеллу',
        type: 'error',
      })
      return false
    } finally {
      isLoading.value = false
    }
  }

  return {
    isLoading,
    parseRanobeHub,
  }
}
