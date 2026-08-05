<script setup lang="ts">
import type { DashboardUser } from '@/types/frontend/dashboard/users'

import { staticImage } from '@/utils/str'
import { getUserRoleColor, formatRole } from '@/utils/dashboard'

defineProps<{ user: DashboardUser }>()
</script>

<template>
  <div class="flex justify-center items-center gap-4 mb-4">
    <NuxtLink :to="{ path: `/profile/${user.id}` }">
      <div class="flex justify-center items-center gap-4">
        <img
          :src="staticImage(user.picture)"
          :alt="user.name"
          class="w-16 h-16 rounded-2xl object-cover shadow-sm shrink-0"
        />
        <div class="min-w-0">
          <h3 class="text-lg font-bold text-zinc-900 dark:text-white truncate">
            {{ user.name }}
          </h3>
          <div class="flex items-center gap-2 mt-1">
            <span
              class="inline-flex px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider"
              :class="getUserRoleColor(user.role)"
            >
              {{ formatRole(user.role) }}
            </span>
            <span v-if="user.email" class="text-xs text-zinc-400 truncate hidden sm:inline">
              {{ user.email }}
            </span>
          </div>
        </div>
      </div>
    </NuxtLink>
  </div>
</template>
