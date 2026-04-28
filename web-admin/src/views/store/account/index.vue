<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">门店记账</h2>
      <div class="flex flex-col sm:flex-row flex-wrap gap-2 w-full md:w-auto">
        <BaseInput v-model="rangeStart" class="w-full sm:w-40" type="date" />
        <BaseInput v-model="rangeEnd" class="w-full sm:w-40" type="date" />
        <BaseButton variant="primary" @click="reloadAll">查询</BaseButton>
        <BaseButton v-permission="'store:account:add'" variant="primary" @click="openCreate">快速记账</BaseButton>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">区间内金额合计</span></template>
        <p class="text-2xl font-semibold text-indigo-700 m-0">{{ statsAmount }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">记账笔数</span></template>
        <p class="text-2xl font-semibold text-slate-800 m-0">{{ statsCount }}</p>
      </BaseCard>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading" min-width="960px">
      <template #cell-account_date="{ row }">
        {{ formatDate((row as StoreAccount).account_date) }}
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
          <BaseButton v-permission="'store:account:list'" variant="link" size="sm" @click="openView(row as StoreAccount)">详情</BaseButton>
          <BaseButton v-permission="'store:account:edit'" variant="link" size="sm" @click="openEdit(row as StoreAccount)">编辑</BaseButton>
          <BaseButton v-permission="'store:account:delete'" variant="link" size="sm" @click="onDelete(row as StoreAccount)">删除</BaseButton>
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

    <BaseDialog v-model="createDlg" title="快速记账（单商品）" max-width="min(480px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="渠道" required>
          <BaseInput v-model="cForm.channel" placeholder="字典 sales_channel 等" />
        </BaseFormItem>
        <BaseFormItem label="记账日期">
          <BaseInput v-model="cForm.account_date" type="date" />
        </BaseFormItem>
        <BaseFormItem label="订单号">
          <BaseInput v-model="cForm.order_no" />
        </BaseFormItem>
        <BaseFormItem label="商品" required>
          <a-cascader
            v-model="cForm.product_path"
            :options="productCascaderOptions"
            placeholder="先选分类，再选商品"
            allow-clear
            :path-mode="true"
            :check-strictly="false"
            @change="onCreateProductChange"
          />
        </BaseFormItem>
        <BaseFormItem label="数量" required>
          <BaseNumberInput v-model="cForm.quantity" :min="0.01" :step="0.01" />
        </BaseFormItem>
        <BaseFormItem label="单位">
          <BaseSelect
            v-model="cForm.unit"
            :options="createUnitOptions"
            :disabled="createUnitOptions.length <= 1"
            placeholder="单位"
          />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="cForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="createDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitCreate">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="editDlg" title="编辑记账" max-width="min(440px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="渠道">
          <BaseInput v-model="eForm.channel" />
        </BaseFormItem>
        <BaseFormItem label="订单号">
          <BaseInput v-model="eForm.order_no" />
        </BaseFormItem>
        <BaseFormItem label="标签编码">
          <BaseInput v-model="eForm.tag_code" />
        </BaseFormItem>
        <BaseFormItem label="标签名称">
          <BaseInput v-model="eForm.tag_name" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="eForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="editDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitEdit">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="viewDlg" title="记账详情" max-width="min(560px, 96vw)">
      <pre v-if="viewJson" class="text-xs overflow-auto max-h-[60vh] m-0 p-3 rounded bg-[var(--color-fill-2)]">{{ viewJson }}</pre>
      <template #footer>
        <BaseButton variant="ghost" @click="viewDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
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
import type { BaseTableColumn } from '@/components/base/types'
import {
  createStoreAccount,
  deleteStoreAccount,
  getStoreAccount,
  getStoreAccountStats,
  listStoreAccounts,
  updateStoreAccount,
} from '@/api/storeAccount'
import { listDictDataByTypeCode } from '@/api/dict'
import { listProductUnitSpecs } from '@/api/supplierProduct'
import { listPurchasableProducts } from '@/api/storeSupplier'
import type { DictData, ProductUnitSpec, StoreAccount } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'
import { useUserStore } from '@/store/user'

const qc = useQueryClient()
const userStore = useUserStore()
const tenantStoreId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)

function monthRange(): { start: string; end: string } {
  const t = new Date()
  const y = t.getFullYear()
  const m = String(t.getMonth() + 1).padStart(2, '0')
  const d = String(t.getDate()).padStart(2, '0')
  return { start: `${y}-${m}-01`, end: `${y}-${m}-${d}` }
}

const r = monthRange()
const rangeStart = ref(r.start)
const rangeEnd = ref(r.end)

const stats = ref<{ total_amount?: number; count?: number }>({})
const statsAmount = computed(() => (stats.value.total_amount ?? 0).toFixed(2))
const statsCount = computed(() => String(stats.value.count ?? 0))

async function loadStats(): Promise<void> {
  try {
    stats.value = await getStoreAccountStats({
      store_id: tenantStoreId.value,
      start_date: rangeStart.value,
      end_date: rangeEnd.value,
    })
  } catch {
    stats.value = {}
  }
}

onMounted(() => {
  void loadStats()
})

const page = ref(1)
const pageSize = ref(10)
const listKey = computed(
  () => ['store-accounts', tenantStoreId.value, page.value, pageSize.value, rangeStart.value, rangeEnd.value] as const,
)

const { data: pageData, isLoading: loading } = useQuery({
  queryKey: listKey,
  queryFn: () =>
    listStoreAccounts({
      page: page.value,
      page_size: pageSize.value,
      store_id: tenantStoreId.value,
      start_date: rangeStart.value,
      end_date: rangeEnd.value,
    }),
})

const list = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)

