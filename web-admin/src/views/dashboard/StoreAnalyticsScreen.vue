<template>
  <section id="dash-analytics" ref="screenRoot"
    class="dash-screen box-border w-full min-w-0 max-w-full overflow-hidden border border-cyan-500/25 shadow-xl"
    :class="{ 'dash-screen--fullscreen': fullscreen }">
    <div class="dash-screen__bg absolute inset-0 pointer-events-none opacity-90" aria-hidden="true" />
    <div class="dash-scanline absolute inset-x-0 top-0 pointer-events-none" aria-hidden="true" />
    <div class="dash-screen__content relative z-10 box-border w-full min-w-0">
      <div class="dash-screen__topbar">
        <div>
          <h3 class="dash-screen__title">
            <span class="dash-screen__pulse" />
            经营数据大屏
          </h3>
        </div>
        <div class="flex flex-wrap items-center gap-2 dash-controls">
          <select v-model="period" class="dash-native-control dash-native-select" aria-label="统计周期">
            <option v-for="option in periodOptions" :key="option.value" :value="option.value">
              {{ option.label }}
            </option>
          </select>
          <BaseButton size="sm" class="dash-btn-period" @click="refreshDash">刷新周期</BaseButton>
          <BaseButton size="sm" class="dash-btn-admin" @click="goAdmin">后台入口</BaseButton>
        </div>
      </div>

      <p v-if="dashError" class="m-0 text-sm text-rose-300">{{ dashError }} · 可点「刷新周期」重试</p>

      <div class="dash-chart-grid min-h-0 min-w-0">
        <div class="dash-panel dash-panel--overview min-w-0 rounded-xl border border-cyan-500/15 bg-slate-950/40 p-4">
          <div v-if="overview || dash" class="dash-overview-focus">
            <div class="dash-primary-stat">
              <strong class="tabular-nums"
                :class="Number(overview?.net_profit_amount ?? 0) >= 0 ? 'tone-good' : 'tone-bad'">
                <CountUpNumber :value="overview?.net_profit_amount" :decimals="2" />
              </strong>
              <span>净利</span>
            </div>
            <div class="dash-overview-metrics">
              <div v-for="(metric, idx) in overviewCards" :key="metric.label" class="dash-overview-card"
                :style="{ '--delay': `${idx * 45}ms` }">
                <span>{{ metric.label }}</span>
                <strong class="tabular-nums" :class="metric.tone">
                  <CountUpNumber :value="metric.value" :decimals="metric.decimals" :suffix="metric.suffix || ''" />
                </strong>
              </div>
            </div>
          </div>
        </div>

        <div class="dash-panel dash-panel--pie min-w-0 rounded-xl border border-violet-500/20 bg-slate-950/50 p-2">
          <div ref="pieRef" class="dash-chart w-full min-w-0" />
        </div>

        <div class="dash-panel dash-panel--radar min-w-0 rounded-xl border border-violet-500/20 bg-slate-950/50 p-3">
          <div class="dash-radar-layout">
            <div ref="radarRef" class="dash-chart dash-radar-chart w-full min-w-0" />
          </div>
        </div>

        <div class="dash-panel dash-panel--member min-w-0 rounded-xl border border-cyan-500/15 bg-slate-950/40 p-4">
          <div class="dash-member-layout">
            <div class="dash-section-title">会员消费排行</div>
            <div ref="memberRankRef" class="dash-chart dash-member-chart w-full min-w-0" />
          </div>
        </div>

        <div class="dash-panel dash-panel--flow min-w-0 rounded-xl border border-cyan-500/15 bg-slate-950/40 p-4">
          <div class="dash-flow-layout">
            <div class="dash-category-flow">
              <div class="dash-section-title">品类流动</div>
              <div v-if="topCategories.length" class="dash-flow-list">
                <div v-for="item in topCategories" :key="item.category_name" class="dash-flow-row">
                  <span class="dash-flow-name">{{ item.category_name }}</span>
                  <span class="dash-flow-track">
                    <i :style="{ width: `${item.percent}%` }" />
                  </span>
                  <strong class="tabular-nums">
                    <CountUpNumber :value="item.amount" :decimals="2" />
                  </strong>
                </div>
              </div>
              <p v-else class="m-0 text-sm text-cyan-100/70">暂无品类流动数据</p>
            </div>
          </div>
        </div>

        <div class="dash-panel dash-panel--line min-w-0 rounded-xl border border-violet-500/20 bg-slate-950/50 p-2">
          <div ref="lineRef" class="dash-chart w-full min-w-0" />
        </div>
      </div>
    </div>

    <div v-show="screenLoading" class="dash-loading-overlay" aria-live="polite" aria-busy="true">
      <div class="dash-loading-core">
        <span class="dash-loading-orbit" />
        <span class="dash-loading-orbit dash-loading-orbit--second" />
        <span class="dash-loading-dot" />
        <strong>数据同步中</strong>
        <small>正在刷新经营大屏</small>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import * as echarts from 'echarts'
