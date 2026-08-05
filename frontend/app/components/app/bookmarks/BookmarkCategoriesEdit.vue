<script setup lang="ts">
import { ref, useTemplateRef } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { useBookmark } from '~/composables/api/novela/useBookmark'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'
import { Tooltip, ModalWindow, Input, Toggle, Button } from '@kyokusu-ui/vue'

import ModalConfirm from '~/components/common/ModalConfirm.vue'

const { user } = useAuthStore()
const { notify } = useNotificationStore()
const {
  bookmarkCategories,
  createBookmarkCategory,
  updateBookmarkCategory,
  deleteBookmarkCategory,
} = useBookmark()

const isOpen = ref(false)
const containerRef = useTemplateRef('containerRef')

onClickOutside(containerRef, () => {
  isOpen.value = false
})

const newCategoryName = ref('')

const isDeleteModalOpen = ref(false)
const categoryToDelete = ref<number | null>(null)

const isEditModalOpen = ref(false)
const editingCategory = ref<any>(null)
const editCategoryName = ref('')
const isCategoryVisible = ref(true)

const handleCreate = async () => {
  if (!newCategoryName.value.trim()) return
  try {
    await createBookmarkCategory(user!.id, newCategoryName.value.trim())
    newCategoryName.value = ''
    notify({
      title: 'Категория добавлена',
      type: 'success',
    })
  } catch (e) {
    console.error(e)
    notify({
      title: 'Произошла ошибка',
      content: `${e}`,
      type: 'error',
    })
  }
}

const startEdit = (cat: any) => {
  if (!cat.user_id) return
  editingCategory.value = cat
  editCategoryName.value = cat.label
  isCategoryVisible.value = cat.visible
  isEditModalOpen.value = true
}

const handleUpdate = async () => {
  if (!editingCategory.value || !editCategoryName.value.trim()) return
  try {
    await updateBookmarkCategory(
      editingCategory.value.id,
      editCategoryName.value.trim(),
      isCategoryVisible.value,
    )
    isEditModalOpen.value = false
    editingCategory.value = null
    notify({
      title: 'Категория обновлена',
      type: 'success',
    })
  } catch (e) {
    console.error(e)
    notify({
      title: 'Произошла ошибка',
      content: `${e}`,
      type: 'error',
    })
  }
}

const promptDelete = (id: number) => {
  categoryToDelete.value = id
  isDeleteModalOpen.value = true
}

const handleDelete = async () => {
  if (!categoryToDelete.value) return
  try {
    await deleteBookmarkCategory(categoryToDelete.value)
  } catch (e) {
    console.error(e)
    notify({
      title: 'Произошла ошибка',
      content: `${e}`,
      type: 'error',
    })
  } finally {
    categoryToDelete.value = null
    notify({
      title: 'Категория удалена',
      type: 'success',
    })
  }
}
</script>

<template>
  <div ref="containerRef" class="relative z-30">
    <button
      @click="isOpen = !isOpen"
      class="w-10 h-10 rounded-2xl -mt-1 flex items-center justify-center bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 text-zinc-500 hover:text-zinc-900 dark:hover:text-white hover:border-zinc-300 dark:hover:border-zinc-700 transition-all shadow-sm cursor-pointer"
    >
      <Icon name="ph:faders-bold" size="20" />
    </button>

    <Transition name="shrink">
      <div
        v-if="isOpen"
        class="absolute right-0 top-full mt-2 w-80 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-3xl shadow-2xl p-4 origin-top-right flex flex-col gap-4"
      >
        <div
          class="flex items-center justify-center pb-2 border-b border-zinc-100 dark:border-zinc-800"
        >
          <span class="text-xs font-black uppercase tracking-widest text-zinc-900 dark:text-white"
            >Редактор категорий</span
          >
        </div>

        <div class="flex flex-col gap-2 no-scrollbar max-h-60 overflow-y-auto relative">
          <div
            v-for="cat in bookmarkCategories"
            :key="cat.id"
            class="flex items-center justify-between p-2 rounded-xl bg-zinc-50 dark:bg-zinc-800/50 group"
          >
            <div class="flex items-center justify-between w-full">
              <div
                class="flex items-center gap-2 text-sm text-zinc-700 dark:text-zinc-300 font-medium"
              >
                <Icon :name="cat.icon" size="16" class="text-zinc-400" />
                {{ cat.label }}
              </div>

              <div
                v-if="cat.user_id"
                class="flex items-center gap-1 group-hover:opacity-100 transition-opacity"
              >
                <Tooltip text="Настройки" position="left">
                  <button
                    @click="startEdit(cat)"
                    class="p-1 text-zinc-400 hover:text-yellow-500 transition-colors cursor-pointer flex"
                  >
                    <Icon name="ph:gear-bold" size="16" />
                  </button>
                </Tooltip>
                <Tooltip text="Удалить" position="left">
                  <button
                    @click="promptDelete(cat.id)"
                    class="p-1 text-zinc-400 hover:text-red-500 transition-colors cursor-pointer flex"
                  >
                    <Icon name="ph:trash-bold" size="16" />
                  </button>
                </Tooltip>
              </div>
            </div>
          </div>
        </div>

        <div class="pt-2 border-t border-zinc-100 dark:border-zinc-800 flex gap-2 relative">
          <input
            v-model="newCategoryName"
            @keyup.enter="handleCreate"
            placeholder="Новая категория..."
            class="flex-1 bg-zinc-50 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700 rounded-xl px-3 py-2 text-xs text-zinc-900 dark:text-white outline-none focus:border-yellow-500 transition-colors"
          />
          <Tooltip text="Добавить" position="top">
            <button
              @click="handleCreate"
              :disabled="!newCategoryName.trim()"
              class="px-3 py-2 rounded-xl bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 disabled:opacity-50 transition-opacity flex items-center justify-center cursor-pointer h-full"
            >
              <Icon name="ph:plus-bold" size="16" />
            </button>
          </Tooltip>
        </div>
      </div>
    </Transition>

    <ModalWindow
      v-model="isEditModalOpen"
      title="Настройки категории"
      center-title
      width="w-full max-w-md"
    >
      <div class="flex flex-col gap-6 py-2">
        <Input
          id="category-name"
          label="Название категории"
          v-model="editCategoryName"
          placeholder="Введите название"
        />

        <div class="flex justify-between gap-2 ml-2">
          <label for="category-visible" class="text-[14px] font-medium"
            >Показывать категорию в профиле</label
          >
          <Toggle id="category-visible" v-model="isCategoryVisible" />
        </div>

        <div class="flex justify-between gap-3 mt-2">
          <Button variant="outline" @click="isEditModalOpen = false">Отмена</Button>
          <Button variant="primary" @click="handleUpdate">Сохранить</Button>
        </div>
      </div>
    </ModalWindow>

    <ModalConfirm
      v-model="isDeleteModalOpen"
      title="Удаление категории"
      description="Вы уверены, что хотите удалить эту категорию? Все новеллы в ней будут перемещены в 'В планах'."
      confirm-text="Удалить"
      cancel-text="Отмена"
      @confirm="handleDelete"
    />
  </div>
</template>

<style scoped>
.shrink-enter-active,
.shrink-leave-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}
.shrink-enter-from,
.shrink-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
