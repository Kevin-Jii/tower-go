<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">物流路线图</h2>
      <div class="flex gap-2">
        <BaseButton variant="ghost" @click="goBack">返回账号池</BaseButton>
        <BaseButton variant="primary" @click="openCreate">新增路线</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(rows as unknown) as Record<string, unknown>[]" :loading="loading"
      min-width="980px">
      <template #cell-stores="{ row }">
        <div class="max-w-[340px] truncate" :title="formatStores(row as ThirdPartyRoute)">
          {{ formatStores(row as ThirdPartyRoute) }}
        </div>
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="routeRowActions(row as ThirdPartyRoute)" />
      </template>
    </BaseTable>

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑路线' : '新增路线'" max-width="min(680px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="路线名称" required>
          <BaseInput v-model="form.name" />
        </BaseFormItem>
        <BaseFormItem label="门店编排">
          <div class="max-h-52 overflow-y-auto rounded border border-[var(--color-border-2)] p-2">
            <label v-for="opt in storeOptions" :key="String(opt.value)"
              class="mb-1 flex cursor-pointer items-center gap-2 text-sm last:mb-0">
              <a-checkbox :model-value="form.store_ids.includes(Number(opt.value))"
                @change="toggleStore(Number(opt.value), $event)" />
              <span>{{ opt.label }}</span>
            </label>
          </div>
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="form.remark" :rows="2" />
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
import { useRouter } from 'vue-router'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import { BaseButton, BaseDialog, BaseFormItem, BaseInput, BaseTable, BaseTableRowActions, BaseTextarea } from '@/components/base'
import type { BaseSelectOption, BaseTableColumn, TableRowAction } from '@/components/base/types'
import type { ThirdPartyRoute } from '@/api/types'
import { deleteThirdPartyRoute, listThirdPartyRoutes, createThirdPartyRoute, updateThirdPartyRoute } from '@/api/thirdPartyRoute'
import { listAllStores } from '@/api/store'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const router = useRouter()
const qc = useQueryClient()
const { data, isLoading: loading } = useQuery({ queryKey: ['third-party-routes'], queryFn: listThirdPartyRoutes })
const rows = computed(() => data.value ?? [])

const { data: storesData } = useQuery({ queryKey: ['stores', 'all'], queryFn: listAllStores })
const storeOptions = computed<BaseSelectOption[]>(() => (storesData.value ?? []).map((s) => ({ label: `${s.name}（${s.store_code || s.id}）`, value: s.id })))

const columns: BaseTableColumn[] = [
  { key: 'name', label: '路线名称', prop: 'name', width: '120px' },
  { key: 'stores', label: '门店编排', minWidth: '320px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '220px', ellipsis: true },
  { key: 'actions', label: '操作', width: '250px', align: 'right' },
]
const dlg = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const saving = ref(false)
const form = reactive<{ name: string; remark: string; store_ids: number[] }>({ name: '', remark: '', store_ids: [] })


function routeRowActions(row: ThirdPartyRoute): TableRowAction[] {
  return [
    { label: '编辑', onClick: () => openEdit(row) },
    { label: '按日期导入', onClick: () => openImport(row) },
    { label: '历史物流单', onClick: () => openHistory(row) },
    { label: '删除', danger: true, onClick: () => void onDelete(row) },
  ]
}

function formatStores(route: ThirdPartyRoute): string {
  return (route.stores ?? []).map((s) => s.store?.name || `门店#${s.store_id}`).join(' -> ') || '-'
}
function goBack(): void { void router.push({ name: 'ThirdPartyAccount' }) }
function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  form.name = ''
  form.remark = ''
  form.store_ids = []
  dlg.value = true
}
function openEdit(row: ThirdPartyRoute): void {
  isEdit.value = true
  editId.value = row.id
  form.name = row.name || ''
  form.remark = row.remark || ''
  form.store_ids = (row.stores ?? []).map((s) => s.store_id)
  dlg.value = true
}
function toggleStore(storeID: number, checked: boolean | string | number | Array<string | number | boolean>): void {
  const on = Array.isArray(checked) ? checked.length > 0 : Boolean(checked)
  const exists = form.store_ids.includes(storeID)
  if (on && !exists) form.store_ids.push(storeID)
  if (!on && exists) form.store_ids = form.store_ids.filter((id) => id !== storeID)
}
async function submit(): Promise<void> {
  if (!form.name.trim()) { toast.warning('请填写路线名称'); return }
  saving.value = true
  try {
    const payload = { name: form.name.trim(), remark: form.remark.trim(), store_ids: form.store_ids }
    if (isEdit.value) await updateThirdPartyRoute(editId.value, payload)
    else await createThirdPartyRoute(payload)
    toast.success('已保存')
    dlg.value = false
    await qc.invalidateQueries({ queryKey: ['third-party-routes'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}
async function onDelete(row: ThirdPartyRoute): Promise<void> {
  const ok = await confirmDialog({ message: `确认删除路线「${row.name}」？` })
  if (!ok) return
  try {
    await deleteThirdPartyRoute(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['third-party-routes'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
function openImport(row: ThirdPartyRoute): void {
  void router.push({
    name: 'ThirdPartyRouteImport',
    query: { route_id: String(row.id) },
  })
}

function openHistory(row: ThirdPartyRoute): void {
  void router.push({
    name: 'ThirdPartyRouteHistory',
    query: { route_id: String(row.id) },
  })
}
</script>
