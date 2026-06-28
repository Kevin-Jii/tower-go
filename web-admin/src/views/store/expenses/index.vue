<template>
  <div class="flex min-h-0 flex-1 flex-col gap-4">
    <BaseCard flush-body class="flex min-h-0 flex-1 flex-col">
      <template #header>
        <div class="flex w-full flex-col gap-3 xl:flex-row xl:items-center xl:justify-between">
          <div>
            <h2 class="m-0 text-base font-semibold text-slate-900">门店支出</h2>
            <p class="m-0 mt-1 text-xs text-slate-500">记录外卖推广、平台费用、维修维护等已支付支出</p>
          </div>
          <div class="flex flex-wrap items-center gap-2">
            <BaseInput v-model="filters.keyword" class="w-full sm:w-48" placeholder="单号 / 分类 / 备注 / 操作人" clearable @enter="reloadAll" />
            <BaseSelect v-model="filters.category_code" class="w-full sm:w-44" :options="categoryFilterOptions" placeholder="全部分类" />
            <a-date-picker v-model="filters.start_date" value-format="YYYY-MM-DD" class="w-full sm:w-36" />
            <a-date-picker v-model="filters.end_date" value-format="YYYY-MM-DD" class="w-full sm:w-36" />
            <BaseButton variant="primary" @click="reloadAll">查询</BaseButton>
            <BaseButton variant="secondary" @click="openExportDlg">导出Excel</BaseButton>
            <BaseButton v-permission="'store:expenses:add'" variant="primary" @click="openCreate">新增支出</BaseButton>
          </div>
        </div>
      </template>

      <div class="flex min-h-0 flex-1 flex-col gap-3 p-4">
        <div class="grid grid-cols-1 gap-3 md:grid-cols-3">
          <div v-for="item in summaryCards" :key="item.label" class="rounded border border-slate-200 bg-slate-50 px-4 py-3">
            <div class="text-xs font-medium text-slate-500">{{ item.label }}</div>
            <div class="mt-1 text-lg font-semibold text-slate-900">{{ item.value }}</div>
          </div>
        </div>

        <BaseTable :columns="columns" :data="(rows as unknown) as Record<string, unknown>[]" :loading="loading" min-width="1020px" class="min-h-0 flex-1">
          <template #cell-store="{ row }">
            {{ (row as StoreExpense).store?.name || '-' }}
          </template>
          <template #cell-expense_date="{ row }">
            {{ formatDate((row as StoreExpense).expense_date) }}
          </template>
          <template #cell-category_name="{ row }">
            <BaseTag variant="info">{{ (row as StoreExpense).category_name || '-' }}</BaseTag>
          </template>
          <template #cell-amount="{ row }">
            <span class="font-semibold text-rose-600">{{ formatMoney((row as StoreExpense).amount) }}</span>
          </template>
          <template #cell-created_at="{ row }">
            {{ formatDateTime((row as StoreExpense).created_at) }}
          </template>
          <template #cell-actions="{ row }">
            <BaseTableRowActions :actions="rowActions(row as StoreExpense)" :max-inline="2" />
          </template>
        </BaseTable>

        <div class="flex shrink-0 justify-end">
          <BasePagination
            :page="page"
            :page-size="pageSize"
            :total="total"
            @update:page="(p) => (page = p)"
            @update:page-size="(s) => (pageSize = s)"
          />
        </div>
      </div>
    </BaseCard>

    <BaseDialog v-model="dlg" :title="editingId ? '编辑支出' : '新增支出'" max-width="min(520px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="支出分类" required>
          <BaseSelect v-model="form.category_code" :options="categoryOptions" placeholder="请选择支出分类" searchable />
        </BaseFormItem>
        <BaseFormItem label="支出金额" required>
          <BaseNumberInput v-model="form.amount" :min="0.01" :step="0.01" placeholder="请输入支出金额" />
        </BaseFormItem>
        <BaseFormItem label="备注说明">
          <BaseTextarea v-model="form.remark" :rows="3" placeholder="例如：美团推广充值、平台服务费、维修维护说明" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submit">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="exportDlg" title="导出门店支出" max-width="min(420px, 96vw)">
      <BaseFormItem label="导出日期" required>
        <a-date-picker v-model="exportDate" value-format="YYYY-MM-DD" class="w-full" :allow-clear="false" />
      </BaseFormItem>
      <template #footer>
        <BaseButton variant="ghost" @click="exportDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="exporting" @click="submitExport">导出</BaseButton>
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
  BaseTableRowActions,
  BaseTag,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn, TableRowAction } from '@/components/base/types'
import { listDictDataByTypeCode } from '@/api/dict'
import {
  createStoreExpense,
  deleteStoreExpense,
  exportStoreExpenses,
  getStoreExpenseStats,
  listStoreExpenses,
  updateStoreExpense,
} from '@/api/storeExpense'
import type { StoreExpense } from '@/api/types'
import { confirmDialog } from '@/feedback/confirm'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

const qc = useQueryClient()
const userStore = useUserStore()
const tenantStoreId = computed(() => Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)

const filters = reactive({
  keyword: '',
  category_code: '',
  start_date: '',
  end_date: '',
})
const page = ref(1)
const pageSize = ref(10)
const dlg = ref(false)
const saving = ref(false)
const exportDlg = ref(false)
const exportDate = ref(new Date().toISOString().slice(0, 10))
const exporting = ref(false)
const editingId = ref<number | null>(null)
const form = reactive({
  category_code: '',
  amount: undefined as number | undefined,
  remark: '',
})

