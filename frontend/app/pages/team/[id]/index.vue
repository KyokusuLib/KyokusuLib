<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useTeam } from '~/composables/api/teams/useTeams'

import TeamCard from '~/components/app/teams/TeamCard.vue'
import TeamMembersRow from '~/components/app/teams/TeamMembersRow.vue'
import TeamSubscribersRow from '~/components/app/teams/TeamSubscribersRow.vue'
import TeamJoinRequestsRow from '~/components/app/teams/TeamJoinRequestsRow.vue'
import type { CustomRoleNames, Team } from '~/types/frontend/teams'

const route = useRoute()
const slug = computed(() => route.params.id as string)

const { getTeamBySlug } = useTeam()

const { data: team, status } = await useAsyncData(`team-${slug.value}`, () =>
  getTeamBySlug(slug.value),
)

const updatedTeam = (payload: Team) => {
  if (!team.value) return
  team.value = payload
}

const updateMembers = () => {
  refreshNuxtData(`team-members-${slug.value}-preview`)
}

const joinedTeam = async (payload: Team) => {
  updatedTeam(payload)
  await refreshNuxtData(`team-members-${slug.value}-preview`)
}

const leaveTeam = async (payload: Team | null) => {
  if (payload === null) {
    return navigateTo('/profile')
  }
  updatedTeam(payload)
  await refreshNuxtData(`team-members-${slug.value}-preview`)
}

const customRoleNames = computed<CustomRoleNames[]>(() => {
  if (!team.value) return []
  return [
    { value: 'member', label: team.value.role_names.member },
    { value: 'moderator', label: team.value.role_names.moderator },
    { value: 'owner', label: team.value.role_names.owner },
  ]
})

useSeoMeta({
  title: () =>
    team.value?.name ? `${team.value.name} - Команды - KyokusuLib` : 'Команда - KyokusuLib',
  ogTitle: () => team.value?.name,
  ogImage: () => team.value?.avatar_url,
  description: () => team.value?.description || 'Страница команды переводчиков',
})
</script>

<template>
  <div
    class="min-h-screen bg-zinc-50 dark:bg-[#0f0f0f] text-zinc-900 dark:text-zinc-200 flex flex-col font-sans transition-colors duration-300"
  >
    <div class="w-full max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 space-y-6">
      <div v-if="status === 'pending'" class="flex justify-center items-center py-32">
        <Icon name="ph:spinner-gap-bold" size="48" class="animate-spin text-yellow-500 mb-4" />
      </div>
      <template v-else-if="team">
        <TeamCard :team="team" @updated="updatedTeam" @join="joinedTeam" @leave="leaveTeam" />
        <TeamJoinRequestsRow
          :slug="team.slug"
          :owner-id="team.owner_id"
          @update:members="updateMembers"
        />
        <TeamMembersRow :customRoleNames="customRoleNames" :slug="team.slug" />
        <TeamSubscribersRow :slug="team.slug" />
      </template>
      <div
        v-else
        class="flex flex-col items-center justify-center py-32 text-center border-2 border-dashed border-zinc-200 dark:border-zinc-800 rounded-3xl bg-white dark:bg-[#18181b]"
      >
        <Icon name="ph:users-three-bold" size="64" class="text-zinc-300 dark:text-zinc-700 mb-4" />
        <h2 class="text-xl font-bold text-zinc-900 dark:text-zinc-100 mb-2">Команда не найдена</h2>
        <p class="text-sm font-medium text-zinc-500">
          Запрошенная команда не существует или была удалена.
        </p>
      </div>
    </div>
  </div>
</template>
