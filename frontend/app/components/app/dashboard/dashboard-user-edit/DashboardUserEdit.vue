<script setup lang="ts">
import { computed, ref, watch, nextTick } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import { ModalWindow, Label } from '@kyokusu-ui/vue'

import DashboardUserImageEdit from './DashboardUserImageEdit.vue'
import DashboardUserBasicEdit from './DashboardUserBasicEdit.vue'
import DashboardUserRoleEdit from './DashboardUserRoleEdit.vue'
import DashboardUserVisibilityEdit from './DashboardUserVisibilityEdit.vue'
import DashboardUserLevelEdit from './DashboardUserLevelEdit.vue'
import DashboardUserTagsEdit from './DashboardUserTagsEdit.vue'

import { useUserEdit } from '@/composables/api/dashboard/useUserEdit'
import { useUserTag } from '@/composables/api/profile/useUserTag'
import { useRolePermissions } from '@/composables/api/role/useRolePermissions'
import { KyokusuAppRole } from '~/types/enums/role-enum'

import type { DashboardUser } from '@/types/frontend/dashboard/users'

const props = defineProps<{
  modelValue: boolean
  user: DashboardUser | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  saved: []
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const { form, loadUser } = useUserEdit()
const { hasPermission } = useRolePermissions()
const { updateUserTags } = useUserTag()

const userLevel = ref(1)
const userExperience = ref(0)
const userTags = ref<DashboardUser['tags']>([])

let suppressTagSave = false

const debouncedSaveTags = useDebounceFn(async (userId: number) => {
  const tagIds = userTags.value.map((t) => t.tag_id)
  const ok = await updateUserTags(userId, tagIds)
  if (ok) {
    emit('saved')
  }
}, 800)

watch(userTags, () => {
  if (suppressTagSave) return
  if (!props.user) return
  debouncedSaveTags(props.user.id)
})

watch(
  () => props.user,
  (u) => {
    if (!u) return

    loadUser(u)
    userLevel.value = u.user_level?.level ?? 1
    userExperience.value = u.user_level?.experience ?? 0

    suppressTagSave = true
    userTags.value = [...(u.tags ?? [])]
    void nextTick(() => {
      suppressTagSave = false
    })
  },
  { immediate: true },
)
</script>

<template>
  <ModalWindow
    v-model="visible"
    :title="'Редактирование: ' + user?.name"
    center-title
    width="w-full max-w-250"
  >
    <div v-if="user" class="flex flex-col gap-6">
      <div class="flex flex-col gap-2">
        <Label label="Настройка изображения" />
        <div
          class="flex justify-center p-4 rounded-lg border bg-white dark:bg-zinc-900/50 dark:border-zinc-700/50"
        >
          <DashboardUserImageEdit v-model:avatar="form.picture" v-model:banner="form.banner" />
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <Label label="Основные данные" />
        <div
          class="flex p-4 rounded-lg border bg-white dark:bg-zinc-900/50 dark:border-zinc-700/50"
        >
          <DashboardUserBasicEdit :form="form" />
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <Label label="Роль и статус" />
        <div
          class="flex p-4 rounded-lg border bg-white dark:bg-zinc-900/50 dark:border-zinc-700/50"
        >
          <DashboardUserRoleEdit :form="form" />
        </div>
      </div>

      <div v-if="hasPermission(KyokusuAppRole.ADMIN)" class="flex flex-col gap-2">
        <Label label="Уровень и опыт" />
        <div
          class="flex p-4 rounded-lg border bg-white dark:bg-zinc-900/50 dark:border-zinc-700/50"
        >
          <DashboardUserLevelEdit
            :key="user?.id"
            v-model:level="userLevel"
            v-model:experience="userExperience"
          />
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <Label label="Теги" />
        <div
          class="flex p-4 rounded-lg border bg-white dark:bg-zinc-900/50 dark:border-zinc-700/50"
        >
          <DashboardUserTagsEdit :key="user?.id" v-model:tags="userTags" />
        </div>
      </div>

      <div v-if="hasPermission(KyokusuAppRole.ADMIN)" class="flex flex-col gap-2">
        <Label label="Настройки видимости" />
        <div
          class="flex p-4 rounded-lg border bg-white dark:bg-zinc-900/50 dark:border-zinc-700/50"
        >
          <DashboardUserVisibilityEdit :form="form" />
        </div>
      </div>
    </div>
  </ModalWindow>
</template>
