<script setup lang="ts">
import LogoutIcon from '@/assets/images/special/log-out.png'

interface Props {
  name: string
  iconSrc: string
  isConnected: boolean
  bgClasses: string
}

const props = defineProps<Props>()
const emit = defineEmits(['click'])
</script>

<template>
  <div
    @click="emit('click')"
    class="group flex flex-col items-start justify-between p-4 rounded-lg border transition-all duration-300 bg-white border-zinc-200 hover:border-zinc-300 hover:shadow-md cursor-pointer relative overflow-hidden"
    :class="bgClasses"
  >
    <slot name="background"></slot>

    <div class="space-y-1 flex flex-col w-full relative z-10">
      <div class="flex items-center justify-between w-full">
        <span class="text-lg font-medium block text-zinc-700 dark:text-zinc-200">{{ name }}</span>
        <img
          :src="iconSrc"
          class="w-8 h-8 transition-transform group-hover:scale-110 dark:invert"
          :class="{ 'w-6 h-6 invert-0': name === 'Google' }"
        />
      </div>

      <div v-if="isConnected" class="flex items-center gap-4">
        <img
          :src="LogoutIcon"
          class="w-6 h-6 transition-transform group-hover:scale-110 dark:invert"
        />
        <p class="text-xs text-zinc-500 dark:text-stone-300">Уже подключен</p>
      </div>
      <div v-else>
        <p class="text-xs text-zinc-500 dark:text-stone-300">Подключить аккаунт {{ name }}</p>
      </div>
    </div>
  </div>
</template>
