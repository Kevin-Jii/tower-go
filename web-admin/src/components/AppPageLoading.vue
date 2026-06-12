<template>
  <div class="app-page-loading" role="status" aria-live="polite">
    <MathCurveLoader :text="text" size="lg" />
  </div>
</template>

<script setup lang="ts">
import { MathCurveLoader } from '@/components/loading'

withDefaults(
  defineProps<{
    text?: string
  }>(),
  { text: '加载中…' },
)
</script>

<style scoped>
.app-page-loading {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 200px;
  padding: 32px 16px;
}
</style>
