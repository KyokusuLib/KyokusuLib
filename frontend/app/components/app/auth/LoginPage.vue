<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DiscordIcon from '@/assets/images/social/discord.png'
import GoogleIcon from '@/assets/images/social/google.png'
import ReferIcon from '@/assets/images/special/refer.png'
import FooterApp from '@/components/common/FooterApp.vue'
import { useLogin } from '@/composables/api/auth/useLogin'

const route = useRoute()
const router = useRouter()

const {
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
} = useLogin()

const isVerified = ref(false)
const showPassword = ref(false)

onMounted(() => {
  if (route.query.verified === 'true') {
    isVerified.value = true
    router.replace({ query: {} })
  }

  const token = route.query.token as string
  if (token) {
    resetToken.value = token
    isResetPasswordMode.value = true
    router.replace({ query: {} })
  }
})
</script>

<template>
  <main
    class="min-h-[100dvh] bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-white flex flex-col relative overflow-hidden font-sans transition-colors duration-300"
  >
    <div
      class="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] bg-zinc-400/20 rounded-full blur-[100px] pointer-events-none"
    ></div>
    <div
      class="absolute top-[80%] left-[70%] w-[40%] h-[50%] bg-zinc-400/20 rounded-full blur-[100px] pointer-events-none"
    ></div>

    <div class="flex-1 flex items-center justify-center p-4 sm:p-6 z-10">
      <div
        class="w-full max-w-md bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-xl p-8 backdrop-blur-sm transition-colors duration-300"
      >
        <div class="text-center mb-8">
          <div class="flex justify-center items-center gap-3 mb-3">
            <img :src="ReferIcon" class="w-6 h-6 object-contain dark:invert" />
            <h1 class="text-2xl font-bold tracking-tight uppercase">
              {{
                isResetPasswordMode
                  ? 'Смена пароля'
                  : isForgotPasswordMode
                    ? 'Восстановление'
                    : 'Вход'
              }}
            </h1>
          </div>
          <p class="text-zinc-500 dark:text-zinc-400 text-sm">
            {{
              isResetPasswordMode
                ? 'Введите новый пароль'
                : isForgotPasswordMode
                  ? 'Мы отправим ссылку для сброса'
                  : 'Войдите, чтобы продолжить'
            }}
          </p>
        </div>

        <transition name="fade">
          <div
            v-if="isVerified"
            class="mb-6 p-3 rounded-lg bg-green-500/10 border border-green-500/20 text-green-600 dark:text-green-400 text-sm flex items-center gap-2"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"></path>
              <polyline points="22 4 12 14.01 9 11.01"></polyline>
            </svg>
            Email успешно подтвержден! Теперь вы можете войти.
          </div>
        </transition>

        <transition name="fade">
          <div
            v-if="successMessage"
            class="mb-6 p-3 rounded-lg bg-green-500/10 border border-green-500/20 text-green-600 dark:text-green-400 text-sm flex items-center gap-2"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <polyline points="20 6 9 17 4 12"></polyline>
            </svg>
            {{ successMessage }}
          </div>
        </transition>

        <transition name="fade">
          <div
            v-if="errorMessage"
            class="mb-6 p-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-600 dark:text-red-400 text-sm flex items-center gap-2"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <circle cx="12" cy="12" r="10" />
              <line x1="12" x2="12" y1="8" y2="12" />
              <line x1="12" x2="12.01" y1="16" y2="16" />
            </svg>
            {{ errorMessage }}
          </div>
        </transition>

        <div v-if="isResetPasswordMode">
          <form @submit.prevent="handleResetPassword" class="space-y-5">
            <div class="space-y-1.5">
              <label
                for="new-password"
                class="block text-sm font-medium text-zinc-700 dark:text-zinc-300 ml-1 text-left"
                >Новый пароль</label
              >
              <div class="relative group">
                <input
                  id="new-password"
                  v-model="form.password"
                  type="password"
                  required
                  autocomplete="off"
                  placeholder="••••••••"
                  class="w-full bg-zinc-50 dark:bg-[#27272a] border border-zinc-300 dark:border-[#3f3f46] text-zinc-900 dark:text-white text-sm rounded-lg focus:ring-2 focus:ring-zinc-500 dark:focus:ring-white/20 focus:border-zinc-500 block p-2.5 pl-10 placeholder-zinc-400 dark:placeholder-zinc-600 transition-all outline-none"
                />
                <div
                  class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-zinc-400 dark:text-zinc-500 group-focus-within:text-zinc-800 dark:group-focus-within:text-white transition-colors"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                    <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                  </svg>
                </div>
              </div>
            </div>

            <button
              type="submit"
              :disabled="isLoading"
              class="w-full flex items-center cursor-pointer justify-center bg-zinc-900 hover:bg-zinc-800 dark:bg-white dark:hover:bg-zinc-200 text-white dark:text-black font-semibold rounded-lg text-sm px-5 py-2.5 text-center transition-all disabled:opacity-70 disabled:cursor-not-allowed mt-2"
            >
              <span v-if="isLoading" class="flex items-center gap-2">
                <svg
                  class="animate-spin h-4 w-4"
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
                Сохранение...
              </span>
              <span v-else>Сохранить пароль</span>
            </button>

            <button
              type="button"
              @click="$emit('update:isResetPasswordMode', false)"
              class="w-full flex items-center cursor-pointer justify-center text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-white font-medium text-sm px-5 py-2.5 text-center transition-all mt-2"
            >
              Отмена
            </button>
          </form>
        </div>

        <div v-else-if="isForgotPasswordMode">
          <form @submit.prevent="handleSendResetLink" class="space-y-5">
            <div class="space-y-1.5">
              <label
                for="reset-email"
                class="block text-sm font-medium text-zinc-700 dark:text-zinc-300 ml-1 text-left"
                >Почта</label
              >
              <div class="relative group">
                <input
                  id="reset-email"
                  v-model="form.email"
                  type="email"
                  required
                  placeholder="name@example.com"
                  class="w-full bg-zinc-50 dark:bg-[#27272a] border border-zinc-300 dark:border-[#3f3f46] text-zinc-900 dark:text-white text-sm rounded-lg focus:ring-2 focus:ring-zinc-500 dark:focus:ring-white/20 focus:border-zinc-500 block p-2.5 pl-10 placeholder-zinc-400 dark:placeholder-zinc-600 transition-all outline-none"
                />
                <div
                  class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-zinc-400 dark:text-zinc-500 group-focus-within:text-zinc-800 dark:group-focus-within:text-white transition-colors"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <rect width="20" height="16" x="2" y="4" rx="2" />
                    <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
                  </svg>
                </div>
              </div>
            </div>

            <button
              type="submit"
              :disabled="isLoading"
              class="w-full flex items-center cursor-pointer justify-center bg-zinc-900 hover:bg-zinc-800 dark:bg-white dark:hover:bg-zinc-200 text-white dark:text-black font-semibold rounded-lg text-sm px-5 py-2.5 text-center transition-all disabled:opacity-70 disabled:cursor-not-allowed mt-2"
            >
              <span v-if="isLoading" class="flex items-center gap-2">
                <svg
                  class="animate-spin h-4 w-4"
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
                Отправка...
              </span>
              <span v-else>Отправить ссылку</span>
            </button>

            <button
              type="button"
              @click="toggleForgotPasswordMode"
              class="w-full flex items-center cursor-pointer justify-center text-zinc-500 hover:text-zinc-900 dark:text-zinc-400 dark:hover:text-white font-medium text-sm px-5 py-2.5 text-center transition-all mt-2"
            >
              Назад ко входу
            </button>
          </form>
        </div>

        <div v-else>
          <form @submit.prevent="handleManualLogin" class="space-y-5">
            <div class="space-y-1.5">
              <label
                for="email"
                class="block text-sm font-medium text-zinc-700 dark:text-zinc-300 ml-1 text-left"
                >Почта</label
              >
              <div class="relative group">
                <input
                  id="email"
                  v-model="form.email"
                  type="email"
                  required
                  placeholder="name@example.com"
                  class="w-full bg-zinc-50 dark:bg-[#27272a] border border-zinc-300 dark:border-[#3f3f46] text-zinc-900 dark:text-white text-sm rounded-lg focus:ring-2 focus:ring-zinc-500 dark:focus:ring-white/20 focus:border-zinc-500 block p-2.5 pl-10 placeholder-zinc-400 dark:placeholder-zinc-600 transition-all outline-none"
                />
                <div
                  class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-zinc-400 dark:text-zinc-500 group-focus-within:text-zinc-800 dark:group-focus-within:text-white transition-colors"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <rect width="20" height="16" x="2" y="4" rx="2" />
                    <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
                  </svg>
                </div>
              </div>
            </div>

            <div class="space-y-1.5">
              <div class="flex items-center justify-between ml-1">
                <label
                  for="password"
                  class="block text-sm font-medium text-zinc-700 dark:text-zinc-300"
                  >Пароль</label
                >
                <button
                  type="button"
                  @click="toggleForgotPasswordMode"
                  class="text-xs text-zinc-500 hover:text-zinc-900 dark:text-zinc-500 dark:hover:text-zinc-300 transition-colors"
                >
                  Забыли пароль?
                </button>
              </div>

              <div class="relative group">
                <input
                  id="password"
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  required
                  autocomplete="off"
                  placeholder="••••••••"
                  class="w-full bg-zinc-50 dark:bg-[#27272a] border border-zinc-300 dark:border-[#3f3f46] text-zinc-900 dark:text-white text-sm rounded-lg focus:ring-2 focus:ring-zinc-500 dark:focus:ring-white/20 focus:border-zinc-500 block p-2.5 pl-10 pr-10 placeholder-zinc-400 dark:placeholder-zinc-600 transition-all outline-none"
                />

                <div
                  class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-zinc-400 dark:text-zinc-500 group-focus-within:text-zinc-800 dark:group-focus-within:text-white transition-colors"
                >
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <rect width="18" height="11" x="3" y="11" rx="2" ry="2" />
                    <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                  </svg>
                </div>

                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute inset-y-0 right-0 flex items-center pr-3 text-zinc-400 hover:text-zinc-800 dark:text-zinc-500 dark:hover:text-white cursor-pointer transition-colors focus:outline-none"
                >
                  <svg
                    v-if="!showPassword"
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M2 12s3-7 10-7 10 7 10 7-3 7-10 7-10-7-10-7Z" />
                    <circle cx="12" cy="12" r="3" />
                  </svg>

                  <svg
                    v-else
                    xmlns="http://www.w3.org/2000/svg"
                    width="18"
                    height="18"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                  >
                    <path d="M9.88 9.88a3 3 0 1 0 4.24 4.24" />
                    <path
                      d="M10.73 5.08A10.43 10.43 0 0 1 12 5c7 0 10 7 10 7a13.16 13.16 0 0 1-1.67 2.68"
                    />
                    <path
                      d="M6.61 6.61A13.526 13.526 0 0 0 2 12s3 7 10 7c.68 0 1.35-.06 1.99-.17"
                    />
                    <line x1="1" x2="23" y1="1" y2="23" />
                  </svg>
                </button>
              </div>
            </div>

            <button
              type="submit"
              :disabled="isLoading"
              class="w-full flex items-center cursor-pointer justify-center bg-zinc-900 hover:bg-zinc-800 dark:bg-white dark:hover:bg-zinc-200 text-white dark:text-black font-semibold rounded-lg text-sm px-5 py-2.5 text-center transition-all disabled:opacity-70 disabled:cursor-not-allowed mt-2"
            >
              <span v-if="isLoading" class="flex items-center gap-2">
                <svg
                  class="animate-spin h-4 w-4"
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
                Вход...
              </span>
              <span v-else>Войти</span>
            </button>
          </form>

          <div class="mt-8">
            <div class="relative flex py-2 items-center">
              <div class="flex-grow border-t border-zinc-200 dark:border-[#3f3f46]"></div>
              <span class="flex-shrink-0 mx-4 text-xs text-zinc-500 uppercase">Или через</span>
              <div class="flex-grow border-t border-zinc-200 dark:border-[#3f3f46]"></div>
            </div>

            <div class="flex justify-center items-center gap-4 mt-4">
              <button
                @click="handleGoogleLogin"
                class="cursor-pointer w-12 h-12 flex items-center justify-center rounded-full hover:scale-110 transition-all duration-300 bg-zinc-300 hover:bg-zinc-400 dark:bg-transparent dark:hover:bg-transparent"
                title="Войти через Google"
              >
                <img :src="GoogleIcon" alt="Google" class="w-5 h-5 dark:invert" />
              </button>
              <button
                @click="handleDiscordLogin"
                class="cursor-pointer w-12 h-12 flex items-center justify-center rounded-full hover:scale-110 transition-all duration-300 bg-zinc-300 hover:bg-zinc-400 dark:bg-transparent dark:hover:bg-transparent"
                title="Войти через Discord"
              >
                <img :src="DiscordIcon" alt="Discord" class="w-8 h-8 dark:invert" />
              </button>
            </div>
          </div>

          <p
            class="flex gap-2 justify-center mt-8 text-center text-sm text-zinc-600 dark:text-zinc-500"
          >
            <span>Нет аккаунта?</span>
            <RouterLink to="/register">
              <span class="font-medium text-zinc-900 dark:text-white hover:underline transition-all"
                >Регистрация</span
              >
            </RouterLink>
          </p>
        </div>
      </div>
    </div>
  </main>
</template>

<style lang="css" scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
