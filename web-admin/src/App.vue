<template>
  <router-view />
  <FeedbackHost />
  <GlobalLoadingOverlay />
</template>

<script setup lang="ts">
import FeedbackHost from '@/feedback/FeedbackHost.vue'
import { GlobalLoadingOverlay } from '@/components/loading'
</script>
