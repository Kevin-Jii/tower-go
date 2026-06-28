<template>
  <div class="flex min-h-0 flex-1 flex-col gap-4">
    <BaseCard flush-body class="flex min-h-0 flex-1 flex-col">
      <template #header>
        <div class="flex w-full flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h2 class="m-0 text-base font-semibold text-slate-900">报损 / 自用 / 赠送</h2>
            <p class="m-0 mt-1 text-xs text-slate-500">只扣库存、计入成本，不计入销售收入</p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <BaseInput v-model="filters.keyword" class="w-44" placeholder="单号 / 商品 / 原因" clearable @enter="reload" />
            <BaseSelect v-model="filters.type" class="w-32" :options="typeFilterOptions" />
            <BaseInput v-model="filters.start_date" class="w-36" type="date" />
            <BaseInput v-model="filters.end_date" class="w-36" type="date" />
            <BaseButton variant="primary" @click="reload">查询</BaseButton>
            <BaseButton variant="secondary" @click="openExportDlg">导出Excel</BaseButton>
            <BaseButton variant="primary" @click="openCreate">新增单据</BaseButton>
          </div>
        </div>
      </template>

      <div class="flex min-h-0 flex-1 flex-col gap-3 p-4">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
          <div v-for="item in summaryCards" :key="item.label"
            class="rounded border border-slate-200 bg-slate-50 px-4 py-3">
            <div class="text-xs font-medium text-slate-500">{{ item.label }}</div>
            <div class="mt-1 text-lg font-semibold text-slate-900">{{ item.value }}</div>
          </div>
        </div>

        <BaseTable :columns="columns" :data="(rows as unknown) as Record<string, unknown>[]" :loading="loading"
          min-width="1080px" class="min-h-0 flex-1">
          <template #cell-type="{ row }">
            <span :class="typeClass((row as InventoryLossOrder).type)">
              {{ typeLabel((row as InventoryLossOrder).type) }}
            </span>
          </template>
          <template #cell-member="{ row }">
            {{ memberLabel(row as InventoryLossOrder) }}
          </template>
          <template #cell-total_cost="{ row }">
            {{ formatMoney((row as InventoryLossOrder).total_cost) }}
          </template>
          <template #cell-created_at="{ row }">
            {{ formatDateTime((row as InventoryLossOrder).created_at) }}
          </template>
          <template #cell-actions="{ row }">
            <div class="flex justify-end gap-3" @click.stop>
              <BaseButton variant="link" size="sm" @click="openDetail(row as InventoryLossOrder)">详情</BaseButton>
              <BaseButton v-if="!(row as InventoryLossOrder).is_canceled" variant="link" size="sm"
                @click="openEdit(row as InventoryLossOrder)">编辑</BaseButton>
              <BaseButton variant="link" size="sm" @click="onCancel(row as InventoryLossOrder)">撤销</BaseButton>
            </div>
          </template>
        </BaseTable>

        <div class="flex shrink-0 justify-end">
          <BasePagination :page="page" :page-size="pageSize" :total="total" @update:page="(p) => (page = p)"
            @update:page-size="(s) => (pageSize = s)" />
        </div>
      </div>
    </BaseCard>

    <BaseDialog v-model="createDlg" title="新增报损 / 自用 / 赠送" max-width="min(1080px, 96vw)">
      <div class="space-y-4">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
          <BaseFormItem label="类型" required>
            <BaseSelect v-model="form.type" :options="typeOptions" />
          </BaseFormItem>
          <BaseFormItem v-if="form.type === 'gift'" label="赠送会员" required>
            <BaseSelect v-model="form.member_id" :options="memberOptions" placeholder="请选择会员" />
          </BaseFormItem>
          <BaseFormItem label="原因说明" :class="form.type === 'gift' ? '' : 'md:col-span-2'" required>
            <BaseSelect v-model="form.reason" :options="reasonOptions" placeholder="请选择原因说明" searchable />
          </BaseFormItem>
        </div>

        <div class="rounded border border-[var(--color-border-2)] p-3">
          <div class="mb-3 flex items-center justify-between">
            <h3 class="m-0 text-sm font-semibold text-slate-900">商品明细</h3>
            <BaseButton variant="secondary" size="sm" @click="addLine">添加商品</BaseButton>
          </div>
          <div class="loss-line-editor">
            <div class="loss-line-editor__head">
              <span>商品</span>
              <span>规格</span>
              <span>数量</span>
              <span>备注</span>
              <span>操作</span>
            </div>
            <div v-for="(line, idx) in form.lines" :key="idx" class="loss-line-editor__row">
              <BaseSelect v-model="line.product_id" class="loss-line-editor__control" :options="productOptions"
                placeholder="选择商品" @update:model-value="onLineProductChange(idx)" />
              <BaseSelect v-model="line.unit" class="loss-line-editor__control" :options="lineUnitOptions(line)"
                placeholder="选择规格" />
              <BaseNumberInput v-model="line.quantity" class="loss-line-editor__control" :min="0.01" :step="1" />
              <BaseInput v-model="line.remark" class="loss-line-editor__control" placeholder="可选" />
              <BaseButton variant="ghost" size="sm" :disabled="form.lines.length <= 1" @click="removeLine(idx)">移除
              </BaseButton>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="createDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitCreate">提交</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="editDlg" title="编辑原因说明" max-width="min(460px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="原因说明" required>
          <BaseSelect v-model="editForm.reason" :options="reasonOptions" placeholder="请选择原因说明" searchable />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="editDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="savingEdit" @click="submitEdit">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="detailDlg" title="单据详情" max-width="min(860px, 96vw)">
      <div v-if="detail" class="space-y-4">
        <div class="grid grid-cols-1 gap-3 rounded border border-[var(--color-border-2)] p-3 text-sm md:grid-cols-3">
          <div><span class="text-slate-500">单号：</span>{{ detail.order_no }}</div>
          <div><span class="text-slate-500">类型：</span>{{ typeLabel(detail.type) }}</div>
          <div><span class="text-slate-500">总成本：</span>{{ formatMoney(detail.total_cost) }}</div>
          <div><span class="text-slate-500">会员：</span>{{ memberLabel(detail) }}</div>
          <div><span class="text-slate-500">操作人：</span>{{ detail.operator_name || '-' }}</div>
          <div><span class="text-slate-500">时间：</span>{{ formatDateTime(detail.created_at) }}</div>
          <div class="md:col-span-3"><span class="text-slate-500">原因：</span>{{ detail.reason || '-' }}</div>
        </div>
        <BaseTable :columns="itemColumns" :data="(detail.items as unknown) as Record<string, unknown>[]"
          min-width="760px">
          <template #cell-cost_price="{ row }">{{ formatMoney((row as InventoryLossOrderItem).cost_price) }}</template>
          <template #cell-cost_amount="{ row }">{{ formatMoney((row as InventoryLossOrderItem).cost_amount)
            }}</template>
        </BaseTable>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="detailDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="exportDlg" title="导出报损 / 自用 / 赠送" max-width="min(420px, 96vw)">
      <BaseFormItem label="导出日期" required>
        <a-date-picker v-model="exportDate" value-format="YYYY-MM-DD" class="w-full" :allow-clear="false" />
      </BaseFormItem>
      <template #footer>
        <BaseButton variant="ghost" @click="exportDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="exporting" @click="submitExport">导出</BaseButton>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseCard,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseNumberInput,
  BasePagination,
  BaseSelect,
  BaseTable,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn } from '@/components/base/types'
