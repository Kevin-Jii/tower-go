<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">导入第三方订单</h2>
      <div class="flex gap-2">
        <BaseButton variant="ghost" @click="goBack">返回物流路线图</BaseButton>
        <BaseButton variant="primary" :disabled="hotData.length === 0" :loading="saving" @click="saveSheet">保存物流单</BaseButton>
        <BaseButton variant="primary" :disabled="hotData.length === 0" @click="printSheet">打印</BaseButton>
      </div>
    </div>

    <div class="rounded border border-[var(--color-border-2)] bg-white p-3">
      <div class="flex flex-wrap items-center gap-2">
        <BaseSelect v-model="selectedRouteId" :options="routeOptions" placeholder="请选择路线" />
        <BaseInput v-model="startDate" type="date" class="w-44" />
        <span class="text-gray-400">至</span>
        <BaseInput v-model="endDate" type="date" class="w-44" />
        <BaseButton variant="primary" :loading="loading" @click="loadData">导入并生成表格</BaseButton>
      </div>
    </div>

    <div class="rounded border border-[var(--color-border-2)] bg-white p-2">
      <HotTable ref="hotRef" :settings="hotSettings" :theme-name="'ht-theme-main'" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { HotTable } from '@handsontable/vue3'
import { registerAllModules } from 'handsontable/registry'
import 'handsontable/styles/handsontable.min.css'
import 'handsontable/styles/ht-theme-main.min.css'
import { BaseButton, BaseInput, BaseSelect } from '@/components/base'
import type { BaseSelectOption } from '@/components/base/types'
import { listThirdPartyRoutes, importThirdPartyRouteByDateRange, saveThirdPartyLogisticsSheet } from '@/api/thirdPartyRoute'
import type { ThirdPartyRoute } from '@/api/types'
import { toast } from '@/feedback/toast'

registerAllModules()

const router = useRouter()
const route = useRoute()
const { data: routesData } = useQuery({ queryKey: ['third-party-routes'], queryFn: listThirdPartyRoutes })
const routeRows = computed(() => routesData.value ?? [])
const routeOptions = computed<BaseSelectOption[]>(() => routeRows.value.map((r) => ({ label: r.name, value: r.id })))

const selectedRouteId = ref<number | string | undefined>(route.query.route_id ? Number(route.query.route_id) : undefined)
const startDate = ref(new Date().toISOString().slice(0, 10))
const endDate = ref(new Date().toISOString().slice(0, 10))
const loading = ref(false)
const saving = ref(false)

const hotRef = ref<InstanceType<typeof HotTable> | null>(null)
const hotHeaders = ref<string[]>([])
const hotData = ref<Array<Array<string | number>>>([])
const hotSettings = computed(() => ({
  data: hotData.value,
  colHeaders: hotHeaders.value,
  rowHeaders: true,
  height: 560,
  licenseKey: 'non-commercial-and-evaluation',
  contextMenu: true,
  manualColumnResize: true,
  manualRowResize: true,
  stretchH: 'all' as const,
  cells: cellMeta,
  afterChange: handleAfterChange,
}))

const currentRoute = computed<ThirdPartyRoute | undefined>(() => {
  const id = Number(selectedRouteId.value || 0)
  return routeRows.value.find((x) => x.id === id)
})

function goBack(): void {
  void router.push({ name: 'ThirdPartyRoutes' })
}

async function loadData(): Promise<void> {
  const routeId = Number(selectedRouteId.value || 0)
  if (!routeId) {
    toast.warning('请先选择路线')
    return
  }
  loading.value = true
  try {
    const res = await importThirdPartyRouteByDateRange(routeId, startDate.value, endDate.value)
    const stores = currentRoute.value?.stores ?? []
    const headers = ['品类', '大总', ...stores.map((s) => s.store?.name || `门店#${s.store_id}`)]
    const rows = (res.list || []).map((item) => {
      const map = new Map((item.store_qty || []).map((x) => [x.store_id, x.quantity]))
      const perStore = stores.map((s) => Number(map.get(s.store_id) || 0))
      const total = perStore.reduce((sum, n) => sum + n, 0)
      return [item.product_name, total, ...perStore]
    })
    hotHeaders.value = headers
    hotData.value = rows
    await nextTick()
    hotRef.value?.hotInstance?.loadData(rows)
    hotRef.value?.hotInstance?.updateSettings({ colHeaders: headers })
    hotRef.value?.hotInstance?.render()
    toast.success(`已生成表格，共 ${rows.length} 行商品`)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '导入失败')
  } finally {
    loading.value = false
  }
}

function toNum(v: unknown): number {
  const n = Number(v)
  return Number.isFinite(n) ? n : 0
}

function cellMeta(_row: number, col: number): Record<string, unknown> {
  return {
    readOnly: col <= 1,
  }
}

function handleAfterChange(changes: Array<[number, string | number, unknown, unknown]> | null, source: string): void {
  if (!changes || source === 'loadData' || source === 'recalc') return
  const hot = hotRef.value?.hotInstance
  if (!hot) return
  const rowsToRecalc = new Set<number>()
  for (const ch of changes) {
    const row = Number(ch[0])
    const col = typeof ch[1] === 'number' ? Number(ch[1]) : Number(hot.propToCol(ch[1]))
    if (col >= 2) rowsToRecalc.add(row)
  }
  rowsToRecalc.forEach((row) => {
    let total = 0
    for (let col = 2; col < hotHeaders.value.length; col += 1) {
      total += toNum(hot.getDataAtCell(row, col))
    }
    hot.setDataAtCell(row, 1, Number(total.toFixed(2)), 'recalc')
  })
}

function printSheet(): void {
  if (!hotData.value.length) {
    toast.warning('暂无可打印表格')
    return
  }
  const body = hotData.value
    .map((row) => `<tr>${row.map((x) => `<td>${x}</td>`).join('')}</tr>`)
    .join('')
  const html = `<!doctype html><html><head><meta charset="utf-8"><title>导入商品打印</title><style>
      body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif;padding:16px;}
      h2{margin:0 0 8px;}
      table{border-collapse:collapse;width:100%;}
      th,td{border:1px solid #222;padding:6px 8px;text-align:center;font-size:12px;}
      th:first-child,td:first-child{text-align:left;}
    </style></head><body>
      <h2>${currentRoute.value?.name || '物流路线图'}（${startDate.value} ~ ${endDate.value}）</h2>
      <table><thead><tr>${hotHeaders.value.map((h) => `<th>${h}</th>`).join('')}</tr></thead><tbody>${body}</tbody></table>
    </body></html>`
  const win = window.open('', '_blank')
  if (!win) {
    toast.error('浏览器拦截了打印窗口')
    return
  }
  win.document.write(html)
  win.document.close()
  win.focus()
  win.print()
}

async function saveSheet(): Promise<void> {
  const routeId = Number(selectedRouteId.value || 0)
  if (!routeId || !hotData.value.length) {
    toast.warning('暂无可保存数据')
    return
  }
  saving.value = true
  try {
    const headers = [...hotHeaders.value]
    const products = hotData.value.map((r) => String(r[0] ?? ''))
    const rows = hotData.value.map((r) => r.slice(1).map((x) => toNum(x)))
    await saveThirdPartyLogisticsSheet(routeId, {
      start_date: startDate.value,
      end_date: endDate.value,
      headers,
      rows,
      products,
    })
    toast.success('物流单已保存（同一路线与相同导入日期区间将覆盖原记录）')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}
</script>
