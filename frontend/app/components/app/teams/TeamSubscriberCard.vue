<script setup lang="ts">
import { staticImage } from '@/utils/str'
import type { TeamSubscriber } from '@/types/frontend/teams'
import { parseDateToLocale } from '@/utils/date'

defineProps<{
  subscriber: TeamSubscriber
}>()
</script>

<template>
  <div
    class="flex items-center gap-4 bg-zinc-50 dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4 transition-colors hover:bg-zinc-100 dark:hover:bg-zinc-800/40 cursor-pointer"
  >
    <NuxtLink :to="`/profile/${subscriber.user.id}`" class="flex items-center gap-4 w-full">
      <img
        v-if="subscriber.user.picture"
        :src="staticImage(subscriber.user.picture)"
        alt="Avatar"
        class="w-12 h-12 rounded-full object-cover"
      />
      <Icon v-else name="ph:user-bold" size="20" class="text-zinc-700/40" />

      <div class="flex-1 min-w-0 flex flex-col justify-center">
        <span class="font-bold text-sm text-zinc-900 dark:text-zinc-100 truncate transition-colors">
          {{ subscriber.user.name }}
        </span>
        <div class="flex items-center gap-0.5 mt-1">
          <span class="text-xs font-semibold text-yellow-600 dark:text-yellow-500">
            Подписан с
          </span>
          <span class="text-xs font-semibold text-zinc-800 dark:text-zinc-200">
            {{ parseDateToLocale(subscriber.subscribed_at) }}
          </span>
          <span
            v-if="subscriber.user.user_level"
            class="text-[10px] font-bold px-2 py-0.5 ml-1 rounded-full bg-zinc-200 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400"
          >
            Ур. {{ subscriber.user.user_level.level }}
          </span>
        </div>
      </div>
    </NuxtLink>
  </div>
</template>
