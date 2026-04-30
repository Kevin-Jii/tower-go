<template>
  <div class="h-full overflow-auto p-4 md:p-6 space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h2 class="page-title mb-1">已同步订单</h2>
        <p class="text-sm text-gray-500">{{ accountName ? `账号：${accountName}` : `账号ID：${accountId}` }}</p>
      </div>
      <div class="flex items-center gap-2">
        <BaseButton variant="ghost" @click="goBack">返回账号池</BaseButton>
        <BaseButton variant="primary" :loading="loading" @click="load">刷新</BaseButton>
      </div>
    </div>

    <div v-if="account" class="rounded-lg border border-gray-200 bg-white p-3 md:p-4">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-2 text-sm">
        <div><span class="text-gray-500">账号名称：</span><span class="text-gray-800">{{ account.name || '-' }}</span></div>
        <div><span class="text-gray-500">shopId：</span><span class="text-gray-800">{{ account.shop_id || '-' }}</span></div>
        <div><span class="text-gray-500">customerId：</span><span class="text-gray-800">{{ account.customer_id || '-' }}</span></div>
        <div><span class="text-gray-500">最后同步：</span><span class="text-gray-800">{{ formatDateTime(account.last_sync_at) }}</span></div>
        <div><span class="text-gray-500">最近同步条数：</span><span class="text-gray-800">{{ account.last_sync_count ?? 0 }}</span></div>
        <div><span class="text-gray-500">同步备注：</span><span class="text-gray-800">{{ account.last_sync_msg || '-' }}</span></div>
      </div>
    </div>

    <div v-if="!loading && groupedRows.length === 0" class="text-sm text-gray-500">暂无已同步订单</div>

    <div
      v-for="g in groupedRows"
      :key="g.date"
      :ref="(el) => setGroupRef(g.date, el)"
      class="rounded-lg border border-gray-200 bg-white p-3 md:p-4 space-y-2"
    >
      <div class="flex items-center justify-between gap-3">
        <div class="font-medium text-gray-800">
          {{ formatDateLabel(g.date) }}（{{ g.rows.length }}单）
          <span class="ml-3 text-sm font-normal text-red-500">
            支付合计：{{ formatMoney(g.payTotal) }} / 订单合计：{{ formatMoney(g.orderTotal) }}
          </span>
        </div>
        <div class="flex items-center gap-2" data-export-exclude="1">
          <BaseButton variant="ghost" size="sm" @click="toggleGroup(g.date)">
            {{ isCollapsed(g.date) ? '展开' : '折叠' }}
          </BaseButton>
          <BaseButton variant="ghost" size="sm" :loading="exportingDate === g.date" @click="exportGroupAsImage(g.date)">
            导出图片
          </BaseButton>
        </div>
      </div>
      <div v-show="!isCollapsed(g.date)" class="overflow-x-auto">
        <BaseTable :columns="columns" :data="(g.rows as unknown) as Record<string, unknown>[]" :loading="loading" min-width="980px">
          <template #cell-order_items="{ row }">
            <div class="max-w-[420px] truncate" :title="formatOrderItems(row as ThirdPartyOrder)">
              {{ formatOrderItems(row as ThirdPartyOrder) || '-' }}
            </div>
          </template>
          <template #cell-pay_amount="{ row }">
            <span class="text-red-500">{{ formatMoney((row as ThirdPartyOrder).pay_amount) }}</span>
          </template>
          <template #cell-total_amount="{ row }">
            <span class="text-red-500">{{ formatMoney((row as ThirdPartyOrder).total_amount) }}</span>
          </template>
          <template #cell-synced_at="{ row }">
            {{ formatDateTime((row as ThirdPartyOrder).synced_at) }}
          </template>
        </BaseTable>
      </div>
      <div v-show="isCollapsed(g.date)" class="text-xs text-gray-400">
        当前分组已折叠
      </div>
    </div>

    <div class="flex items-center justify-between text-sm text-gray-600">
      <span>共 {{ total }} 条</span>
      <div class="flex items-center gap-2">
        <BaseButton variant="ghost" :disabled="page <= 1 || loading" @click="changePage(page - 1)">上一页</BaseButton>
        <span>第 {{ page }} / {{ totalPages }} 页</span>
        <BaseButton variant="ghost" :disabled="page >= totalPages || loading" @click="changePage(page + 1)">下一页</BaseButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import type { ComponentPublicInstance } from 'vue'
import html2canvas from 'html2canvas'
import { useRoute, useRouter } from 'vue-router'
import { BaseButton, BaseTable } from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import type { ThirdPartyAccount, ThirdPartyOrder } from '@/api/types'
import { getThirdPartyAccount, listThirdPartySyncedOrders } from '@/api/thirdPartyAccount'
import { toast } from '@/feedback/toast'

const route = useRoute()
const router = useRouter()
const accountId = Number(route.params.id || 0)
const accountName = String(route.query.name || '')

const loading = ref(false)
const rows = ref<ThirdPartyOrder[]>([])
const account = ref<ThirdPartyAccount | null>(null)
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const totalPages = computed(() => Math.max(1, Math.ceil(total.value / pageSize.value)))
const exportingDate = ref('')
const groupRefs = ref<Record<string, HTMLElement>>({})
const collapsedMap = ref<Record<string, boolean>>({})

