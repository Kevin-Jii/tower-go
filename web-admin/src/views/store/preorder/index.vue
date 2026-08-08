<template>
  <div class="flex min-h-0 flex-1 flex-col gap-3">
    <div class="flex flex-col gap-3 border-b border-[var(--color-border-2)] pb-3 lg:flex-row lg:items-center lg:justify-between">
      <div>
        <h2 class="m-0 text-lg font-semibold text-[var(--color-text-1)]">预订单</h2>
        <div class="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-[var(--color-text-3)]">
          <span class="inline-flex items-center gap-1"><IconNotification /> 提前一天</span>
          <span>09:30</span>
          <span>16:00</span>
          <span class="text-[var(--color-border-3)]">|</span>
          <span>配送当天</span>
          <span>09:30</span>
          <span>16:00</span>
        </div>
      </div>
      <BaseButton v-permission="'preorder:add'" variant="primary" @click="openCreate">
        <IconPlus class="mr-1" />新增预订单
      </BaseButton>
    </div>

    <div class="flex flex-wrap items-center gap-2">
      <BaseInput v-model="filters.keyword" class="w-full sm:w-60" placeholder="单号 / 客户 / 电话 / 地址" clearable @enter="reload" />
      <BaseSelect v-model="filters.status" class="w-full sm:w-36" :options="statusFilterOptions" placeholder="全部状态" />
      <a-date-picker v-model="filters.start_date" value-format="YYYY-MM-DD" class="w-full sm:w-36" placeholder="开始日期" />
      <a-date-picker v-model="filters.end_date" value-format="YYYY-MM-DD" class="w-full sm:w-36" placeholder="结束日期" />
      <BaseButton variant="secondary" @click="reload"><IconSearch class="mr-1" />查询</BaseButton>
    </div>

    <BaseTable
      :columns="columns"
      :data="(rows as unknown) as Record<string, unknown>[]"
      :loading="loading"
      min-width="1160px"
      class="min-h-0 flex-1"
    >
      <template #cell-scheduled_at="{ row }">
        <div class="leading-tight">
          <div :class="scheduleClass(row as PreOrder)" class="font-semibold">{{ formatDateTime((row as PreOrder).scheduled_at) }}</div>
          <div class="mt-1 text-xs text-[var(--color-text-3)]">{{ scheduleLabel(row as PreOrder) }}</div>
        </div>
      </template>
      <template #cell-customer="{ row }">
        <div class="min-w-0">
          <div class="truncate font-medium text-[var(--color-text-1)]">{{ (row as PreOrder).customer_name }}</div>
          <div class="mt-1 truncate text-xs text-[var(--color-text-3)]">
            {{ contactText(row as PreOrder) }}
          </div>
        </div>
      </template>
      <template #cell-items="{ row }">
        <div class="min-w-0 text-sm">
          <div class="truncate text-[var(--color-text-1)]">{{ itemSummary(row as PreOrder) }}</div>
          <div v-if="((row as PreOrder).items?.length || 0) > 2" class="mt-1 text-xs text-[var(--color-text-3)]">
            共 {{ (row as PreOrder).items?.length }} 项
          </div>
        </div>
      </template>
      <template #cell-delivery_address="{ row }">
        <span class="line-clamp-2 text-sm text-[var(--color-text-2)]">{{ (row as PreOrder).delivery_address || '-' }}</span>
      </template>
      <template #cell-status="{ row }">
        <BaseTag :variant="statusVariant((row as PreOrder).status)">{{ statusText((row as PreOrder).status) }}</BaseTag>
      </template>
      <template #cell-reminders="{ row }">
        <div class="text-xs">
          <span class="font-medium text-[var(--color-text-2)]">已提醒 {{ sentReminderCount(row as PreOrder) }} 次</span>
          <span v-if="failedReminderCount(row as PreOrder)" class="ml-2 text-rose-600">{{ failedReminderCount(row as PreOrder) }} 次失败</span>
        </div>
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="rowActions(row as PreOrder)" :max-inline="2" />
      </template>
    </BaseTable>

    <div class="flex shrink-0 justify-end">
      <BasePagination
        :page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="(value) => (page = value)"
        @update:page-size="(value) => (pageSize = value)"
      />
    </div>

    <a-drawer
      :visible="drawerOpen"
      placement="right"
      width="min(780px, 96vw)"
      :mask-closable="false"
      unmount-on-close
      @cancel="drawerOpen = false"
      @update:visible="drawerOpen = $event"
    >
      <template #title>{{ editingId ? '编辑预订单' : '新增预订单' }}</template>
      <div class="space-y-5">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-2">
          <BaseFormItem label="客户" required>
            <BaseSelect v-model="form.customer_id" :options="customerOptions" searchable placeholder="选择客户" @update:model-value="onCustomerChange" />
          </BaseFormItem>
          <BaseFormItem label="计划配送时间" required>
            <a-date-picker
              v-model="form.scheduled_at"
              show-time
              format="YYYY-MM-DD HH:mm"
              value-format="YYYY-MM-DD HH:mm:ss"
              class="w-full"
              placeholder="选择配送时间"
            />
          </BaseFormItem>
          <BaseFormItem label="联系人">
            <BaseInput v-model="form.contact_person" />
          </BaseFormItem>
          <BaseFormItem label="联系电话">
            <BaseInput v-model="form.contact_phone" />
          </BaseFormItem>
          <BaseFormItem label="配送地址" class="md:col-span-2">
            <BaseInput v-model="form.delivery_address" />
          </BaseFormItem>
        </div>

        <div>
          <div class="mb-2 flex items-center justify-between">
            <h3 class="m-0 text-sm font-semibold text-[var(--color-text-1)]">预订商品</h3>
            <BaseButton variant="secondary" size="sm" :disabled="!form.customer_id" @click="addLine">
              <IconPlus class="mr-1" />添加商品
            </BaseButton>
          </div>
          <div class="overflow-hidden rounded border border-[var(--color-border-2)]">
            <div class="hidden grid-cols-[1.3fr_1.1fr_100px_1fr_40px] gap-2 bg-[var(--color-fill-2)] px-3 py-2 text-xs font-medium text-[var(--color-text-3)] md:grid">
              <span>商品</span><span>规格</span><span>数量</span><span>明细备注</span><span></span>
            </div>
            <div v-if="form.items.length" class="divide-y divide-[var(--color-border-2)]">
              <div v-for="(line, index) in form.items" :key="line.key" class="grid grid-cols-1 gap-2 px-3 py-3 md:grid-cols-[1.3fr_1.1fr_100px_1fr_40px] md:items-center">
                <BaseSelect v-model="line.product_id" :options="productOptions" searchable placeholder="商品" @update:model-value="onProductChange(line)" />
                <BaseSelect v-model="line.unit_spec_id" :options="specOptions(line.product_id)" placeholder="规格" :disabled="!line.product_id" />
                <BaseNumberInput v-model="line.quantity" :min="0.01" :step="1" placeholder="数量" />
                <BaseInput v-model="line.remark" placeholder="备注" />
                <a-tooltip content="移除商品">
                  <BaseButton variant="ghost" size="sm" :disabled="form.items.length <= 1" aria-label="移除商品" @click="removeLine(index)">
                    <IconDelete />
                  </BaseButton>
                </a-tooltip>
              </div>
            </div>
            <div v-else class="px-4 py-8 text-center text-sm text-[var(--color-text-3)]">暂无商品</div>
          </div>
        </div>

        <BaseFormItem label="整单备注">
          <BaseTextarea v-model="form.remark" :rows="3" />
        </BaseFormItem>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <BaseButton variant="ghost" @click="drawerOpen = false">取消</BaseButton>
          <BaseButton variant="primary" :loading="saving" @click="submit">保存预订单</BaseButton>
        </div>
      </template>
    </a-drawer>

    <BaseDialog v-model="detailOpen" title="预订单详情" max-width="min(760px, 96vw)">
      <div v-if="detail" class="space-y-4">
        <div class="grid grid-cols-1 gap-x-6 gap-y-3 border-b border-[var(--color-border-2)] pb-4 text-sm md:grid-cols-2">
          <div><span class="detail-label">预订单号</span>{{ detail.order_no }}</div>
          <div><span class="detail-label">状态</span><BaseTag :variant="statusVariant(detail.status)">{{ statusText(detail.status) }}</BaseTag></div>
          <div><span class="detail-label">客户</span>{{ detail.customer_name }}</div>
          <div><span class="detail-label">配送时间</span>{{ formatDateTime(detail.scheduled_at) }}</div>
          <div><span class="detail-label">联系人</span>{{ contactText(detail) }}</div>
          <div><span class="detail-label">创建人</span>{{ detail.creator?.nickname || detail.creator?.username || '-' }}</div>
          <div class="md:col-span-2"><span class="detail-label">配送地址</span>{{ detail.delivery_address || '-' }}</div>
          <div v-if="detail.remark" class="md:col-span-2"><span class="detail-label">备注</span>{{ detail.remark }}</div>
        </div>
        <BaseTable :columns="detailColumns" :data="(detail.items || []) as unknown as Record<string, unknown>[]" min-width="560px" />
        <div>
          <h3 class="m-0 mb-2 text-sm font-semibold">钉钉提醒</h3>
          <div class="grid grid-cols-1 gap-2 text-xs sm:grid-cols-2">
            <div v-for="item in reminderTimeline(detail)" :key="item.key" class="flex items-center justify-between border-b border-[var(--color-border-2)] py-2">
              <span class="text-[var(--color-text-2)]">{{ item.label }}</span>
              <span :class="item.className">{{ item.status }}</span>
            </div>
          </div>
        </div>
      </div>
      <template #footer><BaseButton variant="ghost" @click="detailOpen = false">关闭</BaseButton></template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { IconDelete, IconNotification, IconPlus, IconSearch } from '@arco-design/web-vue/es/icon'
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
  BaseTag,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn, TableRowAction } from '@/components/base/types'
