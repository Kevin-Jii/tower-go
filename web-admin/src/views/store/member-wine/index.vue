<template>
  <div class="member-wine-page flex flex-col gap-4">
    <div class="flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
      <div>
        <h2 class="page-title">会员存酒</h2>
        <p class="m-0 mt-1 text-sm text-[var(--color-text-3)]">管理会员当前存酒数量和存取流水</p>
      </div>
      <div class="flex w-full flex-col gap-2 sm:flex-row lg:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-64" placeholder="会员 / 手机 / 酒品" clearable @enter="reload" />
        <BaseSelect v-model="onlyStock" class="w-full sm:w-36" :options="stockOptions" />
        <BaseButton variant="primary" @click="reload">查询</BaseButton>
        <BaseButton v-permission="'store:member:edit'" variant="primary" @click="openDeposit()">存入</BaseButton>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
      <div class="summary-tile">
        <span>存酒会员</span>
        <strong>{{ memberCount }}</strong>
      </div>
      <div class="summary-tile">
        <span>酒品记录</span>
        <strong>{{ total }}</strong>
      </div>
      <div class="summary-tile">
        <span>当前总数量</span>
        <strong>{{ totalQuantityText }}</strong>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading" min-width="980px">
      <template #cell-member="{ row }">
        <div class="leading-tight">
          <div class="font-medium text-[var(--color-text-1)]">{{ memberName((row as MemberWineStorage).member) }}</div>
          <div class="mt-1 text-xs text-[var(--color-text-3)]">{{ (row as MemberWineStorage).member?.phone || '-' }}</div>
        </div>
      </template>
      <template #cell-quantity="{ row }">
        <span class="font-semibold text-emerald-700">{{ formatQty((row as MemberWineStorage).quantity) }}</span>
        <span class="ml-1 text-xs text-[var(--color-text-3)]">{{ (row as MemberWineStorage).unit || '瓶' }}</span>
      </template>
      <template #cell-updated_at="{ row }">
        {{ formatDateTime((row as MemberWineStorage).updated_at) }}
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="rowActions(row as MemberWineStorage)" :max-inline="3" />
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

    <BaseDialog v-model="adjustDlg" :title="adjustMode === 'deposit' ? '会员存酒' : '会员取酒'" max-width="min(520px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="会员" required>
          <BaseSelect
            v-model="adjustForm.member_id"
            :disabled="adjustMode === 'withdraw'"
            :options="memberOptions"
            placeholder="选择会员"
            searchable
          />
        </BaseFormItem>
        <BaseFormItem v-if="adjustMode === 'deposit'" label="酒品" required>
          <a-cascader
            v-model="adjustForm.product_path"
            :options="productCascaderOptions"
            placeholder="先选分类，再选商品"
            allow-clear
            :path-mode="true"
            :check-strictly="false"
            @change="onProductChange"
          />
        </BaseFormItem>
        <BaseFormItem v-else label="酒品" required>
          <BaseInput v-model="adjustForm.wine_name" disabled />
        </BaseFormItem>
        <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
          <BaseFormItem label="数量" required>
            <BaseNumberInput v-model="adjustForm.quantity" :min="0.01" :step="1" :hide-button="false" />
          </BaseFormItem>
          <BaseFormItem label="规格" required>
            <BaseSelect
              v-model="adjustForm.unit"
              :disabled="adjustMode === 'withdraw' || specOptions.length <= 1"
              :options="specOptions"
              placeholder="选择规格"
            />
          </BaseFormItem>
        </div>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="adjustForm.remark" placeholder="可填写批次、来源或取酒说明" />
        </BaseFormItem>
        <div v-if="adjustMode === 'withdraw' && selectedStorage" class="rounded border border-[var(--color-border-2)] bg-slate-50 p-3 text-sm">
          当前剩余：
          <span class="font-semibold text-emerald-700">{{ formatQty(selectedStorage.quantity) }}{{ selectedStorage.unit || '瓶' }}</span>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="adjustDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitAdjust">提交</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="txnDlg" title="存取流水" max-width="min(960px, 96vw)">
      <div class="space-y-4">
        <div class="flex flex-col gap-2 sm:flex-row sm:items-end">
          <BaseFormItem label="类型" class="w-full sm:w-36">
            <BaseSelect v-model="txnType" :options="txnTypeOptions" />
          </BaseFormItem>
          <BaseFormItem label="开始日期" class="w-full sm:w-44">
            <BaseInput v-model="txnStart" type="date" />
          </BaseFormItem>
          <BaseFormItem label="结束日期" class="w-full sm:w-44">
            <BaseInput v-model="txnEnd" type="date" />
          </BaseFormItem>
          <BaseButton variant="primary" :loading="txnLoading" @click="reloadTransactions">查询</BaseButton>
        </div>

        <BaseTable :columns="txnColumns" :data="(txnRows as unknown) as Record<string, unknown>[]" :loading="txnLoading" min-width="860px">
          <template #cell-member="{ row }">
            {{ memberName((row as MemberWineTransaction).member) }}
          </template>
          <template #cell-type="{ row }">
            <span :class="(row as MemberWineTransaction).type === 1 ? 'text-emerald-700' : 'text-amber-700'">
              {{ (row as MemberWineTransaction).type === 1 ? '存入' : '取出' }}
            </span>
          </template>
          <template #cell-quantity="{ row }">
            {{ formatQty((row as MemberWineTransaction).quantity) }}{{ (row as MemberWineTransaction).unit || '瓶' }}
          </template>
          <template #cell-balance_after="{ row }">
            {{ formatQty((row as MemberWineTransaction).balance_after) }}{{ (row as MemberWineTransaction).unit || '瓶' }}
          </template>
          <template #cell-created_at="{ row }">
            {{ formatDateTime((row as MemberWineTransaction).created_at) }}
          </template>
        </BaseTable>
        <div class="flex justify-end">
          <BasePagination
            :page="txnPage"
            :page-size="txnPageSize"
            :total="txnTotal"
            @update:page="(p) => (txnPage = p)"
            @update:page-size="(s) => (txnPageSize = s)"
          />
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="txnDlg = false">关闭</BaseButton>
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
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn, TableRowAction } from '@/components/base/types'
import {
  depositMemberWine,
  listMemberWineStorages,
  listMemberWineTransactions,
  listMembers,
  withdrawMemberWine,
} from '@/api/member'
import { listPurchasableProducts } from '@/api/storeSupplier'
import { batchListProductUnitSpecs } from '@/api/supplierProduct'
import type { MemberRow, MemberWineStorage, MemberWineTransaction, ProductUnitSpec } from '@/api/types'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

