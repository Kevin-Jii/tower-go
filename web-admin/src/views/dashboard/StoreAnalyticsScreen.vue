<template>
  <section
    v-if="canStats"
    id="dash-analytics"
    ref="screenRoot"
    class="dash-screen rounded-2xl border border-cyan-500/25 shadow-xl"
  >
    <div
      class="dash-screen__bg absolute inset-0 pointer-events-none opacity-90"
      aria-hidden="true"
    />
    <div class="relative z-10 p-4 md:p-6 space-y-5">
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-3">
        <div>
          <h3 class="m-0 text-lg md:text-xl font-semibold tracking-wide text-cyan-100 flex items-center gap-2">
            <span class="inline-block w-2 h-2 rounded-full bg-cyan-400 shadow-[0_0_12px_#22d3ee]" />
            经营数据大屏
          </h3>
          <p class="m-0 mt-1 text-sm md:text-base text-cyan-100/95 leading-snug">销售趋势 · 渠道占比 · 经营雷达 · 库存概览</p>
        </div>
        <div class="flex flex-wrap items-center gap-2 dash-controls">
          <a-select
            v-model="period"
            class="dash-select w-40"
            :options="periodOptions"
            :allow-clear="false"
            :popup-container="screenRoot ?? undefined"
          />
          <BaseButton size="sm" class="dash-btn-period" @click="refreshDash">刷新周期</BaseButton>
        </div>
      </div>

      <p v-if="dashPending" class="m-0 text-sm text-cyan-100/90">正在加载周期概览…</p>
      <p v-else-if="dashError" class="m-0 text-sm text-rose-300">{{ dashError }} · 可点「刷新周期」重试</p>
      <div v-else-if="dash" class="grid grid-cols-2 md:grid-cols-3 xl:grid-cols-6 gap-3">
        <div v-for="k in kpiCards" :key="k.label" class="dash-kpi rounded-lg px-3 py-3.5 border border-cyan-400/35 bg-slate-950/75">
          <p class="m-0 text-sm font-medium text-cyan-100 tracking-wide">{{ k.label }}</p>
          <p class="m-0 mt-2 text-xl md:text-2xl font-semibold text-white tabular-nums drop-shadow-sm">{{ k.value }}</p>
        </div>
      </div>

      <div class="grid grid-cols-1 xl:grid-cols-12 gap-4 min-h-0">
        <div class="xl:col-span-8 rounded-xl border border-violet-500/20 bg-slate-950/50 p-2 min-h-[280px]">
          <div ref="lineRef" class="w-full h-[280px] md:h-[320px]" />
        </div>
        <div class="xl:col-span-4 rounded-xl border border-violet-500/20 bg-slate-950/50 p-2 min-h-[280px]">
          <div ref="pieRef" class="w-full h-[280px] md:h-[320px]" />
        </div>
        <div class="xl:col-span-5 rounded-xl border border-violet-500/20 bg-slate-950/50 p-2 min-h-[260px]">
          <div ref="radarRef" class="w-full h-[260px]" />
        </div>
        <div class="xl:col-span-7 rounded-xl border border-cyan-500/15 bg-slate-950/40 p-4">
          <div class="flex flex-wrap gap-2 items-center mb-3">
            <span class="text-base font-semibold text-cyan-50">经营总览区间</span>
            <BaseInput v-model="ovStart" type="date" class="dash-date w-40" />
            <span class="text-cyan-200 text-base font-medium px-0.5">—</span>
            <BaseInput v-model="ovEnd" type="date" class="dash-date w-40" />
            <a-select
              v-model="granularity"
              class="dash-select w-32"
              :options="granularityOptions"
              :allow-clear="false"
              :popup-container="screenRoot ?? undefined"
            />
            <BaseButton size="sm" class="dash-btn-apply" @click="loadHomeCharts">应用</BaseButton>
          </div>
          <div v-if="overview" class="grid grid-cols-2 md:grid-cols-4 gap-4 text-base mb-4">
            <div>
              <span class="block text-base font-semibold text-cyan-200 mb-1">销售收入</span>
              <strong class="text-xl text-white tabular-nums">{{ num(overview.sales_amount) }}</strong>
            </div>
            <div>
              <span class="block text-base font-semibold text-cyan-200 mb-1">毛利</span>
              <strong class="text-xl text-white tabular-nums">{{ num(overview.gross_profit_amount) }}</strong>
            </div>
            <div>
              <span class="block text-base font-semibold text-cyan-200 mb-1">净利</span>
              <strong class="text-xl text-emerald-200 tabular-nums">{{ num(overview.net_profit_amount) }}</strong>
            </div>
            <div>
              <span class="block text-base font-semibold text-cyan-200 mb-1">销售单数</span>
              <strong class="text-xl text-white tabular-nums">{{ overview.sales_order_count ?? 0 }}</strong>
            </div>
          </div>
          <div v-if="(overview?.categories?.length ?? 0) > 0" class="overflow-x-auto rounded-lg border border-cyan-400/25">
            <table class="w-full min-w-[480px] text-base border-collapse">
              <thead>
                <tr class="bg-slate-900/95 text-cyan-100">
                  <th class="text-left px-3 py-2.5 font-semibold">品类</th>
                  <th class="text-right px-3 py-2.5 font-semibold">入库金额</th>
                  <th class="text-right px-3 py-2.5 font-semibold">出库金额</th>
                  <th class="text-right px-3 py-2.5 font-semibold">净额</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="(c, i) in overview?.categories" :key="i" class="border-t border-cyan-500/20 text-slate-100">
                  <td class="px-3 py-2 font-medium">{{ c.category_name }}</td>
                  <td class="px-3 py-2 text-right tabular-nums">{{ num(c.in_amount) }}</td>
                  <td class="px-3 py-2 text-right tabular-nums">{{ num(c.out_amount) }}</td>
                  <td class="px-3 py-2 text-right tabular-nums font-medium">{{ num(c.net_amount) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
          <p v-else-if="homeLoading" class="m-0 text-base text-cyan-100">加载图表数据…</p>
        </div>
      </div>
    </div>
  </section>
  <div v-else class="flex flex-col items-center justify-center min-h-[320px] px-4">
    <BaseCard class="max-w-lg w-full">
      <template #header>
        <span class="font-semibold text-slate-800">无法显示经营数据大屏</span>
      </template>
      <p class="m-0 text-slate-600 text-sm leading-relaxed">
        当前账号缺少 <code class="text-indigo-700">statistics:dashboard</code> 权限，且非总部/超级管理员角色时不会展示大屏。
      </p>
      <p class="mt-3 mb-0 text-slate-500 text-sm leading-relaxed">
        请在「系统管理 → 角色」中为该角色勾选菜单「数据统计」，或联系管理员分配权限。
      </p>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import * as echarts from 'echarts'
import { BaseButton, BaseCard, BaseInput } from '@/components/base'
import { getHomeCharts, getStatisticsDashboard } from '@/api/statistics'
import type { HomeChartsStats } from '@/api/types'
import { useUserStore } from '@/store/user'
import { toast } from '@/feedback/toast'

const userStore = useUserStore()
const qc = useQueryClient()

/** 与后端一致：总部管理员 / 超级管理员走全量接口；其余角色依赖权限列表中的 statistics:dashboard */
const roleCode = computed(() => String(userStore.userInfo?.role?.code ?? ''))
const canStats = computed(() => {
  const code = roleCode.value
  if (code === 'super_admin' || code === 'admin') return true
  return userStore.permissions.includes('statistics:dashboard')
})

const period = ref('today')
const periodOptions = [
  { label: '今日', value: 'today' },
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' },
  { label: '本季', value: 'quarter' },
  { label: '本年', value: 'year' },
]

const granularityOptions = [
  { label: '按日', value: 'day' },
  { label: '按月', value: 'month' },
]

const screenRoot = ref<HTMLElement | null>(null)

const storeParam = computed(() => {
  const tid = Number(userStore.tenantId ?? 0)
  return tid > 0 ? { store_id: tid } : {}
})

const dashKey = computed(() => ['statistics', 'dashboard', period.value, userStore.tenantId] as const)
const {
  data: dash,
  isPending: dashPending,
  isError: dashIsError,
  error: dashQueryError,
  refetch: refetchDash,
} = useQuery({
  queryKey: dashKey,
  queryFn: () => getStatisticsDashboard({ period: period.value, ...storeParam.value }),
  enabled: () => canStats.value,
})

const dashError = computed(() => {
  if (!dashIsError.value) return ''
  const e = dashQueryError.value
  return e instanceof Error ? e.message : '周期概览加载失败'
})

function refreshDash(): void {
  void refetchDash()
  void qc.invalidateQueries({ queryKey: ['statistics', 'dashboard'] })
}

watch(period, () => {
  void qc.invalidateQueries({ queryKey: ['statistics', 'dashboard'] })
})

const kpiCards = computed(() => {
  const d = dash.value
  if (!d) return []
  return [
    { label: '库存 SKU', value: String(d.inventory.total_products ?? 0) },
    { label: '库存数量', value: String(d.inventory.total_quantity ?? 0) },
    { label: '今日入库', value: String(d.inventory.today_in ?? 0) },
    { label: '今日出库', value: String(d.inventory.today_out ?? 0) },
    { label: d.sales.period_label ? `${d.sales.period_label}销售` : '周期销售', value: num(d.sales.total_amount) },
    { label: '今日销售', value: num(d.sales.today_amount) },
  ]
})

function monthRange(): { start: string; end: string } {
  const t = new Date()
  const y = t.getFullYear()
  const m = String(t.getMonth() + 1).padStart(2, '0')
  const day = String(t.getDate()).padStart(2, '0')
  return { start: `${y}-${m}-01`, end: `${y}-${m}-${day}` }
}

const r = monthRange()
const ovStart = ref(r.start)
const ovEnd = ref(r.end)
const granularity = ref<'day' | 'month'>('day')

const homeCharts = ref<HomeChartsStats | null>(null)
const homeLoading = ref(false)
const overview = computed(() => homeCharts.value?.overview ?? null)

const lineRef = ref<HTMLElement | null>(null)
const pieRef = ref<HTMLElement | null>(null)
const radarRef = ref<HTMLElement | null>(null)
let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null
let radarChart: echarts.ECharts | null = null

const axisText = '#e2e8f0'
const splitLine = 'rgba(148, 163, 184, 0.22)'
const axisFont = 13

function num(v: number | undefined): string {
  if (v == null || Number.isNaN(v)) return '-'
  return Number(v).toFixed(2)
}

function disposeCharts(): void {
  lineChart?.dispose()
  pieChart?.dispose()
  radarChart?.dispose()
  lineChart = pieChart = radarChart = null
}

function ensureCharts(): void {
  if (!lineRef.value || !pieRef.value || !radarRef.value) return
  if (!lineChart) lineChart = echarts.init(lineRef.value)
  if (!pieChart) pieChart = echarts.init(pieRef.value)
  if (!radarChart) radarChart = echarts.init(radarRef.value)
}

function applyChartOptions(hc: HomeChartsStats): void {
  ensureCharts()
  if (!lineChart || !pieChart || !radarChart) return

  const line = hc.line ?? []
  lineChart.setOption({
    backgroundColor: 'transparent',
    tooltip: { trigger: 'axis', backgroundColor: 'rgba(15, 23, 42, 0.92)', borderColor: '#22d3ee', textStyle: { color: '#e2e8f0' } },
    grid: { left: 56, right: 20, top: 36, bottom: line.length > 10 ? 52 : 44 },
    xAxis: {
      type: 'category',
      data: line.map((x) => x.date),
      axisLine: { lineStyle: { color: 'rgba(226,232,240,0.5)' } },
      axisLabel: { color: axisText, fontSize: axisFont, fontWeight: 500, rotate: line.length > 10 ? 28 : 0 },
    },
    yAxis: {
      type: 'value',
      splitLine: { lineStyle: { color: splitLine } },
      axisLabel: { color: axisText, fontSize: axisFont, fontWeight: 500 },
    },
    series: [
      {
        name: '销售额',
        type: 'line',
        smooth: true,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: { width: 2, color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [{ offset: 0, color: '#22d3ee' }, { offset: 1, color: '#a78bfa' }]) },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(34, 211, 238, 0.35)' },
            { offset: 1, color: 'rgba(15, 23, 42, 0.02)' },
          ]),
        },
        data: line.map((x) => x.amount),
      },
    ],
  })

  const pie = hc.pie ?? []
  pieChart.setOption({
    backgroundColor: 'transparent',
    tooltip: { trigger: 'item', backgroundColor: 'rgba(15, 23, 42, 0.92)', borderColor: '#a78bfa', textStyle: { color: '#e2e8f0' } },
    legend: { bottom: 4, textStyle: { color: axisText, fontSize: 13, fontWeight: 600 } },
    series: [
      {
        name: '渠道',
        type: 'pie',
        radius: ['42%', '68%'],
        center: ['50%', '46%'],
        itemStyle: {
          borderRadius: 6,
          borderColor: '#0f172a',
          borderWidth: 2,
        },
        label: { color: '#f8fafc', fontSize: 13, fontWeight: 600 },
        data: pie.map((p, i) => ({
          name: p.channel_name || p.channel || `渠道${i}`,
          value: p.amount,
        })),
      },
    ],
    color: ['#22d3ee', '#a78bfa', '#34d399', '#fbbf24', '#fb7185', '#38bdf8'],
  })

  const radar = hc.radar ?? []
  const indicators = radar.map((x) => ({ name: x.name, max: Math.max(Math.abs(Number(x.value)) * 1.35, 1) }))
  radarChart.setOption({
    backgroundColor: 'transparent',
    tooltip: { backgroundColor: 'rgba(15, 23, 42, 0.92)', borderColor: '#22d3ee', textStyle: { color: '#e2e8f0' } },
    radar: {
      indicator: indicators.length ? indicators : [{ name: '暂无', max: 1 }],
      splitLine: { lineStyle: { color: splitLine } },
      splitArea: { show: true, areaStyle: { color: ['rgba(34,211,238,0.05)', 'rgba(167,139,250,0.06)'] } },
      axisName: { color: '#e2e8f0', fontSize: 13, fontWeight: 600, padding: [2, 4] },
    },
    series: [
      {
        type: 'radar',
        data: [
          {
            value: radar.length ? radar.map((x) => x.value) : [0],
            name: '经营指标',
            areaStyle: { color: 'rgba(34, 211, 238, 0.25)' },
            lineStyle: { color: '#22d3ee' },
            itemStyle: { color: '#a78bfa' },
          },
        ],
      },
    ],
  })
}

