<script setup lang="ts">
import { ref } from 'vue'
import { staticImage, roleConv } from '@/utils/str'
import { useProfile } from '@/composables/api/profile/useProfile'
import { useUserActivity } from '@/composables/api/profile/useUserActivity'
import { useInterfaceSettings } from '@/composables/api/settings/useInterfaceSettings'

import { Separator, ModalWindow } from '@kyokusu-ui/vue'
import TabActivity from '@/components/app/profile/activity/TabActivity.vue'
import TagSelector from '@/components/features/TagSelector/TagSelector.vue'
import UserExperiance from '@/components/features/UserExperience/UserExperience.vue'
import ExperienceInfo from './experience/ExperienceInfo.vue'
import UserTagId from '@/components/features/UserTagID/UserTagId.vue'
import TabBookmarks from './tabs/TabBookmarks.vue'
import TabUserComments from './tabs/TabUserComments.vue'
import TabTeams from './tabs/TabTeams.vue'
import type { KyokusuAppRole } from '@/types/enums/role-enum'
import { TabProfile } from '@/types/enums/tab-profile'

const { profileData, accountCreated, profileTabs, userRoleColor, userGender, lastLogin, isLogin } =
  useProfile()

const { activities, fetchActivities, isLoadingActivities } = useUserActivity()
const { isShowTag, syncSettingWithBackend } = useInterfaceSettings()
const activeTab = ref<TabProfile>(TabProfile.Overview)
const isModalOpen = ref(false)

const tabComponents: Record<string, any> = {
  [TabProfile.Overview]: TabActivity,
  [TabProfile.Bookmarks]: TabBookmarks,
  [TabProfile.Comments]: TabUserComments,
  [TabProfile.Teams]: TabTeams,
}

const activeComponent = computed(() => tabComponents[activeTab.value])

const activeTabProps = computed(() => {
  if (activeTab.value === TabProfile.Overview) {
    return { activities: activities.value, isLoading: isLoadingActivities.value }
  }
  if (
    activeTab.value === TabProfile.Bookmarks ||
    activeTab.value === TabProfile.Comments ||
    activeTab.value === TabProfile.Teams
  ) {
    return { userId: profileData.value!.id }
  }
  return {}
})

const { pending } = await useAsyncData('profile-init', async () => {
  await Promise.all([fetchActivities(profileData.value!.id), syncSettingWithBackend()])
  return { success: true }
})
</script>

