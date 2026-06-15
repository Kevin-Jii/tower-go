<template>
  <div class="flex min-h-0 flex-1 flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
      <div>
        <h2 class="page-title">B2B 供货</h2>
        <p class="m-0 mt-1 text-sm text-slate-500">餐饮客户、供货价格、供货单和应收状态统一管理。</p>
      </div>
      <div class="inline-flex rounded-lg border border-[var(--color-border-2)] bg-[var(--color-bg-2)] p-0.5 shadow-sm">
        <BaseButton size="sm" :variant="tab === 'customers' ? 'primary' : 'ghost'" @click="tab = 'customers'">客户</BaseButton>
        <BaseButton size="sm" :variant="tab === 'prices' ? 'primary' : 'ghost'" @click="tab = 'prices'">供货价</BaseButton>
        <BaseButton size="sm" :variant="tab === 'orders' ? 'primary' : 'ghost'" @click="tab = 'orders'">供货单</BaseButton>
      </div>
    </div>

    <template v-if="tab === 'customers'">
      <div class="flex flex-wrap items-center gap-2">
        <BaseInput v-model="customerKeyword" class="w-56" placeholder="客户 / 电话 / 联系人" clearable @enter="reloadCustomers" />
        <BaseSelect v-model="customerStatus" class="w-32" :options="customerStatusOptions" />
        <BaseButton variant="primary" @click="reloadCustomers">查询</BaseButton>
        <BaseButton variant="primary" @click="openCustomerCreate">新增客户</BaseButton>
      </div>
      <BaseTable :columns="customerColumns" :data="(customers as unknown) as Record<string, unknown>[]" :loading="customerLoading" min-width="1080px">
        <template #cell-settlement="{ row }">{{ settlementLabel((row as B2BCustomer).settlement) }}</template>
        <template #cell-status="{ row }">
          <a-tag :color="(row as B2BCustomer).status === 1 ? 'green' : 'gray'">{{ (row as B2BCustomer).status === 1 ? '启用' : '停用' }}</a-tag>
        </template>
        <template #cell-receivable="{ row }">{{ money((row as B2BCustomer).receivable) }}</template>
        <template #cell-actions="{ row }">
          <BaseTableRowActions :actions="customerActions(row as B2BCustomer)" />
        </template>
      </BaseTable>
      <div class="flex justify-end">
        <BasePagination :page="customerPage" :page-size="customerPageSize" :total="customerTotal" @update:page="(p) => (customerPage = p)" @update:page-size="(s) => (customerPageSize = s)" />
      </div>
    </template>

    <template v-else-if="tab === 'prices'">
      <div class="flex flex-wrap items-center gap-2">
        <BaseInput v-model="priceKeyword" class="w-56" placeholder="商品 / 客户 / 规格" clearable @enter="reloadPrices" />
        <BaseButton variant="primary" @click="reloadPrices">查询</BaseButton>
        <BaseButton variant="primary" @click="openPriceCreate">新增供货价</BaseButton>
      </div>
      <BaseTable :columns="priceColumns" :data="(groupedPrices as unknown) as Record<string, unknown>[]" :loading="priceLoading" min-width="1120px" row-key="id">
        <template #cell-owner="{ row }">{{ (row as PriceGroupRow).owner }}</template>
        <template #cell-product="{ row }">{{ (row as PriceGroupRow).productName }}</template>
        <template #cell-specs="{ row }">
          <div class="price-spec-grid">
            <div class="price-spec-head">规格</div>
            <div class="price-spec-head">供货价</div>
            <div class="price-spec-head">起订</div>
            <div class="price-spec-head">状态</div>
            <div class="price-spec-head text-right">操作</div>
            <template v-for="item in (row as PriceGroupRow).items" :key="item.id">
              <div class="price-spec-cell font-medium text-slate-800">{{ item.unit_name }}</div>
              <div class="price-spec-cell">{{ money(item.supply_price) }}</div>
              <div class="price-spec-cell">{{ item.min_quantity }}</div>
              <div class="price-spec-cell">
                <a-tag :color="item.is_enabled ? 'green' : 'gray'">{{ item.is_enabled ? '启用' : '停用' }}</a-tag>
              </div>
              <div class="price-spec-cell text-right">
                <BaseButton variant="link" size="sm" class="text-red-500" @click="onDeletePrice(item)">删除</BaseButton>
              </div>
            </template>
          </div>
        </template>
        <template #cell-actions="{ row }">
          <BaseButton variant="link" size="sm" @click="openPriceGroupEdit(row as PriceGroupRow)">编辑</BaseButton>
        </template>
      </BaseTable>
      <div class="flex justify-end">
        <BasePagination :page="pricePage" :page-size="pricePageSize" :total="priceTotal" @update:page="(p) => (pricePage = p)" @update:page-size="(s) => (pricePageSize = s)" />
      </div>
    </template>

    <template v-else>
      <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
        <div v-for="item in orderSummary" :key="item.label" class="rounded border border-slate-200 bg-white px-4 py-3">
          <div class="text-xs font-medium text-slate-500">{{ item.label }}</div>
          <div class="mt-1 text-lg font-semibold text-slate-900">{{ item.value }}</div>
        </div>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <BaseInput v-model="orderKeyword" class="w-52" placeholder="单号 / 客户" clearable @enter="reloadOrders" />
        <BaseSelect v-model="paymentStatus" class="w-32" :options="paymentStatusOptions" />
        <BaseInput v-model="startDate" class="w-36" type="date" />
        <BaseInput v-model="endDate" class="w-36" type="date" />
        <BaseButton variant="primary" @click="reloadOrders">查询</BaseButton>
        <BaseButton variant="primary" @click="openOrderCreate">新增供货单</BaseButton>
      </div>
      <BaseTable :columns="orderColumns" :data="(orders as unknown) as Record<string, unknown>[]" :loading="orderLoading" min-width="1180px">
        <template #cell-total_amount="{ row }">{{ money((row as B2BSupplyOrder).total_amount) }}</template>
        <template #cell-paid_amount="{ row }">{{ money((row as B2BSupplyOrder).paid_amount) }}</template>
        <template #cell-unpaid_amount="{ row }">{{ money((row as B2BSupplyOrder).unpaid_amount) }}</template>
        <template #cell-profit_amount="{ row }">{{ money((row as B2BSupplyOrder).profit_amount) }}</template>
        <template #cell-payment_status="{ row }">{{ paymentStatusLabel((row as B2BSupplyOrder).payment_status) }}</template>
        <template #cell-delivery_status="{ row }">{{ deliveryStatusLabel((row as B2BSupplyOrder).delivery_status) }}</template>
        <template #cell-actions="{ row }">
          <BaseButton variant="link" size="sm" @click="openOrderDetail(row as B2BSupplyOrder)">详情</BaseButton>
        </template>
      </BaseTable>
      <div class="flex justify-end">
        <BasePagination :page="orderPage" :page-size="orderPageSize" :total="orderTotal" @update:page="(p) => (orderPage = p)" @update:page-size="(s) => (orderPageSize = s)" />
      </div>
    </template>

    <BaseDialog v-model="customerDlg" :title="customerEditId ? '编辑客户' : '新增客户'" max-width="min(640px, 96vw)">
      <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <BaseFormItem label="客户名称" required><BaseInput v-model="customerForm.name" /></BaseFormItem>
        <BaseFormItem label="客户类型"><BaseInput v-model="customerForm.customer_type" placeholder="餐饮店 / 酒吧" /></BaseFormItem>
        <BaseFormItem label="联系人"><BaseInput v-model="customerForm.contact_person" /></BaseFormItem>
        <BaseFormItem label="电话"><BaseInput v-model="customerForm.phone" /></BaseFormItem>
        <BaseFormItem label="结算方式"><BaseSelect v-model="customerForm.settlement" :options="settlementOptions" /></BaseFormItem>
        <BaseFormItem label="价格等级"><BaseInput v-model="customerForm.price_level" placeholder="如 A / B / VIP" /></BaseFormItem>
        <BaseFormItem label="信用额度"><BaseNumberInput v-model="customerForm.credit_limit" :min="0" :step="100" /></BaseFormItem>
        <BaseFormItem v-if="customerEditId" label="状态"><BaseSelect v-model="customerForm.status" :options="[{ label: '启用', value: 1 }, { label: '停用', value: 2 }]" /></BaseFormItem>
        <BaseFormItem label="地址" class="sm:col-span-2"><BaseInput v-model="customerForm.address" /></BaseFormItem>
        <BaseFormItem label="备注" class="sm:col-span-2"><BaseTextarea v-model="customerForm.remark" :rows="2" /></BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="customerDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitCustomer">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="priceDlg" title="供货价格" max-width="min(920px, 96vw)">
      <div class="space-y-4">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
        <BaseFormItem label="客户专属价">
          <BaseSelect v-model="priceForm.customer_id" :options="customerOptionsWithNone" placeholder="不选则为价格等级价" />
        </BaseFormItem>
        <BaseFormItem label="价格等级">
          <BaseInput v-model="priceForm.price_level" placeholder="客户专属价可不填；等级价必填" />
        </BaseFormItem>
        <BaseFormItem label="商品" required class="md:col-span-2">
          <BaseSelect v-model="priceForm.product_id" :options="productOptions" @update:model-value="onPriceProductChange" />
        </BaseFormItem>
        </div>
        <div class="rounded border border-slate-200 bg-white">
          <div class="grid grid-cols-[1.2fr_130px_110px_90px_1fr] gap-2 border-b border-slate-200 bg-slate-50 px-3 py-2 text-sm font-medium text-slate-600">
            <div>规格</div>
            <div>供货价</div>
            <div>起订</div>
            <div>启用</div>
            <div>备注</div>
          </div>
          <div v-if="priceLines.length" class="divide-y divide-slate-100">
            <div v-for="line in priceLines" :key="line.unit_spec_id" class="grid grid-cols-1 gap-2 px-3 py-3 md:grid-cols-[1.2fr_130px_110px_90px_1fr] md:items-center">
              <div>
                <div class="font-medium text-slate-900">{{ line.unit_name }}</div>
                <div class="mt-0.5 text-xs text-slate-500">换算 {{ line.factor_to_base }} / 售价 {{ money(line.sale_price) }}</div>
              </div>
              <BaseNumberInput v-model="line.supply_price" :min="0" :step="0.01" />
              <BaseNumberInput v-model="line.min_quantity" :min="0" :step="1" />
              <BaseSwitch v-model="line.is_enabled" :active-value="true" :inactive-value="false" />
              <BaseInput v-model="line.remark" placeholder="可选" />
            </div>
          </div>
          <div v-else class="px-3 py-8 text-center text-sm text-slate-500">
            选择商品后会显示该商品的所有规格
          </div>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="priceDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitPrice">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="orderDlg" title="新增供货单" max-width="min(980px, 96vw)">
      <div class="space-y-4">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
          <BaseFormItem label="客户" required><BaseSelect v-model="orderForm.customer_id" :options="customerOptions" @update:model-value="onOrderCustomerChange" /></BaseFormItem>
          <BaseFormItem label="供货日期"><a-date-picker v-model="orderForm.order_date" value-format="YYYY-MM-DD" class="w-full" /></BaseFormItem>
          <BaseFormItem label="已收金额"><BaseNumberInput v-model="orderForm.paid_amount" :min="0" :step="0.01" /></BaseFormItem>
        </div>
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium text-slate-700">商品明细</span>
          <BaseButton variant="secondary" size="sm" @click="addOrderLine">加一行</BaseButton>
        </div>
        <div class="space-y-2">
          <div v-for="(line, idx) in orderForm.items" :key="idx" class="grid grid-cols-1 gap-2 rounded border border-slate-200 p-3 md:grid-cols-[1.4fr_1.3fr_110px_120px_80px] md:items-end">
            <BaseFormItem label="商品" required><BaseSelect v-model="line.product_id" :options="orderProductOptions" :disabled="!orderForm.customer_id || orderPriceLoading" placeholder="先选择客户" @update:model-value="onOrderProductChange(idx)" /></BaseFormItem>
            <BaseFormItem label="规格" required><BaseSelect v-model="line.unit_spec_id" :options="lineSpecOptions(line)" :disabled="!orderForm.customer_id || !line.product_id || orderPriceLoading" @update:model-value="fillLinePrice(idx)" /></BaseFormItem>
            <BaseFormItem label="数量" required><BaseNumberInput v-model="line.quantity" :min="0.01" :step="1" /></BaseFormItem>
            <BaseFormItem label="供货价"><BaseNumberInput v-model="line.supply_price" :min="0" :step="0.01" disabled placeholder="自动获取" /></BaseFormItem>
            <BaseButton variant="ghost" size="sm" :disabled="orderForm.items.length <= 1" @click="removeOrderLine(idx)">移除</BaseButton>
          </div>
        </div>
        <div class="rounded bg-slate-50 px-3 py-2 text-sm text-slate-700">
          合计：<span class="font-semibold text-slate-900">{{ money(orderTotalAmount) }}</span>
        </div>
        <BaseFormItem label="备注"><BaseTextarea v-model="orderForm.remark" :rows="2" /></BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="orderDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitOrder">提交供货单</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="detailDlg" title="供货单详情" max-width="min(760px, 96vw)">
      <pre class="m-0 max-h-[60vh] overflow-auto rounded bg-slate-50 p-3 text-xs">{{ detailJson }}</pre>
      <template #footer><BaseButton variant="ghost" @click="detailDlg = false">关闭</BaseButton></template>
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
  BaseSwitch,
  BaseTable,
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn, TableRowAction } from '@/components/base/types'
import {
  createB2BCustomer,
  createB2BSupplyOrder,
  deleteB2BPrice,
  getB2BSupplyOrder,
  listB2BCustomers,
  listB2BPrices,
  listB2BSupplyOrders,
  updateB2BCustomer,
  upsertB2BPrice,
} from '@/api/b2b'
import { listPurchasableProducts } from '@/api/storeSupplier'
import type { B2BCustomer, B2BCustomerProductPrice, B2BSupplyOrder, ProductUnitSpec } from '@/api/types'
import { confirmDialog } from '@/feedback/confirm'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

