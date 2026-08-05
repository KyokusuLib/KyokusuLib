<script setup lang="ts">
import type { PrivacySubCard } from '@/types/frontend/privacy'

defineOptions({ inheritAttrs: false })

const props = withDefaults(
  defineProps<{
    index?: number
    cardTitle?: string
    description?: string
    text?: string
    sectionId?: string

    subcards?: PrivacySubCard[]
  }>(),
  {
    cardTitle: 'Заголовок',
  },
)
</script>

<template>
  <div
    :id="sectionId"
    class="flex flex-col bg-white dark:bg-zinc-900/40 rounded-[2.5rem] border border-zinc-200 dark:border-zinc-800 p-6 shadow-sm mb-8 scroll-mt-24"
  >
    <h2 class="text-2xl font-bold mb-4">
      <span v-if="index">{{ index }}. </span>{{ cardTitle }}
    </h2>
    <p v-if="description" class="text-zinc-600 dark:text-zinc-400 mb-4 whitespace-pre-wrap">
      {{ description }}
    </p>
    <p v-if="text" class="text-zinc-600 dark:text-zinc-400 mb-4 whitespace-pre-wrap">
      {{ text }}
    </p>

    <div v-if="subcards && subcards.length" class="mt-6 space-y-6">
      <div
        v-for="(subcard, subIndex) in subcards"
        :key="subIndex"
        class="pl-4 border-l-2 border-zinc-100 dark:border-zinc-800"
      >
        <div class="font-bold text-zinc-900 dark:text-zinc-100 mb-2">
          <span v-if="index">{{ index }}.</span>{{ subIndex + 1 }} {{ subcard.title }}
        </div>
        <div class="text-sm text-zinc-600 dark:text-zinc-400 whitespace-pre-wrap">
          {{ subcard.text }}
        </div>
      </div>
    </div>
  </div>
</template>
