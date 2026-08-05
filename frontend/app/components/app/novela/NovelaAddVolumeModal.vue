<script setup lang="ts">
import { ModalWindow } from '@kyokusu-ui/vue'
import { ref, watch } from 'vue'
import type { NovelaVolume } from '@/types/backend/novela'
import { useNovela } from '~/composables/api/novela/useNovela'

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  modelValue: boolean
  novelaId: number
  volumes: NovelaVolume[]
}>()

const emit = defineEmits(['update:modelValue', 'volume-added'])

const { addVolume } = useNovela()
const loading = ref(false)

const title = ref('')

const maxVolumeNumber = computed(() => {
  if (!props.volumes || props.volumes.length === 0) return 1
  const numbers = props.volumes.map((v) => v.number)
  return Math.max(...numbers) + 1
})

const volumeNumber = ref(maxVolumeNumber.value)

watch(maxVolumeNumber, (next) => {
  if (next > volumeNumber.value) {
    volumeNumber.value = next
  }
})

watch(
  () => props.modelValue,
  (open) => {
    if (open) {
      volumeNumber.value = maxVolumeNumber.value
      title.value = ''
    }
  },
)

const handleSubmit = async () => {
  loading.value = true
  try {
    const volume = await addVolume(props.novelaId, volumeNumber.value, title.value)
    if (volume) {
      emit('volume-added', volume)
    } else {
      emit('volume-added')
    }
    emit('update:modelValue', false)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <ModalWindow
    v-if="modelValue"
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="Добавить том"
    width="w-full max-w-md"
    center-title
  >
    <div class="py-6 space-y-4">
      <div class="space-y-2 flex flex-col gap-0.5">
        <label class="text-sm font-semibold text-zinc-700 dark:text-zinc-300 ml-1">
          Номер тома
        </label>
        <input
          v-model.number="volumeNumber"
          type="number"
          :min="maxVolumeNumber"
          class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700 rounded-xl text-zinc-900 dark:text-zinc-100 outline-none focus:ring-2 focus:ring-yellow-500 transition-all"
          placeholder="Например, 1"
        />
      </div>

      <div class="space-y-2 flex flex-col gap-0.5">
        <label class="text-sm font-semibold text-zinc-700 dark:text-zinc-300 ml-1">
          Название (опционально)
        </label>
        <input
          v-model="title"
          type="text"
          class="w-full px-4 py-2 bg-zinc-100 dark:bg-zinc-800/50 border border-zinc-200 dark:border-zinc-700 rounded-xl text-zinc-900 dark:text-zinc-100 outline-none focus:ring-2 focus:ring-yellow-500 transition-all"
          placeholder="Например, Начало пути"
        />
      </div>

      <div class="pt-4 flex justify-end gap-3">
        <button
          @click="$emit('update:modelValue', false)"
          class="px-4 py-2 text-sm font-semibold text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-xl transition-colors cursor-pointer"
        >
          Отмена
        </button>
        <button
          @click="handleSubmit"
          :disabled="loading || volumeNumber < 1"
          class="px-4 py-2 text-sm font-bold bg-yellow-600 hover:bg-yellow-700 text-white rounded-xl transition-all disabled:opacity-50 flex items-center gap-2 cursor-pointer"
        >
          <span
            v-if="loading"
            class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"
          ></span>
          Добавить
        </button>
      </div>
    </div>
  </ModalWindow>
</template>