import { useRouter } from 'vue-router'
import { BaseButton } from '@/components/base'
import CountUpNumber from '@/components/CountUpNumber.vue'
import { getHomeCharts, getStatisticsDashboard } from '@/api/statistics'
import type { HomeChartsStats } from '@/api/types'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const router = useRouter()

defineProps<{
  fullscreen?: boolean
}>()

const period = ref('week')
const periodOptions = [
  { label: '今日', value: 'today' },
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' },
  { label: '本季', value: 'quarter' },
  { label: '本年', value: 'year' },
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
  isFetching: dashFetching,
  refetch: refetchDash,
} = useQuery({
  queryKey: dashKey,
  queryFn: () => getStatisticsDashboard({ period: period.value, ...storeParam.value }),
})

const dashError = computed(() => {
  if (!dashIsError.value) return ''
  const e = dashQueryError.value
  return e instanceof Error ? e.message : '周期概览加载失败'
})

async function refreshDash(): Promise<void> {
  await Promise.all([refetchDash(), loadHomeCharts()])
}

function goAdmin(): void {
  void router.push('/admin')
}

watch(period, () => {
  void loadHomeCharts()
})

function dateText(t: Date): string {
  const y = t.getFullYear()
  const m = String(t.getMonth() + 1).padStart(2, '0')
  const d = String(t.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function periodRange(value: string): { start: string; end: string } {
  const t = new Date()
  const end = dateText(t)
  if (value === 'today') return { start: end, end }
  if (value === 'week') {
    const weekday = t.getDay() || 7
    const start = new Date(t)
    start.setDate(t.getDate() - weekday + 1)
    return { start: dateText(start), end }
  }
  if (value === 'quarter') {
    const quarter = Math.floor(t.getMonth() / 3)
    return { start: dateText(new Date(t.getFullYear(), quarter * 3, 1)), end }
  }
  if (value === 'year') return { start: `${t.getFullYear()}-01-01`, end }
  return { start: dateText(new Date(t.getFullYear(), t.getMonth(), 1)), end }
}

function rangeGranularity(range: { start: string; end: string }): 'day' | 'month' {
  const start = new Date(`${range.start}T00:00:00`)
  const end = new Date(`${range.end}T00:00:00`)
  const days = Math.max(1, Math.round((end.getTime() - start.getTime()) / 86_400_000) + 1)
  return days > 92 ? 'month' : 'day'
}

const activeRange = computed(() => periodRange(period.value))

const homeCharts = ref<HomeChartsStats | null>(null)
const homeLoading = ref(false)
const overview = computed(() => homeCharts.value?.overview ?? null)
const screenLoading = computed(() => dashPending.value || dashFetching.value || homeLoading.value)
const overviewCards = computed(() => {
  const o = overview.value
  const d = dash.value
  return [
    { label: '销售收入', value: Number(o?.sales_amount ?? d?.sales.total_amount ?? 0), decimals: 2, tone: '' },
    { label: '毛利', value: Number(o?.gross_profit_amount ?? 0), decimals: 2, tone: '' },
    { label: '消耗品金额', value: Number(o?.consumable_amount ?? 0), decimals: 2, tone: '' },
    { label: 'B2B供货收入', value: Number(o?.b2b_supply_amount ?? 0), decimals: 2, tone: '' },
    { label: 'B2B订单数', value: Number(o?.b2b_supply_order_count ?? 0), decimals: 0, tone: '' },
    { label: '返罐押金', value: Number(o?.return_deposit_amount ?? 0), decimals: 2, tone: '' },
    { label: '物流费用', value: Number(o?.return_logistics_fee ?? 0), decimals: 2, tone: '' },
    { label: '跑腿费用', value: Number(o?.errand_fee_amount ?? 0), decimals: 2, tone: '' },
    { label: '抹零金额', value: Number(o?.round_amount ?? 0), decimals: 2, tone: '' },
    { label: '赠酒成本', value: Number(o?.gift_wine_cost_amount ?? 0), decimals: 2, tone: '' },
    { label: '门店支出', value: Number(o?.store_expense_amount ?? 0), decimals: 2, tone: '' },
    { label: '推广费用', value: Number(o?.takeout_promotion_amount ?? 0), decimals: 2, tone: '' },
    { label: '推广ROI', value: Number(o?.takeout_promotion_roi ?? 0), decimals: 2, tone: Number(o?.takeout_promotion_roi ?? 0) >= 1 ? 'tone-good' : 'tone-bad' },
    { label: '报损金额', value: Number(o?.inventory_loss_amount ?? 0), decimals: 2, tone: '' },
    { label: '自用金额', value: Number(o?.inventory_self_use_amount ?? 0), decimals: 2, tone: '' },
    { label: '销售单数', value: Number(o?.sales_order_count ?? d?.sales.total_orders ?? 0), decimals: 0, tone: '' },
    { label: '库存数量', value: Number(d?.inventory.total_quantity ?? 0), decimals: 0, tone: '' },
    { label: '今日销售', value: Number(d?.sales.today_amount ?? 0), decimals: 2, tone: '' },
    { label: '今日入库', value: Number(d?.inventory.today_in ?? 0), decimals: 0, tone: '' },
    { label: '毛利率', value: grossMargin.value, decimals: 1, suffix: '%', tone: grossMargin.value >= 0 ? 'tone-good' : 'tone-bad' },
    { label: '客单价', value: avgOrderAmount.value, decimals: 2, tone: '' },
    { label: '入库单数', value: Number(o?.inventory_in_count ?? 0), decimals: 0, tone: '' },
    { label: '出库单数', value: Number(o?.inventory_out_count ?? 0), decimals: 0, tone: '' },
    { label: '出库金额', value: Number(o?.outbound_amount ?? 0), decimals: 2, tone: '' },
    { label: '其他支出', value: Number(o?.other_expense_amount ?? 0), decimals: 2, tone: '' },
  ]
})
const grossMargin = computed(() => {
  const sales = Number(overview.value?.sales_amount ?? 0)
  const gross = Number(overview.value?.gross_profit_amount ?? 0)
  return sales > 0 ? (gross / sales) * 100 : 0
})
const avgOrderAmount = computed(() => {
  const sales = Number(overview.value?.sales_amount ?? 0)
  const orders = Number(overview.value?.sales_order_count ?? 0)
  return orders > 0 ? sales / orders : 0
})
const topCategories = computed(() => {
  const list = [...(overview.value?.categories ?? [])]
    .map((item) => ({
      ...item,
      amount: Math.max(Number(item.in_amount ?? 0), Number(item.out_amount ?? 0), Math.abs(Number(item.net_amount ?? 0))),
    }))
    .sort((a, b) => b.amount - a.amount)
    .slice(0, 4)
  const max = Math.max(...list.map((x) => x.amount), 1)
  return list.map((item) => ({ ...item, percent: Math.max(8, Math.round((item.amount / max) * 100)) }))
})
function formatAmount(value: number): string {
  return Number(value || 0).toLocaleString('zh-CN', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  })
}

const lineRef = ref<HTMLElement | null>(null)
const pieRef = ref<HTMLElement | null>(null)
const radarRef = ref<HTMLElement | null>(null)
const memberRankRef = ref<HTMLElement | null>(null)
let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null
let radarChart: echarts.ECharts | null = null
let memberRankChart: echarts.ECharts | null = null
let resizeObserver: ResizeObserver | null = null

const axisText = '#e2e8f0'
const splitLine = 'rgba(148, 163, 184, 0.22)'
const axisFont = 13

function disposeCharts(): void {
  lineChart?.dispose()
  pieChart?.dispose()
  radarChart?.dispose()
  memberRankChart?.dispose()
  lineChart = pieChart = radarChart = memberRankChart = null
}

function ensureCharts(): void {
  if (!lineRef.value || !pieRef.value || !radarRef.value || !memberRankRef.value) return
  if (!lineChart) lineChart = echarts.init(lineRef.value)
  if (!pieChart) pieChart = echarts.init(pieRef.value)
  if (!radarChart) radarChart = echarts.init(radarRef.value)
  if (!memberRankChart) memberRankChart = echarts.init(memberRankRef.value)
}

function applyChartOptions(hc: HomeChartsStats): void {
  ensureCharts()
  if (!lineChart || !pieChart || !radarChart || !memberRankChart) return

  const line = hc.line ?? []
  lineChart.setOption({
    backgroundColor: 'transparent',
    animationDuration: 900,
    animationEasing: 'cubicOut',
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
    animationDuration: 900,
    animationEasing: 'quarticOut',
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
    animationDuration: 900,
    animationEasing: 'cubicOut',
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

  const rank = hc.overview?.member_consumption_rank ?? []
  const rankLabels = rank.map((x) => x.member_name || x.member_phone || '未知会员')
  memberRankChart.setOption({
    backgroundColor: 'transparent',
    animationDuration: 900,
    animationEasing: 'cubicOut',
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      backgroundColor: 'rgba(15, 23, 42, 0.92)',
      borderColor: '#22d3ee',
      textStyle: { color: '#e2e8f0' },
      valueFormatter: (value: unknown) => `¥${formatAmount(Number(value))}`,
    },
    grid: { left: 92, right: 74, top: 18, bottom: 16 },
    xAxis: {
      type: 'value',
      show: false,
      max: (value: { max: number }) => Math.max(value.max * 1.18, 1),
    },
    yAxis: {
      type: 'category',
      inverse: true,
      data: rankLabels.length ? rankLabels : ['暂无'],
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: {
        color: axisText,
        fontSize: 13,
        fontWeight: 800,
        overflow: 'truncate',
        width: 82,
      },
    },
    series: [
      {
        name: '消费金额',
        type: 'bar',
        barWidth: 12,
        barGap: '45%',
        data: rank.length ? rank.map((x) => Number(x.amount || 0)) : [0],
        backgroundStyle: {
          color: 'rgba(15, 23, 42, 0.82)',
          borderRadius: 999,
        },
        showBackground: true,
        label: {
          show: true,
          position: 'right',
          color: '#a7f3d0',
          fontSize: 12,
          fontWeight: 900,
          formatter: ({ value }: { value: number | string }) => `¥${formatAmount(Number(value))}`,
        },
        itemStyle: {
          borderRadius: 999,
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: '#22d3ee' },
            { offset: 1, color: '#34d399' },
          ]),
        },
      },
    ],
  })
}

