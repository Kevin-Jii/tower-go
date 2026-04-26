<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">数据统计</h2>
      <div class="flex flex-col sm:flex-row flex-wrap gap-2 w-full md:w-auto items-stretch sm:items-center">
        <BaseSelect
          v-model="period"
          class="w-full sm:w-40"
          :options="[
            { label: '今日', value: 'today' },
            { label: '本周', value: 'week' },
            { label: '本月', value: 'month' },
            { label: '本季', value: 'quarter' },
            { label: '本年', value: 'year' },
          ]"
        />
        <BaseButton variant="primary" @click="refreshDash">刷新概览</BaseButton>
      </div>
    </div>

    <div v-if="dash" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">库存 SKU 数</span></template>
        <p class="text-2xl font-semibold m-0">{{ dash.inventory.total_products }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">库存总数量</span></template>
        <p class="text-2xl font-semibold m-0">{{ dash.inventory.total_quantity }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">出入库记录数</span></template>
        <p class="text-2xl font-semibold m-0">{{ dash.inventory.total_records }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">今日入库量</span></template>
        <p class="text-2xl font-semibold m-0">{{ dash.inventory.today_in }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">今日出库量</span></template>
        <p class="text-2xl font-semibold m-0">{{ dash.inventory.today_out }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">{{ dash.sales.period_label || '销售' }}销售额</span></template>
        <p class="text-2xl font-semibold text-indigo-700 m-0">{{ dash.sales.total_amount?.toFixed?.(2) ?? dash.sales.total_amount }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">今日销售额</span></template>
        <p class="text-2xl font-semibold m-0">{{ dash.sales.today_amount?.toFixed?.(2) ?? dash.sales.today_amount }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">本月销售额</span></template>
        <p class="text-2xl font-semibold m-0">{{ dash.sales.month_amount?.toFixed?.(2) ?? dash.sales.month_amount }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">订单数 / 客单价</span></template>
        <p class="text-lg font-semibold m-0">
          {{ dash.sales.total_orders }} 单 · 均 {{ dash.sales.avg_amount?.toFixed?.(2) ?? dash.sales.avg_amount }}
        </p>
      </BaseCard>
    </div>
    <p v-else-if="dashLoading" class="text-slate-500 text-sm m-0">加载中…</p>
    <p v-else class="text-slate-500 text-sm m-0">暂无数据</p>

    <BaseCard>
      <template #header>
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-2 w-full">
          <span class="font-semibold text-slate-800">经营总览（按日期）</span>
          <div class="flex flex-wrap gap-2 items-center">
            <BaseInput v-model="ovStart" type="date" class="w-36" />
            <span class="text-slate-400">至</span>
            <BaseInput v-model="ovEnd" type="date" class="w-36" />
            <BaseButton size="sm" variant="primary" @click="loadOverview">查询</BaseButton>
          </div>
        </div>
      </template>
      <div v-if="overview" class="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
        <div><span class="text-slate-500">销售收入</span><br /><strong>{{ num(overview.sales_amount) }}</strong></div>
        <div><span class="text-slate-500">毛利</span><br /><strong>{{ num(overview.gross_profit_amount) }}</strong></div>
        <div><span class="text-slate-500">净利</span><br /><strong>{{ num(overview.net_profit_amount) }}</strong></div>
        <div><span class="text-slate-500">销售单数</span><br /><strong>{{ overview.sales_order_count ?? 0 }}</strong></div>
      </div>
      <div v-if="(overview?.categories?.length ?? 0) > 0" class="mt-4 min-w-0 overflow-x-auto">
        <BaseTable
          :columns="catCols"
          :data="(overview!.categories as unknown) as Record<string, unknown>[]"
          min-width="520px"
        />
      </div>
      <p v-else-if="overviewLoading" class="text-slate-500 text-sm m-0">加载经营数据…</p>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { BaseButton, BaseCard, BaseInput, BaseSelect, BaseTable } from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import { getBusinessOverview, getStatisticsDashboard } from '@/api/statistics'
import type { BusinessOverviewStats } from '@/api/types'
import { toast } from '@/feedback/toast'

const qc = useQueryClient()
const period = ref('today')

const dashKey = computed(() => ['statistics', 'dashboard', period.value] as const)
const { data: dash, isLoading: dashLoading } = useQuery({
  queryKey: dashKey,
  queryFn: () => getStatisticsDashboard({ period: period.value }),
})

function refreshDash(): void {
  void qc.invalidateQueries({ queryKey: ['statistics', 'dashboard'] })
}

watch(period, () => {
  void qc.invalidateQueries({ queryKey: ['statistics', 'dashboard'] })
})

function monthRange(): { start: string; end: string } {
  const t = new Date()
  const y = t.getFullYear()
  const m = String(t.getMonth() + 1).padStart(2, '0')
  const d = String(t.getDate()).padStart(2, '0')
  return { start: `${y}-${m}-01`, end: `${y}-${m}-${d}` }
}

const r = monthRange()
const ovStart = ref(r.start)
const ovEnd = ref(r.end)

const overview = ref<BusinessOverviewStats | null>(null)
const overviewLoading = ref(false)

async function loadOverview(): Promise<void> {
  if (!ovStart.value || !ovEnd.value) {
    toast.warning('请选择起止日期')
    return
  }
  overviewLoading.value = true
  try {
    overview.value = await getBusinessOverview({
      start_date: ovStart.value,
      end_date: ovEnd.value,
    })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载失败')
    overview.value = null
  } finally {
    overviewLoading.value = false
  }
}

void loadOverview()

const catCols: BaseTableColumn[] = [
  { key: 'category_name', label: '品类', prop: 'category_name', minWidth: '120px' },
  { key: 'in_amount', label: '入库金额', prop: 'in_amount', width: '110px' },
  { key: 'out_amount', label: '出库金额', prop: 'out_amount', width: '110px' },
  { key: 'net_amount', label: '净额', prop: 'net_amount', width: '110px' },
]

function num(v: number | undefined): string {
  if (v == null || Number.isNaN(v)) return '-'
  return v.toFixed(2)
}
</script>