const columns: BaseTableColumn[] = [
  { key: 'order_no', label: '订单号', prop: 'order_no', minWidth: '200px', ellipsis: true },
  { key: 'status_name', label: '状态', prop: 'status_name', width: '120px' },
  { key: 'order_items', label: '行项目明细', minWidth: '360px' },
  { key: 'total_item_num', label: '总件数', prop: 'total_item_num', width: '90px' },
  { key: 'pay_amount', label: '支付金额', prop: 'pay_amount', width: '110px' },
  { key: 'total_amount', label: '订单金额', prop: 'total_amount', width: '110px' },
  { key: 'synced_at', label: '同步时间', prop: 'synced_at', minWidth: '170px' },
]

const groupedRows = computed(() => {
  const map = new Map<string, ThirdPartyOrder[]>()
  for (const r of rows.value) {
    const d = String(r.place_date || '').trim() || '未知日期'
    if (!map.has(d)) map.set(d, [])
    map.get(d)!.push(r)
  }
  return Array.from(map.entries())
    .sort((a, b) => (a[0] > b[0] ? -1 : 1))
    .map(([date, rs]) => ({
      date,
      rows: rs,
      payTotal: rs.reduce((sum, item) => sum + toNumber(item.pay_amount), 0),
      orderTotal: rs.reduce((sum, item) => sum + toNumber(item.total_amount), 0),
    }))
})

function formatDateLabel(d: string): string {
  if (!d || d === '未知日期') return d
  const dt = new Date(`${d}T00:00:00`)
  if (Number.isNaN(dt.getTime())) return d
  return dt.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', weekday: 'short' })
}

function formatDateTime(value?: string): string {
  const text = String(value || '').trim()
  if (!text) return '-'
  const dt = new Date(text)
  if (Number.isNaN(dt.getTime())) return text
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${dt.getFullYear()}-${pad(dt.getMonth() + 1)}-${pad(dt.getDate())} ${pad(dt.getHours())}:${pad(dt.getMinutes())}:${pad(dt.getSeconds())}`
}

function toNumber(v: unknown): number {
  const n = Number(v ?? 0)
  return Number.isFinite(n) ? n : 0
}

function formatMoney(v: unknown): string {
  return toNumber(v).toFixed(2)
}

function setGroupRef(date: string, el: Element | ComponentPublicInstance | null): void {
  if (!el) return
  if (el instanceof Element) {
    groupRefs.value[date] = el as HTMLElement
    return
  }
  const root = (el as ComponentPublicInstance).$el
  if (root instanceof Element) {
    groupRefs.value[date] = root as HTMLElement
  }
}

function isCollapsed(date: string): boolean {
  return !!collapsedMap.value[date]
}

function toggleGroup(date: string): void {
  collapsedMap.value[date] = !isCollapsed(date)
}

async function exportGroupAsImage(date: string): Promise<void> {
  const target = groupRefs.value[date]
  if (!target) {
    toast.warning('未找到可导出的表格区域')
    return
  }
  exportingDate.value = date
  let exportRoot: HTMLDivElement | null = null
  try {
    exportRoot = document.createElement('div')
    exportRoot.style.position = 'fixed'
    exportRoot.style.left = '-99999px'
    exportRoot.style.top = '0'
    exportRoot.style.background = '#fff'
    exportRoot.style.padding = '20px'
    exportRoot.style.width = `${Math.max(target.scrollWidth, target.clientWidth)}px`

    const title = document.createElement('h1')
    title.style.margin = '0 0 16px 0'
    title.style.fontSize = '28px'
    title.style.lineHeight = '1.2'
    title.style.fontWeight = '700'
    title.style.color = '#111827'
    title.textContent = account.value?.name || accountName || '门店'

    const cloned = target.cloneNode(true) as HTMLElement
    cloned.querySelectorAll('[data-export-exclude="1"]').forEach((node) => node.remove())
    exportRoot.appendChild(title)
    exportRoot.appendChild(cloned)
    document.body.appendChild(exportRoot)

    const canvas = await html2canvas(exportRoot, {
      scale: 2,
      backgroundColor: '#ffffff',
      useCORS: true,
    })
    const a = document.createElement('a')
    a.href = canvas.toDataURL('image/png')
    a.download = `synced-orders-${date || 'unknown'}.png`
    a.click()
    toast.success('已导出图片')
  } catch {
    toast.error('导出失败，请重试')
  } finally {
    if (exportRoot && document.body.contains(exportRoot)) {
      document.body.removeChild(exportRoot)
    }
    exportingDate.value = ''
  }
}

function formatOrderItems(row: ThirdPartyOrder): string {
  const raw = String(row.raw_json || '').trim()
  if (!raw) return ''
  try {
    const parsed = JSON.parse(raw) as {
      itemList?: Array<{ itemName?: string; skuName?: string; itemNum?: number }>
      rowItemTypeInfoList?: Array<{ wrapperUnit?: string }>
    }
    const unit = parsed.rowItemTypeInfoList?.[0]?.wrapperUnit || ''
    const items = parsed.itemList || []
    if (!items.length) return ''
    return items
      .map((i) => `${i.itemName || i.skuName || '未知商品'} x${Number(i.itemNum || 0)}${unit}`)
      .join('；')
  } catch {
    return ''
  }
}

async function load(): Promise<void> {
  if (!accountId) return
  loading.value = true
  try {
    if (!account.value) {
      account.value = await getThirdPartyAccount(accountId)
    }
    const res = await listThirdPartySyncedOrders(accountId, page.value, pageSize.value)
    rows.value = res.list || []
    total.value = res.total || 0
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载订单失败')
  } finally {
    loading.value = false
  }
}

async function changePage(p: number): Promise<void> {
  const target = Math.min(Math.max(1, p), totalPages.value)
  if (target === page.value) return
  page.value = target
  await load()
}

function goBack(): void {
  void router.push({ name: 'ThirdPartyAccount' })
}

void load()
</script>
