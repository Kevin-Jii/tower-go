<template>
  <div class="dict-page">
    <aside class="dict-type-panel">
      <div class="dict-panel-head">
        <div>
          <h2 class="dict-panel-title">字典类型</h2>
          <p class="dict-panel-sub">{{ filteredTypes.length }} / {{ types.length }}</p>
        </div>
        <BaseButton v-permission="'system:dict:type:add'" variant="primary" size="sm" @click="openType()">新增</BaseButton>
      </div>

      <BaseInput v-model="typeKeyword" placeholder="搜索编码 / 名称" clearable />

      <div class="dict-type-list">
        <div
          v-for="item in filteredTypes"
          :key="item.id"
          class="dict-type-item"
          :class="{ 'is-active': currentType?.id === item.id }"
          @click="pickType(item)"
        >
          <div class="min-w-0 flex-1">
            <div class="dict-type-main">
              <span class="truncate font-medium">{{ item.name }}</span>
              <span class="dict-status" :class="item.status === 1 ? 'is-on' : 'is-off'">{{ statusText(item.status) }}</span>
            </div>
            <div class="dict-type-code">{{ item.code }}</div>
          </div>
          <div class="dict-type-actions" @click.stop>
            <BaseButton v-permission="'system:dict:type:edit'" variant="link" size="sm" @click="openType(item)">编辑</BaseButton>
            <BaseButton v-permission="'system:dict:type:delete'" variant="link" size="sm" class="!text-rose-600" @click="delType(item)">删除</BaseButton>
          </div>
        </div>

        <div v-if="!typesLoading && filteredTypes.length === 0" class="dict-empty">暂无数据</div>
      </div>
    </aside>

    <section class="dict-data-panel">
      <div class="dict-panel-head">
        <div class="min-w-0">
          <h2 class="dict-panel-title truncate">{{ currentType?.name || '字典数据' }}</h2>
          <p class="dict-panel-sub truncate">
            <template v-if="currentType">{{ currentType.code }}{{ currentType.remark ? ` · ${currentType.remark}` : '' }}</template>
            <template v-else>请选择字典类型</template>
          </p>
        </div>
        <BaseButton v-permission="'system:dict:data:add'" variant="primary" size="sm" :disabled="!currentType" @click="openData()">新增数据</BaseButton>
      </div>

      <BaseTable
        v-if="currentType"
        :columns="dataColumns"
        :data="(dataList as unknown) as Record<string, unknown>[]"
        :loading="dataLoading"
        min-width="720px"
        height="calc(100vh - 260px)"
      >
        <template #cell-style="{ row }">
          <DictTag :type="currentType.code" :value="(row as DictData).value" />
        </template>
        <template #cell-status="{ row }">
          <span class="dict-status" :class="(row as DictData).status === 1 ? 'is-on' : 'is-off'">{{ statusText((row as DictData).status) }}</span>
        </template>
        <template #cell-actions="{ row }">
          <BaseTableRowActions :actions="dictDataRowActions(row as DictData)" :max-inline="2" />
        </template>
      </BaseTable>

      <div v-else class="dict-empty dict-empty-large">请选择字典类型</div>
    </section>

    <BaseDialog v-model="typeDlg" :title="typeEditId ? '编辑类型' : '新增类型'" max-width="min(520px, 96vw)">
      <div class="dict-form-grid">
        <BaseFormItem label="编码" required>
          <BaseInput v-model="typeForm.code" :disabled="!!typeEditId" />
        </BaseFormItem>
        <BaseFormItem label="名称" required>
          <BaseInput v-model="typeForm.name" />
        </BaseFormItem>
        <BaseFormItem label="备注" class="sm:col-span-2">
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

    <BaseDialog v-model="dataDlg" :title="dataEditId ? '编辑数据' : '新增数据'" max-width="min(620px, 96vw)">
      <div class="dict-form-grid">
        <BaseFormItem label="标签" required>
          <BaseInput v-model="dataForm.label" />
        </BaseFormItem>
        <BaseFormItem label="值" required>
          <BaseInput v-model="dataForm.value" />
        </BaseFormItem>
        <BaseFormItem label="排序">
          <BaseNumberInput v-model="dataForm.sort" />
        </BaseFormItem>
        <BaseFormItem label="样式">
          <BaseSelect v-model="dataForm.list_class" :options="listClassOptions" placeholder="success / info / ..." />
        </BaseFormItem>
        <BaseFormItem label="备注" class="sm:col-span-2">
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