const { data: productData } = useQuery({
  queryKey: computed(() => ['store-account-products', tenantStoreId.value] as const),
  queryFn: () =>
    listPurchasableProducts({
      store_id: tenantStoreId.value,
    }),
})
const productList = computed(() => productData.value ?? [])
const productById = computed(() => {
  const map = new Map<number, (typeof productList.value)[number]>()
  for (const p of productList.value) map.set(p.id, p)
  return map
})
const productIdsKey = computed(() =>
  productList.value
    .map((p) => p.id)
    .sort((a, b) => a - b)
    .join(','),
)
const { data: unitSpecsData } = useQuery({
  queryKey: computed(() => ['store-account-product-unit-specs', productIdsKey.value] as const),
  queryFn: async () => {
    const ids = productList.value.map((p) => p.id)
    if (!ids.length) return [] as ProductUnitSpec[]
    const rows = await Promise.all(ids.map((id) => listProductUnitSpecs(id)))
    return rows.flat()
  },
  enabled: computed(() => productList.value.length > 0),
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
const { data: unitData } = useQuery({
  queryKey: ['dict-data', 'product_unit'],
  queryFn: () => listDictDataByTypeCode('product_unit'),
})
const unitDict = computed(() => unitData.value ?? ([] as DictData[]))

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

function reloadAll(): void {
  page.value = 1
  void loadStats()
  void qc.invalidateQueries({ queryKey: ['store-accounts'] })
}

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['store-accounts'] })
})

watch(
  () => tenantStoreId.value,
  () => {
    void loadStats()
    void qc.invalidateQueries({ queryKey: ['store-accounts'] })
  },
)

const columns: BaseTableColumn[] = [
  { key: 'account_no', label: '记账编号', prop: 'account_no', minWidth: '140px', ellipsis: true },
  { key: 'channel', label: '渠道', prop: 'channel', width: '100px' },
  { key: 'order_no', label: '订单号', prop: 'order_no', minWidth: '120px', ellipsis: true },
  { key: 'total_amount', label: '金额', prop: 'total_amount', width: '96px' },
  { key: 'account_date', label: '日期', width: '120px' },
  { key: 'actions', label: '操作', width: '220px', align: 'right' },
]

