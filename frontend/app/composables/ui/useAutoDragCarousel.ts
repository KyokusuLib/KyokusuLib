import { nextTick, onMounted, onUnmounted, type Ref, ref } from 'vue'

export function useAutoDragCarousel(containerRef: Ref<HTMLElement | null>, intervalSeconds = 5) {
  const isDown = ref(false)
  const startX = ref(0)
  const scrollLeft = ref(0)
  let autoPlayInterval: ReturnType<typeof setInterval> | null = null

  const alignToCenter = () => {
    if (!containerRef.value) return
    const container = containerRef.value
    const firstCard = container.firstElementChild as HTMLElement

    if (firstCard) {
      const padding = (container.clientWidth - firstCard.offsetWidth) / 261

      container.style.paddingLeft = `${padding}px`
      container.style.paddingRight = `${padding}px`
    }
  }

  const nextSlide = () => {
    if (!containerRef.value) return
    const container = containerRef.value

    const cardWidth = container.children[0]?.clientWidth || 220
    const gap = 16

    if (container.scrollLeft + container.clientWidth >= container.scrollWidth - 10) {
      container.scrollTo({ left: 0, behavior: 'smooth' })
    } else {
      container.scrollBy({ left: cardWidth + gap, behavior: 'smooth' })
    }
  }

  const startAutoPlay = () => {
    stopAutoPlay()
    autoPlayInterval = setInterval(nextSlide, intervalSeconds * 1000)
  }

  const stopAutoPlay = () => {
    if (autoPlayInterval) {
      clearInterval(autoPlayInterval)
      autoPlayInterval = null
    }
  }

  const onMouseDown = (e: MouseEvent) => {
    if (!containerRef.value) return
    isDown.value = true
    stopAutoPlay()
    startX.value = e.pageX - containerRef.value.offsetLeft
    scrollLeft.value = containerRef.value.scrollLeft

    containerRef.value.style.scrollBehavior = 'auto'
    containerRef.value.style.scrollSnapType = 'none'
    containerRef.value.style.cursor = 'grabbing'
  }

  const onMouseLeaveOrUp = () => {
    if (!isDown.value || !containerRef.value) return
    isDown.value = false
    startAutoPlay()
    containerRef.value.style.scrollBehavior = 'smooth'
    containerRef.value.style.scrollSnapType = 'x mandatory'
    containerRef.value.style.cursor = 'grab'
  }

  const onMouseMove = (e: MouseEvent) => {
    if (!isDown.value || !containerRef.value) return
    e.preventDefault()
    const x = e.pageX - containerRef.value.offsetLeft
    const walk = (x - startX.value) * 1
    containerRef.value.scrollLeft = scrollLeft.value - walk
  }

  onMounted(async () => {
    await nextTick()
    alignToCenter()

    window.addEventListener('resize', alignToCenter)

    startAutoPlay()
  })

  onUnmounted(() => {
    window.removeEventListener('resize', alignToCenter)
    stopAutoPlay()
  })

  return {
    onMouseDown,
    onMouseLeaveOrUp,
    onMouseMove,
    stopAutoPlay,
    startAutoPlay,
  }
}