async function loadHomeCharts(): Promise<void> {
  const range = activeRange.value
  homeLoading.value = true
  try {
    const hc = await getHomeCharts({
      start_date: range.start,
      end_date: range.end,
      granularity: rangeGranularity(range),
      ...storeParam.value,
    })
    homeCharts.value = hc
    await paintChartsWhenReady(hc)
  } catch {
    homeCharts.value = null
  } finally {
    homeLoading.value = false
  }
}

/** 等图表容器 ref 挂载后再 init/setOption（避免 onMounted 首帧 ref 仍为 null 导致永远不画图） */
async function paintChartsWhenReady(hc: HomeChartsStats): Promise<void> {
  for (let i = 0; i < 12; i++) {
    await nextTick()
    if (lineRef.value && pieRef.value && radarRef.value && memberRankRef.value) {
      applyChartOptions(hc)
      requestAnimationFrame(() => {
        lineChart?.resize()
        pieChart?.resize()
        radarChart?.resize()
        memberRankChart?.resize()
      })
      return
    }
    await new Promise<void>((r) => setTimeout(r, 32))
  }
  applyChartOptions(hc)
}

function onWinResize(): void {
  requestAnimationFrame(() => {
    lineChart?.resize()
    pieChart?.resize()
    radarChart?.resize()
    memberRankChart?.resize()
  })
}

