<template>
  <!-- xl 以下单列占满宽，避免双列时表格被压得过窄；双列用 minmax 允许列在网格内正确收缩而不挤爆表格 -->
  <div class="grid grid-cols-1 gap-4 xl:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
    <BaseCard class="min-w-0">
      <template #header>
        <div class="flex items-center justify-between gap-2 flex-wrap w-full">
          <span class="font-semibold text-slate-800">字典类型</span>
          <BaseButton v-permission="'system:dict:type:add'" variant="primary" size="sm" @click="openType()">新增
          </BaseButton>
        </div>
      </template>
      <div class="min-w-0 overflow-x-auto">
        <BaseTable :columns="typeColumns" :data="(types as unknown) as Record<string, unknown>[]"
          :loading="typesLoading" min-width="460px" height="360px" row-key="id"
          :highlight-row-key="currentType?.id ?? null" row-clickable @row-click="onPickTypeRow">
          <template #cell-status="{ row }">
            {{ (row as DictType).status === 1 ? '启用' : '停用' }}
          </template>
          <template #cell-actions="{ row }">
            <BaseTableRowActions :actions="dictTypeRowActions(row as DictType)" />
          </template>
        </BaseTable>
      </div>
    </BaseCard>

    <BaseCard class="min-w-0">
      <template #header>
        <div class="flex items-center justify-between gap-2 flex-wrap w-full">
          <span class="font-semibold text-slate-800">字典数据 {{ currentType ? `(${currentType.code})` : '' }}</span>
          <BaseButton v-permission="'system:dict:data:add'" variant="primary" size="sm" :disabled="!currentType"
            @click="openData()">
            新增
          </BaseButton>
        </div>
      </template>
      <div class=" overflow-x-auto">
        <BaseTable :columns="dataColumns" :data="(dataList as unknown) as Record<string, unknown>[]"
          :loading="dataLoading" height="360px">
          <template #cell-style="{ row }">
            <DictTag v-if="currentType" :type="currentType.code" :value="(row as DictData).value" />
          </template>
          <template #cell-actions="{ row }">
            <BaseTableRowActions :actions="dictDataRowActions(row as DictData)" />
          </template>
        </BaseTable>
      </div>
    </BaseCard>

    <BaseDialog v-model="typeDlg" :title="typeEditId ? '编辑类型' : '新增类型'" max-width="min(440px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="编码" required>
          <BaseInput v-model="typeForm.code" :disabled="!!typeEditId" />
        </BaseFormItem>
        <BaseFormItem label="名称" required>
          <BaseInput v-model="typeForm.name" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="typeForm.remark" :rows="2" />
        </BaseFormItem>
        <BaseFormItem label="状态">
          <BaseSwitch v-model="typeForm.status" :active-value="1" :inactive-value="0" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="typeDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="saveType">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="dataDlg" :title="dataEditId ? '编辑数据' : '新增数据'" max-width="min(480px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="标签" required>
          <BaseInput v-model="dataForm.label" />
        </BaseFormItem>
        <BaseFormItem label="值" required>
          <BaseInput v-model="dataForm.value" />
        </BaseFormItem>
        <BaseFormItem label="排序">
          <BaseNumberInput v-model="dataForm.sort" />
        </BaseFormItem>
        <BaseFormItem label="list_class">
          <BaseSelect v-model="dataForm.list_class" :options="listClassOptions" placeholder="success / info / …" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="dataForm.remark" :rows="2" />
        </BaseFormItem>
        <BaseFormItem label="状态">
          <BaseSwitch v-model="dataForm.status" :active-value="1" :inactive-value="0" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dataDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="saveData">保存</BaseButton>
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
  BaseSelect,
  BaseSwitch,
  BaseTable,
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn, TableRowAction } from '@/components/base/types'
import DictTag from '@/components/DictTag.vue'
import {
  createDictData,
  createDictType,
  deleteDictData,
  deleteDictType,
  listDictDataByTypeCode,
  listDictTypes,
  updateDictData,
  updateDictType,
} from '@/api/dict'
import type { DictData, DictType } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()

const typeColumns: BaseTableColumn[] = [
  { key: 'code', label: '编码', prop: 'code', minWidth: '120px', ellipsis: true },
  { key: 'name', label: '名称', prop: 'name', minWidth: '120px', ellipsis: true },
  { key: 'status', label: '状态', width: '72px' },
  { key: 'actions', label: '操作', width: '140px', align: 'right' },
]

const dataColumns: BaseTableColumn[] = [
  { key: 'label', label: '标签', prop: 'label', width: '100px', ellipsis: true },
  { key: 'value', label: '值', prop: 'value', width: '96px', ellipsis: true },
  { key: 'style', label: '样式', width: '100px' },
  { key: 'actions', label: '操作', width: '140px', align: 'right' },
]