const dataColumns: BaseTableColumn[] = [
  { key: 'label', label: '标签', prop: 'label', minWidth: '140px', ellipsis: true },
  { key: 'value', label: '值', prop: 'value', minWidth: '120px', ellipsis: true },
  { key: 'style', label: '样式', width: '100px' },
  { key: 'sort', label: '排序', prop: 'sort', width: '80px' },
  { key: 'status', label: '状态', width: '80px' },
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
const typeKeyword = ref('')
const filteredTypes = computed(() => {
  const kw = typeKeyword.value.trim().toLowerCase()
  if (!kw) return types.value
  return types.value.filter((item) => `${item.code} ${item.name}`.toLowerCase().includes(kw))
})

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

watch(
  types,
  (list) => {
    if (!list.length) {
      currentType.value = null
      return
    }
    if (!currentType.value || !list.some((item) => item.id === currentType.value?.id)) {
      currentType.value = list[0]
    }
  },
  { immediate: true },
)

function pickType(row: DictType): void {
  currentType.value = row
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

function dictDataRowActions(row: DictData): TableRowAction[] {
  return [
    { label: '编辑', permission: 'system:dict:data:edit', onClick: () => openData(row) },
    { label: '删除', permission: 'system:dict:data:delete', danger: true, onClick: () => void delData(row) },
  ]
}

function statusText(status: number): string {
  return status === 1 ? '启用' : '停用'
}
</script>

<style scoped>
.dict-page {
  display: grid;
  grid-template-columns: minmax(260px, 320px) minmax(0, 1fr);
  gap: 16px;
  height: calc(100vh - 150px);
  min-height: 520px;
  overflow: hidden;
}

.dict-type-panel,
.dict-data-panel {
  min-width: 0;
  min-height: 0;
  border: 1px solid var(--color-border-2);
  border-radius: var(--border-radius-large);
  background: var(--color-bg-2);
}

.dict-type-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 14px;
}

.dict-data-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 14px;
}

.dict-panel-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  min-height: 36px;
}

.dict-panel-title {
  margin: 0;
  color: #0f172a;
  font-size: 16px;
  font-weight: 650;
  line-height: 22px;
}

.dict-panel-sub {
  margin: 2px 0 0;
  color: #64748b;
  font-size: 12px;
  line-height: 18px;
}

.dict-type-list {
  display: flex;
  min-height: 0;
  flex: 1;
  flex-direction: column;
  gap: 6px;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 2px;
}

.dict-type-item {
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 58px;
  padding: 9px 10px;
  border: 1px solid transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.16s ease, border-color 0.16s ease;
}

.dict-type-item:hover {
  background: #f8fafc;
}

.dict-type-item.is-active {
  border-color: #93c5fd;
  background: #eff6ff;
}

.dict-type-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  min-width: 0;
}

.dict-type-code {
  margin-top: 3px;
  overflow: hidden;
  color: #64748b;
  font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dict-type-actions {
  display: flex;
  flex-shrink: 0;
  gap: 4px;
}

.dict-status {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  height: 22px;
  padding: 0 8px;
  border-radius: 999px;
  font-size: 12px;
  line-height: 22px;
}

.dict-status.is-on {
  color: #047857;
  background: #d1fae5;
}

.dict-status.is-off {
  color: #64748b;
  background: #f1f5f9;
}

.dict-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 120px;
  color: #94a3b8;
  font-size: 14px;
}

.dict-empty-large {
  flex: 1;
  min-height: 320px;
  border: 1px dashed var(--color-border-2);
  border-radius: var(--border-radius-large);
}

.dict-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

@media (max-width: 900px) {
  .dict-page {
    grid-template-columns: 1fr;
    height: auto;
    min-height: 0;
    overflow: visible;
  }

  .dict-type-panel {
    max-height: min(520px, 70vh);
  }

  .dict-form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
