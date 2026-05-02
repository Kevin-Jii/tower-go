<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">第三方账号池</h2>
      <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-60" placeholder="账号名/手机号关键词" clearable @enter="reload" />
        <BaseButton variant="primary" @click="reload">查询</BaseButton>
        <BaseButton variant="ghost" @click="goRouteMap">物流路线图</BaseButton>
        <BaseButton variant="primary" @click="openCreate">新增账号</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading"
      min-width="1240px">
      <template #cell-is_enabled="{ row }">
        <BaseSwitch :model-value="(row as ThirdPartyAccount).is_enabled"
          @update:model-value="toggleEnabled(row as ThirdPartyAccount, $event)" />
      </template>
      <template #cell-last_test_ok="{ row }">
        <span :class="(row as ThirdPartyAccount).last_test_ok ? 'text-emerald-600' : 'text-rose-600'">
          {{ (row as ThirdPartyAccount).last_test_ok ? '成功' : '失败/未测' }}
        </span>
      </template>
      <template #cell-last_test_msg="{ row }">
        <div class="max-w-[240px] truncate" :title="(row as ThirdPartyAccount).last_test_msg">
          {{ (row as ThirdPartyAccount).last_test_msg || '-' }}
        </div>
      </template>
      <template #cell-last_sync_msg="{ row }">
        <div class="max-w-[240px] truncate" :title="(row as ThirdPartyAccount).last_sync_msg">
          {{ (row as ThirdPartyAccount).last_sync_msg || '-' }}
        </div>
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-nowrap items-center justify-end gap-3 whitespace-nowrap shrink-0" @click.stop>
          <BaseButton variant="link" size="sm" @click="openEdit(row as ThirdPartyAccount)">编辑</BaseButton>
          <BaseSelect v-model="actionValues[(row as ThirdPartyAccount).id]" class="w-40" placeholder="更多操作"
            :options="actionOptions"
            @update:model-value="onActionSelect(row as ThirdPartyAccount, String($event || ''))" />
        </div>
      </template>
    </BaseTable>

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑第三方账号' : '新增第三方账号'" max-width="min(720px, 96vw)">
      <div class="space-y-4 max-h-[72vh] overflow-y-auto pr-1">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
          <BaseFormItem label="平台标识" required>
            <BaseInput v-model="form.platform_name" placeholder="如 tsbeer" />
          </BaseFormItem>
          <BaseFormItem label="账号名称" required>
            <BaseInput v-model="form.name" placeholder="如 泰山ERP-账号A" />
          </BaseFormItem>
          <BaseFormItem label="登录名" required>
            <BaseInput v-model="form.login_name" />
          </BaseFormItem>
          <BaseFormItem label="手机号">
            <BaseInput v-model="form.phone" />
          </BaseFormItem>
          <BaseFormItem label="密码(加密串)" required>
            <BaseInput v-model="form.password" />
          </BaseFormItem>
          <BaseFormItem label="application-key" required>
            <BaseInput v-model="form.application_key" />
          </BaseFormItem>
          <BaseFormItem label="loginType">
            <BaseInput v-model="form.login_type" placeholder="默认 2" />
          </BaseFormItem>
          <BaseFormItem label="channel">
            <BaseInput v-model="form.channel" placeholder="默认 WEB" />
          </BaseFormItem>
          <BaseFormItem label="shopId" required>
            <BaseInput v-model="form.shop_id" placeholder="第三方店铺ID" />
          </BaseFormItem>
          <BaseFormItem label="customerId">
            <BaseInput v-model="form.customer_id" placeholder="默认与shopId一致" />
          </BaseFormItem>
        </div>
        <BaseFormItem label="启用">
          <BaseSwitch v-model="form.is_enabled" :active-value="true" :inactive-value="false" />
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
import {
  BaseButton,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseSelect,
  BaseSwitch,
  BaseTable,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import type { ThirdPartyAccount } from '@/api/types'
import {
  createThirdPartyAccount,
  deleteThirdPartyAccount,
  listThirdPartyAccounts,
  syncThirdPartyLatestOrders,
  testThirdPartyAccountLogin,
  updateThirdPartyAccount,
} from '@/api/thirdPartyAccount'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

const router = useRouter()
const qc = useQueryClient()
const keyword = ref('')
const queryKey = computed(() => ['third-party-accounts', keyword.value.trim()] as const)

const { data: rowsData, isLoading: loading } = useQuery({
  queryKey,
  queryFn: () => listThirdPartyAccounts(keyword.value.trim()),
})
const list = computed(() => rowsData.value ?? [])

const columns: BaseTableColumn[] = [
  { key: 'name', label: '账号名称', prop: 'name', minWidth: '140px', ellipsis: true },
  { key: 'login_name', label: '登录名', prop: 'login_name', minWidth: '150px', ellipsis: true },
  { key: 'phone', label: '手机号', prop: 'phone', width: '120px' },
  { key: 'last_test_ok', label: '测试状态', width: '110px' },
  { key: 'last_test_msg', label: '测试信息', minWidth: '220px' },
  { key: 'last_sync_msg', label: '同步信息', minWidth: '220px' },
  { key: 'actions', label: '操作', width: '260px', align: 'right' },
]
const actionOptions = [
  { label: '测试登录', value: 'test' },
  { label: '同步最近订单', value: 'sync' },
  { label: '查看已同步订单', value: 'orders' },
  { label: '删除', value: 'delete' },
]

function reload(): void {
  void qc.invalidateQueries({ queryKey: ['third-party-accounts'] })
}

const dlg = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const saving = ref(false)
const actionValues = reactive<Record<number, string>>({})
const form = reactive({
  platform_name: 'tsbeer',
  name: '',
  login_name: '',
  phone: '',
  password: '',
  application_key: '',
  login_type: '2',
  channel: 'WEB',
  shop_id: '',
  customer_id: '',
  is_enabled: true,
  remark: '',
})

function resetForm(): void {
  form.platform_name = 'tsbeer'
  form.name = ''
  form.login_name = ''
  form.phone = ''
  form.password = ''
  form.application_key = ''
  form.login_type = '2'
  form.channel = 'WEB'
  form.shop_id = ''
  form.customer_id = ''
  form.is_enabled = true
  form.remark = ''
}

function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  resetForm()
  dlg.value = true
}

