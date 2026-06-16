<template>
  <div class="flex min-h-0 flex-1 flex-col gap-4">
    <BaseCard flush-body class="flex min-h-0 flex-1 flex-col">
      <template #header>
        <div class="flex w-full flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h2 class="m-0 text-base font-semibold text-slate-900">门店返厂管理</h2>
            <p class="m-0 mt-1 text-xs text-slate-500">维护返厂商品押金、返厂日期与货拉拉费用</p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <BaseInput v-model="filters.keyword" class="w-48" placeholder="单号 / 商品 / 备注" clearable @enter="reloadAll" />
            <a-date-picker v-model="filters.start_date" value-format="YYYY-MM-DD" class="w-36" />
            <a-date-picker v-model="filters.end_date" value-format="YYYY-MM-DD" class="w-36" />
            <BaseButton variant="primary" @click="reloadAll">查询</BaseButton>
            <BaseButton v-permission="'store:return:add'" variant="primary" @click="openCreate">新增返厂</BaseButton>
          </div>
        </div>
      </template>

      <div class="flex min-h-0 flex-1 flex-col gap-3 p-4">
        <div class="flex gap-2 border-b border-slate-200">
          <button :class="tabClass('records')" type="button" @click="activeTab = 'records'">返厂记录</button>
          <button :class="tabClass('products')" type="button" @click="activeTab = 'products'">返厂商品</button>
        </div>

        <template v-if="activeTab === 'records'">
          <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
            <div v-for="item in summaryCards" :key="item.label"
              class="rounded border border-slate-200 bg-slate-50 px-4 py-3">
              <div class="text-xs font-medium text-slate-500">{{ item.label }}</div>
              <div class="mt-1 text-lg font-semibold text-slate-900">{{ item.value }}</div>
            </div>
          </div>

          <BaseTable :columns="columns" :data="(rows as unknown) as Record<string, unknown>[]" :loading="loading"
            min-width="1080px" class="min-h-0 flex-1">
            <template #cell-store="{ row }">
              {{ (row as StoreReturn).store?.name || '-' }}
            </template>
            <template #cell-return_date="{ row }">
              {{ formatDate((row as StoreReturn).return_date) }}
            </template>
            <template #cell-total_deposit="{ row }">
              {{ formatMoney((row as StoreReturn).total_deposit) }}
            </template>
            <template #cell-logistics_fee="{ row }">
              {{ formatMoney((row as StoreReturn).logistics_fee) }}
            </template>
            <template #cell-created_at="{ row }">
              {{ formatDateTime((row as StoreReturn).created_at) }}
            </template>
            <template #cell-actions="{ row }">
              <div class="flex justify-end gap-3" @click.stop>
                <BaseButton variant="link" size="sm" @click="openDetail(row as StoreReturn)">详情</BaseButton>
                <BaseButton v-if="isReturnEditable(row as StoreReturn)" v-permission="'store:return:edit'"
                  variant="link" size="sm" @click="openEdit(row as StoreReturn)">编辑</BaseButton>
                <BaseButton v-if="isReturnEditable(row as StoreReturn)" v-permission="'store:return:delete'"
                  variant="link" size="sm" @click="onDelete(row as StoreReturn)">删除</BaseButton>
              </div>
            </template>
          </BaseTable>

          <div class="flex shrink-0 justify-end">
            <BasePagination :page="page" :page-size="pageSize" :total="total" @update:page="(p) => (page = p)"
              @update:page-size="(s) => (pageSize = s)" />
          </div>
        </template>

        <template v-else>
          <div class="flex flex-wrap items-center justify-between gap-2">
            <BaseInput v-model="productFilters.keyword" class="w-56" placeholder="商品名称 / 备注" clearable
              @enter="reloadProducts" />
            <BaseButton variant="primary" @click="openProductCreate">新增商品</BaseButton>
          </div>
          <BaseTable :columns="productColumns" :data="(productRows as unknown) as Record<string, unknown>[]"
            :loading="productsLoading" min-width="760px" class="min-h-0 flex-1">
            <template #cell-store="{ row }">{{ (row as StoreReturnProduct).store?.name || '-' }}</template>
            <template #cell-deposit="{ row }">{{ formatMoney((row as StoreReturnProduct).deposit) }}</template>
            <template #cell-status="{ row }">{{ (row as StoreReturnProduct).status === 1 ? '启用' : '停用' }}</template>
            <template #cell-actions="{ row }">
              <div class="flex justify-end gap-3" @click.stop>
                <BaseButton variant="link" size="sm" @click="openProductEdit(row as StoreReturnProduct)">编辑</BaseButton>
                <BaseButton variant="link" size="sm" @click="onProductDelete(row as StoreReturnProduct)">删除</BaseButton>
              </div>
            </template>
          </BaseTable>
          <div class="flex shrink-0 justify-end">
            <BasePagination :page="productPage" :page-size="productPageSize" :total="productTotal"
              @update:page="(p) => (productPage = p)" @update:page-size="(s) => (productPageSize = s)" />
          </div>
        </template>
      </div>
    </BaseCard>

    <BaseDialog v-model="formDlg" :title="editingId ? '编辑返厂记录' : '新增返厂记录'" max-width="min(760px, 96vw)">
      <div class="space-y-3">
        <BaseFormItem label="返厂日期" required>
          <a-date-picker v-model="form.return_date" value-format="YYYY-MM-DD" class="w-full" />
        </BaseFormItem>

        <div class="rounded border border-[var(--color-border-2)] p-3">
          <div class="mb-3 flex items-center justify-between">
            <h3 class="m-0 text-sm font-semibold text-slate-900">返厂商品</h3>
            <BaseButton variant="secondary" size="sm" @click="addLine">添加商品</BaseButton>
          </div>
          <div class="return-line-editor">
            <div class="return-line-editor__head">
              <span>商品名称</span>
              <span>数量</span>
              <span>押金小计</span>
              <span>操作</span>
            </div>
            <div v-for="(line, idx) in form.items" :key="idx" class="return-line-editor__row">
              <BaseSelect v-model="line.product_id" class="return-line-editor__control" :options="returnProductOptions"
                placeholder="选择返厂商品" @update:model-value="onLineProductChange(idx)" />
              <BaseNumberInput v-model="line.quantity" class="return-line-editor__control" :min="0.01" :step="1" />
              <div class="return-line-editor__subtotal">{{ formatMoney(line.deposit * line.quantity) }}</div>
              <BaseButton variant="ghost" size="sm" :disabled="form.items.length <= 1" @click="removeLine(idx)">移除
              </BaseButton>
            </div>
          </div>
        </div>

        <BaseFormItem label="货拉拉费用">
          <BaseNumberInput v-model="form.logistics_fee" :min="0" :step="0.01" />
        </BaseFormItem>
        <div class="return-total-row">
          <span>押金合计</span>
          <strong>{{ formatMoney(formTotalDeposit) }}</strong>
        </div>
        <BaseFormItem label="整单备注">
          <BaseTextarea v-model="form.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="formDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitForm">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="productDlg" :title="editingProductId ? '编辑返厂商品' : '新增返厂商品'" max-width="min(520px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="商品名称" required>
          <BaseInput v-model="productForm.product_name" placeholder="请输入商品名称" />
        </BaseFormItem>
        <BaseFormItem label="默认押金" required>
          <BaseNumberInput v-model="productForm.deposit" :min="0" :step="0.01" />
        </BaseFormItem>
        <BaseFormItem label="状态">
          <BaseSelect v-model="productForm.status" :options="statusOptions" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="productForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="productDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="savingProduct" @click="submitProduct">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="detailDlg" title="返厂详情" max-width="min(820px, 96vw)">
      <div v-if="detail" class="space-y-4">
        <div class="grid grid-cols-1 gap-3 rounded border border-[var(--color-border-2)] p-3 text-sm md:grid-cols-3">
          <div><span class="text-slate-500">单号：</span>{{ detail.return_no }}</div>
          <div><span class="text-slate-500">门店：</span>{{ detail.store?.name || '-' }}</div>
          <div><span class="text-slate-500">返厂日期：</span>{{ formatDate(detail.return_date) }}</div>
          <div><span class="text-slate-500">押金合计：</span>{{ formatMoney(detail.total_deposit) }}</div>
          <div><span class="text-slate-500">货拉拉费用：</span>{{ formatMoney(detail.logistics_fee) }}</div>
          <div><span class="text-slate-500">操作人：</span>{{ detail.operator_name || '-' }}</div>
          <div class="md:col-span-3"><span class="text-slate-500">备注：</span>{{ detail.remark || '-' }}</div>
        </div>
        <BaseTable :columns="itemColumns" :data="(detail.items ?? []) as unknown as Record<string, unknown>[]"
          min-width="560px">
          <template #cell-subtotal="{ row }">{{ formatMoney((row as StoreReturnItem).deposit * (row as
            StoreReturnItem).quantity) }}</template>
        </BaseTable>
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
  BaseCard,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseNumberInput,
  BasePagination,
  BaseSelect,
  BaseTable,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn } from '@/components/base/types'