function bindResizeObserver(): void {
  if (!screenRoot.value || typeof ResizeObserver === 'undefined') return
  resizeObserver?.disconnect()
  resizeObserver = new ResizeObserver(() => {
    onWinResize()
  })
  resizeObserver.observe(screenRoot.value)
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
  bindResizeObserver()
  void (async () => {
    await nextTick()
    await nextTick()
    await loadHomeCharts()
  })()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onWinResize)
  resizeObserver?.disconnect()
  resizeObserver = null
  disposeCharts()
})
</script>

<style scoped>
.dash-screen {
  position: relative;
  box-sizing: border-box;
  background: linear-gradient(145deg, #070b14 0%, #0b1220 45%, #0a1628 100%);
  border-radius: 16px;
  height: 100%;
  min-height: 0;
}

.dash-screen--fullscreen {
  width: 100vw;
  height: 100vh;
  border-radius: 0;
  border: 0;
}

.dash-loading-overlay {
  position: absolute;
  z-index: 30;
  inset: 0;
  display: grid;
  place-items: center;
  background:
    linear-gradient(180deg, rgba(7, 11, 20, 0.72), rgba(7, 11, 20, 0.88)),
    radial-gradient(circle at 50% 42%, rgba(34, 211, 238, 0.14), transparent 34%);
  backdrop-filter: blur(6px);
}

.dash-loading-core {
  position: relative;
  display: grid;
  place-items: center;
  width: clamp(180px, 18vw, 260px);
  aspect-ratio: 1;
  color: #ecfeff;
}

.dash-loading-core::before {
  content: "";
  position: absolute;
  inset: 18%;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.72);
  border: 1px solid rgba(34, 211, 238, 0.26);
  box-shadow:
    inset 0 0 28px rgba(34, 211, 238, 0.08),
    0 0 38px rgba(34, 211, 238, 0.14);
}

