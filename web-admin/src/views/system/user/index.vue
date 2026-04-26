<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row gap-3 md:items-end justify-between">
      <h2 class="page-title">用户管理</h2>
      <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-44" placeholder="关键字" clearable @enter="reload" />
        <div class="flex gap-2">
          <BaseButton variant="primary" @click="reload">查询</BaseButton>
          <BaseButton v-permission="'system:user:add'" variant="primary" @click="openCreate">新增</BaseButton>
        </div>
      </div>
    </div>
    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading" min-width="800px">
      <template #cell-role="{ row }">
        {{ (row as User).role?.name || '-' }}
      </template>
      <template #cell-store="{ row }">
        {{ (row as User).store?.name || (row as User).store_id || '-' }}
      </template>
      <template #cell-status="{ row }">
        {{ (row as User).status === 1 ? '正常' : '禁用' }}
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-wrap gap-1 justify-end" @click.stop>
          <BaseButton v-permission="'system:user:edit'" variant="link" size="sm" @click="openEdit(row as User)">编辑</BaseButton>
          <BaseButton v-permission="'system:user:edit'" variant="link" size="sm" @click="openRole(row as User)">角色</BaseButton>
          <BaseButton v-permission="'system:user:delete'" variant="link" size="sm" @click="onDelete(row as User)">删除</BaseButton>
        </div>
      </template>
    </BaseTable>
    <div class="flex justify-end">
      <BasePagination
        :page="page"
        :page-size="pageSize"
        :total="total"
        @update:page="(p) => (page = p)"
        @update:page-size="(s) => (pageSize = s)"
      />
    </div>

    <BaseDialog v-model="dlg" :title="dlgTitle" max-width="min(480px, 96vw)">
      <div v-if="mode !== 'role'" class="space-y-4">
        <BaseFormItem label="用户名" required>
          <BaseInput v-model="form.username" :disabled="mode === 'edit'" />
        </BaseFormItem>
        <BaseFormItem label="手机" required>
          <BaseInput v-model="form.phone" maxlength="11" />
        </BaseFormItem>
        <BaseFormItem :label="mode === 'create' ? '密码' : '新密码'" :required="mode === 'create'">
          <BaseInput v-model="form.password" type="password" show-password placeholder="至少 6 位" />
        </BaseFormItem>
        <BaseFormItem label="昵称">
          <BaseInput v-model="form.nickname" />
        </BaseFormItem>
        <BaseFormItem label="角色">
          <BaseSelect v-model="form.role_id" :options="roleOptions" placeholder="选择角色" />
        </BaseFormItem>
        <BaseFormItem v-if="isAdmin" label="门店" required>
          <BaseSelect
            v-model="form.store_code"
            :options="storeOptions"
            placeholder="请选择门店（按编码提交）"
          />
        </BaseFormItem>
        <BaseFormItem v-if="mode === 'edit'" label="状态">
          <BaseSelect
            v-model="form.status"
            :options="[
              { label: '正常', value: 1 },
              { label: '禁用', value: 2 },
            ]"
          />
        </BaseFormItem>
      </div>
      <div v-else class="space-y-4">
        <BaseFormItem label="分配角色">
          <BaseSelect v-model="roleOnlyId" :options="roleOptions" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="submit">确定</BaseButton>
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
  BaseTable,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn } from '@/components/base/types'
import { assignUserRole, createUser, deleteUser, listUsers, updateUser } from '@/api/user'
import { listRoles } from '@/api/role'
import { listAllStores } from '@/api/store'
import type { User } from '@/api/types'
import { useUserStore } from '@/store/user'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()
const userStore = useUserStore()

const isAdmin = computed(() => {
  const c = userStore.userInfo?.role?.code ?? ''
  return c === 'admin' || c === 'super_admin'
})

const columns: BaseTableColumn[] = [
  { key: 'id', label: 'ID', prop: 'id', width: '72px' },
  { key: 'username', label: '用户名', prop: 'username', minWidth: '100px' },
  { key: 'phone', label: '手机', prop: 'phone', width: '120px' },
  { key: 'nickname', label: '昵称', prop: 'nickname', minWidth: '100px' },
  { key: 'role', label: '角色', minWidth: '100px' },
  { key: 'store', label: '门店', minWidth: '120px' },
  { key: 'status', label: '状态', width: '88px' },
  { key: 'actions', label: '操作', width: '200px', fixed: 'right' },
]

