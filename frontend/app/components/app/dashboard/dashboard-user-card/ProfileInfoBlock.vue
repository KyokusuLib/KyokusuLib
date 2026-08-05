<script setup lang="ts">
import { computed } from 'vue'
import type { DashboardUser } from '~/types/frontend/dashboard/users'

import { fmtDate } from '@/utils/date'
import { formatGender } from '@/utils/dashboard'

const props = defineProps<{ user: DashboardUser }>()

const GENDER_ICONS: Record<string, string> = {
  male: 'ph:gender-male-bold',
  female: 'ph:gender-female-bold',
  hidden: 'ph:gender-intersex-bold',
}

const genderIcon = computed(() => GENDER_ICONS[props.user.gender] ?? 'ph:question-bold')
const profileIcon = computed(() => (props.user.is_public ? 'ph:user-square-bold' : 'ph:lock-bold'))
</script>

<template>
  <div class="grid grid-cols-2 gap-3">
    <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
      <div class="flex gap-2">
        <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 mb-1">Пол</p>
        <Icon :name="genderIcon" size="16" class="text-zinc-400 shrink-0" />
      </div>
      <p class="text-sm font-semibold text-zinc-700 dark:text-zinc-300">
        {{ formatGender(user.gender) }}
      </p>
    </div>
    <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
      <div class="flex gap-2">
        <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 mb-1">
          Дата рождения
        </p>
        <Icon name="ph:calendar-dots-bold" size="16" class="text-zinc-400 shrink-0" />
      </div>
      <p class="text-sm font-semibold text-zinc-700 dark:text-zinc-300">
        {{ fmtDate(user.birthday) }}
      </p>
    </div>
    <div class="bg-zinc-50 dark:bg-zinc-800/50 rounded-xl px-4 py-3 border border-zinc-700">
      <div class="flex gap-2">
        <p class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 mb-1">Профиль</p>
        <Icon :name="profileIcon" size="16" class="text-zinc-400 shrink-0" />
      </div>
      <div class="flex items-center gap-1.5 text-sm font-medium text-zinc-700 dark:text-zinc-300">
        <Icon v-if="user.is_public" name="ph:globe-bold" size="14" class="text-green-500" />
        <span class="font-semibold text-zinc-700 dark:text-zinc-300">{{
          user.is_public ? 'Публичный' : 'Приватный'
        }}</span>
      </div>
    </div>
  </div>
</template>
