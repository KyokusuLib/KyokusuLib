<script setup lang="ts">
import type { DashboardUserStats } from '@/types/frontend/dashboard/users'

import { fmtNumber } from '@/utils/date'

defineProps<{ stats: DashboardUserStats }>()
</script>

<template>
  <div class="grid grid-cols-2 gap-3">
    <div
      class="flex justify-center bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700"
    >
      <div class="flex items-center gap-2">
        <div>
          <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400">Комментарии</p>
          <div
            class="flex justify-center items-center gap-2 text-sm font-semibold text-zinc-700 dark:text-zinc-300"
          >
            <span class="text-[40px] ml-4">{{ fmtNumber(stats.total_comments) }}</span>
            <Icon name="ph:chat-circle-bold" size="24" class="text-blue-500 shrink-0" />
          </div>
        </div>
      </div>
    </div>
    <div
      class="flex justify-center bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700"
    >
      <div class="flex gap-2">
        <div>
          <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400">
            Прочитано глав
          </p>
          <div
            class="flex justify-center items-center gap-2 text-sm font-semibold text-zinc-700 dark:text-zinc-300"
          >
            <span class="text-[40px] ml-4">{{ fmtNumber(stats.read_chapters) }}</span>
            <Icon name="ph:book-open-bold" size="24" class="text-emerald-500 shrink-0" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
