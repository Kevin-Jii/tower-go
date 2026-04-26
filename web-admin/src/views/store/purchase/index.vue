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
          <BaseButton v-permission="'purchase:list'" variant="link" size="sm" @click="openDetail(row as PurchaseOrder)">详情</BaseButton>
          <BaseButton v-permission="'purchase:edit'" variant="link" size="sm" @click="runAction(row as PurchaseOrder, 'confirm')">确认</BaseButton>
          <BaseButton v-permission="'purchase:edit'" variant="link" size="sm" @click="runAction(row as PurchaseOrder, 'complete')">完成</BaseButton>
          <BaseButton v-permission="'purchase:edit'" variant="link" size="sm" @click="runAction(row as PurchaseOrder, 'cancel')">取消</BaseButton>
          <BaseButton v-permission="'purchase:delete'" variant="link" size="sm" @click="onDelete(row as PurchaseOrder)">删除</BaseButton>
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
        <BaseFormItem label="报菜日期" required>
          <BaseInput v-model="createForm.order_date" type="date" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="createForm.remark" :rows="2" />
        </BaseFormItem>
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">明细</span>
          <BaseButton variant="secondary" size="sm" @click="addLine">加一行</BaseButton>
        </div>
        <div v-for="(line, idx) in createLines" :key="idx" class="flex flex-wrap gap-2 items-end border border-[var(--color-border-2)] rounded p-3">
          <BaseFormItem label="商品" class="min-w-[200px] flex-1">
            <BaseSelect
              v-model="line.product_id"
              :options="productOptions"
              placeholder="选择可采商品"
              allow-clear
            />
          </BaseFormItem>
          <BaseFormItem label="数量" class="w-28">
            <BaseNumberInput v-model="line.quantity" :min="0.01" :step="0.01" />
          </BaseFormItem>
          <BaseButton variant="ghost" size="sm" @click="removeLine(idx)">移除</BaseButton>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="createDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitCreate">提交</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="detailDlg" title="采购单详情" max-width="min(640px, 96vw)">
      <div v-if="detail" class="space-y-3 text-sm">
        <p class="m-0"><span class="text-slate-500">单号</span> {{ detail.order_no }}</p>
        <p class="m-0"><span class="text-slate-500">状态</span> {{ statusLabel(detail.status) }}</p>
        <p class="m-0"><span class="text-slate-500">日期</span> {{ formatDate(detail.order_date) }}</p>
        <p class="m-0"><span class="text-slate-500">金额</span> {{ detail.total_amount?.toFixed?.(2) ?? detail.total_amount }}</p>
        <BaseTable
          v-if="detailRows.length > 0"
          :columns="itemCols"
          :data="(detailRows as unknown) as Record<string, unknown>[]"
          min-width="520px"
        />
        <p v-else class="text-slate-500 m-0">暂无明细</p>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="detailDlg = false">关闭</BaseButton>
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
import type { BaseSelectOption } from '@/components/base/types'
import type { PurchaseOrder, PurchaseOrderItem } from '@/api/types'
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
  { key: 'product_name', label: '商品', prop: 'product_name', minWidth: '140px', ellipsis: true },
  { key: 'quantity', label: '数量', prop: 'quantity', width: '80px' },
  { key: 'unit_price', label: '单价', prop: 'unit_price', width: '80px' },
  { key: 'amount', label: '金额', prop: 'amount', width: '88px' },
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

const { data: productsData } = useQuery({
  queryKey: ['store-suppliers', 'products'],
  queryFn: () => listPurchasableProducts({}),
  enabled: computed(() => createDlg.value),
})

const productOptions = computed<BaseSelectOption[]>(() => {
  const rows = productsData.value ?? []
  return rows.map((p) => ({
    label: `${p.name}（#${p.id}）`,
    value: p.id,
  }))
})

const createDlg = ref(false)
const saving = ref(false)
const createForm = reactive({ order_date: '', remark: '' })
const createLines = ref<{ product_id: number | undefined; quantity: number }[]>([{ product_id: undefined, quantity: 1 }])

function openCreate(): void {
  const t = new Date()
  createForm.order_date = `${t.getFullYear()}-${String(t.getMonth() + 1).padStart(2, '0')}-${String(t.getDate()).padStart(2, '0')}`
  createForm.remark = ''
  createLines.value = [{ product_id: undefined, quantity: 1 }]
  createDlg.value = true
  void qc.invalidateQueries({ queryKey: ['store-suppliers', 'products'] })
}

function addLine(): void {
  createLines.value.push({ product_id: undefined, quantity: 1 })
}

function removeLine(i: number): void {
  createLines.value = createLines.value.filter((_, j) => j !== i)
  if (createLines.value.length === 0) createLines.value.push({ product_id: undefined, quantity: 1 })
}

async function submitCreate(): Promise<void> {
  if (!createForm.order_date) {
    toast.warning('请选择报菜日期')
    return
  }
  const items = createLines.value
    .filter((l) => l.product_id != null && l.quantity > 0)
    .map((l) => ({
      product_id: l.product_id as number,
      quantity: l.quantity,
      remark: '',
    }))
  if (!items.length) {
    toast.warning('请至少选择一条有效明细')
    return
  }
  saving.value = true
  try {
    await createPurchaseOrder({
      order_date: createForm.order_date,
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

const detailRows = computed(() =>
  (detail.value?.items ?? []).map((i: PurchaseOrderItem) => ({
    ...i,
    product_name: i.product?.name ?? `商品#${i.product_id}`,
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