import {
  cancelInventoryLossOrder,
  createInventoryLossOrder,
  exportInventoryLossOrders,
  getInventoryLossOrder,
  listInventoryLossOrders,
  updateInventoryLossOrder,
} from '@/api/inventoryLoss'
import { listDictDataByTypeCode } from '@/api/dict'
import { listMembers } from '@/api/member'
import { listPurchasableProducts } from '@/api/storeSupplier'
import { batchListProductUnitSpecs } from '@/api/supplierProduct'
import type {
  InventoryLossOrder,
  InventoryLossOrderDetail,
  InventoryLossOrderItem,
  InventoryLossType,
  MemberRow,
  ProductUnitSpec,
} from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'
import { useUserStore } from '@/store/user'

type LossLine = {
  product_id: number | ''
  unit: string
  quantity: number
  remark: string
}

const qc = useQueryClient()
const userStore = useUserStore()
const tenantStoreId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)

const typeOptions: Array<BaseSelectOption & { value: InventoryLossType }> = [
  { label: '报损', value: 'loss' },
  { label: '自用', value: 'self_use' },
  { label: '赠送', value: 'gift' },
]
const typeFilterOptions: BaseSelectOption[] = [{ label: '全部类型', value: '' }, ...typeOptions]

const { data: reasonDict } = useQuery({
  queryKey: ['dict-data', 'PERSONALUSE'] as const,
  queryFn: () => listDictDataByTypeCode('PERSONALUSE'),
})
const reasonOptions = computed<BaseSelectOption[]>(() =>
  (reasonDict.value ?? []).filter((item) => item.status === 1).map((item) => ({ label: item.label, value: item.value })),
)
const reasonValueByLabel = computed(() => {
  const map = new Map<string, string>()
  for (const item of reasonDict.value ?? []) map.set(item.label, item.value)
  return map
})

