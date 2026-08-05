<script setup lang="ts">
import { ref } from 'vue'
import { onClickOutside } from '@vueuse/core'
import { useProfile } from '@/composables/api/profile/useProfile'
import { useUserTag } from '@/composables/api/profile/useUserTag'
import type { UserTagDTO } from '@/types/backend/user'

interface Props {
  modelValue: string
  tags: UserTagDTO[]
}

const props = defineProps<Props>()
const emit = defineEmits(['update:modelValue'])

const { isSelfProfile } = useProfile()
const { updateUserTag } = useUserTag()

const isOpen = ref(false)
const containerRef = ref<HTMLElement | null>(null)

onClickOutside(containerRef, () => (isOpen.value = false))

const selectTag = async (dto: UserTagDTO) => {
  isOpen.value = false
  emit('update:modelValue', dto.tag)
  await updateUserTag(dto)
}
</script>

<template>
  <div class="relative inline-block mt-2" ref="containerRef">
    <button
      @click="isSelfProfile ? (isOpen = !isOpen) : null"
      type="button"
      class="flex items-center gap-2 px-4 py-1.5 h-9 rounded-full border transition-all duration-300 select-none group"
      :class="[
        isSelfProfile
          ? 'cursor-pointer bg-white dark:bg-zinc-900 border-zinc-200 dark:border-zinc-800 hover:border-yellow-500/50 hover:shadow-lg hover:shadow-yellow-500/5'
          : 'cursor-default bg-zinc-50 dark:bg-zinc-800/30 border-transparent text-zinc-500',
      ]"
    >
      <Icon
        name="ph:hash-bold"
        size="14"
        class="text-zinc-400 group-hover:text-yellow-500 transition-colors"
      />

      <span class="text-sm font-bold tracking-tight text-zinc-700 dark:text-zinc-200">
        {{ modelValue || 'Выбрать тег' }}
      </span>

      <Icon
        v-if="isSelfProfile"
        name="ph:caret-down-bold"
        size="12"
        class="text-zinc-400 transition-transform duration-300"
        :class="{ 'rotate-180 text-yellow-500': isOpen }"
      />
    </button>

    <!-- ВЫПАДАЮЩИЙ СПИСОК -->
    <Transition name="dropdown">
      <div
        v-if="isOpen && isSelfProfile"
        class="absolute left-0 mt-3 w-64 bg-white/90 dark:bg-zinc-900/90 backdrop-blur-xl border border-zinc-200 dark:border-zinc-800 rounded-2xl shadow-2xl z-50 overflow-hidden shadow-black/20"
      >
        <div class="px-4 py-3 border-b flex justify-center border-zinc-100 dark:border-zinc-800/50">
          <p class="text-[10px] font-black uppercase tracking-[0.2em] text-zinc-400">
            Ваши доступные теги
          </p>
        </div>

        <div class="max-h-64 overflow-y-auto custom-scrollbar p-1.5">
          <div v-if="tags.length === 0" class="py-8 text-center">
            <Icon
              name="ph:seal-question-light"
              size="32"
              class="text-zinc-300 dark:text-zinc-700 mb-2 mx-auto"
            />
            <p class="text-xs text-zinc-500">Список пуст</p>
          </div>

          <div v-else class="space-y-1">
            <button
              v-for="tagItem in tags"
              :key="tagItem.id"
              @click="selectTag(tagItem)"
              class="w-full flex items-center cursor-pointer justify-center px-3 py-2.5 rounded-xl text-sm transition-all duration-200 group/item"
              :class="[
                tagItem.tag === modelValue
                  ? 'bg-yellow-500/10 text-yellow-600 dark:text-yellow-500 font-bold'
                  : 'text-zinc-600 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 hover:text-zinc-900 dark:hover:text-zinc-100',
              ]"
            >
              <div class="flex items-center gap-2">
                <span
                  class="opacity-0 group-hover/item:opacity-100 transition-opacity text-yellow-500"
                  >#</span
                >
                <span class="mr-3">{{ tagItem.tag }}</span>
              </div>

              <Icon v-if="tagItem.tag === modelValue" name="ph:check-circle-fill" size="16" />
            </button>
          </div>
        </div>

        <div
          class="p-2 bg-zinc-50 dark:bg-zinc-800/40 border-t border-zinc-100 dark:border-zinc-800/50"
        >
          <p class="text-[9px] text-center text-zinc-400">Теги выдаются за особые достижения</p>
        </div>
      </div>
    </Transition>
  </div>
</template>

<style scoped>
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.3s cubic-bezier(0.23, 1, 0.32, 1);
}
.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.98);
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #3f3f46;
  border-radius: 10px;
}
</style>
