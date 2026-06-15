<template>
  <div class="flex flex-col gap-4">
    <a-card :bordered="true" class="min-w-0">
      <template #title>
        <span class="text-base font-semibold">门店列表</span>
      </template>
      <template #extra>
        <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
          <BaseInput v-model="keyword" class="w-full sm:w-56" placeholder="名称 / 编码 / 电话" clearable @enter="page = 1" />
          <BaseButton variant="primary" @click="page = 1">查询</BaseButton>
          <BaseButton v-permission="'store:add'" variant="primary" @click="openCreate">新增门店</BaseButton>
        </div>
      </template>

      <div v-if="pagedRows.length" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        <a-card v-for="row in pagedRows" :key="row.id"
          class="store-card rounded-2xl border border-slate-200/80 shadow-sm hover:shadow-md transition-shadow"
          :body-style="{ padding: '14px 14px 10px' }">
          <div class="flex items-start justify-between gap-2">
            <div class="min-w-0">
              <p class="m-0 text-base font-semibold text-slate-900 truncate">
                {{ row.name || '-' }}
                <span class="text-slate-500 font-medium">【{{ row.store_code || '-' }}】</span>
              </p>
              <p class="m-0 mt-1 text-sm text-slate-500 truncate">负责人：{{ row.contact_person || '-' }}</p>
            </div>
            <span class="inline-flex shrink-0 items-center rounded-full px-2.5 py-1 text-xs font-semibold"
              :class="row.status === 1 ? 'bg-emerald-100 text-emerald-700' : 'bg-rose-100 text-rose-700'">
              {{ statusLabel(row.status) }}
            </span>
          </div>

          <div class="mt-3 space-y-1.5 text-sm text-slate-700">
            <p class="m-0"><span class="text-slate-500"><b>营业时间：</b></span>{{ row.business_hours || '-' }}</p>
            <p class="m-0"><span class="text-slate-500"><b>电话：</b></span>{{ row.phone || '-' }}</p>
            <p class="m-0"><span class="text-slate-500"><b>地址：</b></span>{{ row.address || '-' }}</p>
            <p class="m-0"><span class="text-slate-500"><b>第三方账号：</b></span>{{ row.third_party_account?.name || '-' }}
            </p>
          </div>

          <div class="mt-4 flex flex-wrap gap-2 border-t border-slate-100 pt-3">
            <BaseButton size="sm" class="store-card-btn store-card-btn--third" v-permission="'store:menu'"
              @click="openBindThirdAccount(row)">
              第三方绑定
            </BaseButton>
            <BaseButton size="sm" class="store-card-btn store-card-btn--supplier" v-permission="'store:menu'"
              @click="openBindSupplier(row)">
              绑定供应商
            </BaseButton>
            <BaseButton size="sm" class="store-card-btn store-card-btn--permission" v-permission="'store:menu'"
              @click="openPermissionDrawer(row)">
              配置权限
            </BaseButton>
            <BaseButton size="sm" class="store-card-btn store-card-btn--edit" v-permission="'store:edit'"
              @click="openEdit(row)">
              编辑
            </BaseButton>
            <BaseButton size="sm" class="store-card-btn store-card-btn--status" v-permission="'store:edit'"
              :loading="statusChanging[row.id] === true" @click="setBusinessStatus(row, row.status === 1 ? 2 : 1)">
              {{ row.status === 1 ? '设为停业' : '设为正常' }}
            </BaseButton>
          </div>
        </a-card>
      </div>
      <a-empty v-else-if="!loading" description="暂无门店数据" />
      <div v-else class="py-8">
        <a-skeleton :animation="true">
          <a-skeleton-line :rows="4" />
        </a-skeleton>
      </div>

      <div class="flex justify-end mt-3">
        <BasePagination :page="page" :page-size="pageSize" :total="filteredList.length" @update:page="(p) => (page = p)"
          @update:page-size="(s) => (pageSize = s)" />
      </div>
    </a-card>

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑门店' : '新增门店'" max-width="min(520px, 96vw)">
      <div class="max-h-[70vh] overflow-y-auto space-y-4 pr-1">
        <p v-if="!isEdit" class="m-0 text-xs text-[var(--color-text-3)]">门店编码由系统自动生成，创建后即可在列表中查看。</p>
        <BaseFormItem v-if="isEdit" label="门店编码">
          <BaseInput :model-value="codeDisplay" disabled />
        </BaseFormItem>
        <BaseFormItem label="门店名称" required>
          <BaseInput v-model="form.name" placeholder="门店名称" />
        </BaseFormItem>
        <BaseFormItem label="联系电话">
          <BaseInput v-model="form.phone" placeholder="可选" />
        </BaseFormItem>
        <BaseFormItem label="地址">
          <BaseInput v-model="form.address" placeholder="可选" />
        </BaseFormItem>
        <BaseFormItem label="归属区">
          <BaseSelect v-model="form.administrative_unit" :options="administrativeUnitOptions" placeholder="请选择归属区"
            clearable />
        </BaseFormItem>
        <BaseFormItem label="营业时间">
          <BaseInput v-model="form.business_hours" placeholder="如 10:00-22:00" />
        </BaseFormItem>
        <BaseFormItem label="联系人">
          <BaseInput v-model="form.contact_person" placeholder="可选" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="form.remark" :rows="2" placeholder="可选" />
        </BaseFormItem>
        <BaseFormItem v-if="isEdit" label="状态">
          <BaseSelect v-model="form.status" :options="[
            { label: '正常', value: 1 },
            { label: '停业', value: 2 },
          ]" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="save">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="bindDlg" title="绑定供应商" max-width="min(520px, 96vw)">
      <div class="space-y-4">
        <p class="m-0 text-sm text-slate-600">
          当前门店：<span class="font-medium text-slate-900">{{ bindStore?.name || '-' }}</span>
        </p>
        <BaseFormItem label="新增绑定">
          <div class="flex gap-2">
            <BaseSelect v-model="bindSupplierId" class="flex-1" :options="supplierOptions" placeholder="请选择供应商" />
            <BaseButton variant="primary" :loading="bindSaving" @click="onBindSupplier">绑定</BaseButton>
          </div>
        </BaseFormItem>
        <BaseFormItem label="已绑定供应商">
          <div class="min-h-14 rounded border border-[var(--color-border-2)] p-3">
            <div v-if="boundSuppliers.length" class="flex flex-wrap gap-2">
              <span v-for="b in boundSuppliers" :key="b.id"
                class="inline-flex items-center gap-2 rounded bg-[var(--color-fill-2)] px-2 py-1 text-xs">
                {{ b.supplier?.supplier_name || `供应商#${b.supplier_id}` }}
                <button v-permission="'store:menu'"
                  class="cursor-pointer border-none bg-transparent p-0 text-[var(--color-danger-6)]"
                  @click="onUnbindSupplier(b.supplier_id)">
                  解绑
                </button>
              </span>
            </div>
            <p v-else class="m-0 text-xs text-slate-400">暂无已绑定供应商</p>
          </div>
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="bindDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="thirdAccountDlg" title="绑定第三方账号" max-width="min(520px, 96vw)">
      <div class="space-y-4">
        <p class="m-0 text-sm text-slate-600">
          当前门店：<span class="font-medium text-slate-900">{{ bindStore?.name || '-' }}</span>
        </p>
        <BaseFormItem label="第三方账号">
          <BaseSelect v-model="bindThirdAccountId" :options="thirdAccountOptions" placeholder="请选择第三方账号" />
        </BaseFormItem>
        <p class="m-0 text-xs text-slate-400">一个门店有且只能绑定一个第三方账号，保存时会覆盖原绑定。</p>
      </div>
      <template #footer>
        <BaseButton variant="ghost" :loading="bindThirdSaving" @click="onUnbindThirdAccount">解绑</BaseButton>
        <BaseButton variant="ghost" @click="thirdAccountDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="bindThirdSaving" @click="onBindThirdAccount">保存</BaseButton>
      </template>
    </BaseDialog>

    <a-drawer
      :visible="permissionDrawer"
      placement="right"
      :width="620"
      :drawer-style="{ maxWidth: '96vw' }"
      :mask-closable="true"
      unmount-on-close
      @cancel="permissionDrawer = false"
      @update:visible="permissionDrawer = $event"
    >
      <template #title>配置门店权限</template>
      <div class="space-y-4">
        <div class="rounded border border-slate-200 bg-slate-50 px-3 py-2 text-sm text-slate-600">
          当前门店：<span class="font-semibold text-slate-900">{{ permissionStore?.name || '-' }}</span>
        </div>
        <BaseFormItem label="角色">
          <BaseSelect v-model="permissionRoleId" :options="roleOptions" @update:model-value="onPermissionRoleChange" />
        </BaseFormItem>
        <div class="permission-tree-panel">
          <BaseTreeCheck v-model="checkedMenuIds" :nodes="menuTreeNodes" :check-strictly="false" />
        </div>
      </div>
      <template #footer>
        <div class="flex justify-end gap-2">
          <BaseButton variant="ghost" @click="permissionDrawer = false">取消</BaseButton>
          <BaseButton variant="primary" :loading="permissionSaving" @click="saveStoreRoleMenus">保存</BaseButton>
        </div>
      </template>
    </a-drawer>
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
  BaseTextarea,
  BaseTreeCheck,
} from '@/components/base'
import { assignStoreRoleMenus, fetchMenuTree, fetchStoreRoleMenuIds, fetchStoreRoleMenuPermissions } from '@/api/menu'
import { bindStoreThirdPartyAccount, createStore, listStores, updateStore } from '@/api/store'
import { bindStoreSuppliers, listStoreBoundSuppliers, unbindStoreSuppliers } from '@/api/storeSupplier'
import { listSuppliers } from '@/api/supplier'
import { listThirdPartyAccounts } from '@/api/thirdPartyAccount'
import { listRoles } from '@/api/role'
import { listDictDataByTypeCode } from '@/api/dict'
import type { DictData, Menu, Role, Store, StoreSupplierBinding } from '@/api/types'
import type { BaseSelectOption, BaseTreeNode } from '@/components/base/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()

