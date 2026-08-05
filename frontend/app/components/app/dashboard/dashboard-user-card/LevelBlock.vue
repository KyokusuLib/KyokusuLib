<script setup lang="ts">
import type { DashboardUserLevel } from '@/types/frontend/dashboard/users'

import UserExperience from '@/components/features/UserExperience/UserExperience.vue'
import InnerBlock from '@/components/ui/InnerBlock.vue'

defineProps<{ level: DashboardUserLevel }>()
</script>

<template>
  <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
    <div class="flex gap-2 px-4 mb-4">
      <span class="text-[10px] font-bold uppercase tracking-widest text-zinc-400"
        >Информация о уровне</span
      >
      <Icon name="ph:flask-bold" size="16" class="text-zinc-400 shrink-0" />
    </div>
    <InnerBlock>
      <div class="flex flex-col gap-2">
        <div class="flex justify-start items-center px-4 gap-2 cursor-default">
          <span class="text-[10px] font-bold uppercase tracking-widest text-zinc-400"
            >Текущий уровень:</span
          >
          <span class="text-[10px] text-zinc-300 mb-0.5">{{ level.level }}</span>
          <span class="text-[10px] font-bold mb-0.5">{{ level.level_title }}</span>
        </div>
        <div class="flex justify-start items-center px-4 gap-2 cursor-default">
          <span class="text-[10px] font-bold uppercase tracking-widest text-zinc-400"
            >До следующего уровня:
          </span>
          <span class="text-[10px] font-bold mb-0.5">{{
            level.experience + level.xp_needed_for_next
          }}</span>
          <span class="text-[10px] text-zinc-500 shadow-black shadow-2xs">XP</span>
        </div>
        <div class="flex justify-start items-center px-4 gap-2 cursor-default">
          <span class="text-[10px] font-bold uppercase tracking-widest text-zinc-400"
            >Текущий опыт:
          </span>
          <span class="text-[10px] font-bold mb-0.5">{{ level.experience }}</span>
          <span class="text-[10px] text-zinc-500 shadow-black shadow-2xs">XP</span>
        </div>
      </div>
    </InnerBlock>
    <UserExperience
      :level="level.level"
      :currentExp="level.experience"
      :expToNextLevel="level.xp_needed_for_next"
      :levelTitle="level.level_title"
    />
  </div>
</template>
