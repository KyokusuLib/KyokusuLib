<script setup lang="ts">
import type { SearchCategory, SearchResultItem } from '@/composables/api/search/useSearch'
import SearchResultItemCard from './SearchResultItemCard.vue'

defineProps<{
  results: SearchResultItem[]
  activeCategory: SearchCategory
}>()

const emit = defineEmits(['select'])
</script>

<template>
  <template v-for="item in results" :key="item.id">
    <SearchResultItemCard
      :item="item"
      :activeCategory="activeCategory"
      @select="emit('select', $event)"
    />
  </template>
</template>
