<script setup lang="ts">
interface Props {
  text?: string | number
  icon?: string
  border?: boolean
  colorClass?: string
  size?: string
}

const props = withDefaults(defineProps<Props>(), {
  text: '',
  icon: '',
  border: false,
  colorClass: '',
  size: 'w-8 h-8',
})
</script>

<template>
  <div
    class="flex items-center justify-center rounded-full shadow-sm shrink-0"
    :class="[
      size,
      colorClass || 'bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900',
      border ? 'ring-2 ring-zinc-300 dark:ring-amber-500' : '',
    ]"
  >
    <span v-if="text" class="font-black italic leading-none">{{ text }}</span>
    <img v-if="icon && !text" :src="icon" class="p-1.5" alt="" />
  </div>
</template>
