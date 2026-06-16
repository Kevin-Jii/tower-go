<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-center gap-3 justify-between">
      <h2 class="page-title">菜单管理</h2>
      <BaseButton v-permission="'system:menu:add'" variant="primary" @click="openCreate()">新增</BaseButton>
    </div>
    <BaseTable :columns="columns" :data="(tree as unknown) as Record<string, unknown>[]" :loading="loading"
      min-width="960px" row-key="id" tree-children-key="children" :tree-default-expand-all="true">
      <template #cell-type="{ row }">
        {{ typeLabel(Number((row as Menu).type)) }}
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="menuRowActions(row as Menu)" />
      </template>
    </BaseTable>

    <BaseDialog v-model="visible" :title="isEdit ? '编辑菜单' : '新增菜单'" max-width="min(520px, 96vw)">
      <div class="max-h-[70vh] overflow-y-auto space-y-4 pr-1">
        <BaseFormItem label="父级">
          <a-tree-select
            v-model="form.parent_id"
            :data="parentTreeOptions"
            :field-names="{ key: 'id', title: 'title', children: 'children' }"
            placeholder="请选择父级目录或菜单"
            allow-clear
            allow-search
            class="w-full"
            @clear="form.parent_id = 0"
          />
        </BaseFormItem>
        <BaseFormItem label="名称 name" required>
          <BaseInput v-model="form.name" />
        </BaseFormItem>
        <BaseFormItem label="标题 title" required>
          <BaseInput v-model="form.title" />
        </BaseFormItem>
        <BaseFormItem label="图标">
          <a-select
            v-model="form.icon"
            placeholder="请选择图标"
            allow-clear
            allow-search
            :filter-option="filterIconOption"
            class="w-full"
          >
            <a-option v-for="item in iconOptions" :key="item.value" :value="item.value" :label="item.label">
              <span class="inline-flex items-center gap-2">
                <component :is="item.component" class="text-[16px] text-slate-500" />
                <span>{{ item.label }}</span>
              </span>
            </a-option>
          </a-select>
        </BaseFormItem>
        <BaseFormItem label="路由 path">
          <BaseInput v-model="form.path" placeholder="/system/user" />
        </BaseFormItem>
        <BaseFormItem label="组件">
          <a-cascader
            v-model="componentPath"
            :options="viewPathOptions"
            placeholder="请选择 src/views 下的文件"
            allow-clear
            allow-search
            :path-mode="true"
            class="w-full"
            @change="onComponentPathChange"
          />
        </BaseFormItem>
        <BaseFormItem label="权限码">
          <BaseInput v-model="form.permission" />
        </BaseFormItem>
        <BaseFormItem label="类型">
          <BaseSelect v-model="form.type" :options="[
            { label: '目录', value: 1 },
            { label: '菜单', value: 2 },
            { label: '按钮', value: 3 },
          ]" />
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
import type { Component } from 'vue'
import { computed, reactive, ref } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import * as ArcoIcons from '@arco-design/web-vue/es/icon'
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
import type { BaseTableColumn, TableRowAction } from '@/components/base/types'
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
const viewModules = import.meta.glob('../../**/*.vue')

interface CascaderOption {
  label: string
  value: string
  children?: CascaderOption[]
}

interface IconOption {
  label: string
  value: string
  component: Component
}

const columns: BaseTableColumn[] = [
  { key: 'title', label: '名称', prop: 'title', minWidth: '160px' },
  { key: 'name', label: '标识', prop: 'name', width: '120px' },
  { key: 'path', label: '路由', prop: 'path', minWidth: '140px' },
  { key: 'component', label: '组件路径', prop: 'component', minWidth: '160px', ellipsis: true },
  { key: 'permission', label: '权限码', prop: 'permission', minWidth: '140px', ellipsis: true },
  { key: 'type', label: '类型', width: '80px' },
  { key: 'actions', label: '操作', width: '140px', align: 'right' },
]

const visible = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const componentPath = ref<string[]>([])
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

const viewPathOptions = computed<CascaderOption[]>(() => {
  const root: CascaderOption[] = []
  for (const rawKey of Object.keys(viewModules).sort()) {
    const path = rawKey.replace(/^\.\.\/\.\.\//, '').replace(/\.vue$/, '')
    insertPath(root, path.split('/'))
  }
  return root
})

const iconComponents = ArcoIcons as Record<string, Component>
const iconOptions = computed<IconOption[]>(() =>
  Object.entries(iconComponents)
    .filter(([name, cmp]) => name.startsWith('Icon') && Boolean(cmp))
    .map(([name, component]) => ({
      label: iconExportToValue(name),
      value: iconExportToValue(name),
      component,
    }))
    .sort((a, b) => a.label.localeCompare(b.label)),
)

interface ParentTreeOption {
  id: number
  title: string
  children?: ParentTreeOption[]
}

function buildParentTreeOptions(menus: Menu[], editingID = 0): ParentTreeOption[] {
  return (menus ?? [])
    .filter((m) => m.type !== 3 && m.id !== editingID)
    .map((m) => {
      const children = buildParentTreeOptions(m.children ?? [], editingID)
      return {
        id: m.id,
        title: `${m.title}（${typeLabel(m.type)}）`,
        children: children.length ? children : undefined,
      }
    })
}

const parentTreeOptions = computed<ParentTreeOption[]>(() => [
  { id: 0, title: '顶级' },
  ...buildParentTreeOptions(tree.value, isEdit.value ? editId.value : 0),
])

function menuRowActions(row: Menu): TableRowAction[] {
  const actions: TableRowAction[] = []
  if (row.type === 1 || row.type === 2) {
    actions.push({ label: '新增下级', permission: 'system:menu:add', onClick: () => openCreate(row) })
  }
  actions.push({ label: '编辑', permission: 'system:menu:edit', onClick: () => openEdit(row) })
  actions.push({ label: '删除', permission: 'system:menu:delete', danger: true, onClick: () => void onDelete(row) })
  return actions
}

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
  componentPath.value = []
  form.permission = ''
  form.type = 2
  form.sort = 0
  form.visible = 1
  form.status = 1
  form.remark = ''
}

function openCreate(parent?: Menu): void {
  isEdit.value = false
  resetForm()
  if (parent) {
    form.parent_id = parent.id
    form.type = parent.type === 1 ? 2 : 3
  }
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
  componentPath.value = form.component ? form.component.split('/') : []
  form.permission = row.permission ?? ''
  form.type = row.type
  form.sort = row.sort ?? 0
  form.visible = row.visible ?? 1
  form.status = row.status ?? 1
  form.remark = row.remark ?? ''
  visible.value = true
}

function insertPath(options: CascaderOption[], segments: string[]): void {
  if (!segments.length) return
  const [head, ...tail] = segments
  let node = options.find((x) => x.value === head)
  if (!node) {
    node = { label: head, value: head }
    options.push(node)
  }
  if (tail.length) {
    node.children = node.children ?? []
    insertPath(node.children, tail)
  }
}

function onComponentPathChange(value: unknown): void {
  if (Array.isArray(value)) {
    form.component = value.map(String).join('/')
    return
  }
  form.component = value ? String(value) : ''
}

function iconExportToValue(exportName: string): string {
  const raw = exportName.replace(/^Icon/, '')
  return raw
    .replace(/([a-z0-9])([A-Z])/g, '$1-$2')
    .replace(/([A-Z])([A-Z][a-z])/g, '$1-$2')
    .toLowerCase()
}

function filterIconOption(input: string, option: unknown): boolean {
  const label = String((option as { label?: string })?.label ?? '')
  return label.toLowerCase().includes(input.toLowerCase())
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
