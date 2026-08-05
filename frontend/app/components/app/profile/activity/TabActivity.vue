<script setup lang="ts">
import { ref, computed } from 'vue'
import { formatDateUserActivity } from '@/utils/str'
import ActivityIcon from './ActivityIcon.vue'
import type { UserActivity } from '@/types/backend/user_activity'
import { useActivityConfig } from '@/composables/api/profile/useActivityConfig'
import { ACTIVITY_PAGE_SIZE } from '@/constants/data'
import { naviagateFromActivity } from '~/constants/urls'

interface Props {
  activities: UserActivity[]
  isLoading?: boolean
}

const props = defineProps<Props>()
const { getActivityConfig } = useActivityConfig()

const isExpanded = ref(false)

const processedActivities = computed(() => {
  return props.activities.map((activity) =>
    Object.assign(Object.create(activity), { ui: getActivityConfig(activity) }),
  )
})

const visibleActivities = computed(() => {
  if (isExpanded.value) return processedActivities.value
  return processedActivities.value.slice(0, ACTIVITY_PAGE_SIZE)
})
</script>

<template>
  <div
    class="space-y-4 bg-zinc-50 p-4 border border-zinc-200 dark:border-zinc-700 rounded-2xl dark:bg-radial-[at_center] dark:from-zinc-900 dark:to-zinc-950 dark:to-90% transition-colors duration-300 min-h-[300px]"
  >
    <div v-if="isLoading" class="space-y-3">
      <div
        v-for="i in 3"
        :key="i"
        class="flex items-center gap-4 p-4 border border-zinc-200 dark:border-zinc-800 rounded-xl animate-pulse"
      >
        <div class="shrink-0 w-10 h-10 rounded-full bg-zinc-200 dark:bg-zinc-800"></div>
        <div class="flex-1 space-y-2">
          <div class="h-4 bg-zinc-200 dark:bg-zinc-800 rounded w-1/3"></div>
          <div class="h-3 bg-zinc-200 dark:bg-zinc-800 rounded w-3/4"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div
      v-else-if="activities.length === 0"
      class="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed border-zinc-300 dark:border-zinc-800 rounded-xl h-full"
    >
      <div
        class="w-12 h-12 rounded-full bg-zinc-200 dark:bg-zinc-800 flex items-center justify-center mb-3"
      >
        <ActivityIcon name="activity" class="text-zinc-400" />
      </div>
      <p class="text-zinc-500 dark:text-zinc-400 font-medium">История активности пуста</p>
    </div>

    <!-- Activities List -->
    <div v-else class="space-y-3">
      <div
        class="space-y-3 transition-all duration-500 custom-scrollbar"
        :class="isExpanded ? 'max-h-[780px] overflow-y-auto pr-2' : 'max-h-full'"
      >
        <div
          v-for="item in visibleActivities"
          :key="item.id"
          @click="naviagateFromActivity(item)"
          class="group flex items-center gap-4 p-4 bg-white dark:bg-zinc-900/30 border border-zinc-200 dark:border-zinc-800 rounded-xl hover:border-zinc-300 dark:hover:border-zinc-700 transition-all duration-200 cursor-pointer"
        >
          <div
            class="shrink-0 w-10 h-10 rounded-full flex items-center justify-center transition-transform group-hover:scale-110 duration-300"
            :class="item.ui.color"
          >
            <ActivityIcon :name="item.ui.icon" />
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex justify-between items-baseline gap-2 mb-0.5">
              <h4
                class="text-sm font-bold text-zinc-900 dark:text-zinc-100 truncate group-hover:text-blue-500 transition-colors"
              >
                {{ item.ui.title }}
              </h4>
              <span class="text-[10px] uppercase tracking-wider text-zinc-400 whitespace-nowrap">
                {{ formatDateUserActivity(item.timestamp) }}
              </span>
            </div>

            <p class="text-sm text-zinc-600 dark:text-zinc-400 leading-snug line-clamp-2">
              {{ item.ui.description }}
            </p>
          </div>

          <div
            class="hidden md:block opacity-0 group-hover:opacity-100 transition-opacity translate-x-2 group-hover:translate-x-0 duration-300 text-zinc-400"
          >
            <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 5l7 7-7 7"
              />
            </svg>
          </div>
        </div>
      </div>

      <!-- Toggle Button -->
      <button
        v-if="activities.length > 3"
        @click="isExpanded = !isExpanded"
        class="w-full py-2 text-center text-xs font-bold uppercase tracking-widest text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300 transition-colors cursor-pointer mt-2"
      >
        {{ isExpanded ? 'Свернуть' : 'Вся история' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}
.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-10px);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 5px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #3f3f46;
  border-radius: 10px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #fbbf24;
}

.transition-all {
  transition-property: all;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 500ms;
}
</style>
