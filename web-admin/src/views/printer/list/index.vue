<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
      <div>
        <h2 class="page-title">打印机管理</h2>
        <p class="m-0 mt-1 text-sm text-slate-500">绑定门店打印机，维护默认设备、启停状态和在线检测。</p>
      </div>
      <div class="flex w-full flex-col gap-2 sm:flex-row md:w-auto">
        <BaseSelect
          v-if="storeOptions.length > 1"
          v-model="selectedStoreId"
          class="w-full sm:w-44"
          :options="storeOptions"
          placeholder="选择门店"
        />
        <BaseInput v-model="keyword" class="w-full sm:w-56" placeholder="SN / 名称 / 备注" clearable @enter="reload" />
        <BaseSelect v-model="statusFilter" class="w-full sm:w-32" :options="statusOptions" placeholder="状态" />
        <BaseButton variant="secondary" :loading="loading" @click="reload">刷新</BaseButton>
        <BaseButton variant="secondary" :loading="syncing" :disabled="!currentStoreId" @click="syncStatus">刷新状态</BaseButton>
        <BaseButton v-permission="'printer:bind'" variant="primary" @click="openCreate">绑定打印机</BaseButton>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-3 md:grid-cols-4">
      <div class="printer-stat">
        <span>打印机总数</span>
        <strong>{{ filteredRows.length }}</strong>
      </div>
      <div class="printer-stat">
        <span>在线设备</span>
        <strong class="text-emerald-600">{{ onlineCount }}</strong>
      </div>
      <div class="printer-stat">
        <span>默认设备</span>
        <strong class="text-blue-600">{{ defaultCount }}</strong>
      </div>
      <div class="printer-stat">
        <span>停用设备</span>
        <strong class="text-rose-600">{{ disabledCount }}</strong>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(pagedRows as unknown) as Record<string, unknown>[]" :loading="loading" min-width="1040px">
      <template #cell-store_id="{ row }">
        {{ storeName((row as PrinterRow).store_id) }}
      </template>
      <template #cell-type="{ row }">
        <a-tag :color="(row as PrinterRow).type === 2 ? 'purple' : 'blue'">
          {{ printerTypeLabel((row as PrinterRow).type) }}
        </a-tag>
      </template>
      <template #cell-status="{ row }">
        <BaseSwitch
          v-permission="'printer:edit'"
          :model-value="(row as PrinterRow).status"
          :active-value="1"
          :inactive-value="2"
          @update:model-value="toggleStatus(row as PrinterRow, Number($event))"
        />
      </template>
      <template #cell-is_default="{ row }">
        <BaseSwitch
          v-permission="'printer:edit'"
          :model-value="(row as PrinterRow).is_default"
          :active-value="1"
          :inactive-value="0"
          @update:model-value="toggleDefault(row as PrinterRow, Number($event))"
        />
      </template>
      <template #cell-online="{ row }">
        <a-tag :color="onlineColor((row as PrinterRow).online)">
          {{ onlineLabel((row as PrinterRow).online) }}
        </a-tag>
      </template>
      <template #cell-created_at="{ row }">
        {{ formatDate((row as PrinterRow).created_at) }}
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="rowActions(row as PrinterRow)" :max-inline="2" />
      </template>
    </BaseTable>

    <div class="flex justify-end">
      <BasePagination
        :page="page"
        :page-size="pageSize"
        :total="filteredRows.length"
        @update:page="(p) => (page = p)"
        @update:page-size="(s) => (pageSize = s)"
      />
    </div>

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑打印机' : '绑定打印机'" max-width="min(560px, 96vw)">
      <div class="max-h-[72vh] space-y-4 overflow-y-auto pr-1">
        <BaseFormItem label="所属门店" required>
          <BaseSelect v-model="form.store_id" :options="storeOptions" placeholder="请选择门店" :disabled="isEdit" />
        </BaseFormItem>
        <BaseFormItem label="打印机 SN" required>
          <BaseInput v-model="form.sn" placeholder="请输入打印机底部 SN" :disabled="isEdit" />
        </BaseFormItem>
        <BaseFormItem label="打印机名称">
          <BaseInput v-model="form.name" placeholder="如 前台小票机" />
        </BaseFormItem>
        <BaseFormItem label="打印机类型">
          <BaseSelect
            v-model="form.type"
            :options="[
              { label: '小票打印机', value: 1 },
              { label: '标签打印机', value: 2 },
            ]"
          />
        </BaseFormItem>
        <BaseFormItem v-if="isEdit" label="启用状态">
          <BaseSwitch v-model="form.status" :active-value="1" :inactive-value="2" />
        </BaseFormItem>
        <BaseFormItem label="默认打印机">
          <BaseSwitch v-model="form.is_default" :active-value="1" :inactive-value="0" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="form.remark" :rows="3" placeholder="可填写摆放位置、用途等" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submit">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="testDlg" title="测试打印" max-width="min(560px, 96vw)">
      <div class="space-y-4">
        <p class="m-0 text-sm text-slate-600">
          当前打印机：<span class="font-medium text-slate-900">{{ testingRow?.name || '-' }}</span>
        </p>
        <BaseFormItem label="打印份数">
          <a-input-number v-model="testForm.copies" :min="1" :max="5" class="w-full" />
        </BaseFormItem>
        <BaseFormItem label="测试内容">
          <BaseTextarea v-model="testForm.content" :rows="5" placeholder="留空则打印系统默认测试小票" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="testDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="testing" @click="submitTest">开始打印</BaseButton>
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
  BasePagination,
  BaseSelect,
  BaseSwitch,
  BaseTable,
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn, TableRowAction } from '@/components/base/types'
import {
  batchQueryPrinterStatus,
  bindPrinter,
  listAllPrinters,
  listStorePrinters,
  testPrint,
  unbindPrinter,
  updatePrinter,
  type PrinterRow,
} from '@/api/printer'
import { listAllStores } from '@/api/store'
import type { Store } from '@/api/types'
import { confirmDialog } from '@/feedback/confirm'
import { toast } from '@/feedback/toast'
import { useUserStore } from '@/store/user'

