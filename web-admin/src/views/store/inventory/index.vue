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
        <BaseInput v-model="stockKeyword" class="w-full sm:w-56" placeholder="分类 / 商品名称" clearable @enter="reloadStock" />
        <BaseButton variant="primary" @click="reloadStock">查询</BaseButton>
        <BaseButton v-permission="'inventory:in'" variant="secondary" @click="openOrderDlg(1)">入库登记</BaseButton>
        <BaseButton v-permission="'inventory:out'" variant="secondary" @click="openOrderDlg(2)">出库登记</BaseButton>
      </div>

      <div v-if="stockLoading || productLoading" class="rounded border border-[var(--color-border-2)] p-6 text-center text-slate-500">
        库存数据加载中...
      </div>
      <div v-else-if="groupedStockCards.length === 0" class="rounded border border-[var(--color-border-2)] p-6 text-center text-slate-400">
        暂无库存商品
      </div>
      <div v-else class="space-y-4">
        <div class="flex flex-wrap gap-2">
          <BaseButton
            size="sm"
            :variant="activeCategory === '' ? 'primary' : 'ghost'"
            @click="activeCategory = ''"
          >
            全部分类
          </BaseButton>
          <BaseButton
            v-for="c in categoryTabs"
            :key="c"
            size="sm"
            :variant="activeCategory === c ? 'primary' : 'ghost'"
            @click="activeCategory = c"
          >
            {{ c }}
          </BaseButton>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-3">
          <div
            v-for="item in displayStockCards"
            :key="item.product_id"
            :class="[
              'rounded-xl border px-4 py-3 transition shadow-sm cursor-pointer select-none',
              item.low
                ? 'border-rose-200 bg-rose-50 hover:bg-rose-100'
                : 'border-[var(--color-border-2)] bg-[var(--color-bg-2)] hover:bg-[var(--color-fill-1)]',
            ]"
            @dblclick="openQtyFromCard(item)"
          >
            <div class="flex items-center justify-between gap-2">
              <p class="m-0 text-sm font-semibold text-slate-800 truncate">{{ item.product_name }}</p>
              <span class="text-xs text-slate-500">{{ item.category }}</span>
            </div>
            <div class="mt-3 flex items-end justify-between">
              <div class="space-y-1">
                <p class="m-0 text-xs text-slate-500">
                  小规格（{{ item.small_unit }}）：<span class="font-semibold">{{ item.small_qty }}</span>
                </p>
                <p v-if="item.large_unit && item.large_qty !== undefined" class="m-0 text-xs text-slate-500">
                  大规格（{{ item.large_unit }}）：<span class="font-semibold">{{ item.large_qty }}</span>
                </p>
              </div>
              <p :class="['m-0 text-xl leading-none font-bold', item.low ? 'text-rose-600' : 'text-slate-800']">
                {{ item.small_qty }}
              </p>
            </div>
            <p class="m-0 mt-2 text-[11px]" :class="item.low ? 'text-rose-600' : 'text-slate-400'">
              {{ item.low ? '库存偏低，请及时补货（双击可改数量）' : '双击卡片可修改库存数量' }}
            </p>
          </div>
        </div>
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
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">商品明细</span>
          <BaseButton variant="secondary" size="sm" @click="addOrderLine">加一行</BaseButton>
        </div>
        <div
          v-for="(line, idx) in orderLines"
          :key="idx"
          class="rounded border border-[var(--color-border-2)] p-3 flex flex-wrap items-end gap-2"
        >
          <BaseFormItem label="商品" required class="min-w-[220px] flex-1">
            <a-cascader
              v-model="line.product_path"
              :options="productCascaderOptions"
              placeholder="先选分类，再选商品"
              allow-clear
              :path-mode="true"
              :check-strictly="false"
              @change="onOrderLineProductChange(idx)"
            />
          </BaseFormItem>
          <BaseFormItem label="数量" required class="w-28">
            <BaseNumberInput v-model="line.quantity" :min="0.01" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="单位" class="w-28">
            <BaseSelect
              v-model="line.unit"
              :options="lineUnitOptions(line)"
              :disabled="lineUnitDisabled(line)"
              placeholder="单位"
            />
          </BaseFormItem>
          <BaseButton variant="ghost" size="sm" @click="removeOrderLine(idx)">移除</BaseButton>
        </div>
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
  BaseSelect,
  BaseTable,
} from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import { createInventoryOrder, getInventoryOrder, listInventories, listInventoryOrders, updateInventoryQuantity } from '@/api/inventory'
import { listPurchasableProducts } from '@/api/storeSupplier'
import { listDictDataByTypeCode } from '@/api/dict'
import { listProductUnitSpecs } from '@/api/supplierProduct'
import type { DictData, InventoryOrder, InventoryRow, ProductUnitSpec } from '@/api/types'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