const qc = useQueryClient()
const userStore = useUserStore()
const storeId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)
const tab = ref<'customers' | 'prices' | 'orders'>('customers')
const saving = ref(false)

const settlementOptions: BaseSelectOption[] = [
  { label: '现结', value: 'cash' },
  { label: '周结', value: 'week' },
  { label: '月结', value: 'month' },
]
const customerStatusOptions: BaseSelectOption[] = [
  { label: '全部', value: '' },
  { label: '启用', value: 1 },
  { label: '停用', value: 2 },
]
const paymentStatusOptions: BaseSelectOption[] = [
  { label: '全部收款', value: '' },
  { label: '未收', value: 1 },
  { label: '部分', value: 2 },
  { label: '已收', value: 3 },
]

const customerKeyword = ref('')
const customerStatus = ref<number | ''>('')
const customerPage = ref(1)
const customerPageSize = ref(10)
const customerQueryKey = computed(() => ['b2b-customers', customerKeyword.value, customerStatus.value, customerPage.value, customerPageSize.value] as const)
const { data: customerPageData, isLoading: customerLoading } = useQuery({
  queryKey: customerQueryKey,
  queryFn: () => listB2BCustomers({ keyword: customerKeyword.value.trim() || undefined, status: customerStatus.value || undefined, page: customerPage.value, page_size: customerPageSize.value }),
})
const customers = computed(() => customerPageData.value?.list ?? [])
const customerTotal = computed(() => customerPageData.value?.total ?? 0)
const customerOptions = computed<BaseSelectOption[]>(() => customers.value.map((x) => ({ label: x.name, value: x.id })))
const customerOptionsWithNone = computed<BaseSelectOption[]>(() => [{ label: '不指定客户', value: 0 }, ...customerOptions.value])

