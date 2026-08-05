<script setup lang="ts">
import { ref, nextTick, onMounted, onUnmounted } from 'vue'

const props = withDefaults(
  defineProps<{
    text: string
    position?: 'top' | 'bottom' | 'left' | 'right'
    delay?: number
  }>(),
  {
    position: 'top',
    delay: 200,
  },
)

const triggerRef = ref<HTMLElement | null>(null)
const visible = ref(false)
const pos = ref({})

let timer: ReturnType<typeof setTimeout> | null = null

function reposition() {
  if (!triggerRef.value || !visible.value) return
  const r = triggerRef.value.getBoundingClientRect()
  const gap = 8

  const s: Record<string, string> = {
    position: 'fixed',
    zIndex: '100',
    pointerEvents: 'none',
    whiteSpace: 'nowrap',
  }

  switch (props.position) {
    case 'top':
      s.left = `${r.left + r.width / 2}px`
      s.top = `${r.top - gap}px`
      s.transform = 'translate(-50%, -100%)'
      break
    case 'bottom':
      s.left = `${r.left + r.width / 2}px`
      s.top = `${r.bottom + gap}px`
      s.transform = 'translateX(-50%)'
      break
    case 'left':
      s.top = `${r.top + r.height / 2}px`
      s.left = `${r.left - gap}px`
      s.transform = 'translate(-100%, -50%)'
      break
    case 'right':
      s.top = `${r.top + r.height / 2}px`
      s.left = `${r.right + gap}px`
      s.transform = 'translateY(-50%)'
      break
  }

  pos.value = s
}

function show() {
  if (timer) clearTimeout(timer)
  timer = setTimeout(() => {
    visible.value = true
    nextTick(reposition)
  }, props.delay)
}

function hide() {
  if (timer) clearTimeout(timer)
  visible.value = false
}

onMounted(() => {
  window.addEventListener('scroll', reposition, { capture: true })
  window.addEventListener('resize', reposition)
})

onUnmounted(() => {
  window.removeEventListener('scroll', reposition, { capture: true })
  window.removeEventListener('resize', reposition)
  if (timer) clearTimeout(timer)
})
</script>

<template>
  <div ref="triggerRef" class="inline-block w-fit" @mouseenter="show" @mouseleave="hide">
    <slot />
    <Teleport to="body">
      <Transition name="ttp-fade">
        <div
          v-if="visible && text"
          :style="pos"
          class="relative bg-zinc-800 dark:bg-zinc-700 text-white text-[14px] font-bold leading-none px-3 py-1.5 rounded-md shadow-lg"
        >
          {{ text }}
          <div
            class="absolute -bottom-1 left-1/2 -translate-x-1/2 rotate-45 w-2 h-2 bg-zinc-800 dark:bg-zinc-700 rounded-sm"
          ></div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<style>
.ttp-fade-enter-active,
.ttp-fade-leave-active {
  transition: opacity 0.15s ease;
}
.ttp-fade-enter-from,
.ttp-fade-leave-to {
  opacity: 0;
}
</style>
