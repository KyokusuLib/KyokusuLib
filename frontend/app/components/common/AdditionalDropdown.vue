<script setup lang="ts">
import { onClickOutside } from '@vueuse/core'
import { ref } from 'vue'
import { Separator } from '@kyokusu-ui/vue'

const router = useRouter()

const isOpen = ref(false)
const dropdownRef = ref<HTMLElement | null>(null)
let closeTimer: ReturnType<typeof setTimeout> | null = null

onClickOutside(dropdownRef, () => {
  clearCloseTimer()
  isOpen.value = false
})

const clearCloseTimer = () => {
  if (closeTimer) {
    clearTimeout(closeTimer)
    closeTimer = null
  }
}

const scheduleClose = () => {
  clearCloseTimer()
  closeTimer = setTimeout(() => {
    isOpen.value = false
    closeTimer = null
  }, 200)
}

const cancelClose = () => {
  clearCloseTimer()
}

const goToNews = () => {
  clearCloseTimer()
  isOpen.value = false
  router.push('/news')
}
</script>

<template>
  <div
    class="relative"
    ref="dropdownRef"
    @mouseleave="scheduleClose"
    @mouseenter="cancelClose"
  >
    <div
      @click="isOpen = !isOpen"
      class="flex items-center gap-2 px-6 py-3 rounded-2xl cursor-pointer bg-zinc-300 hover:bg-zinc-200 dark:bg-zinc-800/80 dark:hover:bg-zinc-700 transition-colors font-medium text-zinc-700 hover:text-zinc-900 dark:text-zinc-200 dark:hover:text-white"
    >
      <span>...</span>
    </div>

    <Transition name="fade">
      <div
        v-if="isOpen"
        class="absolute left-0 mt-3 w-96 gap-2 bg-white dark:bg-[#18181b] border border-zinc-200 dark:border-zinc-700 rounded-xl shadow-xl py-2 px-2 flex flex-col z-60"
      >
        <div
          @click="goToNews"
          class="
          flex items-center
          rounded-lg cursor-pointer 
          px-4 py-2.5
          text-zinc-600 
          hover:bg-zinc-100 hover:text-zinc-900 dark:text-zinc-300 dark:hover:bg-zinc-700 dark:hover:text-white transition-colors
          border border-zinc-600
          "
        >
          <Icon name="ph:newspaper-bold" size="24" />
          <Separator orientation="vertical" style="height: 40px;"/>
          <div class="flex flex-col">
              <span class="text-sm font-extrabold">Новости</span>
              <span class="text-xs italic">просмотреть новости с разных соцсетей</span>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(-5px);
}
</style>
