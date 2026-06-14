<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col lg:flex-row lg:items-end justify-between gap-3">
      <div>
        <h2 class="page-title">AI运营</h2>
        <p class="m-0 text-sm text-[var(--color-text-3)]">美团外卖半自动运营：接口同步订单，生成建议，人工确认后执行。</p>
      </div>
      <div class="flex flex-col sm:flex-row gap-2 w-full lg:w-auto">
        <BaseInput v-model="rangeStart" class="w-full sm:w-40" type="date" />
        <BaseInput v-model="rangeEnd" class="w-full sm:w-40" type="date" />
        <BaseSelect v-model="activeAccountId" class="w-full sm:w-56" :options="accountOptions" placeholder="选择美团店铺" />
        <BaseButton variant="primary" @click="reloadAll">刷新</BaseButton>
        <BaseButton variant="secondary" @click="openAccountDlg()">店铺配置</BaseButton>
      </div>
    </div>

    <div class="grid grid-cols-2 lg:grid-cols-4 gap-3">
      <BaseCard v-for="card in statCards" :key="card.label">
        <template #header><span class="font-semibold text-slate-700">{{ card.label }}</span></template>
        <p class="m-0 text-2xl font-semibold" :class="card.tone">{{ card.value }}</p>
      </BaseCard>
    </div>

    <div class="grid grid-cols-1 xl:grid-cols-[380px_1fr] gap-4">
      <div class="flex flex-col gap-4">
        <BaseCard>
          <template #header>
            <div class="flex items-center justify-between gap-2">
              <span class="font-semibold text-slate-800">数据导入</span>
              <BaseButton variant="link" size="sm" @click="fillDemoData">填入示例</BaseButton>
            </div>
          </template>
          <div class="space-y-3">
            <BaseFormItem label="同步美团订单">
              <div class="space-y-2">
                <BaseTextarea v-model="syncOrderIds" :rows="3" placeholder="输入美团订单号，多个订单号用逗号或换行分隔" />
                <BaseButton variant="primary" class="w-full" :loading="syncingOpenAPIOrders" @click="syncOpenAPIOrders">从美团接口同步</BaseButton>
              </div>
            </BaseFormItem>
            <BaseFormItem label="文件导入备用">
              <div class="space-y-2">
                <input class="block w-full text-sm" type="file" accept=".csv,.json,text/csv,application/json" @change="onPickOrderFile" />
                <p v-if="orderFileName" class="m-0 text-xs text-[var(--color-text-3)]">{{ orderFileName }}</p>
                <BaseButton variant="secondary" class="w-full" :loading="syncingOrders" @click="syncOrders">上传订单文件</BaseButton>
              </div>
            </BaseFormItem>
            <BaseFormItem label="订单 JSON">
              <BaseTextarea v-model="orderJson" :rows="8" placeholder='[{"order_no":"MT001","order_time":"2026-06-14 12:30:00","product_summary":"2L精酿","actual_amount":68}]' />
            </BaseFormItem>
            <BaseButton variant="primary" class="w-full" :loading="importingOrders" @click="submitOrders">导入订单</BaseButton>
            <BaseFormItem label="评价 JSON">
              <BaseTextarea v-model="reviewJson" :rows="8" placeholder='[{"review_id":"R001","rating":3,"content":"配送有点慢","review_time":"2026-06-14 13:00:00"}]' />
            </BaseFormItem>
            <BaseButton variant="primary" class="w-full" :loading="importingReviews" @click="submitReviews">导入评价</BaseButton>
          </div>
        </BaseCard>

        <BaseCard>
          <template #header><span class="font-semibold text-slate-800">运营动作</span></template>
          <div class="space-y-2">
            <BaseButton variant="primary" class="w-full" :loading="generating" @click="generateSuggestions">生成AI运营建议</BaseButton>
            <BaseButton variant="secondary" class="w-full" @click="openAccountDlg()">新增/编辑美团店铺</BaseButton>
          </div>
        </BaseCard>
      </div>

      <div class="flex flex-col gap-4 min-w-0">
        <BaseCard>
          <template #header><span class="font-semibold text-slate-800">AI建议</span></template>
          <div v-if="suggestions.length === 0" class="text-sm text-[var(--color-text-3)]">暂无建议，导入订单/评价后点击生成。</div>
          <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-3">
            <div v-for="item in suggestions" :key="item.id" class="rounded border border-[var(--color-border-2)] p-3 bg-[var(--color-bg-2)]">
              <div class="flex items-start justify-between gap-2">
                <div class="min-w-0">
                  <div class="font-semibold text-slate-800">{{ item.title }}</div>
                  <div class="mt-1 text-xs text-[var(--color-text-3)]">影响分 {{ item.impact_score }} · {{ suggestionTypeLabel(item.type) }}</div>
                </div>
                <span class="text-xs px-2 py-1 rounded bg-[var(--color-fill-2)]">{{ statusLabel(item.status) }}</span>
              </div>
              <p class="mt-3 mb-2 text-sm text-slate-600 whitespace-pre-wrap">{{ item.reason }}</p>
              <p class="m-0 text-sm text-slate-800 whitespace-pre-wrap">{{ item.content }}</p>
              <div class="mt-3 flex flex-wrap justify-end gap-2">
                <BaseButton v-if="item.status === 'pending'" size="sm" variant="secondary" @click="setSuggestionStatus(item, 'ignored')">忽略</BaseButton>
                <BaseButton v-if="item.status === 'pending'" size="sm" variant="primary" @click="setSuggestionStatus(item, 'approved')">确认执行</BaseButton>
                <BaseButton v-if="item.status === 'approved'" size="sm" variant="primary" @click="setSuggestionStatus(item, 'done')">标记完成</BaseButton>
              </div>
            </div>
          </div>
        </BaseCard>

        <BaseCard>
          <template #header><span class="font-semibold text-slate-800">最近订单</span></template>
          <BaseTable :columns="orderColumns" :data="(orders as unknown) as Record<string, unknown>[]" min-width="820px">
            <template #cell-order_time="{ row }">{{ formatDateTime((row as MeituanAIOrder).order_time) }}</template>
            <template #cell-actual_amount="{ row }">{{ formatMoney((row as MeituanAIOrder).actual_amount) }}</template>
          </BaseTable>
        </BaseCard>

        <BaseCard>
          <template #header><span class="font-semibold text-slate-800">最近评价</span></template>
          <BaseTable :columns="reviewColumns" :data="(reviews as unknown) as Record<string, unknown>[]" min-width="760px">
            <template #cell-rating="{ row }">{{ (row as MeituanAIReview).rating }} 星</template>
            <template #cell-tags="{ row }">{{ (row as MeituanAIReview).tags || '-' }}</template>
            <template #cell-review_time="{ row }">{{ formatDateTime((row as MeituanAIReview).review_time) }}</template>
          </BaseTable>
        </BaseCard>
      </div>
    </div>

    <BaseDialog v-model="accountDlg" title="美团店铺配置" max-width="min(560px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="店铺名称" required>
          <BaseInput v-model="accountForm.shop_name" placeholder="如 井科伟美团外卖店" />
        </BaseFormItem>
        <BaseFormItem label="美团店铺ID">
          <BaseInput v-model="accountForm.shop_id" placeholder="后续官方API授权使用" />
        </BaseFormItem>
        <BaseFormItem label="登录账号">
          <BaseInput v-model="accountForm.login_name" placeholder="仅备注，不自动登录" />
        </BaseFormItem>
        <BaseFormItem label="DeveloperId">
          <BaseInput v-model="accountForm.developer_id" placeholder="美团开放平台 DeveloperId" />
        </BaseFormItem>
        <BaseFormItem label="SignKey">
          <BaseInput v-model="accountForm.sign_key" placeholder="美团开放平台 SignKey" />
        </BaseFormItem>
        <BaseFormItem label="appAuthToken">
          <BaseInput v-model="accountForm.app_auth_token" placeholder="门店授权后获得的 appAuthToken" />
        </BaseFormItem>
        <BaseFormItem label="接口地址">
          <BaseInput v-model="accountForm.api_base_url" placeholder="https://api-open-cater.meituan.com" />
        </BaseFormItem>
        <BaseFormItem label="订单接口路径">
          <BaseInput v-model="accountForm.query_order_path" placeholder="/api/order/queryById" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="accountForm.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="accountDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="savingAccount" @click="submitAccount">保存</BaseButton>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseCard,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseSelect,
  BaseTable,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import {
  createMeituanAIAccount,
  generateMeituanAISuggestions,
  getMeituanAIDashboard,
  importMeituanAIOrders,
  importMeituanAIReviews,
  listMeituanAIAccounts,
  listMeituanAIOrders,
  listMeituanAIReviews,
  listMeituanAISuggestions,
  syncMeituanAIOpenAPIOrders,
  syncMeituanAIOrders,
  updateMeituanAIAccount,
  updateMeituanAISuggestionStatus,
} from '@/api/meituanAI'
import type { MeituanAIAccount, MeituanAIDashboard, MeituanAIOrder, MeituanAIReview, MeituanAISuggestion } from '@/api/types'
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
const activeAccountId = ref<number | undefined>(undefined)
const accounts = ref<MeituanAIAccount[]>([])
const dashboard = ref<MeituanAIDashboard | null>(null)
const orders = ref<MeituanAIOrder[]>([])
const reviews = ref<MeituanAIReview[]>([])
const suggestions = ref<MeituanAISuggestion[]>([])