.dash-loading-core strong,
.dash-loading-core small {
  position: relative;
  z-index: 2;
  text-align: center;
}

.dash-loading-core strong {
  margin-top: 18px;
  font-size: clamp(18px, 1.45vw, 26px);
  font-weight: 900;
}

.dash-loading-core small {
  margin-top: 48px;
  color: rgba(207, 250, 254, 0.68);
  font-size: 13px;
  font-weight: 700;
}

.dash-loading-orbit {
  position: absolute;
  inset: 0;
  border-radius: 999px;
  border: 1px solid rgba(34, 211, 238, 0.18);
  border-top-color: rgba(34, 211, 238, 0.95);
  border-right-color: rgba(167, 139, 250, 0.78);
  animation: dash-loading-spin 1.35s linear infinite;
}

.dash-loading-orbit--second {
  inset: 14%;
  animation-duration: 1.9s;
  animation-direction: reverse;
  border-top-color: rgba(52, 211, 153, 0.9);
  border-right-color: rgba(34, 211, 238, 0.46);
}

.dash-loading-dot {
  position: absolute;
  z-index: 2;
  width: 10px;
  height: 10px;
  border-radius: 999px;
  background: #22d3ee;
  box-shadow: 0 0 16px #22d3ee;
  animation: dash-loading-pulse 1s ease-in-out infinite alternate;
}

.dash-screen__content {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  gap: clamp(10px, 1.1vh, 18px);
  padding: clamp(12px, 1.45vw, 28px);
  overflow: hidden;
}

.dash-screen__bg {
  background-image:
    radial-gradient(ellipse 80% 50% at 50% -20%, rgba(34, 211, 238, 0.18), transparent),
    linear-gradient(rgba(34, 211, 238, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(34, 211, 238, 0.04) 1px, transparent 1px);
  background-size: 100% 100%, 24px 24px, 24px 24px;
  animation: dash-bg-drift 18s linear infinite;
}

.dash-screen::before {
  content: "";
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 16% 18%, rgba(34, 211, 238, 0.16), transparent 20%),
    radial-gradient(circle at 82% 38%, rgba(167, 139, 250, 0.12), transparent 24%);
  opacity: 0.68;
  animation: dash-glow-drift 12s ease-in-out infinite alternate;
}

.dash-scanline {
  height: 120px;
  background: linear-gradient(180deg, transparent, rgba(34, 211, 238, 0.08), transparent);
  animation: dash-scan 7s ease-in-out infinite;
}

.dash-screen__topbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 0;
}

.dash-screen__title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 0;
  color: #cffafe;
  font-size: clamp(18px, 1.55vw, 30px);
  font-weight: 700;
  line-height: 1.15;
}

