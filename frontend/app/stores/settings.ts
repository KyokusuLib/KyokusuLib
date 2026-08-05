import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useSettingsStore = defineStore('settings', () => {
  const activeTab = ref<string>('account')

  function setActiveTab(id: string) {
    activeTab.value = id
  }

  return { activeTab, setActiveTab }
})
