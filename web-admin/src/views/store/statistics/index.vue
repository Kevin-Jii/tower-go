<template>
  <div class="flex flex-col items-center justify-center min-h-[38vh] gap-4 px-4">
    <p class="m-0 text-center text-slate-600 text-sm md:text-base max-w-md">
      「数据统计」已合并到工作台的数据大屏。正在为您跳转…
    </p>
    <BaseButton variant="primary" @click="go">立即前往工作台</BaseButton>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { BaseButton } from '@/components/base'

const router = useRouter()

function go(): void {
  void router.replace({ name: 'Dashboard', query: { section: 'analytics' } })
}

onMounted(() => {
  go()
})
</script>
