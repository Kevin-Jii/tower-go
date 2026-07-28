<template>
  <section id="dash-analytics" ref="screenRoot" class="dash-screen" :class="{ 'dash-screen--fullscreen': fullscreen }"
    aria-labelledby="dash-title">
    <div class="dash-screen__grid" aria-hidden="true" />
    <div class="dash-screen__content">
      <header class="dash-header">
        <div class="dash-heading">
          <h1 id="dash-title">经营数据大屏</h1>
          <span class="dash-heading__info" title="经营统计以当前门店和所选周期为准">
            <IconInfoCircle />
          </span>
          <span class="dash-live-status"><i />实时数据</span>
        </div>

        <div class="dash-toolbar" aria-label="数据大屏操作">
          <div class="dash-periods" role="tablist" aria-label="统计周期">
            <button v-for="option in periodOptions" :key="option.value" type="button" class="dash-period"
              :class="{ 'dash-period--active': period === option.value }" role="tab"
              :aria-selected="period === option.value" @click="period = option.value">
              {{ option.label }}
            </button>
          </div>

          <div class="dash-range-control" title="当前统计日期范围">
            <IconCalendar />
            <span>{{ activeRange.start }} ~ {{ activeRange.end }}</span>
          </div>

          <button class="dash-icon-button" type="button" title="刷新数据" aria-label="刷新数据" @click="refreshDash">
            <IconRefresh :class="{ 'dash-icon-spin': screenLoading }" />
          </button>
          <button class="dash-icon-button" type="button" :title="browserFullscreen ? '退出全屏' : '全屏显示'"
            :aria-label="browserFullscreen ? '退出全屏' : '全屏显示'" @click="toggleFullscreen">
            <IconFullscreenExit v-if="browserFullscreen" />
            <IconFullscreen v-else />
          </button>
          <button class="dash-admin-button" type="button" @click="goAdmin">
            <IconApps />
            后台入口
          </button>
          <button class="dash-export-button" type="button" @click="exportData">
            <IconDownload />
            导出数据
          </button>
        </div>
      </header>

      <p v-if="dashError" class="dash-error" role="alert">
        {{ dashError }}
        <button type="button" @click="refreshDash">重新加载</button>
      </p>

      <main class="dash-main">
        <section class="dash-hero-panel dash-panel" aria-label="经营关键指标">
          <div class="dash-hero-summary">
            <span class="dash-eyebrow">{{ periodLabel }}净利润</span>
            <strong class="dash-hero-value" :class="profitTone">
              <CountUpNumber :value="overview?.net_profit_amount" prefix="¥" :decimals="2" :use-grouping="true" />
            </strong>
            <span class="dash-hero-meta">{{ activeRange.start }} 至 {{ activeRange.end }} · 统计已同步</span>
            <div class="dash-hero-chart" ref="heroChartRef" aria-hidden="true" />
          </div>

          <div class="dash-metric-grid">
            <article v-for="metric in metricCards" :key="metric.label" class="dash-metric-card">
              <span class="dash-metric-icon" :class="metric.iconTone">
                <component :is="metric.icon" />
              </span>
              <div class="dash-metric-copy">
                <span class="dash-metric-label">{{ metric.label }}</span>
                <strong :class="metric.tone">
                  <CountUpNumber :value="metric.value" :prefix="metric.prefix" :suffix="metric.suffix"
                    :decimals="metric.decimals" :use-grouping="true" />
                </strong>
                <small>{{ metric.note }}</small>
              </div>
            </article>
          </div>
        </section>

        <section class="dash-middle-grid">
          <article class="dash-chart-panel dash-panel dash-panel--trend">
            <div class="dash-panel-heading">
              <div>
                <h2>经营趋势</h2>
                <p>销售额与订单量 · {{ periodLabel }}</p>
              </div>
              <div class="dash-chart-legend" aria-label="图表图例">
                <span><i class="legend-dot legend-dot--sales" />销售额</span>
                <span><i class="legend-dot legend-dot--orders" />订单数</span>
              </div>
            </div>
            <div ref="lineRef" class="dash-chart" role="img" aria-label="经营趋势折线图" />
            <p v-if="!lineRows.length" class="dash-empty">暂无趋势数据</p>
          </article>

        </section>

        <section class="dash-lower-grid">
          <article class="dash-data-panel dash-panel">
            <div class="dash-panel-heading">
              <div>
                <h2>销售渠道占比</h2>
              </div>
              <span class="dash-panel-more">{{ channelRows.length }} 个渠道</span>
            </div>
            <div class="dash-channel-layout">
              <div ref="pieRef" class="dash-chart dash-pie-chart" role="img" aria-label="销售渠道占比环形图" />
              <div class="dash-channel-list">
                <div v-for="item in channelRows" :key="item.key" class="dash-channel-row">
                  <span class="dash-channel-name"><i :style="{ background: item.color }" />{{ item.name }}</span>
                  <span>{{ item.percent.toFixed(1) }}%</span>
                  <strong>{{ formatMoney(item.amount) }}</strong>
                </div>
                <p v-if="!channelRows.length" class="dash-empty">暂无渠道数据</p>
              </div>
            </div>
          </article>

          <article class="dash-data-panel dash-panel">
            <div class="dash-panel-heading">
              <div>
                <h2>品类销售 TOP10</h2>
              </div>
              <span class="dash-panel-more">{{ categoryRows.length }} 个品类</span>
            </div>
            <div class="dash-ranking-list">
              <div v-for="(item, index) in categoryRows" :key="item.key" class="dash-ranking-row">
                <span class="dash-ranking-index" :class="{ 'dash-ranking-index--top': index < 3 }">{{ index + 1
                  }}</span>
                <span class="dash-ranking-name">{{ item.name }}</span>
                <span class="dash-ranking-track"><i :style="{ width: `${item.percent}%` }" /></span>
                <strong>{{ formatMoney(item.amount) }}</strong>
              </div>
              <p v-if="!categoryRows.length" class="dash-empty">暂无品类数据</p>
            </div>
          </article>

          <article class="dash-data-panel dash-panel">
            <div class="dash-panel-heading">
              <div>
                <h2>会员消费排行</h2>
              </div>
              <span class="dash-panel-more">{{ memberRows.length }} 位会员</span>
            </div>
            <div class="dash-ranking-list dash-ranking-list--member">
              <div v-for="(item, index) in memberRows" :key="item.key" class="dash-ranking-row">
                <span class="dash-ranking-index" :class="{ 'dash-ranking-index--top': index < 3 }">{{ index + 1
                  }}</span>
                <span class="dash-ranking-name">{{ item.name }}</span>
                <span class="dash-ranking-orders">{{ item.orders }} 单</span>
                <strong>{{ formatMoney(item.amount) }}</strong>
              </div>
              <p v-if="!memberRows.length" class="dash-empty">暂无会员消费数据</p>
            </div>
          </article>

        </section>

        <section class="dash-footer-grid" aria-label="库存与经营指标">
          <article v-for="metric in footerCards" :key="metric.label" class="dash-footer-card dash-panel"
            :class="metric.cardTone">
            <span class="dash-footer-icon">
              <component :is="metric.icon" />
            </span>
            <div>
              <span>{{ metric.label }}</span>
              <strong>
                <CountUpNumber :value="metric.value" :prefix="metric.prefix" :suffix="metric.suffix"
                  :decimals="metric.decimals" :use-grouping="true" />
              </strong>
              <small>{{ metric.note }}</small>
            </div>
          </article>
        </section>
      </main>
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
import { computed, markRaw, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import * as echarts from 'echarts'
import {
  IconApps,
  IconCalendar,
  IconDownload,
  IconFullscreen,
  IconFullscreenExit,
  IconInfoCircle,
  IconRefresh,
} from '@arco-design/web-vue/es/icon'
import { useRouter } from 'vue-router'
import CountUpNumber from '@/components/CountUpNumber.vue'
import { getHomeCharts, getStatisticsDashboard } from '@/api/statistics'
import type { HomeChartsStats } from '@/api/types'
import { useUserStore } from '@/store/user'
import {
  DashboardAverageIcon,
  DashboardDepositIcon,
  DashboardInboundIcon,
  DashboardLossIcon,
  DashboardMarginIcon,
  DashboardOrdersIcon,
  DashboardOutboundIcon,
  DashboardProfitIcon,
  DashboardRoundIcon,
  DashboardSalesIcon,
  DashboardStockIcon,
} from './components/icons'

defineProps<{
  fullscreen?: boolean
}>()

const userStore = useUserStore()
const router = useRouter()
const screenRoot = ref<HTMLElement | null>(null)
const period = ref('month')
const browserFullscreen = ref(false)
const periodOptions = [
  { label: '今日', value: 'today' },
  { label: '昨日', value: 'yesterday' },
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' },
  { label: '本季度', value: 'quarter' },
  { label: '今年', value: 'year' },
]

const currentStoreId = computed(() => userStore.currentStoreId || 0)
const storeParam = computed(() => (currentStoreId.value > 0 ? { store_id: currentStoreId.value } : {}))
const dashboardPeriod = computed(() => (period.value === 'yesterday' ? 'today' : period.value))
const dashKey = computed(() => ['statistics', 'dashboard', dashboardPeriod.value, period.value, currentStoreId.value] as const)
const {
  data: dash,
  isPending: dashPending,
  isError: dashIsError,
  error: dashQueryError,
  isFetching: dashFetching,
  refetch: refetchDash,
} = useQuery({
  queryKey: dashKey,
  queryFn: () => getStatisticsDashboard({ period: dashboardPeriod.value, ...storeParam.value }),
})

const dashError = computed(() => {
  if (!dashIsError.value) return ''
  return dashQueryError.value instanceof Error ? dashQueryError.value.message : '周期概览加载失败'
})

function dateText(value: Date): string {
  const y = value.getFullYear()
  const m = String(value.getMonth() + 1).padStart(2, '0')
  const d = String(value.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

function periodRange(value: string): { start: string; end: string } {
  const now = new Date()
  const end = dateText(now)
  if (value === 'today') return { start: end, end }
  if (value === 'yesterday') {
    const yesterday = new Date(now)
    yesterday.setDate(now.getDate() - 1)
    const day = dateText(yesterday)
    return { start: day, end: day }
  }
  if (value === 'week') {
    const weekday = now.getDay() || 7
    const start = new Date(now)
    start.setDate(now.getDate() - weekday + 1)
    return { start: dateText(start), end }
  }
  if (value === 'quarter') {
    const quarter = Math.floor(now.getMonth() / 3)
    return { start: dateText(new Date(now.getFullYear(), quarter * 3, 1)), end }
  }
  if (value === 'year') return { start: `${now.getFullYear()}-01-01`, end }
  return { start: dateText(new Date(now.getFullYear(), now.getMonth(), 1)), end }
}

function rangeGranularity(range: { start: string; end: string }): 'day' | 'month' {
  const start = new Date(`${range.start}T00:00:00`)
  const end = new Date(`${range.end}T00:00:00`)
  const days = Math.max(1, Math.round((end.getTime() - start.getTime()) / 86_400_000) + 1)
  return days > 92 ? 'month' : 'day'
}

const activeRange = computed(() => periodRange(period.value))
const periodLabel = computed(() => periodOptions.find((item) => item.value === period.value)?.label || '本月')
const homeCharts = ref<HomeChartsStats | null>(null)
const homeLoading = ref(false)
const overview = computed(() => homeCharts.value?.overview ?? null)
const screenLoading = computed(() => dashPending.value || dashFetching.value || homeLoading.value)

function goAdmin(): void {
  void router.push('/admin')
}

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
const profitTone = computed(() => (Number(overview.value?.net_profit_amount ?? 0) >= 0 ? 'tone-good' : 'tone-bad'))

const iconSet = {
  sales: markRaw(DashboardSalesIcon),
  profit: markRaw(DashboardProfitIcon),
  margin: markRaw(DashboardMarginIcon),
  average: markRaw(DashboardAverageIcon),
  orders: markRaw(DashboardOrdersIcon),
  inbound: markRaw(DashboardInboundIcon),
  outbound: markRaw(DashboardOutboundIcon),
  stock: markRaw(DashboardStockIcon),
  loss: markRaw(DashboardLossIcon),
  round: markRaw(DashboardRoundIcon),
  deposit: markRaw(DashboardDepositIcon),
}

const metricCards = computed(() => [
  {
    label: '销售收入',
    value: Number(overview.value?.sales_amount ?? dash.value?.sales.total_amount ?? 0),
    prefix: '¥',
    suffix: '',
    decimals: 2,
    note: '本期累计',
    icon: iconSet.sales,
    iconTone: 'dash-icon-tone--blue',
    tone: '',
  },
  {
    label: '毛利额',
    value: Number(overview.value?.gross_profit_amount ?? 0),
    prefix: '¥',
    suffix: '',
    decimals: 2,
    note: '本期累计',
    icon: iconSet.profit,
    iconTone: 'dash-icon-tone--purple',
    tone: '',
  },
  {
    label: '毛利率',
    value: grossMargin.value,
    prefix: '',
    suffix: '%',
    decimals: 1,
    note: grossMargin.value >= 0 ? '经营效率稳定' : '需要重点关注',
    icon: iconSet.margin,
    iconTone: 'dash-icon-tone--green',
    tone: grossMargin.value >= 0 ? 'tone-good' : 'tone-bad',
  },
  {
    label: '客单价',
    value: avgOrderAmount.value,
    prefix: '¥',
    suffix: '',
    decimals: 2,
    note: `${Number(overview.value?.sales_order_count ?? 0)} 笔订单`,
    icon: iconSet.average,
    iconTone: 'dash-icon-tone--orange',
    tone: '',
  },
  {
    label: '销售单数',
    value: Number(overview.value?.sales_order_count ?? dash.value?.sales.total_orders ?? 0),
    prefix: '',
    suffix: ' 单',
    decimals: 0,
    note: '本期累计',
    icon: iconSet.orders,
    iconTone: 'dash-icon-tone--blue',
    tone: '',
  },
])

const footerCards = computed(() => [
  { label: '入库金额', value: Number(overview.value?.inbound_amount ?? 0), prefix: '¥', suffix: '', decimals: 2, note: `${Number(overview.value?.inventory_in_count ?? 0)} 单入库`, icon: iconSet.inbound, cardTone: 'dash-footer-card--purple' },
  { label: '出库金额', value: Number(overview.value?.outbound_amount ?? 0), prefix: '¥', suffix: '', decimals: 2, note: `${Number(overview.value?.inventory_out_count ?? 0)} 单出库`, icon: iconSet.outbound, cardTone: 'dash-footer-card--blue' },
  { label: '库存金额', value: Number(overview.value?.all_category_amount ?? 0), prefix: '¥', suffix: '', decimals: 2, note: `${Number(dash.value?.inventory.total_quantity ?? 0)} 件库存`, icon: iconSet.stock, cardTone: 'dash-footer-card--orange' },
  { label: '库存数量', value: Number(dash.value?.inventory.total_quantity ?? 0), prefix: '', suffix: ' 件', decimals: 0, note: `今日入库 ${Number(dash.value?.inventory.today_in ?? 0)} 件`, icon: iconSet.stock, cardTone: 'dash-footer-card--green' },
  { label: '报损金额', value: Number(overview.value?.inventory_loss_amount ?? 0), prefix: '¥', suffix: '', decimals: 2, note: `${Number(overview.value?.inventory_loss_count ?? 0)} 笔报损`, icon: iconSet.loss, cardTone: 'dash-footer-card--red' },
  { label: '抹零金额', value: Number(overview.value?.round_amount ?? 0), prefix: '¥', suffix: '', decimals: 2, note: '本期累计', icon: iconSet.round, cardTone: 'dash-footer-card--blue' },
  { label: '退货押金', value: Number(overview.value?.return_deposit_amount ?? 0), prefix: '¥', suffix: '', decimals: 2, note: `物流费用 ${formatMoney(Number(overview.value?.return_logistics_fee ?? 0))}`, icon: iconSet.deposit, cardTone: 'dash-footer-card--purple' },
])

const lineRows = computed(() => homeCharts.value?.line ?? [])
const channelColors = ['#2f8cff', '#25c7a4', '#f6b73c', '#e36b9a', '#7b61ff', '#35c8f4']
const channelRows = computed(() => {
  return [...(homeCharts.value?.pie ?? [])]
    .map((item, index) => ({
      key: `${item.channel}-${index}`,
      name: item.channel_name || item.channel || `渠道 ${index + 1}`,
      amount: Number(item.amount ?? 0),
      percent: Number(item.percent ?? 0),
      color: channelColors[index % channelColors.length],
    }))
    .sort((a, b) => b.amount - a.amount)
    .slice(0, 6)
})

const categoryRows = computed(() => {
  const rows = [...(overview.value?.categories ?? [])]
    .map((item, index) => ({
      key: `${item.category_id}-${index}`,
      name: item.category_name || '未分类',
      amount: Math.max(Number(item.in_amount ?? 0), Number(item.out_amount ?? 0), Math.abs(Number(item.net_amount ?? 0))),
    }))
    .sort((a, b) => b.amount - a.amount)
    .slice(0, 10)
  const max = Math.max(...rows.map((item) => item.amount), 1)
  return rows.map((item) => ({ ...item, percent: Math.max(8, Math.round((item.amount / max) * 100)) }))
})

const memberRows = computed(() => {
  return [...(overview.value?.member_consumption_rank ?? [])]
    .map((item, index) => ({
      key: `${item.member_id}-${index}`,
      name: String(item.member_name || item.member_phone || '').trim(),
      amount: Number(item.amount ?? 0),
      orders: Number(item.orders ?? 0),
    }))
    .filter((item) => item.name && item.name !== '未知会员')
    .sort((a, b) => b.amount - a.amount)
    .slice(0, 10)
})

const lineRef = ref<HTMLElement | null>(null)
const pieRef = ref<HTMLElement | null>(null)
const heroChartRef = ref<HTMLElement | null>(null)
let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null
let heroChart: echarts.ECharts | null = null
let resizeObserver: ResizeObserver | null = null

const axisText = '#aebdd7'
const gridLine = 'rgba(117, 153, 211, 0.14)'

function isCompactScreen(): boolean {
  return typeof window !== 'undefined' && window.matchMedia('(max-width: 820px)').matches
}

function formatMoney(value: number): string {
  return `¥${Number(value || 0).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}`
}

function disposeCharts(): void {
  lineChart?.dispose()
  pieChart?.dispose()
  heroChart?.dispose()
  lineChart = pieChart = heroChart = null
}

function ensureCharts(): void {
  if (lineRef.value && !lineChart) lineChart = echarts.init(lineRef.value)
  if (pieRef.value && !pieChart) pieChart = echarts.init(pieRef.value)
  if (heroChartRef.value && !heroChart) heroChart = echarts.init(heroChartRef.value)
}

function applyChartOptions(hc: HomeChartsStats): void {
  ensureCharts()
  const compact = isCompactScreen()
  const line = hc.line ?? []

  if (lineChart) {
    lineChart.setOption({
      backgroundColor: 'transparent',
      animationDuration: 520,
      animationEasing: 'cubicOut',
      tooltip: {
        trigger: 'axis',
        appendToBody: true,
        confine: false,
        extraCssText: 'z-index: 1000; box-shadow: 0 8px 24px rgba(0, 0, 0, 0.28);',
        backgroundColor: '#0a1832',
        borderColor: '#1f5ca5',
        textStyle: { color: '#e7f0ff' },
        valueFormatter: (value: unknown) => Number(value || 0).toLocaleString('zh-CN', { maximumFractionDigits: 2 }),
      },
      grid: { left: compact ? 38 : 56, right: compact ? 38 : 54, top: 20, bottom: compact ? 34 : 28 },
      xAxis: {
        type: 'category',
        data: line.map((item) => item.date),
        boundaryGap: false,
        axisLine: { lineStyle: { color: 'rgba(155, 185, 231, 0.32)' } },
        axisTick: { show: false },
        axisLabel: { color: axisText, fontSize: compact ? 10 : 11, margin: 12, rotate: compact || line.length > 12 ? 28 : 0 },
      },
      yAxis: [
        {
          type: 'value',
          splitLine: { lineStyle: { color: gridLine } },
          axisLabel: { color: axisText, fontSize: compact ? 10 : 11, formatter: (value: number) => Number(value).toLocaleString('zh-CN') },
        },
        {
          type: 'value',
          splitLine: { show: false },
          axisLabel: { color: '#8b7bea', fontSize: compact ? 10 : 11 },
        },
      ],
      series: [
        {
          name: '销售额',
          type: 'line',
          smooth: 0.32,
          symbol: 'circle',
          symbolSize: compact ? 4 : 6,
          yAxisIndex: 0,
          lineStyle: { width: 2.5, color: '#2f9bff' },
          itemStyle: { color: '#2f9bff', borderColor: '#d1eaff', borderWidth: 1 },
          areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{ offset: 0, color: 'rgba(47, 155, 255, 0.24)' }, { offset: 1, color: 'rgba(47, 155, 255, 0.01)' }]) },
          data: line.map((item) => item.amount),
        },
        {
          name: '订单数',
          type: 'line',
          smooth: 0.32,
          symbol: 'circle',
          symbolSize: compact ? 3 : 5,
          yAxisIndex: 1,
          lineStyle: { width: 2, color: '#7b61ff' },
          itemStyle: { color: '#7b61ff', borderColor: '#e4ddff', borderWidth: 1 },
          data: line.map((item) => item.orders),
        },
      ],
    })
  }

  if (pieChart) {
    const pie = channelRows.value.length ? channelRows.value : [{ name: '暂无数据', amount: 1, percent: 100, color: '#273755', key: 'empty' }]
    pieChart.setOption({
      backgroundColor: 'transparent',
      animationDuration: 520,
      tooltip: { trigger: 'item', appendToBody: true, confine: false, extraCssText: 'z-index: 1000; box-shadow: 0 8px 24px rgba(0, 0, 0, 0.28);', backgroundColor: '#0a1832', borderColor: '#1f5ca5', textStyle: { color: '#e7f0ff' }, formatter: (params: { name: string; value: number; percent: number }) => `${params.name}<br/>${formatMoney(params.value)} · ${params.percent}%` },
      series: [{
        type: 'pie',
        radius: compact ? ['48%', '70%'] : ['54%', '76%'],
        center: compact ? ['50%', '48%'] : ['50%', '50%'],
        avoidLabelOverlap: true,
        label: { show: false },
        labelLine: { show: false },
        itemStyle: { borderColor: '#0b1830', borderWidth: 3 },
        data: pie.map((item) => ({ name: item.name, value: item.amount, itemStyle: { color: item.color } })),
      }],
    })
  }

  if (heroChart) {
    heroChart.setOption({
      backgroundColor: 'transparent',
      animationDuration: 700,
      grid: { left: 0, right: 0, top: 6, bottom: 2 },
      xAxis: { type: 'category', show: false, data: line.map((item) => item.date) },
      yAxis: { type: 'value', show: false, scale: true },
      series: [{
        type: 'line',
        smooth: 0.42,
        symbol: 'none',
        lineStyle: { width: 2.5, color: '#6955e7' },
        areaStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{ offset: 0, color: 'rgba(105, 85, 231, 0.35)' }, { offset: 1, color: 'rgba(105, 85, 231, 0)' }]) },
        data: line.map((item) => item.amount),
      }],
    })
  }
}

