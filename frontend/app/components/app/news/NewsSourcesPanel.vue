<script setup lang="ts">
import { onMounted } from 'vue'
import { Card } from '@kyokusu-ui/vue'
import { useNewsSources } from '@/composables/news/useNewsSources'

const { sources, isSelected, toggleSource, isLoading, loadSources } = useNewsSources()

onMounted(() => {
  loadSources()
})
</script>

<template>
  <Card>
    <template #header>
      <div class="flex items-center justify-center gap-4">
        <h2 class="text-xl font-bold flex justify-center items-center">Источники</h2>
        <Icon name="ph:article-bold" size="24" />
      </div>
    </template>

    <div v-if="isLoading" class="flex flex-col gap-4" aria-hidden="true">
      <div
        v-for="i in 2"
        :key="i"
        class="flex items-center gap-2 px-4 py-3 rounded-xl"
      >
        <div class="h-6 w-6 rounded-full bg-zinc-200 dark:bg-zinc-800 animate-pulse shrink-0"></div>
        <div class="h-4 flex-1 rounded bg-zinc-200 dark:bg-zinc-800 animate-pulse"></div>
        <div class="h-5 w-5 rounded-full bg-zinc-200 dark:bg-zinc-800 animate-pulse shrink-0"></div>
      </div>
    </div>

    <div v-else class="flex flex-col gap-4">
      <button
        v-for="source in sources"
        :key="source.id"
        type="button"
        :aria-pressed="isSelected(source.id)"
        class="flex items-center gap-2 px-4 py-3 rounded-xl cursor-pointer border transition-colors text-left"
        :class="
          isSelected(source.id)
            ? 'bg-zinc-200 dark:bg-zinc-800 border-zinc-400 dark:border-zinc-600 text-zinc-900 dark:text-white'
            : 'bg-transparent border-transparent hover:bg-zinc-100 dark:hover:bg-zinc-800/50 text-zinc-600 dark:text-zinc-300'
        "
        @click="toggleSource(source.id)"
      >
        <Icon
          :name="source.icon"
          size="24"
          :class="isSelected(source.id) ? 'text-yellow-500' : 'opacity-70'"
        />
        <span class="font-medium flex-1">{{ source.name }}</span>
        <Icon
          :name="isSelected(source.id) ? 'ph:check-circle-fill' : 'ph:circle'"
          size="20"
          :class="isSelected(source.id) ? 'text-yellow-500' : 'text-zinc-400 dark:text-zinc-600'"
        />
      </button>
    </div>
  </Card>
</template>