import { listB2BCustomers, listB2BPrices } from '@/api/b2b'
import { createPreOrder, deletePreOrder, listPreOrders, updatePreOrder, updatePreOrderStatus } from '@/api/preorder'
import type { B2BCustomer, PreOrder } from '@/api/types'
import { confirmDialog } from '@/feedback/confirm'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

interface FormLine {
  key: number
  product_id: number | undefined
  unit_spec_id: number | undefined
  quantity: number | undefined
  remark: string
}

const qc = useQueryClient()
const userStore = useUserStore()
const storeId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)
const page = ref(1)
const pageSize = ref(20)
const filters = reactive({ keyword: '', status: '' as number | '', start_date: '', end_date: '' })
const queryParams = computed(() => ({
  store_id: storeId.value,
  keyword: filters.keyword.trim() || undefined,
  status: filters.status || undefined,
  start_date: filters.start_date || undefined,
  end_date: filters.end_date || undefined,
  page: page.value,
  page_size: pageSize.value,
}))
const { data: pageData, isLoading: loading } = useQuery({
  queryKey: computed(() => ['pre-orders', queryParams.value] as const),
  queryFn: () => listPreOrders(queryParams.value),
})
const rows = computed(() => pageData.value?.list || [])
const total = computed(() => pageData.value?.total || 0)

