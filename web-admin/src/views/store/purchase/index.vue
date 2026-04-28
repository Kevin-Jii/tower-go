<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">采购管理</h2>
      <div class="flex flex-col sm:flex-row flex-wrap gap-2 w-full md:w-auto">
        <BaseInput v-model="filterDate" class="w-full sm:w-40" type="date" />
        <BaseSelect
          v-model="filterStatus"
          class="w-full sm:w-36"
          :options="[
            { label: '全部状态', value: '' },
            { label: '待确认', value: 1 },
            { label: '已确认', value: 2 },
            { label: '已完成', value: 3 },
            { label: '已取消', value: 4 },
          ]"
          placeholder="状态"
        />
        <BaseButton variant="primary" @click="reload">查询</BaseButton>
        <BaseButton v-permission="'purchase:add'" variant="primary" @click="openCreate">新增采购单</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading" min-width="1200px">
      <template #cell-status="{ row }">
        {{ statusLabel((row as PurchaseOrder).status) }}
      </template>
      <template #cell-order_date="{ row }">
        {{ formatDate((row as PurchaseOrder).order_date) }}
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
          <BaseButton v-permission="'purchase:list'" variant="link" size="sm" @click="openDetail(row as PurchaseOrder)">详情单</BaseButton>
          <BaseButton
            v-if="canAction((row as PurchaseOrder).status, 'confirm')"
            v-permission="'purchase:edit'"
            variant="link"
            size="sm"
            @click="runAction(row as PurchaseOrder, 'confirm')"
          >
            确认
          </BaseButton>
          <BaseButton
            v-if="canAction((row as PurchaseOrder).status, 'complete')"
            v-permission="'purchase:edit'"
            variant="link"
            size="sm"
            @click="runAction(row as PurchaseOrder, 'complete')"
          >
            完成
          </BaseButton>
          <BaseButton
            v-if="canAction((row as PurchaseOrder).status, 'cancel')"
            v-permission="'purchase:edit'"
            variant="link"
            size="sm"
            @click="runAction(row as PurchaseOrder, 'cancel')"
          >
            取消
          </BaseButton>
          <BaseButton
            v-if="canDelete((row as PurchaseOrder).status)"
            v-permission="'purchase:delete'"
            variant="link"
            size="sm"
            @click="onDelete(row as PurchaseOrder)"
          >
            删除
          </BaseButton>
        </div>
      </template>
    </BaseTable>

    <div class="flex justify-end">
      <BasePagination
        :page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="(p) => (page = p)"
        @update:page-size="(s) => (pageSize = s)"
      />
    </div>

    <BaseDialog v-model="createDlg" title="新增采购单" max-width="min(560px, 96vw)">
      <div class="space-y-4 max-h-[70vh] overflow-y-auto pr-1">
        <p class="m-0 text-xs text-slate-500">报菜日期由后端自动生成（当天日期）</p>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="createForm.remark" :rows="2" />
        </BaseFormItem>
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">明细</span>
          <BaseButton variant="secondary" size="sm" @click="addLine">加一行</BaseButton>
        </div>
        <div v-for="(line, idx) in createLines" :key="idx" class="flex flex-wrap gap-2 items-end border border-[var(--color-border-2)] rounded p-3">
          <BaseFormItem label="分类 / 商品" class="min-w-[220px] flex-1">
            <a-cascader
              v-model="line.product_path"
              :options="productCascaderOptions"
              placeholder="先选分类，再选商品"
              allow-clear
              :path-mode="true"
              :check-strictly="false"
            />
          </BaseFormItem>
          <BaseFormItem label="数量 / 单位" class="min-w-[240px]">
            <div class="flex items-center gap-2">
              <BaseNumberInput v-model="line.quantity" :min="0.01" :step="0.01" class="w-28" />
              <BaseSelect v-model="line.unit_code" :options="unitOptions" placeholder="单位" class="w-28" />
            </div>
          </BaseFormItem>
          <BaseButton variant="ghost" size="sm" @click="removeLine(idx)">移除</BaseButton>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="createDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitCreate">提交</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="detailDlg" title="采购详情单" max-width="min(860px, 98vw)">
      <div v-if="detail" class="space-y-3 text-sm">
        <div class="rounded-lg border border-[var(--color-border-2)] bg-[var(--color-fill-1)] p-4">
          <h3 class="m-0 text-lg font-semibold text-slate-800">采购详情单</h3>
          <div class="mt-3 grid grid-cols-1 md:grid-cols-2 gap-y-2 gap-x-6">
            <p class="m-0"><span class="text-slate-500">单号：</span>{{ detail.order_no }}</p>
            <p class="m-0"><span class="text-slate-500">状态：</span>{{ statusLabel(detail.status) }}</p>
            <p class="m-0"><span class="text-slate-500">日期：</span>{{ formatDate(detail.order_date) }}</p>
            <p class="m-0"><span class="text-slate-500">门店：</span>{{ detail.store?.name ?? '-' }}</p>
            <p class="m-0"><span class="text-slate-500">创建人：</span>{{ detail.creator?.nickname || detail.creator?.username || '-' }}</p>
            <p class="m-0"><span class="text-slate-500">创建时间：</span>{{ formatDateTime(detail.created_at) }}</p>
          </div>
        </div>
        <BaseTable
          v-if="detailRows.length > 0"
          :columns="itemCols"
          :data="(detailRows as unknown) as Record<string, unknown>[]"
          min-width="760px"
        >
          <template #cell-quantity="{ row }">
            {{ (row as PurchaseOrderItem).quantity }} {{ (row as PurchaseOrderItem).product?.unit || '' }}
          </template>
          <template #cell-unit_price="{ row }">
            {{ money((row as PurchaseOrderItem).unit_price) }}
          </template>
          <template #cell-amount="{ row }">
            {{ money((row as PurchaseOrderItem).amount) }}
          </template>
        </BaseTable>
        <p v-else class="text-slate-500 m-0">暂无明细</p>
        <div class="flex justify-end">
          <div class="min-w-[240px] rounded-lg border border-[var(--color-border-2)] bg-[var(--color-fill-1)] px-4 py-3">
            <p class="m-0 flex items-center justify-between">
              <span class="text-slate-500">明细行数</span>
              <span class="font-semibold text-slate-800">{{ detailRows.length }}</span>
            </p>
            <p class="m-0 mt-2 flex items-center justify-between">
              <span class="text-slate-500">采购总额</span>
              <span class="font-semibold text-slate-800">{{ money(detail.total_amount) }}</span>
            </p>
          </div>
        </div>
        <div v-if="detail.remark" class="rounded border border-[var(--color-border-2)] p-3">
          <p class="m-0 text-slate-500 mb-1">备注</p>
          <p class="m-0 text-slate-700 whitespace-pre-wrap break-words">{{ detail.remark }}</p>
        </div>
      </div>
      <template #footer>
        <BaseButton
          v-permission="'printer:query'"
          variant="secondary"
          :loading="printing"
          @click="openPrintDialog"
        >
          打印详情单
        </BaseButton>
        <BaseButton variant="ghost" @click="detailDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="printDlg" title="打印采购详情单" max-width="min(460px, 96vw)">
      <div class="space-y-3">
        <p class="m-0 text-sm text-slate-600">
          单号：<span class="font-medium text-slate-800">{{ detail?.order_no || '-' }}</span>
        </p>
        <BaseFormItem label="选择打印机" required>
          <BaseSelect
            v-model="printForm.printer_id"
            :options="printerOptions"
            placeholder="请选择打印机"
            allow-clear
          />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="printDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="printing" @click="submitPrint">确认打印</BaseButton>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseNumberInput,
  BasePagination,
  BaseSelect,
  BaseTable,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import {
  cancelPurchaseOrder,
  completePurchaseOrder,
  confirmPurchaseOrder,
  createPurchaseOrder,
  deletePurchaseOrder,
  getPurchaseOrder,
  listPurchaseOrders,
  listPurchasableProducts,
} from '@/api/purchase'
import { listStorePrinters, printPurchaseOrder, type PrinterRow } from '@/api/printer'
import { listDictDataByTypeCode } from '@/api/dict'
import type { BaseSelectOption } from '@/components/base/types'
import type { DictData, PurchaseOrder, PurchaseOrderItem, StorePurchasableProduct } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()
const page = ref(1)
const pageSize = ref(10)
const filterDate = ref('')
const filterStatus = ref<number | ''>('')

