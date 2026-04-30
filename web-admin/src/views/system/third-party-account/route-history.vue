<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">历史物流单</h2>
      <div class="flex gap-2">
        <BaseButton variant="ghost" @click="goBack">返回物流路线图</BaseButton>
      </div>
    </div>

    <div class="rounded border border-[var(--color-border-2)] bg-white p-3">
      <div class="flex flex-wrap items-center gap-2">
        <BaseSelect v-model="selectedRouteId" :options="routeOptions" placeholder="请选择路线" />
        <a-range-picker
          v-model="filterDateRange"
          value-format="YYYY-MM-DD"
          format="YYYY-MM-DD"
          class="route-history-range"
          style="width: 280px"
          allow-clear
          :placeholder="['开始日期', '结束日期']"
        />
        <BaseButton variant="primary" @click="loadHistory">查询</BaseButton>
        <BaseButton variant="ghost" @click="resetDateFilter">清空日期</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(filteredSheetRows as unknown) as Record<string, unknown>[]" :loading="loading" min-width="960px">
      <template #cell-updated_at="{ row }">
        {{ formatDateTime((row as ThirdPartyLogisticsSheet).updated_at) }}
      </template>
      <template #cell-actions="{ row }">
        <BaseButton variant="link" size="sm" @click="viewSheet(row as ThirdPartyLogisticsSheet)">查看</BaseButton>
      </template>
    </BaseTable>

    <a-drawer
      :visible="viewDrawer"
      placement="right"
      :width="1180"
      :drawer-style="{ maxWidth: '96vw' }"
      :mask-closable="true"
      unmount-on-close
      @cancel="viewDrawer = false"
    >
      <template #title>物流单详情（只读）</template>
      <div class="space-y-3">
        <div class="text-sm text-gray-600">
          保存日期：{{ currentSheet?.sheet_date || '-' }}，
          区间：{{ currentSheet?.start_date || '-' }} ~ {{ currentSheet?.end_date || '-' }}
        </div>
        <div class="overflow-x-auto">
          <table class="w-full min-w-[960px] border-collapse text-sm">
            <thead>
              <tr>
                <th v-for="h in currentHeaders" :key="h" class="border border-gray-300 bg-gray-50 px-2 py-2 text-center">{{ h }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(r, idx) in currentTableRows" :key="idx">
                <td v-for="(c, cIdx) in r" :key="`${idx}-${cIdx}`" class="border border-gray-300 px-2 py-1 text-center">
                  {{ c }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end">
          <BaseButton variant="ghost" @click="viewDrawer = false">关闭</BaseButton>
        </div>
      </template>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuery } from '@tanstack/vue-query'
import { BaseButton, BaseSelect, BaseTable } from '@/components/base'
import type { BaseSelectOption, BaseTableColumn } from '@/components/base/types'
import { listThirdPartyLogisticsSheets, listThirdPartyRoutes } from '@/api/thirdPartyRoute'
import type { ThirdPartyLogisticsSheet } from '@/api/types'
import { toast } from '@/feedback/toast'

const router = useRouter()
const route = useRoute()
const { data: routesData } = useQuery({ queryKey: ['third-party-routes'], queryFn: listThirdPartyRoutes })
const routeRows = computed(() => routesData.value ?? [])
const routeOptions = computed<BaseSelectOption[]>(() => routeRows.value.map((r) => ({ label: r.name, value: r.id })))
const selectedRouteId = ref<number | string | undefined>(route.query.route_id ? Number(route.query.route_id) : undefined)
const loading = ref(false)
const sheetRows = ref<ThirdPartyLogisticsSheet[]>([])
/** Arco RangePicker：value-format 为 YYYY-MM-DD 时为 string[] */
const filterDateRange = ref<string[]>([])

const filterStartDate = computed(() => filterDateRange.value?.[0] ?? '')
const filterEndDate = computed(() => filterDateRange.value?.[1] ?? '')

const columns: BaseTableColumn[] = [
  { key: 'sheet_date', label: '保存日期', prop: 'sheet_date', width: '120px' },
  { key: 'start_date', label: '开始日期', prop: 'start_date', width: '120px' },
  { key: 'end_date', label: '结束日期', prop: 'end_date', width: '120px' },
  { key: 'updated_at', label: '更新时间', prop: 'updated_at', minWidth: '200px' },
  { key: 'actions', label: '操作', width: '90px', align: 'right' },
]

const viewDrawer = ref(false)
const currentSheet = ref<ThirdPartyLogisticsSheet | null>(null)
const currentHeaders = computed(() => currentSheet.value?.headers || [])
const currentTableRows = computed(() => {
  if (!currentSheet.value) return []
  const products = currentSheet.value.products || []
  const rows = currentSheet.value.rows || []
  return rows.map((r, i) => [products[i] || '-', ...r])
})
const filteredSheetRows = computed(() => {
  return sheetRows.value.filter((x) => {
    if (filterStartDate.value && x.sheet_date < filterStartDate.value) return false
    if (filterEndDate.value && x.sheet_date > filterEndDate.value) return false
    return true
  })
})

function goBack(): void {
  void router.push({ name: 'ThirdPartyRoutes' })
}

async function loadHistory(): Promise<void> {
  const routeId = Number(selectedRouteId.value || 0)
  if (!routeId) {
    toast.warning('请先选择路线')
    return
  }
  loading.value = true
  try {
    sheetRows.value = await listThirdPartyLogisticsSheets(routeId)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载失败')
  } finally {
    loading.value = false
  }
}

function formatDateTime(v: string | undefined): string {
  const s = String(v || '').trim()
  if (!s) return '-'
  const d = new Date(s)
  if (Number.isNaN(d.getTime())) {
    return s.length >= 19 ? s.slice(0, 19).replace('T', ' ') : s
  }
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function viewSheet(row: ThirdPartyLogisticsSheet): void {
  currentSheet.value = row
  viewDrawer.value = true
}

function resetDateFilter(): void {
  filterDateRange.value = []
}

if (selectedRouteId.value) {
  void loadHistory()
}
</script>

<style scoped>
.route-history-range :deep(.arco-picker) {
  width: 100%;
}
</style>
