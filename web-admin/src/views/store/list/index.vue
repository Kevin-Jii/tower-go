<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">门店列表</h2>
      <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-56" placeholder="名称 / 编码 / 电话" clearable @enter="page = 1" />
        <BaseButton variant="primary" @click="page = 1">查询</BaseButton>
        <BaseButton v-permission="'store:add'" variant="primary" @click="openCreate">新增门店</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(pagedRows as unknown) as Record<string, unknown>[]" :loading="loading" min-width="960px">
      <template #cell-status="{ row }">
        {{ statusLabel((row as Store).status) }}
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-wrap gap-1 justify-end" @click.stop>
          <BaseButton v-permission="'store:edit'" variant="link" size="sm" @click="openEdit(row as Store)">编辑</BaseButton>
          <BaseButton v-permission="'store:delete'" variant="link" size="sm" @click="onDelete(row as Store)">删除</BaseButton>
        </div>
      </template>
    </BaseTable>

    <div class="flex justify-end">
      <BasePagination
        :page="page"
        :page-size="pageSize"
        :total="filteredList.length"
        @update:page="(p) => (page = p)"
        @update:page-size="(s) => (pageSize = s)"
      />
    </div>

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
          <BaseSelect
            v-model="form.status"
            :options="[
              { label: '正常', value: 1 },
              { label: '停业', value: 2 },
            ]"
          />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="save">保存</BaseButton>
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
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import { createStore, deleteStore, listStores, updateStore } from '@/api/store'
import type { Store } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()

const columns: BaseTableColumn[] = [
  { key: 'id', label: 'ID', prop: 'id', width: '72px' },
  { key: 'store_code', label: '编码', prop: 'store_code', width: '100px' },
  { key: 'name', label: '名称', prop: 'name', minWidth: '120px' },
  { key: 'phone', label: '电话', prop: 'phone', width: '120px' },
  { key: 'address', label: '地址', prop: 'address', minWidth: '160px', ellipsis: true },
  { key: 'business_hours', label: '营业时间', prop: 'business_hours', width: '120px' },
  { key: 'contact_person', label: '联系人', prop: 'contact_person', width: '100px' },
  { key: 'status', label: '状态', width: '88px' },
  { key: 'actions', label: '操作', width: '140px', fixed: 'right' },
]

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

const form = reactive({
  name: '',
  phone: '',
  address: '',
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

async function onDelete(row: Store): Promise<void> {
  const ok = await confirmDialog({ message: `删除门店「${row.name}」？此操作不可恢复。` })
  if (!ok) return
  try {
    await deleteStore(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['stores'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>
