import { onMounted, onUnmounted, type Ref, ref } from 'vue'

export function useCarousel(itemCount: Ref<number>, intervalDuration = 5000) {
  const currentSlide = ref(0)
  let autoPlayTimer: ReturnType<typeof setInterval> | null = null

  const touchStartX = ref(0)
  const touchEndX = ref(0)
  const minSwipeDistance = 50

  const next = () => {
    currentSlide.value = (currentSlide.value + 1) % itemCount.value
  }

  const prev = () => {
    currentSlide.value = (currentSlide.value - 1 + itemCount.value) % itemCount.value
  }

  const goTo = (index: number) => {
    currentSlide.value = index
  }

  const startAutoPlay = () => {
    if (intervalDuration <= 0) return
    stopAutoPlay()
    autoPlayTimer = setInterval(next, intervalDuration)
  }

  const stopAutoPlay = () => {
    if (autoPlayTimer) {
      clearInterval(autoPlayTimer)
      autoPlayTimer = null
    }
  }

  const onTouchStart = (e: TouchEvent) => {
    touchStartX.value = e.changedTouches[0]!.screenX
    stopAutoPlay()
  }

  const onTouchEnd = (e: TouchEvent) => {
    touchEndX.value = e.changedTouches[0]!.screenX
    handleSwipe()
    startAutoPlay()
  }

  const handleSwipe = () => {
    const distance = touchStartX.value - touchEndX.value
    if (distance > minSwipeDistance) next()
    if (distance < -minSwipeDistance) prev()
  }

  onMounted(() => {
    startAutoPlay()
  })

  onUnmounted(() => {
    stopAutoPlay()
  })

  return {
    currentSlide,
    next,
    prev,
    goTo,
    startAutoPlay,
    stopAutoPlay,
    onTouchStart,
    onTouchEnd,
  }
}
