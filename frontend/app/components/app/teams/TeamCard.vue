<script setup lang="ts">
import { ref, computed } from 'vue'
import { type Team } from '@/types/frontend/teams'
import { staticImage } from '@/utils/str'
import { useAuthStore } from '@/stores/auth'
import { useTeam } from '@/composables/api/teams/useTeams'
import { useRolePermissions } from '@/composables/api/role/useRolePermissions'
import { KyokusuAppRole } from '@/types/enums/role-enum'
import TeamSettings from './TeamSettings.vue'
import ModalConfirm from '~/components/common/ModalConfirm.vue'
import { Tooltip } from '@kyokusu-ui/vue'

const props = defineProps<{
  team: Team
}>()

const emit = defineEmits(['updated', 'join', 'leave'])

const authStore = useAuthStore()
const { joinTeam, leaveTeam, subscribeTeam, unsubscribeTeam, isLoading } = useTeam()
const { hasPermission } = useRolePermissions()

const isOpenTeamSettings = ref(false)
const isOpenLeaveConfirm = ref(false)

const textToLeave = computed(() => {
  return props.team.owner_id === authStore.user?.id
    ? `Вы владелец, если выпокинете комаду, то она исчезнет, участники и подписчики автоматически покинут её. Вы уверены, что хотите покинуть команду "${props.team.name}"? 
     Это действие нельзя будет отменить`
    : `Вы уверены, что хотите покинуть команду "${props.team.name}"?`
})

const canEdit = computed(() => {
  if (!authStore.isAuthenticated || !authStore.user) return false
  return authStore.user.id === props.team.owner_id || hasPermission(KyokusuAppRole.MODERATOR)
})

const updatedTeam = (payload: Team) => {
  emit('updated', payload)
}

const openModalLeaveConfirm = () => {
  isOpenLeaveConfirm.value = true
}

const handleJoin = async () => {
  const status = await joinTeam(props.team.slug)
  if (status === 'joined') {
    emit('join', {
      ...props.team,
      is_member: true,
      stats: {
        ...props.team.stats,
        member_count: props.team.stats.member_count + 1,
      },
    })
  } else if (status === 'request_sent') {
    emit('updated', {
      ...props.team,
      has_requested: true,
    })
  }
}

const handleLeave = async () => {
  const success = await leaveTeam(props.team.slug)
  if (success && !(props.team.owner_id === authStore.user?.id)) {
    emit('leave', {
      ...props.team,
      is_member: false,
      stats: {
        ...props.team.stats,
        member_count: props.team.stats.member_count - 1,
      },
    })
  } else if (success) {
    emit('leave', null)
  }
}

const handleSubscribe = async () => {
  if (props.team.is_subscriber) {
    await unsubscribeTeam(props.team.id)
    emit('updated', {
      ...props.team,
      is_subscriber: false,
      stats: { ...props.team.stats, subscribers_count: props.team.stats.subscribers_count - 1 },
    })
  } else {
    await subscribeTeam(props.team.id)
    emit('updated', {
      ...props.team,
      is_subscriber: true,
      stats: { ...props.team.stats, subscribers_count: props.team.stats.subscribers_count + 1 },
    })
  }
}
</script>