import {
  createStoreReturnProduct,
  createStoreReturn,
  deleteStoreReturnProduct,
  deleteStoreReturn,
  getStoreReturn,
  getStoreReturnStats,
  listStoreReturnProducts,
  listStoreReturns,
  updateStoreReturnProduct,
  updateStoreReturn,
} from '@/api/storeReturn'
import type { StoreReturn, StoreReturnItem, StoreReturnProduct } from '@/api/types'
import { confirmDialog } from '@/feedback/confirm'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

type ReturnLine = {
  product_id: number | ''
  product_name?: string
  quantity: number
  deposit: number
  remark: string
}

const qc = useQueryClient()
const userStore = useUserStore()
const tenantStoreId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)
const activeTab = ref<'records' | 'products'>('records')

const filters = reactive({
  keyword: '',
  start_date: '',
  end_date: '',
})
const page = ref(1)
const pageSize = ref(10)
const productPage = ref(1)
const productPageSize = ref(10)
const productFilters = reactive({ keyword: '' })
const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  store_id: tenantStoreId.value,
  keyword: filters.keyword.trim() || undefined,
  start_date: filters.start_date || undefined,
  end_date: filters.end_date || undefined,
}))
const queryKey = computed(() => ['store-returns', queryParams.value] as const)
const { data: pageData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: () => listStoreReturns(queryParams.value),
})
const { data: statsData } = useQuery({
  queryKey: computed(() => ['store-return-stats', tenantStoreId.value, filters.start_date, filters.end_date] as const),
  queryFn: () =>
    getStoreReturnStats({
      store_id: tenantStoreId.value,
      start_date: filters.start_date || undefined,
      end_date: filters.end_date || undefined,
    }),
})
const productQueryParams = computed(() => ({
  page: productPage.value,
  page_size: productPageSize.value,
  store_id: tenantStoreId.value,
  keyword: productFilters.keyword.trim() || undefined,
}))
const { data: productsPageData, isLoading: productsLoading } = useQuery({
  queryKey: computed(() => ['store-return-products', productQueryParams.value] as const),
  queryFn: () => listStoreReturnProducts(productQueryParams.value),
})
const { data: enabledProductsPage } = useQuery({
  queryKey: computed(() => ['store-return-products-enabled', tenantStoreId.value] as const),
  queryFn: () => listStoreReturnProducts({ page: 1, page_size: 500, store_id: tenantStoreId.value, status: 1 }),
})

