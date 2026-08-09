import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { SocialNetwork } from '@/types/enums/social-network-enum'

export const useNewsStore = defineStore('news', () => {
  const selectedSources = ref<SocialNetwork[]>([])

  function isSelected(source: SocialNetwork): boolean {
    return selectedSources.value.includes(source)
  }

  function toggleSource(source: SocialNetwork) {
    if (isSelected(source)) {
      selectedSources.value = selectedSources.value.filter((s) => s !== source)
    } else {
      selectedSources.value.push(source)
    }
  }

  function clearSources() {
    selectedSources.value = []
  }

  return {
    selectedSources,
    isSelected,
    toggleSource,
    clearSources,
  }
})
