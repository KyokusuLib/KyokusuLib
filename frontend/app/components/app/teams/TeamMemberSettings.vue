<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ModalWindow, Select } from '@kyokusu-ui/vue'
import { useTeam } from '@/composables/api/teams/useTeams'
import type { CustomRoleNames, TeamMember } from '@/types/frontend/teams'
import { staticImage } from '@/utils/str'

const props = defineProps<{
  modelValue: boolean
  member: TeamMember
  customRoleNames: CustomRoleNames[]
  slug: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  updated: [member: TeamMember]
}>()

const { updateTeamMember, isLoading } = useTeam()

const selectedRole = ref(props.member.role || '')

const selectedLabel = computed(() => {
  const option = props.customRoleNames.find((o) => o.value === selectedRole.value)
  return option?.label || ''
})

watch(
  () => props.modelValue,
  (val) => {
    if (val) {
      selectedRole.value = props.member.role || ''
    }
  },
)

const handleSave = async () => {
  const label = selectedLabel.value
  const success = await updateTeamMember(props.slug, props.member.user.id, {
    role: selectedRole.value || undefined,
    custom_role_name: label || undefined,
  })

  if (success) {
    emit('updated', {
      ...props.member,
      role: selectedRole.value,
      role_name: label || props.member.role_name,
    })
    emit('update:modelValue', false)
  }
}

const handleCancel = () => {
  selectedRole.value = props.member.role || ''
  emit('update:modelValue', false)
}
</script>

<template>
  <ModalWindow
    v-if="modelValue"
    :model-value="modelValue"
    @update:model-value="handleCancel"
    title="Настройка участника"
    width="w-[95%] md:w-full md:max-w-lg"
    center-title
  >
    <div class="flex flex-col gap-6 px-1 sm:px-2 pb-2 team-member-settings-modal">
      <div
        class="flex items-center gap-4 bg-zinc-50 dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-2xl p-4"
      >
        <img
          v-if="member.user.picture"
          :src="staticImage(member.user.picture)"
          :alt="member.user.name"
          class="w-12 h-12 rounded-full object-cover shrink-0"
        />
        <Icon v-else name="ph:user-bold" size="24" class="text-zinc-700/40 shrink-0" />

        <div class="flex-1 min-w-0">
          <p class="font-bold text-sm text-zinc-900 dark:text-zinc-100 truncate">
            {{ member.user.name }}
          </p>
          <p class="text-xs font-semibold text-zinc-500 dark:text-zinc-400 mt-0.5">
            {{ member.role_name }}
          </p>
        </div>
      </div>

      <div class="flex flex-col">
        <label
          for="member-custom-role"
          class="block text-sm font-bold text-zinc-700 dark:text-zinc-300 mb-2 ml-1"
        >
          Кастомная роль
        </label>
        <p class="text-xs text-zinc-500 mb-3 ml-1">
          Выберите роль для участника. Кастомное название роли будет установлено соответственно
          выбранной роли.
        </p>
        <Select
          id="member-custom-role"
          v-model="selectedRole"
          :options="customRoleNames"
          :disabled="isLoading"
        />
      </div>

      <div class="flex flex-col sm:flex-row-reverse gap-3 pt-2">
        <button
          @click="handleSave"
          :disabled="isLoading"
          class="flex-1 sm:flex-none sm:w-44 py-2.5 cursor-pointer bg-zinc-900 dark:bg-white hover:bg-zinc-700 dark:hover:bg-zinc-200 text-white dark:text-zinc-900 font-bold rounded-xl active:scale-[0.98] transition-all text-sm disabled:opacity-50 flex items-center justify-center gap-2 shadow-sm"
        >
          <Icon v-if="isLoading" name="ph:spinner-gap-bold" class="animate-spin" />
          {{ isLoading ? 'Сохранение...' : 'Сохранить' }}
        </button>
        <button
          @click="handleCancel"
          :disabled="isLoading"
          class="flex-1 sm:flex-none sm:w-32 py-2.5 cursor-pointer bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 font-bold rounded-xl hover:bg-zinc-200 dark:hover:bg-zinc-700 transition-colors text-sm"
        >
          Отмена
        </button>
      </div>
    </div>
  </ModalWindow>
</template>

<style>
.k-modal-content:has(.team-member-settings-modal) {
  overflow: visible !important;
}
</style>