const rows = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)
const productRows = computed(() => productsPageData.value?.list ?? [])
const productTotal = computed(() => productsPageData.value?.total ?? 0)
const returnProducts = computed(() => enabledProductsPage.value?.list ?? [])
const returnProductOptions = computed<BaseSelectOption[]>(() =>
  returnProducts.value.map((p) => ({ label: `${p.product_name}（${formatMoney(p.deposit)}）`, value: p.id })),
)
const returnProductMap = computed(() => {
  const map = new Map<number, StoreReturnProduct>()
  for (const product of returnProducts.value) map.set(product.id, product)
  return map
})
const summaryCards = computed(() => {
  const stats = statsData.value
  return [
    { label: '返厂单数', value: `${stats?.return_count ?? 0} 单` },
    { label: '商品明细', value: `${stats?.item_count ?? 0} 条` },
    { label: '押金总额', value: formatMoney(stats?.total_deposit ?? 0) },
    { label: '货拉拉费用', value: formatMoney(stats?.logistics_fee ?? 0) },
  ]
})

const columns: BaseTableColumn[] = [
  { key: 'return_no', label: '返厂单号', prop: 'return_no', minWidth: '150px', ellipsis: true },
  { key: 'store', label: '门店', minWidth: '140px', ellipsis: true },
  { key: 'return_date', label: '返厂日期', width: '120px' },
  { key: 'item_count', label: '商品数', prop: 'item_count', width: '80px', align: 'right' },
  { key: 'total_deposit', label: '押金合计', width: '120px', align: 'right' },
  { key: 'logistics_fee', label: '货拉拉费用', width: '120px', align: 'right' },
  { key: 'operator_name', label: '操作人', prop: 'operator_name', width: '110px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '160px', ellipsis: true },
  { key: 'created_at', label: '创建时间', width: '170px' },
  { key: 'actions', label: '操作', width: '150px', align: 'right' },
]

const itemColumns: BaseTableColumn[] = [
  { key: 'product_name', label: '商品名称', prop: 'product_name', minWidth: '220px', ellipsis: true },
  { key: 'quantity', label: '数量', prop: 'quantity', width: '90px', align: 'right' },
  { key: 'subtotal', label: '押金小计', width: '120px', align: 'right' },
]
const productColumns: BaseTableColumn[] = [
  { key: 'product_name', label: '商品名称', prop: 'product_name', minWidth: '220px', ellipsis: true },
  { key: 'store', label: '门店', minWidth: '140px', ellipsis: true },
  { key: 'deposit', label: '默认押金', width: '120px', align: 'right' },
  { key: 'status', label: '状态', width: '80px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '180px', ellipsis: true },
  { key: 'actions', label: '操作', width: '120px', align: 'right' },
]
const statusOptions: BaseSelectOption[] = [
  { label: '启用', value: 1 },
  { label: '停用', value: 0 },
]

const formDlg = ref(false)
const productDlg = ref(false)
const detailDlg = ref(false)
const saving = ref(false)
const savingProduct = ref(false)
const editingId = ref<number | null>(null)
const editingProductId = ref<number | null>(null)
const detail = ref<StoreReturn | null>(null)
const form = reactive({
  client_request_id: '',
  return_date: today(),
  logistics_fee: 0,
  remark: '',
  items: [] as ReturnLine[],
})
const productForm = reactive({
  product_name: '',
  deposit: 0,
  remark: '',
  status: 1,
})
const formTotalDeposit = computed(() =>
  form.items.reduce((sum, x) => sum + Number(x.deposit || 0) * Number(x.quantity || 0), 0),
)

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['store-returns'] })
})
watch([productPage, productPageSize], () => {
  void qc.invalidateQueries({ queryKey: ['store-return-products'] })
})

