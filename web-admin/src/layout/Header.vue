<template>
  <div class="layout-header">
    <a-button type="text" class="md:!hidden" @click="$emit('toggle-mobile-menu')">
      <template #icon><IconMenu /></template>
    </a-button>

    <a-button type="text" class="!hidden md:!inline-flex" @click="layout.toggleSidebar()">
      <template #icon>
        <IconMenuFold v-if="!layout.sidebarCollapsed" />
        <IconMenuUnfold v-else />
      </template>
    </a-button>

    <a-breadcrumb class="min-w-0 flex-1 !hidden sm:!flex">
      <a-breadcrumb-item v-for="(b, i) in crumbs" :key="i">
        <RouterLink v-if="b.to && i < crumbs.length - 1" :to="b.to" class="text-[var(--color-text-2)] hover:text-[rgb(var(--primary-6))]">
          {{ b.label }}
        </RouterLink>
        <span v-else class="text-[var(--color-text-1)] font-medium">{{ b.label }}</span>
      </a-breadcrumb-item>
    </a-breadcrumb>

    <div class="flex items-center gap-2 shrink-0 ml-auto">
      <TenantSwitcher />
      <a-button type="outline" size="small" @click="onLogout">退出</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { IconMenu, IconMenuFold, IconMenuUnfold } from '@arco-design/web-vue/es/icon'
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import TenantSwitcher from './TenantSwitcher.vue'
import { useLayoutStore } from '@/store/layout'
import { useUserStore } from '@/store/user'

defineEmits<{ 'toggle-mobile-menu': [] }>()

const route = useRoute()
const layout = useLayoutStore()
const userStore = useUserStore()

const crumbs = computed(() =>
  route.matched
    .filter((r) => r.meta?.title)
    .map((r) => ({
      label: String(r.meta.title),
      to: r.path && r.path !== '/' ? { path: r.path } : undefined,
    })),
)

async function onLogout(): Promise<void> {
  await userStore.logout()
}
</script>
