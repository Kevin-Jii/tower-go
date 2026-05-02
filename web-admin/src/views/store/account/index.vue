<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">门店记账</h2>
      <div class="flex flex-col sm:flex-row flex-wrap gap-2 w-full md:w-auto">
        <BaseInput v-model="rangeStart" class="w-full sm:w-40" type="date" />
        <BaseInput v-model="rangeEnd" class="w-full sm:w-40" type="date" />
        <BaseButton variant="primary" @click="reloadAll">查询</BaseButton>
        <BaseButton v-permission="'store:account:add'" variant="primary" @click="openCreate">快速记账</BaseButton>
        <BaseButton v-permission="'store:account:add'" variant="secondary" @click="openCustomCreate">自定义记账</BaseButton>
      </div>
    </div>

    <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">总销售额</span></template>
        <p class="text-2xl font-semibold text-indigo-700 m-0">{{ statsAmount }}</p>
      </BaseCard>
      <BaseCard>
        <template #header><span class="font-semibold text-slate-800">净利润</span></template>
        <p class="text-2xl font-semibold text-emerald-700 m-0">{{ statsNetIncome }}</p>
      </BaseCard>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading" min-width="960px">
      <template #cell-channel="{ row }">
        {{ channelLabel((row as StoreAccount).channel) }}
      </template>
      <template #cell-member="{ row }">
        {{ memberLabel(row as StoreAccount) }}
      </template>
      <template #cell-payment_status="{ row }">
        {{ paymentStatusLabel((row as StoreAccount).payment_status) }}
      </template>
      <template #cell-net_income_amount="{ row }">
        {{ formatMoney((row as StoreAccount).net_income_amount) }}
      </template>
      <template #cell-account_date="{ row }">
        {{ formatDate((row as StoreAccount).account_date) }}
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="accountRowActions(row as StoreAccount)" />
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

    <BaseDialog v-model="customCreateDlg" title="自定义记账" max-width="min(720px, 96vw)">
      <p class="m-0 mb-3 text-xs text-[var(--color-text-3)]">
        商品明细为手写描述，不关联系统商品与库存；请填写单价与小计依据。渠道、会员、支付状态、备注与快速记账一致。
      </p>
      <div class="space-y-4">
        <BaseFormItem label="渠道" required>
          <BaseSelect v-model="customForm.channel" :options="channelOptions" placeholder="请选择销售渠道" />
        </BaseFormItem>
        <BaseFormItem label="绑定会员">
          <BaseSelect v-model="customForm.member_id" :options="memberOptionsWithNone" placeholder="可选，默认不绑定" />
        </BaseFormItem>
        <BaseFormItem label="支付状态">
          <BaseSelect v-model="customForm.payment_status" :options="paymentStatusOptions" />
        </BaseFormItem>
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">商品明细（任意描述）</span>
          <BaseButton variant="secondary" size="sm" @click="addCustomLine">加一行</BaseButton>
        </div>
        <div
          v-for="(line, idx) in customForm.lines"
          :key="idx"
          class="rounded border border-[var(--color-border-2)] p-3 flex flex-col gap-3"
        >
          <BaseFormItem label="明细描述" required class="w-full">
            <BaseTextarea v-model="line.description" :rows="2" placeholder="可填写任意商品或服务说明" />
          </BaseFormItem>
          <div class="flex flex-wrap items-end gap-2">
            <BaseFormItem label="数量" required class="w-28">
              <BaseNumberInput v-model="line.quantity" :min="0.01" :step="0.01" />
            </BaseFormItem>
            <BaseFormItem label="单位" required class="w-28">
              <BaseInput v-model="line.unit" placeholder="如 瓶、次、项" />
            </BaseFormItem>
            <BaseFormItem label="单价" required class="w-32">
              <BaseNumberInput v-model="line.price" :min="0.01" :step="0.01" />
            </BaseFormItem>
            <BaseFormItem label="行备注" class="min-w-[140px] flex-1">
              <BaseInput v-model="line.line_remark" placeholder="可选" />
            </BaseFormItem>
            <BaseButton variant="ghost" size="sm" class="shrink-0" @click="removeCustomLine(idx)">移除</BaseButton>
          </div>
        </div>
        <BaseFormItem label="整单备注">
          <BaseTextarea v-model="customForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="customCreateDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submitCustomCreate">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="createDlg" title="快速记账（多商品）" max-width="min(720px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="渠道" required>
          <BaseSelect v-model="cForm.channel" :options="channelOptions" placeholder="请选择销售渠道" />
        </BaseFormItem>
        <BaseFormItem label="绑定会员">
          <BaseSelect v-model="cForm.member_id" :options="memberOptionsWithNone" placeholder="可选，默认不绑定" />
        </BaseFormItem>
        <BaseFormItem label="支付状态">
          <BaseSelect v-model="cForm.payment_status" :options="paymentStatusOptions" />
        </BaseFormItem>
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">商品明细</span>
          <BaseButton variant="secondary" size="sm" @click="addCreateLine">加一行</BaseButton>
        </div>
        <div
          v-for="(line, idx) in cForm.lines"
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
              @change="onCreateProductChange(idx)"
            />
          </BaseFormItem>
          <BaseFormItem label="数量" required class="w-28">
            <BaseNumberInput v-model="line.quantity" :min="0.01" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="单位" class="w-28">
            <BaseSelect
              v-model="line.unit"
              :options="lineUnitOptions(line)"
              :disabled="lineUnitOptions(line).length <= 1"
              placeholder="单位"
            />
          </BaseFormItem>
          <BaseButton variant="ghost" size="sm" @click="removeCreateLine(idx)">移除</BaseButton>
        </div>
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
          <BaseSelect v-model="eForm.channel" :options="channelOptions" placeholder="请选择销售渠道" />
        </BaseFormItem>
        <BaseFormItem label="绑定会员">
          <BaseSelect v-model="eForm.member_id" :options="memberOptionsWithNone" placeholder="可选，默认不绑定" />
        </BaseFormItem>
        <BaseFormItem label="支付状态">
          <BaseSelect v-model="eForm.payment_status" :options="paymentStatusOptions" />
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

    <BaseDialog v-model="viewDlg" title="记账详情" max-width="min(720px, 96vw)">
      <div v-if="viewAccount" class="max-h-[70vh] overflow-y-auto space-y-4 pr-1">
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-2 text-sm">
          <div><span class="text-[var(--color-text-3)]">单号</span>：{{ viewAccount.account_no }}</div>
          <div><span class="text-[var(--color-text-3)]">记账日期</span>：{{ formatDate(viewAccount.account_date) }}</div>
          <div><span class="text-[var(--color-text-3)]">渠道</span>：{{ channelLabel(viewAccount.channel) }}</div>
          <div><span class="text-[var(--color-text-3)]">会员</span>：{{ memberLabel(viewAccount) }}</div>
          <div><span class="text-[var(--color-text-3)]">支付状态</span>：{{ paymentStatusLabel(viewAccount.payment_status) }}</div>
          <div><span class="text-[var(--color-text-3)]">订单号</span>：{{ viewAccount.order_no || '-' }}</div>
          <div><span class="text-[var(--color-text-3)]">总金额</span>：{{ formatMoney(viewAccount.total_amount) }}</div>
          <div><span class="text-[var(--color-text-3)]">其他支出</span>：{{ formatMoney(viewAccount.other_expense_amount) }}</div>
          <div><span class="text-[var(--color-text-3)]">商品成本</span>：{{ formatMoney(accountItemCost(viewAccount)) }}</div>
          <div><span class="text-[var(--color-text-3)]">耗材金额</span>：{{ formatMoney(accountConsumableAmount(viewAccount)) }}</div>
          <div><span class="text-[var(--color-text-3)]">净收入</span>：{{ formatMoney(viewAccount.net_income_amount) }}</div>
          <div><span class="text-[var(--color-text-3)]">明细条数</span>：{{ viewAccount.item_count ?? (viewAccount.items?.length ?? 0) }}</div>
          <div class="sm:col-span-2">
            <span class="text-[var(--color-text-3)]">标签</span>：{{ viewAccount.tag_name || viewAccount.tag_code || '-' }}
          </div>
          <div class="sm:col-span-2"><span class="text-[var(--color-text-3)]">备注</span>：{{ viewAccount.remark || '-' }}</div>
          <div class="sm:col-span-2 text-xs text-[var(--color-text-3)]">
            创建时间：{{ viewAccount.created_at ? formatDateTime(viewAccount.created_at) : '-' }}
          </div>
        </div>

        <div>
          <p class="m-0 mb-2 text-sm font-medium text-slate-800">商品明细</p>
          <div v-if="(viewAccount.items?.length ?? 0) === 0" class="text-sm text-[var(--color-text-3)]">暂无商品明细</div>
          <div v-else class="overflow-x-auto rounded border border-[var(--color-border-2)]">
            <table class="w-full min-w-[560px] border-collapse text-sm">
              <thead>
                <tr class="bg-[var(--color-fill-2)]">
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-left font-medium">商品</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-left w-24">规格</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-20">数量</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-center w-16">单位</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-24">单价</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-24">小计</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-left min-w-[80px]">备注</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="it in viewAccount.items" :key="it.id">
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5">
                    {{ it.product_name || (it.product_id ? `商品#${it.product_id}` : '—') }}
                  </td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5">{{ it.spec || '-' }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ it.quantity }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-center">{{ it.unit || '-' }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ formatMoney(it.price) }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ formatMoney(it.amount) }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5">{{ it.remark || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div v-if="(viewAccount.consumables?.length ?? 0) > 0">
          <p class="m-0 mb-2 text-sm font-medium text-slate-800">消耗品</p>
          <div class="overflow-x-auto rounded border border-[var(--color-border-2)]">
            <table class="w-full min-w-[480px] border-collapse text-sm">
              <thead>
                <tr class="bg-[var(--color-fill-2)]">
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-left font-medium">商品</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-20">数量</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-center w-16">单位</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-24">单价</th>
                  <th class="border-b border-[var(--color-border-2)] px-2 py-2 text-right w-24">小计</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="c in viewAccount.consumables" :key="c.id">
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5">{{ c.product_name || `商品#${c.product_id}` }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ c.quantity }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-center">{{ c.unit || '-' }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ formatMoney(c.price) }}</td>
                  <td class="border-b border-[var(--color-border-2)] px-2 py-1.5 text-right">{{ formatMoney(c.amount) }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <div class="rounded border border-[var(--color-border-2)] bg-[var(--color-fill-1)] px-3 py-2 text-sm">
          <span class="text-[var(--color-text-3)]">净利润口径：</span>
          销售额 {{ formatMoney(viewAccount.total_amount) }} - 其他支出 {{ formatMoney(viewAccount.other_expense_amount) }} - 商品成本
          {{ formatMoney(accountItemCost(viewAccount)) }} - 耗材金额 {{ formatMoney(accountConsumableAmount(viewAccount)) }} =
          <span class="font-semibold text-emerald-700">{{ formatMoney(accountNetProfitBreakdown(viewAccount)) }}</span>
        </div>
      </div>
      <p v-else class="m-0 text-sm text-[var(--color-text-3)]">加载中…</p>
      <template #footer>
        <BaseButton variant="ghost" @click="viewDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="consumableDlg" title="绑定消耗品" max-width="min(720px, 96vw)">
      <div class="space-y-4">
        <p class="m-0 text-sm text-slate-600">记账单：{{ consumableTarget?.account_no }}（绑定后会计入成本并扣减净利润）</p>
        <div class="flex items-center justify-between gap-2">
          <span class="text-sm font-medium text-slate-700">消耗品明细</span>
          <BaseButton variant="secondary" size="sm" @click="addConsumableLine">加一行</BaseButton>
        </div>
        <div
          v-for="(line, idx) in consumableLines"
          :key="idx"
          class="rounded border border-[var(--color-border-2)] p-3 flex flex-wrap items-end gap-2"
        >
          <BaseFormItem label="消耗品" required class="min-w-[220px] flex-1">
            <a-cascader
              v-model="line.product_path"
              :options="productCascaderOptions"
              placeholder="先选分类，再选商品"
              allow-clear
              :path-mode="true"
              :check-strictly="false"
              @change="onConsumableProductChange(idx)"
            />
          </BaseFormItem>
          <BaseFormItem label="数量" required class="w-28">
            <BaseNumberInput v-model="line.quantity" :min="0.01" :step="0.01" />
          </BaseFormItem>
          <BaseFormItem label="单位" class="w-28">
            <BaseSelect
              v-model="line.unit"
              :options="lineUnitOptions(line)"
              :disabled="lineUnitOptions(line).length <= 1"
              placeholder="单位"
            />
          </BaseFormItem>
          <BaseButton variant="ghost" size="sm" @click="removeConsumableLine(idx)">移除</BaseButton>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="consumableDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="consumableSaving" @click="submitConsumables">保存消耗品</BaseButton>
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
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn, TableRowAction } from '@/components/base/types'
import {
  bindStoreAccountConsumables,
  createStoreAccount,
  getStoreAccount,
  getStoreAccountStats,
  listStoreAccounts,
  updateStoreAccount,
} from '@/api/storeAccount'
import { listDictDataByTypeCode } from '@/api/dict'
import { listProductUnitSpecs } from '@/api/supplierProduct'
import { listPurchasableProducts } from '@/api/storeSupplier'
import { listMembers } from '@/api/member'
import type { DictData, MemberRow, ProductUnitSpec, StoreAccount } from '@/api/types'
import { toast } from '@/feedback/toast'
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

