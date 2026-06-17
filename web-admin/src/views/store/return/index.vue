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
import { useRouter } from 'vue-router'
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
  deleteStoreReturnProduct,
  deleteStoreReturn,
  getStoreReturn,
  getStoreReturnStats,
  listStoreReturnProducts,
  listStoreReturns,
  updateStoreReturnProduct,
} from '@/api/storeReturn'
import type { StoreReturn, StoreReturnItem, StoreReturnProduct } from '@/api/types'
import { confirmDialog } from '@/feedback/confirm'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

const qc = useQueryClient()
const router = useRouter()
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
const rows = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)
const productRows = computed(() => productsPageData.value?.list ?? [])
const productTotal = computed(() => productsPageData.value?.total ?? 0)
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

const productDlg = ref(false)
const detailDlg = ref(false)
const savingProduct = ref(false)
const editingProductId = ref<number | null>(null)
const detail = ref<StoreReturn | null>(null)
const productForm = reactive({
  product_name: '',
  deposit: 0,
  remark: '',
  status: 1,
})

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
}

function openCreate(): void {
  void router.push('/store/return/form')
}

function openEdit(row: StoreReturn): void {
  void router.push({ path: '/store/return/form', query: { id: row.id } })
}

async function openDetail(row: StoreReturn): Promise<void> {
  detail.value = await getStoreReturn(row.id)
  detailDlg.value = true
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

function tabClass(tab: 'records' | 'products'): string {
  const active = activeTab.value === tab
  return [
    'border-b-2 px-3 py-2 text-sm font-medium transition',
    active ? 'border-indigo-600 text-indigo-700' : 'border-transparent text-slate-500 hover:text-slate-800',
  ].join(' ')
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