const qc = useQueryClient()
const userStore = useUserStore()
const tenantStoreId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)
const keyword = ref('')
const onlyStock = ref<number | string>(1)
const page = ref(1)
const pageSize = ref(10)

const stockOptions: BaseSelectOption[] = [
  { label: '仅有库存', value: 1 },
  { label: '全部记录', value: 0 },
]

const queryKey = computed(() => ['member-wines', page.value, pageSize.value, keyword.value.trim(), onlyStock.value] as const)
const { data: pageData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: () =>
    listMemberWineStorages({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value.trim() || undefined,
      only_stock: Number(onlyStock.value),
    }),
})

const list = computed(() => pageData.value?.list ?? [])
const total = computed(() => Number(pageData.value?.total ?? 0))
const memberCount = computed(() => new Set(list.value.map((row) => row.member_id)).size)
const totalQuantityText = computed(() => formatQty(list.value.reduce((sum, row) => sum + Number(row.quantity || 0), 0)))

const memberQuery = useQuery({
  queryKey: ['members', 'member-wine-options'],
  queryFn: () => listMembers({ page: 1, page_size: 200 }),
})
const memberOptions = computed<BaseSelectOption[]>(() =>
  (memberQuery.data.value?.list ?? []).map((m) => ({
    label: `${m.phone}${m.name ? ` / ${m.name}` : ''}`,
    value: m.id,
  })),
)

const productQuery = useQuery({
  queryKey: computed(() => ['member-wine-products', tenantStoreId.value] as const),
  queryFn: () => listPurchasableProducts({ store_id: tenantStoreId.value }),
})
const productList = computed(() => productQuery.data.value ?? [])
const productIdsKey = computed(() =>
  productList.value
    .map((p) => p.id)
    .sort((a, b) => a - b)
    .join(','),
)
const unitSpecsQuery = useQuery({
  queryKey: computed(() => ['member-wine-product-specs', productIdsKey.value] as const),
  queryFn: async () => {
    const ids = productList.value.map((p) => p.id)
    if (!ids.length) return [] as ProductUnitSpec[]
    return batchListProductUnitSpecs(ids)
  },
  enabled: computed(() => productList.value.length > 0),
})
const specsByProduct = computed(() => {
  const map = new Map<number, ProductUnitSpec[]>()
  for (const spec of unitSpecsQuery.data.value ?? []) {
    if (!spec.is_enabled) continue
    if (!map.has(spec.product_id)) map.set(spec.product_id, [])
    map.get(spec.product_id)!.push(spec)
  }
  for (const [, arr] of map.entries()) {
    arr.sort((a, b) => Number(a.factor_to_base) - Number(b.factor_to_base))
  }
  return map
})
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

