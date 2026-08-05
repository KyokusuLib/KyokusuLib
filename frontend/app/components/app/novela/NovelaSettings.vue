<script setup lang="ts">
import { computed } from 'vue'
import {
  ModalWindow,
  Input as BaseInput,
  SearchSelect,
  MultiSelect as BaseMultiSelect,
  RichText as BaseRichTextEditor,
} from '@kyokusu-ui/vue'

import { useNotificationStore } from '~/stores/notification'

defineOptions({
  inheritAttrs: false,
})

import YearSelect from '~/components/features/YearSelector/YearSelector.vue'
import UploadIcon from '@/assets/images/special/add.png'

import { useNovelaSettings } from '~/composables/api/novela/useNovelaSettings'
import { useAuthors } from '~/composables/api/authors/useAuthors'
import { useNovela } from '~/composables/api/novela/useNovela'
import {
  NOVELA_STATUSES,
  TRANSLATION_STATUSES,
  NOVELA_TYPES,
  AGE_RATINGS,
  COUNTRIES_LIST,
  NOVELA_GENRES,
  NOVELA_CATEGORIES,
} from '~/constants/data'
import type { NovelaDetails } from '~/types/backend/novela'

interface Props {
  modelValue: boolean
  novela: NovelaDetails
}

const props = defineProps<Props>()
const emit = defineEmits(['update:modelValue', 'updated'])

const { notify } = useNotificationStore()

const { updateNovela, isUpdating } = useNovela()
const { searchAuthors, foundAuthors, isSearching } = useAuthors()
const {
  form,
  fileInput,
  previewUrl,
  selectedFile,
  handleImageClick,
  handleImageChange,
  removeImage,
} = useNovelaSettings(props.novela)

const initialAuthors = props.novela.authors?.map((a) => ({ id: a.id, label: a.name })) || []

const authorsOptions = computed(() => {
  const map = new Map()
  for (const a of initialAuthors) map.set(a.id, a)
  for (const a of foundAuthors.value) map.set(a.id, a)
  return Array.from(map.values())
})

const arrStrWOutObjects = (arr: any[]) => {
  return arr.map((s) => {
    if (typeof s === 'object') return s.id
    return s
  })
}

const handleUpdate = async () => {
  const payload = {
    ...form,
    authors: form.authors.map((a: any) => a.id),
    age_rating: form.ageRating,
    release_date: form.releaseYear,
    translation_status: form.translationStatus,
    genres: arrStrWOutObjects(form.genres),
    categories: arrStrWOutObjects(form.categories),
    alternative_titles:
      typeof form.alternativeTitles === 'string'
        ? form.alternativeTitles
            .split('/')
            .map((t) => t.trim())
            .filter(Boolean)
        : form.alternativeTitles,
  }

  try {
    await updateNovela(props.novela.id, payload, selectedFile.value)
    const uiUpdatePayload = {
      ...payload,
      poster_url: previewUrl.value,
      alternative_titles: payload.alternative_titles,
      authors: form.authors.map((a: any) => ({ id: a.id, name: a.label, role: 'Author' })),
    }
    emit('updated', uiUpdatePayload)
  } catch (e) {
    notify({
      type: 'error',
      title: 'Ошибка',
      content: e instanceof Error ? e.message : String(e),
    })
  }
}
</script>

