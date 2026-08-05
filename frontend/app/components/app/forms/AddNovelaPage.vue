<script setup lang="ts">
import HeaderApp from '@/components/common/HeaderApp.vue'
import FooterApp from '@/components/common/FooterApp.vue'
import {
  Input as BaseInput,
  SearchSelect,
  RichText as BaseRichTextEditor,
  MultiSelect as BaseMultiSelect,
} from '@kyokusu-ui/vue'
import YearSelect from '@/components/features/YearSelector/YearSelector.vue'
import UploadIcon from '@/assets/images/special/add.png'

import { useNovelaForm } from '@/composables/api/forms/useNovelaForm'

import {
  NOVELA_TYPES,
  AGE_RATINGS,
  NOVELA_STATUSES,
  NOVELA_GENRES,
  TRANSLATION_STATUSES,
  NOVELA_CATEGORIES,
  COUNTRIES_LIST,
} from '@/constants/data'

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
} = useNovelaForm()
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <main
      class="grow bg-white dark:bg-radial-[at_center] dark:from-zinc-800 dark:to-zinc-950 dark:to-90% text-zinc-900 dark:text-zinc-200 transition-colors duration-300"
    >
      <div
        class="max-w-4xl mx-auto m-6 p-6 bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-sm transition-colors duration-300"
      >
        <div class="flex flex-col items-center mb-8">
          <h2 class="text-4xl font-bold text-zinc-900 dark:text-white mb-2">Добавить Новеллу</h2>
          <p class="text-sm text-zinc-500 dark:text-zinc-400">
            Заполните информацию о новом произведении.
          </p>
        </div>

        <form @submit.prevent="submit" class="space-y-8">
          <!-- Cover Image Upload -->
          <div
            class="flex flex-col items-center sm:flex-row gap-6 pb-6 border-b border-zinc-200 dark:border-zinc-800"
          >
            <div
              class="relative group cursor-pointer shrink-0"
              @click="!isLoading && handleImageClick()"
            >
              <div
                class="w-40 h-56 rounded-xl overflow-hidden border-2 border-zinc-200 dark:border-zinc-700 group-hover:border-zinc-400 dark:group-hover:border-zinc-500 transition-colors bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center relative"
              >
                <img
                  v-if="previewUrl"
                  :src="previewUrl"
                  alt="Preview"
                  class="w-full h-full object-cover"
                />
                <div v-else class="text-zinc-400 dark:text-zinc-600 flex flex-col items-center">
                  <img :src="UploadIcon" class="w-8 h-8 opacity-50 dark:invert mb-2" />
                  <span class="text-xs font-medium">Обложка</span>
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
                  >Обложка тайтла</span
                >
                <p class="text-xs text-zinc-500 dark:text-zinc-400 mb-3">
                  Рекомендуемый размер 200x300. JPG, PNG. Макс. 5 МБ.
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

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-left">
            <!-- Main Info -->
            <div class="md:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-6">
              <BaseInput
                id="novela-title"
                label="Название"
                v-model="form.title"
                placeholder="Например: Поднятие уровня в одиночку"
                :disabled="isLoading"
                :error="errors.title"
              />

              <BaseInput
                id="novela-alt-title"
                label="Альтернативные названия"
                v-model="form.alternativeTitles"
                placeholder="Solo Leveling, Ore dake Level Up na Ken"
                :disabled="isLoading"
                :error="errors.alternativeTitles"
              />
            </div>

            <SearchSelect
              id="novela-type"
              label="Тип"
              v-model="form.type"
              :options="NOVELA_TYPES"
              placeholder="Выберите тип"
              :disabled="isLoading"
              :error="errors.type"
            />

            <SearchSelect
              id="novela-age"
              label="Возрастной рейтинг"
              v-model="form.ageRating"
              :options="AGE_RATINGS"
              placeholder="Выберите рейтинг"
              :disabled="isLoading"
              :error="errors.ageRating"
            />

            <YearSelect
              id="ranobe-year"
              label="Год выпуска"
              v-model="form.releaseYear"
              :from="1980"
              :to="2030"
              placeholder="Выберите год"
              :disabled="isLoading"
              :error="errors.releaseYear"
            />

            <SearchSelect
              id="novela-status"
              label="Статус тайтла"
              v-model="form.status"
              :options="NOVELA_STATUSES"
              placeholder="Онгоинг / Завершен"
              :disabled="isLoading"
              :error="errors.status"
            />

            <SearchSelect
              id="novela-tl-status"
              label="Статус перевода"
              v-model="form.translationStatus"
              :options="TRANSLATION_STATUSES"
              placeholder="Продолжается / Заморожен"
              :disabled="isLoading"
              :error="errors.translationStatus"
            />
            <SearchSelect
              id="novela-country"
              label="Страна"
              v-model="form.country"
              :options="COUNTRIES_LIST"
              placeholder="Выберите страну"
              :disabled="isLoading"
              :error="errors.country"
            />

            <div class="md:col-span-2 space-y-6">
              <BaseMultiSelect
                id="novela-genres"
                label="Жанры"
                v-model="form.genres"
                :options="NOVELA_GENRES"
                placeholder="Выберите жанры"
                :disabled="isLoading"
                :error="errors.genres"
              />

              <BaseMultiSelect
                id="novela-categories"
                label="Категории"
                v-model="form.categories"
                :options="NOVELA_CATEGORIES"
                placeholder="Выберите категории"
                :disabled="isLoading"
                :error="errors.categories"
              />

              <BaseRichTextEditor
                id="novela-desc"
                label="Описание"
                v-model="form.description"
                :maxLength="5000"
                placeholder="Описание сюжета..."
                :disabled="isLoading"
                :error="errors.description"
              />
            </div>
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
              {{ isLoading ? 'Сохранение...' : 'Создать тайтл' }}
            </button>
          </div>
        </form>
      </div>
    </main>
  </div>
</template>
