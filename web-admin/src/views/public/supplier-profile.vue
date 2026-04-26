<template>
  <div class="min-h-screen bg-slate-50 py-6 px-4">
    <div class="mx-auto max-w-3xl">
      <BaseCard>
        <template #header>
          <div class="flex items-center justify-between gap-2">
            <h2 class="m-0 text-xl font-semibold text-slate-800">供应商档案信息</h2>
            <BaseButton variant="ghost" size="sm" @click="goBack">返回</BaseButton>
          </div>
        </template>

        <div v-if="loading" class="py-10 text-center text-slate-500">加载中...</div>
        <div v-else-if="!supplier" class="py-10 text-center text-slate-400">未找到该供应商</div>
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div class="rounded bg-[var(--color-fill-1)] p-3"><span class="text-slate-500">供应商编码：</span>{{ supplier.supplier_code || '-' }}</div>
          <div class="rounded bg-[var(--color-fill-1)] p-3"><span class="text-slate-500">供应商名称：</span>{{ supplier.supplier_name || '-' }}</div>
          <div class="rounded bg-[var(--color-fill-1)] p-3"><span class="text-slate-500">联系人：</span>{{ supplier.contact_person || '-' }}</div>
          <div class="rounded bg-[var(--color-fill-1)] p-3"><span class="text-slate-500">联系电话：</span>{{ supplier.contact_phone || '-' }}</div>
          <div class="rounded bg-[var(--color-fill-1)] p-3"><span class="text-slate-500">联系邮箱：</span>{{ supplier.contact_email || '-' }}</div>
          <div class="rounded bg-[var(--color-fill-1)] p-3"><span class="text-slate-500">状态：</span>{{ supplier.status === 1 ? '启用' : '禁用' }}</div>
          <div class="md:col-span-2 rounded bg-[var(--color-fill-1)] p-3">
            <span class="text-slate-500">地址：</span>{{ supplier.supplier_address || '-' }}
          </div>
          <div class="md:col-span-2 rounded bg-[var(--color-fill-1)] p-3 whitespace-pre-wrap">
            <span class="text-slate-500">备注：</span>{{ supplier.remark || '-' }}
          </div>
        </div>
      </BaseCard>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BaseButton, BaseCard } from '@/components/base'
import { getPublicSupplier } from '@/api/supplier'
import type { Supplier } from '@/api/types'

const route = useRoute()
const router = useRouter()
const supplier = ref<Supplier | null>(null)
const loading = ref(true)

function goBack(): void {
  if (window.history.length > 1) router.back()
  else void router.push('/')
}

async function init(): Promise<void> {
  const id = Number(route.params.id)
  if (!Number.isFinite(id) || id <= 0) {
    loading.value = false
    return
  }
  try {
    supplier.value = await getPublicSupplier(id)
  } finally {
    loading.value = false
  }
}

void init()
</script>