const queryKey = computed(
  () => ['purchase-orders', page.value, pageSize.value, filterDate.value, filterStatus.value] as const,
)

const { data: pageData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: () =>
    listPurchaseOrders({
      page: page.value,
      page_size: pageSize.value,
      date: filterDate.value || undefined,
      ...(filterStatus.value === '' ? {} : { status: filterStatus.value as number }),
    }),
})

const list = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)

function reload(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['purchase-orders'] })
}

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['purchase-orders'] })
})

const columns: BaseTableColumn[] = [
  { key: 'order_no', label: '单号', prop: 'order_no', minWidth: '140px', ellipsis: true },
  { key: 'order_date', label: '日期', width: '120px' },
  { key: 'total_amount', label: '金额', prop: 'total_amount', width: '96px' },
  { key: 'status', label: '状态', width: '88px' },
  { key: 'actions', label: '操作', width: '300px', minWidth: '300px', align: 'right' },
]

const itemCols: BaseTableColumn[] = [
  { key: 'supplier_name', label: '供应商', minWidth: '140px' },
  { key: 'product_name', label: '商品', prop: 'product_name', minWidth: '140px', ellipsis: true },
  { key: 'quantity', label: '数量', width: '120px' },
  { key: 'unit_price', label: '单价', width: '100px' },
  { key: 'amount', label: '金额', width: '110px' },
]

