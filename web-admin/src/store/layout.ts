import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLayoutStore = defineStore('layout', () => {
  const sidebarCollapsed = ref(false)
  const mobileDrawerOpen = ref(false)

  function toggleSidebar(): void {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function setSidebarCollapsed(v: boolean): void {
    sidebarCollapsed.value = v
  }

  function setMobileDrawer(open: boolean): void {
    mobileDrawerOpen.value = open
  }

  return {
    sidebarCollapsed,
    mobileDrawerOpen,
    toggleSidebar,
    setSidebarCollapsed,
    setMobileDrawer,
  }
})
