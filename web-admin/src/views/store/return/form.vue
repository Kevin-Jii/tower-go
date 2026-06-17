<template>
  <div class="flex min-h-0 flex-1 flex-col gap-4">
    <BaseCard flush-body class="flex min-h-0 flex-1 flex-col">
      <template #header>
        <div class="flex w-full flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
          <div>
            <h2 class="m-0 text-base font-semibold text-slate-900">{{ isEdit ? '编辑返厂记录' : '新增返厂记录' }}</h2>
            <p class="m-0 mt-1 text-xs text-slate-500">选择返厂商品、数量、日期与物流费用</p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <BaseButton variant="secondary" @click="goBack">返回列表</BaseButton>
            <BaseButton variant="primary" :loading="saving" @click="submitForm">保存</BaseButton>
          </div>
        </div>
      </template>

      <div class="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto p-4">
        <div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
          <BaseFormItem label="返厂日期" required>
            <a-date-picker v-model="form.return_date" value-format="YYYY-MM-DD" class="w-full" :allow-clear="false" />
          </BaseFormItem>
          <BaseFormItem label="货拉拉费用">
            <BaseNumberInput v-model="form.logistics_fee" :min="0" :step="0.01" />
          </BaseFormItem>
          <div class="return-total-row">
            <span>押金合计</span>
            <strong>{{ formatMoney(formTotalDeposit) }}</strong>
          </div>
        </div>

        <div class="rounded border border-[var(--color-border-2)] bg-white p-4">
          <div class="mb-3 flex flex-col gap-3 lg:flex-row lg:items-end lg:justify-between">
            <BaseFormItem label="返厂商品" required class="m-0 min-w-0 flex-1">
              <a-select
                v-model="selectedProductIds"
                multiple
                allow-search
                allow-clear
                :loading="productsLoading"
                placeholder="请选择返厂商品，可多选"
                class="w-full"
                @change="syncSelectedProducts"
              >
                <a-option v-for="item in returnProducts" :key="item.id" :value="item.id" :label="productOptionLabel(item)">
                  <div class="flex items-center justify-between gap-4">
                    <span class="min-w-0 truncate">{{ item.product_name }}</span>
                    <span class="shrink-0 text-xs font-semibold text-slate-500">{{ formatMoney(item.deposit) }}</span>
                  </div>
                </a-option>
              </a-select>
            </BaseFormItem>
          </div>

          <div v-if="form.items.length" class="return-line-editor">
            <div class="return-line-editor__head">
              <span>商品名称</span>
              <span>默认押金</span>
              <span>数量</span>
              <span>押金小计</span>
              <span>操作</span>
            </div>
            <div v-for="line in form.items" :key="line.product_id" class="return-line-editor__row">
              <div class="min-w-0 truncate font-medium text-slate-800">{{ line.product_name }}</div>
              <div class="text-right font-semibold text-slate-700">{{ formatMoney(line.deposit) }}</div>
              <BaseNumberInput v-model="line.quantity" :min="0.01" :step="1" :hide-button="false" />
              <div class="return-line-editor__subtotal">{{ formatMoney(line.deposit * line.quantity) }}</div>
              <BaseButton variant="ghost" size="sm" @click="removeLine(line.product_id)">移除</BaseButton>
            </div>
          </div>
          <div v-else class="rounded border border-dashed border-slate-200 bg-slate-50 px-4 py-8 text-center text-sm text-slate-500">
            请选择需要返厂的商品
          </div>
        </div>

        <BaseFormItem label="整单备注">
          <BaseTextarea v-model="form.remark" :rows="3" />
        </BaseFormItem>
      </div>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useQuery } from '@tanstack/vue-query'
import { useRoute, useRouter } from 'vue-router'
import {
  BaseButton,
  BaseCard,
  BaseFormItem,
  BaseNumberInput,
  BaseTextarea,
} from '@/components/base'
import { createStoreReturn, getStoreReturn, listStoreReturnProducts, updateStoreReturn } from '@/api/storeReturn'
import type { StoreReturnProduct } from '@/api/types'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

type ReturnLine = {
  product_id: number
  product_name: string
  quantity: number
  deposit: number
  remark: string
}

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const tenantStoreId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)
const editId = computed(() => Number(route.query.id || 0))
const isEdit = computed(() => editId.value > 0)

