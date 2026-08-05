<script setup lang="ts">
import { useTeam } from '~/composables/api/teams/useTeams'
import { parseDateToLocale } from '@/utils/date'
import { staticImage } from '@/utils/str'
import { Tooltip } from '@kyokusu-ui/vue'

const props = defineProps<{
  request: any
  slug: string
}>()

const emit = defineEmits(['accepted', 'rejected'])

const { acceptJoinRequest, rejectJoinRequest, isLoading } = useTeam()

const accept = async () => {
  const success = await acceptJoinRequest(props.slug, props.request.user.id)
  if (success) {
    emit('accepted', props.request.user.id)
  }
}

const reject = async () => {
  const success = await rejectJoinRequest(props.slug, props.request.user.id)
  if (success) {
    emit('rejected', props.request.user.id)
  }
}
</script>

<template>
  <div
    class="flex items-center w-90 gap-4 bg-zinc-50 dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4"
  >
    <NuxtLink :to="`/profile/${request.user.id}`" class="flex items-center gap-4 flex-1 min-w-0">
      <img
        v-if="request.user.picture"
        :src="staticImage(request.user.picture)"
        alt="Avatar"
        class="w-12 h-12 rounded-full object-cover shrink-0"
      />

      <Icon v-else name="ph:user-bold" size="20" class="text-zinc-700/40 shrink-0" />

      <div class="flex-1 min-w-0">
        <div class="font-bold text-sm text-zinc-900 dark:text-zinc-100 truncate">
          {{ request.user.name }}
        </div>

        <div class="flex items-center gap-2 mt-1 flex-wrap">
          <span class="text-xs font-semibold text-yellow-600 dark:text-yellow-500">
            Заявка от
          </span>

          <span class="text-xs font-semibold text-zinc-800 dark:text-zinc-200">
            {{ parseDateToLocale(request.created_at) }}
          </span>

          <span
            v-if="request.user.user_level"
            class="text-[10px] font-bold px-2 py-0.5 rounded-full bg-zinc-200 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400"
          >
            Ур. {{ request.user.user_level.level }}
          </span>
        </div>
      </div>
    </NuxtLink>

    <div class="flex items-center gap-2 shrink-0">
      <Tooltip text="Принять">
        <button @click="accept" :disabled="isLoading" class="cursor-pointer">
          <Icon name="ph:check-bold" size="16" class="text-white hover:text-green-500" />
        </button>
      </Tooltip>

      <Tooltip text="Отклонить">
        <button @click="reject" :disabled="isLoading" class="cursor-pointer">
          <Icon name="ph:x-bold" size="16" class="text-white hover:text-red-500" />
        </button>
      </Tooltip>
    </div>
  </div>
</template>
