<template>
  <div class="flex flex-col h-full min-h-0">
    <div class="layout-sider-logo shrink-0">
      <IconApps v-if="collapsed" class="text-[22px] text-[rgb(var(--primary-6))]" />
      <span v-else class="truncate px-2">Tower 后台</span>
    </div>
    <div class="flex-1 min-h-0 overflow-y-auto py-2">
      <ArcoProMenu :items="userStore.menus" :collapsed="collapsed" @navigate="emit('navigate')" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { IconApps } from '@arco-design/web-vue/es/icon'
import ArcoProMenu from './ArcoProMenu.vue'
import { useUserStore } from '@/store/user'

defineProps<{ collapsed: boolean }>()
const emit = defineEmits<{ navigate: [] }>()

const userStore = useUserStore()
</script>
