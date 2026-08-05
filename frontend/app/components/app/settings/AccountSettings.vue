<script lang="ts" setup>
import { staticImage } from '@/utils/str'
import {
  DatePicker as BaseDateInput,
  Input as BaseInput,
  Select as BaseSelect,
  Button,
} from '@kyokusu-ui/vue'
import { useProfileSettings } from '@/composables/api/settings/useProfileSettings'
import { GenderSetting } from '@/types/enums/gender-enum'

import TabToggleSettings from './TabToggleSettings.vue'
import AccountDeleteBlock from './AccountDeleteBlock.vue'

const {
  profile,
  fileInput,
  bannerInput,
  errors,
  isLoading,
  isSuccess,
  isBirthDateDisable,
  handleAvatarClick,
  handleFileChange,
  handleBannerClick,
  handleBannerChange,
  handleDeleteAccount,
  saveChanges,
  removeAvatar,
  removeBanner,
} = useProfileSettings()

const genderOptions = [
  { value: GenderSetting.MALE, label: 'Мужской' },
  { value: GenderSetting.FEMALE, label: 'Женский' },
  { value: GenderSetting.HIDDEN, label: 'Не указывать' },
]
</script>

<template>
  <div class="space-y-8 cursor-default">
    <div class="flex flex-col items-center">
      <h2 class="text-3xl font-bold text-zinc-900 dark:text-white mb-1">Профиль</h2>
      <p class="text-sm text-zinc-500 dark:text-stone-400">
        Управляйте публичной информацией о себе.
      </p>
    </div>

    <transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="transform -translate-y-2 opacity-0"
      enter-to-class="transform translate-y-0 opacity-100"
    >
      <div
        v-if="isSuccess"
        class="p-3 bg-green-500/10 border border-green-500/20 rounded-md text-green-600 dark:text-green-400 text-sm"
      >
        Изменения успешно сохранены!
      </div>
      <div
        v-else-if="errors.global"
        class="p-3 bg-red-500/10 border border-red-500/20 rounded-md text-red-600 dark:text-red-400 text-sm"
      >
        {{ errors.global }}
      </div>
    </transition>

    <div
      class="flex flex-col sm:flex-row items-start sm:items-center gap-6 pb-6 border-b border-zinc-200 dark:border-zinc-800"
    >
      <div
        class="relative group cursor-pointer shrink-0"
        @click="!isLoading && handleAvatarClick()"
      >
        <div
          class="w-24 h-24 rounded-full overflow-hidden border-2 border-zinc-200 dark:border-zinc-700 group-hover:border-zinc-400 dark:group-hover:border-zinc-500 transition-colors relative"
        >
          <div
            v-if="isLoading"
            class="absolute inset-0 bg-black/60 flex items-center justify-center z-10"
          >
            <svg
              class="animate-spin h-6 w-6 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
          </div>
          <img
            :src="staticImage(profile.picture)"
            alt="Avatar"
            class="w-full h-full object-cover bg-zinc-100 dark:bg-zinc-800"
          />
        </div>
        <div
          v-if="!isLoading"
          class="absolute inset-0 bg-black/50 rounded-full flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <span class="text-xs text-white font-medium">Изменить</span>
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-zinc-900 dark:text-zinc-200 text-left"
          >Фотография профиля</span
        >
        <div class="flex gap-3">
          <button
            @click="handleAvatarClick"
            type="button"
            :disabled="isLoading"
            class="px-3 py-1.5 text-xs font-medium bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-900 dark:text-white rounded border border-zinc-200 dark:border-zinc-700 transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Загрузить фото
          </button>
          <button
            @click="removeAvatar"
            type="button"
            :disabled="isLoading"
            class="px-3 py-1.5 text-xs font-medium text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Удалить
          </button>
        </div>
        <p class="text-xs text-zinc-500 dark:text-stone-500 mt-1">
          Рекомендуется: квадратное изображение, мин. 400x400px.
        </p>
      </div>
      <input
        ref="fileInput"
        type="file"
        accept="image/*"
        class="hidden"
        @change="handleFileChange"
      />
    </div>

    <div
      class="flex flex-col sm:flex-row items-start sm:items-center gap-6 pb-6 border-b border-zinc-200 dark:border-zinc-800"
    >
      <div
        class="relative group cursor-pointer shrink-0 w-full sm:w-auto"
        @click="!isLoading && handleBannerClick()"
      >
        <div
          class="w-full sm:w-64 h-32 rounded-xl overflow-hidden border-2 border-zinc-200 dark:border-zinc-700 group-hover:border-zinc-400 dark:group-hover:border-zinc-500 transition-colors relative"
        >
          <div
            v-if="isLoading"
            class="absolute inset-0 bg-black/60 flex items-center justify-center z-10"
          >
            <svg
              class="animate-spin h-6 w-6 text-white"
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
            >
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"
              ></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
          </div>
          <img
            v-if="profile.banner"
            :src="staticImage(profile.banner)"
            alt="Banner"
            class="w-full h-full object-cover bg-zinc-100 dark:bg-zinc-800"
          />
          <div
            v-else
            class="w-full h-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center text-zinc-400 dark:text-zinc-600"
          >
            <span class="text-xs">Нет баннера</span>
          </div>
        </div>
        <div
          v-if="!isLoading"
          class="absolute inset-0 bg-black/50 rounded-xl flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
        >
          <span class="text-xs text-white font-medium">Изменить</span>
        </div>
      </div>

      <div class="flex flex-col gap-2">
        <span class="text-sm font-medium text-zinc-900 dark:text-zinc-200 text-left"
          >Баннер профиля</span
        >
        <div class="flex gap-3">
          <button
            @click="handleBannerClick"
            type="button"
            :disabled="isLoading"
            class="px-3 py-1.5 text-xs font-medium bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-900 dark:text-white rounded border border-zinc-200 dark:border-zinc-700 transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Загрузить фото
          </button>
          <button
            @click="removeBanner"
            type="button"
            :disabled="isLoading"
            class="px-3 py-1.5 text-xs font-medium text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Удалить
          </button>
        </div>
        <p class="text-left text-xs text-zinc-500 dark:text-stone-500 mt-1">
          Рекомендуется: 1200x320px.
        </p>
      </div>
      <input
        ref="bannerInput"
        type="file"
        accept="image/*"
        class="hidden"
        @change="handleBannerChange"
      />
    </div>

    <form @submit.prevent="saveChanges" class="space-y-6 text-left">
      <div class="grid grid-cols-1 gap-6">
        <BaseInput
          id="nickname"
          label="Никнейм"
          v-model="profile.name"
          :disabled="isLoading"
          :error="errors.name"
          description="Ограничение: мин. 6 символов, макс. 50 символов, латиница, мин. 1 большая буква, мин. 1 маленькая."
        />

        <div class="space-y-2 text-left">
          <label for="about" class="text-sm font-medium text-zinc-700 dark:text-stone-300 ml-0.75"
            >О себе</label
          >
          <textarea
            id="about"
            v-model="profile.about"
            rows="4"
            maxlength="500"
            :disabled="isLoading"
            class="w-full bg-zinc-50 dark:bg-zinc-950/50 border rounded-md py-2 px-3 text-sm text-zinc-900 dark:text-zinc-200 placeholder-zinc-400 dark:placeholder-zinc-600 focus:outline-none focus:ring-2 transition-all resize-none disabled:opacity-50"
            :class="[
              errors.about
                ? 'border-red-500 focus:ring-red-500'
                : 'border-zinc-300 dark:border-zinc-700 focus:ring-zinc-500 dark:focus:ring-zinc-600 focus:border-transparent',
            ]"
            placeholder="Расскажите немного о себе..."
          ></textarea>
          <p v-if="errors.about" class="text-xs text-red-500 dark:text-red-400 text-right">
            {{ errors.about }}
          </p>
          <p v-else class="text-xs text-zinc-500 dark:text-stone-500 text-right">
            {{ profile.about ? profile.about.length : 0 }}/500
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
          <div class="flex flex-col">
            <BaseSelect
              id="gender"
              label="Пол"
              v-model="profile.gender"
              :options="genderOptions"
              :disabled="isLoading"
            />
            <p class="text-xs text-zinc-500 dark:text-stone-500 mt-0.5 ml-0.75">
              Поможет подбирать соответствующий контент.
            </p>
          </div>

          <div class="flex flex-col">
            <ClientOnly>
              <BaseDateInput
                id="birthdate"
                label="Дата рождения"
                v-model="profile.birthday"
                :disabled="isBirthDateDisable || isLoading"
              />
            </ClientOnly>
            <p class="text-xs text-zinc-500 dark:text-stone-500 mt-0.5 ml-0.75">
              Выставляется единожды. В этот день сайт выдаст вам награду.
            </p>
          </div>
        </div>

        <TabToggleSettings
          v-model="profile.is_public"
          id="show_profile_toggle"
          title="Публичный профиль"
          imgSrc="ph:link-simple"
          description="Переключите для отключения публичного профиля на сайте."
        />

        <AccountDeleteBlock @confirm-delete="handleDeleteAccount" />
      </div>

      <div class="pt-6 border-t border-zinc-200 dark:border-zinc-800 flex justify-end gap-3">
        <button
          type="submit"
          :disabled="isLoading"
          class="px-4 py-2 text-sm font-medium bg-zinc-900 hover:bg-zinc-800 dark:bg-white dark:hover:bg-zinc-300 text-white dark:text-zinc-900 rounded transition-colors focus:ring-2 focus:ring-offset-2 focus:ring-offset-white dark:focus:ring-offset-zinc-900 focus:ring-zinc-900 dark:focus:ring-white cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          <svg
            v-if="isLoading"
            class="animate-spin -ml-1 mr-2 h-4 w-4 text-white dark:text-zinc-900"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          {{ isLoading ? 'Сохранение...' : 'Сохранить изменения' }}
        </button>
      </div>
    </form>
  </div>
</template>
