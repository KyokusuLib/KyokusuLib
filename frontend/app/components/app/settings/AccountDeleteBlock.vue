<script setup lang="ts">
import { Button } from '@kyokusu-ui/vue'
import InnerBlock from '~/components/ui/InnerBlock.vue'
import ModalConfirm from '@/components/common/ModalConfirm.vue'

const emit = defineEmits<{
  confirmDelete: []
}>()

defineOptions({ inheritAttrs: false })

const showConfirmModal = ref(false)
</script>

<template>
  <InnerBlock>
    <div class="flex flex-col gap-2">
      <div class="flex items-center gap-4">
        <span class="text-sm font-medium">Удаление аккаунта</span>
        <Icon name="ph:user-minus-bold" size="18" class="text-zinc-700 dark:text-zinc-200" />
      </div>
      <p class="text-xs text-zinc-500 dark:text-stone-500">
        Удаление аккаунта приведет к потери всех ваших данных. Аккаунт нельзя будет восстановить.
      </p>
    </div>

    <div class="flex justify-center items-center px-2">
      <Button variant="danger" @click="showConfirmModal = true">
        <span class="font-bold">Удалить аккаунт?</span>
      </Button>
    </div>
  </InnerBlock>

  <ModalConfirm
    v-model="showConfirmModal"
    title="Удаление аккаунта"
    description="Вы уверены, что хотите удалить аккаунт? Это действие необратимо — все ваши данные будут безвозвратно удалены."
    confirm-text="Да, удалить"
    cancel-text="Отмена"
    @confirm="emit('confirmDelete')"
  />
</template>
