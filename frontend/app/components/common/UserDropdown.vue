<script setup lang="ts">
import { onClickOutside } from '@vueuse/core'
import { ref } from 'vue'
import { staticImage } from '@/utils/str'
import { useRolePermissions } from '@/composables/api/role/useRolePermissions'
import { KyokusuAppRole } from '@/types/enums/role-enum'
import type { GetUserDto } from '@/types/backend/user'

const props = defineProps<{
  user: GetUserDto | null
  isAuthenticated: boolean
}>()

const emit = defineEmits<{
  logout: []
  login: []
}>()

const { hasPermission } = useRolePermissions()

const isOpen = ref(false)
const isContentSubmenuOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)

onClickOutside(dropdownRef, () => {
  isOpen.value = false
  isContentSubmenuOpen.value = false
})

const handleLogout = () => {
  isOpen.value = false
  emit('logout')
}

const goToLogin = () => {
  emit('login')
}
</script>

<template>
  <div class="relative" ref="dropdownRef">
    <button
      v-if="isAuthenticated && user"
      @click="isOpen = !isOpen"
      class="flex items-center focus:outline-none"
    >
      <div
        class="h-10 w-10 rounded-full cursor-pointer bg-zinc-200 dark:bg-zinc-700 flex items-center justify-center overflow-hidden border-2 border-transparent hover:border-zinc-400 dark:hover:border-zinc-300/50 transition-colors"
      >
        <img
          :src="staticImage(user.picture)"
          class="w-full h-full object-cover"
          alt="Avatar"
        />
      </div>
    </button>

    <button
      v-else
      @click="goToLogin"
      class="hidden md:flex justify-center gap-2 items-center bg-zinc-300 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 transition-colors px-6 py-3 rounded-2xl cursor-pointer"
    >
      <Icon name="ph:sign-in-bold" size="20" />
      <span class="font-medium">Войти</span>
    </button>

    <Transition name="fade">
      <div
        v-if="isOpen"
        class="absolute right-0 mt-3 w-56 bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-700 rounded-xl shadow-xl py-2 px-2 flex flex-col z-60"
      >
        <div
          class="px-4 py-2 border-b border-zinc-200 dark:border-zinc-700/50 mb-1 text-center"
        >
          <p class="text-sm font-semibold text-zinc-900 dark:text-white truncate">
            {{ user?.name }}
          </p>
          <p class="text-xs text-zinc-500 truncate">{{ user?.email }}</p>
        </div>

        <NuxtLink
          :to="`/profile/${user?.id}`"
          class="flex justify-between items-center rounded-full px-4 py-2 text-sm text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-300 dark:hover:bg-zinc-700 dark:hover:text-white transition-colors"
          @click="isOpen = false"
        >
          <span>Профиль</span>
          <Icon name="ph:user-bold" size="18" />
        </NuxtLink>

        <NuxtLink
          to="/profile/settings"
          class="flex justify-between items-center rounded-full px-4 py-2 text-sm text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-300 dark:hover:bg-zinc-700 dark:hover:text-white transition-colors"
          @click="isOpen = false"
        >
          <span>Настройки</span>
          <Icon name="ph:gear-six-bold" size="18" />
        </NuxtLink>

        <NuxtLink
          v-if="hasPermission(KyokusuAppRole.MODERATOR)"
          to="/dashboard"
          class="flex justify-between items-center rounded-full px-4 py-2 text-sm text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-300 dark:hover:bg-zinc-700 dark:hover:text-white transition-colors"
          @click="isOpen = false"
        >
          <span>Панель управления</span>
          <Icon name="ph:chalkboard-simple-bold" size="18" />
        </NuxtLink>

        <div
          class="relative"
          @mouseenter="isContentSubmenuOpen = true"
          @mouseleave="isContentSubmenuOpen = false"
        >
          <div
            v-if="isAuthenticated"
            class="flex justify-between items-center rounded-full cursor-pointer px-4 py-2 text-sm text-zinc-600 hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-300 dark:hover:bg-zinc-700 dark:hover:text-white transition-colors"
          >
            <span>Добавить контент</span>
            <Icon name="ph:palette-bold" size="18" />
          </div>

          <Transition name="fade">
            <div
              v-if="isContentSubmenuOpen"
              class="absolute top-0 left-full ml-2 w-48 bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-700 rounded-xl shadow-xl py-2 px-2 z-70"
            >
              <NuxtLink
                v-if="hasPermission(KyokusuAppRole.MODERATOR)"
                to="/novela/add"
                class="flex justify-between items-center rounded-full px-4 py-2 text-sm text-zinc-600 hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                @click="isOpen = false"
              >
                <span>Новелла</span>
                <Icon name="ph:book-open-bold" size="18" />
              </NuxtLink>

              <NuxtLink
                v-if="hasPermission(KyokusuAppRole.MODERATOR)"
                to="/author/add"
                class="flex justify-between items-center rounded-full px-4 py-2 text-sm text-zinc-600 hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                @click="isOpen = false"
              >
                <span>Автора</span>
                <Icon name="ph:pen-nib-bold" size="18" />
              </NuxtLink>

              <NuxtLink
                to="/team/add"
                class="flex justify-between items-center rounded-full px-4 py-2 text-sm text-zinc-600 hover:bg-zinc-100 dark:text-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                @click="isOpen = false"
              >
                <span>Команду</span>
                <Icon name="ph:users-three-bold" size="18" />
              </NuxtLink>
            </div>
          </Transition>
        </div>

        <div class="h-px bg-zinc-200 dark:bg-zinc-700/50 my-1"></div>

        <button
          @click="handleLogout"
          class="flex justify-between items-center rounded-full cursor-pointer px-4 py-2 text-sm text-zinc-600 hover:bg-red-50 hover:text-red-500 dark:text-zinc-300 dark:hover:bg-red-500/10 dark:hover:text-red-300 transition-colors"
        >
          <span>Выйти</span>
          <Icon name="ph:sign-out-bold" size="18" />
        </button>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}
</style>
