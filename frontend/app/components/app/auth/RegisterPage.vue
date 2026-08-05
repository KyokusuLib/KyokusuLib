<script setup lang="ts">
import DiscordIcon from '@/assets/images/social/discord.png'
import GoogleIcon from '@/assets/images/social/google.png'
import ReferIcon from '@/assets/images/special/add.png'
import FooterApp from '@/components/common/FooterApp.vue'
import { useLogin } from '@/composables/api/auth/useLogin'
import { useRegister } from '@/composables/api/auth/useRegister'
import { ref } from 'vue'

const { handleGoogleLogin, handleDiscordLogin } = useLogin()
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const { form, isLoading, errorMessage, isSuccess, registeredEmail, handleManualRegister } =
  useRegister()
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
        <transition name="fade" mode="out-in">
          <div v-if="isSuccess" class="text-center py-4" key="success">
            <div
              class="mx-auto flex h-20 w-20 items-center justify-center rounded-full bg-zinc-100 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700 mb-6 animate-pulse"
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="40"
                height="40"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                stroke-linecap="round"
                stroke-linejoin="round"
                class="text-zinc-900 dark:text-white"
              >
                <rect width="20" height="16" x="2" y="4" rx="2" />
                <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7" />
              </svg>
            </div>

            <h2 class="text-2xl font-bold tracking-tight mb-3">Проверьте почту</h2>

            <p class="text-zinc-500 dark:text-zinc-400 text-sm mb-6 leading-relaxed">
              Мы отправили письмо с ссылкой для подтверждения на адрес:
              <br />
              <span
                class="text-zinc-900 dark:text-white font-medium bg-zinc-100 dark:bg-zinc-800 px-2 py-0.5 rounded mt-2 inline-block break-all"
                >{{ registeredEmail }}</span
              >
            </p>

            <div
              class="bg-zinc-50 dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 rounded-lg p-4 text-xs text-zinc-500 mb-8 text-left"
            >
              <p>• Не пришло письмо? Проверьте папку "Спам".</p>
              <p class="mt-1">• Ссылка действительна в течение 24 часов.</p>
            </div>

            <RouterLink
              to="/login"
              class="w-full flex items-center justify-center bg-zinc-900 hover:bg-zinc-800 dark:bg-white dark:hover:bg-zinc-200 text-white dark:text-black font-semibold rounded-lg text-sm px-5 py-2.5 text-center transition-all"
            >
              Вернуться ко входу
            </RouterLink>
          </div>

          <div v-else key="form">
            <div class="text-center mb-8">
              <div class="flex justify-center items-center gap-3 mb-3">
                <img :src="ReferIcon" class="w-6 h-6 object-contain dark:invert" />
                <h1 class="text-2xl font-bold tracking-tight uppercase">Регистрация</h1>
              </div>
              <p class="text-zinc-500 dark:text-zinc-400 text-sm">
                Создайте аккаунт, чтобы начать пользоваться сервисом
              </p>
            </div>

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

            <form @submit.prevent="handleManualRegister" class="space-y-5">
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
                <label
                  for="username"
                  class="block text-sm font-medium text-zinc-700 dark:text-zinc-300 ml-1 text-left"
                  >Логин</label
                >
                <div class="relative group">
                  <input
                    id="username"
                    v-model="form.username"
                    type="text"
                    required
                    placeholder="логин"
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
                      <path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"></path>
                      <circle cx="12" cy="7" r="4"></circle>
                    </svg>
                  </div>
                </div>
                <p class="text-xs text-stone-500 text-left text-wrap ml-[3px]">
                  Ограничение: мин. 6 символов, макс. 50 символов, латиница, мин. 1 большая буква,
                  мин. 1 маленькая.
                </p>
              </div>

              <div class="space-y-1.5">
                <label
                  for="password"
                  class="block text-sm font-medium text-zinc-700 dark:text-zinc-300 ml-1 text-left"
                  >Пароль</label
                >
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
                <p class="text-xs text-stone-500 text-left text-wrap ml-[3px]">
                  Ограничение: мин. 6 символов, макс. 50 символов, латиница, мин. 1 большая буква,
                  мин. 1 маленькая.
                </p>
              </div>

              <div class="space-y-1.5">
                <label
                  for="password_confirm"
                  class="block text-sm font-medium text-zinc-700 dark:text-zinc-300 ml-1 text-left"
                  >Повторите пароль</label
                >
                <div class="relative group">
                  <input
                    id="password_confirm"
                    v-model="form.passwordConfirm"
                    :type="showConfirmPassword ? 'text' : 'password'"
                    required
                    autocomplete="new-password"
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
                      <path
                        d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 1 1-7.778 7.778 5.5 5.5 0 0 1 7.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"
                      />
                    </svg>
                  </div>
                  <button
                    type="button"
                    @click="showConfirmPassword = !showConfirmPassword"
                    class="absolute inset-y-0 right-0 flex items-center pr-3 text-zinc-400 hover:text-zinc-800 dark:text-zinc-500 dark:hover:text-white cursor-pointer transition-colors focus:outline-none"
                  >
                    <svg
                      v-if="!showConfirmPassword"
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
                    class="animate-spin h-4 w-4 text-white dark:text-black"
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
                  Регистрация...
                </span>
                <span v-else>Зарегистрироваться</span>
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
                  class="cursor-pointer w-12 h-12 flex items-center justify-center rounded-full hover:scale-110 transition-all duration-300 bg-zinc-100 hover:bg-zinc-200 dark:bg-transparent dark:hover:bg-transparent"
                  title="Войти через Google"
                >
                  <img :src="GoogleIcon" alt="Google" class="w-5 h-5 dark:invert" />
                </button>
                <button
                  @click="handleDiscordLogin"
                  class="cursor-pointer w-12 h-12 flex items-center justify-center rounded-full hover:scale-110 transition-all duration-300 bg-zinc-100 hover:bg-zinc-200 dark:bg-transparent dark:hover:bg-transparent"
                  title="Войти через Discord"
                >
                  <img :src="DiscordIcon" alt="Discord" class="w-8 h-8 dark:invert" />
                </button>
              </div>
            </div>

            <p
              class="flex gap-2 justify-center mt-8 text-center text-sm text-zinc-600 dark:text-zinc-500"
            >
              <span>Уже есть аккаунт?</span>
              <RouterLink to="/login">
                <span
                  class="font-medium text-zinc-900 dark:text-white hover:underline transition-all"
                  >Войти</span
                >
              </RouterLink>
            </p>
          </div>
        </transition>
      </div>
    </div>
  </main>
  <FooterApp />
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