const statusFilterOptions: BaseSelectOption[] = [
  { label: '全部状态', value: '' },
  { label: '待备货', value: 1 },
  { label: '已备货', value: 2 },
  { label: '已配送', value: 3 },
  { label: '已取消', value: 4 },
]
const columns: BaseTableColumn[] = [
  { key: 'scheduled_at', label: '计划配送', width: '150px' },
  { key: 'customer', label: '客户', minWidth: '150px' },
  { key: 'items', label: '商品明细', minWidth: '220px' },
  { key: 'delivery_address', label: '配送地址', minWidth: '190px' },
  { key: 'status', label: '状态', width: '88px' },
  { key: 'reminders', label: '提醒记录', width: '150px' },
  { key: 'actions', label: '操作', width: '190px', align: 'right', fixed: 'right' },
]
const detailColumns: BaseTableColumn[] = [
  { key: 'product_name', label: '商品', prop: 'product_name', minWidth: '220px' },
  { key: 'unit_name', label: '规格', prop: 'unit_name', width: '120px' },
  { key: 'quantity', label: '数量', prop: 'quantity', width: '100px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '140px' },
]

function reload(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['pre-orders'] })
}

const { data: customerData } = useQuery({
  queryKey: computed(() => ['preorder-customers', storeId.value] as const),
  queryFn: () => listB2BCustomers({ store_id: storeId.value, status: 1, page: 1, page_size: 100 }),
})
const customers = computed(() => customerData.value?.list || [])
const customerOptions = computed<BaseSelectOption[]>(() => customers.value.map((item) => ({
  label: item.contact_person ? `${item.name} · ${item.contact_person}` : item.name,
  value: item.id,
})))