function reloadAll(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['store-returns'] })
  void qc.invalidateQueries({ queryKey: ['store-return-stats'] })
}

function reloadProducts(): void {
  productPage.value = 1
  void qc.invalidateQueries({ queryKey: ['store-return-products'] })
  void qc.invalidateQueries({ queryKey: ['store-return-products-enabled'] })
}

function openCreate(): void {
  editingId.value = null
  resetForm()
  form.client_request_id = createClientRequestId()
  formDlg.value = true
}

function openEdit(row: StoreReturn): void {
  editingId.value = row.id
  form.client_request_id = row.client_request_id || ''
  form.return_date = normalizeDate(row.return_date)
  form.logistics_fee = Number(row.logistics_fee || 0)
  form.remark = row.remark || ''
  form.items = (row.items?.length ? row.items : [{ product_id: 0, product_name: '', quantity: 1, deposit: 0, remark: '' }]).map((x) => ({
    product_id: x.product_id || '',
    product_name: x.product_name || '',
    quantity: Number(x.quantity || 1),
    deposit: Number(x.deposit || 0),
    remark: x.remark || '',
  }))
  formDlg.value = true
}

async function openDetail(row: StoreReturn): Promise<void> {
  detail.value = await getStoreReturn(row.id)
  detailDlg.value = true
}

