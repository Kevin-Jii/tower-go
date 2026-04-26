<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">库存管理</h2>
      <div class="inline-flex rounded-lg border border-[var(--color-border-2)] p-0.5 bg-[var(--color-bg-2)]">
        <BaseButton
          size="sm"
          :variant="tab === 'stock' ? 'primary' : 'ghost'"
          class="!min-w-[7rem]"
          @click="tab = 'stock'"
        >
          库存快照
        </BaseButton>
        <BaseButton
          size="sm"
          :variant="tab === 'orders' ? 'primary' : 'ghost'"
          class="!min-w-[7rem]"
          @click="tab = 'orders'"
        >
          出入库单
        </BaseButton>
      </div>
    </div>

    <template v-if="tab === 'stock'">
      <div class="flex flex-col sm:flex-row gap-2">
        <BaseInput v-model="stockKeyword" class="w-full sm:w-56" placeholder="商品名称" clearable @enter="reloadStock" />
        <BaseButton variant="primary" @click="reloadStock">查询</BaseButton>
        <BaseButton v-permission="'inventory:in'" variant="secondary" @click="openOrderDlg(1)">入库登记</BaseButton>
        <BaseButton v-permission="'inventory:out'" variant="secondary" @click="openOrderDlg(2)">出库登记</BaseButton>
      </div>
      <BaseTable :columns="stockCols" :data="(stockList as unknown) as Record<string, unknown>[]" :loading="stockLoading" min-width="880px">
        <template #cell-actions="{ row }">
          <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
            <BaseButton
              v-permission="['inventory:in', 'inventory:out']"
              variant="link"
              size="sm"
              @click="openQty(row as InventoryRow)"
            >
              改数量
            </BaseButton>
          </div>
        </template>
      </BaseTable>
      <div class="flex justify-end">
        <BasePagination
          :page="stockPage"
          :page-size="stockPageSize"
          :total="stockTotal"
          @update:page="(p) => (stockPage = p)"
          @update:page-size="(s) => (stockPageSize = s)"
        />
      </div>
    </template>

    <template v-else>
      <div class="flex flex-col sm:flex-row flex-wrap gap-2">
        <BaseInput v-model="orderNo" class="w-full sm:w-44" placeholder="单号" clearable @enter="reloadOrders" />
        <BaseInput v-model="orderDate" class="w-full sm:w-40" type="date" />
        <BaseSelect
          v-model="orderType"
          class="w-full sm:w-32"
          :options="[
            { label: '全部类型', value: '' },
            { label: '入库', value: 1 },
            { label: '出库', value: 2 },
          ]"
        />
        <BaseButton variant="primary" @click="reloadOrders">查询</BaseButton>
      </div>
      <BaseTable :columns="orderCols" :data="(orderList as unknown) as Record<string, unknown>[]" :loading="orderLoading" min-width="920px">
        <template #cell-type="{ row }">
          {{ (row as InventoryOrder).type === 1 ? '入库' : '出库' }}
        </template>
        <template #cell-actions="{ row }">
          <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
            <BaseButton variant="link" size="sm" @click="openOrderDetail(row as InventoryOrder)">详情</BaseButton>
          </div>
        </template>
      </BaseTable>
      <div class="flex justify-end">
        <BasePagination
          :page="orderPage"
          :page-size="orderPageSize"
          :total="orderTotal"
          @update:page="(p) => (orderPage = p)"
          @update:page-size="(s) => (orderPageSize = s)"
        />
      </div>
    </template>

    <BaseDialog v-model="qtyDlg" title="调整库存数量" max-width="min(400px, 96vw)">
      <div class="space-y-4">
        <p class="m-0 text-sm text-slate-600">{{ qtyRow?.product_name }}（当前 {{ qtyRow?.quantity }} {{ qtyRow?.unit }}）</p>
        <BaseFormItem label="新数量" required>
          <BaseNumberInput v-model="qtyForm.quantity" :min="0" :step="0.01" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseInput v-model="qtyForm.remark" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="qtyDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="qtySaving" @click="saveQty">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="orderDlg" :title="orderTypeForm === 1 ? '入库单' : '出库单'" max-width="min(480px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="原因" required>
          <BaseInput v-model="orderForm.reason" placeholder="如 采购入库 / 销售出库" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseInput v-model="orderForm.remark" />
        </BaseFormItem>
        <BaseFormItem label="商品ID" required>
          <BaseNumberInput v-model="orderForm.product_id" :min="1" :step="1" />
        </BaseFormItem>
        <BaseFormItem label="数量" required>
          <BaseNumberInput v-model="orderForm.quantity" :min="0.01" :step="0.01" />
        </BaseFormItem>
        <BaseFormItem label="单位">
          <BaseInput v-model="orderForm.unit" placeholder="瓶 / 箱" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="orderDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="orderSaving" @click="submitOrder">提交</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="detailDlg" title="出入库单详情" max-width="min(560px, 96vw)">
      <pre v-if="orderDetailJson" class="text-xs overflow-auto max-h-[60vh] m-0 p-3 rounded bg-[var(--color-fill-2)]">{{ orderDetailJson }}</pre>
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
} from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import { createInventoryOrder, getInventoryOrder, listInventories, listInventoryOrders, updateInventoryQuantity } from '@/api/inventory'
import type { InventoryOrder, InventoryRow } from '@/api/types'
import { toast } from '@/feedback/toast'

const qc = useQueryClient()
const tab = ref<'stock' | 'orders'>('stock')

