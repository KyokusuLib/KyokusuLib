<script setup lang="ts">
import { ref, computed } from 'vue'
import { useUsers } from '@/composables/api/dashboard/useUsers'
import { useUserApi } from '@/composables/api/user/userApi'
import {
  mapToRows,
  getUserRoleColor,
  getUserStatusColor,
  formatRole,
  formatStatus,
  getLevelColor,
} from '@/utils/dashboard'

import UiTable from '@/components/ui/Table.vue'
import UiTableSearch from '@/components/ui/TableSearch.vue'

import DashboardUserCard from '@/components/app/dashboard/dashboard-user-card/DashboardUserCard.vue'
import DashboardUserEdit from '@/components/app/dashboard/dashboard-user-edit/DashboardUserEdit.vue'

import CircleBlock from '@/components/ui/CircleBlock.vue'
import ModalConfirm from '@/components/common/ModalConfirm.vue'

import type { UiTableColumn } from '@/components/ui/Table.vue'
import type { DashboardUser } from '@/types/frontend/dashboard/users'

const { users, isLoading, refreshUsers } = useUsers()
const { deleteUser } = useUserApi()

const searchQuery = ref('')
const selectedUserId = ref<number | null>(null)

const showDeleteConfirm = ref(false)
const showEdit = ref(false)

const targetDeleteUser = ref<DashboardUser | null>(null)
const targetEditUser = ref<DashboardUser | null>(null)

const tableColumns: UiTableColumn[] = [
  { key: 'id', label: 'ID', align: 'center', sortable: true },
  { key: 'name', label: 'Имя пользователя', sortable: true },
  { key: 'role', label: 'Роль', align: 'center' },
  { key: 'status', label: 'Статус', align: 'center' },
  { key: 'level', label: 'Уровень', align: 'center', sortable: true },
  { key: 'registered', label: 'Зарегистрирован', align: 'center' },
]

const tableRows = computed<Record<string, any>[]>(() => {
  return (users.value ?? []).map(mapToRows)
})

const selectedUser = computed<DashboardUser | null>(() => {
  if (selectedUserId.value === null) return null
  return users.value?.find((u) => u.id === selectedUserId.value) ?? null
})

async function onEdit(userId: number) {
  const stale = users.value
  await refreshUsers()
  targetEditUser.value =
    users.value.find((u) => u.id === userId) ?? stale.find((u) => u.id === userId) ?? null
  showEdit.value = true
}

function onDelete(userId: number) {
  targetDeleteUser.value = users.value.find((u) => u.id === userId)!
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  if (targetDeleteUser.value === null) return
  await deleteUser(targetDeleteUser.value.id)
  targetDeleteUser.value = null
}
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between gap-4 flex-wrap">
      <div>
        <h2 class="text-xl font-bold text-zinc-900 dark:text-white">Управление пользователями</h2>
        <p class="text-sm text-zinc-500 dark:text-zinc-400 mt-1">
          Просмотр, поиск и управление учётными записями
        </p>
      </div>

      <UiTableSearch
        v-model="searchQuery"
        placeholder="Поиск по таблице..."
        class="w-full sm:w-64"
      />
    </div>

    <UiTable
      :columns="tableColumns"
      :items="tableRows"
      :loading="isLoading"
      v-model:search="searchQuery"
      empty-text="Пользователи не найдены"
    >
      <template #cell-name="{ item }">
        <button
          class="font-medium text-zinc-900 dark:text-white hover:text-amber-500 dark:hover:text-amber-500 transition-colors cursor-pointer text-left"
          @click="selectedUserId = item.id"
        >
          {{ item.name }}
        </button>
      </template>

      <template #cell-role="{ item }">
        <span
          class="inline-flex px-2.5 py-1 rounded-lg text-[11px] font-bold uppercase tracking-wider"
          :class="getUserRoleColor(item.role)"
        >
          {{ formatRole(item.role) }}
        </span>
      </template>

      <template #cell-status="{ item }">
        <span class="inline-flex items-center gap-1.5">
          <span class="w-2 h-2 rounded-full" :class="getUserStatusColor(item.status)" />
          <span class="text-zinc-500 dark:text-zinc-400 text-sm">
            {{ formatStatus(item.status) }}
          </span>
        </span>
      </template>

      <template #cell-level="{ item }">
        <CircleBlock
          class="ml-4"
          :text="item.levelDisplay"
          :colorClass="getLevelColor(item.level)"
          border
        />
      </template>

      <template #cell-registered="{ item }">
        <span class="font-semibold p-2 border border-amber-500 rounded-xl">
          {{ item.registered }}
        </span>
      </template>

      <template #actions="{ item }">
        <div class="flex items-center justify-center gap-1">
          <button
            class="w-8 h-8 rounded-xl bg-zinc-100 dark:bg-zinc-800 hover:bg-blue-100 dark:hover:bg-blue-500/20 text-zinc-500 hover:text-blue-600 dark:hover:text-blue-400 transition-colors flex items-center justify-center cursor-pointer"
            title="Редактировать"
            @click.stop="onEdit(item.id)"
          >
            <Icon name="ph:pencil-bold" size="14" />
          </button>
          <button
            class="w-8 h-8 rounded-xl bg-zinc-100 dark:bg-zinc-800 hover:bg-red-100 dark:hover:bg-red-500/20 text-zinc-500 hover:text-red-600 dark:hover:text-red-400 transition-colors flex items-center justify-center cursor-pointer"
            title="Удалить"
            @click.stop="onDelete(item.id)"
          >
            <Icon name="ph:trash-bold" size="14" />
          </button>
        </div>
      </template>
    </UiTable>
  </div>

  <DashboardUserCard
    :model-value="selectedUserId != null"
    :user="selectedUser"
    @update:model-value="selectedUserId = null"
  />

  <DashboardUserEdit
    :model-value="targetEditUser != null"
    :user="targetEditUser"
    @update:model-value="targetEditUser = null"
    @saved="refreshUsers()"
  />

  <ModalConfirm
    v-model="showDeleteConfirm"
    title="Удаление пользователя"
    :description="`Вы уверены, что хотите удалить пользователя ${targetDeleteUser?.name}? Это действие нельзя отменить.`"
    confirm-text="Удалить"
    @confirm="confirmDelete"
  />
</template>