function openEdit(row: ThirdPartyAccount): void {
  isEdit.value = true
  editId.value = row.id
  form.platform_name = row.platform_name || 'tsbeer'
  form.name = row.name || ''
  form.login_name = row.login_name || ''
  form.phone = row.phone || ''
  form.password = row.password || ''
  form.application_key = row.application_key || ''
  form.login_type = row.login_type || '2'
  form.channel = row.channel || 'WEB'
  form.shop_id = row.shop_id || ''
  form.customer_id = row.customer_id || ''
  form.is_enabled = !!row.is_enabled
  form.remark = row.remark || ''
  dlg.value = true
}

async function submit(): Promise<void> {
  if (!form.platform_name.trim() || !form.name.trim() || !form.login_name.trim() || !form.password.trim() || !form.application_key.trim() || !form.shop_id.trim()) {
    toast.warning('请填写完整的必填信息')
    return
  }
  saving.value = true
  try {
    const payload = {
      platform_name: form.platform_name.trim(),
      name: form.name.trim(),
      login_name: form.login_name.trim(),
      phone: form.phone.trim(),
      password: form.password.trim(),
      application_key: form.application_key.trim(),
      login_type: form.login_type.trim() || '2',
      channel: form.channel.trim() || 'WEB',
      shop_id: form.shop_id.trim(),
      customer_id: form.customer_id.trim(),
      is_enabled: form.is_enabled,
      remark: form.remark.trim(),
    }
    if (isEdit.value) await updateThirdPartyAccount(editId.value, payload)
    else await createThirdPartyAccount(payload)

    toast.success('已保存')
    dlg.value = false
    await qc.invalidateQueries({ queryKey: ['third-party-accounts'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function onDelete(row: ThirdPartyAccount): Promise<void> {
  const ok = await confirmDialog({ message: `确定删除账号「${row.name}」？` })
  if (!ok) return
  try {
    await deleteThirdPartyAccount(row.id)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['third-party-accounts'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

async function onTest(row: ThirdPartyAccount): Promise<void> {
  try {
    const res = await testThirdPartyAccountLogin(row.id)
    const msg = String((res as { resultMsg?: string }).resultMsg || '登录测试成功')
    toast.success(msg)
    await qc.invalidateQueries({ queryKey: ['third-party-accounts'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '登录测试失败')
    await qc.invalidateQueries({ queryKey: ['third-party-accounts'] })
  }
}

async function onSyncLatestOrders(row: ThirdPartyAccount): Promise<void> {
  try {
    const res = await syncThirdPartyLatestOrders(row.id)
    const msg = String((res as { message?: string }).message || '同步完成')
    toast.success(msg)
    await qc.invalidateQueries({ queryKey: ['third-party-accounts'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '同步失败')
    await qc.invalidateQueries({ queryKey: ['third-party-accounts'] })
  }
}

async function onActionSelect(row: ThirdPartyAccount, action: string): Promise<void> {
  if (!action) return
  try {
    if (action === 'test') await onTest(row)
    else if (action === 'sync') await onSyncLatestOrders(row)
    else if (action === 'orders') await openOrders(row)
    else if (action === 'delete') await onDelete(row)
  } finally {
    actionValues[row.id] = ''
  }
}

function openOrders(row: ThirdPartyAccount): void {
  void router.push({
    name: 'ThirdPartyAccountOrders',
    params: { id: row.id },
    query: { name: row.name || '' },
  })
}

function goRouteMap(): void {
  void router.push({ name: 'ThirdPartyRoutes' })
}

async function toggleEnabled(row: ThirdPartyAccount, value: boolean | number): Promise<void> {
  try {
    await updateThirdPartyAccount(row.id, { is_enabled: Boolean(value) })
    row.is_enabled = Boolean(value)
    toast.success('状态已更新')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '状态更新失败')
  }
}
</script>