const { data: productData } = useQuery({
  queryKey: computed(() => ['b2b-products', storeId.value] as const),
  queryFn: () => listPurchasableProducts({ store_id: storeId.value }),
})
const products = computed(() => productData.value ?? [])
const productOptions = computed<BaseSelectOption[]>(() => products.value.map((p) => ({ label: p.name, value: p.id })))
const specsByProduct = computed(() => {
  const map = new Map<number, ProductUnitSpec[]>()
  for (const p of products.value) {
    map.set(p.id, (p.unit_specs ?? []).filter((s) => s.is_enabled))
  }
  return map
})

const customerColumns: BaseTableColumn[] = [
  { key: 'name', label: '客户', prop: 'name', minWidth: '150px', ellipsis: true },
  { key: 'customer_type', label: '类型', prop: 'customer_type', width: '100px' },
  { key: 'contact_person', label: '联系人', prop: 'contact_person', width: '100px' },
  { key: 'phone', label: '电话', prop: 'phone', width: '130px' },
  { key: 'settlement', label: '结算', width: '80px' },
  { key: 'price_level', label: '等级', prop: 'price_level', width: '80px' },
  { key: 'receivable', label: '应收', width: '100px' },
  { key: 'status', label: '状态', width: '80px' },
  { key: 'actions', label: '操作', width: '120px', align: 'right' },
]

