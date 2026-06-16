<template>
  <a-menu
    v-if="!nested"
    mode="vertical"
    :collapsed="collapsed"
    :selected-keys="selectedKeys"
    auto-open-selected
    class="!border-none !bg-transparent"
    @menu-item-click="onMenuItemClick"
  >
    <ArcoProMenu :nested="true" :items="items" :collapsed="collapsed" @navigate="emit('navigate')" />
  </a-menu>
  <template v-else>
    <template v-for="node in items" :key="'node-' + node.id">
      <a-sub-menu v-if="node.type === 1 && visibleChildren(node).length" :key="'sub-' + node.id">
        <template #title>
          <span class="inline-flex items-center gap-2">
            <component :is="iconCmp(node.icon)" v-if="iconCmp(node.icon)" class="text-[16px]" />
            {{ node.title }}
          </span>
        </template>
        <ArcoProMenu :nested="true" :items="visibleChildren(node)" :collapsed="collapsed" @navigate="emit('navigate')" />
      </a-sub-menu>
      <a-menu-item v-else-if="node.type === 2 && node.path" :key="node.path">
        <template #icon>
          <component :is="iconCmp(node.icon)" v-if="iconCmp(node.icon)" />
        </template>
        {{ node.title }}
      </a-menu-item>
    </template>
  </template>
</template>

<script setup lang="ts">
import type { Component } from 'vue'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import * as ArcoIcons from '@arco-design/web-vue/es/icon'
import type { Menu } from '@/api/types'
import ArcoProMenu from './ArcoProMenu.vue'

withDefaults(
  defineProps<{
    items: Menu[]
    collapsed?: boolean
    nested?: boolean
  }>(),
  { collapsed: false, nested: false },
)

const emit = defineEmits<{ navigate: [] }>()

const route = useRoute()
const router = useRouter()

const selectedKeys = computed(() => {
  const p = route.path
  return p && p !== '/' ? [p] : []
})

const iconComponents = ArcoIcons as Record<string, Component>
const legacyIconMap: Record<string, string> = {
  setting: 'IconSettings',
  user: 'IconUser',
  usergroup: 'IconUserGroup',
  read: 'IconBook',
  picture: 'IconImage',
  document: 'IconFile',
  apps: 'IconApps',
}

function iconCmp(name?: string): Component | undefined {
  if (!name) return undefined
  const exportName = legacyIconMap[name] ?? iconValueToExport(name)
  return iconComponents[exportName] ?? iconComponents.IconFile
}

function iconValueToExport(value: string): string {
  if (value.startsWith('Icon')) return value
  return `Icon${value
    .split(/[-_\s]+/)
    .filter(Boolean)
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join('')}`
}

function visibleChildren(node: Menu): Menu[] {
  return (node.children ?? []).filter((c) => c.visible !== 0 && c.status !== 0)
}

function onMenuItemClick(key: string): void {
  if (key && key.startsWith('/')) {
    void router.push(key)
    emit('navigate')
  }
}
</script>
