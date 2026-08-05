<script setup lang="ts">
import { ref } from 'vue'
import { useNovelaRating } from '@/composables/api/novela/useNovelaRating'
import { ModalWindow } from '@kyokusu-ui/vue'
import { useNotificationStore } from '@/stores/notification'

defineOptions({
  inheritAttrs: false,
})

interface Props {
  modelValue: boolean
  novelaId: number
  initialRating?: number
}

const props = defineProps<Props>()
const emit = defineEmits(['update:modelValue', 'rated'])

const { setNovelaRating, loading } = useNovelaRating()
const { notify } = useNotificationStore()

const selectedRating = ref(props.initialRating || 10)
const ratings = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

const handleRatingSubmit = async () => {
  try {
    await setNovelaRating({
      novela_id: props.novelaId,
      rating: selectedRating.value,
    })
    emit('rated', selectedRating.value)
    emit('update:modelValue', false)

    notify({
      title: 'Успех',
      content: 'Спасибо за оценку!',
      type: 'success',
    })
  } catch (e) {
    console.error(e)
  }
}

const changeRating = (delta: number) => {
  const newVal = selectedRating.value + delta
  if (newVal >= 1 && newVal <= 10) {
    selectedRating.value = newVal
  }
}
</script>

<template>
  <ModalWindow
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    title="Оцените произведение"
    width="w-full max-w-xl"
    center-title
  >
    <div class="py-6 space-y-8 flex flex-col items-center">
      <div class="hidden md:block space-y-3">
        <div class="flex justify-between gap-2">
          <button
            v-for="rate in ratings"
            :key="rate"
            @click="selectedRating = rate"
            class="w-10 h-10 lg:w-11 lg:h-11 rounded-full border-2 flex items-center justify-center font-bold text-base transition-all cursor-pointer"
            :class="[
              selectedRating === rate
                ? 'bg-yellow-500 border-yellow-500 text-white shadow-lg shadow-yellow-500/20'
                : 'border-zinc-200 dark:border-zinc-800 hover:border-yellow-500 text-zinc-600 dark:text-zinc-400',
            ]"
          >
            {{ rate }}
          </button>
        </div>

        <div class="flex justify-between px-1">
          <span
            class="text-[10px] uppercase tracking-tighter font-bold text-zinc-400 dark:text-zinc-500"
          >
            Очень плохо
          </span>
          <span
            class="text-[10px] uppercase tracking-tighter font-bold text-zinc-400 dark:text-zinc-500"
          >
            Невероятно
          </span>
        </div>
      </div>

      <div
        class="flex md:hidden items-center justify-between bg-zinc-100 dark:bg-zinc-800/50 p-4 gap-4 rounded-2xl"
      >
        <button
          @click="changeRating(-1)"
          :disabled="selectedRating <= 1"
          class="w-12 h-12 flex items-center justify-center rounded-xl bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 disabled:opacity-30"
        >
          <Icon name="ph:star-fill" class="text-yellow-500" size="12" />
        </button>

        <div class="text-center">
          <div class="text-4xl font-black text-zinc-900 dark:text-white">{{ selectedRating }}</div>
          <div class="text-[10px] uppercase tracking-widest text-zinc-500 mt-1">Ваша оценка</div>
        </div>

        <button
          @click="changeRating(1)"
          :disabled="selectedRating >= 10"
          class="w-12 h-12 flex items-center justify-center rounded-xl bg-white dark:bg-zinc-800 border border-zinc-200 dark:border-zinc-700 disabled:opacity-30"
        >
          <Icon name="ph:star-fill" class="text-yellow-500" size="12" />
        </button>
      </div>

      <button
        @click="handleRatingSubmit"
        :disabled="loading"
        class="md:w-auto w-full py-2 px-4 bg-yellow-500 hover:bg-yellow-600 text-white font-bold rounded-2xl transition-all active:scale-95 disabled:opacity-50 flex items-center justify-center gap-3 shadow-xl shadow-yellow-500/10"
      >
        <span
          v-if="loading"
          class="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin"
        ></span>
        Подтвердить оценку
      </button>
    </div>
  </ModalWindow>
</template>

<style scoped>
button {
  cursor: pointer;
  -webkit-tap-highlight-color: transparent;
}
</style>