function reloadCustomers(): void {
  customerPage.value = 1
  void qc.invalidateQueries({ queryKey: ['b2b-customers'] })
}

const customerDlg = ref(false)
const customerEditId = ref(0)
const customerForm = reactive({
  name: '',
  customer_type: '',
  contact_person: '',
  phone: '',
  address: '',
  settlement: 'cash',
  price_level: '',
  credit_limit: 0,
  status: 1,
  remark: '',
})

function resetCustomer(): void {
  customerForm.name = ''
  customerForm.customer_type = ''
  customerForm.contact_person = ''
  customerForm.phone = ''
  customerForm.address = ''
  customerForm.settlement = 'cash'
  customerForm.price_level = ''
  customerForm.credit_limit = 0
  customerForm.status = 1
  customerForm.remark = ''
}

function openCustomerCreate(): void {
  customerEditId.value = 0
  resetCustomer()
  customerDlg.value = true
}

function openCustomerEdit(row: B2BCustomer): void {
  customerEditId.value = row.id
  customerForm.name = row.name
  customerForm.customer_type = row.customer_type || ''
  customerForm.contact_person = row.contact_person || ''
  customerForm.phone = row.phone || ''
  customerForm.address = row.address || ''
  customerForm.settlement = row.settlement || 'cash'
  customerForm.price_level = row.price_level || ''
  customerForm.credit_limit = Number(row.credit_limit || 0)
  customerForm.status = row.status || 1
  customerForm.remark = row.remark || ''
  customerDlg.value = true
}

