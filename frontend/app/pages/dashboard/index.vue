<script setup lang="ts">
import { useDashboardStore } from '@/stores/dashboard'
import { DashboardTab } from '@/types/enums/dashboard-tab'
import DashboardModeration from '@/components/app/dashboard/DashboardModeration.vue'
import DashboardParser from '@/components/app/dashboard/DashboardParser.vue'
import DashboardUsers from '@/components/app/dashboard/DashboardUsers.vue'
import { Separator } from '@kyokusu-ui/vue'
import { KyokusuAppRole } from '@/types/enums/role-enum'

definePageMeta({
  middleware: ['auth', 'role'],
  minRole: KyokusuAppRole.MODERATOR,
})

const dashboardStore = useDashboardStore()

const categories = [
  { id: DashboardTab.Moderation, label: 'Модерация', icon: 'ph:shield-check-bold' },
  { id: DashboardTab.Parser, label: 'Парсер', icon: 'ph:codepen-logo-bold' },
  { id: DashboardTab.Users, label: 'Управление пользователями', icon: 'ph:users-bold' },
]
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-radial-[at_center] dark:from-zinc-800 dark:to-zinc-950 dark:to-90% text-zinc-900 dark:text-zinc-200 transition-colors duration-300"
  >
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="mb-8">
        <div class="flex items-start sm:justify-start justify-center gap-2">
          <h1
            class="cs-text text-left text-4xl leading-xl font-semibold text-zinc-900 dark:text-white mb-4"
          >
            Панель управления
          </h1>
          <Icon name="ph:layout-bold" class="mt-2 text-yellow-500" size="24" />
        </div>
        <Separator />
      </div>

      <div class="flex flex-col md:flex-row gap-8">
        <aside class="w-full md:w-64 shrink-0 hidden lg:block">
          <div class="sticky top-24">
            <nav class="flex flex-col space-y-2">
              <button
                v-for="category in categories"
                :key="category.id"
                @click="dashboardStore.setActiveTab(category.id)"
                :class="[
                  dashboardStore.activeTab === category.id
                    ? 'bg-zinc-200 text-zinc-900 dark:bg-zinc-800 dark:text-white ring-1 ring-zinc-300 dark:ring-zinc-700'
                    : 'text-zinc-500 hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-400 dark:hover:bg-zinc-800/50 dark:hover:text-zinc-200',
                  'group flex items-center justify-between px-4 py-2 text-sm font-bold rounded-2xl transition-all duration-200 ease-in-out w-full text-left cursor-pointer',
                ]"
              >
                <span>{{ category.label }}</span>
                <Icon
                  :name="category.icon"
                  size="20"
                  :class="
                    dashboardStore.activeTab === category.id
                      ? 'text-yellow-500'
                      : 'text-zinc-400 dark:text-zinc-600'
                  "
                />
              </button>
            </nav>
          </div>
        </aside>

        <div
          class="flex-1 bg-white dark:bg-zinc-900/40 rounded-[2.5rem] border border-zinc-200 dark:border-zinc-800 p-8 min-h-150 shadow-sm transition-all duration-300"
        >
          <div v-if="dashboardStore.activeTab === DashboardTab.Moderation" class="h-full">
            <DashboardModeration />
          </div>

          <div v-else-if="dashboardStore.activeTab === DashboardTab.Parser" class="h-full">
            <DashboardParser />
          </div>

          <div v-else-if="dashboardStore.activeTab === DashboardTab.Users" class="h-full">
            <DashboardUsers />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
