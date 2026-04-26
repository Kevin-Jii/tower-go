<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">会员管理</h2>
      <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-48" placeholder="手机 / UID" clearable @enter="reload" />
        <BaseButton variant="primary" @click="reload">查询</BaseButton>
        <BaseButton v-permission="'store:member:add'" variant="primary" @click="openCreate">新增会员</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading" min-width="1000px">
      <template #cell-balance="{ row }">
        {{ fmtMoney((row as MemberRow).balance) }}
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
          <BaseButton v-permission="'store:member:edit'" variant="link" size="sm" @click="openEdit(row as MemberRow)">编辑</BaseButton>
          <BaseButton v-permission="'store:member:balance'" variant="link" size="sm" @click="openAdjust(row as MemberRow)">调余额</BaseButton>
          <BaseButton v-permission="'store:member:delete'" variant="link" size="sm" @click="onDelete(row as MemberRow)">删除</BaseButton>
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

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑会员' : '新增会员'" max-width="min(440px, 96vw)">
      <div class="space-y-4">
        <BaseFormItem label="手机号" required>
          <BaseInput v-model="form.phone" maxlength="20" />
        </BaseFormItem>
        <BaseFormItem label="姓名">
          <BaseInput v-model="form.name" />
        </BaseFormItem>
        <BaseFormItem v-if="!isEdit" label="等级">
          <BaseNumberInput v-model="form.level" :min="1" :step="1" />
        </BaseFormItem>
        <BaseFormItem v-if="isEdit" label="积分">
          <BaseNumberInput v-model="form.points" :min="0" :step="1" />
        </BaseFormItem>
        <BaseFormItem v-if="isEdit" label="等级">
          <BaseNumberInput v-model="form.level" :min="1" :step="1" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="dlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="saving" @click="save">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="adjDlg" title="调整余额（调增/调减）" max-width="min(400px, 96vw)">
      <p class="text-sm text-slate-600 m-0 mb-3">当前余额：{{ adjMember ? fmtMoney(adjMember.balance) : '-' }}</p>
      <div class="space-y-4">
        <BaseFormItem label="类型" required>
          <BaseSelect
            v-model="adjForm.type"
            :options="[
              { label: '调增', value: 4 },
              { label: '调减', value: 5 },
            ]"
          />
        </BaseFormItem>
        <BaseFormItem label="金额" required>
          <BaseInput v-model="adjForm.amount" placeholder="如 10.00" />
        </BaseFormItem>
        <BaseFormItem label="备注">
          <BaseInput v-model="adjForm.remark" />
        </BaseFormItem>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="adjDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="adjSaving" @click="submitAdjust">提交</BaseButton>
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
  BasePagination,
  BaseSelect,
  BaseTable,
} from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import { adjustMemberBalance, createMember, deleteMember, listMembers, updateMember } from '@/api/member'
import type { MemberRow } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const qc = useQueryClient()
const keyword = ref('')
const page = ref(1)
const pageSize = ref(10)

const queryKey = computed(() => ['members', page.value, pageSize.value, keyword.value.trim()] as const)

const { data: pageData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: () =>
    listMembers({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value.trim() || undefined,
    }),
})

const list = computed(() => pageData.value?.list ?? [])
const total = computed(() => pageData.value?.total ?? 0)

function reload(): void {
  page.value = 1
  void qc.invalidateQueries({ queryKey: ['members'] })
}

watch([page, pageSize], () => {
  void qc.invalidateQueries({ queryKey: ['members'] })
})

const columns: BaseTableColumn[] = [
  { key: 'id', label: 'ID', prop: 'id', width: '72px' },
  { key: 'phone', label: '手机', prop: 'phone', width: '120px' },
  { key: 'name', label: '姓名', prop: 'name', minWidth: '100px' },
  { key: 'balance', label: '余额', width: '100px' },
  { key: 'points', label: '积分', prop: 'points', width: '72px' },
  { key: 'level', label: '等级', prop: 'level', width: '72px' },
  { key: 'actions', label: '操作', width: '260px', align: 'right' },
]

function fmtMoney(v: string | number): string {
  if (v === '' || v == null) return '-'
  const n = typeof v === 'number' ? v : Number(v)
  if (Number.isFinite(n)) return n.toFixed(2)
  return String(v)
}

const dlg = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const form = reactive({
  phone: '',
  name: '',
  level: 1,
  points: 0,
})

function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  form.phone = ''
  form.name = ''
  form.level = 1
  dlg.value = true
}

function openEdit(row: MemberRow): void {
  isEdit.value = true
  editId.value = row.id
  form.phone = row.phone ?? ''
  form.name = row.name ?? ''
  form.level = row.level ?? 1
  form.points = row.points ?? 0
  dlg.value = true
}

async function save(): Promise<void> {
  if (!form.phone.trim()) {
    toast.warning('请填写手机号')
    return
  }
  saving.value = true
  try {
    if (isEdit.value) {
      await updateMember(editId.value, {
        phone: form.phone.trim(),
        name: form.name.trim() || undefined,
        points: form.points,
        level: form.level,
      })
    } else {
      await createMember({
        phone: form.phone.trim(),
        name: form.name.trim() || undefined,
        level_id: form.level,
      })
    }
    toast.success('已保存')
    dlg.value = false
    await qc.invalidateQueries({ queryKey: ['members'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

const adjDlg = ref(false)
const adjSaving = ref(false)
const adjMember = ref<MemberRow | null>(null)
const adjForm = reactive({
  type: 4,
  amount: '',
  remark: '',
})

function openAdjust(row: MemberRow): void {
  adjMember.value = row
  adjForm.type = 4
  adjForm.amount = ''
  adjForm.remark = ''
  adjDlg.value = true
}

async function submitAdjust(): Promise<void> {
  if (!adjMember.value) return
  if (!adjForm.amount.trim()) {
    toast.warning('请填写金额')
    return
  }
  adjSaving.value = true
  try {
    const updated = await adjustMemberBalance(adjMember.value.id, {
      amount: adjForm.amount.trim(),
      type: adjForm.type,
      remark: adjForm.remark.trim(),
      version: adjMember.value.version,
    })
    toast.success('已调整')
    adjDlg.value = false
    adjMember.value = updated
    await qc.invalidateQueries({ queryKey: ['members'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    adjSaving.value = false
  }
}

async function onDelete(row: MemberRow): Promise<void> {
  const ok = await confirmDialog({ message: `删除会员「${row.phone}」？` })
  if (!ok) return
  try {
    await deleteMember(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['members'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}
</script>
