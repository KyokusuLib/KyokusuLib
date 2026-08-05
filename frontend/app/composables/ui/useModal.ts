import { onMounted, onUnmounted, watch, type Ref } from 'vue'

export function useModalLogic(isOpen: Ref<boolean>, closeCallback: () => void) {
  const handleKeydown = (e: KeyboardEvent) => {
    if (e.key === 'Escape' && isOpen.value) {
      closeCallback()
    }
  }

  watch(isOpen, (val) => {
    if (val) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
  })

  onMounted(() => {
    document.addEventListener('keydown', handleKeydown)
  })

  onUnmounted(() => {
    document.removeEventListener('keydown', handleKeydown)
    document.body.style.overflow = ''
  })
}
