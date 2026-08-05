<script setup lang="ts">
import {
  ModalWindow,
  Input as BaseInput,
  SearchSelect,
  RichText as BaseRichTextEditor,
} from '@kyokusu-ui/vue'
import UploadIcon from '@/assets/images/special/add.png'
import { useAuthorSettings } from '@/composables/api/authors/useAuthorSettings'
import { COUNTRIES_LIST, METIER_LIST } from '@/constants/data'
import type { NovelaAuthorDetails } from '@/types/backend/novela'

defineOptions({
  inheritAttrs: false,
})

const props = defineProps<{
  modelValue: boolean
  author: NovelaAuthorDetails
}>()

const emit = defineEmits(['update:modelValue', 'updated'])

const {
  form,
  errors,
  isUpdating,
  fileInput,
  previewUrl,
  handleImageClick,
  handleImageChange,
  removeImage,
  updateAuthor,
} = useAuthorSettings(props.author)

const handleUpdate = async () => {
  const updatedAuthor = await updateAuthor(props.author.id)
  if (updatedAuthor) {
    emit('updated', updatedAuthor)
    emit('update:modelValue', false)
  }
}
</script>

<template>
  <ModalWindow
    v-if="modelValue"
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="Настройки автора"
    width="w-[95%] md:w-full md:max-w-3xl"
    center-title
  >
    <div class="max-h-[85vh] overflow-y-auto pr-2 custom-scrollbar">
      <div class="space-y-8 pb-6 px-2">
        <div
          class="flex flex-col items-center sm:flex-row gap-6 pb-8 border-b border-zinc-100 dark:border-zinc-800"
        >
          <div class="relative group cursor-pointer shrink-0" @click="handleImageClick">
            <div
              class="w-32 h-32 rounded-2xl overflow-hidden border-2 border-dashed border-zinc-200 dark:border-zinc-700 group-hover:border-yellow-500 transition-colors bg-zinc-50 dark:bg-zinc-900/50 flex items-center justify-center"
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
            <h4 class="text-sm font-bold">Фотография автора</h4>
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
          <div class="md:col-span-2">
            <BaseInput
              id="edit-author-name"
              label="Имя автора"
              v-model="form.name"
              placeholder="Например: Макото Синкай"
              :disabled="isUpdating"
              :error="errors.name"
            />
          </div>

          <SearchSelect
            id="edit-author-country"
            label="Страна"
            v-model="form.country"
            :options="COUNTRIES_LIST"
            placeholder="Выберите страну"
            :disabled="isUpdating"
            :error="errors.country"
          />

          <SearchSelect
            id="edit-author-metier"
            label="Род деятельности"
            v-model="form.metier"
            :options="METIER_LIST"
            placeholder="Выберите род деятельности"
            :disabled="isUpdating"
            :error="errors.metier"
          />

          <div class="md:col-span-2">
            <BaseRichTextEditor
              id="edit-author-about"
              label="Биография / Описание"
              v-model="form.about"
              :maxLength="2000"
              placeholder="Краткая биография автора..."
              :disabled="isUpdating"
              :error="errors.about"
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
