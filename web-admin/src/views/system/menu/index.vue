<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-center gap-3 justify-between">
      <h2 class="page-title">菜单管理</h2>
      <BaseButton v-permission="'system:menu:add'" variant="primary" @click="openCreate()">新增</BaseButton>
    </div>
    <BaseTable
      :columns="columns"
      :data="(tree as unknown) as Record<string, unknown>[]"
      :loading="loading"
      min-width="960px"
      row-key="id"
      tree-children-key="children"
      :tree-default-expand-all="true"
    >
      <template #cell-type="{ row }">
        {{ typeLabel(Number((row as Menu).type)) }}
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
          <BaseButton v-permission="'system:menu:edit'" variant="link" size="sm" @click="openEdit(row as Menu)">编辑</BaseButton>
          <BaseButton v-permission="'system:menu:delete'" variant="link" size="sm" @click="onDelete(row as Menu)">删除</BaseButton>
        </div>
      </template>
    </BaseTable>

    <BaseDialog v-model="visible" :title="isEdit ? '编辑菜单' : '新增菜单'" max-width="min(520px, 96vw)">
      <div class="max-h-[70vh] overflow-y-auto space-y-4 pr-1">
        <BaseFormItem label="父级">
          <BaseSelect v-model="form.parent_id" :options="parentOptions" placeholder="0 为顶级" />
        </BaseFormItem>
        <BaseFormItem label="名称 name" required>
          <BaseInput v-model="form.name" />
        </BaseFormItem>
        <BaseFormItem label="标题 title" required>
          <BaseInput v-model="form.title" />
        </BaseFormItem>
        <BaseFormItem label="图标">
          <BaseInput v-model="form.icon" placeholder="如 setting / user" />
        </BaseFormItem>
        <BaseFormItem label="路由 path">
          <BaseInput v-model="form.path" placeholder="/system/user" />
        </BaseFormItem>
        <BaseFormItem label="组件">
          <BaseInput v-model="form.component" placeholder="system/user/index" />
        </BaseFormItem>
        <BaseFormItem label="权限码">
          <BaseInput v-model="form.permission" />
        </BaseFormItem>
        <BaseFormItem label="类型">
          <BaseSelect
            v-model="form.type"
            :options="[
              { label: '目录', value: 1 },
              { label: '菜单', value: 2 },
              { label: '按钮', value: 3 },
            ]"
          />
        </BaseFormItem>
        <BaseFormItem label="排序">
          <BaseNumberInput v-model="form.sort" />
        </BaseFormItem>
        <BaseFormItem label="可见">
          <BaseSwitch v-model="form.visible" :active-value="1" :inactive-value="0" />
        </BaseFormItem>
        <BaseFormItem label="状态">
          <BaseSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="form.remark" :rows="2" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="visible = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="save">保存</BaseButton>
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
  BaseNumberInput,
  BaseSelect,
  BaseSwitch,
  BaseTable,
  BaseTextarea,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn } from '@/components/base/types'
import { createMenu, deleteMenu, fetchMenuTree, updateMenu } from '@/api/menu'
import type { Menu } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()
const { data: treeData, isLoading: loading } = useQuery({
  queryKey: ['menus', 'tree'],
  queryFn: fetchMenuTree,
})
const tree = computed(() => treeData.value ?? [])

const columns: BaseTableColumn[] = [
  { key: 'title', label: '名称', prop: 'title', minWidth: '160px' },
  { key: 'name', label: '标识', prop: 'name', width: '120px' },
  { key: 'path', label: '路由', prop: 'path', minWidth: '140px' },
  { key: 'component', label: '组件路径', prop: 'component', minWidth: '160px', ellipsis: true },
  { key: 'permission', label: '权限码', prop: 'permission', minWidth: '140px', ellipsis: true },
  { key: 'type', label: '类型', width: '80px' },
  { key: 'actions', label: '操作', width: '200px', align: 'right' },
]

const visible = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const form = reactive({
  parent_id: 0 as number | undefined,
  name: '',
  title: '',
  icon: '',
  path: '',
  component: '',
  permission: '',
  type: 2,
  sort: 0 as number | undefined,
  visible: 1,
  status: 1,
  remark: '',
})

function flatten(menus: Menu[], acc: { id: number; label: string }[] = [], prefix = ''): { id: number; label: string }[] {
  for (const m of menus) {
    acc.push({ id: m.id, label: `${prefix}${m.title} (#${m.id})` })
    if (m.children?.length) flatten(m.children, acc, `${prefix}  `)
  }
  return acc
}

const parentOptions = computed<BaseSelectOption[]>(() => {
  const flat = flatten(tree.value)
  return [{ label: '顶级', value: 0 }, ...flat.map((x) => ({ label: x.label, value: x.id }))]
})

function typeLabel(t: number): string {
  if (t === 1) return '目录'
  if (t === 2) return '菜单'
  if (t === 3) return '按钮'
  return String(t)
}

function resetForm(): void {
  form.parent_id = 0
  form.name = ''
  form.title = ''
  form.icon = ''
  form.path = ''
  form.component = ''
  form.permission = ''
  form.type = 2
  form.sort = 0
  form.visible = 1
  form.status = 1
  form.remark = ''
}

function openCreate(): void {
  isEdit.value = false
  resetForm()
  visible.value = true
}

function openEdit(row: Menu): void {
  isEdit.value = true
  editId.value = row.id
  form.parent_id = row.parent_id ?? 0
  form.name = row.name
  form.title = row.title
  form.icon = row.icon ?? ''
  form.path = row.path ?? ''
  form.component = row.component ?? ''
  form.permission = row.permission ?? ''
  form.type = row.type
  form.sort = row.sort ?? 0
  form.visible = row.visible ?? 1
  form.status = row.status ?? 1
  form.remark = row.remark ?? ''
  visible.value = true
}

async function save(): Promise<void> {
  if (!form.name || !form.title) {
    toast.warning('请填写名称与标题')
    return
  }
  saving.value = true
  try {
    const body = { ...form, parent_id: form.parent_id ?? 0, sort: form.sort ?? 0 }
    if (isEdit.value) await updateMenu(editId.value, body)
    else await createMenu(body)
    toast.success('已保存')
    visible.value = false
    await qc.invalidateQueries({ queryKey: ['menus'] })
    await qc.invalidateQueries({ queryKey: ['menus', 'tree'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function onDelete(row: Menu): Promise<void> {
  const ok = await confirmDialog({ message: `确定删除「${row.title}」？` })
  if (!ok) return
  try {
    await deleteMenu(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['menus'] })
    await qc.invalidateQueries({ queryKey: ['menus', 'tree'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>
