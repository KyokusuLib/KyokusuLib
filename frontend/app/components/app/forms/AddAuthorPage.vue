<script setup lang="ts">
import { Input as BaseInput, SearchSelect, RichText as BaseRichTextEditor } from '@kyokusu-ui/vue'
import UploadIcon from '@/assets/images/special/add.png'

import { useAuthorForm } from '@/composables/api/forms/useAuthorForm'
import { COUNTRIES_LIST, METIER_LIST } from '@/constants/data'

const {
  form,
  errors,
  isLoading,
  fileInput,
  previewUrl,
  handleImageClick,
  handleImageChange,
  removeImage,
  submit,
} = useAuthorForm()
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <main
      class="grow bg-white dark:bg-radial-[at_center] dark:from-zinc-800 dark:to-zinc-950 dark:to-90% text-zinc-900 dark:text-zinc-200 transition-colors duration-300"
    >
      <div
        class="max-w-2xl mx-auto m-6 p-6 bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-sm transition-colors duration-300"
      >
        <div class="flex flex-col items-center mb-8">
          <h2 class="text-4xl font-bold text-zinc-900 dark:text-white mb-2">Добавить автора</h2>
          <p class="text-sm text-zinc-500 dark:text-zinc-400">
            Заполните информацию об авторе манги или ранобэ.
          </p>
        </div>

        <form @submit.prevent="submit" class="space-y-6">
          <div
            class="flex flex-col items-center sm:flex-row gap-6 pb-6 border-b border-zinc-200 dark:border-zinc-800"
          >
            <div
              class="relative group cursor-pointer shrink-0"
              @click="!isLoading && handleImageClick()"
            >
              <div
                class="w-32 h-32 rounded-xl overflow-hidden border-2 border-zinc-200 dark:border-zinc-700 group-hover:border-zinc-400 dark:group-hover:border-zinc-500 transition-colors bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center relative"
              >
                <img
                  v-if="previewUrl"
                  :src="previewUrl"
                  alt="Preview"
                  class="w-full h-full object-cover"
                />
                <div v-else class="text-zinc-400 dark:text-zinc-600 flex flex-col items-center">
                  <img :src="UploadIcon" class="w-8 h-8 opacity-50 dark:invert mb-2" />
                  <span class="text-xs font-medium">Фото</span>
                </div>

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
                <div
                  v-else
                  class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center"
                >
                  <span class="text-xs text-white font-medium">Изменить</span>
                </div>
              </div>
            </div>

            <div class="flex flex-col gap-3 w-full">
              <div class="flex flex-col text-left">
                <span class="text-sm font-medium text-zinc-900 dark:text-zinc-200 mb-1"
                  >Изображение</span
                >
                <p class="text-xs text-zinc-500 dark:text-zinc-400 mb-3">
                  Поддерживаются JPG, PNG. Макс. 5 МБ.
                </p>
              </div>
              <div class="flex gap-3">
                <button
                  type="button"
                  @click="handleImageClick"
                  :disabled="isLoading"
                  class="px-4 py-2 text-xs font-medium bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-900 dark:text-white rounded-lg cursor-pointer border border-zinc-200 dark:border-zinc-700 transition-colors disabled:opacity-50"
                >
                  Загрузить
                </button>
                <button
                  v-if="previewUrl"
                  type="button"
                  @click="removeImage"
                  :disabled="isLoading"
                  class="px-4 py-2 text-xs font-medium text-red-500 hover:text-red-600 dark:text-red-400 dark:hover:text-red-300 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors disabled:opacity-50"
                >
                  Удалить
                </button>
              </div>
              <input
                ref="fileInput"
                type="file"
                accept="image/*"
                class="hidden"
                @change="handleImageChange"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 gap-6 text-left">
            <BaseInput
              id="author-name"
              label="Имя автора"
              v-model="form.name"
              placeholder="Например: Макото Синкай"
              :disabled="isLoading"
              :error="errors.name"
            />

            <SearchSelect
              id="author-country"
              label="Страна"
              v-model="form.country"
              :options="COUNTRIES_LIST"
              placeholder="Выберите страну"
              :disabled="isLoading"
              :error="errors.country"
            />

            <SearchSelect
              id="author-metier"
              label="Род деятельности"
              v-model="form.metier"
              :options="METIER_LIST"
              placeholder="Выберите род деятельности"
              :disabled="isLoading"
              :error="errors.metier"
            />

            <BaseRichTextEditor
              id="author-about"
              label="Биография / Описание"
              v-model="form.about"
              :maxLength="2000"
              placeholder="Краткая биография автора..."
              :disabled="isLoading"
              :error="errors.about"
            />
          </div>

          <div class="pt-4 flex justify-end">
            <button
              type="submit"
              :disabled="isLoading"
              class="px-6 py-2.5 text-sm cursor-pointer font-medium bg-zinc-900 hover:bg-zinc-800 dark:bg-zinc-100 dark:hover:bg-white/50 text-white dark:text-zinc-900 rounded-lg transition-colors focus:ring-2 focus:ring-offset-2 focus:ring-offset-white dark:focus:ring-offset-zinc-900 focus:ring-zinc-900 dark:focus:ring-white disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
            >
              <svg
                v-if="isLoading"
                class="animate-spin -ml-1 mr-2 h-4 w-4"
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
              {{ isLoading ? 'Сохранение...' : 'Создать автора' }}
            </button>
          </div>
        </form>
      </div>
    </main>
  </div>
</template>