const { data: pageData, isLoading: loading } = useQuery({
  queryKey: ['stores', 'list'],
  queryFn: () => listStores(),
})

const list = computed(() => pageData.value?.list ?? [])

const keyword = ref('')
const page = ref(1)
const pageSize = ref(10)

const filteredList = computed(() => {
  const k = keyword.value.trim().toLowerCase()
  if (!k) return list.value
  return list.value.filter((s) => {
    const parts = [s.name, s.store_code ?? '', s.phone ?? '', s.address ?? ''].join(' ').toLowerCase()
    return parts.includes(k)
  })
})

const pagedRows = computed(() => {
  const start = (page.value - 1) * pageSize.value
  return filteredList.value.slice(start, start + pageSize.value)
})

watch([keyword, pageSize], () => {
  page.value = 1
})

watch(filteredList, (rows) => {
  const maxPage = Math.max(1, Math.ceil(rows.length / pageSize.value) || 1)
  if (page.value > maxPage) page.value = maxPage
})

const dlg = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const statusChanging = reactive<Record<number, boolean>>({})

const form = reactive({
  name: '',
  phone: '',
  address: '',
  administrative_unit: '',
  business_hours: '',
  contact_person: '',
  remark: '',
  status: 1,
})

const codeDisplay = ref('')

