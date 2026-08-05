<script setup lang="ts">
import { ref } from 'vue'
import { Input, RichText, DatePicker, Label } from '@kyokusu-ui/vue'
import TabToggleSettings from '../../settings/TabToggleSettings.vue'
import type { UserEditForm } from '@/composables/api/dashboard/useUserEdit'
import { GENDER_OPTIONS } from '@/composables/api/dashboard/useUserEdit'

defineProps<{
  form: UserEditForm
}>()

const isGenderOpen = ref(false)
</script>

<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-5 w-full">
    <!-- name + gender in one row -->
    <div class="flex flex-col gap-2">
      <Label label="Имя пользователя" />
      <Input v-model="form.name" placeholder="Введите имя" />
    </div>

    <div class="flex flex-col gap-2">
      <Label label="Пол" />
      <div class="relative">
        <select
          v-model="form.gender"
          class="w-full h-10 rounded-lg border pl-3 pr-10 text-sm outline-none appearance-none transition focus:ring-2"
          :style="{
            backgroundColor: 'var(--k-editor-bg)',
            borderColor: 'var(--k-editor-border)',
            color: 'var(--k-editor-text)',
          }"
          @focus="isGenderOpen = true"
          @blur="isGenderOpen = false"
        >
          <option v-for="opt in GENDER_OPTIONS" :key="opt.value" :value="opt.value">
            {{ opt.label }}
          </option>
        </select>
        <Icon
          name="ph:caret-down-bold"
          size="14"
          class="absolute right-3 top-1/2 -translate-y-1/2 text-zinc-400 pointer-events-none transition-transform duration-200"
          :class="{ 'rotate-180': isGenderOpen }"
        />
      </div>
    </div>

    <!-- birthday -->
    <div class="sm:col-span-2">
      <ClientOnly>
        <div class="flex flex-col gap-2">
          <Label label="Дата рождения" />
          <DatePicker id="birthday" v-model="form.birthday" />
        </div>
        <template #fallback>
          <div class="flex flex-col gap-2">
            <Label label="Дата рождения" />
            <Input v-model="form.birthday" />
          </div>
        </template>
      </ClientOnly>
    </div>

    <!-- is_public — full width -->
    <div class="sm:col-span-2">
      <TabToggleSettings
        v-model="form.isPublic"
        id="show_profile_toggle"
        title="Публичный профиль"
        imgSrc="ph:link-simple"
        description="Переключите для отключения публичного профиля на сайте."
      />
    </div>

    <!-- about — full width -->
    <div class="flex flex-col gap-1.5 sm:col-span-2">
      <RichText
        id="about"
        label="О себе"
        v-model="form.about"
        placeholder="Введите описание"
        :max-length="500"
      />
    </div>
  </div>
</template>
