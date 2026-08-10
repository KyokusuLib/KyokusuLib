<script setup lang="ts">
import { computed, ref, watch } from 'vue'

defineProps<{
  alt?: string
}>()

const model = defineModel<string | null>({ default: null })
const isZoomed = ref(false)
const isOpen = computed(() => model.value !== null)

useEventListener(window, 'keydown', (e: KeyboardEvent) => {
  if (e.key === 'Escape' && isOpen.value) {
    close()
  }
})

const isBodyLocked = useScrollLock(() => document.body)
watch(isOpen, (open) => {
  isBodyLocked.value = open
})

watch(model, () => {
  isZoomed.value = false
})

function close() {
  model.value = null
}

function toggleZoom() {
  isZoomed.value = !isZoomed.value
}
</script>

<template>
  <Teleport to="body">
    <Transition name="lightbox-fade">
      <div
        v-if="isOpen"
        role="dialog"
        aria-modal="true"
        aria-label="Просмотр изображения"
        class="fixed inset-0 z-200 bg-black/90 flex overflow-auto p-4"
        @click="close"
      >
        <button
          class="fixed top-4 right-4 p-2 text-white/70 hover:text-white transition-colors cursor-pointer"
          aria-label="Закрыть"
          @click.stop="close"
        >
          <Icon name="ph:x-bold" size="28" />
        </button>
        <img
          :src="model ?? undefined"
          :alt="alt ?? 'Изображение'"
          :class="[
            'm-auto max-w-full max-h-full object-contain rounded-lg transition-transform duration-300 select-none cursor-zoom-in',
            isZoomed && 'scale-150 cursor-zoom-out',
          ]"
          @click.stop="toggleZoom"
        />
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.lightbox-fade-enter-active,
.lightbox-fade-leave-active {
  transition: opacity 0.3s ease;
}
.lightbox-fade-enter-from,
.lightbox-fade-leave-to {
  opacity: 0;
}
</style>