<template>
  <div
    class="min-h-screen flex flex-col relative overflow-hidden font-sans transition-all duration-300"
    :class="{ 'opacity-100': !pending, 'opacity-0': pending }"
  >
    <div
      class="h-48 md:h-64 w-full bg-linear-to-r from-zinc-800 to-zinc-700 relative overflow-hidden group"
    >
      <img
        v-if="profileData?.banner"
        :src="staticImage(profileData.banner)"
        alt="Profile Banner"
        class="absolute inset-0 w-full h-full transition-transform duration-700"
      />
      <div class="absolute inset-0 bg-black/20"></div>
      <div
        class="absolute -bottom-10 -right-10 w-64 h-64 bg-white/5 rounded-full blur-3xl pointer-events-none"
      ></div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-10 w-full">
      <div
        class="relative -mt-16 md:-mt-20 mb-8 flex flex-col md:flex-row items-start md:items-end gap-6"
      >
        <div class="relative shrink-0">
          <div
            class="w-32 h-32 md:w-40 md:h-40 rounded-2xl border-4 border-white dark:border-zinc-900 bg-zinc-800 overflow-hidden shadow-xl"
          >
            <img
              v-if="profileData?.picture"
              :src="staticImage(profileData.picture)"
              class="w-full h-full object-cover"
              alt="Avatar"
            />
          </div>
          <div
            class="absolute bottom-2 right-2 w-5 h-5 border-4 border-white dark:border-zinc-900 rounded-full transition-colors duration-300"
            :class="isLogin ? 'bg-green-500' : 'bg-red-500'"
          ></div>
        </div>

        <div class="flex-1 w-full pt-16 md:pt-20">
          <div
            class="flex flex-col md:flex-row justify-between items-start md:items-center gap-6 md:gap-4"
          >
            <div class="w-full md:w-auto">
              <h1
                class="text-3xl md:text-4xl font-bold text-zinc-900 dark:text-white flex flex-wrap items-center gap-3"
              >
                {{ profileData?.name || 'Загрузка...' }}
                <span
                  v-if="profileData?.role"
                  :class="[
                    'text-[10px] md:text-xs px-2 py-0.5 rounded border uppercase tracking-wider font-semibold',
                    userRoleColor,
                  ]"
                >
                  {{ roleConv(profileData.role as KyokusuAppRole, 'ru') }}
                </span>
              </h1>

              <div class="flex flex-wrap items-center gap-4 mt-3 md:mt-2">
                <UserTagId :userID="profileData?.id!" />
                <TagSelector
                  v-if="isShowTag && profileData?.active_tag"
                  v-model="profileData.active_tag"
                  :tags="profileData.tags"
                />

                <div class="flex flex-row items-center gap-3">
                  <UserExperiance
                    v-if="profileData?.user_level"
                    :level="profileData.user_level.level"
                    :currentExp="profileData.user_level.experience"
                    :expToNextLevel="profileData.user_level.xp_needed_for_next"
                    :levelTitle="profileData.user_level.level_title"
                  />

                  <button
                    class="flex justify-center cursor-pointer items-center h-8 w-8 mt-2 rounded-full bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                    @click="isModalOpen = true"
                  >
                    <Icon name="ph:info-bold" size="20" />
                  </button>
                </div>
              </div>
            </div>

            <NuxtLink
              to="/profile/settings"
              class="w-full md:w-auto flex items-center justify-center gap-2 px-5 py-2.5 md:py-2 bg-zinc-400 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 rounded-xl md:rounded-lg text-sm font-medium transition-colors border border-zinc-300 dark:border-zinc-700"
            >
              <Icon name="ph:gear-six-bold" size="18" />
              <span>Редактировать</span>
            </NuxtLink>
          </div>
        </div>
      </div>

      <Separator class="mb-8" />

      <div class="grid grid-cols-1 lg:grid-cols-12 gap-8">
        <div class="lg:col-span-4 xl:col-span-3 space-y-6 cursor-default">
          <div
            class="bg-white dark:bg-zinc-900/5 border border-zinc-200 dark:border-zinc-800 rounded-xl p-6 shadow-sm"
          >
            <div class="flex flex-col items-center">
              <h3 class="text-lg font-semibold text-zinc-900 dark:text-white mb-4">О себе</h3>
              <p class="text-sm text-zinc-600 dark:text-zinc-400 leading-relaxed text-center">
                {{ profileData?.about || 'Пользователь не оставил описания.' }}
              </p>
            </div>

            <div class="mt-6 pt-6 border-t border-zinc-200 dark:border-zinc-800 space-y-3 text-sm">
              <div class="flex justify-between">
                <span class="text-zinc-500">Регистрация</span>
                <span class="text-zinc-900 dark:text-zinc-300 font-medium">{{
                  accountCreated
                }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-zinc-500">Последний вход</span>
                <span class="text-zinc-900 dark:text-zinc-300 font-medium">{{ lastLogin }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-zinc-500">Пол</span>
                <span class="text-zinc-900 dark:text-zinc-300 font-medium capitalize">{{
                  userGender
                }}</span>
              </div>
            </div>
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div
              class="bg-zinc-100 dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 p-4 rounded-xl text-center"
            >
              <div class="text-2xl font-bold text-zinc-900 dark:text-white">
                {{ profileData?.user_stats?.read_chapters || 0 }}
              </div>
              <div class="text-[10px] text-zinc-500 uppercase tracking-widest mt-1">
                <p>Прочитано</p>
                <p>Глав</p>
              </div>
            </div>
            <div
              class="bg-zinc-100 dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 p-4 rounded-xl text-center"
            >
              <div class="text-2xl font-bold text-zinc-900 dark:text-white">
                {{ profileData?.user_stats?.total_comments || 0 }}
              </div>
              <div class="text-[10px] text-zinc-500 uppercase tracking-widest mt-1">
                <p>Комментариев</p>
              </div>
            </div>
          </div>
        </div>

        <div class="lg:col-span-8 xl:col-span-9">
          <div
            class="flex gap-1 border-b border-zinc-200 dark:border-zinc-800 mb-6 overflow-x-auto no-scrollbar"
          >
            <button
              v-for="tab in profileTabs"
              :key="tab.id"
              @click="activeTab = tab.id"
              class="px-4 py-3 text-sm cursor-pointer font-medium border-b-2 transition-all whitespace-nowrap"
              :class="
                activeTab === tab.id
                  ? 'border-zinc-900 dark:border-white text-zinc-900 dark:text-white'
                  : 'border-transparent text-zinc-500 hover:text-zinc-700 dark:hover:text-zinc-300'
              "
            >
              {{ tab.label }}
            </button>
          </div>

          <div class="min-h-75">
            <component v-if="activeComponent" :is="activeComponent" v-bind="activeTabProps" />
            <div
              v-else
              class="p-12 rounded-xl border border-dashed border-zinc-300 dark:border-zinc-700 text-zinc-500 text-center"
            >
              Здесь пока ничего нет
            </div>
          </div>
        </div>
      </div>
    </div>

    <ModalWindow
      v-if="isModalOpen"
      v-model="isModalOpen"
      title="Информация об уровнях"
      width="w-full max-w-2xl"
    >
      <ExperienceInfo />
    </ModalWindow>
  </div>
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