const drawerOpen = ref(false)
const editingId = ref(0)
const saving = ref(false)
let nextLineKey = 1
const form = reactive({
  customer_id: undefined as number | undefined,
  scheduled_at: '',
  contact_person: '',
  contact_phone: '',
  delivery_address: '',
  remark: '',
  items: [] as FormLine[],
})

const { data: priceData, isFetching: priceLoading } = useQuery({
  queryKey: computed(() => ['preorder-customer-products', storeId.value, form.customer_id] as const),
  queryFn: () => listB2BPrices({ store_id: storeId.value, customer_id: form.customer_id, page: 1, page_size: 100 }),
  enabled: computed(() => drawerOpen.value && Number(form.customer_id || 0) > 0),
})
const prices = computed(() => priceData.value?.list || [])
const productOptions = computed<BaseSelectOption[]>(() => {
  const seen = new Set<number>()
  const options: BaseSelectOption[] = []
  for (const item of prices.value) {
    if (seen.has(item.product_id)) continue
    seen.add(item.product_id)
    options.push({ label: item.product?.name || `商品#${item.product_id}`, value: item.product_id })
  }
  return options
})

function specOptions(productId?: number): BaseSelectOption[] {
  if (!productId) return []
  return prices.value
    .filter((item) => item.product_id === productId && item.is_enabled)
    .map((item) => ({ label: item.unit_name || item.unit_spec?.unit_name || `规格#${item.unit_spec_id}`, value: item.unit_spec_id }))
}

function blankLine(): FormLine {
  return { key: nextLineKey++, product_id: undefined, unit_spec_id: undefined, quantity: 1, remark: '' }
}

function resetForm(): void {
  form.customer_id = undefined
  form.scheduled_at = ''
  form.contact_person = ''
  form.contact_phone = ''
  form.delivery_address = ''
  form.remark = ''
  form.items = [blankLine()]
}

