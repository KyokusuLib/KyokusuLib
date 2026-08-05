<script setup lang="ts">
import { ref } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { SearchSelect } from '@kyokusu-ui/vue'
import { REPORT_COMMNET_REASONS } from '~/constants/data'

const props = defineProps<{ commentId: number; reported: boolean }>()
const emit = defineEmits(['close', 'submit', 'cancel'])

const target = ref(null)
const selectedReason = ref<string | null>(null)
const error = ref(false)

onClickOutside(target, () => emit('close'), { ignore: ['.report-toggle-btn'] })

const handleSend = () => {
  if (!selectedReason.value) {
    error.value = true
    return
  }
  emit('submit', { commentId: props.commentId, reason: selectedReason.value })
  emit('close')
}

const handleSendCansel = () => {
  emit('cancel', props.commentId)
  emit('close')
}
</script>

<template>
  <div
    ref="target"
    class="absolute right-0 top-8 z-50 w-64 bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-2xl p-4 animate-in fade-in zoom-in duration-200"
  >
    <div class="space-y-4">
      <div class="flex flex-col gap-1">
        <span class="text-[10px] font-black uppercase tracking-widest text-zinc-400 ml-1"
          >Причина жалобы</span
        >
        <SearchSelect
          id="report-reason"
          v-model="selectedReason"
          :options="REPORT_COMMNET_REASONS"
          placeholder="Выберите причину"
          :error="error && !selectedReason ? 'Обязательно выберите причину' : ''"
        />
      </div>

      <div class="flex flex-col gap-2">
        <button
          @click="!reported ? handleSend() : handleSendCansel()"
          class="w-full py-1.5 text-white text-[10px] font-black uppercase tracking-widest rounded-xl transition-colors cursor-pointer"
          :class="[!reported ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700']"
        >
          {{ !reported ? 'Отправить жалобу' : 'Отменить жалобу' }}
        </button>

        <button
          @click="emit('close')"
          class="w-full py-1.5 bg-zinc-100 dark:bg-zinc-800 text-zinc-500 text-[10px] font-bold uppercase rounded-xl hover:bg-zinc-200 dark:hover:bg-zinc-700 transition-colors cursor-pointer"
        >
          Отмена
        </button>
      </div>
    </div>
  </div>
</template>