const filters = reactive({
  keyword: '',
  type: '' as InventoryLossType | '',
  start_date: '',
  end_date: '',
})
const page = ref(1)
const pageSize = ref(10)
const exportDlg = ref(false)
const exportDate = ref(new Date().toISOString().slice(0, 10))
const exporting = ref(false)
const queryKey = computed(
  () =>
    [
      'inventory-loss-orders',
      page.value,
      pageSize.value,
      filters.keyword.trim(),
      filters.type,
      filters.start_date,
      filters.end_date,
      tenantStoreId.value,
    ] as const,
)
const { data: pageData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: () =>
    listInventoryLossOrders({
      page: page.value,
      page_size: pageSize.value,
      store_id: tenantStoreId.value,
      keyword: filters.keyword.trim() || undefined,
      type: filters.type || undefined,
      start_date: filters.start_date || undefined,
      end_date: filters.end_date || undefined,
    }),
})
const rows = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)

const { data: productsData } = useQuery({
  queryKey: computed(() => ['inventory-loss-products', tenantStoreId.value] as const),
  queryFn: () => listPurchasableProducts({ store_id: tenantStoreId.value }),
})
const products = computed(() => productsData.value ?? [])
const productOptions = computed<BaseSelectOption[]>(() => products.value.map((p) => ({ label: p.name, value: p.id })))
const productIdsKey = computed(() => products.value.map((p) => p.id).sort((a, b) => a - b).join(','))
const { data: specsData } = useQuery({
  queryKey: computed(() => ['inventory-loss-product-specs', productIdsKey.value] as const),
  queryFn: async () => {
    const ids = products.value.map((p) => p.id)
    if (!ids.length) return [] as ProductUnitSpec[]
    return batchListProductUnitSpecs(ids)
  },
  enabled: computed(() => products.value.length > 0),
})
const specsByProduct = computed(() => {
  const map = new Map<number, ProductUnitSpec[]>()
  for (const s of specsData.value ?? []) {
    if (!s.is_enabled) continue
    if (!map.has(s.product_id)) map.set(s.product_id, [])
    map.get(s.product_id)!.push(s)
  }
  for (const [, specs] of map.entries()) {
    specs.sort((a, b) => Number(a.factor_to_base) - Number(b.factor_to_base))
  }
  return map
})

const { data: membersPage } = useQuery({
  queryKey: ['inventory-loss-members'],
  queryFn: () => listMembers({ page: 1, page_size: 300 }),
})
const members = computed(() => membersPage.value?.list ?? ([] as MemberRow[]))
const memberOptions = computed<BaseSelectOption[]>(() =>
  members.value.map((m) => ({
    label: `${m.phone}${m.name ? `（${m.name}）` : ''}`,
    value: m.id,
  })),
)

const columns: BaseTableColumn[] = [
  { key: 'order_no', label: '单号', prop: 'order_no', minWidth: '150px', ellipsis: true },
  { key: 'type', label: '类型', width: '90px' },
  { key: 'member', label: '会员', minWidth: '150px', ellipsis: true },
  { key: 'reason', label: '原因', prop: 'reason', minWidth: '180px', ellipsis: true },
  { key: 'item_count', label: '明细数', prop: 'item_count', width: '80px' },
  { key: 'total_cost', label: '总成本', width: '110px', align: 'right' },
  { key: 'operator_name', label: '操作人', prop: 'operator_name', width: '110px' },
  { key: 'created_at', label: '时间', width: '170px' },
  { key: 'actions', label: '操作', width: '128px', align: 'right' },
]