async function submitCustomer(): Promise<void> {
  if (!customerForm.name.trim()) {
    toast.warning('请填写客户名称')
    return
  }
  saving.value = true
  try {
    if (customerEditId.value) await updateB2BCustomer(customerEditId.value, { ...customerForm })
    else await createB2BCustomer({ ...customerForm, store_id: storeId.value })
    toast.success('已保存')
    customerDlg.value = false
    reloadCustomers()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

function customerActions(row: B2BCustomer): TableRowAction[] {
  return [{ label: '编辑', permission: 'b2b:customer:edit', onClick: () => openCustomerEdit(row) }]
}

const priceKeyword = ref('')
const pricePage = ref(1)
const pricePageSize = ref(10)
const priceQueryKey = computed(() => ['b2b-prices', priceKeyword.value, pricePage.value, pricePageSize.value] as const)
const { data: pricePageData, isLoading: priceLoading } = useQuery({
  queryKey: priceQueryKey,
  queryFn: () => listB2BPrices({ keyword: priceKeyword.value.trim() || undefined, page: pricePage.value, page_size: pricePageSize.value }),
})
const prices = computed(() => pricePageData.value?.list ?? [])
const priceTotal = computed(() => pricePageData.value?.total ?? 0)
const priceColumns: BaseTableColumn[] = [
  { key: 'owner', label: '适用对象', minWidth: '160px', ellipsis: true },
  { key: 'product', label: '商品', minWidth: '160px', ellipsis: true },
  { key: 'specs', label: '规格价格', minWidth: '520px' },
  { key: 'actions', label: '操作', width: '100px', align: 'right' },
]

interface PriceGroupRow {
  id: string
  owner: string
  productName: string
  customer_id: number | null
  price_level: string
  product_id: number
  items: B2BCustomerProductPrice[]
}

const groupedPrices = computed<PriceGroupRow[]>(() => {
  const map = new Map<string, PriceGroupRow>()
  for (const item of prices.value) {
    const ownerKey = item.customer_id ? `c:${item.customer_id}` : `l:${item.price_level || ''}`
    const key = `${ownerKey}|p:${item.product_id}`
    const group = map.get(key)
    if (group) {
      group.items.push(item)
    } else {
      map.set(key, {
        id: key,
        owner: priceOwner(item),
        productName: item.product?.name || `商品#${item.product_id}`,
        customer_id: item.customer_id ?? null,
        price_level: item.price_level || '',
        product_id: item.product_id,
        items: [item],
      })
    }
  }
  return Array.from(map.values()).map((group) => ({
    ...group,
    items: [...group.items].sort((a, b) => Number(a.unit_spec?.factor_to_base || 0) - Number(b.unit_spec?.factor_to_base || 0)),
  }))
})

function reloadPrices(): void {
  pricePage.value = 1
  void qc.invalidateQueries({ queryKey: ['b2b-prices'] })
}

const priceDlg = ref(false)
const priceForm = reactive({
  customer_id: 0,
  price_level: '',
  product_id: undefined as number | undefined,
})

interface PriceLine {
  unit_spec_id: number
  unit_name: string
  factor_to_base: number
  sale_price: number
  supply_price: number
  min_quantity: number
  is_enabled: boolean
  remark: string
}

const priceLines = ref<PriceLine[]>([])

function openPriceCreate(): void {
  priceForm.customer_id = 0
  priceForm.price_level = ''
  priceForm.product_id = undefined
  priceLines.value = []
  priceDlg.value = true
}

function onPriceProductChange(): void {
  syncPriceLines()
}

function syncPriceLines(existing: B2BCustomerProductPrice[] = []): void {
  if (!priceForm.product_id) {
    priceLines.value = []
    return
  }
  const bySpec = new Map(existing.map((item) => [item.unit_spec_id, item]))
  priceLines.value = (specsByProduct.value.get(Number(priceForm.product_id)) ?? []).map((spec) => {
    const old = bySpec.get(spec.id)
    return {
      unit_spec_id: spec.id,
      unit_name: spec.unit_name || spec.unit_code,
      factor_to_base: Number(spec.factor_to_base || 1),
      sale_price: Number(spec.sale_price || 0),
      supply_price: Number(old?.supply_price || spec.sale_price || 0),
      min_quantity: Number(old?.min_quantity || 1),
      is_enabled: old?.is_enabled ?? true,
      remark: old?.remark || '',
    }
  })
}

function openPriceGroupEdit(row: PriceGroupRow): void {
  priceForm.customer_id = row.customer_id || 0
  priceForm.price_level = row.price_level || ''
  priceForm.product_id = row.product_id
  syncPriceLines(row.items)
  priceDlg.value = true
}

async function submitPrice(): Promise<void> {
  if (!priceForm.product_id) {
    toast.warning('请选择商品')
    return
  }
  if (!priceForm.customer_id && !priceForm.price_level.trim()) {
    toast.warning('请选择客户或填写价格等级')
    return
  }
  const lines = priceLines.value.filter((line) => Number(line.supply_price) > 0)
  if (!lines.length) {
    toast.warning('至少填写一个规格的供货价')
    return
  }
  saving.value = true
  try {
    await Promise.all(
      lines.map((line) =>
        upsertB2BPrice({
          customer_id: priceForm.customer_id || null,
          price_level: priceForm.price_level.trim(),
          product_id: priceForm.product_id,
          unit_spec_id: line.unit_spec_id,
          supply_price: line.supply_price,
          min_quantity: line.min_quantity,
          is_enabled: line.is_enabled,
          remark: line.remark,
        }),
      ),
    )
    toast.success(`已保存 ${lines.length} 个规格价格`)
    priceDlg.value = false
    reloadPrices()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function onDeletePrice(row: B2BCustomerProductPrice): Promise<void> {
  const ok = await confirmDialog({ message: '确定删除这条供货价？' })
  if (!ok) return
  await deleteB2BPrice(row.id)
  toast.success('已删除')
  reloadPrices()
}

const orderKeyword = ref('')
const paymentStatus = ref<number | ''>('')
const startDate = ref('')
const endDate = ref('')
const orderPage = ref(1)
const orderPageSize = ref(10)
const orderQueryKey = computed(() => ['b2b-orders', orderKeyword.value, paymentStatus.value, startDate.value, endDate.value, orderPage.value, orderPageSize.value] as const)
const { data: orderPageData, isLoading: orderLoading } = useQuery({
  queryKey: orderQueryKey,
  queryFn: () => listB2BSupplyOrders({ keyword: orderKeyword.value.trim() || undefined, payment_status: paymentStatus.value || undefined, start_date: startDate.value || undefined, end_date: endDate.value || undefined, page: orderPage.value, page_size: orderPageSize.value }),
})
const orders = computed(() => orderPageData.value?.list ?? [])
const orderTotal = computed(() => orderPageData.value?.total ?? 0)
const orderColumns: BaseTableColumn[] = [
  { key: 'order_no', label: '单号', prop: 'order_no', minWidth: '150px', ellipsis: true },
  { key: 'customer_name', label: '客户', prop: 'customer_name', minWidth: '140px', ellipsis: true },
  { key: 'order_date', label: '日期', prop: 'order_date', width: '120px' },
  { key: 'total_amount', label: '金额', width: '100px' },
  { key: 'paid_amount', label: '已收', width: '100px' },
  { key: 'unpaid_amount', label: '未收', width: '100px' },
  { key: 'profit_amount', label: '毛利', width: '100px' },
  { key: 'payment_status', label: '收款', width: '90px' },
  { key: 'delivery_status', label: '配送', width: '90px' },
  { key: 'actions', label: '操作', width: '90px', align: 'right' },
]

const orderSummary = computed(() => [
  { label: '本页供货额', value: money(orders.value.reduce((sum, x) => sum + Number(x.total_amount || 0), 0)) },
  { label: '本页已收', value: money(orders.value.reduce((sum, x) => sum + Number(x.paid_amount || 0), 0)) },
  { label: '本页未收', value: money(orders.value.reduce((sum, x) => sum + Number(x.unpaid_amount || 0), 0)) },
  { label: '本页毛利', value: money(orders.value.reduce((sum, x) => sum + Number(x.profit_amount || 0), 0)) },
])

interface OrderLine {
  product_id?: number
  unit_spec_id?: number
  quantity: number
  supply_price: number
}

const orderDlg = ref(false)
const orderForm = reactive({
  customer_id: undefined as number | undefined,
  order_date: new Date().toISOString().slice(0, 10),
  paid_amount: 0,
  remark: '',
  items: [] as OrderLine[],
})

const selectedCustomerId = computed(() => Number(orderForm.customer_id || 0))
const selectedCustomer = computed(() => customers.value.find((item) => item.id === selectedCustomerId.value))
const selectedPriceLevel = computed(() => (selectedCustomer.value?.price_level || '').trim())
const { data: orderCustomerPriceData, isFetching: orderCustomerPriceFetching } = useQuery({
  queryKey: computed(() => ['b2b-order-customer-prices', selectedCustomerId.value] as const),
  queryFn: () => listB2BPrices({ customer_id: selectedCustomerId.value, page: 1, page_size: 100 }),
  enabled: computed(() => selectedCustomerId.value > 0),
})
const { data: orderLevelPriceData, isFetching: orderLevelPriceFetching } = useQuery({
  queryKey: computed(() => ['b2b-order-level-prices', selectedPriceLevel.value] as const),
  queryFn: () => listB2BPrices({ price_level: selectedPriceLevel.value, page: 1, page_size: 100 }),
  enabled: computed(() => selectedPriceLevel.value !== ''),
})
const orderCustomerPrices = computed(() => orderCustomerPriceData.value?.list ?? [])
const orderLevelPrices = computed(() => orderLevelPriceData.value?.list ?? [])
const orderPriceLoading = computed(() => orderCustomerPriceFetching.value || orderLevelPriceFetching.value)
const availableOrderPrices = computed(() => {
  const bySpec = new Map<string, B2BCustomerProductPrice>()
  for (const item of orderLevelPrices.value) {
    if (!item.is_enabled || Number(item.supply_price || 0) <= 0) continue
    bySpec.set(`${item.product_id}:${item.unit_spec_id}`, item)
  }
  for (const item of orderCustomerPrices.value) {
    if (!item.is_enabled || Number(item.supply_price || 0) <= 0) continue
    bySpec.set(`${item.product_id}:${item.unit_spec_id}`, item)
  }
  return Array.from(bySpec.values())
})
const orderProductOptions = computed<BaseSelectOption[]>(() => {
  const map = new Map<number, string>()
  for (const item of availableOrderPrices.value) {
    map.set(item.product_id, item.product?.name || `商品#${item.product_id}`)
  }
  return Array.from(map.entries()).map(([value, label]) => ({ label, value }))
})

function reloadOrders(): void {
  orderPage.value = 1
  void qc.invalidateQueries({ queryKey: ['b2b-orders'] })
}

function openOrderCreate(): void {
  orderForm.customer_id = undefined
  orderForm.order_date = new Date().toISOString().slice(0, 10)
  orderForm.paid_amount = 0
  orderForm.remark = ''
  orderForm.items = [{ quantity: 1, supply_price: 0 }]
  orderDlg.value = true
}

function addOrderLine(): void {
  if (!orderForm.customer_id) {
    toast.warning('请先选择客户')
    return
  }
  if (!orderPriceLoading.value && !orderProductOptions.value.length) {
    toast.warning('该客户暂无可供货商品，请先配置供货价')
    return
  }
  orderForm.items.push({ quantity: 1, supply_price: 0 })
}

function removeOrderLine(idx: number): void {
  orderForm.items.splice(idx, 1)
}

function onOrderProductChange(idx: number): void {
  const line = orderForm.items[idx]
  if (!orderForm.customer_id) {
    line.product_id = undefined
    line.unit_spec_id = undefined
    line.supply_price = 0
    toast.warning('请先选择客户')
    return
  }
  if (!orderProductOptions.value.some((item) => item.value === line.product_id)) {
    line.product_id = undefined
    line.unit_spec_id = undefined
    line.supply_price = 0
    toast.warning('请选择供货价中已配置的商品')
    return
  }
  const first = lineSpecOptions(line)[0]
  line.unit_spec_id = first ? Number(first.value) : undefined
  fillLinePrice(idx)
}

function fillLinePrice(idx: number): void {
  const line = orderForm.items[idx]
  if (!line?.product_id || !line.unit_spec_id) return
  line.supply_price = resolveConfiguredSupplyPrice(line.product_id, line.unit_spec_id)
  if (line.supply_price <= 0) {
    toast.warning('该客户未配置这个规格的供货价')
  }
}

function lineSpecOptions(line: OrderLine): BaseSelectOption[] {
  if (!orderForm.customer_id || !line.product_id) return []
  return availableOrderPrices.value
    .filter((item) => item.product_id === line.product_id)
    .sort((a, b) => Number(a.unit_spec?.factor_to_base || 0) - Number(b.unit_spec?.factor_to_base || 0))
    .map((item) => ({
      label: `${item.unit_name || item.unit_spec?.unit_name || `规格#${item.unit_spec_id}`} / ${money(item.supply_price)}`,
      value: item.unit_spec_id,
    }))
}

function selectedOrderCustomer(): B2BCustomer | undefined {
  return selectedCustomer.value
}

function resolveConfiguredSupplyPrice(productID: number, unitSpecID: number): number {
  const customer = selectedOrderCustomer()
  if (!customer) return 0
  const customerPrice = orderCustomerPrices.value.find((item) => item.product_id === productID && item.unit_spec_id === unitSpecID && item.is_enabled)
  if (customerPrice) return Number(customerPrice.supply_price || 0)

  const levelPrice = orderLevelPrices.value.find((item) => item.product_id === productID && item.unit_spec_id === unitSpecID && item.is_enabled)
  return Number(levelPrice?.supply_price || 0)
}

function onOrderCustomerChange(): void {
  orderForm.items = [{ quantity: 1, supply_price: 0 }]
}

const orderTotalAmount = computed(() => orderForm.items.reduce((sum, line) => sum + Number(line.quantity || 0) * Number(line.supply_price || 0), 0))

async function submitOrder(): Promise<void> {
  if (!orderForm.customer_id) {
    toast.warning('请选择客户')
    return
  }
  const items = orderForm.items.filter((x) => x.product_id && x.unit_spec_id && x.quantity > 0).map((x) => ({
    product_id: x.product_id,
    unit_spec_id: x.unit_spec_id,
    quantity: x.quantity,
    supply_price: x.supply_price,
  }))
  if (!items.length) {
    toast.warning('请添加商品明细')
    return
  }
  if (!orderProductOptions.value.length) {
    toast.warning('该客户暂无可供货商品，请先配置供货价')
    return
  }
  if (items.some((x) => Number(x.supply_price || 0) <= 0)) {
    toast.warning('存在未配置供货价的商品规格，请先到供货价中配置')
    return
  }
  saving.value = true
  try {
    await createB2BSupplyOrder({
      customer_id: orderForm.customer_id,
      order_date: orderForm.order_date,
      paid_amount: orderForm.paid_amount,
      remark: orderForm.remark,
      items,
    })
    toast.success('供货单已创建，库存已扣减')
    orderDlg.value = false
    reloadOrders()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '提交失败')
  } finally {
    saving.value = false
  }
}

const detailDlg = ref(false)
const detailJson = ref('')
async function openOrderDetail(row: B2BSupplyOrder): Promise<void> {
  const data = await getB2BSupplyOrder(row.id)
  detailJson.value = JSON.stringify(data, null, 2)
  detailDlg.value = true
}

function priceOwner(row: B2BCustomerProductPrice): string {
  if (row.customer?.name) return row.customer.name
  if (row.price_level) return `等级：${row.price_level}`
  return '通用'
}

function money(v: number | string | undefined): string {
  return `￥${Number(v || 0).toFixed(2)}`
}

function settlementLabel(v?: string): string {
  if (v === 'week') return '周结'
  if (v === 'month') return '月结'
  return '现结'
}

function paymentStatusLabel(v: number): string {
  if (v === 3) return '已收'
  if (v === 2) return '部分'
  return '未收'
}

function deliveryStatusLabel(v: number): string {
  if (v === 2) return '已配送'
  if (v === 3) return '已取消'
  return '待配送'
}

watch(tab, () => {
  if (tab.value === 'customers') reloadCustomers()
  if (tab.value === 'prices') reloadPrices()
  if (tab.value === 'orders') reloadOrders()
})
</script>

<style scoped>
.price-spec-grid {
  display: grid;
  grid-template-columns: minmax(100px, 1.2fr) minmax(90px, 0.8fr) minmax(70px, 0.6fr) minmax(70px, 0.6fr) minmax(70px, 0.7fr);
  align-items: center;
  overflow: hidden;
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
}

.price-spec-head,
.price-spec-cell {
  min-height: 36px;
  padding: 8px 10px;
  border-bottom: 1px solid var(--color-border-2);
}

.price-spec-head {
  background: var(--color-fill-2);
  color: var(--color-text-2);
  font-size: 13px;
  font-weight: 600;
}

.price-spec-cell {
  background: var(--color-bg-2);
}

.price-spec-grid > :nth-last-child(-n + 5) {
  border-bottom: 0;
}
</style>
