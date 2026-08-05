<script setup lang="ts">
import { Input as BaseInput } from '@kyokusu-ui/vue'

const modelValue = defineModel<string>({ required: true })

const emit = defineEmits<{
  update: [value: string]
}>()

function onUpdate(value: string) {
  emit('update', value)
}
</script>

<template>
  <div class="flex flex-col space-y-1.5">
    <label
      class="text-xs font-bold uppercase tracking-wider text-zinc-500 dark:text-zinc-400 ml-0.5"
    >
      Название
    </label>
    <BaseInput
      id="catalog-search"
      :model-value="modelValue"
      placeholder="Название новеллы..."
      @update:model-value="onUpdate"
    />
  </div>
</template>
