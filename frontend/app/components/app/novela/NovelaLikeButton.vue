<script setup lang="ts">
import { ref } from 'vue'
import { useNovelaLike } from '@/composables/api/novela/useNovelaLike'
import { useAuthStore } from '@/stores/auth'
import AuthRequiredModal from '@/components/common/AuthRequiredModal.vue'

interface Props {
  modelValue: boolean
  novelaId: number
}

const props = defineProps<Props>()
const emit = defineEmits(['update:modelValue'])

const { setNovelaLike, loading } = useNovelaLike()
const authStore = useAuthStore()

const isAuthModalOpen = ref(false)

const handleLikeClick = async () => {
  if (!authStore.isAuthenticated) {
    isAuthModalOpen.value = true
    return
  }

  try {
    await setNovelaLike({
      novela_id: props.novelaId,
      has_liked: !props.modelValue,
    })
    emit('update:modelValue', !props.modelValue)
  } catch (e) {
    console.error('Ошибка при установке лайка:', e)
  }
}
</script>

<template>
  <div class="w-full">
    <button
      @click="handleLikeClick"
      :disabled="loading"
      type="button"
      class="w-full py-2.5 cursor-pointer rounded-xl border transition-all flex items-center justify-center gap-2 text-sm font-semibold group"
      :class="[
        modelValue
          ? 'bg-yellow-50 border-yellow-200 text-yellow-600 dark:bg-yellow-500/10 dark:border-yellow-500/20 dark:text-yellow-400'
          : 'bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 text-zinc-700 dark:text-zinc-200 hover:bg-zinc-50 dark:hover:bg-zinc-800',
      ]"
    >
      <Icon
        :name="modelValue ? 'ph:heart-fill' : 'ph:heart-bold'"
        size="22"
        :class="[
          'transition-all duration-300',
          modelValue ? 'text-yellow-500 scale-110' : 'text-zinc-400',
        ]"
      />
      <span>{{ modelValue ? 'Понравилось' : 'Лайк' }}</span>
      <div
        v-if="loading"
        class="ml-1 w-3 h-3 border-2 border-current border-t-transparent rounded-full animate-spin"
      ></div>
    </button>

    <AuthRequiredModal v-model="isAuthModalOpen" action-text="ставить лайки" />
  </div>
</template>