async function loadHomeCharts(): Promise<void> {
  if (!canStats.value) return
  if (!ovStart.value || !ovEnd.value) {
    toast.warning('请选择起止日期')
    return
  }
  homeLoading.value = true
  try {
    const hc = await getHomeCharts({
      start_date: ovStart.value,
      end_date: ovEnd.value,
      granularity: granularity.value,
      ...storeParam.value,
    })
    homeCharts.value = hc
    await paintChartsWhenReady(hc)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载失败')
    homeCharts.value = null
  } finally {
    homeLoading.value = false
  }
}

/** 等图表容器 ref 挂载后再 init/setOption（避免 onMounted 首帧 ref 仍为 null 导致永远不画图） */
async function paintChartsWhenReady(hc: HomeChartsStats): Promise<void> {
  for (let i = 0; i < 12; i++) {
    await nextTick()
    if (lineRef.value && pieRef.value && radarRef.value) {
      applyChartOptions(hc)
      requestAnimationFrame(() => {
        lineChart?.resize()
        pieChart?.resize()
        radarChart?.resize()
      })
      return
    }
    await new Promise<void>((r) => setTimeout(r, 32))
  }
  applyChartOptions(hc)
}

function onWinResize(): void {
  lineChart?.resize()
  pieChart?.resize()
  radarChart?.resize()
}