function formatDate(v: string): string {
  if (!v) return '-'
  return String(v).slice(0, 10)
}

const createDlg = ref(false)
const saving = ref(false)
const cForm = reactive({
  channel: '',
  account_date: '',
  order_no: '',
  product_path: [] as Array<string | number>,
  quantity: 1,
  unit: '',
  remark: '',
})

function openCreate(): void {
  const t = new Date()
  cForm.channel = ''
  cForm.account_date = `${t.getFullYear()}-${String(t.getMonth() + 1).padStart(2, '0')}-${String(t.getDate()).padStart(2, '0')}`
  cForm.order_no = ''
  cForm.product_path = []
  cForm.quantity = 1
  cForm.unit = ''
  cForm.remark = ''
  createDlg.value = true
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

const createUnitOptions = computed(() => {
  const pid = getProductId(cForm.product_path)
  if (!pid) return []
  const specs = specsByProduct.value.get(pid) ?? []
  if (specs.length > 0) {
    return specs.map((s) => ({
      label: s.unit_name,
      value: s.unit_code,
    }))
  }
  const product = productById.value.get(pid)
  const defaultUnit = product?.unit || unitDict.value[0]?.value || '件'
  return [{ label: defaultUnit, value: defaultUnit }]
})

function onCreateProductChange(): void {
  const options = createUnitOptions.value
  cForm.unit = String(options[0]?.value || '')
}

async function submitCreate(): Promise<void> {
  if (!cForm.channel.trim()) {
    toast.warning('请填写渠道')
    return
  }
  const productId = getProductId(cForm.product_path)
  if (!productId) {
    toast.warning('请选择商品')
    return
  }
  if (cForm.quantity <= 0) {
    toast.warning('数量必须大于0')
    return
  }
  if (!cForm.unit.trim()) {
    toast.warning('请选择单位')
    return
  }
  saving.value = true
  try {
    await createStoreAccount({
      store_id: tenantStoreId.value,
      channel: cForm.channel.trim(),
      order_no: cForm.order_no.trim(),
      remark: cForm.remark.trim(),
      account_date: cForm.account_date || undefined,
      other_expense_amount: 0,
      items: [
        {
          product_id: productId,
          quantity: cForm.quantity,
          unit: cForm.unit.trim(),
          spec: '',
          price: 0,
          amount: 0,
          remark: '',
        },
      ],
    })
    toast.success('已保存')
    createDlg.value = false
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

const editDlg = ref(false)
const editId = ref(0)
const eForm = reactive({
  channel: '',
  order_no: '',
  tag_code: '',
  tag_name: '',
  remark: '',
})

function openEdit(row: StoreAccount): void {
  editId.value = row.id
  eForm.channel = row.channel ?? ''
  eForm.order_no = row.order_no ?? ''
  eForm.tag_code = row.tag_code ?? ''
  eForm.tag_name = row.tag_name ?? ''
  eForm.remark = row.remark ?? ''
  editDlg.value = true
}

async function submitEdit(): Promise<void> {
  saving.value = true
  try {
    await updateStoreAccount(editId.value, {
      channel: eForm.channel.trim(),
      order_no: eForm.order_no.trim(),
      tag_code: eForm.tag_code.trim(),
      tag_name: eForm.tag_name.trim(),
      remark: eForm.remark.trim(),
    })
    toast.success('已保存')
    editDlg.value = false
    await qc.invalidateQueries({ queryKey: ['store-accounts'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

const viewDlg = ref(false)
const viewJson = ref('')

async function openView(row: StoreAccount): Promise<void> {
  try {
    const full = await getStoreAccount(row.id)
    viewJson.value = JSON.stringify(full, null, 2)
    viewDlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载失败')
  }
}

async function onDelete(row: StoreAccount): Promise<void> {
  const ok = await confirmDialog({ message: `删除记账「${row.account_no}」？` })
  if (!ok) return
  try {
    await deleteStoreAccount(row.id)
    toast.success('已删除')
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>
