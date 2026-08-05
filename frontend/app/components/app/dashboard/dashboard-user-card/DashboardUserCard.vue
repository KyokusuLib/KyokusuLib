<script setup lang="ts">
import { computed } from 'vue'
import { ModalWindow } from '@kyokusu-ui/vue'

import type { DashboardUser } from '@/types/frontend/dashboard/users'

import ImageBlock from './ImageBlock.vue'
import MetaDataBlock from './MetaDataBlock.vue'
import AboutBlock from './AboutBlock.vue'
import ProfileInfoBlock from './ProfileInfoBlock.vue'
import LevelBlock from './LevelBlock.vue'
import StatsBlock from './StatsBlock.vue'
import TagsBlock from './TagsBlock.vue'
import SettingsBlock from './SettingsBlock.vue'

const props = defineProps<{
  modelValue: boolean
  user: DashboardUser | null
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})
</script>

<template>
  <ModalWindow
    v-model="visible"
    :title="user?.name ?? 'Пользователь'"
    center-title
    width="w-full max-w-lg"
  >
    <div v-if="user" class="space-y-6 py-2">
      <ImageBlock :user="user" />
      <MetaDataBlock :user="user" />

      <AboutBlock :about="user.about" />

      <ProfileInfoBlock :user="user" />

      <LevelBlock v-if="user.user_level" :level="user.user_level" />

      <StatsBlock v-if="user.user_stats" :stats="user.user_stats" />

      <TagsBlock
        v-if="user.tags && user.tags.length > 0"
        :tags="user.tags"
        :activeTag="user.active_tag"
      />

      <SettingsBlock :settings="user.settings" />
    </div>
  </ModalWindow>
</template>
