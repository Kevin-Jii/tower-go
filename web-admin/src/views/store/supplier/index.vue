<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">供应商管理</h2>
      <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-56" placeholder="名称 / 编码" clearable @enter="reload" />
        <BaseButton variant="primary" @click="reload">查询</BaseButton>
        <BaseButton v-permission="'supplier:add'" variant="primary" @click="openCreate">新增</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading" min-width="880px">
      <template #cell-status="{ row }">
        {{ (row as Supplier).status === 1 ? '启用' : '禁用' }}
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-nowrap items-center justify-end gap-2 shrink-0" @click.stop>
          <BaseButton v-permission="'supplier:edit'" variant="link" size="sm" @click="openEdit(row as Supplier)">编辑</BaseButton>
          <BaseButton v-permission="'supplier:delete'" variant="link" size="sm" @click="onDelete(row as Supplier)">删除</BaseButton>
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

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑供应商' : '新增供应商'" max-width="min(520px, 96vw)">
      <div class="space-y-4 max-h-[70vh] overflow-y-auto pr-1">
        <BaseFormItem label="名称" required>
          <BaseInput v-model="form.supplier_name" placeholder="供应商名称" />
        </BaseFormItem>
        <BaseFormItem label="联系人">
          <BaseInput v-model="form.contact_person" />
        </BaseFormItem>
        <BaseFormItem label="电话">
          <BaseInput v-model="form.contact_phone" />
        </BaseFormItem>
        <BaseFormItem label="邮箱">
          <BaseInput v-model="form.contact_email" />
        </BaseFormItem>
        <BaseFormItem label="地址">
          <BaseInput v-model="form.supplier_address" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseTextarea v-model="form.remark" :rows="2" />
        </BaseFormItem>
        <BaseFormItem v-if="isEdit" label="状态">
          <BaseSelect
            v-model="form.status"
            :options="[
              { label: '启用', value: 1 },
              { label: '禁用', value: 0 },
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
import { createSupplier, deleteSupplier, listSuppliers, updateSupplier } from '@/api/supplier'
import type { Supplier } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()
const keyword = ref('')
const page = ref(1)
const pageSize = ref(10)
const queryKey = computed(() => ['suppliers', page.value, pageSize.value, keyword.value.trim()] as const)

const { data: pageData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: () =>
    listSuppliers({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value.trim() || undefined,
    }),
})

const list = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)

function reload(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['suppliers'] })
}

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['suppliers'] })
})

const columns: BaseTableColumn[] = [
  { key: 'supplier_code', label: '编码', prop: 'supplier_code', width: '120px' },
  { key: 'supplier_name', label: '名称', prop: 'supplier_name', minWidth: '140px', ellipsis: true },
  { key: 'contact_person', label: '联系人', prop: 'contact_person', width: '100px' },
  { key: 'contact_phone', label: '电话', prop: 'contact_phone', width: '120px' },
  { key: 'status', label: '状态', width: '72px' },
  { key: 'actions', label: '操作', width: '140px', align: 'right' },
]

const dlg = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref(0)

const form = reactive({
  supplier_name: '',
  contact_person: '',
  contact_phone: '',
  contact_email: '',
  supplier_address: '',
  remark: '',
  status: 1,
})

function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  form.supplier_name = ''
  form.contact_person = ''
  form.contact_phone = ''
  form.contact_email = ''
  form.supplier_address = ''
  form.remark = ''
  form.status = 1
  dlg.value = true
}

function openEdit(row: Supplier): void {
  isEdit.value = true
  editId.value = row.id
  form.supplier_name = row.supplier_name ?? ''
  form.contact_person = row.contact_person ?? ''
  form.contact_phone = row.contact_phone ?? ''
  form.contact_email = row.contact_email ?? ''
  form.supplier_address = row.supplier_address ?? ''
  form.remark = row.remark ?? ''
  form.status = row.status === 1 ? 1 : 0
  dlg.value = true
}

async function save(): Promise<void> {
  if (!form.supplier_name.trim()) {
    toast.warning('请填写供应商名称')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateSupplier(editId.value, {
        supplier_name: form.supplier_name.trim(),
        contact_person: form.contact_person.trim(),
        contact_phone: form.contact_phone.trim(),
        contact_email: form.contact_email.trim() || undefined,
        supplier_address: form.supplier_address.trim(),
        remark: form.remark.trim(),
        status: form.status,
      })
    } else {
      await createSupplier({
        supplier_name: form.supplier_name.trim(),
        contact_person: form.contact_person.trim(),
        contact_phone: form.contact_phone.trim(),
        contact_email: form.contact_email.trim() || undefined,
        supplier_address: form.supplier_address.trim(),
        remark: form.remark.trim(),
      })
    }
    toast.success('已保存')
    dlg.value = false
    await qc.invalidateQueries({ queryKey: ['suppliers'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function onDelete(row: Supplier): Promise<void> {
  const ok = await confirmDialog({ message: `删除供应商「${row.supplier_name}」？` })
  if (!ok) return
  try {
    await deleteSupplier(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['suppliers'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>