const accountOptions = computed(() =>
  accounts.value.map((a) => ({ label: a.shop_name, value: a.id })),
)

const statCards = computed(() => [
  { label: '美团销售额', value: formatMoney(dashboard.value?.sales_amount), tone: 'text-indigo-700' },
  { label: '订单数', value: String(dashboard.value?.order_count ?? 0), tone: 'text-slate-800' },
  { label: '客单价', value: formatMoney(dashboard.value?.avg_order_amount), tone: 'text-emerald-700' },
  { label: '待处理建议', value: String(dashboard.value?.pending_suggestions ?? 0), tone: 'text-orange-600' },
  { label: '平台费用', value: formatMoney(dashboard.value?.platform_fee), tone: 'text-slate-800' },
  { label: '退款金额', value: formatMoney(dashboard.value?.refund_amount), tone: 'text-rose-600' },
  { label: '评价数', value: String(dashboard.value?.review_count ?? 0), tone: 'text-slate-800' },
  { label: '差评率', value: `${Number(dashboard.value?.negative_rate ?? 0).toFixed(1)}%`, tone: 'text-rose-600' },
])

const orderColumns: BaseTableColumn[] = [
  { key: 'order_no', label: '订单号', prop: 'order_no', minWidth: '140px', ellipsis: true },
  { key: 'order_time', label: '时间', width: '160px' },
  { key: 'product_summary', label: '商品', prop: 'product_summary', minWidth: '220px', ellipsis: true },
  { key: 'actual_amount', label: '实收', width: '90px', align: 'right' },
  { key: 'status', label: '状态', prop: 'status', width: '100px' },
]