watch(
  () => homeCharts.value,
  (hc) => {
    if (!hc) return
    void paintChartsWhenReady(hc)
  },
)

onMounted(() => {
  window.addEventListener('resize', onWinResize)
  void (async () => {
    await nextTick()
    await nextTick()
    if (canStats.value) await loadHomeCharts()
  })()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onWinResize)
  disposeCharts()
})
</script>

<style scoped>
.dash-screen {
  position: relative;
  background: linear-gradient(145deg, #070b14 0%, #0b1220 45%, #0a1628 100%);
}
.dash-screen__bg {
  background-image:
    radial-gradient(ellipse 80% 50% at 50% -20%, rgba(34, 211, 238, 0.18), transparent),
    linear-gradient(rgba(34, 211, 238, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(34, 211, 238, 0.04) 1px, transparent 1px);
  background-size: 100% 100%, 24px 24px, 24px 24px;
}
.dash-kpi {
  box-shadow: 0 0 24px rgba(34, 211, 238, 0.06);
}

/* 大屏内 Arco 控件：高对比深色，避免白底下拉与灰字 */
.dash-screen :deep(.arco-select-view-single) {
  min-height: 36px;
  background: rgba(15, 23, 42, 0.95) !important;
  border: 1px solid rgba(34, 211, 238, 0.55) !important;
  color: #f0fdfa !important;
}
.dash-screen :deep(.arco-select-view-value),
.dash-screen :deep(.arco-select-view-suffix) {
  color: #f0fdfa !important;
  font-size: 15px;
  font-weight: 600;
}
.dash-screen :deep(.arco-select-dropdown) {
  background: #0f172a !important;
  border: 1px solid rgba(34, 211, 238, 0.5) !important;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.55);
}
.dash-screen :deep(.arco-select-option) {
  color: #f1f5f9 !important;
  font-size: 15px;
  font-weight: 500;
  line-height: 1.5;
}
.dash-screen :deep(.arco-select-option-selected) {
  background: rgba(34, 211, 238, 0.22) !important;
  color: #fff !important;
  font-weight: 700;
}
.dash-screen :deep(.arco-select-option:not(.arco-select-option-disabled):hover) {
  background: rgba(34, 211, 238, 0.14) !important;
}

.dash-screen :deep(.dash-date .arco-input-wrapper) {
  background: rgba(15, 23, 42, 0.95) !important;
  border: 1px solid rgba(34, 211, 238, 0.5) !important;
  min-height: 36px;
}
.dash-screen :deep(.dash-date .arco-input) {
  color: #f8fafc !important;
  font-size: 15px;
  font-weight: 600;
}

.dash-btn-period {
  background: rgba(8, 145, 178, 0.45) !important;
  border: 1px solid rgba(34, 211, 238, 0.75) !important;
  color: #ecfeff !important;
  font-weight: 700 !important;
  font-size: 14px !important;
}
.dash-btn-apply {
  background: rgba(109, 40, 217, 0.55) !important;
  border: 1px solid rgba(167, 139, 250, 0.75) !important;
  color: #faf5ff !important;
  font-weight: 700 !important;
  font-size: 14px !important;
}
</style>
