<script setup lang="ts">
import { Label } from '@kyokusu-ui/vue'
import { getUserRoleColor, getUserStatusColor } from '@/utils/dashboard'
import type { UserEditForm } from '@/composables/api/dashboard/useUserEdit'
import { STATUS_OPTIONS, ROLE_OPTIONS } from '@/composables/api/dashboard/useUserEdit'

import { useRolePermissions } from '~/composables/api/role/useRolePermissions'
import { KyokusuAppRole } from '@/types/enums/role-enum'

defineProps<{
  form: UserEditForm
}>()

const { hasPermission } = useRolePermissions()

const showRole = computed(() => hasPermission(KyokusuAppRole.ADMIN))
</script>

<template>
  <div class="grid grid-cols-1 gap-4 w-full" :class="showRole && 'sm:grid-cols-2'">
    <template v-if="showRole">
      <div class="flex flex-col gap-2">
        <Label label="Роль в системе" />
        <div class="relative">
          <select
            v-model="form.role"
            class="w-full h-10 rounded-lg border pl-3 pr-10 text-sm outline-none appearance-none transition focus:ring-2"
            :style="{
              backgroundColor: 'var(--k-editor-bg)',
              borderColor: 'var(--k-editor-border)',
              color: 'var(--k-editor-text)',
            }"
          >
            <option v-for="opt in ROLE_OPTIONS" :key="opt.value" :value="opt.value">
              {{ opt.label }}
            </option>
          </select>
          <span
            v-if="form.role"
            class="absolute right-3 top-1/2 -translate-y-1/2 inline-flex px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-wider pointer-events-none"
            :class="getUserRoleColor(form.role)"
          >
            {{ form.role }}
          </span>
        </div>
      </div>
    </template>

    <!-- status -->
    <div class="flex flex-col gap-2">
      <Label label="Статус" />
      <div class="relative">
        <select
          v-model="form.status"
          class="w-full h-10 rounded-lg border pl-3 pr-10 text-sm outline-none appearance-none transition focus:ring-2"
          :style="{
            backgroundColor: 'var(--k-editor-bg)',
            borderColor: 'var(--k-editor-border)',
            color: 'var(--k-editor-text)',
          }"
        >
          <option v-for="opt in STATUS_OPTIONS" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <span
          class="absolute right-3 top-1/2 -translate-y-1/2 w-2 h-2 rounded-full pointer-events-none"
          :class="getUserStatusColor(form.status)"
        />
      </div>
    </div>
  </div>
</template>
