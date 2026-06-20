<template>
  <div class="app-page-loading">
    <DashboardLoadingOverlay :title="text" subtitle="正在准备页面数据" />
  </div>
</template>

<script setup lang="ts">
import { DashboardLoadingOverlay } from '@/components/loading'

withDefaults(
  defineProps<{
    text?: string
  }>(),
  { text: '数据加载中...' },
)
</script>

<style scoped>
.app-page-loading {
  position: fixed;
  inset: 0;
  z-index: 9990;
  background: #050914;
}
</style>