const qc = useQueryClient()
const userStore = useUserStore()

const selectedStoreId = ref<number | undefined>(Number(userStore.tenantId || userStore.userInfo?.store_id || 0) || undefined)
const keyword = ref('')
const statusFilter = ref<number | ''>('')
const page = ref(1)
const pageSize = ref(10)
const syncing = ref(false)

const stores = ref<Store[]>([])
void listAllStores()
  .then((rows) => {
    stores.value = rows
    if (!selectedStoreId.value && rows.length === 1) selectedStoreId.value = rows[0].id
  })
  .catch(() => {
    stores.value = []
  })

const storeOptions = computed<BaseSelectOption[]>(() => {
  const opts = stores.value.map((s) => ({ label: s.name, value: s.id }))
  if (opts.length) return opts
  if (selectedStoreId.value) return [{ label: `门店 #${selectedStoreId.value}`, value: selectedStoreId.value }]
  return []
})

const currentStoreId = computed(() => Number(selectedStoreId.value || 0))
const queryKey = computed(() => ['printers', currentStoreId.value] as const)

const { data: rowsData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: async () => {
    if (currentStoreId.value > 0) return listStorePrinters(currentStoreId.value)
    const pageData = await listAllPrinters()
    return pageData.list
  },
})

const onlineOverrides = ref<Record<string, number>>({})
const rows = computed<PrinterRow[]>(() =>
  (rowsData.value ?? []).map((row) => ({
    ...row,
    online: onlineOverrides.value[row.sn] ?? row.online,
  })),
)

const statusOptions: BaseSelectOption[] = [
  { label: '全部状态', value: '' },
  { label: '正常', value: 1 },
  { label: '停用', value: 2 },
]

const columns: BaseTableColumn[] = [
  { key: 'store_id', label: '门店', minWidth: '130px', ellipsis: true },
  { key: 'name', label: '名称', prop: 'name', minWidth: '150px', ellipsis: true },
  { key: 'sn', label: 'SN', prop: 'sn', minWidth: '150px', ellipsis: true },
  { key: 'type', label: '类型', width: '112px' },
  { key: 'status', label: '启用', width: '86px', align: 'center' },
  { key: 'is_default', label: '默认', width: '86px', align: 'center' },
  { key: 'online', label: '在线', width: '96px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '160px', ellipsis: true },
  { key: 'created_at', label: '创建时间', width: '150px' },
  { key: 'actions', label: '操作', width: '190px', align: 'right', fixed: 'right' },
]

const filteredRows = computed(() =>
  rows.value.filter((row) => {
    const kw = keyword.value.trim().toLowerCase()
    const hit =
      !kw ||
      [row.sn, row.name, row.remark ?? '', storeName(row.store_id)]
        .join(' ')
        .toLowerCase()
        .includes(kw)
    const statusHit = statusFilter.value === '' || row.status === statusFilter.value
    return hit && statusHit
  }),
)

const pagedRows = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredRows.value.slice(start, start + pageSize.value)
})

const onlineCount = computed(() => filteredRows.value.filter((row) => row.online === 1).length)
const defaultCount = computed(() => filteredRows.value.filter((row) => row.is_default === 1).length)
const disabledCount = computed(() => filteredRows.value.filter((row) => row.status === 2).length)

watch([keyword, statusFilter, selectedStoreId, pageSize], () => {
  page.value = 1
})

watch(selectedStoreId, () => {
  onlineOverrides.value = {}
})

function reload(): void {
  void qc.invalidateQueries({ queryKey: ['printers'] })
}

function storeName(id: number): string {
  return stores.value.find((s) => s.id === id)?.name || `门店 #${id}`
}

function printerTypeLabel(v: number): string {
  return v === 2 ? '标签打印机' : '小票打印机'
}

