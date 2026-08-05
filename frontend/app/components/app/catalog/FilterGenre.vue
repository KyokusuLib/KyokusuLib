<script setup lang="ts">
import { computed } from 'vue'
import { MultiSelect } from '@kyokusu-ui/vue'
import type { MultiSelectOption } from '@kyokusu-ui/vue'
import type { NovelaFilterOption } from '~/types/frontend/query/novela-filters'

const props = defineProps<{
  options: NovelaFilterOption[]
  modelValue: string[]
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string[]]
  search: [query: string]
}>()

const multiOptions = computed<MultiSelectOption[]>(() =>
  props.options.map((o) => ({ id: o.value, label: o.label })),
)

function onUpdate(value: (string | number | MultiSelectOption)[]) {
  const ids = value.map((v) => (typeof v === 'object' ? String(v.id) : String(v)))
  emit('update:modelValue', ids)
}

function onSearch(query: string) {
  emit('search', query)
}
</script>

<template>
  <div>
    <MultiSelect
      id="filter-genres"
      label="Жанры"
      :model-value="modelValue"
      :options="multiOptions"
      placeholder="Выберите жанры..."
      @update:model-value="onUpdate"
      @search="onSearch"
    />
  </div>
</template>