const saving = ref(false)
const selectedProductIds = ref<number[]>([])
const form = reactive({
  client_request_id: '',
  return_date: today(),
  logistics_fee: 0,
  remark: '',
  items: [] as ReturnLine[],
})

const { data: enabledProductsPage, isLoading: productsLoading } = useQuery({
  queryKey: computed(() => ['store-return-products-enabled', tenantStoreId.value] as const),
  queryFn: () => listStoreReturnProducts({ page: 1, page_size: 500, store_id: tenantStoreId.value, status: 1 }),
})
const returnProducts = computed(() => enabledProductsPage.value?.list ?? [])
const returnProductMap = computed(() => {
  const map = new Map<number, StoreReturnProduct>()
  for (const product of returnProducts.value) map.set(product.id, product)
  return map
})
const formTotalDeposit = computed(() =>
  form.items.reduce((sum, x) => sum + Number(x.deposit || 0) * Number(x.quantity || 0), 0),
)

watch(returnProducts, () => {
  if (!form.items.length) return
  for (const line of form.items) {
    const product = returnProductMap.value.get(line.product_id)
    if (!product) continue
    line.product_name = product.product_name
    line.deposit = Number(product.deposit || 0)
  }
})

onMounted(() => {
  if (isEdit.value) void loadEditData()
  else form.client_request_id = createClientRequestId()
})

async function loadEditData(): Promise<void> {
  try {
    const row = await getStoreReturn(editId.value)
    form.client_request_id = row.client_request_id || ''
    form.return_date = normalizeDate(row.return_date)
    form.logistics_fee = Number(row.logistics_fee || 0)
    form.remark = row.remark || ''
    form.items = (row.items ?? [])
      .filter((x) => Number(x.product_id || 0) > 0)
      .map((x) => ({
        product_id: Number(x.product_id),
        product_name: x.product_name || '',
        quantity: Number(x.quantity || 1),
        deposit: Number(x.deposit || 0),
        remark: x.remark || '',
      }))
    selectedProductIds.value = form.items.map((x) => x.product_id)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载返厂记录失败')
    goBack()
  }
}

function syncSelectedProducts(): void {
  const selected = selectedProductIds.value.map(Number).filter((id) => id > 0)
  const existing = new Map(form.items.map((line) => [line.product_id, line]))
  form.items = selected
    .map((id) => {
      const old = existing.get(id)
      if (old) return old
      const product = returnProductMap.value.get(id)
      return {
        product_id: id,
        product_name: product?.product_name || '',
        quantity: 1,
        deposit: Number(product?.deposit || 0),
        remark: '',
      }
    })
    .filter((line) => line.product_name)
}

function removeLine(productID: number): void {
  selectedProductIds.value = selectedProductIds.value.filter((id) => Number(id) !== Number(productID))
  syncSelectedProducts()
}

async function submitForm(): Promise<void> {
  if (saving.value) return
  const items = form.items.map((x) => ({
    product_id: Number(x.product_id || 0),
    product_name: x.product_name.trim(),
    quantity: Number(x.quantity || 0),
    deposit: Number(x.deposit || 0),
    remark: x.remark.trim(),
  }))
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
    if (isEdit.value) {
      await updateStoreReturn(editId.value, body)
      toast.success('返厂记录已更新')
    } else {
      await createStoreReturn(body)
      toast.success('返厂记录已创建')
    }
    goBack()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

function goBack(): void {
  void router.push('/store/return')
}

function productOptionLabel(item: StoreReturnProduct): string {
  return `${item.product_name}（${formatMoney(item.deposit)}）`
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
  grid-template-columns: minmax(220px, 1fr) minmax(110px, 0.4fr) minmax(120px, 0.42fr) minmax(130px, 0.46fr) 72px;
  gap: 10px;
  align-items: center;
  min-width: 720px;
}

.return-line-editor__head {
  color: var(--color-text-3);
  font-size: 12px;
  font-weight: 600;
}

.return-line-editor__row {
  border-top: 1px solid var(--color-border-2);
  padding-top: 8px;
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
  min-height: 70px;
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