function statusLabel(s: number): string {
  switch (s) {
    case 1:
      return '待确认'
    case 2:
      return '已确认'
    case 3:
      return '已完成'
    case 4:
      return '已取消'
    default:
      return String(s)
  }
}

function formatDate(v: string): string {
  if (!v) return '-'
  return v.slice(0, 10)
}

function formatDateTime(v?: string): string {
  if (!v) return '-'
  return v.slice(0, 19).replace('T', ' ')
}

function money(v?: number): string {
  if (typeof v !== 'number' || Number.isNaN(v)) return '-'
  return v.toFixed(2)
}

function canAction(status: number, action: 'confirm' | 'complete' | 'cancel'): boolean {
  if (action === 'confirm') return status === 1
  if (action === 'complete') return status === 2
  if (action === 'cancel') return status === 1 || status === 2
  return false
}

function canDelete(status: number): boolean {
  return status === 1 || status === 4
}

const createDlg = ref(false)
const saving = ref(false)
const createForm = reactive({ remark: '' })
const createLines = ref<{ product_path: Array<string | number> | string | number | undefined; quantity: number; unit_code?: string }[]>([
  { product_path: [], quantity: 1, unit_code: undefined },
])
const unitDict = ref<DictData[]>([])

const { data: productsData } = useQuery({
  queryKey: ['store-suppliers', 'products'],
  queryFn: () => listPurchasableProducts({}),
  enabled: computed(() => createDlg.value),
})
const { data: unitData } = useQuery({
  queryKey: ['dict-data', 'product_unit'],
  queryFn: () => listDictDataByTypeCode('product_unit'),
  enabled: computed(() => createDlg.value),
})
watch(unitData, (v) => {
  unitDict.value = v ?? []
})

const unitOptions = computed<BaseSelectOption[]>(() =>
  unitDict.value.map((d) => ({
    label: d.label,
    value: d.value,
  })),
)

const productCascaderOptions = computed(() => {
  const rows = productsData.value ?? []
  const grouped = new Map<string, StorePurchasableProduct[]>()
  for (const p of rows) {
    const cat = p.category?.name?.trim() || '未分类'
    if (!grouped.has(cat)) grouped.set(cat, [])
    grouped.get(cat)!.push(p)
  }
  let idx = 0
  return Array.from(grouped.entries()).map(([cat, products]) => {
    idx += 1
    return {
      label: cat,
      value: `cat-${idx}`,
      children: products.map((p) => ({
        label: `${p.name}（#${p.id}）`,
        value: p.id,
      })),
    }
  })
})

function openCreate(): void {
  createForm.remark = ''
  const firstUnit = unitOptions.value[0]?.value
  createLines.value = [{ product_path: [], quantity: 1, unit_code: firstUnit ? String(firstUnit) : undefined }]
  createDlg.value = true
  void qc.invalidateQueries({ queryKey: ['store-suppliers', 'products'] })
  void qc.invalidateQueries({ queryKey: ['dict-data', 'product_unit'] })
}

function addLine(): void {
  const firstUnit = unitOptions.value[0]?.value
  createLines.value.push({ product_path: [], quantity: 1, unit_code: firstUnit ? String(firstUnit) : undefined })
}

function removeLine(i: number): void {
  createLines.value = createLines.value.filter((_, j) => j !== i)
  if (createLines.value.length === 0) {
    const firstUnit = unitOptions.value[0]?.value
    createLines.value.push({ product_path: [], quantity: 1, unit_code: firstUnit ? String(firstUnit) : undefined })
  }
}

