<script setup lang="ts">
import { ModalWindow, Input as BaseInput, RichText as BaseRichTextEditor } from '@kyokusu-ui/vue'
import UploadIcon from '@/assets/images/special/add.png'
import InnerBlock from '@/components/ui/InnerBlock.vue'
import { Select, Separator } from '@kyokusu-ui/vue'

import { useTeamSettings } from '@/composables/api/teams/useTeamSettings'
import type { Team } from '@/types/frontend/teams'

defineOptions({
  inheritAttrs: false,
})

const props = defineProps<{
  modelValue: boolean
  team: Team
}>()

const emit = defineEmits(['update:modelValue', 'updated'])

const {
  form,
  errors,
  isUpdating,
  fileInput,
  previewUrl,
  bannerInput,
  previewBannerUrl,
  handleImageClick,
  handleImageChange,
  removeImage,
  handleBannerClick,
  handleBannerChange,
  removeBanner,
  submitUpdate,
} = useTeamSettings(props.team)

const handleUpdate = async () => {
  const updatedTeam = await submitUpdate(props.team.slug)
  if (updatedTeam) {
    emit('updated', updatedTeam)
    emit('update:modelValue', false)
  }
}

const teamTypeOptions = [
  { value: 'open', label: 'Открытая' },
  { value: 'private', label: 'Закрытая' },
]
</script>

