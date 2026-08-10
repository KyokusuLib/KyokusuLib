<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { ref } from 'vue'
import { staticImage } from '@/utils/str'
import LogoIcon from '@/assets/images/Kyokusu/kyokusulib_logo_2.png'
import { useAuthStore } from '@/stores/auth'
import { useRolePermissions } from '@/composables/api/role/useRolePermissions'
import { useSearch } from '@/composables/api/search/useSearch'
import { useNotifications } from '@/composables/api/notifications/useNotifications'
import SearchModal from '@/components/common/SearchModal.vue'
import UserDropdown from '@/components/common/UserDropdown.vue'
import AdditionalDropdown from '@/components/common/AdditionalDropdown.vue'
import { KyokusuAppRole } from '@/types/enums/role-enum'
import { textMaxNumValue } from '@/utils/str'
import { MAX_UNREAD_COUNT } from '@/constants/data'

const router = useRouter()
const authStore = useAuthStore()
const { user, isAuthenticated } = storeToRefs(authStore)
const { hasPermission } = useRolePermissions()
const { openSearch } = useSearch()
const { unreadCount } = useNotifications()

const isMobileMenuOpen = ref(false)

const toggleMobileMenu = () => (isMobileMenuOpen.value = !isMobileMenuOpen.value)
const closeMobileMenu = () => (isMobileMenuOpen.value = false)

const handleLogout = async () => {
  await authStore.logout()
  closeMobileMenu()
}

const goToLogin = () => {
  closeMobileMenu()
  router.push('/login')
}

const handleGlobalSearchShortcut = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    openSearch()
    closeMobileMenu()
  }
}

onMounted(async () => {
  window.addEventListener('keydown', handleGlobalSearchShortcut)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleGlobalSearchShortcut)
})
</script>

