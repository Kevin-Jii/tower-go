<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col sm:flex-row gap-3 sm:items-center justify-between">
      <h2 class="page-title">角色管理</h2>
      <BaseButton v-permission="'system:role:add'" variant="primary" @click="openCreate">新增角色</BaseButton>
    </div>
    <BaseTable :columns="columns" :data="(roles as unknown) as Record<string, unknown>[]" :loading="loading" min-width="960px">
      <template #cell-data_scope="{ row }">
        {{ scopeLabel(Number((row as Role).data_scope)) }}
      </template>
      <template #cell-status="{ row }">
        {{ (row as Role).status === 1 ? '启用' : '禁用' }}
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
          <BaseButton v-permission="'system:role:menu'" variant="link" size="sm" @click="openAssign(row as Role)">分配菜单</BaseButton>
          <BaseButton v-permission="'system:role:edit'" variant="link" size="sm" @click="openEdit(row as Role)">编辑</BaseButton>
          <BaseButton v-permission="'system:role:delete'" variant="link" size="sm" @click="onDelete(row as Role)">删除</BaseButton>
        </div>
      </template>
    </BaseTable>

    <BaseDialog v-model="dlg" :title="dlgTitle" :max-width="mode === 'assign' ? 'min(560px, 96vw)' : 'min(480px, 96vw)'">
      <div v-if="mode !== 'assign'" class="space-y-4">
        <BaseFormItem label="名称" required>
          <BaseInput v-model="form.name" />
        </BaseFormItem>
        <BaseFormItem label="编码" required>
          <BaseInput v-model="form.code" :disabled="mode === 'edit'" />
        </BaseFormItem>
        <BaseFormItem label="描述">
          <BaseTextarea v-model="form.description" :rows="2" />
        </BaseFormItem>
        <BaseFormItem label="数据范围">
          <BaseSelect
            v-model="form.data_scope"
            :options="[
              { label: '全部', value: 1 },
              { label: '租户', value: 2 },
              { label: '本门店', value: 3 },
              { label: '仅本人', value: 4 },
            ]"
          />
        </BaseFormItem>
        <BaseFormItem label="状态">
          <BaseSwitch v-model="form.status" :active-value="1" :inactive-value="0" />
        </BaseFormItem>
      </div>
      <div v-else class="max-h-[60vh] overflow-y-auto pr-1 rounded-xl border border-slate-100 bg-slate-50/40 p-3">
        <BaseTreeCheck v-model="checkedMenuIds" :nodes="menuTreeNodes" />
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submit">确定</BaseButton>
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
  BaseTextarea,
  BaseTreeCheck,
} from '@/components/base'
import type { BaseTableColumn, BaseTreeNode } from '@/components/base/types'
import { assignRoleMenus, fetchMenuTree, fetchRoleMenuIds, fetchRoleMenuPermissions } from '@/api/menu'
import { createRole, deleteRole, listRoles, updateRole } from '@/api/role'
import type { Menu, Role } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()
const { data: rolesData, isLoading: loading } = useQuery({
  queryKey: ['roles'],
  queryFn: () => listRoles(),
})
const roles = computed(() => rolesData.value ?? [])

const columns: BaseTableColumn[] = [
  { key: 'name', label: '名称', prop: 'name', minWidth: '120px' },
  { key: 'code', label: '编码', prop: 'code', width: '140px' },
  { key: 'description', label: '描述', prop: 'description', minWidth: '160px', ellipsis: true },
  { key: 'data_scope', label: '数据范围', width: '100px' },
  { key: 'status', label: '状态', width: '88px' },
  { key: 'actions', label: '操作', width: '280px', align: 'right' },
]

const { data: rawMenuTree } = useQuery({
  queryKey: ['menus', 'tree'],
  queryFn: fetchMenuTree,
})

