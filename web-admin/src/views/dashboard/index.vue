<template>
  <div class="dashboard-fullscreen-shell">
    <StoreAnalyticsScreen fullscreen />
  </div>
</template>

<script setup lang="ts">
import { defineAsyncComponent, nextTick, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import AppPageLoading from '@/components/AppPageLoading.vue'
import DashboardChunkError from './DashboardChunkError.vue'

const StoreAnalyticsScreen = defineAsyncComponent({
  loader: () => import('./StoreAnalyticsScreen.vue'),
  loadingComponent: AppPageLoading,
  delay: 80,
  errorComponent: DashboardChunkError,
  onError(err, retry, fail, attempts) {
    console.error('[dashboard] chunk load error', err)
    if (attempts <= 2) {
      retry()
    } else {
      fail()
    }
  },
})

const route = useRoute()

function scrollToAnalytics(): void {
  if (route.query.section !== 'analytics') return
  const run = (): void => {
    document.getElementById('dash-analytics')?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
  void nextTick(() => {
    run()
    setTimeout(run, 400)
  })
}

onMounted(scrollToAnalytics)
watch(() => route.query.section, scrollToAnalytics)
</script>

<style scoped>
.dashboard-fullscreen-shell {
  width: 100vw;
  height: 100vh;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
  background: #050914;
}
</style>