function addLine(): void {
  form.items.push({ product_id: '', product_name: '', quantity: 1, deposit: 0, remark: '' })
}

function removeLine(idx: number): void {
  if (form.items.length <= 1) return
  form.items.splice(idx, 1)
}

async function submitForm(): Promise<void> {
  if (saving.value) return
  const items = form.items
    .map((x) => ({
      product_id: Number(x.product_id || 0),
      product_name: x.product_name?.trim() || '',
      quantity: Number(x.quantity || 0),
      deposit: Number(x.deposit || 0),
      remark: x.remark.trim(),
    }))
    .filter((x) => x.product_id > 0)
  if (!form.return_date) {
    toast.error('请选择返厂日期')
    return
  }
  if (!items.length) {
    toast.error('请至少选择一个返厂商品')
    return
  }
  if (items.some((x) => x.quantity <= 0)) {
    toast.error('返厂商品数量必须大于0')
    return
  }

  saving.value = true
  try {
    const body = {
      store_id: tenantStoreId.value,
      client_request_id: form.client_request_id || undefined,
      return_date: form.return_date,
      logistics_fee: Number(form.logistics_fee || 0),
      remark: form.remark.trim(),
      items,
    }
    if (editingId.value) {
      await updateStoreReturn(editingId.value, body)
      toast.success('返厂记录已更新')
    } else {
      await createStoreReturn(body)
      toast.success('返厂记录已创建')
    }
    formDlg.value = false
    reloadAll()
  } finally {
    saving.value = false
  }
}

function onLineProductChange(idx: number): void {
  const line = form.items[idx]
  const productID = Number(line.product_id || 0)
  const product = returnProductMap.value.get(productID)
  if (!product) {
    line.product_name = ''
    line.deposit = 0
    return
  }
  line.product_name = product.product_name
  line.deposit = Number(product.deposit || 0)
}

function openProductCreate(): void {
  editingProductId.value = null
  productForm.product_name = ''
  productForm.deposit = 0
  productForm.remark = ''
  productForm.status = 1
  productDlg.value = true
}

function openProductEdit(row: StoreReturnProduct): void {
  editingProductId.value = row.id
  productForm.product_name = row.product_name || ''
  productForm.deposit = Number(row.deposit || 0)
  productForm.remark = row.remark || ''
  productForm.status = Number(row.status ?? 1)
  productDlg.value = true
}

async function submitProduct(): Promise<void> {
  const name = productForm.product_name.trim()
  if (!name) {
    toast.error('请输入商品名称')
    return
  }
  savingProduct.value = true
  try {
    const body = {
      store_id: tenantStoreId.value,
      product_name: name,
      deposit: Number(productForm.deposit || 0),
      remark: productForm.remark.trim(),
      status: Number(productForm.status),
    }
    if (editingProductId.value) {
      await updateStoreReturnProduct(editingProductId.value, body)
      toast.success('返厂商品已更新')
    } else {
      await createStoreReturnProduct(body)
      toast.success('返厂商品已创建')
    }
    productDlg.value = false
    reloadProducts()
  } finally {
    savingProduct.value = false
  }
}