async function loadHomeCharts(): Promise<void> {
  const range = activeRange.value
  const storeId = currentStoreId.value
  const requestKey = `${period.value}:${storeId}:${range.start}:${range.end}`
  homeLoading.value = true
  try {
    const stats = await getHomeCharts({
      start_date: range.start,
      end_date: range.end,
      granularity: rangeGranularity(range),
      ...(storeId > 0 ? { store_id: storeId } : {}),
    })
    const latestKey = `${period.value}:${currentStoreId.value}:${activeRange.value.start}:${activeRange.value.end}`
    if (requestKey !== latestKey) return
    homeCharts.value = stats
    await paintChartsWhenReady(stats)
  } catch {
    homeCharts.value = null
  } finally {
    homeLoading.value = false
  }
}

async function paintChartsWhenReady(stats: HomeChartsStats): Promise<void> {
  for (let index = 0; index < 12; index += 1) {
    await nextTick()
    if (lineRef.value && pieRef.value && heroChartRef.value) {
      applyChartOptions(stats)
      requestAnimationFrame(() => {
        lineChart?.resize()
        pieChart?.resize()
        heroChart?.resize()
      })
      return
    }
    await new Promise<void>((resolve) => setTimeout(resolve, 32))
  }
}

async function refreshDash(): Promise<void> {
  await Promise.all([refetchDash(), loadHomeCharts()])
}