const reviewColumns: BaseTableColumn[] = [
  { key: 'rating', label: '评分', width: '80px' },
  { key: 'content', label: '评价内容', prop: 'content', minWidth: '260px', ellipsis: true },
  { key: 'tags', label: '标签', width: '140px' },
  { key: 'review_time', label: '时间', width: '160px' },
]

async function loadAccounts(): Promise<void> {
  accounts.value = await listMeituanAIAccounts()
  if (!activeAccountId.value && accounts.value.length) {
    activeAccountId.value = accounts.value[0].id
  }
}

function queryParams(): Record<string, unknown> {
  return {
    store_id: tenantStoreId.value,
    account_id: activeAccountId.value,
    start_date: rangeStart.value,
    end_date: rangeEnd.value,
  }
}

async function reloadAll(): Promise<void> {
  await loadAccounts()
  if (!activeAccountId.value) {
    dashboard.value = null
    orders.value = []
    reviews.value = []
    suggestions.value = []
    return
  }
  const params = queryParams()
  const [dash, orderPage, reviewPage, suggestionPage] = await Promise.all([
    getMeituanAIDashboard(params),
    listMeituanAIOrders({ ...params, page: 1, page_size: 10 }),
    listMeituanAIReviews({ ...params, page: 1, page_size: 10 }),
    listMeituanAISuggestions({ ...params, page: 1, page_size: 20 }),
  ])
  dashboard.value = dash
  orders.value = orderPage.list
  reviews.value = reviewPage.list
  suggestions.value = suggestionPage.list
  await qc.invalidateQueries({ queryKey: ['meituan-ai'] })
}

