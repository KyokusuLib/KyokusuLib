<script setup lang="ts">
import { Badge } from '@kyokusu-ui/vue'
import type { NewsSource } from '@/types/frontend/news'

defineProps<{
  sources: NewsSource[]
}>()
</script>

<template>
  <div class="flex flex-wrap gap-2 justify-center">
    <Badge
      v-for="source in sources"
      :key="source.id"
      variant="secondary"
      size="sm"
      :text="source.name"
      :style="[source.styles, { '--k-badge-text': '#fff' }]"
    />
  </div>
</template>
