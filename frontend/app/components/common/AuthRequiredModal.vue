<script setup lang="ts">
import { ModalWindow } from '@kyokusu-ui/vue'

defineOptions({
  inheritAttrs: false,
})

interface Props {
  modelValue: boolean
  actionText?: string
}

defineProps<Props>()
defineEmits(['update:modelValue'])
</script>

<template>
  <ModalWindow
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="Требуется авторизация"
    :center-title="true"
    width="w-full max-w-md"
  >
    <div class="space-y-4 py-2 text-center">
      <p class="text-zinc-600 dark:text-zinc-400">
        Чтобы {{ actionText || 'продолжить' }}, пожалуйста, водите в аккаунт или зарегистрируйтесь.
      </p>

      <div class="flex flex-col items-center gap-2">
        <NuxtLink
          to="/login"
          class="flex justify-center items-center w-1/2 py-1.5 bg-yellow-500 hover:bg-yellow-600 text-white font-bold rounded-xl transition-colors"
        >
          Войти
        </NuxtLink>

        <div class="text-center text-sm text-zinc-500">
          Нет аккаунта?
          <NuxtLink
            to="/register"
            class="text-yellow-600 underline underline-offset-4 hover:text-yellow-700 transition-colors"
          >
            Зарегистрироваться
          </NuxtLink>
        </div>
      </div>
    </div>
  </ModalWindow>
</template>
