<script setup lang="ts">
import { Toggle as BaseToggle } from '@kyokusu-ui/vue'

interface Props {
  title: string
  description: string
  imgSrc?: string
  id?: string
}

withDefaults(defineProps<Props>(), {
  imgSrc: '',
})

const modelValue = defineModel<boolean>({ default: false })
</script>

<template>
  <div
    class="flex flex-row items-center justify-between gap-4 p-4 rounded-lg border bg-white dark:bg-zinc-900/50 dark:border-zinc-800"
  >
    <div class="flex flex-col gap-2">
      <div class="flex items-center gap-4">
        <span class="text-sm font-medium">{{ title }}</span>
        <Icon :name="imgSrc" size="18" class="text-zinc-700 dark:text-zinc-200" />
      </div>
      <p class="text-xs text-zinc-500 dark:text-stone-500">
        {{ description }}
      </p>
    </div>
    <BaseToggle v-model="modelValue" :id="id" />
  </div>
</template>
