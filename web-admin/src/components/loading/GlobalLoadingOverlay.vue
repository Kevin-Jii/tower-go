<template>
  <Teleport to="body">
    <DashboardLoadingOverlay
      v-if="loadingStore.visible"
      fixed
      :title="loadingStore.message"
      subtitle="正在准备页面数据"
    />
  </Teleport>
</template>

<script setup lang="ts">
import DashboardLoadingOverlay from './DashboardLoadingOverlay.vue'
import { useLoadingStore } from '@/store/loading'

const loadingStore = useLoadingStore()
</script>
