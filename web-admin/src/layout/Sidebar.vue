<template>
  <div class="flex flex-col h-full min-h-0">
    <div class="layout-sider-logo shrink-0">
      <span class="layout-sider-logo__mark">
        <img :src="logoUrl" alt="" />
      </span>
      <span v-if="!collapsed" class="layout-sider-logo__text">Tower 后台</span>
    </div>
    <div class="flex-1 min-h-0 overflow-y-auto py-2">
      <ArcoProMenu :items="userStore.menus" :collapsed="collapsed" @navigate="emit('navigate')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import ArcoProMenu from './ArcoProMenu.vue'
import { useUserStore } from '@/store/user'
import logoUrl from '@/assets/logo.svg'

defineProps<{ collapsed: boolean }>()
const emit = defineEmits<{ navigate: [] }>()

const userStore = useUserStore()
</script>
