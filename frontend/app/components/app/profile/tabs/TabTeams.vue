<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTeam } from '@/composables/api/teams/useTeams'
import type { Team } from '@/types/frontend/teams'
import { staticImage } from '@/utils/str'
import FilterBar from '~/components/ui/FilterBar.vue'
import type { FilterConfig } from '~/components/ui/FilterBar.vue'

const props = defineProps<{
  userId: number
}>()

const { getUserTeams } = useTeam()
const teams = ref<Team[]>([])
const sortedTeams = ref<Team[]>([])
const isLoading = ref(true)

const teamFilterConfig: FilterConfig = {
  defaultField: 'name',
  sortOptions: [
    {
      field: 'name',
      label: 'По названию',
      icon: 'ph:text-aa',
      compare: (a, b) => (a as Team).name.localeCompare((b as Team).name),
    },
    {
      field: 'member_count',
      label: 'По участникам',
      icon: 'ph:users',
      compare: (a, b) => (a as Team).stats.member_count - (b as Team).stats.member_count,
    },
    {
      field: 'subscribers_count',
      label: 'По подписчикам',
      icon: 'ph:users-three',
      compare: (a, b) => (a as Team).stats.subscribers_count - (b as Team).stats.subscribers_count,
    },
    {
      field: 'created_at',
      label: 'По дате создания',
      icon: 'ph:calendar',
      compare: (a, b) =>
        new Date((a as Team).created_at).getTime() - new Date((b as Team).created_at).getTime(),
    },
  ],
}

onMounted(async () => {
  isLoading.value = true
  try {
    const fetchedTeams = await getUserTeams(props.userId)
    teams.value = fetchedTeams || []
    sortedTeams.value = [...(fetchedTeams || [])]
  } catch (e) {
    console.error('Failed to fetch user teams', e)
  } finally {
    isLoading.value = false
  }
})
</script>

<template>
  <div class="space-y-6">
    <FilterBar
      v-if="!isLoading && teams.length > 0"
      :items="teams"
      v-model="sortedTeams"
      :config="teamFilterConfig"
    />

    <div v-if="isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="i in 3"
        :key="i"
        class="h-24.5 rounded-2xl bg-zinc-100 dark:bg-zinc-800/50 animate-pulse border border-zinc-200 dark:border-zinc-800"
      ></div>
    </div>

    <div
      v-else-if="sortedTeams.length > 0"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-2 gap-4"
    >
      <NuxtLink
        v-for="team in sortedTeams"
        :key="team.id"
        :to="`/team/${team.slug}`"
        class="flex items-center gap-4 p-4 rounded-2xl bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 hover:border-yellow-500/50 hover:shadow-md transition-all duration-300 group cursor-pointer"
      >
        <div
          class="w-16 h-16 shrink-0 rounded-xl bg-zinc-100 dark:bg-zinc-800 overflow-hidden border border-zinc-200/50 dark:border-zinc-700/50 flex items-center justify-center relative"
        >
          <img
            v-if="team.avatar_url"
            :src="staticImage(team.avatar_url)"
            class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-500"
          />
          <Icon
            v-else
            name="ph:users-three-bold"
            size="24"
            class="text-zinc-400 group-hover:text-yellow-500 transition-colors"
          />
          <div
            class="absolute inset-0 bg-yellow-500/10 opacity-0 group-hover:opacity-100 transition-opacity"
          ></div>
        </div>

        <div class="flex flex-col min-w-1 flex-1 justify-center">
          <span
            class="text-base font-bold text-zinc-900 dark:text-zinc-100 truncate group-hover:text-yellow-600 dark:group-hover:text-yellow-500 transition-colors"
          >
            {{ team.name }}
          </span>
          <span
            class="text-xs font-semibold tracking-wide text-zinc-500 dark:text-zinc-400 truncate mt-0.5 group-hover:text-zinc-700 dark:group-hover:text-zinc-300 transition-colors"
          >
            @{{ team.slug }}
          </span>
          <div
            class="flex gap-2 text-xs font-semibold tracking-wide text-zinc-500 dark:text-zinc-400 truncate mt-0.5 group-hover:text-zinc-700 dark:group-hover:text-zinc-300 transition-colors"
          >
            <span class="flex items-center gap-2">
              <Icon name="ph:users" size="12" />
              <span>Участников: {{ team.stats.member_count }}</span>
            </span>
            <span class="flex items-center gap-2">
              <Icon name="ph:users-three" size="12" />
              <span>Подписчиков: {{ team.stats.subscribers_count }}</span>
            </span>
          </div>
        </div>
        <Icon
          name="ph:caret-right-bold"
          size="16"
          class="text-zinc-300 dark:text-zinc-600 opacity-0 group-hover:opacity-100 -translate-x-2 group-hover:translate-x-0 transition-all"
        />
      </NuxtLink>
    </div>

    <div
      v-else
      class="flex flex-col items-center justify-center py-20 text-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-3xl bg-zinc-50/50 dark:bg-zinc-900/20"
    >
      <Icon name="ph:users-three-bold" size="48" class="text-zinc-300 dark:text-zinc-700 mb-4" />
      <p class="text-xs font-bold text-zinc-400 uppercase tracking-widest">Команд не найдено</p>
    </div>
  </div>
</template>