const listClassOptions: BaseSelectOption[] = [
  { label: '（无）', value: '' },
  { label: 'success', value: 'success' },
  { label: 'info', value: 'info' },
  { label: 'warning', value: 'warning' },
  { label: 'danger', value: 'danger' },
]

const { data: typesData, isLoading: typesLoading } = useQuery({
  queryKey: ['dict-types'],
  queryFn: listDictTypes,
})

const types = computed(() => typesData.value ?? [])

const currentType = ref<DictType | null>(null)
const dataList = ref<DictData[]>([])
const dataLoading = ref(false)

async function loadData(): Promise<void> {
  if (!currentType.value) {
    dataList.value = []
    return
  }
  dataLoading.value = true
  try {
    dataList.value = await listDictDataByTypeCode(currentType.value.code)
  } catch {
    dataList.value = []
  } finally {
    dataLoading.value = false
  }
}

watch(
  currentType,
  () => {
    void loadData()
  },
  { immediate: true },
)

function onPickTypeRow(row: Record<string, unknown>): void {
  currentType.value = (row as unknown as DictType) ?? null
}

const typeDlg = ref(false)
const dataDlg = ref(false)
const saving = ref(false)
const typeEditId = ref(0)
const dataEditId = ref(0)

const typeForm = reactive({
  code: '',
  name: '',
  remark: '',
  status: 1,
})

const dataForm = reactive({
  label: '',
  value: '',
  sort: 0 as number | undefined,
  list_class: '' as string | undefined,
  remark: '',
  status: 1,
})

function openType(row?: DictType): void {
  typeEditId.value = row?.id ?? 0
  typeForm.code = row?.code ?? ''
  typeForm.name = row?.name ?? ''
  typeForm.remark = row?.remark ?? ''
  typeForm.status = row?.status ?? 1
  typeDlg.value = true
}

async function saveType(): Promise<void> {
  if (!typeForm.code || !typeForm.name) {
    toast.warning('请填写编码与名称')
    return
  }
  saving.value = true
  try {
    if (typeEditId.value) await updateDictType(typeEditId.value, { ...typeForm })
    else await createDictType({ ...typeForm })
    toast.success('已保存')
    typeDlg.value = false
    await qc.invalidateQueries({ queryKey: ['dict-types'] })
    await qc.invalidateQueries({ queryKey: ['dict', 'all'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

async function delType(row: DictType): Promise<void> {
  const ok = await confirmDialog({ message: `删除类型「${row.name}」？` })
  if (!ok) return
  try {
    await deleteDictType(row.id)
    toast.success('已删除')
    if (currentType.value?.id === row.id) currentType.value = null
    await qc.invalidateQueries({ queryKey: ['dict-types'] })
    await qc.invalidateQueries({ queryKey: ['dict', 'all'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

function openData(row?: DictData): void {
  if (!currentType.value) return
  dataEditId.value = row?.id ?? 0
  dataForm.label = row?.label ?? ''
  dataForm.value = row?.value ?? ''
  dataForm.sort = row?.sort ?? 0
  dataForm.list_class = row?.list_class ?? ''
  dataForm.remark = row?.remark ?? ''
  dataForm.status = row?.status ?? 1
  dataDlg.value = true
}

async function saveData(): Promise<void> {
  if (!currentType.value) return
  if (!dataForm.label || !dataForm.value) {
    toast.warning('请填写标签与值')
    return
  }
  saving.value = true
  try {
    const body = {
      type_code: currentType.value.code,
      label: dataForm.label,
      value: dataForm.value,
      sort: dataForm.sort ?? 0,
      list_class: dataForm.list_class ?? '',
      remark: dataForm.remark,
      status: dataForm.status,
    }
    if (dataEditId.value) await updateDictData(dataEditId.value, body)
    else await createDictData(body)
    toast.success('已保存')
    dataDlg.value = false
    await loadData()
    await qc.invalidateQueries({ queryKey: ['dict', 'all'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

async function delData(row: DictData): Promise<void> {
  const ok = await confirmDialog({ message: `删除「${row.label}」？` })
  if (!ok) return
  try {
    await deleteDictData(row.id)
    toast.success('已删除')
    await loadData()
    await qc.invalidateQueries({ queryKey: ['dict', 'all'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

function dictTypeRowActions(row: DictType): TableRowAction[] {
  return [
    { label: '编辑', permission: 'system:dict:type:edit', onClick: () => openType(row) },
    { label: '删除', permission: 'system:dict:type:delete', danger: true, onClick: () => void delType(row) },
  ]
}

function dictDataRowActions(row: DictData): TableRowAction[] {
  return [
    { label: '编辑', permission: 'system:dict:data:edit', onClick: () => openData(row) },
    { label: '删除', permission: 'system:dict:data:delete', danger: true, onClick: () => void delData(row) },
  ]
}
</script>