function statusLabel(v?: number): string {
  if (v === 2) return '停业'
  return '正常'
}

function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  form.name = ''
  form.phone = ''
  form.address = ''
  form.administrative_unit = ''
  form.business_hours = ''
  form.contact_person = ''
  form.remark = ''
  form.status = 1
  codeDisplay.value = ''
  dlg.value = true
}

function openEdit(row: Store): void {
  isEdit.value = true
  editId.value = row.id
  form.name = row.name ?? ''
  form.phone = row.phone ?? ''
  form.address = row.address ?? ''
  form.administrative_unit = row.administrative_unit ?? ''
  form.business_hours = row.business_hours ?? ''
  form.contact_person = row.contact_person ?? ''
  form.remark = row.remark ?? ''
  form.status = row.status === 2 ? 2 : 1
  codeDisplay.value = row.store_code != null && String(row.store_code) !== '' ? String(row.store_code) : '-'
  dlg.value = true
}

async function save(): Promise<void> {
  if (!form.name.trim()) {
    toast.warning('请填写门店名称')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateStore(editId.value, {
        name: form.name.trim(),
        phone: form.phone.trim(),
        address: form.address.trim(),
        administrative_unit: form.administrative_unit.trim(),
        business_hours: form.business_hours.trim(),
        contact_person: form.contact_person.trim(),
        remark: form.remark.trim(),
        status: form.status,
      })
    } else {
      await createStore({
        name: form.name.trim(),
        phone: form.phone.trim(),
        address: form.address.trim(),
        administrative_unit: form.administrative_unit.trim(),
        business_hours: form.business_hours.trim(),
        contact_person: form.contact_person.trim(),
        remark: form.remark.trim(),
      })
    }
    toast.success('已保存')
    dlg.value = false
    await qc.invalidateQueries({ queryKey: ['stores'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function setBusinessStatus(row: Store, targetStatus: 1 | 2): Promise<void> {
  if (row.status === targetStatus) return
  statusChanging[row.id] = true
  try {
    await updateStore(row.id, { status: targetStatus })
    row.status = targetStatus
    toast.success(`已设为${targetStatus === 1 ? '正常营业' : '停业'}`)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '状态更新失败')
  } finally {
    statusChanging[row.id] = false
  }
}

const bindDlg = ref(false)
const bindSaving = ref(false)
const bindStore = ref<Store | null>(null)
const bindSupplierId = ref<number | undefined>(undefined)
const boundSuppliers = ref<StoreSupplierBinding[]>([])
const supplierOptions = ref<BaseSelectOption[]>([])
const thirdAccountDlg = ref(false)
const bindThirdSaving = ref(false)
const bindThirdAccountId = ref<number | undefined>(undefined)
const thirdAccountOptions = ref<BaseSelectOption[]>([])

const { data: administrativeUnitDictData } = useQuery({
  queryKey: ['dict-data', 'ADMINISTRATIVEUNIT'],
  queryFn: () => listDictDataByTypeCode('ADMINISTRATIVEUNIT'),
})

const { data: rawMenuTree } = useQuery({
  queryKey: ['menus', 'tree'],
  queryFn: fetchMenuTree,
})

const { data: rolesData } = useQuery({
  queryKey: ['roles', 'store-permission'],
  queryFn: () => listRoles(),
})

const roleOptions = computed<BaseSelectOption[]>(() =>
  (rolesData.value ?? [])
    .filter((role: Role) => role.code === 'store_admin' || role.code === 'staff')
    .map((role: Role) => ({ label: `${role.name}（${role.code}）`, value: role.id })),
)

const menuTreeNodes = computed(() => filterMenuTree(rawMenuTree.value ?? []) as unknown as BaseTreeNode[])

const administrativeUnitOptions = computed<BaseSelectOption[]>(() => {
  const rows = (administrativeUnitDictData.value ?? []) as DictData[]
  return rows
    .filter((x) => Number(x.status ?? 1) === 1)
    .sort((a, b) => Number(a.sort ?? 0) - Number(b.sort ?? 0))
    .map((x) => ({
      label: x.label || x.value,
      value: x.value,
    }))
})

async function loadSupplierOptions(): Promise<void> {
  const pageData = await listSuppliers({ page: 1, page_size: 100 })
  supplierOptions.value = (pageData.list ?? []).map((s) => ({
    label: `${s.supplier_name}（${s.supplier_code}）`,
    value: s.id,
  }))
}

async function loadBoundSuppliers(): Promise<void> {
  if (!bindStore.value?.id) return
  boundSuppliers.value = await listStoreBoundSuppliers({ store_id: bindStore.value.id })
}

async function openBindSupplier(row: Store): Promise<void> {
  bindStore.value = row
  bindSupplierId.value = undefined
  bindDlg.value = true
  try {
    await Promise.all([loadSupplierOptions(), loadBoundSuppliers()])
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载绑定信息失败')
  }
}

async function openBindThirdAccount(row: Store): Promise<void> {
  bindStore.value = row
  bindThirdAccountId.value = row.third_party_account_id ?? undefined
  thirdAccountDlg.value = true
  try {
    const rows = await listThirdPartyAccounts('')
    thirdAccountOptions.value = rows.map((x) => ({
      label: `${x.name}（${x.login_name}）`,
      value: x.id,
    }))
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载第三方账号失败')
  }
}

async function onBindSupplier(): Promise<void> {
  if (!bindStore.value?.id || !bindSupplierId.value) {
    toast.warning('请选择要绑定的供应商')
    return
  }
  bindSaving.value = true
  try {
    await bindStoreSuppliers({
      store_id: bindStore.value.id,
      supplier_ids: [bindSupplierId.value],
    })
    toast.success('绑定成功')
    bindSupplierId.value = undefined
    await loadBoundSuppliers()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '绑定失败')
  } finally {
    bindSaving.value = false
  }
}

async function onUnbindSupplier(supplierId: number): Promise<void> {
  if (!bindStore.value?.id) return
  const ok = await confirmDialog({ message: '确认解绑该供应商？' })
  if (!ok) return
  try {
    await unbindStoreSuppliers({
      store_id: bindStore.value.id,
      supplier_ids: [supplierId],
    })
    toast.success('已解绑')
    await loadBoundSuppliers()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '解绑失败')
  }
}

async function onBindThirdAccount(): Promise<void> {
  if (!bindStore.value?.id) return
  if (!bindThirdAccountId.value) {
    toast.warning('请选择第三方账号')
    return
  }
  bindThirdSaving.value = true
  try {
    await bindStoreThirdPartyAccount(bindStore.value.id, bindThirdAccountId.value)
    toast.success('绑定成功')
    thirdAccountDlg.value = false
    await qc.invalidateQueries({ queryKey: ['stores'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '绑定失败')
  } finally {
    bindThirdSaving.value = false
  }
}

async function onUnbindThirdAccount(): Promise<void> {
  if (!bindStore.value?.id) return
  bindThirdSaving.value = true
  try {
    await bindStoreThirdPartyAccount(bindStore.value.id, null)
    toast.success('已解绑')
    thirdAccountDlg.value = false
    await qc.invalidateQueries({ queryKey: ['stores'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '解绑失败')
  } finally {
    bindThirdSaving.value = false
  }
}

const permissionDrawer = ref(false)
const permissionSaving = ref(false)
const permissionStore = ref<Store | null>(null)
const permissionRoleId = ref<number | undefined>(undefined)
const checkedMenuIds = ref<number[]>([])
const storeRolePermSnapshot = ref<Record<number, number>>({})

function normalizeCheckedMenuIds(value: unknown): number[] {
  if (Array.isArray(value)) return value.map((x) => Number(x)).filter((x) => Number.isFinite(x))
  if (value && typeof value === 'object' && 'checked' in value) {
    const checked = (value as { checked?: unknown }).checked
    if (Array.isArray(checked)) return checked.map((x: unknown) => Number(x)).filter((x: number) => Number.isFinite(x))
  }
  return []
}

function filterMenuTree(nodes: Menu[]): Menu[] {
  return (nodes ?? [])
    .filter((n) => n.status !== 0)
    .map((n) => ({
      ...n,
      children: n.children?.length ? filterMenuTree(n.children) : [],
    }))
}

async function openPermissionDrawer(row: Store): Promise<void> {
  permissionStore.value = row
  permissionDrawer.value = true
  const defaultRole = (rolesData.value ?? []).find((role) => role.code === 'store_admin') ?? (rolesData.value ?? [])[0]
  permissionRoleId.value = defaultRole?.id
  if (permissionRoleId.value) {
    await loadStoreRoleMenus()
  } else {
    checkedMenuIds.value = []
  }
}

function onPermissionRoleChange(): void {
  void loadStoreRoleMenus()
}

async function loadStoreRoleMenus(): Promise<void> {
  if (!permissionStore.value?.id || !permissionRoleId.value) return
  try {
    const [ids, permMap] = await Promise.all([
      fetchStoreRoleMenuIds(permissionStore.value.id, permissionRoleId.value),
      fetchStoreRoleMenuPermissions(permissionStore.value.id, permissionRoleId.value).catch(() => ({} as Record<number, number>)),
    ])
    checkedMenuIds.value = ids.map((x) => Number(x))
    storeRolePermSnapshot.value = permMap
  } catch (e: unknown) {
    checkedMenuIds.value = []
    storeRolePermSnapshot.value = {}
    toast.error(e instanceof Error ? e.message : '加载门店权限失败')
  }
}

async function saveStoreRoleMenus(): Promise<void> {
  if (!permissionStore.value?.id || !permissionRoleId.value) {
    toast.warning('请选择门店和角色')
    return
  }
  const keys = normalizeCheckedMenuIds(checkedMenuIds.value)
  const perms: Record<number, number> = {}
  for (const id of keys) {
    perms[id] = storeRolePermSnapshot.value[id] ?? 15
  }
  permissionSaving.value = true
  try {
    await assignStoreRoleMenus({
      store_id: permissionStore.value.id,
      role_id: permissionRoleId.value,
      menu_ids: keys,
      perms,
    })
    toast.success('门店权限已保存，相关账号重新登录后生效')
    permissionDrawer.value = false
    await qc.invalidateQueries({ queryKey: ['menus'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存门店权限失败')
  } finally {
    permissionSaving.value = false
  }
}
</script>

<style scoped>
.store-card {
  border-radius: 16px;
}

:deep(.store-card-btn) {
  border: none !important;
  font-weight: 600 !important;
}

:deep(.store-card-btn--third) {
  background: #eef2ff !important;
  color: #4338ca !important;
}

:deep(.store-card-btn--supplier) {
  background: #ecfeff !important;
  color: #0e7490 !important;
}

:deep(.store-card-btn--permission) {
  background: #f5f3ff !important;
  color: #6d28d9 !important;
}

:deep(.store-card-btn--edit) {
  background: #f0fdf4 !important;
  color: #166534 !important;
}

:deep(.store-card-btn--status) {
  background: #fff1f2 !important;
  color: #be123c !important;
}

.permission-tree-panel {
  min-height: calc(100vh - 230px);
  max-height: calc(100vh - 230px);
  overflow: auto;
  padding: 12px;
  border: 1px solid var(--color-border-2);
  border-radius: 8px;
  background: var(--color-fill-1);
}
</style>