.dash-screen__pulse {
  display: inline-block;
  width: 8px;
  height: 8px;
  flex: 0 0 auto;
  border-radius: 999px;
  background: #22d3ee;
  box-shadow: 0 0 12px #22d3ee;
}

.dash-screen__subtitle {
  margin: 6px 0 0;
  color: rgba(207, 250, 254, 0.94);
  font-size: clamp(12px, 1.1vw, 18px);
  line-height: 1.25;
}

.dash-chart-grid {
  flex: 1 1 auto;
  min-height: 0;
  display: grid;
  grid-template-columns: repeat(12, minmax(0, 1fr));
  grid-template-rows: minmax(0, 0.95fr) minmax(0, 0.82fr) minmax(0, 0.82fr);
  gap: clamp(10px, 1vw, 16px);
}

.dash-panel {
  position: relative;
  box-sizing: border-box;
  min-height: 0;
  overflow: hidden;
  animation: dash-panel-rise 0.7s ease both;
}

.dash-panel::before {
  content: "";
  position: absolute;
  inset: 0;
  pointer-events: none;
  border-radius: inherit;
  background: linear-gradient(135deg, rgba(34, 211, 238, 0.08), transparent 38%, rgba(167, 139, 250, 0.08));
  opacity: 0.52;
}

.dash-panel--line {
  grid-column: 6 / 13;
  grid-row: 3;
}

.dash-panel--pie {
  grid-column: 6 / 10;
  grid-row: 1;
}

.dash-panel--radar {
  grid-column: 10 / 13;
  grid-row: 1;
}

.dash-panel--overview {
  grid-column: 1 / 6;
  grid-row: 1 / 4;
  overflow: hidden;
}

.dash-panel--member {
  grid-column: 6 / 10;
  grid-row: 2;
  overflow: hidden;
}

.dash-panel--flow {
  grid-column: 10 / 13;
  grid-row: 2;
  overflow: hidden;
}

.dash-chart {
  position: relative;
  z-index: 1;
  height: 100%;
  min-height: 0;
}

.dash-radar-layout {
  position: relative;
  z-index: 1;
  height: 100%;
  min-height: 0;
  display: block;
}

.dash-radar-chart {
  height: 100%;
}

.dash-pulse-item,
.dash-overview-card,
.dash-mini-stat {
  box-sizing: border-box;
  border: 1px solid rgba(34, 211, 238, 0.18);
  border-radius: 8px;
  background: rgba(15, 23, 42, 0.44);
  box-shadow: inset 0 0 18px rgba(34, 211, 238, 0.025);
}

.dash-pulse-item {
  padding: 10px 12px;
}

.dash-pulse-item span,
.dash-overview-card span,
.dash-mini-stat span {
  display: block;
  color: rgba(207, 250, 254, 0.76);
  font-size: 12px;
  font-weight: 700;
}

.dash-pulse-item strong {
  display: block;
  margin-top: 5px;
  color: #f8fafc;
  font-size: clamp(17px, 1.25vw, 24px);
  line-height: 1.1;
}

.tone-good {
  color: #a7f3d0 !important;
}

.tone-bad {
  color: #fca5a5 !important;
}

.dash-category-flow {
  position: relative;
  z-index: 1;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.dash-section-title {
  color: #ecfeff;
  font-size: 15px;
  font-weight: 800;
  letter-spacing: 0.02em;
}

.dash-flow-list {
  display: grid;
  gap: clamp(8px, 0.8vh, 12px);
  margin-top: 12px;
}

.dash-flow-row {
  display: grid;
  grid-template-columns: minmax(4em, 0.8fr) minmax(80px, 1fr) minmax(64px, auto);
  align-items: center;
  gap: 8px;
  color: #e2e8f0;
  font-size: clamp(13px, 0.8vw, 16px);
}

.dash-flow-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dash-flow-track {
  height: 9px;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.9);
}

.dash-flow-track i {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #22d3ee, #a78bfa);
  box-shadow: 0 0 12px rgba(34, 211, 238, 0.35);
  animation: dash-bar-grow 0.9s ease both;
}

.dash-overview-head {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 12px;
}

.dash-section-title--main {
  font-size: clamp(20px, 1.75vw, 34px);
}

.dash-live-pill {
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 9px;
  border: 1px solid rgba(34, 211, 238, 0.34);
  border-radius: 999px;
  color: #a5f3fc;
  background: rgba(8, 145, 178, 0.14);
  font-size: 12px;
  font-weight: 800;
}

