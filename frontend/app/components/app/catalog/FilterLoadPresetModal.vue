<script setup lang="ts">
import { computed } from 'vue'
import { ModalWindow, Button } from '@kyokusu-ui/vue'
import ModalConfirm from '~/components/common/ModalConfirm.vue'
import type { CatalogFilterPreset } from '~/types/backend/catalog-filters'

const props = defineProps<{
  modelValue: boolean
  presets: CatalogFilterPreset[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  close: []
  load: [preset: CatalogFilterPreset]
  delete: [id: number]
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const deletingPresetId = ref<number | null>(null)
const deletingPresetName = ref('')

const deleteConfirmVisible = computed(() => deletingPresetId.value !== null)

function onClose() {
  emit('update:modelValue', false)
  emit('close')
}

function onLoad(preset: CatalogFilterPreset) {
  emit('load', preset)
  onClose()
}

function confirmDelete(preset: CatalogFilterPreset) {
  deletingPresetId.value = preset.id
  deletingPresetName.value = preset.name
}

function onDeleteConfirmed() {
  if (deletingPresetId.value !== null) {
    emit('delete', deletingPresetId.value)
  }
  deletingPresetId.value = null
  deletingPresetName.value = ''
}

function onDeleteCancelled() {
  deletingPresetId.value = null
  deletingPresetName.value = ''
}

function formatDate(dateStr: string): string {
  const d = new Date(dateStr)
  return d.toLocaleDateString('ru-RU', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<template>
  <ModalWindow
    v-model="visible"
    title="Сохранённые шаблоны"
    center-title
    width="w-full max-w-lg"
    @close="onClose"
  >
    <div class="flex flex-col gap-2 py-2">
      <div v-if="presets.length === 0" class="text-center py-8 text-zinc-400 text-sm">
        <Icon name="ph:file-search-duotone" size="40" class="mx-auto mb-3 opacity-50" />
        <p>У вас пока нет сохранённых шаблонов</p>
      </div>

      <div v-else class="flex flex-col gap-2 max-h-80 overflow-y-auto pr-1">
        <div
          v-for="preset in presets"
          :key="preset.id"
          class="flex items-center gap-3 px-4 py-3 rounded-xl border border-zinc-200 dark:border-zinc-700 hover:border-zinc-400 dark:hover:border-zinc-500 bg-white dark:bg-zinc-900/50 transition-colors cursor-pointer group"
          @click="onLoad(preset)"
        >
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-zinc-800 dark:text-zinc-200 truncate">
              {{ preset.name }}
            </p>
            <p class="text-xs text-zinc-400 mt-0.5">
              {{ formatDate(preset.created_at) }}
            </p>
          </div>

          <Button
            variant="outline"
            class="opacity-0 group-hover:opacity-100 transition-opacity shrink-0"
            @click.stop="confirmDelete(preset)"
          >
            <Icon name="ph:trash-bold" size="14" class="text-red-500" />
          </Button>
        </div>
      </div>
    </div>

    <ModalConfirm
      v-model="deleteConfirmVisible"
      title="Удаление шаблона"
      :description="`Вы уверены, что хотите удалить шаблон «${deletingPresetName}»?`"
      confirm-text="Удалить"
      cancel-text="Отмена"
      @confirm="onDeleteConfirmed"
      @cancel="onDeleteCancelled"
    />
  </ModalWindow>
</template>
