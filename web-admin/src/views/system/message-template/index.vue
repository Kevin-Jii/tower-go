<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">消息模板</h2>
      <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-56" placeholder="编码/名称关键词" clearable @enter="reload" />
        <BaseSelect
          v-model="enabledFilter"
          class="w-full sm:w-36"
          :options="[
            { label: '全部状态', value: '' },
            { label: '启用', value: 1 },
            { label: '禁用', value: 0 },
          ]"
        />
        <BaseButton variant="primary" @click="reload">查询</BaseButton>
        <BaseButton v-permission="'message:template:add'" variant="primary" @click="openCreate">新增模板</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(filteredList as unknown) as Record<string, unknown>[]" :loading="loading" min-width="1120px">
      <template #cell-is_enabled="{ row }">
        <BaseSwitch
          v-permission="'message:template:edit'"
          :model-value="(row as MessageTemplate).is_enabled"
          :active-value="true"
          :inactive-value="false"
          @update:model-value="toggleEnabled(row as MessageTemplate, $event)"
        />
      </template>
      <template #cell-content="{ row }">
        <div class="max-w-[460px] truncate" :title="(row as MessageTemplate).content">{{ (row as MessageTemplate).content }}</div>
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="messageTemplateRowActions(row as MessageTemplate)" />
      </template>
    </BaseTable>

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑消息模板' : '新增消息模板'" max-width="min(760px, 96vw)">
      <div class="space-y-4 max-h-[72vh] overflow-y-auto pr-1">
        <BaseFormItem label="模板编码" required>
          <BaseInput v-model="form.code" :disabled="isEdit" placeholder="如 purchase_created" />
        </BaseFormItem>
        <BaseFormItem label="模板名称" required>
          <BaseInput v-model="form.name" />
        </BaseFormItem>
        <BaseFormItem label="消息标题">
          <BaseInput v-model="form.title" placeholder="支持变量，如 {{store_name}}" />
        </BaseFormItem>
        <BaseFormItem label="消息内容" required>
          <BaseTextarea v-model="form.content" :rows="6" placeholder="模板正文，可包含变量占位符" />
        </BaseFormItem>
        <BaseFormItem label="模板说明">
          <BaseTextarea v-model="form.description" :rows="2" />
        </BaseFormItem>
        <BaseFormItem label="可用变量(JSON)">
          <BaseTextarea v-model="form.variables" :rows="3" placeholder='如 ["store_name","order_no"]' />
        </BaseFormItem>
        <BaseFormItem label="是否启用">
          <BaseSwitch v-model="form.is_enabled" :active-value="true" :inactive-value="false" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submit">保存</BaseButton>
      </template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseSelect,
  BaseSwitch,
  BaseTable,
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn, TableRowAction } from '@/components/base/types'
import {
  createMessageTemplate,
  deleteMessageTemplate,
  listMessageTemplates,
  updateMessageTemplate,
} from '@/api/messageTemplate'
import type { MessageTemplate } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()

const { data: rowsData, isLoading: loading } = useQuery({
  queryKey: ['message-templates'],
  queryFn: listMessageTemplates,
})

const list = computed(() => rowsData.value ?? [])

const columns: BaseTableColumn[] = [
  { key: 'code', label: '编码', prop: 'code', minWidth: '170px', ellipsis: true },
  { key: 'name', label: '名称', prop: 'name', minWidth: '140px', ellipsis: true },
  { key: 'title', label: '标题', prop: 'title', minWidth: '180px', ellipsis: true },
  { key: 'content', label: '内容', minWidth: '320px' },
  { key: 'is_enabled', label: '启用', width: '90px' },
  { key: 'actions', label: '操作', width: '120px', align: 'right' },
]

const keyword = ref('')
const enabledFilter = ref<number | ''>('')

const filteredList = computed(() =>
  list.value.filter((x) => {
    const kw = keyword.value.trim().toLowerCase()
    const hit = !kw || x.code.toLowerCase().includes(kw) || x.name.toLowerCase().includes(kw)
    const enabled = enabledFilter.value === '' ? true : Number(x.is_enabled ? 1 : 0) === enabledFilter.value
    return hit && enabled
  }),
)

function reload(): void {
  void qc.invalidateQueries({ queryKey: ['message-templates'] })
}

const dlg = ref(false)
const isEdit = ref(false)
const editId = ref<number>(0)
const saving = ref(false)

const form = reactive({
  code: '',
  name: '',
  title: '',
  content: '',
  description: '',
  variables: '',
  is_enabled: true,
})

function resetForm(): void {
  form.code = ''
  form.name = ''
  form.title = ''
  form.content = ''
  form.description = ''
  form.variables = ''
  form.is_enabled = true
}

function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  resetForm()
  dlg.value = true
}

function openEdit(row: MessageTemplate): void {
  isEdit.value = true
  editId.value = row.id
  form.code = row.code
  form.name = row.name
  form.title = row.title || ''
  form.content = row.content || ''
  form.description = row.description || ''
  form.variables = row.variables || ''
  form.is_enabled = !!row.is_enabled
  dlg.value = true
}

async function submit(): Promise<void> {
  if (!form.code.trim() || !form.name.trim() || !form.content.trim()) {
    toast.warning('请填写编码、名称、内容')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateMessageTemplate(editId.value, {
        name: form.name.trim(),
        title: form.title.trim(),
        content: form.content.trim(),
        description: form.description.trim(),
        variables: form.variables.trim(),
        is_enabled: form.is_enabled,
      })
    } else {
      await createMessageTemplate({
        code: form.code.trim(),
        name: form.name.trim(),
        title: form.title.trim(),
        content: form.content.trim(),
        description: form.description.trim(),
        variables: form.variables.trim(),
        is_enabled: form.is_enabled,
      })
    }
    toast.success('已保存')
    dlg.value = false
    await qc.invalidateQueries({ queryKey: ['message-templates'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function toggleEnabled(row: MessageTemplate, value: number | boolean): Promise<void> {
  try {
    await updateMessageTemplate(row.id, { is_enabled: Boolean(value) })
    row.is_enabled = Boolean(value)
    toast.success('状态已更新')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '状态更新失败')
  }
}

async function onDelete(row: MessageTemplate): Promise<void> {
  const ok = await confirmDialog({ message: `确定删除模板「${row.name}」？` })
  if (!ok) return
  try {
    await deleteMessageTemplate(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['message-templates'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

function messageTemplateRowActions(row: MessageTemplate): TableRowAction[] {
  return [
    { label: '编辑', permission: 'message:template:edit', onClick: () => openEdit(row) },
    { label: '删除', permission: 'message:template:delete', danger: true, onClick: () => void onDelete(row) },
  ]
}
</script>
