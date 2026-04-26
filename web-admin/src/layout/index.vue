<template>
  <a-layout class="layout-root layout-root--fill">
    <a-layout-sider
      v-show="isDesktop"
      class="layout-sider layout-sider--fill !bg-[var(--color-bg-1)]"
      :width="220"
      :collapsed-width="64"
      collapsible
      breakpoint="lg"
      :collapsed="layout.sidebarCollapsed"
      :hide-trigger="true"
      @update:collapsed="layout.setSidebarCollapsed"
    >
      <LayoutSidebar :collapsed="layout.sidebarCollapsed" @navigate="onSidebarNavigate" />
    </a-layout-sider>

    <a-layout class="layout-main-column">
      <a-layout-header
        class="layout-main-header !h-[52px] shrink-0 !p-0 !leading-[52px] !bg-[var(--color-bg-2)] !border-b !border-[var(--color-border-2)]"
      >
        <LayoutHeader @toggle-mobile-menu="layout.setMobileDrawer(true)" />
      </a-layout-header>
      <a-layout-content class="layout-content layout-content--scroll">
        <LayoutAppMain />
      </a-layout-content>
    </a-layout>

    <a-drawer
      :visible="layout.mobileDrawerOpen"
      placement="left"
      :width="268"
      :footer="false"
      unmount-on-close
      @update:visible="layout.setMobileDrawer"
    >
      <template #title>菜单</template>
      <LayoutSidebar :collapsed="false" @navigate="onSidebarNavigate" />
    </a-drawer>
  </a-layout>
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { useBreakpoints } from '@vueuse/core'
import LayoutAppMain from './AppMain.vue'
import LayoutHeader from './Header.vue'
import LayoutSidebar from './Sidebar.vue'
import { useLayoutStore } from '@/store/layout'

const layout = useLayoutStore()
const breakpoints = useBreakpoints({ md: 768 })
const isDesktop = breakpoints.greaterOrEqual('md')

function onSidebarNavigate(): void {
  layout.setMobileDrawer(false)
}

watch(isDesktop, (v) => {
  if (v) layout.setMobileDrawer(false)
})
</script>