function defaultScheduledAt(): string {
  const date = new Date()
  date.setDate(date.getDate() + 1)
  date.setHours(12, 0, 0, 0)
  const pad = (value: number) => String(value).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:00`
}

function openCreate(): void {
  editingId.value = 0
  resetForm()
  form.scheduled_at = defaultScheduledAt()
  drawerOpen.value = true
}

function openEdit(row: PreOrder): void {
  editingId.value = row.id
  form.customer_id = row.customer_id
  form.scheduled_at = formatForPicker(row.scheduled_at)
  form.contact_person = row.contact_person || ''
  form.contact_phone = row.contact_phone || ''
  form.delivery_address = row.delivery_address || ''
  form.remark = row.remark || ''
  form.items = (row.items || []).map((item) => ({
    key: nextLineKey++,
    product_id: item.product_id,
    unit_spec_id: item.unit_spec_id,
    quantity: Number(item.quantity),
    remark: item.remark || '',
  }))
  if (!form.items.length) form.items = [blankLine()]
  drawerOpen.value = true
}

function onCustomerChange(value: string | number | undefined): void {
  const customer = customers.value.find((item) => item.id === Number(value))
  if (!customer) return
  fillCustomer(customer)
  form.items = [blankLine()]
}

function fillCustomer(customer: B2BCustomer): void {
  form.contact_person = customer.contact_person || ''
  form.contact_phone = customer.phone || ''
  form.delivery_address = customer.address || ''
}

function onProductChange(line: FormLine): void {
  line.unit_spec_id = undefined
  const options = specOptions(line.product_id)
  if (options.length === 1) line.unit_spec_id = Number(options[0].value)
}

function addLine(): void {
  form.items.push(blankLine())
}

function removeLine(index: number): void {
  if (form.items.length <= 1) return
  form.items.splice(index, 1)
}

async function submit(): Promise<void> {
  if (!form.customer_id) {
    toast.warning('请选择客户')
    return
  }
  if (!form.scheduled_at) {
    toast.warning('请选择计划配送时间')
    return
  }
  if (priceLoading.value) {
    toast.warning('客户商品仍在加载，请稍后')
    return
  }
  if (form.items.some((item) => !item.product_id || !item.unit_spec_id || !item.quantity || item.quantity <= 0)) {
    toast.warning('请完整填写商品、规格和数量')
    return
  }
  const unique = new Set(form.items.map((item) => `${item.product_id}:${item.unit_spec_id}`))
  if (unique.size !== form.items.length) {
    toast.warning('同一商品规格不能重复添加')
    return
  }
  const body = {
    store_id: storeId.value,
    customer_id: Number(form.customer_id),
    scheduled_at: form.scheduled_at,
    contact_person: form.contact_person.trim(),
    contact_phone: form.contact_phone.trim(),
    delivery_address: form.delivery_address.trim(),
    remark: form.remark.trim(),
    items: form.items.map((item) => ({
      product_id: Number(item.product_id),
      unit_spec_id: Number(item.unit_spec_id),
      quantity: Number(item.quantity),
      remark: item.remark.trim(),
    })),
  }
  saving.value = true
  try {
    if (editingId.value) await updatePreOrder(editingId.value, body)
    else await createPreOrder(body)
    toast.success(editingId.value ? '预订单已更新' : '预订单已创建')
    drawerOpen.value = false
    reload()
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '保存失败')
  } finally {
    saving.value = false
  }
}

const detailOpen = ref(false)
const detail = ref<PreOrder | null>(null)
function openDetail(row: PreOrder): void {
  detail.value = row
  detailOpen.value = true
}

async function changeStatus(row: PreOrder, status: number): Promise<void> {
  const label = statusText(status)
  const ok = await confirmDialog({ title: `标记为${label}`, message: `确认将 ${row.order_no} 标记为“${label}”吗？` })
  if (!ok) return
  try {
    await updatePreOrderStatus(row.id, status)
    toast.success(`已标记为${label}`)
    reload()
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '状态更新失败')
  }
}

async function removeOrder(row: PreOrder): Promise<void> {
  const ok = await confirmDialog({ title: '删除预订单', message: `确认删除 ${row.order_no} 吗？删除后无法恢复。` })
  if (!ok) return
  try {
    await deletePreOrder(row.id)
    toast.success('预订单已删除')
    reload()
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '删除失败')
  }
}

function rowActions(row: PreOrder): TableRowAction[] {
  const actions: TableRowAction[] = [
    { label: '详情', permission: 'preorder:list', onClick: () => openDetail(row), place: 'inline' },
  ]
  if (row.status === 1) {
    actions.push({ label: '已备货', permission: 'preorder:edit', onClick: () => void changeStatus(row, 2), place: 'inline' })
    actions.push({ label: '编辑', permission: 'preorder:edit', onClick: () => openEdit(row), place: 'more' })
    actions.push({ label: '取消', permission: 'preorder:edit', danger: true, onClick: () => void changeStatus(row, 4), place: 'more' })
    actions.push({ label: '删除', permission: 'preorder:delete', danger: true, onClick: () => void removeOrder(row), place: 'more' })
  } else if (row.status === 2) {
    actions.push({ label: '已配送', permission: 'preorder:edit', onClick: () => void changeStatus(row, 3), place: 'inline' })
    actions.push({ label: '编辑', permission: 'preorder:edit', onClick: () => openEdit(row), place: 'more' })
    actions.push({ label: '退回待备货', permission: 'preorder:edit', onClick: () => void changeStatus(row, 1), place: 'more' })
    actions.push({ label: '取消', permission: 'preorder:edit', danger: true, onClick: () => void changeStatus(row, 4), place: 'more' })
  } else if (row.status === 4) {
    actions.push({ label: '删除', permission: 'preorder:delete', danger: true, onClick: () => void removeOrder(row), place: 'more' })
  }
  return actions
}

function statusText(status: number): string {
  return ({ 1: '待备货', 2: '已备货', 3: '已配送', 4: '已取消' } as Record<number, string>)[status] || '未知'
}

function statusVariant(status: number): 'success' | 'warning' | 'info' | 'danger' | 'neutral' {
  if (status === 1) return 'warning'
  if (status === 2) return 'info'
  if (status === 3) return 'success'
  return 'neutral'
}

function formatDateTime(value?: string): string {
  if (!value) return '-'
  return value.replace('T', ' ').slice(0, 16)
}

function formatForPicker(value: string): string {
  const formatted = value.replace('T', ' ').slice(0, 19)
  return formatted.length === 16 ? `${formatted}:00` : formatted
}

function parseDate(value: string): Date {
  return new Date(value.replace(' ', 'T'))
}

function scheduleLabel(row: PreOrder): string {
  if (row.status === 3) return '配送完成'
  if (row.status === 4) return '已取消'
  const scheduled = parseDate(row.scheduled_at)
  const today = new Date()
  const start = new Date(today.getFullYear(), today.getMonth(), today.getDate()).getTime()
  const target = new Date(scheduled.getFullYear(), scheduled.getMonth(), scheduled.getDate()).getTime()
  const days = Math.round((target - start) / 86400000)
  if (days < 0) return `已逾期 ${Math.abs(days)} 天`
  if (days === 0) return '今天配送'
  if (days === 1) return '明天配送'
  return `${days} 天后配送`
}

function scheduleClass(row: PreOrder): string {
  if ((row.status === 1 || row.status === 2) && parseDate(row.scheduled_at).getTime() < Date.now()) return 'text-rose-600'
  return 'text-[var(--color-text-1)]'
}

function contactText(row: PreOrder): string {
  return [row.contact_person, row.contact_phone].filter(Boolean).join(' · ') || '-'
}

function itemSummary(row: PreOrder): string {
  const items = row.items || []
  if (!items.length) return '-'
  return items.slice(0, 2).map((item) => `${item.product_name} ${item.quantity}${item.unit_name}`).join('；')
}

function sentReminderCount(row: PreOrder): number {
  return (row.reminder_logs || []).filter((item) => item.status === 2).length
}

function failedReminderCount(row: PreOrder): number {
  return (row.reminder_logs || []).filter((item) => item.status === 3).length
}

const reminderLabels: Record<string, string> = {
  previous_day_0930: '提前一天 09:30',
  previous_day_1600: '提前一天 16:00',
  due_day_0930: '配送当天 09:30',
  due_day_1600: '配送当天 16:00',
}

function reminderTimeline(row: PreOrder): Array<{ key: string; label: string; status: string; className: string }> {
  const logs = new Map((row.reminder_logs || []).map((item) => [item.reminder_key, item]))
  return Object.entries(reminderLabels).map(([key, label]) => {
    const log = logs.get(key)
    if (log?.status === 2) return { key, label, status: formatDateTime(log.sent_at), className: 'text-emerald-700' }
    if (log?.status === 3) return { key, label, status: '发送失败', className: 'text-rose-600' }
    if (log?.status === 1) return { key, label, status: '发送中', className: 'text-blue-600' }
    return { key, label, status: '待提醒', className: 'text-[var(--color-text-3)]' }
  })
}
</script>

<style scoped>
.detail-label {
  display: inline-block;
  min-width: 72px;
  color: var(--color-text-3);
}
</style>
