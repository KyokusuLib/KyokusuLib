import { type Ref, ref } from 'vue'

export function useDraggableScroll(scrollContainer: Ref<HTMLElement | null>) {
  const isDragging = ref(false)
  const startX = ref(0)
  const scrollLeft = ref(0)

  const onMouseDown = (e: MouseEvent) => {
    if (!scrollContainer.value) return

    isDragging.value = true
    // Get initial position
    startX.value = e.pageX - scrollContainer.value.offsetLeft
    scrollLeft.value = scrollContainer.value.scrollLeft

    // Stop snap effects while dragging for smoothness
    scrollContainer.value.style.scrollSnapType = 'none'
    scrollContainer.value.style.cursor = 'grabbing'
  }

  const onMouseLeaveOrUp = () => {
    if (!scrollContainer.value) return
    isDragging.value = false

    // Re-enable snap and restore cursor
    scrollContainer.value.style.scrollSnapType = 'x mandatory'
    scrollContainer.value.style.cursor = 'grab'
  }

  const onMouseMove = (e: MouseEvent) => {
    if (!isDragging.value || !scrollContainer.value) return

    e.preventDefault()
    const x = e.pageX - scrollContainer.value.offsetLeft
    const walk = (x - startX.value) * 1.5 // * 1.5 makes scrolling faster/snappier
    scrollContainer.value.scrollLeft = scrollLeft.value - walk
  }

  // Helper to scroll Programmatically (Buttons)
  const scroll = (direction: 'left' | 'right', amount = 300) => {
    if (!scrollContainer.value) return
    scrollContainer.value.scrollBy({
      left: direction === 'left' ? -amount : amount,
      behavior: 'smooth',
    })
  }

  return {
    isDragging,
    onMouseDown,
    onMouseLeaveOrUp,
    onMouseMove,
    scroll,
  }
}
