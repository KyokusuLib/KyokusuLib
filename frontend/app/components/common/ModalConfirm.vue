<script setup lang="ts">
import { ModalWindow, Button } from '@kyokusu-ui/vue'

defineOptions({
  inheritAttrs: false,
})

interface Props {
  modelValue: boolean
  title?: string
  description?: string
  confirmText?: string
  cancelText?: string
}

withDefaults(defineProps<Props>(), {
  title: 'Подтверждение',
  description: 'Вы уверены, что хотите выполнить это действие?',
  confirmText: 'Да',
  cancelText: 'Отмена',
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel'])

const onConfirm = () => {
  emit('confirm')
  emit('update:modelValue', false)
}

const onCancel = () => {
  emit('cancel')
  emit('update:modelValue', false)
}
</script>

<template>
  <div class="modal-confirm">
    <ModalWindow
      :model-value="modelValue"
      @update:model-value="$emit('update:modelValue', $event)"
      :title="title"
      :center-title="true"
      width="w-full max-w-sm"
    >
      <div class="space-y-6 py-2">
        <p class="text-zinc-600 dark:text-zinc-400 text-center text-sm">
          {{ description }}
        </p>

        <div class="flex items-center justify-center gap-3">
          <Button variant="outline" @click="onCancel">
            {{ cancelText }}
          </Button>
          <Button variant="danger" @click="onConfirm">
            {{ confirmText }}
          </Button>
        </div>
      </div>
    </ModalWindow>
  </div>
</template>