function onlineLabel(v: number): string {
  if (v === 1) return '在线'
  if (v === 2) return '异常'
  return '离线'
}

function onlineColor(v: number): string {
  if (v === 1) return 'green'
  if (v === 2) return 'orange'
  return 'gray'
}

function formatDate(v?: string): string {
  if (!v) return '-'
  return v.replace('T', ' ').replace(/\.\d+.*$/, '').replace(/\+\d\d:\d\d$/, '')
}

const dlg = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const saving = ref(false)

const form = reactive({
  store_id: undefined as number | undefined,
  sn: '',
  name: '',
  type: 1,
  status: 1,
  is_default: 0,
  remark: '',
})

function resetForm(): void {
  form.store_id = currentStoreId.value || undefined
  form.sn = ''
  form.name = ''
  form.type = 1
  form.status = 1
  form.is_default = 0
  form.remark = ''
}

function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  resetForm()
  dlg.value = true
}

function openEdit(row: PrinterRow): void {
  isEdit.value = true
  editId.value = row.id
  form.store_id = row.store_id
  form.sn = row.sn
  form.name = row.name
  form.type = row.type || 1
  form.status = row.status || 1
  form.is_default = row.is_default || 0
  form.remark = row.remark || ''
  dlg.value = true
}

async function submit(): Promise<void> {
  if (!form.store_id) {
    toast.warning('请选择门店')
    return
  }
  if (!form.sn.trim()) {
    toast.warning('请填写打印机 SN')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updatePrinter(editId.value, {
        name: form.name.trim(),
        type: form.type,
        status: form.status,
        is_default: form.is_default,
        remark: form.remark.trim(),
      })
    } else {
      await bindPrinter({
        store_id: form.store_id,
        sn: form.sn.trim(),
        name: form.name.trim(),
        type: form.type,
        is_default: form.is_default,
        remark: form.remark.trim(),
      })
      selectedStoreId.value = form.store_id
    }
    toast.success('已保存')
    dlg.value = false
    reload()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function toggleStatus(row: PrinterRow, status: number): Promise<void> {
  try {
    await updatePrinter(row.id, { status })
    row.status = status
    toast.success('状态已更新')
    reload()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '状态更新失败')
  }
}

async function toggleDefault(row: PrinterRow, isDefault: number): Promise<void> {
  try {
    await updatePrinter(row.id, { is_default: isDefault })
    toast.success(isDefault === 1 ? '已设为默认打印机' : '已取消默认')
    reload()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '默认打印机更新失败')
  }
}

async function syncStatus(): Promise<void> {
  if (!currentStoreId.value) {
    toast.warning('请选择门店')
    return
  }
  syncing.value = true
  try {
    const list = await batchQueryPrinterStatus(currentStoreId.value)
    onlineOverrides.value = Object.fromEntries(list.map((x) => [x.sn, x.online ?? x.status ?? 0]))
    toast.success('在线状态已刷新')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '刷新状态失败')
  } finally {
    syncing.value = false
  }
}

async function onUnbind(row: PrinterRow): Promise<void> {
  const ok = await confirmDialog({ message: `确定解绑打印机「${row.name || row.sn}」？` })
  if (!ok) return
  try {
    await unbindPrinter(row.id)
    toast.success('已解绑')
    reload()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '解绑失败')
  }
}

const testDlg = ref(false)
const testing = ref(false)
const testingRow = ref<PrinterRow | null>(null)
const testForm = reactive({
  content: '',
  copies: 1,
})

function openTest(row: PrinterRow): void {
  testingRow.value = row
  testForm.content = ''
  testForm.copies = 1
  testDlg.value = true
}

async function submitTest(): Promise<void> {
  if (!testingRow.value) return
  testing.value = true
  try {
    const res = await testPrint(testingRow.value.id, {
      content: testForm.content.trim(),
      copies: testForm.copies,
    })
    toast.success(`测试打印已发送：${res.order_id}`)
    testDlg.value = false
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '测试打印失败')
  } finally {
    testing.value = false
  }
}

function rowActions(row: PrinterRow): TableRowAction[] {
  return [
    { label: '编辑', permission: 'printer:edit', onClick: () => openEdit(row) },
    { label: '测试打印', onClick: () => openTest(row), disabled: row.status !== 1 },
    {
      label: row.is_default === 1 ? '取消默认' : '设为默认',
      permission: 'printer:edit',
      place: 'more',
      onClick: () => void toggleDefault(row, row.is_default === 1 ? 0 : 1),
    },
    { label: '解绑', permission: 'printer:unbind', danger: true, place: 'more', onClick: () => void onUnbind(row) },
  ]
}
</script>

<style scoped>
.printer-stat {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  border: 1px solid var(--color-border-2);
  border-radius: 10px;
  background: var(--color-bg-1);
  padding: 14px 16px;
}

.printer-stat span {
  color: #64748b;
  font-size: 13px;
}

.printer-stat strong {
  color: #0f172a;
  font-size: 22px;
  line-height: 1;
}
</style>
