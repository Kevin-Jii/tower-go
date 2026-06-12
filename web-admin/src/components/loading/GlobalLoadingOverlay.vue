<template>
  <Teleport to="body">
    <MathCurveLoader
      v-if="loadingStore.visible"
      overlay
      size="lg"
      :text="loadingStore.message"
    />
  </Teleport>
</template>

<script setup lang="ts">
import MathCurveLoader from './MathCurveLoader.vue'
import { useLoadingStore } from '@/store/loading'

const loadingStore = useLoadingStore()
</script>