const itemColumns: BaseTableColumn[] = [
  { key: 'product_name', label: '商品', prop: 'product_name', minWidth: '180px', ellipsis: true },
  { key: 'unit', label: '规格', prop: 'unit', width: '120px' },
  { key: 'quantity', label: '数量', prop: 'quantity', width: '90px', align: 'right' },
  { key: 'base_quantity', label: '扣减库存', prop: 'base_quantity', width: '100px', align: 'right' },
  { key: 'base_unit', label: '库存单位', prop: 'base_unit', width: '100px' },
  { key: 'cost_price', label: '成本价', width: '100px', align: 'right' },
  { key: 'cost_amount', label: '成本金额', width: '110px', align: 'right' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '160px', ellipsis: true },
]

const summaryCards = computed(() => {
  const list = rows.value
  const sum = (type?: InventoryLossType) =>
    list
      .filter((x) => !type || x.type === type)
      .reduce((total, x) => total + Number(x.total_cost || 0), 0)
  return [
    { label: '当前页单据', value: `${list.length} 单` },
    { label: '报损成本', value: formatMoney(sum('loss')) },
    { label: '自用成本', value: formatMoney(sum('self_use')) },
    { label: '赠送成本', value: formatMoney(sum('gift')) },
  ]
})

function reload(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['inventory-loss-orders'] })
}

function openExportDlg(): void {
  exportDate.value = filters.start_date || filters.end_date || new Date().toISOString().slice(0, 10)
  exportDlg.value = true
}

async function submitExport(): Promise<void> {
  if (!exportDate.value) {
    toast.warning('请选择导出日期')
    return
  }
  exporting.value = true
  try {
    await exportInventoryLossOrders({ date: exportDate.value, store_id: tenantStoreId.value })
    exportDlg.value = false
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '导出失败')
  } finally {
    exporting.value = false
  }
}

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['inventory-loss-orders'] })
})

const createDlg = ref(false)
const saving = ref(false)
const form = reactive({
  type: 'loss' as InventoryLossType,
  member_id: '' as number | '',
  reason: '',
  lines: [] as LossLine[],
})

function makeLine(): LossLine {
  return { product_id: '', unit: '', quantity: 1, remark: '' }
}

function openCreate(): void {
  form.type = 'loss'
  form.member_id = ''
  form.reason = ''
  form.lines = [makeLine()]
  createDlg.value = true
}

function addLine(): void {
  form.lines.push(makeLine())
}

function removeLine(idx: number): void {
  form.lines = form.lines.filter((_, i) => i !== idx)
  if (!form.lines.length) form.lines.push(makeLine())
}

function lineUnitOptions(line: LossLine): BaseSelectOption[] {
  const pid = Number(line.product_id || 0)
  if (!pid) return []
  const specs = specsByProduct.value.get(pid) ?? []
  if (specs.length) {
    return specs.map((s) => ({
      label: specOptionLabel(s),
      value: String(s.unit_name || s.unit_code || '').trim(),
    }))
  }
  const product = products.value.find((p) => p.id === pid)
  const unit = product?.unit || '件'
  return [{ label: unit, value: unit }]
}

function onLineProductChange(idx: number): void {
  const line = form.lines[idx]
  if (!line) return
  const options = lineUnitOptions(line)
  line.unit = String(options[0]?.value || '')
}

function specOptionLabel(s: ProductUnitSpec): string {
  const name = s.unit_name || s.unit_code
  const factor = Number(s.factor_to_base || 0)
  const cost = Number(s.cost_price || 0)
  const parts = [name]
  if (factor > 0) parts.push(`扣减${formatQty(factor)}`)
  if (cost > 0) parts.push(`成本${formatMoney(cost)}`)
  return parts.join(' / ')
}