async function submitCreate(): Promise<void> {
  const getProductId = (path: Array<string | number> | string | number | undefined): number | null => {
    if (Array.isArray(path)) {
      const leaf = path[path.length - 1]
      const id = Number(leaf)
      return Number.isFinite(id) && id > 0 ? id : null
    }
    if (typeof path === 'number' || typeof path === 'string') {
      const id = Number(path)
      return Number.isFinite(id) && id > 0 ? id : null
    }
    return null
  }

  const items = createLines.value
    .map((l) => ({
      product_id: getProductId(l.product_path),
      quantity: l.quantity,
      unit: l.unit_code,
      remark: '',
    }))
    .filter((l) => l.product_id != null && l.quantity > 0 && l.unit)
    .map((l) => ({
      product_id: l.product_id as number,
      quantity: l.quantity,
      unit: l.unit as string,
      remark: '',
    }))
  if (!items.length) {
    toast.warning('请至少填写一条有效明细（商品、数量、单位）')
    return
  }
  saving.value = true
  try {
    await createPurchaseOrder({
      remark: createForm.remark.trim(),
      items,
    })
    toast.success('已创建')
    createDlg.value = false
    await qc.invalidateQueries({ queryKey: ['purchase-orders'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '创建失败')
  } finally {
    saving.value = false
  }
}

const detailDlg = ref(false)
const detail = ref<PurchaseOrder | null>(null)
const printDlg = ref(false)
const printing = ref(false)
const printForm = reactive<{ printer_id?: number }>({ printer_id: undefined })

const detailRows = computed(() =>
  (detail.value?.items ?? []).map((i: PurchaseOrderItem) => ({
    ...i,
    supplier_name: i.supplier?.supplier_name ?? `供应商#${i.supplier_id}`,
    product_name: i.product?.name ?? `商品#${i.product_id}`,
  })),
)

const { data: printersData } = useQuery({
  queryKey: computed(() => ['printers', detail.value?.store_id ?? 0] as const),
  queryFn: async () => {
    const storeId = detail.value?.store_id ?? 0
    if (!storeId) return []
    return listStorePrinters(storeId)
  },
  enabled: computed(() => printDlg.value && !!detail.value?.store_id),
})

const printerOptions = computed<BaseSelectOption[]>(() =>
  (printersData.value ?? [])
    .filter((p: PrinterRow) => p.status === 1)
    .map((p: PrinterRow) => ({
      label: `${p.name || p.sn}${p.is_default === 1 ? '（默认）' : ''}`,
      value: p.id,
    })),
)

async function openDetail(row: PurchaseOrder): Promise<void> {
  try {
    detail.value = await getPurchaseOrder(row.id)
    detailDlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载失败')
  }
}

function openPrintDialog(): void {
  if (!detail.value) {
    toast.warning('请先打开采购详情单')
    return
  }
  printForm.printer_id = undefined
  printDlg.value = true
}

async function submitPrint(): Promise<void> {
  if (!detail.value) return
  if (!printForm.printer_id) {
    toast.warning('请选择打印机')
    return
  }
  printing.value = true
  try {
    await printPurchaseOrder(printForm.printer_id, detail.value.id)
    toast.success('已发送打印任务')
    printDlg.value = false
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '打印失败')
  } finally {
    printing.value = false
  }
}

async function runAction(row: PurchaseOrder, act: 'confirm' | 'complete' | 'cancel'): Promise<void> {
  const labels = { confirm: '确认', complete: '完成', cancel: '取消' }
  const ok = await confirmDialog({ message: `${labels[act]}采购单「${row.order_no}」？` })
  if (!ok) return
  try {
    if (act === 'confirm') await confirmPurchaseOrder(row.id)
    else if (act === 'complete') await completePurchaseOrder(row.id)
    else await cancelPurchaseOrder(row.id)
    toast.success('已提交')
    await qc.invalidateQueries({ queryKey: ['purchase-orders'] })
    if (detail.value?.id === row.id) {
      detail.value = await getPurchaseOrder(row.id)
    }
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '操作失败')
  }
}

async function onDelete(row: PurchaseOrder): Promise<void> {
  const ok = await confirmDialog({ message: `删除采购单「${row.order_no}」？` })
  if (!ok) return
  try {
    await deletePurchaseOrder(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['purchase-orders'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>