const queryParams = computed(() => ({
  page: page.value,
  page_size: pageSize.value,
  store_id: tenantStoreId.value,
  keyword: filters.keyword.trim() || undefined,
  category_code: filters.category_code || undefined,
  start_date: filters.start_date || undefined,
  end_date: filters.end_date || undefined,
}))

const { data: pageData, isLoading: loading } = useQuery({
  queryKey: computed(() => ['store-expenses', queryParams.value] as const),
  queryFn: () => listStoreExpenses(queryParams.value),
})

const { data: statsData } = useQuery({
  queryKey: computed(() => ['store-expense-stats', tenantStoreId.value, filters.category_code, filters.start_date, filters.end_date] as const),
  queryFn: () =>
    getStoreExpenseStats({
      store_id: tenantStoreId.value,
      category_code: filters.category_code || undefined,
      start_date: filters.start_date || undefined,
      end_date: filters.end_date || undefined,
    }),
})

const { data: categoryDict } = useQuery({
  queryKey: ['dict-data', 'EXPENDITURECLASS'] as const,
  queryFn: () => listDictDataByTypeCode('EXPENDITURECLASS'),
})

const rows = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)
const categoryOptions = computed<BaseSelectOption[]>(() =>
  (categoryDict.value ?? []).filter((item) => item.status === 1).map((item) => ({ label: item.label, value: item.value })),
)
const categoryFilterOptions = computed<BaseSelectOption[]>(() => [{ label: '全部分类', value: '' }, ...categoryOptions.value])
const summaryCards = computed(() => {
  const stats = statsData.value
  const avg = (stats?.count ?? 0) > 0 ? Number(stats?.total_amount || 0) / Number(stats?.count || 1) : 0
  return [
    { label: '支出总额', value: formatMoney(stats?.total_amount ?? 0) },
    { label: '支出记录', value: `${stats?.count ?? 0} 笔` },
    { label: '单笔均值', value: formatMoney(avg) },
  ]
})

const columns: BaseTableColumn[] = [
  { key: 'expense_no', label: '支出单号', prop: 'expense_no', minWidth: '170px', ellipsis: true },
  { key: 'store', label: '门店', minWidth: '140px', ellipsis: true },
  { key: 'expense_date', label: '支出日期', width: '120px' },
  { key: 'category_name', label: '分类', minWidth: '150px' },
  { key: 'amount', label: '金额', width: '120px', align: 'right' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '180px', ellipsis: true },
  { key: 'operator_name', label: '操作人', prop: 'operator_name', width: '110px' },
  { key: 'created_at', label: '创建时间', width: '170px' },
  { key: 'actions', label: '操作', width: '120px', align: 'right' },
]

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['store-expenses'] })
})

function reloadAll(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['store-expenses'] })
  void qc.invalidateQueries({ queryKey: ['store-expense-stats'] })
}

function openExportDlg(): void {
  exportDate.value = filters.start_date || filters.end_date || new Date().toISOString().slice(0, 10)
  exportDlg.value = true
}

async function submitExport(): Promise<void> {
  if (!exportDate.value) {
    toast.warning('请选择导出日期')
    return
  }
  exporting.value = true
  try {
    await exportStoreExpenses({ date: exportDate.value, store_id: tenantStoreId.value })
    exportDlg.value = false
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '导出失败')
  } finally {
    exporting.value = false
  }
}

function resetForm(): void {
  editingId.value = null
  form.category_code = ''
  form.amount = undefined
  form.remark = ''
}

function openCreate(): void {
  resetForm()
  dlg.value = true
}

function openEdit(row: StoreExpense): void {
  editingId.value = row.id
  form.category_code = row.category_code || ''
  form.amount = Number(row.amount || 0)
  form.remark = row.remark || ''
  dlg.value = true
}

async function submit(): Promise<void> {
  const categoryCode = String(form.category_code || '').trim()
  const amount = Number(form.amount || 0)
  if (!categoryCode) {
    toast.error('请选择支出分类')
    return
  }
  if (!Number.isFinite(amount) || amount <= 0) {
    toast.error('请输入大于 0 的支出金额')
    return
  }
  saving.value = true
  try {
    const body = {
      store_id: tenantStoreId.value,
      category_code: categoryCode,
      amount,
      remark: form.remark.trim(),
    }
    if (editingId.value) {
      await updateStoreExpense(editingId.value, body)
      toast.success('支出记录已更新')
    } else {
      await createStoreExpense(body)
      toast.success('支出记录已创建')
    }
    dlg.value = false
    reloadAll()
  } finally {
    saving.value = false
  }
}

async function onDelete(row: StoreExpense): Promise<void> {
  const ok = await confirmDialog({
    title: '删除支出记录',
    message: `确认删除支出单 ${row.expense_no} 吗？`,
  })
  if (!ok) return
  await deleteStoreExpense(row.id)
  toast.success('已删除')
  reloadAll()
}

function rowActions(row: StoreExpense): TableRowAction[] {
  return [
    { label: '编辑', permission: 'store:expenses:edit', onClick: () => openEdit(row), place: 'inline' },
    { label: '删除', permission: 'store:expenses:delete', danger: true, onClick: () => void onDelete(row), place: 'inline' },
  ]
}

function formatDate(value?: string): string {
  if (!value) return '-'
  return value.slice(0, 10)
}

function formatDateTime(value?: string): string {
  if (!value) return '-'
  return value.replace('T', ' ').slice(0, 19)
}

function formatMoney(value?: number): string {
  return `¥${Number(value || 0).toFixed(2)}`
}
</script>
