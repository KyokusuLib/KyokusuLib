<script setup lang="ts">
import TabToggleSettings from '../../settings/TabToggleSettings.vue'
import type { UserEditForm } from '@/composables/api/dashboard/useUserEdit'

defineProps<{
  form: UserEditForm
}>()
</script>

<template>
  <div class="flex flex-col gap-3 w-full">
    <TabToggleSettings
      v-model="form.isShowTag"
      id="show_tag_toggle"
      title="Показывать теги"
      imgSrc="ph:tag-bold"
      description="Отображать теги пользователя в публичном профиле."
    />
    <TabToggleSettings
      v-model="form.isShowBookmark"
      id="show_bookmark_toggle"
      title="Показывать закладки"
      imgSrc="ph:bookmark-bold"
      description="Отображать закладки пользователя в публичном профиле."
    />
  </div>
</template>