const orderJson = ref('')
const reviewJson = ref('')
const syncOrderIds = ref('')
const orderFile = ref<File | undefined>(undefined)
const orderFileName = ref('')
const syncingOpenAPIOrders = ref(false)
const syncingOrders = ref(false)
const importingOrders = ref(false)
const importingReviews = ref(false)
const generating = ref(false)

function parseJsonArray(v: string): unknown[] {
  const data = JSON.parse(v)
  if (!Array.isArray(data)) throw new Error('请填写 JSON 数组')
  return data
}

function onPickOrderFile(e: Event): void {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  orderFile.value = file
  orderFileName.value = file?.name ?? ''
}

async function syncOpenAPIOrders(): Promise<void> {
  if (!activeAccountId.value) {
    toast.warning('请先配置美团店铺')
    return
  }
  const orderIds = syncOrderIds.value
    .split(/[\n,，\s]+/)
    .map((v) => v.trim())
    .filter(Boolean)
  if (orderIds.length === 0) {
    toast.warning('请输入美团订单号')
    return
  }
  syncingOpenAPIOrders.value = true
  try {
    const res = await syncMeituanAIOpenAPIOrders(activeAccountId.value, { order_ids: orderIds })
    toast.success(`已从美团同步 ${res.imported} 条订单${res.skipped ? `，失败/跳过 ${res.skipped} 条` : ''}`)
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '同步失败')
  } finally {
    syncingOpenAPIOrders.value = false
  }
}

async function syncOrders(): Promise<void> {
  if (!activeAccountId.value) {
    toast.warning('请先配置美团店铺')
    return
  }
  if (!orderFile.value) {
    toast.warning('请选择美团订单导出文件')
    return
  }
  const fd = new FormData()
  fd.append('file', orderFile.value)
  syncingOrders.value = true
  try {
    const res = await syncMeituanAIOrders(activeAccountId.value, fd)
    toast.success(`已同步 ${res.imported} 条订单${res.skipped ? `，跳过 ${res.skipped} 行` : ''}`)
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '同步失败')
  } finally {
    syncingOrders.value = false
  }
}

async function submitOrders(): Promise<void> {
  if (!activeAccountId.value) {
    toast.warning('请先配置美团店铺')
    return
  }
  importingOrders.value = true
  try {
    const orders = parseJsonArray(orderJson.value)
    const res = await importMeituanAIOrders(activeAccountId.value, { orders })
    toast.success(`已导入 ${res.imported} 条订单`)
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '导入失败')
  } finally {
    importingOrders.value = false
  }
}

async function submitReviews(): Promise<void> {
  if (!activeAccountId.value) {
    toast.warning('请先配置美团店铺')
    return
  }
  importingReviews.value = true
  try {
    const reviews = parseJsonArray(reviewJson.value)
    const res = await importMeituanAIReviews(activeAccountId.value, { reviews })
    toast.success(`已导入 ${res.imported} 条评价`)
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '导入失败')
  } finally {
    importingReviews.value = false
  }
}

async function generateSuggestions(): Promise<void> {
  if (!activeAccountId.value) {
    toast.warning('请先配置美团店铺')
    return
  }
  generating.value = true
  try {
    const res = await generateMeituanAISuggestions(activeAccountId.value, queryParams())
    toast.success(`已生成 ${res.generated} 条建议${res.ai_enabled ? '（DeepSeek）' : '（规则引擎）'}`)
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '生成失败')
  } finally {
    generating.value = false
  }
}

async function setSuggestionStatus(item: MeituanAISuggestion, status: MeituanAISuggestion['status']): Promise<void> {
  try {
    await updateMeituanAISuggestionStatus(item.id, status)
    toast.success('已更新')
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '更新失败')
  }
}

