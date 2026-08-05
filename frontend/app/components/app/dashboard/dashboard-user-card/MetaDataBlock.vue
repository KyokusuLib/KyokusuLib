<script setup lang="ts">
import type { DashboardUser } from '@/types/frontend/dashboard/users'
import type { DashboardRowUserStatus } from '@/types/enums/dashboard-table'

import { getUserStatusColor, formatStatus } from '@/utils/dashboard'

import { fmtDate, fmtDateTime } from '@/utils/date'

defineProps<{ user: DashboardUser }>()
</script>

<template>
  <div class="grid grid-cols-2 gap-3">
    <div
      class="flex flex-col items-start bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700"
    >
      <div class="flex gap-2 mb-1">
        <span class="text-[10px] font-bold uppercase tracking-widest text-zinc-400">Статус</span>
        <Icon name="ph:eye-bold" size="16" class="text-zinc-400 shrink-0" />
      </div>
      <span
        class="inline-flex items-center gap-1.5 text-sm font-semibold text-zinc-700 dark:text-zinc-300"
      >
        <span
          class="w-2 h-2 rounded-full"
          :class="getUserStatusColor(user.status as DashboardRowUserStatus)"
        />
        {{ formatStatus(user.status as DashboardRowUserStatus) }}
      </span>
    </div>
    <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
      <div class="flex gap-2">
        <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 mb-1">ID</p>
        <Icon name="ph:identification-card-bold" size="16" class="text-zinc-400 shrink-0" />
      </div>
      <p class="text-sm font-semibold text-zinc-700 dark:text-zinc-300">#{{ user.id }}</p>
    </div>
    <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
      <div class="flex gap-2">
        <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 mb-1">
          Зарегистрирован
        </p>
        <Icon name="ph:calendar-bold" size="16" class="text-zinc-400 shrink-0" />
      </div>
      <p class="text-sm font-semibold text-zinc-700 dark:text-zinc-300">
        {{ fmtDate(user.create_at) }}
      </p>
    </div>
    <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
      <div class="flex gap-2">
        <p class="text-[10px] font-semibold uppercase tracking-widest text-zinc-400 mb-1">
          Последний вход
        </p>
        <Icon name="ph:clock-clockwise-bold" size="16" class="text-zinc-400 shrink-0" />
      </div>
      <p class="text-sm font-medium text-zinc-700 dark:text-zinc-300">
        {{ fmtDateTime(user.last_login) }}
      </p>
    </div>
  </div>
</template>