const dlg = ref(false)
const saving = ref(false)
const mode = ref<'create' | 'edit' | 'assign'>('create')
const editId = ref(0)
const assignRoleId = ref(0)
const checkedMenuIds = ref<number[]>([])
/** 打开「分配菜单」时从服务端拉取的各菜单权限位，提交时保留未改动项的权限 */
const roleMenuPermSnapshot = ref<Record<number, number>>({})

const form = reactive({
  name: '',
  code: '',
  description: '',
  data_scope: 3 as number | undefined,
  status: 1 as number,
})

const dlgTitle = computed(() => {
  if (mode.value === 'assign') return '分配菜单'
  return mode.value === 'edit' ? '编辑角色' : '新增角色'
})

const menuTreeData = computed(() => filterTree(rawMenuTree.value ?? []))

const menuTreeNodes = computed(() => menuTreeData.value as unknown as BaseTreeNode[])

function filterTree(nodes: Menu[]): Menu[] {
  return (nodes ?? [])
    .filter((n) => n.status !== 0)
    .map((n) => ({
      ...n,
      children: n.children?.length ? filterTree(n.children) : [],
    }))
}

function scopeLabel(v?: number): string {
  const m: Record<number, string> = { 1: '全部', 2: '租户', 3: '本门店', 4: '本人' }
  return m[v ?? 3] ?? '-'
}

function openCreate(): void {
  mode.value = 'create'
  editId.value = 0
  form.name = ''
  form.code = ''
  form.description = ''
  form.data_scope = 3
  form.status = 1
  dlg.value = true
}

function openEdit(row: Role): void {
  mode.value = 'edit'
  editId.value = row.id
  form.name = row.name
  form.code = row.code
  form.description = row.description ?? ''
  form.data_scope = row.data_scope ?? 3
  form.status = row.status
  dlg.value = true
}

async function openAssign(row: Role): Promise<void> {
  mode.value = 'assign'
  assignRoleId.value = row.id
  roleMenuPermSnapshot.value = {}
  try {
    const [ids, permMap] = await Promise.all([
      fetchRoleMenuIds(row.id),
      fetchRoleMenuPermissions(row.id).catch(() => ({} as Record<number, number>)),
    ])
    checkedMenuIds.value = ids.map((x) => Number(x))
    roleMenuPermSnapshot.value = permMap
  } catch {
    checkedMenuIds.value = []
  }
  dlg.value = true
}

async function submit(): Promise<void> {
  saving.value = true
  try {
    if (mode.value === 'assign') {
      const keys = checkedMenuIds.value.map((x) => Number(x))
      const perms: Record<number, number> = {}
      for (const id of keys) {
        perms[id] = roleMenuPermSnapshot.value[id] ?? 15
      }
      await assignRoleMenus({ role_id: assignRoleId.value, menu_ids: keys, perms })
      toast.success('已保存菜单权限')
      dlg.value = false
      await qc.invalidateQueries({ queryKey: ['menus'] })
      await qc.invalidateQueries({ queryKey: ['menus', 'tree'] })
      await qc.invalidateQueries({ queryKey: ['roles'] })
      return
    }
    if (!form.name || !form.code) {
      toast.warning('请填写名称与编码')
      return
    }
    if (mode.value === 'create') {
      await createRole({
        name: form.name,
        code: form.code,
        description: form.description,
        data_scope: form.data_scope ?? 3,
        status: form.status,
      })
    } else {
      await updateRole(editId.value, {
        name: form.name,
        description: form.description,
        data_scope: form.data_scope ?? 3,
        status: form.status,
      })
    }
    toast.success('已保存')
    dlg.value = false
    await qc.invalidateQueries({ queryKey: ['roles'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '操作失败')
  } finally {
    saving.value = false
  }
}

async function onDelete(row: Role): Promise<void> {
  const ok = await confirmDialog({ message: `删除角色「${row.name}」？` })
  if (!ok) return
  try {
    await deleteRole(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['roles'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>