const stockKeyword = ref('')
const stockPage = ref(1)
const stockPageSize = ref(10)
const stockKey = computed(
  () => ['inventories', stockPage.value, stockPageSize.value, stockKeyword.value] as const,
)
const { data: stockData, isLoading: stockLoading } = useQuery({
  queryKey: stockKey,
  queryFn: () =>
    listInventories({
      page: stockPage.value,
      page_size: stockPageSize.value,
      keyword: stockKeyword.value.trim() || undefined,
    }),
  enabled: computed(() => tab.value === 'stock'),
})
const stockList = computed(() => stockData.value?.list ?? [])
const stockTotal = computed(() => stockData.value?.total ?? 0)

function reloadStock(): void {
  stockPage.value = 1
  void qc.invalidateQueries({ queryKey: ['inventories'] })
}

watch([stockPage, stockPageSize], () => {
  if (tab.value === 'stock') void qc.invalidateQueries({ queryKey: ['inventories'] })
})

const stockCols: BaseTableColumn[] = [
  { key: 'product_name', label: '商品', prop: 'product_name', minWidth: '160px', ellipsis: true },
  { key: 'quantity', label: '数量', prop: 'quantity', width: '96px' },
  { key: 'unit', label: '单位', prop: 'unit', width: '72px' },
  { key: 'store_name', label: '门店', prop: 'store_name', width: '120px', ellipsis: true },
  { key: 'actions', label: '操作', width: '120px', align: 'right' },
]

const orderNo = ref('')
const orderDate = ref('')
const orderType = ref<number | ''>('')
const orderPage = ref(1)
const orderPageSize = ref(10)
const orderKey = computed(
  () => ['inventory-orders', orderPage.value, orderPageSize.value, orderNo.value, orderDate.value, orderType.value] as const,
)
const { data: orderData, isLoading: orderLoading } = useQuery({
  queryKey: orderKey,
  queryFn: () =>
    listInventoryOrders({
      page: orderPage.value,
      page_size: orderPageSize.value,
      order_no: orderNo.value.trim() || undefined,
      date: orderDate.value || undefined,
      ...(orderType.value === '' ? {} : { type: orderType.value as number }),
    }),
  enabled: computed(() => tab.value === 'orders'),
})
const orderList = computed(() => orderData.value?.list ?? [])
const orderTotal = computed(() => orderData.value?.total ?? 0)

function reloadOrders(): void {
  orderPage.value = 1
  void qc.invalidateQueries({ queryKey: ['inventory-orders'] })
}

watch([orderPage, orderPageSize], () => {
  if (tab.value === 'orders') void qc.invalidateQueries({ queryKey: ['inventory-orders'] })
})

watch(tab, (t) => {
  if (t === 'stock') void qc.invalidateQueries({ queryKey: ['inventories'] })
  else void qc.invalidateQueries({ queryKey: ['inventory-orders'] })
})

const orderCols: BaseTableColumn[] = [
  { key: 'order_no', label: '单号', prop: 'order_no', minWidth: '140px', ellipsis: true },
  { key: 'type', label: '类型', width: '72px' },
  { key: 'reason', label: '原因', prop: 'reason', minWidth: '100px', ellipsis: true },
  { key: 'total_quantity', label: '总数量', prop: 'total_quantity', width: '88px' },
  { key: 'operator_name', label: '操作人', prop: 'operator_name', width: '100px' },
  { key: 'created_at', label: '时间', prop: 'created_at', width: '168px' },
  { key: 'actions', label: '操作', width: '100px', align: 'right' },
]

const qtyDlg = ref(false)
const qtySaving = ref(false)
const qtyRow = ref<InventoryRow | null>(null)
const qtyForm = reactive({ quantity: 0, remark: '' })

function openQty(row: InventoryRow): void {
  qtyRow.value = row
  qtyForm.quantity = row.quantity
  qtyForm.remark = ''
  qtyDlg.value = true
}

async function saveQty(): Promise<void> {
  if (!qtyRow.value) return
  qtySaving.value = true
  try {
    await updateInventoryQuantity(qtyRow.value.id, {
      quantity: qtyForm.quantity,
      remark: qtyForm.remark.trim(),
    })
    toast.success('已更新')
    qtyDlg.value = false
    await qc.invalidateQueries({ queryKey: ['inventories'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    qtySaving.value = false
  }
}

const orderDlg = ref(false)
const orderSaving = ref(false)
const orderTypeForm = ref<1 | 2>(1)
const orderForm = reactive({
  reason: '',
  remark: '',
  product_id: 1,
  quantity: 1,
  unit: '瓶',
})

function openOrderDlg(type: 1 | 2): void {
  orderTypeForm.value = type
  orderForm.reason = type === 1 ? '采购入库' : '销售出库'
  orderForm.remark = ''
  orderForm.product_id = 1
  orderForm.quantity = 1
  orderForm.unit = '瓶'
  orderDlg.value = true
}

async function submitOrder(): Promise<void> {
  if (!orderForm.reason.trim()) {
    toast.warning('请填写原因')
    return
  }
  orderSaving.value = true
  try {
    await createInventoryOrder({
      type: orderTypeForm.value,
      reason: orderForm.reason.trim(),
      remark: orderForm.remark.trim(),
      items: [
        {
          product_id: orderForm.product_id,
          quantity: orderForm.quantity,
          unit: orderForm.unit.trim() || '瓶',
          production_date: '',
          expiry_date: '',
          remark: '',
        },
      ],
    })
    toast.success('已提交')
    orderDlg.value = false
    await qc.invalidateQueries({ queryKey: ['inventories'] })
    await qc.invalidateQueries({ queryKey: ['inventory-orders'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    orderSaving.value = false
  }
}

const detailDlg = ref(false)
const orderDetailJson = ref('')

async function openOrderDetail(row: InventoryOrder): Promise<void> {
  try {
    const full = await getInventoryOrder(row.id)
    orderDetailJson.value = JSON.stringify(full, null, 2)
    detailDlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载失败')
  }
}
</script>
