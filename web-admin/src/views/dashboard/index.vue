<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <BaseCard>
      <template #header>
        <span class="font-semibold text-slate-800">欢迎</span>
      </template>
      <p class="text-slate-600 text-sm md:text-base leading-relaxed m-0">
        {{ userStore.userInfo?.nickname || userStore.userInfo?.username }}，当前角色：
        <strong class="text-indigo-700">{{ userStore.userInfo?.role?.name || '-' }}</strong>
      </p>
    </BaseCard>
    <BaseCard>
      <template #header>
        <span class="font-semibold text-slate-800">字典标签示例</span>
      </template>
      <p class="text-sm text-slate-500 mb-3 m-0">组件 &lt;DictTag type="inbound_status" :value="v" /&gt;</p>
      <div class="flex flex-wrap gap-2 items-center">
        <DictTag type="inbound_status" :value="0" />
        <DictTag type="inbound_status" :value="1" />
        <DictTag type="inbound_status" :value="2" />
        <span class="text-xs text-slate-400">（无字典数据时回退显示原始值）</span>
      </div>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { BaseCard } from '@/components/base'
import DictTag from '@/components/DictTag.vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
</script>