const qc = useQueryClient()
const tab = ref<'stock' | 'orders'>('stock')
const userStore = useUserStore()
const tenantStoreId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)

const stockKeyword = ref('')
const LOW_STOCK_THRESHOLD = 5
const stockKey = computed(() => ['inventories', stockKeyword.value] as const)
const { data: stockData, isLoading: stockLoading } = useQuery({
  queryKey: stockKey,
  queryFn: () =>
    listInventories({
      page: 1,
      page_size: 100,
      store_id: tenantStoreId.value,
      keyword: stockKeyword.value.trim() || undefined,
    }),
  enabled: computed(() => tab.value === 'stock'),
})
const stockList = computed(() => stockData.value?.list ?? [])

const { data: productData, isLoading: productLoading } = useQuery({
  queryKey: computed(() => ['store-products', stockKeyword.value] as const),
  queryFn: () =>
    listPurchasableProducts({
      store_id: tenantStoreId.value,
      keyword: stockKeyword.value.trim() || undefined,
    }),
  enabled: computed(() => tab.value === 'stock'),
})
const productList = computed(() => productData.value ?? [])
const productIdsKey = computed(() =>
  productList.value
    .map((p) => p.id)
    .sort((a, b) => a - b)
    .join(','),
)
const { data: unitSpecsData } = useQuery({
  queryKey: computed(() => ['product-unit-specs', productIdsKey.value] as const),
  queryFn: async () => {
    const ids = productList.value.map((p) => p.id)
    if (!ids.length) return [] as ProductUnitSpec[]
    const rows = await Promise.all(ids.map((id) => listProductUnitSpecs(id)))
    return rows.flat()
  },
  enabled: computed(() => productList.value.length > 0 && tab.value === 'stock'),
})
const specsByProduct = computed(() => {
  const map = new Map<number, ProductUnitSpec[]>()
  for (const s of unitSpecsData.value ?? []) {
    if (!s.is_enabled) continue
    if (!map.has(s.product_id)) map.set(s.product_id, [])
    map.get(s.product_id)!.push(s)
  }
  for (const [, arr] of map.entries()) {
    arr.sort((a, b) => Number(a.factor_to_base) - Number(b.factor_to_base))
  }
  return map
})
const unitDict = ref<DictData[]>([])

const { data: unitData } = useQuery({
  queryKey: ['dict-data', 'product_unit'],
  queryFn: () => listDictDataByTypeCode('product_unit'),
})

watch(unitData, (v) => {
  unitDict.value = v ?? []
})

const unitOptions = computed(() =>
  unitDict.value.map((d) => ({
    label: d.label,
    value: d.value,
  })),
)