.dash-overview-metrics {
  position: relative;
  z-index: 1;
  min-height: 0;
  overflow: hidden;
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: clamp(7px, 0.58vw, 10px);
  align-content: start;
}

.dash-overview-focus {
  position: relative;
  z-index: 1;
  display: grid;
  height: 100%;
  min-height: 0;
  grid-template-columns: minmax(0, 1fr);
  grid-template-rows: clamp(132px, 18vh, 210px) minmax(0, 1fr);
  gap: clamp(10px, 0.78vw, 14px);
  align-items: stretch;
}

.dash-primary-stat {
  box-sizing: border-box;
  display: flex;
  min-height: 0;
  height: 100%;
  overflow: hidden;
  flex-direction: column;
  justify-content: center;
  border: 1px solid rgba(34, 211, 238, 0.22);
  border-radius: 12px;
  background:
    radial-gradient(circle at 50% 20%, rgba(34, 211, 238, 0.14), transparent 55%),
    rgba(15, 23, 42, 0.48);
  padding: clamp(18px, 1.4vw, 28px);
  box-shadow: inset 0 0 30px rgba(34, 211, 238, 0.04), 0 0 26px rgba(34, 211, 238, 0.05);
}

.dash-primary-stat span {
  display: block;
  color: rgba(207, 250, 254, 0.78);
  font-size: clamp(14px, 0.9vw, 18px);
  font-weight: 800;
}

.dash-primary-stat strong {
  display: block;
  margin-top: 8px;
  font-size: clamp(36px, 3.15vw, 64px);
  line-height: 0.95;
}

:deep(.dash-primary-stat .count-up-number),
:deep(.dash-overview-card .count-up-number),
:deep(.dash-flow-row .count-up-number) {
  color: inherit;
  font: inherit;
  line-height: inherit;
}

.dash-primary-stat small {
  margin-top: 14px;
  color: rgba(207, 250, 254, 0.72);
  font-size: 13px;
  font-weight: 700;
}

.dash-overview-card {
  position: relative;
  overflow: hidden;
  min-height: 0;
  padding: clamp(8px, 0.62vw, 12px);
  animation: dash-card-in 0.65s ease both;
  animation-delay: var(--delay, 0ms);
}

.dash-overview-card::after {
  content: "";
  position: absolute;
  inset: 0;
  pointer-events: none;
  transform: translateX(-120%);
  background: linear-gradient(100deg, transparent, rgba(125, 211, 252, 0.12), transparent);
  animation: dash-shimmer 4.6s ease-in-out infinite;
}

.dash-overview-card strong {
  display: block;
  margin-top: 5px;
  color: #fff;
  font-size: clamp(16px, 1.05vw, 23px);
  line-height: 1.05;
}

.dash-flow-layout {
  position: relative;
  z-index: 1;
  height: 100%;
  min-height: 0;
  display: block;
}

.dash-member-layout {
  position: relative;
  z-index: 1;
  display: grid;
  height: 100%;
  min-height: 0;
  grid-template-rows: auto minmax(0, 1fr);
  gap: 10px;
}

.dash-member-chart {
  min-height: 0;
}

.dash-native-control {
  height: 36px;
  box-sizing: border-box;
  border: 1px solid rgba(34, 211, 238, 0.58);
  border-radius: 4px;
  outline: none;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.98), rgba(8, 13, 27, 0.98));
  color: #f0fdfa;
  font: inherit;
  font-size: 15px;
  font-weight: 700;
  box-shadow: inset 0 0 0 1px rgba(125, 211, 252, 0.04), 0 0 18px rgba(34, 211, 238, 0.05);
}

.dash-native-control:focus {
  border-color: rgba(125, 211, 252, 0.95);
  box-shadow: 0 0 0 2px rgba(34, 211, 238, 0.16), 0 0 18px rgba(34, 211, 238, 0.08);
}