async function submitCreate(): Promise<void> {
  if (form.type === 'gift' && !Number(form.member_id || 0)) {
    toast.warning('赠送类型必须绑定会员')
    return
  }
  if (!form.reason.trim()) {
    toast.warning('请填写原因说明')
    return
  }
  const items = form.lines
    .map((line) => ({
      product_id: Number(line.product_id || 0),
      unit: line.unit.trim(),
      quantity: Number(line.quantity || 0),
      remark: line.remark.trim(),
    }))
    .filter((line) => line.product_id > 0 && line.unit && line.quantity > 0)
  if (!items.length) {
    toast.warning('请至少选择一条有效商品明细')
    return
  }
  saving.value = true
  try {
    await createInventoryLossOrder({
      store_id: tenantStoreId.value,
      type: form.type,
      member_id: form.type === 'gift' ? Number(form.member_id) : undefined,
      reason: form.reason.trim(),
      items,
    })
    toast.success('已提交')
    createDlg.value = false
    await qc.invalidateQueries({ queryKey: ['inventory-loss-orders'] })
    await qc.invalidateQueries({ queryKey: ['inventories'] })
    await qc.invalidateQueries({ queryKey: ['store-products'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '提交失败')
  } finally {
    saving.value = false
  }
}

const editDlg = ref(false)
const savingEdit = ref(false)
const editingId = ref<number | null>(null)
const editForm = reactive({
  reason: '',
})

function normalizeReasonValue(reason?: string): string {
  const value = String(reason || '').trim()
  return reasonValueByLabel.value.get(value) || value
}

function openEdit(row: InventoryLossOrder): void {
  editingId.value = row.id
  editForm.reason = normalizeReasonValue(row.reason)
  editDlg.value = true
}

async function submitEdit(): Promise<void> {
  const reason = editForm.reason.trim()
  if (!editingId.value) return
  if (!reason) {
    toast.warning('请选择原因说明')
    return
  }
  savingEdit.value = true
  try {
    await updateInventoryLossOrder(editingId.value, { reason })
    toast.success('原因说明已更新')
    editDlg.value = false
    await qc.invalidateQueries({ queryKey: ['inventory-loss-orders'] })
    if (detail.value?.id === editingId.value) {
      detail.value = await getInventoryLossOrder(editingId.value)
    }
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    savingEdit.value = false
  }
}

const detailDlg = ref(false)
const detail = ref<InventoryLossOrderDetail | null>(null)

async function openDetail(row: InventoryLossOrder): Promise<void> {
  try {
    detail.value = await getInventoryLossOrder(row.id)
    detailDlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载详情失败')
  }
}

async function onCancel(row: InventoryLossOrder): Promise<void> {
  const ok = await confirmDialog({ message: `撤销单据「${row.order_no}」？撤销后应由后端反向恢复库存。` })
  if (!ok) return
  try {
    await cancelInventoryLossOrder(row.id)
    toast.success('已撤销')
    await qc.invalidateQueries({ queryKey: ['inventory-loss-orders'] })
    await qc.invalidateQueries({ queryKey: ['inventories'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '撤销失败')
  }
}

function typeLabel(type: InventoryLossType | string): string {
  if (type === 'loss') return '报损'
  if (type === 'self_use') return '自用'
  if (type === 'gift') return '赠送'
  return String(type || '-')
}

function typeClass(type: InventoryLossType | string): string {
  const base = 'inline-flex rounded px-2 py-0.5 text-xs font-medium'
  if (type === 'loss') return `${base} bg-rose-50 text-rose-700`
  if (type === 'self_use') return `${base} bg-amber-50 text-amber-700`
  return `${base} bg-indigo-50 text-indigo-700`
}

function memberLabel(row: InventoryLossOrder): string {
  if (row.member) {
    const phone = String(row.member.phone || '').trim()
    const name = String(row.member.name || '').trim()
    return phone && name ? `${phone}（${name}）` : phone || name || `会员#${row.member.id}`
  }
  return row.member_id ? `会员#${row.member_id}` : '-'
}

function formatMoney(v: number | string | undefined | null): string {
  const n = Number(v ?? 0)
  return Number.isFinite(n) ? n.toFixed(2) : '0.00'
}

function formatQty(v: number | string | undefined | null): string {
  const n = Number(v ?? 0)
  if (!Number.isFinite(n)) return '0'
  return Number.isInteger(n) ? String(n) : String(Number(n.toFixed(2)))
}

function formatDateTime(v?: string): string {
  const s = String(v || '').trim()
  if (!s) return '-'
  return s.slice(0, 19).replace('T', ' ')
}
</script>

<style scoped>
.loss-line-editor {
  width: 100%;
  overflow-x: auto;
}

.loss-line-editor__head,
.loss-line-editor__row {
  display: grid;
  grid-template-columns: minmax(220px, 1.5fr) minmax(170px, 1fr) minmax(110px, 0.55fr) minmax(180px, 1fr) 72px;
  gap: 12px;
  align-items: center;
  min-width: 860px;
}

.loss-line-editor__head {
  margin-bottom: 8px;
  color: var(--color-text-2);
  font-size: 13px;
  font-weight: 600;
}

.loss-line-editor__row {
  margin-bottom: 10px;
}

.loss-line-editor__control {
  width: 100%;
  min-width: 0;
}
</style>
