<script setup lang="ts">
import { Tooltip, Button } from '@kyokusu-ui/vue'
import { useNovelaFilters } from '@/composables/api/novela/useNovelaFilters'
import { useAuthStore } from '@/stores/auth'

import CatalogFilters from '@/components/app/catalog/CatalogFilters.vue'
import FilterSavePresetModal from './FilterSavePresetModal.vue'
import FilterLoadPresetModal from './FilterLoadPresetModal.vue'
import type { CatalogFilterPreset } from '@/types/backend/catalog-filters'

const emit = defineEmits<{
  apply: [params: Record<string, any>]
  reset: []
  'search-update': [params: Record<string, any>]
}>()

const saveFilterPresetOpen = ref(false)
const loadFilterPresetOpen = ref(false)

const { isAuthenticated } = useAuthStore()

const {
  savedFilters,
  hasSavedFilters,
  hasActiveFilters,
  saveFilters,
  getFilters,
  deleteFilterPreset,
  loadPreset,
} = useNovelaFilters()

const openFilterSaveModal = () => {
  saveFilterPresetOpen.value = true
}

const openFilterLoadModal = () => {
  loadFilterPresetOpen.value = true
}

const saveFilterPreset = async (presetName: string) => {
  await saveFilters(presetName)
}

const onLoadPreset = (preset: CatalogFilterPreset) => {
  loadPreset(preset)
  emit('apply', {})
}

const onDeletePreset = async (id: number) => {
  await deleteFilterPreset(id)
}

onMounted(async () => {
  await getFilters()
})
</script>

<template>
  <div
    class="min-h-150 bg-white dark:bg-zinc-900/50 border border-zinc-200 dark:border-zinc-800 rounded-xl p-6 shadow-sm hidden lg:block"
  >
    <div
      class="flex items-center justify-center text-sm font-bold text-zinc-900 dark:text-white mb-6 gap-2"
    >
      <div class="flex items-center gap-2 mb-2">
        <Icon name="ph:funnel-bold" size="16" class="text-yellow-500" />
        <span>Фильтры</span>
      </div>

      <Tooltip v-if="hasActiveFilters && isAuthenticated" text="Сохранить шаблон">
        <Button size="md" variant="outline" @click="openFilterSaveModal">
          <Icon name="ph:archive-bold" size="16" class="text-yellow-500" />
        </Button>
      </Tooltip>

      <Tooltip v-if="hasSavedFilters" text="Сохраненные шаблоны">
        <Button size="md" variant="outline" @click="openFilterLoadModal">
          <Icon name="ph:hard-drives-bold" size="16" class="text-yellow-500" />
        </Button>
      </Tooltip>
    </div>
    <CatalogFilters
      @apply="(params) => emit('apply', params)"
      @search-update="(params) => emit('search-update', params)"
      @reset="() => emit('reset')"
    />

    <FilterSavePresetModal v-model="saveFilterPresetOpen" @save="saveFilterPreset" />

    <FilterLoadPresetModal
      v-model="loadFilterPresetOpen"
      :presets="savedFilters"
      @load="onLoadPreset"
      @delete="onDeletePreset"
    />
  </div>
</template>
