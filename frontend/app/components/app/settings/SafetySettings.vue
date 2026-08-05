<script lang="ts" setup>
import { useSecuritySettings } from '@/composables/api/settings/useSecuritySettings'
import { useAuthStore } from '@/stores/auth'
import { Input as BaseInput } from '@kyokusu-ui/vue'
import { ref } from 'vue'

const { user } = useAuthStore()
const {
  form,
  errors,
  isLoading,
  isSuccess,
  resetSuccessMessage,
  saveChanges,
  handleSendResetLink,
} = useSecuritySettings()

const showPassword = ref(false)
</script>

<template>
  <div class="space-y-8 cursor-default">
    <div class="flex flex-col items-center">
      <h2 class="text-3xl font-bold text-dark dark:text-white mb-1">Безопасность</h2>
    </div>

    <transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="transform -translate-y-2 opacity-0"
      enter-to-class="transform translate-y-0 opacity-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isSuccess"
        class="p-3 bg-green-500/10 border border-green-500/20 rounded-md text-green-400 text-sm"
      >
        Настройки безопасности успешно обновлены!
      </div>
      <div
        v-else-if="resetSuccessMessage"
        class="p-3 bg-blue-500/10 border border-blue-500/20 rounded-md text-blue-400 text-sm"
      >
        {{ resetSuccessMessage }}
      </div>
      <div
        v-else-if="errors.global"
        class="p-3 bg-red-500/10 border border-red-500/20 rounded-md text-red-400 text-sm"
      >
        {{ errors.global }}
      </div>
    </transition>

    <form @submit.prevent="saveChanges" class="space-y-6 text-left">
      <div class="pb-6 border-b border-zinc-800 space-y-4">
        <h3 class="text-base text-zinc-700 dark:text-zinc-200 leading-xl font-semibold">
          Смена электронной почты
        </h3>

        <div class="space-y-2">
          <div class="flex gap-2 items-center mb-2">
            <span class="text-xs text-stone-500">Текущая почта:</span>
            <span class="text-xs text-zinc-300 font-bold">{{ user?.email || '' }}</span>
          </div>

          <BaseInput
            id="email"
            v-model="form.email"
            type="email"
            autocomplete="email"
            placeholder="example@kyokusu.lib"
            :disabled="isLoading"
            :error="errors.email"
            description="На этот адрес будут приходить уведомления и письма для восстановления доступа."
          />
        </div>
      </div>

      <div class="space-y-6">
        <h3 class="text-base text-zinc-700 dark:text-zinc-200 leading-xl font-semibold">
          Смена пароля
        </h3>

        <div class="grid grid-cols-1 gap-6">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
            <BaseInput
              id="new_password"
              label="Новый пароль"
              v-model="form.newPassword"
              type="password"
              autocomplete="new-password"
              :show-text="showPassword"
              :disabled="isLoading"
              :error="errors.newPassword"
            />

            <BaseInput
              id="repeat_password"
              label="Повторите пароль"
              v-model="form.repeatPassword"
              type="password"
              autocomplete="new-password"
              :disabled="isLoading"
              :error="errors.repeatPassword"
            />
          </div>

          <p class="text-xs text-stone-500 text-left text-wrap ml-[3px]">
            Ограничение на пароль: мин. 6 символов, макс. 50 символов, латиница, мин. 1 большая
            буква, мин. 1 маленькая.
          </p>
        </div>
      </div>

      <div class="pt-6 flex flex-col border-t border-zinc-800">
        <BaseInput
          id="current_password"
          label="Текущий пароль"
          v-model="form.currentPassword"
          type="password"
          autocomplete="current-password"
          :disabled="isLoading"
          :error="errors.currentPassword"
        />
        <p class="text-xs text-stone-500 text-left text-wrap ml-[3px]">
          При смене email или пароля, необходим ввод текущего пароля.
        </p>
        <div class="mt-2 ml-[3px]">
          <button
            type="button"
            @click="handleSendResetLink"
            :disabled="isLoading"
            class="text-xs text-zinc-600 hover:text-zinc-800 dark:text-zinc-500 dark:hover:text-zinc-300 transition-colors cursor-pointer disabled:opacity-50"
          >
            Забыли пароль?
          </button>
        </div>
      </div>

      <div class="pt-6 border-t border-zinc-800 flex justify-end gap-3">
        <button
          type="submit"
          :disabled="isLoading"
          class="px-4 py-2 text-sm font-medium bg-black dark:bg-white dark:hover:bg-zinc-300 text-white dark:text-zinc-900 rounded transition-colors focus:ring-2 focus:ring-offset-2 focus:ring-offset-zinc-900 focus:ring-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          <svg
            v-if="isLoading"
            class="animate-spin -ml-1 mr-2 h-4 w-4 text-zinc-900"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          {{ isLoading ? 'Сохранение...' : 'Сохранить изменения' }}
        </button>
      </div>
    </form>
  </div>
</template>