<template>
  <header
    class="w-full border-b-2 border-zinc-200 dark:border-zinc-700 bg-white dark:bg-[#0f0f0f] z-50 sticky top-0 text-zinc-900 dark:text-white transition-colors duration-300"
  >
    <div
      class="mx-auto flex items-center justify-between px-4 py-4 md:px-8 lg:px-12 max-w-screen-2xl"
    >
      <div class="flex items-center gap-4 md:gap-8 flex-1">
        <NuxtLink
          to="/"
          class="bg-white dark:bg-zinc-800 hover:dark:bg-zinc-700 rounded-full shrink-0"
          @click="closeMobileMenu"
        >
          <img :src="LogoIcon" alt="Logo" class="h-12 w-12 md:h-12 md:w-12 dark:invert" />
        </NuxtLink>

        <nav class="hidden md:flex items-center gap-2 lg:gap-4">
          <NuxtLink
            to="/catalog"
            class="flex items-center gap-2 px-6 py-3 rounded-2xl bg-zinc-300 hover:bg-zinc-200 dark:bg-zinc-800/80 dark:hover:bg-zinc-700 transition-colors font-medium text-zinc-700 hover:text-zinc-900 dark:text-zinc-200 dark:hover:text-white"
          >
            <Icon name="ph:books-bold" size="16" class="text-zinc-700 dark:text-zinc-200" />
            Каталог
          </NuxtLink>
          <NuxtLink
            to="/top"
            class="flex items-center gap-2 px-6 py-3 rounded-2xl bg-zinc-300 hover:bg-zinc-200 dark:bg-zinc-800/80 dark:hover:bg-zinc-700 transition-colors font-medium text-zinc-700 hover:text-zinc-900 dark:text-zinc-200 dark:hover:text-white"
          >
            <Icon name="ph:ranking-bold" size="16" class="text-zinc-700 dark:text-zinc-200" />
            Топы
          </NuxtLink>
          <NuxtLink
            to="/forum"
            class="flex items-center gap-2 px-6 py-3 rounded-2xl bg-zinc-300 hover:bg-zinc-200 dark:bg-zinc-800/80 dark:hover:bg-zinc-700 transition-colors font-medium text-zinc-700 hover:text-zinc-900 dark:text-zinc-200 dark:hover:text-white"
          >
            <Icon name="ph:chats-circle-bold" size="16" class="text-zinc-700 dark:text-zinc-200" />
            Форум
          </NuxtLink>

          <AdditionalDropdown />
        </nav>
      </div>

      <div class="flex items-center gap-3 md:gap-4">
        <div
          @click="openSearch"
          class="hidden w-64 md:flex items-center justify-between px-4 py-3 rounded-full bg-zinc-300 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 transition-colors cursor-pointer group"
        >
          <div class="flex items-center gap-3">
            <Icon
              name="ph:magnifying-glass-bold"
              size="20"
              class="text-zinc-500 group-hover:text-zinc-900 dark:text-zinc-400 dark:group-hover:text-white"
            />
            <span
              class="text-sm lg:text-base font-medium text-zinc-500 group-hover:text-zinc-900 dark:text-zinc-400 dark:group-hover:text-white"
              >Поиск</span
            >
          </div>
          <div
            class="hidden lg:flex items-center gap-1 opacity-70 group-hover:opacity-100 transition-opacity"
          >
            <kbd
              class="px-1.5 py-0.5 text-[10px] font-bold rounded bg-zinc-400/30 text-zinc-600 dark:text-zinc-300"
              >Ctrl</kbd
            >
            <span class="text-xs font-bold text-zinc-500">+</span>
            <kbd
              class="px-1.5 py-0.5 text-[10px] font-bold rounded bg-zinc-400/30 text-zinc-600 dark:text-zinc-300"
              >K</kbd
            >
          </div>
        </div>

        <NuxtLink to="/notifications">
          <button
            class="relative p-3 rounded-full cursor-pointer bg-zinc-300 hover:bg-zinc-200 dark:bg-zinc-800/80 dark:hover:bg-zinc-700 transition-colors flex items-center justify-center"
          >
            <Icon name="ph:bell-bold" size="22" class="text-zinc-700 dark:text-zinc-200" />
            <ClientOnly>
              <span
                v-if="unreadCount"
                class="absolute -top-0.5 -right-0.5 flex items-center justify-center min-w-4.5 h-4.5 px-1 rounded-full bg-red-600 text-[10px] font-bold text-white leading-none"
                >{{ textMaxNumValue(unreadCount, MAX_UNREAD_COUNT) }}</span
              >
            </ClientOnly>
          </button>
        </NuxtLink>

        <NuxtLink
          to="/bookmarks"
          class="hidden sm:flex items-center gap-2 px-4 py-3 rounded-full bg-zinc-300 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 transition-colors"
        >
          <Icon name="ph:bookmarks-bold" size="16" class="text-zinc-700 dark:text-zinc-200" />
          <span class="hidden lg:inline font-medium">Закладки</span>
        </NuxtLink>

        <UserDropdown
          :user="user"
          :is-authenticated="isAuthenticated"
          @logout="handleLogout"
          @login="goToLogin"
        />

        <button
          @click="toggleMobileMenu"
          class="md:hidden p-2 active:scale-95 transition-transform bg-zinc-100 dark:bg-zinc-800 rounded-2xl flex items-center justify-center"
        >
          <Icon
            :name="isMobileMenuOpen ? 'ph:x-bold' : 'ph:list-bold'"
            size="24"
            class="text-zinc-700 dark:text-zinc-200"
          />
        </button>
      </div>
    </div>

    <Transition name="slide-fade">
      <div
        v-if="isMobileMenuOpen"
        class="absolute left-0 w-full h-[calc(100vh-80px)] overflow-y-auto border-t border-zinc-200 dark:border-zinc-700 bg-white/95 dark:bg-[#0f0f0f]/95 backdrop-blur-md md:hidden shadow-xl z-50"
      >
        <nav class="flex flex-col p-4 gap-2 pb-10">
          <div
            v-if="isAuthenticated && user"
            class="flex items-center gap-3 p-3 bg-zinc-100 dark:bg-zinc-800/50 rounded-xl mb-2"
          >
            <NuxtLink
              :to="`/profile/${user.id}`"
              class="h-10 w-10 shrink-0 rounded-full overflow-hidden"
              @click="closeMobileMenu"
            >
              <img :src="staticImage(user.picture)" class="w-full h-full object-cover" />
            </NuxtLink>
            <div class="flex-1 min-w-0 text-left">
              <p class="font-medium text-zinc-900 dark:text-white truncate">{{ user.name }}</p>
              <p class="text-xs text-zinc-500 dark:text-zinc-400 truncate">{{ user.email }}</p>
            </div>
            <button @click="handleLogout" class="p-2 text-zinc-500 hover:text-red-500">
              <Icon name="ph:sign-out-bold" size="22" />
            </button>
          </div>

          <button
            v-else
            @click="goToLogin"
            class="w-full py-3 bg-zinc-100 dark:bg-zinc-800 rounded-xl font-medium mb-4"
          >
            Войти в аккаунт
          </button>

          <NuxtLink
            to="/catalog"
            class="w-full text-center py-3 text-lg font-medium hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-xl"
            @click="closeMobileMenu"
            >Каталог</NuxtLink
          >
          <NuxtLink
            to="/top"
            class="w-full text-center py-3 text-lg font-medium hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-xl"
            @click="closeMobileMenu"
            >Топы</NuxtLink
          >
          <NuxtLink
            to="/forum"
            class="w-full text-center py-3 text-lg font-medium hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-xl"
            @click="closeMobileMenu"
            >Форум</NuxtLink
          >

          <div class="w-full border-t border-zinc-200 dark:border-zinc-700 my-2 opacity-50"></div>

          <button
            @click="openSearch(); closeMobileMenu();"
            class="flex items-center gap-3 px-6 py-3 rounded-xl bg-zinc-100 dark:bg-zinc-800/50 w-full justify-center text-zinc-700 dark:text-zinc-200 cursor-pointer"
          >
            <Icon name="ph:magnifying-glass-bold" size="20" />
            <span>Поиск</span>
          </button>

          <NuxtLink
            to="/bookmarks"
            class="flex items-center gap-3 px-6 py-3 rounded-xl bg-zinc-100 dark:bg-zinc-800/50 w-full justify-center text-zinc-700 dark:text-zinc-200"
            @click="closeMobileMenu"
          >
            <Icon name="ph:bookmark-simple-bold" size="20" />
            <span>Закладки</span>
          </NuxtLink>

          <div
            v-if="hasPermission(KyokusuAppRole.MODERATOR)"
            class="flex flex-col items-center gap-2 mt-2"
          >
            <p class="text-[10px] text-zinc-500 uppercase font-black tracking-widest mb-1">
              Управление
            </p>
            <NuxtLink
              to="/dashboard"
              class="flex justify-center items-center gap-3 px-6 py-3 rounded-xl bg-zinc-50 dark:bg-zinc-900/30 w-full"
              @click="closeMobileMenu"
            >
              <Icon name="ph:layout-bold" size="20" class="opacity-70" />
              <span>Дашборд</span>
            </NuxtLink>

            <p class="text-[10px] text-zinc-500 uppercase font-black tracking-widest mb-1 mt-2">
              Добавить контент
            </p>
            <NuxtLink
              to="/novela/add"
              class="flex justify-center items-center gap-3 px-6 py-3 rounded-xl bg-zinc-50 dark:bg-zinc-900/30 w-full"
              @click="closeMobileMenu"
            >
              <Icon name="ph:book-open-bold" size="20" class="opacity-70" />
              <span>Ранобэ</span>
            </NuxtLink>
            <NuxtLink
              to="/author/add"
              class="flex justify-center items-center gap-3 px-6 py-3 rounded-xl bg-zinc-50 dark:bg-zinc-900/30 w-full"
              @click="closeMobileMenu"
            >
              <Icon name="ph:pen-nib-bold" size="20" class="opacity-70" />
              <span>Автора</span>
            </NuxtLink>
          </div>
        </nav>
      </div>
    </Transition>

    <SearchModal />
  </header>
</template>

<style scoped>
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.slide-fade-enter-from,
.slide-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
