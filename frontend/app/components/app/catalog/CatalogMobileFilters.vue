<script setup lang="ts">
import CatalogFilters from '@/components/app/catalog/CatalogFilters.vue'

defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  close: []
  apply: [params: Record<string, any>]
  reset: []
  'search-update': [params: Record<string, any>]
}>()
</script>

<template>
  <Transition name="drawer-backdrop">
    <div v-if="isOpen" class="fixed inset-0 z-40 bg-black/50 lg:hidden" @click="emit('close')" />
  </Transition>

  <Transition name="drawer-slide">
    <div
      v-if="isOpen"
      class="fixed top-0 right-0 z-50 h-full w-80 max-w-[85vw] bg-white dark:bg-zinc-900 border-l border-zinc-200 dark:border-zinc-800 shadow-2xl overflow-y-auto lg:hidden"
    >
      <div
        class="sticky top-0 z-10 bg-white dark:bg-zinc-900 px-4 pt-4 pb-2 border-b border-zinc-200 dark:border-zinc-800 flex items-center justify-between"
      >
        <h2 class="text-sm font-bold flex items-center gap-2">
          <Icon name="ph:funnel-bold" size="16" class="text-yellow-500" />
          Фильтры
        </h2>
        <button
          type="button"
          class="p-1.5 rounded-lg hover:bg-zinc-100 dark:hover:bg-zinc-800 transition-colors cursor-pointer"
          @click="emit('close')"
        >
          <Icon name="ph:x-bold" size="18" />
        </button>
      </div>
      <div class="p-4">
        <CatalogFilters
          @apply="(params) => emit('apply', params)"
          @search-update="(params) => emit('search-update', params)"
          @reset="() => emit('reset')"
        />
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.drawer-backdrop-enter-active,
.drawer-backdrop-leave-active {
  transition: opacity 0.25s ease;
}
.drawer-backdrop-enter-from,
.drawer-backdrop-leave-to {
  opacity: 0;
}

.drawer-slide-enter-active,
.drawer-slide-leave-active {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
.drawer-slide-enter-from,
.drawer-slide-leave-to {
  transform: translateX(100%);
}
</style>