<template>
  <ModalWindow
    v-if="modelValue"
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="Настройки произведения"
    width="w-[95%] md:w-full md:max-w-5xl"
    center-title
  >
    <div class="max-h-[85vh] overflow-y-auto pr-2 custom-scrollbar">
      <div class="space-y-8 pb-6 px-2">
        <div
          class="flex flex-col items-center sm:flex-row gap-6 pb-8 border-b border-zinc-100 dark:border-zinc-800"
        >
          <div class="relative group cursor-pointer shrink-0" @click="handleImageClick">
            <div
              class="w-36 h-52 rounded-2xl overflow-hidden border-2 border-dashed border-zinc-200 dark:border-zinc-700 group-hover:border-yellow-500 transition-colors bg-zinc-50 dark:bg-zinc-900/50 flex items-center justify-center"
            >
              <img v-if="previewUrl" :src="previewUrl" class="w-full h-full object-cover" />
              <div v-else class="flex flex-col items-center text-zinc-400">
                <img :src="UploadIcon" class="w-8 h-8 dark:invert opacity-20 mb-2" />
                <span class="text-[10px] font-bold uppercase">Загрузить</span>
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
          <div class="flex flex-col items-start gap-2 w-full text-center sm:text-left">
            <h4 class="text-sm font-bold">Обложка произведения</h4>
            <p class="text-xs text-zinc-500">JPG/PNG, до 5МБ.</p>
            <button
              v-if="previewUrl"
              @click.stop="removeImage"
              class="text-xs font-bold text-red-500 uppercase mt-2 cursor-pointer"
            >
              Удалить
            </button>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-left">
          <div class="md:col-span-2 grid grid-cols-1 md:grid-cols-2 gap-6">
            <BaseInput
              id="edit-title"
              label="Название"
              v-model="form.title"
              :disabled="isUpdating"
            />
            <BaseInput
              id="edit-alt-title"
              label="Альтернативные названия"
              v-model="form.alternativeTitles"
              :disabled="isUpdating"
            />
          </div>

          <SearchSelect
            id="edit-type"
            label="Тип"
            v-model="form.type"
            :options="NOVELA_TYPES"
            :disabled="isUpdating"
          />

          <SearchSelect
            id="edit-age"
            label="Рейтинг"
            v-model="form.ageRating"
            :options="AGE_RATINGS"
            :disabled="isUpdating"
          />

          <YearSelect
            id="edit-year"
            label="Год выхода"
            v-model="form.releaseYear"
            :from="1980"
            :to="2030"
            :disabled="isUpdating"
          />

          <SearchSelect
            id="edit-country"
            label="Страна"
            v-model="form.country"
            :options="COUNTRIES_LIST"
            :disabled="isUpdating"
          />

          <SearchSelect
            id="edit-status"
            label="Статус выхода"
            v-model="form.status"
            :options="NOVELA_STATUSES"
            :disabled="isUpdating"
          />

          <SearchSelect
            id="edit-tl-status"
            label="Статус перевода"
            v-model="form.translationStatus"
            :options="TRANSLATION_STATUSES"
            :disabled="isUpdating"
          />

          <div class="md:col-span-2 space-y-6">
            <BaseMultiSelect
              id="edit-authors"
              label="Авторы"
              v-model="form.authors"
              :options="authorsOptions"
              :loading="isSearching"
              :disabled="isUpdating"
              @search="searchAuthors"
              placeholder="Поиск по имени..."
            />

            <BaseMultiSelect
              id="edit-genres"
              label="Жанры"
              v-model="form.genres"
              :options="NOVELA_GENRES"
              :disabled="isUpdating"
            />

            <BaseMultiSelect
              id="edit-categories"
              label="Категории"
              v-model="form.categories"
              :options="NOVELA_CATEGORIES"
              :disabled="isUpdating"
            />

            <BaseRichTextEditor
              id="edit-desc"
              label="Описание"
              v-model="form.description"
              :maxLength="5000"
              :disabled="isUpdating"
            />
          </div>
        </div>

        <div
          class="flex flex-col md:flex-row gap-3 pt-6 border-t border-zinc-100 dark:border-zinc-800"
        >
          <button
            @click="handleUpdate"
            :disabled="isUpdating"
            class="flex-1 py-2 cursor-pointer bg-zinc-900 dark:bg-white hover:bg-zinc-700 dark:hover:bg-zinc-200 text-white dark:text-zinc-900 font-black rounded-2xl active:scale-[0.98] transition-all text-sm disabled:opacity-50 flex items-center justify-center gap-2"
          >
            <div
              v-if="isUpdating"
              class="w-4 h-4 border-2 border-zinc-500 border-t-transparent rounded-full animate-spin"
            ></div>
            {{ isUpdating ? 'Сохранение...' : 'Сохранить изменения' }}
          </button>
          <button
            @click="$emit('update:modelValue', false)"
            :disabled="isUpdating"
            class="md:flex-none px-10 py-2 cursor-pointer bg-zinc-50 dark:bg-zinc-800 text-zinc-500 font-bold rounded-2xl hover:bg-zinc-100 dark:hover:bg-zinc-700 transition-colors text-sm"
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
