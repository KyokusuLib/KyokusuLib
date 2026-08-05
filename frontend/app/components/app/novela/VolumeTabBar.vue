<script setup lang="ts">
import { ref, watch, nextTick } from 'vue'
import type { NovelaVolume } from '~/types/backend/novela'
import { useDraggableScroll } from '~/composables/ui/useDraggablescroll'

const props = defineProps<{
  volumes: NovelaVolume[]
  activeVolumeId: string
}>()

const emit = defineEmits<{
  select: [id: string]
}>()

const scrollContainer = ref<HTMLElement | null>(null)
const { onMouseDown, onMouseLeaveOrUp, onMouseMove } = useDraggableScroll(scrollContainer)

function onWheel(e: WheelEvent) {
  const el = scrollContainer.value
  if (!el) return
  el.scrollLeft = Math.max(0, el.scrollLeft + e.deltaY)
}

watch(
  () => props.activeVolumeId,
  async () => {
    await nextTick()
    if (!scrollContainer.value) return
    const activeBtn = scrollContainer.value.querySelector(
      `button[data-vol-id="${props.activeVolumeId}"]`,
    ) as HTMLElement
    if (activeBtn) {
      const container = scrollContainer.value
      const scrollLeft =
        activeBtn.offsetLeft -
        container.offsetLeft -
        container.clientWidth / 2 +
        activeBtn.clientWidth / 2
      container.scrollTo({ left: scrollLeft, behavior: 'smooth' })
    }
  },
)
</script>

<template>
  <div
    ref="scrollContainer"
    class="flex gap-2 overflow-x-auto no-scrollbar min-w-0 max-w-full items-center cursor-grab active:cursor-grabbing"
    @mousedown="onMouseDown"
    @mouseleave="onMouseLeaveOrUp"
    @mouseup="onMouseLeaveOrUp"
    @mousemove="onMouseMove"
    @wheel.prevent="onWheel"
  >
    <button
      v-for="vol in volumes"
      :key="vol.id"
      :data-vol-id="vol.id"
      @click="emit('select', vol.id)"
      class="px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors border cursor-pointer"
      :class="
        activeVolumeId === vol.id
          ? 'bg-zinc-900 text-white border-zinc-900 dark:bg-zinc-100 dark:text-zinc-900 dark:border-zinc-100'
          : 'bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-700 text-zinc-600 dark:text-zinc-400 hover:border-zinc-400'
      "
    >
      Том {{ vol.number }}
    </button>
  </div>
</template>

<style scoped>
.no-scrollbar::-webkit-scrollbar {
  display: none;
}
.no-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