const productCascaderOptions = computed(() => {
  const grouped = new Map<string, { id: number; name: string }[]>()
  for (const p of productList.value) {
    const cat = p.category?.name?.trim() || '未分类'
    if (!grouped.has(cat)) grouped.set(cat, [])
    grouped.get(cat)!.push({ id: p.id, name: p.name })
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

interface StockCardItem {
  inventory_id: number
  product_id: number
  category: string
  product_name: string
  small_unit: string
  large_unit?: string
  small_qty: number
  large_qty?: number
  quantity: number
  low: boolean
}

const groupedStockCards = computed(() => {
  const invMap = new Map<number, InventoryRow>()
  for (const inv of stockList.value) invMap.set(inv.product_id, inv)

  const grouped = new Map<string, StockCardItem[]>()
  for (const p of productList.value) {
    const inv = invMap.get(p.id)
    const qty = Number(inv?.quantity ?? 0)
    const category = p.category?.name?.trim() || '未分类'
    const specs = specsByProduct.value.get(p.id) ?? []
    const smallSpec = specs[0]
    const largeSpec = specs.length > 1 ? specs[specs.length - 1] : undefined
    const smallFactor = Number(smallSpec?.factor_to_base || 1)
    const largeFactor = Number(largeSpec?.factor_to_base || 0)
    const smallQty = Number((qty / (smallFactor > 0 ? smallFactor : 1)).toFixed(2))
    const largeQty = largeFactor > smallFactor ? Number((qty / largeFactor).toFixed(2)) : undefined
    const item: StockCardItem = {
      inventory_id: Number(inv?.id ?? 0),
      product_id: p.id,
      category,
      product_name: p.name,
      small_unit: smallSpec?.unit_name || p.unit || '件',
      large_unit: largeSpec?.unit_name,
      small_qty: smallQty,
      large_qty: largeQty,
      quantity: qty,
      low: smallQty <= LOW_STOCK_THRESHOLD,
    }
    if (!grouped.has(category)) grouped.set(category, [])
    grouped.get(category)!.push(item)
  }

  return Array.from(grouped.entries())
    .map(([category, items]) => ({ category, items }))
    .filter((g) => g.items.length > 0)
})

const activeCategory = ref('')

const categoryTabs = computed(() => groupedStockCards.value.map((g) => g.category))

const displayStockCards = computed(() => {
  const all = groupedStockCards.value.flatMap((g) => g.items)
  if (!activeCategory.value) return all
  return all.filter((x) => x.category === activeCategory.value)
})

watch(groupedStockCards, (groups) => {
  if (!groups.some((g) => g.category === activeCategory.value)) activeCategory.value = ''
})

function reloadStock(): void {
  void qc.invalidateQueries({ queryKey: ['inventories'] })
  void qc.invalidateQueries({ queryKey: ['store-products'] })
}

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
      store_id: tenantStoreId.value,
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
  if (t === 'stock') {
    void qc.invalidateQueries({ queryKey: ['inventories'] })
    void qc.invalidateQueries({ queryKey: ['store-products'] })
  }
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

function openQtyFromCard(item: StockCardItem): void {
  if (!item.inventory_id) {
    orderTypeForm.value = 1
    orderForm.reason = '采购入库'
    orderForm.remark = ''
    orderLines.value = [makeOrderLine([item.product_id])]
    orderDlg.value = true
    toast.warning('该商品暂无库存记录，请先入库后再改数量')
    return
  }
  openQty({
    id: item.inventory_id,
    store_id: 0,
    product_id: item.product_id,
    product_name: item.product_name,
    quantity: item.small_qty,
    unit: item.small_unit,
  })
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
})

interface OrderLine {
  product_path: Array<string | number> | string | number | undefined
  quantity: number
  unit: string
}

const orderLines = ref<OrderLine[]>([])

function makeOrderLine(productPath?: Array<string | number> | string | number): OrderLine {
  return {
    product_path: productPath ?? [],
    quantity: 1,
    unit: String(unitOptions.value[0]?.value || '瓶'),
  }
}

function getProductId(path: Array<string | number> | string | number | undefined): number | null {
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

function lineUnitOptions(line: OrderLine): Array<{ label: string; value: string | number }> {
  const pid = getProductId(line.product_path)
  if (!pid) return unitOptions.value
  const specs = specsByProduct.value.get(pid) ?? []
  if (!specs.length) return unitOptions.value
  return specs.map((s) => ({
    label: s.unit_name,
    value: s.unit_code,
  }))
}

function lineUnitDisabled(line: OrderLine): boolean {
  return lineUnitOptions(line).length <= 1
}

function onOrderLineProductChange(idx: number): void {
  const line = orderLines.value[idx]
  if (!line) return
  const options = lineUnitOptions(line)
  if (!options.length) return
  line.unit = String(options[0].value)
}

function addOrderLine(): void {
  orderLines.value.push(makeOrderLine())
}

function removeOrderLine(i: number): void {
  orderLines.value = orderLines.value.filter((_, idx) => idx !== i)
  if (orderLines.value.length === 0) orderLines.value.push(makeOrderLine())
}

function openOrderDlg(type: 1 | 2): void {
  orderTypeForm.value = type
  orderForm.reason = type === 1 ? '采购入库' : '销售出库'
  orderForm.remark = ''
  orderLines.value = [makeOrderLine()]
  orderDlg.value = true
}

async function submitOrder(): Promise<void> {
  if (!orderForm.reason.trim()) {
    toast.warning('请填写原因')
    return
  }
  const items = orderLines.value
    .map((line) => ({
      product_id: getProductId(line.product_path),
      quantity: line.quantity,
      unit: line.unit.trim(),
    }))
    .filter((line) => line.product_id && line.quantity > 0 && line.unit)
    .map((line) => ({
      product_id: line.product_id as number,
      quantity: line.quantity,
      unit: line.unit,
      production_date: '',
      expiry_date: '',
      remark: '',
    }))
  if (!items.length) {
    toast.warning('请至少选择一条有效商品明细')
    return
  }
  orderSaving.value = true
  try {
    await createInventoryOrder({
      type: orderTypeForm.value,
      store_id: tenantStoreId.value,
      reason: orderForm.reason.trim(),
      remark: orderForm.remark.trim(),
      items,
    })
    toast.success('已提交')
    orderDlg.value = false
    await qc.invalidateQueries({ queryKey: ['inventories'] })
    await qc.invalidateQueries({ queryKey: ['store-products'] })
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
