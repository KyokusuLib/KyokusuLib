<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import NovelaRatingModal from './NovelaRatingModal.vue'
import AuthRequiredModal from '~/components/common/AuthRequiredModal.vue'

const props = defineProps<{
  novelaId: number
  rating: number
  userRating: number
  count: number
}>()

const { isAuthenticated } = useAuthStore()

const emit = defineEmits(['update:rated'])

const isRatingModalOpen = ref(false)
const isAuthModalOpen = ref(false)

const starColor = computed(() => {
  if (props.rating >= 8) return 'text-green-500'
  if (props.rating >= 6) return 'text-yellow-500'
  return 'text-red-500'
})

const handleOpenRating = () => {
  if (!isAuthenticated) {
    isAuthModalOpen.value = true
    return
  }
  isRatingModalOpen.value = !isRatingModalOpen.value
}

const updateRating = (rating: number) => {
  emit('update:rated', rating)
}
</script>

<template>
  <div class="flex px-1 py-1.5 rounded-lg">
    <div class="flex flex-row gap-2 leading-none">
      <span class="text-3xl font-bold text-zinc-200 text-right" :style="{ color: starColor }">{{
        roundTo(rating, 1)
      }}</span>
      <div class="flex flex-col">
        <span class="text-[10px] text-white dark:text-zinc-400 uppercase">{{ count }} оценок</span>
        <button
          class="flex items-end gap-1 text-xs font-bold text-white dark:text-zinc-400 hover:text-yellow-600 dark:hover:text-yellow-400 cursor-pointer"
          @click="handleOpenRating"
        >
          Оценить
          <Icon name="ph:star-fill" class="text-yellow-500" size="12" />
        </button>
        <NovelaRatingModal
          v-model="isRatingModalOpen"
          :novela-id="novelaId"
          :initial-rating="userRating"
          @rated="updateRating"
        />
        <AuthRequiredModal v-model="isAuthModalOpen" action-text="ставить оценки" />
      </div>
    </div>
  </div>
</template>
