<script setup lang="ts">
import { staticImage } from '@/utils/str'
import UserTagId from '~/components/features/UserTagID/UserTagId.vue'
import { Separator } from '@kyokusu-ui/vue'
import { useProfile } from '@/composables/api/profile/useProfile'
import type { GetUserDto } from '@/types/backend/user'

const props = defineProps<{
  profileData: GetUserDto | null
}>()

const { userRoleColor, isLogin } = useProfile()
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-radial-[at_center] dark:from-zinc-900 cursor-default dark:to-zinc-950 dark:to-90% text-zinc-900 dark:text-zinc-200 transition-colors duration-300 font-sans"
  >
    <div
      class="h-48 md:h-64 w-full bg-linear-to-r from-zinc-800 to-zinc-700 relative overflow-hidden"
    >
      <img
        v-if="profileData?.banner"
        :src="staticImage(profileData.banner)"
        alt="Profile Banner"
        class="absolute inset-0 w-full h-full object-cover opacity-50 blur-sm"
      />
      <div class="absolute inset-0 bg-black/40"></div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-10">
      <!-- Header -->
      <div class="relative -mt-16 md:-mt-20 mb-8 flex flex-col md:flex-row items-end gap-6">
        <div class="relative shrink-0">
          <div
            class="w-32 h-32 md:w-40 md:h-40 rounded-2xl border-4 border-white dark:border-zinc-900 bg-zinc-800 overflow-hidden shadow-xl"
          >
            <img
              :src="staticImage(profileData?.picture || '')"
              class="w-full h-full object-cover"
              alt="Avatar"
            />
          </div>
          <div
            class="absolute bottom-2 right-2 w-5 h-5 border-4 border-white dark:border-zinc-900 rounded-full transition-colors duration-300"
            :class="isLogin ? 'bg-green-500' : 'bg-red-500'"
          ></div>
        </div>

        <div class="flex-1 pb-1 h-20 w-full md:w-auto">
          <div class="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
            <div>
              <h1
                class="text-3xl md:text-4xl font-bold text-zinc-900 dark:text-white flex items-center gap-3"
              >
                {{ profileData?.name }}
                <span
                  :class="[
                    'text-xs px-2 py-0.5 rounded border uppercase tracking-wider font-semibold',
                    userRoleColor,
                  ]"
                >
                  {{ profileData?.role || 'User' }}
                </span>
              </h1>
              <div
                class="flex flex-row gap-4 justify-start items-center text-sm text-zinc-500 dark:text-zinc-400 mt-2"
              >
                <UserTagId :userID="profileData?.id!" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <Separator class="mb-8" />

      <!-- Private Profile Placeholder -->
      <div
        class="flex flex-col items-center justify-center py-24 text-center bg-white dark:bg-zinc-900/30 border border-zinc-200 dark:border-zinc-800 rounded-xl shadow-sm"
      >
        <div
          class="w-20 h-20 bg-zinc-100 dark:bg-zinc-800 rounded-full flex items-center justify-center mb-6"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="w-10 h-10 text-zinc-400"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          >
            <rect x="3" y="11" width="18" height="11" rx="2" ry="2"></rect>
            <path d="M7 11V7a5 5 0 0 1 10 0v4"></path>
          </svg>
        </div>
        <h3 class="text-xl font-bold text-zinc-900 dark:text-white mb-2">Это закрытый профиль</h3>
        <p class="text-zinc-500 dark:text-zinc-400 max-w-sm">
          Пользователь ограничил доступ к информации о себе.
        </p>
      </div>
    </div>
  </div>
</template>
