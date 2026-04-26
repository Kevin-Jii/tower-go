<template>
  <div
    class="flex-1 min-h-0 overflow-y-auto flex flex-col items-center justify-center gap-4 p-8 text-center"
  >
    <IconFaceFrownFill class="text-[56px]" style="color: rgb(var(--warning-6))" />
    <h1 class="text-2xl font-semibold text-[var(--color-text-1)] m-0">404</h1>
    <p class="text-[var(--color-text-3)] m-0">页面不存在</p>
    <a-button type="primary" @click="$router.push({ name: 'Dashboard' })">返回工作台</a-button>
  </div>
</template>

<script setup lang="ts">
import { IconFaceFrownFill } from '@arco-design/web-vue/es/icon'
</script>