function onWinResize(): void {
  requestAnimationFrame(() => {
    if (homeCharts.value) applyChartOptions(homeCharts.value)
    lineChart?.resize()
    pieChart?.resize()
    heroChart?.resize()
  })
}

async function toggleFullscreen(): Promise<void> {
  if (document.fullscreenElement) {
    await document.exitFullscreen()
    return
  }
  await screenRoot.value?.requestFullscreen?.()
}

function onFullscreenChange(): void {
  browserFullscreen.value = Boolean(document.fullscreenElement)
}

function csvCell(value: unknown): string {
  const text = String(value ?? '')
  return `"${text.replaceAll('"', '""')}"`
}

function exportData(): void {
  const rows: unknown[][] = [
    ['经营数据大屏'],
    ['统计周期', periodLabel.value],
    ['日期范围', `${activeRange.value.start} ~ ${activeRange.value.end}`],
    [],
    ['关键指标', '数值'],
    ...metricCards.value.map((item) => [item.label, `${item.prefix}${Number(item.value).toFixed(item.decimals)}${item.suffix}`]),
    [],
    ['销售趋势', '日期', '销售额', '订单数'],
    ...lineRows.value.map((item) => [item.date, item.amount, item.orders]),
  ]
  const csv = rows.map((row) => row.map(csvCell).join(',')).join('\n')
  const blob = new Blob([`\uFEFF${csv}`], { type: 'text/csv;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `经营数据-${activeRange.value.start}-${activeRange.value.end}.csv`
  link.click()
  URL.revokeObjectURL(url)
}

function bindResizeObserver(): void {
  if (!screenRoot.value || typeof ResizeObserver === 'undefined') return
  resizeObserver?.disconnect()
  resizeObserver = new ResizeObserver(onWinResize)
  resizeObserver.observe(screenRoot.value)
}

watch(period, () => {
  void loadHomeCharts()
})

watch(currentStoreId, () => {
  void loadHomeCharts()
})

watch(homeCharts, (stats) => {
  if (stats) void paintChartsWhenReady(stats)
})

onMounted(() => {
  window.addEventListener('resize', onWinResize)
  document.addEventListener('fullscreenchange', onFullscreenChange)
  bindResizeObserver()
  void (async () => {
    await nextTick()
    await loadHomeCharts()
  })()
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', onWinResize)
  document.removeEventListener('fullscreenchange', onFullscreenChange)
  resizeObserver?.disconnect()
  resizeObserver = null
  disposeCharts()
})
</script>

<style scoped>
.dash-screen {
  --dash-bg: #071227;
  --dash-bg-deep: #050e1f;
  --dash-panel: rgba(11, 29, 58, 0.84);
  --dash-panel-soft: rgba(13, 36, 72, 0.72);
  --dash-border: rgba(41, 110, 198, 0.32);
  --dash-border-strong: rgba(55, 133, 229, 0.5);
  --dash-text: #eff5ff;
  --dash-muted: #aebdd7;
  --dash-dim: #7084a6;
  position: relative;
  box-sizing: border-box;
  width: 100%;
  min-width: 0;
  min-height: 100%;
  overflow: hidden;
  color: var(--dash-text);
  background: var(--dash-bg);
  font-family: Inter, "SF Pro Display", "PingFang SC", "Microsoft YaHei", system-ui, sans-serif;
}

.dash-screen--fullscreen {
  width: 100vw;
  height: 100vh;
  border-radius: 0;
}

.dash-screen__grid {
  position: absolute;
  inset: 0;
  pointer-events: none;
  opacity: 0.45;
  background-image:
    linear-gradient(rgba(76, 132, 209, 0.055) 1px, transparent 1px),
    linear-gradient(90deg, rgba(76, 132, 209, 0.055) 1px, transparent 1px);
  background-size: 42px 42px;
  mask-image: linear-gradient(180deg, rgba(0, 0, 0, 0.92), transparent 78%);
}

.dash-screen__content {
  position: relative;
  z-index: 1;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  min-height: 0;
  gap: 10px;
  padding: 18px;
  overflow: hidden;
}

.dash-header,
.dash-heading,
.dash-toolbar,
.dash-periods,
.dash-range-control,
.dash-icon-button,
.dash-export-button,
.dash-live-status,
.dash-panel-heading,
.dash-chart-legend,
.dash-channel-name,
.dash-structure-label,
.dash-insight-title {
  display: flex;
  align-items: center;
}

.dash-header {
  justify-content: space-between;
  gap: 18px;
  min-height: 42px;
}

.dash-heading {
  gap: 10px;
  flex-wrap: wrap;
}

.dash-heading h1 {
  margin: 0;
  color: #f7faff;
  font-size: 26px;
  font-weight: 760;
  line-height: 1.1;
  letter-spacing: 0;
}

.dash-heading__info {
  display: inline-flex;
  color: #a8bbd9;
  font-size: 17px;
}

.dash-live-status {
  gap: 6px;
  color: #62d9b0;
  font-size: 12px;
  font-weight: 700;
}

.dash-live-status i {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #42d29a;
  box-shadow: 0 0 8px rgba(66, 210, 154, 0.8);
}

.dash-toolbar {
  justify-content: flex-end;
  gap: 8px;
  min-width: 0;
  flex-wrap: wrap;
}

.dash-periods {
  gap: 2px;
  padding: 3px;
  border: 1px solid rgba(37, 98, 170, 0.28);
  border-radius: 9px;
  background: rgba(8, 23, 49, 0.8);
}

.dash-period,
.dash-icon-button,
.dash-admin-button,
.dash-export-button,
.dash-error button {
  border: 0;
  font: inherit;
  cursor: pointer;
}

.dash-period {
  min-width: 58px;
  height: 30px;
  padding: 0 12px;
  border-radius: 6px;
  color: #a9b9d4;
  background: transparent;
  font-size: 13px;
  font-weight: 650;
  transition: color 180ms ease, background 180ms ease;
}

.dash-period:hover,
.dash-period:focus-visible {
  color: #f3f7ff;
  background: rgba(40, 92, 164, 0.38);
  outline: none;
}

.dash-period--active {
  color: #f7fbff;
  background: linear-gradient(180deg, rgba(33, 103, 207, 0.92), rgba(22, 70, 151, 0.92));
  box-shadow: 0 0 0 1px rgba(78, 150, 243, 0.35) inset;
}

.dash-range-control {
  gap: 8px;
  height: 36px;
  padding: 0 12px;
  border: 1px solid rgba(37, 98, 170, 0.32);
  border-radius: 8px;
  color: #bdcbe2;
  background: rgba(8, 24, 50, 0.82);
  font-size: 13px;
  white-space: nowrap;
}

.dash-range-control svg {
  color: #90a7c8;
  font-size: 16px;
}

.dash-icon-button,
.dash-admin-button,
.dash-export-button {
  justify-content: center;
  height: 36px;
  border: 1px solid rgba(37, 98, 170, 0.34);
  border-radius: 8px;
  color: #d8e6fa;
  background: rgba(8, 24, 50, 0.84);
  transition: border-color 180ms ease, background 180ms ease, color 180ms ease;
}

.dash-icon-button {
  width: 38px;
  font-size: 17px;
}

.dash-export-button {
  gap: 7px;
  padding: 0 14px;
  font-size: 13px;
  font-weight: 700;
}

.dash-admin-button {
  gap: 7px;
  padding: 0 13px;
  color: #b9cfff;
  background: rgba(19, 50, 100, 0.76);
  font-size: 13px;
  font-weight: 700;
}

.dash-icon-button:hover,
.dash-icon-button:focus-visible,
.dash-admin-button:hover,
.dash-admin-button:focus-visible,
.dash-export-button:hover,
.dash-export-button:focus-visible {
  border-color: rgba(75, 143, 235, 0.72);
  color: #ffffff;
  background: rgba(24, 66, 129, 0.84);
  outline: none;
}

.dash-main {
  display: grid;
  flex: 1 1 auto;
  min-height: 0;
  grid-template-rows: minmax(0, 0.72fr) minmax(0, 1.28fr) minmax(0, 1.08fr) minmax(74px, 0.48fr);
  gap: 10px;
  overflow: hidden;
}

.dash-panel {
  position: relative;
  box-sizing: border-box;
  overflow: hidden;
  border: 1px solid var(--dash-border);
  border-radius: 11px;
  background: var(--dash-panel);
}

.dash-panel::after {
  position: absolute;
  inset: 0;
  pointer-events: none;
  content: "";
  background: linear-gradient(130deg, rgba(47, 132, 246, 0.08), transparent 34%, rgba(94, 70, 225, 0.05));
}

.dash-hero-panel {
  display: grid;
  grid-template-columns: minmax(330px, 1.02fr) minmax(0, 2.08fr);
  min-height: 0;
  padding: 18px;
  background:
    linear-gradient(105deg, rgba(19, 53, 122, 0.42), rgba(10, 28, 61, 0.72) 48%, rgba(11, 31, 60, 0.82)),
    var(--dash-panel);
}

.dash-hero-summary {
  position: relative;
  z-index: 1;
  display: grid;
  align-content: center;
  min-width: 0;
  padding-right: 20px;
}

.dash-eyebrow,
.dash-metric-label,
.dash-footer-card span,
.dash-hero-meta,
.dash-panel-heading p,
.dash-insight-footer span,
.dash-structure-summary span {
  color: var(--dash-muted);
  font-size: 12px;
  font-weight: 650;
}

.dash-hero-value {
  margin-top: 7px;
  color: #f3f7ff;
  font-size: clamp(36px, 3.2vw, 56px);
  font-weight: 760;
  line-height: 1;
  letter-spacing: 0;
}

.dash-hero-meta {
  margin-top: 14px;
  color: #8298bb;
}

.dash-hero-chart {
  position: absolute;
  right: 16px;
  bottom: 14px;
  width: min(44%, 300px);
  height: 80px;
  opacity: 0.9;
}

.dash-metric-grid {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px;
  min-width: 0;
}

.dash-metric-card {
  display: flex;
  min-width: 0;
  gap: 10px;
  padding: 16px 12px 12px;
  border: 1px solid rgba(41, 111, 201, 0.3);
  border-radius: 9px;
  background: rgba(8, 25, 54, 0.58);
}

.dash-metric-icon,
.dash-footer-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex: 0 0 auto;
  width: 36px;
  height: 36px;
  border-radius: 11px;
  font-size: 18px;
}

.dash-icon-tone--blue {
  color: #65adff;
  background: rgba(35, 112, 231, 0.25);
}

.dash-icon-tone--purple {
  color: #9c80ff;
  background: rgba(90, 58, 209, 0.28);
}

.dash-icon-tone--green {
  color: #39d1a4;
  background: rgba(13, 156, 111, 0.24);
}

.dash-icon-tone--orange {
  color: #f5a247;
  background: rgba(192, 99, 28, 0.24);
}

.dash-metric-copy {
  display: grid;
  align-content: start;
  min-width: 0;
}

.dash-metric-label {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dash-metric-copy strong {
  margin-top: 8px;
  color: #f4f7ff;
  font-size: clamp(18px, 1.2vw, 25px);
  font-weight: 720;
  line-height: 1.05;
  white-space: nowrap;
}

.dash-metric-copy small,
.dash-footer-card small {
  margin-top: 9px;
  color: #7890b5;
  font-size: 11px;
  white-space: nowrap;
}

.dash-middle-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr);
  gap: 10px;
  min-height: 0;
}