<template>
  <ModalWindow
    v-if="modelValue"
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="Настройки команды"
    width="w-[95%] md:w-full md:max-w-3xl"
    center-title
  >
    <div class="max-h-[85vh] overflow-y-auto pr-1 sm:pr-2 custom-scrollbar">
      <div class="space-y-6 sm:space-y-8 pb-6 px-1 sm:px-2 mt-2">
        <h3 class="ml-1 text-base font-bold text-zinc-900 dark:text-white mb-3">Изображения</h3>
        <div
          class="flex flex-col md:flex-row justify-between gap-4 pb-6 border-b border-zinc-100 dark:border-zinc-800"
        >
          <InnerBlock class="flex w-full">
            <div class="flex flex-col items-center sm:flex-row gap-4 w-full">
              <div class="relative group cursor-pointer shrink-0" @click="handleImageClick">
                <div
                  class="w-24 h-24 sm:w-28 sm:h-28 rounded-2xl overflow-hidden border-2 border-dashed border-zinc-200 dark:border-zinc-700 group-hover:border-yellow-500 transition-colors bg-zinc-50 dark:bg-zinc-900/50 flex items-center justify-center"
                >
                  <img v-if="previewUrl" :src="previewUrl" class="w-full h-full object-cover" />
                  <div v-else class="flex flex-col items-center text-zinc-400">
                    <img :src="UploadIcon" class="w-6 h-6 dark:invert opacity-20 mb-2" />
                    <span class="text-[10px] font-bold uppercase">Аватар</span>
                  </div>
                </div>
                <input
                  ref="fileInput"
                  type="file"
                  accept="image/*"
                  class="hidden"
                  @change="handleImageChange"
                />
              </div>
              <div class="flex flex-col items-center sm:items-start gap-1 text-center sm:text-left">
                <h4 class="text-sm font-bold text-zinc-900 dark:text-zinc-100">Логотип</h4>
                <p class="text-xs text-zinc-500">JPG/PNG, до 5МБ</p>
                <button
                  v-if="previewUrl"
                  @click.stop="removeImage"
                  class="text-xs font-bold text-red-500 hover:text-red-600 transition-colors uppercase mt-1 cursor-pointer"
                >
                  Удалить
                </button>
              </div>
            </div>
          </InnerBlock>

          <InnerBlock class="flex w-full">
            <div class="flex flex-col items-center sm:flex-row gap-4 w-full">
              <div
                class="relative group cursor-pointer shrink-0 w-full sm:w-48 h-24 sm:h-28"
                @click="handleBannerClick"
              >
                <div
                  class="w-full h-full rounded-2xl overflow-hidden border-2 border-dashed border-zinc-200 dark:border-zinc-700 group-hover:border-yellow-500 transition-colors bg-zinc-50 dark:bg-zinc-900/50 flex items-center justify-center"
                >
                  <img
                    v-if="previewBannerUrl"
                    :src="previewBannerUrl"
                    class="w-full h-full object-cover"
                  />
                  <div v-else class="flex flex-col items-center text-zinc-400">
                    <img :src="UploadIcon" class="w-6 h-6 dark:invert opacity-20 mb-2" />
                    <span class="text-[10px] font-bold uppercase">Баннер</span>
                  </div>
                </div>
                <input
                  ref="bannerInput"
                  type="file"
                  accept="image/*"
                  class="hidden"
                  @change="handleBannerChange"
                />
              </div>
              <div class="flex flex-col items-center sm:items-start gap-1 text-center sm:text-left">
                <h4 class="text-sm font-bold text-zinc-900 dark:text-zinc-100">Баннер команды</h4>
                <p class="text-xs text-zinc-500">Рекомендуется 3:1</p>
                <button
                  v-if="previewBannerUrl"
                  @click.stop="removeBanner"
                  class="text-xs font-bold text-red-500 hover:text-red-600 transition-colors uppercase mt-1 cursor-pointer"
                >
                  Удалить
                </button>
              </div>
            </div>
          </InnerBlock>
        </div>

        <div class="grid grid-cols-1 gap-6 text-left">
          <div>
            <h3 class="ml-1 text-base font-bold text-zinc-900 dark:text-white mb-3">
              Описание команды
            </h3>
            <BaseRichTextEditor
              id="edit-team-description"
              v-model="form.description"
              :maxLength="2000"
              placeholder="Расскажите о вашей команде..."
              :disabled="isUpdating"
              :error="errors.description"
            />
          </div>
        </div>

        <Separator />
        <div class="pt-2">
          <h3 class="ml-1 text-base font-bold text-zinc-900 dark:text-white mt-4 mb-3">
            Названия ролей
          </h3>
          <InnerBlock class="p-4 sm:p-5">
            <div
              class="grid grid-cols-1 sm:grid-cols-3 md:grid-cols-3 gap-4 sm:gap-5 cursor-default"
            >
              <div class="w-full">
                <div class="flex items-center gap-1.5 mb-2 ml-1">
                  <Icon name="ph:crown" size="16" class="text-yellow-500" />
                  <label
                    for="edit-team-owner-role"
                    class="block text-sm font-bold text-zinc-700 dark:text-zinc-300"
                    >Владелец</label
                  >
                </div>
                <BaseInput
                  id="edit-team-owner-role"
                  v-model="form.owner_role_name"
                  placeholder="Например: Создатель"
                  :disabled="isUpdating"
                  :error="errors.owner_role_name"
                />
              </div>
              <div class="w-full">
                <div class="flex items-center gap-1.5 mb-2 ml-1">
                  <Icon name="ph:sword" size="16" class="text-blue-500" />
                  <label
                    for="edit-team-moderator-role"
                    class="block text-sm font-bold text-zinc-700 dark:text-zinc-300"
                    >Модератор</label
                  >
                </div>
                <BaseInput
                  id="edit-team-moderator-role"
                  v-model="form.moderator_role_name"
                  placeholder="Например: Редактор"
                  :disabled="isUpdating"
                  :error="errors.moderator_role_name"
                />
              </div>
              <div class="w-full">
                <div class="flex items-center gap-1.5 mb-2 ml-1">
                  <Icon name="ph:user-bold" size="16" class="text-emerald-500" />
                  <label
                    for="edit-team-member-role"
                    class="block text-sm font-bold text-zinc-700 dark:text-zinc-300"
                    >Участник</label
                  >
                </div>
                <BaseInput
                  id="edit-team-member-role"
                  v-model="form.member_role_name"
                  placeholder="Например: Переводчик"
                  :disabled="isUpdating"
                  :error="errors.member_role_name"
                />
              </div>
            </div>
          </InnerBlock>

          <Separator />

          <h3 class="ml-1 text-base font-bold text-zinc-900 dark:text-white mt-4 mb-3">
            Тип команды
          </h3>
          <InnerBlock class="mt-4">
            <Select
              id="team-type"
              v-model="form.team_type"
              :options="teamTypeOptions"
              :disabled="isUpdating"
            />
          </InnerBlock>
        </div>
        <Separator />
        <div class="flex flex-col sm:flex-row-reverse gap-3 pt-6">
          <button
            @click="handleUpdate"
            :disabled="isUpdating"
            class="flex-1 sm:flex-none sm:w-48 py-2.5 cursor-pointer bg-zinc-900 dark:bg-white hover:bg-zinc-700 dark:hover:bg-zinc-200 text-white dark:text-zinc-900 font-bold rounded-xl active:scale-[0.98] transition-all text-sm disabled:opacity-50 flex items-center justify-center gap-2 shadow-sm"
          >
            <Icon v-if="isUpdating" name="ph:spinner-gap-bold" class="animate-spin" />
            {{ isUpdating ? 'Сохранение...' : 'Сохранить' }}
          </button>
          <button
            @click="$emit('update:modelValue', false)"
            :disabled="isUpdating"
            class="flex-1 sm:flex-none sm:w-32 py-2.5 cursor-pointer bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400 font-bold rounded-xl hover:bg-zinc-200 dark:hover:bg-zinc-700 transition-colors text-sm"
          >
            Отмена
          </button>
        </div>
      </div>
    </div>
  </ModalWindow>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #3f3f46;
  border-radius: 10px;
}
</style>
