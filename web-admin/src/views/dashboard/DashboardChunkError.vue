<template>
  <div class="flex flex-col items-center justify-center min-h-[280px] gap-4 px-4 text-center">
    <p class="m-0 text-slate-700 text-sm max-w-md">
      经营数据模块加载失败（网络中断或静态资源缓存过期）。请点击下方按钮重试。
    </p>
    <BaseButton variant="primary" @click="reload">刷新页面</BaseButton>
  </div>
</template>

<script setup lang="ts">
import { BaseButton } from '@/components/base'

function reload(): void {
  window.location.reload()
}
</script>