const stats = ref<{ total_amount?: number; net_income_amount?: number; count?: number }>({})
const statsAmount = computed(() => (stats.value.total_amount ?? 0).toFixed(2))
const statsNetIncome = computed(() => (stats.value.net_income_amount ?? 0).toFixed(2))

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
const { data: channelData } = useQuery({
  queryKey: ['dict-data', 'sales_channel'],
  queryFn: () => listDictDataByTypeCode('sales_channel'),
})
const channelOptions = computed(() => (channelData.value ?? []).map((d) => ({ label: d.label, value: d.value })))
const channelDictMap = computed(() => {
  const map = new Map<string, string>()
  for (const d of channelData.value ?? []) {
    map.set(String(d.value), d.label || String(d.value))
  }
  return map
})
const { data: membersPageData } = useQuery({
  queryKey: ['store-account-members'],
  queryFn: () => listMembers({ page: 1, page_size: 200 }),
})
const memberList = computed(() => membersPageData.value?.list ?? ([] as MemberRow[]))
const memberOptions = computed(() =>
  memberList.value.map((m) => ({
    label: `${m.phone}${m.name ? `（${m.name}）` : ''}`,
    value: m.id,
  })),
)
const memberOptionsWithNone = computed(() => [{ label: '不绑定会员', value: 0 }, ...memberOptions.value])
const memberMap = computed(() => {
  const map = new Map<number, MemberRow>()
  for (const m of memberList.value) {
    map.set(m.id, m)
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
  { key: 'member', label: '会员', width: '150px', ellipsis: true },
  { key: 'payment_status', label: '支付状态', width: '96px' },
  { key: 'order_no', label: '订单号', prop: 'order_no', minWidth: '120px', ellipsis: true },
  { key: 'total_amount', label: '销售额', prop: 'total_amount', width: '96px' },
  { key: 'net_income_amount', label: '净利润', prop: 'net_income_amount', width: '96px' },
  { key: 'account_date', label: '日期', width: '120px' },
  { key: 'actions', label: '操作', width: '140px', align: 'right' },
]

function formatDate(v: string): string {
  if (!v) return '-'
  return String(v).slice(0, 10)
}

function channelLabel(v: string | undefined): string {
  const key = String(v || '').trim()
  if (!key) return '-'
  return channelDictMap.value.get(key) || key
}

function memberLabel(row: StoreAccount): string {
  if (row.member) {
    const phone = String(row.member.phone || '').trim()
    const name = String(row.member.name || '').trim()
    if (phone && name) return `${phone}（${name}）`
    return phone || name || `会员#${row.member.id}`
  }
  const mid = Number(row.member_id || 0)
  if (mid > 0) {
    const m = memberMap.value.get(mid)
    if (m) {
      return `${m.phone}${m.name ? `（${m.name}）` : ''}`
    }
    return `会员#${mid}`
  }
  return '-'
}

const paymentStatusOptions = [
  { label: '已支付', value: 1 },
  { label: '未支付', value: 2 },
]

function paymentStatusLabel(v: number | undefined): string {
  return Number(v) === 2 ? '未支付' : '已支付'
}

function canEditAccount(row: StoreAccount): boolean {
  const s = String(row.created_at || '').trim()
  if (!s) return false
  const created = new Date(s)
  if (Number.isNaN(created.getTime())) return false
  const now = new Date()
  const dayStart = new Date(created.getFullYear(), created.getMonth(), created.getDate(), 0, 0, 0, 0)
  let cutoff = new Date(dayStart.getTime() + 24 * 60 * 60 * 1000)
  if (created.getHours() >= 23) {
    cutoff = new Date(dayStart.getTime() + 27 * 60 * 60 * 1000)
  }
  return now < cutoff
}

function accountRowActions(row: StoreAccount): TableRowAction[] {
  const editable = canEditAccount(row)
  return [
    { label: '详情', permission: 'store:account:list', onClick: () => openView(row) },
    {
      label: '绑定消耗品',
      permission: 'store:account:edit',
      disabled: !editable,
      onClick: () => openConsumableDlg(row),
    },
    { label: '编辑', permission: 'store:account:edit', disabled: !editable, onClick: () => openEdit(row) },
  ]
}

const createDlg = ref(false)
const customCreateDlg = ref(false)
const saving = ref(false)
interface AccountLine {
  product_path: Array<string | number> | string | number | undefined
  quantity: number
  unit: string
}
const cForm = reactive({
  channel: '',
  member_id: 0,
  payment_status: 1,
  lines: [] as AccountLine[],
  remark: '',
})

interface CustomAccountLine {
  description: string
  quantity: number
  unit: string
  price: number
  line_remark: string
}
const customForm = reactive({
  channel: '',
  member_id: 0,
  payment_status: 1,
  lines: [] as CustomAccountLine[],
  remark: '',
})

function openCreate(): void {
  cForm.channel = ''
  cForm.member_id = 0
  cForm.payment_status = 1
  cForm.lines = [makeAccountLine()]
  cForm.remark = ''
  createDlg.value = true
}

function makeCustomLine(): CustomAccountLine {
  return {
    description: '',
    quantity: 1,
    unit: '',
    price: 0,
    line_remark: '',
  }
}

function openCustomCreate(): void {
  customForm.channel = ''
  customForm.member_id = 0
  customForm.payment_status = 1
  customForm.lines = [makeCustomLine()]
  customForm.remark = ''
  customCreateDlg.value = true
}

function addCustomLine(): void {
  customForm.lines.push(makeCustomLine())
}

function removeCustomLine(idx: number): void {
  customForm.lines = customForm.lines.filter((_, i) => i !== idx)
  if (!customForm.lines.length) customForm.lines.push(makeCustomLine())
}

function makeAccountLine(): AccountLine {
  return {
    product_path: [],
    quantity: 1,
    unit: '',
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

function lineUnitOptions(line: AccountLine): Array<{ label: string; value: string | number }> {
  const pid = getProductId(line.product_path)
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
}

function onCreateProductChange(idx: number): void {
  const line = cForm.lines[idx]
  if (!line) return
  const options = lineUnitOptions(line)
  line.unit = String(options[0]?.value || '')
}

function addCreateLine(): void {
  cForm.lines.push(makeAccountLine())
}

function removeCreateLine(idx: number): void {
  cForm.lines = cForm.lines.filter((_, i) => i !== idx)
  if (!cForm.lines.length) cForm.lines.push(makeAccountLine())
}

async function submitCustomCreate(): Promise<void> {
  if (!customForm.channel.trim()) {
    toast.warning('请选择渠道')
    return
  }
  const items: Array<{
    product_id: number
    product_name: string
    quantity: number
    unit: string
    price: number
    amount: number
    remark: string
    spec: string
  }> = []
  for (const line of customForm.lines) {
    const name = line.description.trim()
    const unit = line.unit.trim()
    if (!name) {
      toast.warning('请填写每条明细的描述')
      return
    }
    if (!unit) {
      toast.warning('请填写每条明细的单位')
      return
    }
    const price = Number(line.price)
    if (!Number.isFinite(price) || price <= 0) {
      toast.warning('自定义明细须填写大于 0 的单价')
      return
    }
    const qty = Number(line.quantity)
    if (!Number.isFinite(qty) || qty <= 0) {
      toast.warning('请填写有效的数量')
      return
    }
    const amount = Number((price * qty).toFixed(2))
    items.push({
      product_id: 0,
      product_name: name,
      quantity: qty,
      unit,
      price,
      amount,
      remark: line.line_remark.trim(),
      spec: '',
    })
  }
  if (!items.length) {
    toast.warning('请至少添加一条明细')
    return
  }
  saving.value = true
  try {
    await createStoreAccount({
      store_id: tenantStoreId.value,
      member_id: customForm.member_id > 0 ? customForm.member_id : undefined,
      payment_status: customForm.payment_status,
      channel: customForm.channel.trim(),
      remark: customForm.remark.trim(),
      other_expense_amount: 0,
      items,
    })
    toast.success('已保存')
    customCreateDlg.value = false
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

async function submitCreate(): Promise<void> {
  if (!cForm.channel.trim()) {
    toast.warning('请选择渠道')
    return
  }
  const items = cForm.lines
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
      spec: '',
      price: 0,
      amount: 0,
      remark: '',
    }))
  if (!items.length) {
    toast.warning('请至少选择一条有效商品明细')
    return
  }
  saving.value = true
  try {
    await createStoreAccount({
      store_id: tenantStoreId.value,
      member_id: cForm.member_id > 0 ? cForm.member_id : undefined,
      payment_status: cForm.payment_status,
      channel: cForm.channel.trim(),
      remark: cForm.remark.trim(),
      other_expense_amount: 0,
      items,
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
  member_id: 0,
  payment_status: 1,
  channel: '',
  tag_code: '',
  tag_name: '',
  remark: '',
})

function openEdit(row: StoreAccount): void {
  if (!canEditAccount(row)) {
    toast.warning('该记录已超过可编辑时间')
    return
  }
  editId.value = row.id
  eForm.member_id = Number(row.member_id || 0)
  eForm.payment_status = Number(row.payment_status || 1)
  eForm.channel = row.channel ?? ''
  eForm.tag_code = row.tag_code ?? ''
  eForm.tag_name = row.tag_name ?? ''
  eForm.remark = row.remark ?? ''
  editDlg.value = true
}

async function submitEdit(): Promise<void> {
  saving.value = true
  try {
    await updateStoreAccount(editId.value, {
      member_id: eForm.member_id,
      payment_status: eForm.payment_status,
      channel: eForm.channel.trim(),
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
const viewAccount = ref<StoreAccount | null>(null)
const consumableDlg = ref(false)
const consumableSaving = ref(false)
const consumableTarget = ref<StoreAccount | null>(null)
const consumableLines = ref<AccountLine[]>([])

function formatMoney(v: number | string | undefined | null): string {
  const n = Number(v ?? 0)
  return Number.isFinite(n) ? n.toFixed(2) : '0.00'
}

function formatDateTime(v: string): string {
  const s = String(v || '').trim()
  if (!s) return '-'
  const d = new Date(s)
  if (Number.isNaN(d.getTime())) return s.slice(0, 19).replace('T', ' ')
  const pad = (x: number) => String(x).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function accountConsumableAmount(account: StoreAccount): number {
  return (account.consumables ?? []).reduce((sum, c) => sum + Number(c.amount || 0), 0)
}

function itemCostPrice(productId: number, unit?: string): number {
  const specs = specsByProduct.value.get(productId) ?? []
  const normalized = String(unit || '').trim().toLowerCase()
  for (const s of specs) {
    if (!s.is_enabled) continue
    const code = String(s.unit_code || '').trim().toLowerCase()
    const name = String(s.unit_name || '').trim().toLowerCase()
    if (normalized && (normalized === code || normalized === name)) return Number(s.cost_price || 0)
  }
  for (const s of specs) {
    if (!s.is_enabled) continue
    const code = String(s.unit_code || '').trim().toLowerCase()
    const name = String(s.unit_name || '').trim().toLowerCase()
    if (normalized && (code.includes(normalized) || name.includes(normalized) || normalized.includes(code) || normalized.includes(name))) {
      return Number(s.cost_price || 0)
    }
  }
  return 0
}

function accountItemCost(account: StoreAccount): number {
  return (account.items ?? []).reduce((sum, it) => {
    const qty = Number(it.quantity || 0)
    if (qty <= 0) return sum
    return sum + qty * itemCostPrice(it.product_id, it.unit)
  }, 0)
}

function accountNetProfitBreakdown(account: StoreAccount): number {
  return (
    Number(account.total_amount || 0) -
    Number(account.other_expense_amount || 0) -
    accountItemCost(account) -
    accountConsumableAmount(account)
  )
}

async function openView(row: StoreAccount): Promise<void> {
  viewAccount.value = null
  viewDlg.value = true
  try {
    const full = await getStoreAccount(row.id)
    viewAccount.value = full
  } catch (e: unknown) {
    viewDlg.value = false
    toast.error(e instanceof Error ? e.message : '加载失败')
  }
}

function makeConsumableLine(): AccountLine {
  return { product_path: [], quantity: 1, unit: '' }
}

function addConsumableLine(): void {
  consumableLines.value.push(makeConsumableLine())
}

function removeConsumableLine(idx: number): void {
  consumableLines.value = consumableLines.value.filter((_, i) => i !== idx)
  if (!consumableLines.value.length) consumableLines.value.push(makeConsumableLine())
}

function onConsumableProductChange(idx: number): void {
  const line = consumableLines.value[idx]
  if (!line) return
  const options = lineUnitOptions(line)
  line.unit = String(options[0]?.value || '')
}

async function openConsumableDlg(row: StoreAccount): Promise<void> {
  if (!canEditAccount(row)) {
    toast.warning('该记录已超过可编辑时间')
    return
  }
  consumableTarget.value = row
  consumableLines.value = [makeConsumableLine()]
  try {
    const full = await getStoreAccount(row.id)
    if (full.consumables?.length) {
      consumableLines.value = full.consumables.map((c) => ({
        product_path: [c.product_id],
        quantity: Number(c.quantity || 1),
        unit: c.unit || '',
      }))
    }
    consumableDlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载消耗品失败')
  }
}

async function submitConsumables(): Promise<void> {
  if (!consumableTarget.value) return
  const consumables = consumableLines.value
    .map((line) => ({
      product_id: getProductId(line.product_path),
      quantity: line.quantity,
      unit: line.unit.trim(),
      price: 0,
      amount: 0,
      remark: '',
    }))
    .filter((x) => x.product_id && x.quantity > 0 && x.unit)
    .map((x) => ({
      product_id: x.product_id as number,
      quantity: x.quantity,
      unit: x.unit,
      price: x.price,
      amount: x.amount,
      remark: x.remark,
    }))
  if (!consumables.length) {
    toast.warning('请至少选择一条有效消耗品明细')
    return
  }
  consumableSaving.value = true
  try {
    await bindStoreAccountConsumables(consumableTarget.value.id, { consumables })
    toast.success('消耗品已绑定')
    consumableDlg.value = false
    await qc.invalidateQueries({ queryKey: ['store-accounts'] })
    await loadStats()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    consumableSaving.value = false
  }
}
</script>