const keyword = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const list = ref<User[]>([])
const loading = ref(false)

const { data: roleListData } = useQuery({
  queryKey: ['roles'],
  queryFn: () => listRoles(),
})

const { data: storeListData } = useQuery({
  queryKey: ['stores', 'all'],
  queryFn: () => listAllStores(),
  enabled: isAdmin,
})

const roleOptions = computed<BaseSelectOption[]>(() =>
  (roleListData.value ?? []).map((r) => ({ label: r.name, value: r.id })),
)

const storeOptions = computed<BaseSelectOption[]>(() =>
  (storeListData.value ?? [])
    .filter((s) => s.store_code && String(s.store_code).trim() !== '')
    .map((s) => ({
      label: `${s.name}（${String(s.store_code)}）`,
      value: String(s.store_code),
    })),
)

async function load(): Promise<void> {
  loading.value = true
  try {
    const res = await listUsers({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value || undefined,
    })
    list.value = res.list ?? []
    total.value = res.total ?? 0
  } catch (e) {
    console.error(e)
    list.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

function reload(): void {
  page.value = 1
  void load()
}

watch([page, pageSize], () => {
  void load()
})

void load()

const dlg = ref(false)
const saving = ref(false)
const mode = ref<'create' | 'edit' | 'role'>('create')
const editId = ref(0)
const roleOnlyId = ref<number | undefined>(undefined)

const form = reactive({
  username: '',
  phone: '',
  password: '',
  nickname: '',
  role_id: undefined as number | undefined,
  /** 总部管理员创建/编辑用户时提交的门店编码 */
  store_code: '' as string,
  status: 1,
})

const dlgTitle = computed(() => {
  if (mode.value === 'role') return '分配角色'
  return mode.value === 'edit' ? '编辑用户' : '新增用户'
})

function openCreate(): void {
  mode.value = 'create'
  form.username = ''
  form.phone = ''
  form.password = ''
  form.nickname = ''
  form.role_id = undefined
  const sc = userStore.userInfo?.store?.store_code
  form.store_code = sc != null && String(sc).trim() !== '' ? String(sc) : ''
  form.status = 1
  dlg.value = true
}

function openEdit(row: User): void {
  mode.value = 'edit'
  editId.value = row.id
  form.username = row.username
  form.phone = row.phone
  form.password = ''
  form.nickname = row.nickname ?? ''
  form.role_id = row.role_id
  const sc = row.store?.store_code
  form.store_code = sc != null && String(sc).trim() !== '' ? String(sc) : ''
  form.status = row.status === 2 ? 2 : 1
  dlg.value = true
}

function openRole(row: User): void {
  mode.value = 'role'
  editId.value = row.id
  roleOnlyId.value = row.role_id ?? undefined
  dlg.value = true
}

async function submit(): Promise<void> {
  saving.value = true
  try {
    if (mode.value === 'role') {
      if (roleOnlyId.value === undefined) {
        toast.warning('请选择角色')
        return
      }
      await assignUserRole(editId.value, roleOnlyId.value)
      toast.success('已更新角色')
      dlg.value = false
      await load()
      return
    }
    if (mode.value === 'create') {
      if (!form.password) {
        toast.warning('请填写密码')
        return
      }
      if (isAdmin.value && !String(form.store_code).trim()) {
        toast.warning('请选择门店')
        return
      }
      const body: Record<string, unknown> = {
        username: form.username,
        phone: form.phone,
        password: form.password,
        nickname: form.nickname,
        role_id: form.role_id,
      }
      if (isAdmin.value && String(form.store_code).trim()) {
        body.store_code = String(form.store_code).trim()
      }
      await createUser(body)
    } else {
      const body: Record<string, unknown> = {
        username: form.username,
        phone: form.phone,
        nickname: form.nickname,
        role_id: form.role_id,
        status: form.status,
      }
      if (form.password) body.password = form.password
      if (isAdmin.value && String(form.store_code).trim()) {
        body.store_code = String(form.store_code).trim()
      }
      await updateUser(editId.value, body)
    }
    toast.success('已保存')
    dlg.value = false
    await load()
    await qc.invalidateQueries({ queryKey: ['users'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '操作失败')
  } finally {
    saving.value = false
  }
}

async function onDelete(row: User): Promise<void> {
  const ok = await confirmDialog({ message: `删除用户「${row.username}」？` })
  if (!ok) return
  try {
    await deleteUser(row.id)
    toast.success('已删除')
    await load()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>
