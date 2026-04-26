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
import {
  IconApps,
  IconBook,
  IconFile,
  IconImage,
  IconMenu,
  IconSettings,
  IconUser,
  IconUserGroup,
} from '@arco-design/web-vue/es/icon'
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

const iconMap: Record<string, Component> = {
  setting: IconSettings,
  user: IconUser,
  usergroup: IconUserGroup,
  'menu-fold': IconMenu,
  read: IconBook,
  picture: IconImage,
  document: IconFile,
  apps: IconApps,
}

function iconCmp(name?: string): Component | undefined {
  if (!name) return undefined
  return iconMap[name] ?? IconFile
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