<template>
  <div
    class="bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-3xl shadow-sm"
  >
    <div
      class="h-32 sm:h-48 w-full bg-linear-to-r from-zinc-800 to-zinc-700 relative overflow-hidden group rounded-t-3xl"
    >
      <div class="absolute inset-0 bg-black/20 z-10 backdrop-blur-[1px]"></div>
      <img
        v-if="team.banner_url"
        :src="staticImage(team.banner_url)"
        alt="Banner Background"
        class="absolute inset-0 w-full h-full"
      />
      <img
        v-else-if="team.avatar_url"
        :src="staticImage(team.avatar_url)"
        alt="Banner Background Fallback"
        class="absolute inset-0 w-full h-full object-cover blur-xl scale-125 opacity-50"
      />
      <div class="absolute inset-0 bg-linear-to-t from-black/60 to-transparent z-10"></div>
    </div>

    <div class="px-4 sm:px-8 pb-8 relative z-20">
      <div
        class="flex flex-col sm:flex-row items-center sm:items-end gap-4 sm:gap-6 -mt-16 sm:-mt-20 mb-6"
      >
        <div
          class="w-32 h-32 sm:w-40 sm:h-40 rounded-3xl border-4 border-white dark:border-[#18181b] bg-zinc-100 dark:bg-zinc-800 overflow-hidden shadow-xl shrink-0 flex items-center justify-center"
        >
          <img
            v-if="team.avatar_url"
            :src="staticImage(team.avatar_url)"
            alt="Team picture"
            class="w-full h-full object-cover"
          />
          <Icon
            v-else
            name="ph:users-three-bold"
            size="64"
            class="text-zinc-300 dark:text-zinc-600"
          />
        </div>
        <div
          class="flex-1 w-full text-center sm:text-left pt-2 sm:pt-0 sm:pb-2 flex flex-col sm:flex-row justify-between items-center sm:items-start gap-4"
        >
          <div>
            <h1
              class="text-3xl sm:text-4xl font-black text-zinc-900 dark:text-white tracking-tight mb-2"
            >
              {{ team.name }}
            </h1>

            <div class="flex flex-wrap items-center justify-center sm:justify-start gap-2">
              <span
                class="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-zinc-300 dark:bg-zinc-800 text-xs font-bold text-zinc-600 dark:text-zinc-300"
              >
                <Icon name="ph:users-bold" size="14" class="text-zinc-400" />
                Участников: {{ team.stats.member_count }}
              </span>

              <span
                class="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-zinc-300 dark:bg-zinc-800 text-xs font-bold text-zinc-600 dark:text-zinc-300"
              >
                <Icon name="ph:star-bold" size="14" class="text-zinc-400" />
                Подписчиков: {{ team.stats.subscribers_count }}
              </span>

              <span
                v-if="team.created_at"
                class="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-yellow-50 dark:bg-yellow-500/10 text-xs font-bold text-yellow-600 dark:text-yellow-500 border border-yellow-200/50 dark:border-yellow-500/20"
              >
                <Icon name="ph:calendar-blank-bold" size="14" />
                Создана: {{ parseDateToLocale(team.created_at) }}
              </span>
            </div>
          </div>

          <div class="flex items-center gap-2 mt-4 sm:mt-0 relative">
            <Tooltip position="top" text="Настройки">
              <button
                v-if="canEdit"
                @click="isOpenTeamSettings = true"
                class="mt-2 p-1 cursor-pointer rounded-full bg-zinc-200 dark:bg-zinc-800 flex items-center text-zinc-400 hover:text-yellow-500 transition-colors duration-300"
              >
                <Icon name="ph:gear-bold" size="16" />
              </button>
            </Tooltip>

            <Tooltip
              position="top"
              :text="team.is_subscriber ? 'Отписаться' : 'Подписаться на обновления'"
            >
              <button
                class="flex items-center mt-2 p-1 cursor-pointer rounded-full transition-colors duration-300"
                :class="[
                  team.is_subscriber
                    ? 'dark:bg-amber-400 bg-amber-200  text-zinc-700  hover:text-zinc-100'
                    : 'dark:bg-zinc-800 bg-zinc-200 text-zinc-400  hover:text-yellow-500',
                ]"
                @click="handleSubscribe"
              >
                <Icon
                  :name="
                    team.is_subscriber
                      ? 'ph:arrow-bend-left-down-bold'
                      : 'ph:arrow-bend-right-up-bold'
                  "
                  size="16"
                />
              </button>
            </Tooltip>

            <button
              v-if="
                authStore.isAuthenticated &&
                team.owner_id !== authStore.user?.id &&
                !team.is_member &&
                !team.has_requested
              "
              @click="handleJoin"
              :disabled="isLoading"
              class="px-4 py-1 mt-2 text-[12px] bg-yellow-500 hover:bg-yellow-400 text-zinc-900 font-bold rounded-full transition-colors duration-300 flex items-center gap-2 cursor-pointer disabled:opacity-50"
            >
              <Icon v-if="isLoading" name="ph:spinner-gap-bold" class="animate-spin" />
              <Icon v-else name="ph:user-plus-bold" />
              Вступить
            </button>
            <button
              v-else-if="
                authStore.isAuthenticated &&
                team.owner_id !== authStore.user?.id &&
                !team.is_member &&
                team.has_requested
              "
              disabled
              class="px-4 py-1 mt-2 text-[12px] bg-zinc-200 dark:bg-zinc-800 text-zinc-500 dark:text-zinc-400 font-bold rounded-full flex items-center gap-2 opacity-70 cursor-not-allowed"
            >
              <Icon name="ph:clock-bold" />
              Заявка отправлена
            </button>
            <button
              v-else-if="authStore.isAuthenticated && team.is_member"
              @click="openModalLeaveConfirm"
              :disabled="isLoading"
              class="px-4 py-1 mt-2 text-[12px] bg-zinc-200 hover:bg-red-500 dark:bg-zinc-800 dark:hover:bg-red-500/20 text-zinc-700 dark:text-zinc-300 hover:text-white dark:hover:text-red-500 font-bold rounded-full transition-colors duration-300 flex items-center gap-2 cursor-pointer disabled:opacity-50"
            >
              <Icon v-if="isLoading" name="ph:spinner-gap-bold" class="animate-spin" />
              <Icon v-else name="ph:sign-out-bold" />
              Выйти
            </button>
          </div>
        </div>
        <TeamSettings v-model="isOpenTeamSettings" :team="team" @updated="updatedTeam" />
      </div>

      <div class="mt-8">
        <h3 class="text-lg font-bold text-zinc-900 dark:text-white mb-4 flex items-center gap-2">
          <Icon name="ph:info-bold" class="text-yellow-500" />
          О команде
        </h3>

        <div
          class="bg-zinc-200/40 dark:bg-zinc-900/50 rounded-2xl p-5 sm:p-6 border border-zinc-300 dark:border-zinc-800"
        >
          <div
            v-if="team.description"
            class="prose dark:prose-invert max-w-none text-zinc-700 dark:text-zinc-300 text-sm sm:text-base leading-relaxed whitespace-pre-line"
          >
            {{ team.description }}
          </div>
          <p v-else class="text-sm font-medium text-zinc-500 italic text-center py-4">
            Информация о команде пока не заполнена.
          </p>
        </div>
      </div>
    </div>
    <ModalConfirm
      v-model="isOpenLeaveConfirm"
      title="Выход из команды"
      confirm-text="Покинуть команду"
      cancel-text="Отмена"
      :description="textToLeave"
      @confirm="handleLeave"
    />
  </div>
</template>