async function onProductDelete(row: StoreReturnProduct): Promise<void> {
  const ok = await confirmDialog({ message: `删除返厂商品「${row.product_name}」？` })
  if (!ok) return
  await deleteStoreReturnProduct(row.id)
  toast.success('已删除')
  reloadProducts()
}

async function onDelete(row: StoreReturn): Promise<void> {
  const ok = await confirmDialog({
    title: '删除返厂记录',
    message: `确认删除返厂单 ${row.return_no} 吗？`,
  })
  if (!ok) return
  await deleteStoreReturn(row.id)
  toast.success('已删除')
  reloadAll()
}

function resetForm(): void {
  form.client_request_id = ''
  form.return_date = today()
  form.logistics_fee = 0
  form.remark = ''
  form.items = [{ product_id: '', product_name: '', quantity: 1, deposit: 0, remark: '' }]
}

function tabClass(tab: 'records' | 'products'): string {
  const active = activeTab.value === tab
  return [
    'border-b-2 px-3 py-2 text-sm font-medium transition',
    active ? 'border-indigo-600 text-indigo-700' : 'border-transparent text-slate-500 hover:text-slate-800',
  ].join(' ')
}

function today(): string {
  const d = new Date()
  const m = `${d.getMonth() + 1}`.padStart(2, '0')
  const day = `${d.getDate()}`.padStart(2, '0')
  return `${d.getFullYear()}-${m}-${day}`
}

function createClientRequestId(): string {
  const cryptoObj = globalThis.crypto
  if (cryptoObj?.randomUUID) return cryptoObj.randomUUID()
  return `sr_${Date.now()}_${Math.random().toString(16).slice(2)}`
}

function normalizeDate(value?: string): string {
  if (!value) return today()
  return value.slice(0, 10)
}

function formatDate(value?: string): string {
  if (!value) return '-'
  return value.slice(0, 10)
}

function formatDateTime(value?: string): string {
  if (!value) return '-'
  return value.replace('T', ' ').slice(0, 19)
}

function isReturnEditable(row: StoreReturn): boolean {
  if (!row.created_at) return false
  const created = new Date(row.created_at)
  if (Number.isNaN(created.getTime())) return false
  const now = new Date()
  return (
    created.getFullYear() === now.getFullYear() &&
    created.getMonth() === now.getMonth() &&
    created.getDate() === now.getDate()
  )
}

function formatMoney(value?: number): string {
  return `¥${Number(value || 0).toFixed(2)}`
}
</script>

<style scoped>
.return-line-editor {
  display: grid;
  gap: 8px;
  overflow-x: auto;
}

.return-line-editor__head,
.return-line-editor__row {
  display: grid;
  grid-template-columns: minmax(220px, 1fr) minmax(90px, 0.35fr) minmax(120px, 0.45fr) 72px;
  gap: 8px;
  align-items: center;
  min-width: 560px;
}

.return-line-editor__head {
  color: var(--color-text-3);
  font-size: 12px;
  font-weight: 600;
}

.return-line-editor__control {
  min-width: 0;
}

.return-line-editor__subtotal {
  min-height: 32px;
  border-radius: 4px;
  background: #f8fafc;
  padding: 6px 10px;
  text-align: right;
  font-weight: 600;
  color: #334155;
}

.return-total-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-radius: 6px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  padding: 10px 12px;
}

.return-total-row span {
  font-size: 14px;
  font-weight: 600;
  color: #991b1b;
}

.return-total-row strong {
  font-size: 22px;
  line-height: 1;
  color: #dc2626;
}
</style>