.dash-chart-panel,
.dash-insight-panel,
.dash-data-panel {
  min-width: 0;
  padding: 16px;
}

.dash-chart-panel,
.dash-insight-panel {
  min-height: 0;
}

.dash-panel-heading {
  position: relative;
  z-index: 1;
  justify-content: space-between;
  gap: 12px;
  min-width: 0;
}

.dash-panel-heading h2 {
  margin: 0;
  color: #edf4ff;
  font-size: 16px;
  font-weight: 720;
  line-height: 1.2;
}

.dash-panel-heading p {
  margin: 5px 0 0;
  color: #7891b6;
  font-size: 11px;
  font-weight: 550;
}

.dash-chart-legend {
  gap: 16px;
  color: #9eb0cc;
  font-size: 11px;
  white-space: nowrap;
}

.dash-chart-legend span {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.legend-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.legend-dot--sales {
  background: #2f9bff;
}

.legend-dot--orders {
  background: #7b61ff;
}

.dash-chart {
  position: relative;
  z-index: 1;
  width: 100%;
  height: calc(100% - 44px);
  min-height: 246px;
  margin-top: 8px;
}

.dash-insight-title {
  gap: 10px;
}

.dash-ai-mark {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: 1px solid rgba(124, 100, 255, 0.54);
  border-radius: 50%;
  color: #b5a8ff;
  background: rgba(75, 61, 183, 0.3);
  font-size: 12px;
  font-weight: 800;
}

.dash-insight-bolt {
  color: #55a7ff;
  font-size: 21px;
}

.dash-insight-list {
  position: relative;
  z-index: 1;
  display: grid;
  gap: 8px;
  margin-top: 14px;
}

.dash-insight-row {
  display: grid;
  grid-template-columns: 18px minmax(0, 1fr);
  gap: 9px;
  min-width: 0;
  padding: 11px 12px;
  border: 1px solid rgba(49, 109, 190, 0.28);
  border-radius: 9px;
  background: rgba(10, 29, 61, 0.72);
}

.dash-insight-row--warning {
  border-color: rgba(222, 154, 44, 0.3);
}

.dash-insight-row--success {
  border-color: rgba(38, 180, 133, 0.28);
}

.dash-insight-icon {
  margin-top: 1px;
  color: #63aafc;
  font-size: 15px;
}

.dash-insight-row--warning .dash-insight-icon {
  color: #efa72f;
}

.dash-insight-row--success .dash-insight-icon {
  color: #35c995;
}

.dash-insight-row strong {
  color: #c7d7ed;
  font-size: 12px;
  font-weight: 700;
}

.dash-insight-row p {
  margin: 4px 0 0;
  color: #8fa4c2;
  font-size: 11px;
  line-height: 1.5;
}

.dash-insight-footer {
  position: relative;
  z-index: 1;
  display: grid;
  align-content: end;
  min-height: 54px;
  margin-top: 10px;
  padding: 10px 12px;
  border-top: 1px solid rgba(53, 108, 184, 0.24);
}

.dash-insight-footer strong {
  margin-top: 4px;
  color: #4c9dff;
  font-size: 23px;
  line-height: 1;
}

.dash-lower-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 10px;
  min-height: 0;
}

