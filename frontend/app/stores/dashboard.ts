import { defineStore } from 'pinia'
import { ref } from 'vue'
import { DashboardTab } from '@/types/enums/dashboard-tab'

export const useDashboardStore = defineStore('dashboard', () => {
  const activeTab = ref<DashboardTab>(DashboardTab.Moderation)

  function setActiveTab(tab: DashboardTab) {
    activeTab.value = tab
  }

  return {
    activeTab,
    setActiveTab,
  }
})