const columns: BaseTableColumn[] = [
  { key: 'member', label: '会员', minWidth: '180px' },
  { key: 'wine_name', label: '酒品名称', prop: 'wine_name', minWidth: '180px', ellipsis: true },
  { key: 'quantity', label: '当前数量', width: '120px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '180px', ellipsis: true },
  { key: 'updated_at', label: '更新时间', width: '170px' },
  { key: 'actions', label: '操作', width: '180px', align: 'right' },
]

const txnColumns: BaseTableColumn[] = [
  { key: 'created_at', label: '时间', width: '170px' },
  { key: 'member', label: '会员', width: '140px' },
  { key: 'type', label: '类型', width: '80px' },
  { key: 'wine_name', label: '酒品名称', prop: 'wine_name', minWidth: '160px', ellipsis: true },
  { key: 'quantity', label: '数量', width: '100px' },
  { key: 'balance_after', label: '剩余', width: '100px' },
  { key: 'operator_name', label: '操作人', prop: 'operator_name', width: '100px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '160px', ellipsis: true },
]

function reload(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['member-wines'] })
}

watch([page, pageSize, onlyStock], () => {
  void qc.invalidateQueries({ queryKey: ['member-wines'] })
})

function memberName(member?: MemberRow): string {
  if (!member) return '-'
  return member.name || member.phone || `会员${member.id}`
}

function formatQty(value: number | string | undefined): string {
  const n = Number(value ?? 0)
  if (!Number.isFinite(n)) return '0'
  return Number.isInteger(n) ? String(n) : n.toFixed(2).replace(/\.?0+$/, '')
}

function formatDateTime(value?: string): string {
  if (!value) return '-'
  return value.replace('T', ' ').replace(/\.\d+.*$/, '').replace(/\+\d{2}:\d{2}$/, '')
}

const adjustDlg = ref(false)
const saving = ref(false)
const adjustMode = ref<'deposit' | 'withdraw'>('deposit')
const selectedStorage = ref<MemberWineStorage | null>(null)
const adjustForm = reactive({
  member_id: undefined as number | undefined,
  product_path: [] as Array<string | number>,
  wine_name: '',
  unit: '瓶',
  quantity: 1,
  remark: '',
})

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

const selectedProduct = computed(() => {
  const id = getProductId(adjustForm.product_path)
  if (!id) return undefined
  return productList.value.find((p) => p.id === id)
})

function specLabel(spec: ProductUnitSpec): string {
  const name = String(spec.unit_name || spec.unit_code || '').trim()
  const factor = Number(spec.factor_to_base || 0)
  return factor > 0 ? `${name} / 换算${formatQty(factor)}` : name
}

function specValue(spec: ProductUnitSpec): string {
  return String(spec.unit_name || spec.unit_code || '').trim()
}

const specOptions = computed<BaseSelectOption[]>(() => {
  const id = getProductId(adjustForm.product_path)
  if (!id) {
    return adjustForm.unit ? [{ label: adjustForm.unit, value: adjustForm.unit }] : []
  }
  const specs = specsByProduct.value.get(id) ?? []
  if (specs.length) return specs.map((s) => ({ label: specLabel(s), value: specValue(s) }))
  const fallback = selectedProduct.value?.unit || '瓶'
  return [{ label: fallback, value: fallback }]
})

function onProductChange(): void {
  const product = selectedProduct.value
  adjustForm.wine_name = product?.name ?? ''
  const first = specOptions.value[0]
  adjustForm.unit = first ? String(first.value) : '瓶'
}

function openDeposit(row?: MemberWineStorage): void {
  adjustMode.value = 'deposit'
  selectedStorage.value = row ?? null
  adjustForm.member_id = row?.member_id
  adjustForm.product_path = []
  adjustForm.wine_name = row?.wine_name ?? ''
  adjustForm.unit = row?.unit || '瓶'
  adjustForm.quantity = 1
  adjustForm.remark = ''
  adjustDlg.value = true
}