.dash-native-select {
  width: 160px;
  padding: 0 34px 0 12px;
  appearance: none;
  background-image:
    linear-gradient(45deg, transparent 50%, #cffafe 50%),
    linear-gradient(135deg, #cffafe 50%, transparent 50%),
    linear-gradient(180deg, rgba(15, 23, 42, 0.98), rgba(8, 13, 27, 0.98));
  background-position:
    calc(100% - 18px) 15px,
    calc(100% - 12px) 15px,
    0 0;
  background-size: 6px 6px, 6px 6px, 100% 100%;
  background-repeat: no-repeat;
}

.dash-native-select--short {
  width: 128px;
}

.dash-native-select option {
  background: #0f172a;
  color: #e0f2fe;
  font-weight: 700;
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

.dash-btn-admin {
  background: rgba(15, 23, 42, 0.88) !important;
  border: 1px solid rgba(125, 211, 252, 0.75) !important;
  color: #ecfeff !important;
  font-weight: 700 !important;
  font-size: 14px !important;
}

@media (max-width: 1180px) {
  .dash-chart-grid {
    grid-template-columns: repeat(8, minmax(0, 1fr));
    grid-template-rows: minmax(0, 1fr) minmax(0, 0.75fr) minmax(0, 0.78fr) minmax(0, 0.78fr);
  }

  .dash-panel--overview {
    grid-column: 1 / 9;
    grid-row: 1;
  }

  .dash-panel--pie {
    grid-column: 1 / 5;
    grid-row: 2;
  }

  .dash-panel--radar {
    grid-column: 5 / 9;
    grid-row: 2;
  }

  .dash-panel--flow {
    grid-column: 5 / 9;
    grid-row: 3;
  }

  .dash-panel--member {
    grid-column: 1 / 5;
    grid-row: 3;
  }

  .dash-panel--line {
    grid-column: 1 / 9;
    grid-row: 4;
  }

  .dash-overview-focus {
    grid-template-columns: minmax(220px, 0.8fr) minmax(0, 1.7fr);
    grid-template-rows: minmax(0, 1fr);
  }
}

@media (max-width: 820px) {
  .dash-screen__content {
    gap: 8px;
    padding: 10px;
  }

  .dash-screen__topbar {
    align-items: flex-start;
  }

  .dash-controls {
    justify-content: flex-end;
  }

  .dash-chart-grid {
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: minmax(0, 1.2fr) minmax(0, 0.8fr) minmax(0, 0.8fr) minmax(0, 0.8fr);
  }

  .dash-overview-focus {
    grid-template-columns: minmax(0, 1fr);
  }

  .dash-overview-metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dash-radar-layout {
    grid-template-columns: minmax(0, 1fr);
    grid-template-rows: minmax(0, 1fr);
  }

  .dash-panel--overview,
  .dash-panel--pie,
  .dash-panel--radar,
  .dash-panel--member,
  .dash-panel--flow,
  .dash-panel--line {
    grid-column: 1;
  }

  .dash-panel--overview {
    grid-row: 1;
  }

  .dash-panel--pie {
    grid-row: 2;
  }

  .dash-panel--member {
    grid-row: 3;
  }

  .dash-panel--line {
    grid-row: 4;
  }

  .dash-panel--radar,
  .dash-panel--flow {
    display: none;
  }

  .dash-native-select {
    width: 116px !important;
  }
}

@keyframes dash-bg-drift {
  from {
    background-position: 0 0, 0 0, 0 0;
  }

  to {
    background-position: 0 0, 0 24px, 24px 0;
  }
}

@keyframes dash-glow-drift {
  from {
    transform: translate3d(-1%, -1%, 0);
  }

  to {
    transform: translate3d(1%, 1%, 0);
  }
}

@keyframes dash-scan {
  0% {
    transform: translateY(-140px);
    opacity: 0;
  }

  18%,
  70% {
    opacity: 1;
  }

  100% {
    transform: translateY(100vh);
    opacity: 0;
  }
}

@keyframes dash-card-in {
  from {
    opacity: 0;
    transform: translateY(10px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes dash-panel-rise {
  from {
    opacity: 0;
    transform: translateY(14px);
  }

  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes dash-shimmer {

  0%,
  58% {
    transform: translateX(-120%);
  }

  100% {
    transform: translateX(120%);
  }
}

@keyframes dash-bar-grow {
  from {
    width: 0;
  }
}

@keyframes dash-loading-spin {
  to {
    transform: rotate(360deg);
  }
}

@keyframes dash-loading-pulse {
  from {
    opacity: 0.48;
    transform: scale(0.72);
  }

  to {
    opacity: 1;
    transform: scale(1.08);
  }
}
</style>
