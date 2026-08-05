<script setup lang="ts">
import ReaderSettingsFontSize from '~/components/app/novela/ReaderSettingsFontSize.vue'
import ReaderSettingsFontFamily from '~/components/app/novela/ReaderSettingsFontFamily.vue'
import ReaderSettingsLineWeight from '~/components/app/novela/ReaderSettingsLineWeight.vue'
import ReaderSettingsAutoScroll from '~/components/app/novela/ReaderSettingsAutoScroll.vue'
import { Tooltip } from '@kyokusu-ui/vue'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
}>()

const close = () => emit('update:modelValue', false)

useEventListener(window, 'keydown', (e: KeyboardEvent) => {
  if (e.key === 'Escape' && props.modelValue) {
    close()
  }
})
</script>

<template>
  <Transition name="slide-fade">
    <div v-if="modelValue" class="fixed inset-0 z-100 flex justify-end">
      <div
        class="absolute inset-0 bg-black/20 dark:bg-black/40 backdrop-blur-[2px]"
        @click="close"
      ></div>

      <div
        class="relative w-80 h-full bg-white dark:bg-zinc-900 border-l border-zinc-200 dark:border-zinc-800 shadow-2xl flex flex-col"
      >
        <div class="shrink-0 px-5 pt-5 pb-3 border-b border-zinc-100 dark:border-zinc-800">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="w-9 h-9 rounded-xl bg-yellow-500/10 flex items-center justify-center">
                <Icon
                  name="ph:gear-six-bold"
                  size="18"
                  class="text-yellow-600 dark:text-yellow-400"
                />
              </div>
              <div class="flex flex-col">
                <h2 class="text-base font-black tracking-tight">Настройки</h2>
                <span class="text-[10px] text-zinc-400 font-bold uppercase tracking-widest"
                  >Режим чтения</span
                >
              </div>
            </div>
            <Tooltip text="Закрыть" position="left">
              <button
                @click="close"
                class="p-2 hover:bg-zinc-100 dark:hover:bg-zinc-800 rounded-xl transition-colors cursor-pointer"
              >
                <Icon name="ph:x-bold" size="18" />
              </button>
            </Tooltip>
          </div>
        </div>

        <div class="flex-1 overflow-y-auto custom-scrollbar p-5 space-y-3">
          <ReaderSettingsFontSize />
          <ReaderSettingsLineWeight />
          <ReaderSettingsFontFamily />
          <ReaderSettingsAutoScroll />
        </div>
      </div>
    </div>
  </Transition>
</template>

<style scoped>
.slide-fade-enter-active,
.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.slide-fade-enter-from,
.slide-fade-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

.slide-fade-enter-to,
.slide-fade-leave-from {
  transform: translateX(0);
  opacity: 1;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #e4e4e7;
  border-radius: 20px;
}
:deep(.dark) .custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #27272a;
}
</style>