const accountDlg = ref(false)
const savingAccount = ref(false)
const editAccountId = ref(0)
const accountForm = reactive({
  shop_name: '',
  shop_id: '',
  login_name: '',
  developer_id: '',
  sign_key: '',
  app_auth_token: '',
  business_id: 2,
  api_version: '2',
  api_base_url: 'https://api-open-cater.meituan.com',
  query_order_path: '/api/order/queryById',
  remark: '',
})

function openAccountDlg(row?: MeituanAIAccount): void {
  const target = row ?? accounts.value.find((a) => a.id === activeAccountId.value)
  editAccountId.value = target?.id ?? 0
  accountForm.shop_name = target?.shop_name ?? ''
  accountForm.shop_id = target?.shop_id ?? ''
  accountForm.login_name = target?.login_name ?? ''
  accountForm.developer_id = target?.developer_id ?? ''
  accountForm.sign_key = target?.sign_key ?? ''
  accountForm.app_auth_token = target?.app_auth_token ?? ''
  accountForm.business_id = target?.business_id ?? 2
  accountForm.api_version = target?.api_version ?? '2'
  accountForm.api_base_url = target?.api_base_url ?? 'https://api-open-cater.meituan.com'
  accountForm.query_order_path = target?.query_order_path ?? '/api/order/queryById'
  accountForm.remark = target?.remark ?? ''
  accountDlg.value = true
}

async function submitAccount(): Promise<void> {
  if (!accountForm.shop_name.trim()) {
    toast.warning('请填写店铺名称')
    return
  }
  savingAccount.value = true
  try {
    const payload = {
      store_id: tenantStoreId.value,
      shop_name: accountForm.shop_name.trim(),
      shop_id: accountForm.shop_id.trim(),
      login_name: accountForm.login_name.trim(),
      developer_id: accountForm.developer_id.trim(),
      sign_key: accountForm.sign_key.trim(),
      app_auth_token: accountForm.app_auth_token.trim(),
      business_id: accountForm.business_id,
      api_version: accountForm.api_version.trim() || '2',
      api_base_url: accountForm.api_base_url.trim() || 'https://api-open-cater.meituan.com',
      query_order_path: accountForm.query_order_path.trim() || '/api/order/queryById',
      remark: accountForm.remark.trim(),
      is_enabled: true,
    }
    if (editAccountId.value > 0) {
      await updateMeituanAIAccount(editAccountId.value, payload)
    } else {
      const row = await createMeituanAIAccount(payload)
      activeAccountId.value = row.id
    }
    accountDlg.value = false
    toast.success('已保存')
    await reloadAll()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    savingAccount.value = false
  }
}

function fillDemoData(): void {
  orderJson.value = JSON.stringify(
    [
      {
        order_no: `MT${Date.now()}`,
        order_time: `${rangeEnd.value} 12:30:00`,
        customer_name: '顾客A',
        product_summary: '2L精酿、杯具',
        original_amount: 78,
        actual_amount: 68,
        delivery_fee: 3,
        platform_fee: 5,
        refund_amount: 0,
        status: '已完成',
      },
    ],
    null,
    2,
  )
  reviewJson.value = JSON.stringify(
    [
      {
        review_id: `R${Date.now()}`,
        order_no: 'MT示例',
        rating: 3,
        content: '配送有点慢，包装有点洒',
        review_time: `${rangeEnd.value} 13:00:00`,
      },
    ],
    null,
    2,
  )
}

function formatMoney(v: number | undefined | null): string {
  const n = Number(v ?? 0)
  return Number.isFinite(n) ? n.toFixed(2) : '0.00'
}

function formatDateTime(v: string | undefined): string {
  const s = String(v || '')
  if (!s) return '-'
  return s.slice(0, 19).replace('T', ' ')
}

function statusLabel(status: MeituanAISuggestion['status']): string {
  return { pending: '待确认', approved: '已确认', done: '已完成', ignored: '已忽略' }[status] || status
}

function suggestionTypeLabel(type: string): string {
  return {
    data: '数据',
    review: '评价',
    bundle: '套餐',
    profit: '利润',
    product: '商品',
    reply: '回复',
    routine: '日常',
  }[type] || type
}

watch([activeAccountId, rangeStart, rangeEnd], () => {
  void reloadAll()
})

onMounted(() => {
  void reloadAll()
})
</script>
