<script setup lang="ts">
import { computed } from 'vue'
import { ModalWindow, Input, Label, Button } from '@kyokusu-ui/vue'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  close: []
  save: [presetName: string]
}>()

const presetName = defineModel<string>('presetName', { default: '' })

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const onClose = () => {
  emit('update:modelValue', false)
  emit('close')
}

const onSave = () => {
  emit('save', presetName.value)
  emit('close')
}
</script>

<template>
  <ModalWindow
    v-model="visible"
    title="Сохранить шаблон поиска"
    center-title
    width="w-full max-w-md"
  >
    <div class="flex flex-col gap-2 py-2">
      <Label label="Название шаблона" />
      <Input id="category-name" v-model="presetName" placeholder="Введите название" />

      <div class="flex justify-between gap-3 mt-2">
        <Button variant="outline" @click="onClose">Отмена</Button>
        <Button variant="primary" @click="onSave">Сохранить</Button>
      </div>
    </div>
  </ModalWindow>
</template>