.dash-data-panel {
  min-height: 0;
}

.dash-panel-more {
  color: #7690b7;
  font-size: 11px;
  white-space: nowrap;
}

.dash-channel-layout {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: minmax(130px, 0.88fr) minmax(0, 1.12fr);
  gap: 8px;
  align-items: center;
  height: 214px;
  margin-top: 6px;
}

.dash-pie-chart {
  min-height: 190px;
  height: 210px;
  margin-top: 0;
}

.dash-channel-list,
.dash-ranking-list,
.dash-structure-list {
  position: relative;
  z-index: 1;
  min-width: 0;
}

.dash-channel-list {
  display: grid;
  gap: 10px;
}

.dash-channel-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 3px 8px;
  color: #a8b9d3;
  font-size: 11px;
}

.dash-channel-row strong {
  grid-column: 1 / -1;
  color: #e5edf9;
  font-size: 12px;
  font-weight: 680;
}

.dash-channel-name {
  min-width: 0;
  gap: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dash-channel-name i {
  width: 7px;
  height: 7px;
  flex: 0 0 auto;
  border-radius: 50%;
}

.dash-ranking-list {
  display: grid;
  gap: 7px;
  margin-top: 16px;
}

.dash-ranking-row {
  display: grid;
  grid-template-columns: 22px minmax(52px, 0.75fr) minmax(45px, 1fr) auto;
  gap: 7px;
  align-items: center;
  min-width: 0;
  color: #a9b9d2;
  font-size: 11px;
}

.dash-ranking-index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 17px;
  height: 17px;
  border-radius: 4px;
  color: #90a4c3;
  background: rgba(47, 74, 117, 0.48);
  font-size: 10px;
  font-weight: 760;
}