function openWithdraw(row: MemberWineStorage): void {
  adjustMode.value = 'withdraw'
  selectedStorage.value = row
  adjustForm.member_id = row.member_id
  adjustForm.product_path = []
  adjustForm.wine_name = row.wine_name
  adjustForm.unit = row.unit || '瓶'
  adjustForm.quantity = 1
  adjustForm.remark = ''
  adjustDlg.value = true
}

async function submitAdjust(): Promise<void> {
  if (!adjustForm.member_id) {
    toast.warning('请选择会员')
    return
  }
  const product = selectedProduct.value
  const wineName = adjustMode.value === 'withdraw' ? adjustForm.wine_name.trim() : String(product?.name || '').trim()
  if (!wineName) {
    toast.warning('请选择酒品')
    return
  }
  const unit = String(adjustForm.unit || '').trim()
  if (!unit) {
    toast.warning('请选择规格')
    return
  }
  if (!adjustForm.quantity || adjustForm.quantity <= 0) {
    toast.warning('数量必须大于0')
    return
  }
  if (adjustMode.value === 'withdraw' && selectedStorage.value && adjustForm.quantity > Number(selectedStorage.value.quantity || 0)) {
    toast.warning('取出数量不能大于当前剩余')
    return
  }

  saving.value = true
  const body = {
    member_id: Number(adjustForm.member_id),
    wine_name: wineName,
    unit,
    quantity: Number(adjustForm.quantity),
    remark: adjustForm.remark.trim() || undefined,
  }
  try {
    if (adjustMode.value === 'deposit') {
      await depositMemberWine(body)
      toast.success('已存入')
    } else {
      await withdrawMemberWine(body)
      toast.success('已取出')
    }
    adjustDlg.value = false
    await qc.invalidateQueries({ queryKey: ['member-wines'] })
    if (txnDlg.value) await loadTransactions()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '提交失败')
  } finally {
    saving.value = false
  }
}

function rowActions(row: MemberWineStorage): TableRowAction[] {
  return [
    { label: '存入', permission: 'store:member:edit', place: 'inline', onClick: () => openDeposit(row) },
    { label: '取出', permission: 'store:member:edit', place: 'inline', disabled: Number(row.quantity || 0) <= 0, onClick: () => openWithdraw(row) },
    { label: '流水', permission: 'store:member:list', place: 'inline', onClick: () => openTransactions(row) },
  ]
}

const txnDlg = ref(false)
const txnLoading = ref(false)
const txnStorage = ref<MemberWineStorage | null>(null)
const txnType = ref<number | string>(0)
const txnStart = ref('')
const txnEnd = ref('')
const txnPage = ref(1)
const txnPageSize = ref(10)
const txnTotal = ref(0)
const txnRows = ref<MemberWineTransaction[]>([])

const txnTypeOptions: BaseSelectOption[] = [
  { label: '全部', value: 0 },
  { label: '存入', value: 1 },
  { label: '取出', value: 2 },
]

function openTransactions(row?: MemberWineStorage): void {
  txnStorage.value = row ?? null
  txnType.value = 0
  txnStart.value = ''
  txnEnd.value = ''
  txnPage.value = 1
  txnPageSize.value = 10
  txnRows.value = []
  txnTotal.value = 0
  txnDlg.value = true
  void loadTransactions()
}

async function loadTransactions(): Promise<void> {
  txnLoading.value = true
  try {
    const data = await listMemberWineTransactions({
      page: txnPage.value,
      page_size: txnPageSize.value,
      storage_id: txnStorage.value?.id,
      type: Number(txnType.value) || undefined,
      start_date: txnStart.value || undefined,
      end_date: txnEnd.value || undefined,
    })
    txnRows.value = data.list ?? []
    txnTotal.value = Number(data.total ?? 0)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载流水失败')
  } finally {
    txnLoading.value = false
  }
}

function reloadTransactions(): void {
  txnPage.value = 1
  void loadTransactions()
}

watch([txnPage, txnPageSize], () => {
  if (txnDlg.value) void loadTransactions()
})
</script>

<style scoped>
.summary-tile {
  display: flex;
  min-height: 72px;
  flex-direction: column;
  justify-content: center;
  gap: 6px;
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  background: linear-gradient(180deg, #ffffff 0%, #f8fafc 100%);
  padding: 14px 16px;
}

.summary-tile span {
  font-size: 12px;
  color: var(--color-text-3);
}

.summary-tile strong {
  font-size: 24px;
  line-height: 1;
  color: var(--color-text-1);
}
</style>
