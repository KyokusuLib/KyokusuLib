<script setup lang="ts">
import { useTeam } from '@/composables/api/teams/useTeams'
import type { TeamSubscriber } from '@/types/frontend/teams'
import TeamSubscriberCard from './TeamSubscriberCard.vue'

const props = defineProps<{
  slug: string
}>()

const { getTeamSubscribers } = useTeam()

const { data: subscribers, pending: isLoading } = useAsyncData<TeamSubscriber[]>(
  `team-subscribers-${props.slug}-preview`,
  async () => {
    try {
      return await getTeamSubscribers(props.slug, 6, 0)
    } catch (e) {
      console.error('Failed to fetch team subscribers preview', e)
      return []
    }
  },
  { default: () => [] },
)
</script>

<template>
  <div
    v-if="isLoading || subscribers.length > 0"
    class="bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-3xl p-6 sm:p-8 shadow-sm w-full mt-6"
  >
    <div class="flex items-end justify-between mb-6">
      <h2
        class="flex items-center gap-3 text-xl md:text-2xl font-black tracking-tight text-zinc-900 dark:text-white uppercase italic"
      >
        <div class="flex flex-col items-start gap-1">
          Подписчики
          <div class="h-1 w-46 bg-yellow-500 rounded-full"></div>
        </div>
      </h2>
      <div class="flex justify-center sm:justify-end">
        <NuxtLink
          :to="`/team/${slug}/subscribers`"
          class="text-center px-4 py-2 bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-300 font-bold rounded-full transition-colors text-[12px] duration-300 flex items-center gap-2"
        >
          Все подписчики
          <Icon name="ph:arrow-right-bold" />
        </NuxtLink>
      </div>
    </div>

    <div class="mt-2">
      <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        <div
          v-for="i in 6"
          :key="i"
          class="h-20 rounded-2xl bg-zinc-100 dark:bg-zinc-800/50 animate-pulse"
        ></div>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
        <TeamSubscriberCard
          v-for="subscriber in subscribers"
          :key="subscriber.user.id"
          :subscriber="subscriber"
        />
      </div>
    </div>
  </div>
</template>