.dash-ranking-index--top {
  color: #fff;
  background: #df8d30;
}

.dash-ranking-name {
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dash-ranking-track,
.dash-structure-track {
  display: block;
  height: 6px;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(27, 55, 98, 0.8);
}

.dash-ranking-track i,
.dash-structure-track i {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #238ef7, #7a5ef2);
}

.dash-ranking-row strong {
  color: #e3ebf8;
  font-size: 11px;
  font-weight: 670;
  white-space: nowrap;
}

.dash-ranking-orders {
  color: #7188ad;
  white-space: nowrap;
}

.dash-structure-list {
  display: grid;
  gap: 14px;
  margin-top: 22px;
}

.dash-structure-label {
  justify-content: space-between;
  gap: 10px;
  color: #a8b9d2;
  font-size: 11px;
}

.dash-structure-label strong {
  color: #e2ebf8;
  font-size: 12px;
}

.dash-structure-track {
  height: 8px;
  margin-top: 7px;
}

.dash-structure-track i.structure-fill--blue {
  background: #2f9bff;
}

.dash-structure-track i.structure-fill--purple {
  background: #7b61ff;
}

.dash-structure-track i.structure-fill--orange {
  background: #eea23b;
}

.dash-structure-track i.structure-fill--red {
  background: #e4617e;
}

.dash-structure-summary {
  position: relative;
  z-index: 1;
  display: grid;
  grid-template-columns: auto 1fr auto 1fr;
  gap: 5px 10px;
  align-items: baseline;
  margin-top: 20px;
  padding-top: 12px;
  border-top: 1px solid rgba(54, 111, 186, 0.22);
}

.dash-structure-summary strong {
  color: #f1f5ff;
  font-size: 13px;
  font-weight: 720;
  white-space: nowrap;
}

.dash-panel-heading-icon {
  color: #8c76f1;
  font-size: 20px;
}

.dash-footer-grid {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  gap: 10px;
  min-height: 0;
}

.dash-footer-card {
  display: flex;
  align-items: flex-start;
  min-width: 0;
  min-height: 0;
  gap: 10px;
  padding: 16px 14px;
  background: rgba(9, 28, 59, 0.76);
}

.dash-footer-card::before {
  position: absolute;
  inset: 0;
  pointer-events: none;
  content: "";
  border-radius: inherit;
  background: linear-gradient(135deg, rgba(78, 89, 217, 0.12), transparent 50%);
}

.dash-footer-card--purple {
  border-color: rgba(121, 85, 232, 0.42);
}

.dash-footer-card--blue {
  border-color: rgba(35, 116, 222, 0.4);
}

.dash-footer-card--orange {
  border-color: rgba(212, 130, 52, 0.38);
}

.dash-footer-card--green {
  border-color: rgba(29, 164, 130, 0.38);
}

.dash-footer-card--red {
  border-color: rgba(211, 74, 114, 0.4);
}

.dash-footer-icon {
  width: 34px;
  height: 34px;
  color: #8d7aff;
  background: rgba(92, 72, 210, 0.25);
}

.dash-footer-card--blue .dash-footer-icon {
  color: #5ba9ff;
  background: rgba(29, 111, 221, 0.25);
}

.dash-footer-card--orange .dash-footer-icon {
  color: #f3a143;
  background: rgba(187, 97, 29, 0.25);
}

.dash-footer-card--green .dash-footer-icon {
  color: #37cba1;
  background: rgba(19, 143, 105, 0.24);
}

.dash-footer-card--red .dash-footer-icon {
  color: #f07191;
  background: rgba(188, 53, 91, 0.24);
}

.dash-footer-card>div {
  position: relative;
  z-index: 1;
  display: grid;
  min-width: 0;
}

.dash-footer-card>div>span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dash-footer-card strong {
  margin-top: 8px;
  color: #f1f5ff;
  font-size: clamp(16px, 1.15vw, 23px);
  font-weight: 720;
  line-height: 1;
  white-space: nowrap;
}

.dash-empty {
  position: relative;
  z-index: 1;
  margin: 28px 0 0;
  color: #7389ad;
  font-size: 12px;
}

.dash-error {
  position: relative;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin: 0;
  padding: 9px 12px;
  border: 1px solid rgba(230, 98, 119, 0.38);
  border-radius: 8px;
  color: #ffb8c5;
  background: rgba(94, 27, 48, 0.35);
  font-size: 12px;
}

.dash-error button {
  padding: 0;
  color: #fff;
  background: transparent;
  font-size: 12px;
  font-weight: 700;
  text-decoration: underline;
}

.dash-loading-overlay {
  position: absolute;
  z-index: 10;
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
  position: absolute;
  inset: 18%;
  content: "";
  border: 1px solid rgba(34, 211, 238, 0.26);
  border-radius: 999px;
  background: rgba(15, 23, 42, 0.72);
  box-shadow: inset 0 0 28px rgba(34, 211, 238, 0.08), 0 0 38px rgba(34, 211, 238, 0.14);
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
  border: 1px solid rgba(34, 211, 238, 0.18);
  border-top-color: rgba(34, 211, 238, 0.95);
  border-right-color: rgba(167, 139, 250, 0.78);
  border-radius: 999px;
  animation: dash-loading-spin 1.35s linear infinite;
}

.dash-loading-orbit--second {
  inset: 14%;
  border-top-color: rgba(52, 211, 153, 0.9);
  border-right-color: rgba(34, 211, 238, 0.46);
  animation-direction: reverse;
  animation-duration: 1.9s;
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

.dash-icon-spin {
  animation: dash-icon-spin 1s linear infinite;
}

:deep(.count-up-number) {
  color: inherit;
  font: inherit;
  line-height: inherit;
}

@media (max-width: 1500px) {
  .dash-hero-panel {
    grid-template-columns: minmax(270px, 0.9fr) minmax(0, 2.1fr);
  }

  .dash-metric-card {
    padding-left: 10px;
    padding-right: 8px;
  }

  .dash-metric-icon {
    width: 31px;
    height: 31px;
    font-size: 16px;
  }

  .dash-footer-grid {
    gap: 10px;
  }

  .dash-footer-card {
    gap: 7px;
    padding-left: 10px;
    padding-right: 8px;
  }
}

@media (max-width: 1220px) {
  .dash-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .dash-toolbar {
    justify-content: flex-start;
    width: 100%;
  }

  .dash-hero-panel {
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .dash-hero-summary {
    min-height: 126px;
    padding-right: 0;
  }

  .dash-hero-chart {
    width: 42%;
  }

  .dash-metric-grid {
    grid-template-columns: repeat(5, minmax(120px, 1fr));
    overflow: hidden;
    padding-bottom: 2px;
  }

  .dash-middle-grid {
    grid-template-columns: 1fr;
  }

  .dash-lower-grid {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .dash-footer-grid {
    grid-template-columns: repeat(4, minmax(0, 1fr));
  }
}

@media (max-width: 820px) {
  .dash-screen {
    min-height: 100dvh;
    overflow: visible;
  }

  .dash-screen__content {
    height: auto;
    min-height: 100dvh;
    padding: 12px;
    overflow: visible;
  }

  .dash-toolbar {
    gap: 6px;
  }

  .dash-periods {
    width: 100%;
    overflow-x: auto;
  }

  .dash-period {
    flex: 1 0 auto;
  }

  .dash-range-control {
    flex: 1 1 180px;
  }

  .dash-middle-grid {
    grid-template-columns: 1fr;
  }

  .dash-lower-grid {
    grid-template-columns: 1fr;
  }

  .dash-footer-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dash-chart-panel,
  .dash-insight-panel,
  .dash-data-panel {
    min-height: 280px;
  }
}

@media (max-width: 520px) {
  .dash-heading h1 {
    font-size: 22px;
  }

  .dash-live-status {
    display: none;
  }

  .dash-range-control {
    order: 3;
    flex-basis: 100%;
  }

  .dash-export-button {
    flex: 1 1 auto;
  }

  .dash-hero-panel {
    padding: 14px;
  }

  .dash-hero-value {
    font-size: 38px;
  }

  .dash-metric-grid {
    grid-template-columns: repeat(2, minmax(150px, 1fr));
    overflow: visible;
  }

  .dash-metric-card:last-child {
    grid-column: 1 / -1;
  }

  .dash-channel-layout {
    grid-template-columns: 1fr;
    height: auto;
  }

  .dash-pie-chart {
    height: 190px;
  }

  .dash-channel-list {
    padding: 0 6px 10px;
  }

  .dash-footer-grid {
    grid-template-columns: 1fr;
  }
}

@media (prefers-reduced-motion: reduce) {

  .dash-screen *,
  .dash-screen *::before,
  .dash-screen *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    scroll-behavior: auto !important;
    transition-duration: 0.01ms !important;
  }
}

@keyframes dash-icon-spin {
  to {
    transform: rotate(360deg);
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
