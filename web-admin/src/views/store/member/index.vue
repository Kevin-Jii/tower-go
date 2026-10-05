<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">会员管理</h2>
      <div class="flex flex-col sm:flex-row gap-2 w-full md:w-auto">
        <BaseInput v-model="keyword" class="w-full sm:w-48" placeholder="手机 / UID" clearable @enter="reload" />
        <BaseButton variant="primary" @click="reload">查询</BaseButton>
        <BaseButton v-permission="'marketing:sms:list'" variant="secondary" @click="openTagManager">标签管理</BaseButton>
        <BaseButton v-permission="'store:member:edit'" variant="secondary" @click="openRuleDialog">会员规则</BaseButton>
        <BaseButton v-permission="'store:member:add'" variant="primary" @click="openCreate">新增会员</BaseButton>
      </div>
    </div>

    <BaseTable :columns="columns" :data="(list as unknown) as Record<string, unknown>[]" :loading="loading"
      min-width="1000px">
      <template #cell-balance="{ row }">
        {{ fmtMoney((row as MemberRow).balance) }}
      </template>
      <template #cell-total_consumption_amount="{ row }">
        <span class="font-semibold text-emerald-700">{{ fmtMoney((row as MemberRow).total_consumption_amount ?? 0) }}</span>
      </template>
      <template #cell-unsettled_amount="{ row }">
        <span class="font-semibold text-red-600">{{ fmtMoney((row as MemberRow).unsettled_amount ?? 0) }}</span>
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="memberRowActions(row as MemberRow)" />
      </template>
    </BaseTable>

    <div class="flex justify-end">
      <BasePagination :page="page" :page-size="pageSize" :total="total" @update:page="(p) => (page = p)"
        @update:page-size="(s) => (pageSize = s)" />
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

    <BaseDialog v-model="tagManagerDlg" title="会员标签管理" max-width="min(980px, 96vw)">
      <div class="grid max-h-[68vh] grid-cols-1 gap-4 overflow-y-auto lg:grid-cols-[280px_minmax(0,1fr)]">
        <aside class="self-start rounded border border-[var(--color-border-2)] p-4">
          <h3 class="mt-0 text-base">{{ tagEditId ? '编辑标签' : '新建标签' }}</h3>
          <BaseFormItem v-if="canChooseTagStore" label="所属门店" required>
            <BaseSelect v-model="tagForm.store_id" :options="tagStoreOptions" searchable :disabled="Boolean(tagEditId)" />
          </BaseFormItem>
          <BaseFormItem label="标签名称" required>
            <BaseInput v-model="tagForm.name" placeholder="如：高频消费" maxlength="80" />
          </BaseFormItem>
          <BaseFormItem label="标识颜色">
            <div class="grid grid-cols-[42px_1fr] items-center gap-2">
              <input v-model="tagForm.color" class="h-8 w-[42px] rounded border border-[var(--color-border-2)] bg-transparent p-0.5" type="color" />
              <BaseInput v-model="tagForm.color" placeholder="#4f46e5" />
            </div>
          </BaseFormItem>
          <BaseFormItem label="说明">
            <BaseTextarea v-model="tagForm.description" :rows="3" maxlength="255" />
          </BaseFormItem>
          <div class="flex justify-end gap-2">
            <BaseButton variant="ghost" @click="resetTagForm">重置</BaseButton>
            <BaseButton v-permission="tagEditId ? 'marketing:sms:edit' : 'marketing:sms:add'" variant="primary" :loading="savingTag" @click="saveTag">保存</BaseButton>
          </div>
        </aside>

        <div class="min-w-0">
          <div v-if="tagListError" class="flex items-center justify-between gap-3 py-4 text-sm text-red-600">
            <span>{{ tagListError }}</span>
            <BaseButton variant="secondary" size="sm" @click="loadManagedTags">重试</BaseButton>
          </div>
          <BaseTable v-else :columns="tagColumns" :data="(managedTags as unknown) as Record<string, unknown>[]" :loading="tagListLoading" min-width="620px">
            <template #cell-name="{ row }">
              <span class="mr-2 inline-block h-2.5 w-2.5 rounded-full" :style="{ backgroundColor: (row as MemberTag).color || '#64748b' }"></span>
              {{ (row as MemberTag).name }}
            </template>
            <template #cell-store_id="{ row }">{{ tagStoreName((row as MemberTag).store_id) }}</template>
            <template #cell-actions="{ row }">
              <BaseTableRowActions :actions="tagActions(row as MemberTag)" :max-inline="2" />
            </template>
          </BaseTable>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="tagManagerDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="tagDlg" :title="`设置标签 · ${tagMember?.phone || ''}`" max-width="min(520px, 96vw)">
      <div class="space-y-3">
        <div v-if="tagLoading" class="py-6 text-center text-sm text-[var(--color-text-3)]">正在加载标签...</div>
        <div v-else-if="tagError" class="flex items-center justify-between gap-3 py-4 text-sm text-red-600">
          <span>{{ tagError }}</span>
          <BaseButton variant="secondary" size="sm" @click="openMemberTags(tagMember!)">重试</BaseButton>
        </div>
        <div v-else-if="!availableTags.length" class="py-6 text-center text-sm text-[var(--color-text-3)]">当前门店暂无可用标签</div>
        <template v-else>
          <div class="max-h-72 space-y-2 overflow-y-auto">
            <label v-for="tag in availableTags" :key="tag.id" class="flex cursor-pointer items-center gap-2 rounded border border-[var(--color-border-2)] px-3 py-2">
              <input v-model="selectedTagIds" type="checkbox" :value="tag.id" />
              <span class="h-2.5 w-2.5 shrink-0 rounded-full" :style="{ backgroundColor: tag.color || '#64748b' }"></span>
              <span>{{ tag.name }}</span>
            </label>
          </div>
          <p v-if="!selectedTagIds.length" class="m-0 text-xs text-[var(--color-text-3)]">未选择标签，保存后将清除该会员的全部标签。</p>
        </template>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="tagDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="tagSaving" :disabled="tagLoading || Boolean(tagError) || !tagMember" @click="saveMemberTags">保存</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="adjDlg" title="调整余额（调增/调减）" max-width="min(400px, 96vw)">
      <p class="text-sm text-slate-600 m-0 mb-3">当前余额：{{ adjMember ? fmtMoney(adjMember.balance) : '-' }}</p>
      <div class="space-y-4">
        <BaseFormItem label="类型" required>
          <BaseSelect v-model="adjForm.type" :options="[
            { label: '调增', value: 4 },
            { label: '调减', value: 5 },
          ]" />
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

    <BaseDialog v-model="consDlg" title="会员消费记录" max-width="min(960px, 96vw)">
      <div class="space-y-4">
        <div class="text-sm text-slate-600">
          会员：{{ consMember ? `${consMember.phone}${consMember.name ? `（${consMember.name}）` : ''}` : '-' }}
        </div>
        <div class="flex flex-wrap items-end gap-2">
          <BaseFormItem label="开始日期" class="w-44">
            <BaseInput v-model="consStart" type="date" />
          </BaseFormItem>
          <BaseFormItem label="结束日期" class="w-44">
            <BaseInput v-model="consEnd" type="date" />
          </BaseFormItem>
          <BaseButton variant="primary" :loading="consLoading" @click="reloadConsumptions">查询</BaseButton>
          <BaseButton variant="secondary" @click="openConsExportDlg">导出Excel</BaseButton>
        </div>
        <div class="grid grid-cols-1 sm:grid-cols-4 gap-3 rounded border border-[var(--color-border-2)] p-3">
          <div>
            <div class="text-xs text-[var(--color-text-3)]">消费笔数</div>
            <div class="text-base font-semibold">{{ consSummary.count ?? 0 }}</div>
          </div>
          <div>
            <div class="text-xs text-[var(--color-text-3)]">销售总额</div>
            <div class="text-base font-semibold text-indigo-700">{{ fmtMoney(consSummary.total_amount ?? 0) }}</div>
          </div>
          <div>
            <div class="text-xs text-[var(--color-text-3)]">消耗品成本</div>
            <div class="text-base font-semibold text-amber-700">{{ fmtMoney(consSummary.consumable_amount ?? 0) }}</div>
          </div>
          <div>
            <div class="text-xs text-[var(--color-text-3)]">净利润</div>
            <div class="text-base font-semibold text-emerald-700">{{ fmtMoney(consSummary.net_income_amount ?? 0) }}
            </div>
          </div>
        </div>
        <BaseTable :columns="consColumns" :data="(consRows as unknown) as Record<string, unknown>[]"
          :loading="consLoading" min-width="860px">
          <template #cell-total_amount="{ row }">{{ fmtMoney((row as MemberConsumptionRecord).total_amount)
          }}</template>
          <template #cell-consumable_amount="{ row }">{{ fmtMoney((row as MemberConsumptionRecord).consumable_amount)
          }}</template>
          <template #cell-net_income_amount="{ row }">{{ fmtMoney((row as MemberConsumptionRecord).net_income_amount)
          }}</template>
        </BaseTable>
        <div class="flex justify-end">
          <BasePagination :page="consPage" :page-size="consPageSize" :total="consTotal"
            @update:page="(p) => (consPage = p)" @update:page-size="(s) => (consPageSize = s)" />
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="consDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="consExportDlg" title="导出会员消费" max-width="min(420px, 96vw)">
      <BaseFormItem label="导出日期" required>
        <a-date-picker v-model="consExportDate" value-format="YYYY-MM-DD" class="w-full" :allow-clear="false" />
      </BaseFormItem>
      <template #footer>
        <BaseButton variant="ghost" @click="consExportDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="consExporting" @click="submitConsExport">导出</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="giftDlg" title="会员赠品记录" max-width="min(960px, 96vw)">
      <div class="space-y-4">
        <div class="text-sm text-slate-600">
          会员：{{ giftMember ? `${giftMember.phone}${giftMember.name ? `（${giftMember.name}）` : ''}` : '-' }}
        </div>
        <BaseTable
          :columns="giftColumns"
          :data="(giftRows as unknown) as Record<string, unknown>[]"
          :loading="giftLoading"
          min-width="860px"
        >
          <template #cell-cost_amount="{ row }">{{ fmtMoney((row as MemberGiftRecord).cost_amount) }}</template>
        </BaseTable>
        <div class="flex justify-end">
          <BasePagination
            :page="giftPage"
            :page-size="giftPageSize"
            :total="giftTotal"
            @update:page="(p) => (giftPage = p)"
            @update:page-size="(s) => (giftPageSize = s)"
          />
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="giftDlg = false">关闭</BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="ruleDlg" title="会员规则" max-width="min(900px, 96vw)">
      <div class="space-y-4">
        <div class="flex flex-col lg:flex-row gap-4">
          <div class="lg:w-72 shrink-0 rounded border border-[var(--color-border-2)] p-3">
            <div class="mb-3 text-sm font-semibold">{{ ruleEditId ? '编辑规则' : '新增规则' }}</div>
            <div class="space-y-3">
              <BaseFormItem label="规则名称" required>
                <BaseInput v-model="ruleForm.name" placeholder="如：消费积分" />
              </BaseFormItem>
              <BaseFormItem label="消费金额" required>
                <BaseNumberInput v-model="ruleForm.spend_amount" :min="0.01" :step="0.01" />
              </BaseFormItem>
              <BaseFormItem label="获得积分" required>
                <BaseNumberInput v-model="ruleForm.points" :min="1" :step="1" />
              </BaseFormItem>
              <BaseFormItem label="状态">
                <BaseSelect v-model="ruleForm.status" :options="ruleStatusOptions" />
              </BaseFormItem>
              <BaseFormItem label="备注">
                <BaseInput v-model="ruleForm.remark" />
              </BaseFormItem>
              <div class="flex justify-end gap-2">
                <BaseButton variant="ghost" @click="resetRuleForm">重置</BaseButton>
                <BaseButton variant="primary" :loading="ruleSaving" @click="saveRule">保存</BaseButton>
              </div>
            </div>
          </div>

          <div class="min-w-0 flex-1 space-y-3">
            <div class="flex flex-wrap items-end gap-2">
              <BaseInput v-model="ruleKeyword" class="w-full sm:w-48" placeholder="规则名称 / 备注" clearable @enter="reloadRules" />
              <BaseButton variant="primary" :loading="ruleLoading" @click="reloadRules">查询</BaseButton>
            </div>
            <BaseTable
              :columns="ruleColumns"
              :data="(ruleRows as unknown) as Record<string, unknown>[]"
              :loading="ruleLoading"
              min-width="620px"
            >
              <template #cell-spend_amount="{ row }">{{ fmtMoney((row as MemberPointRule).spend_amount) }}</template>
              <template #cell-points="{ row }">{{ (row as MemberPointRule).points }}</template>
              <template #cell-status="{ row }">
                <span :class="(row as MemberPointRule).status === 1 ? 'text-emerald-600' : 'text-slate-500'">
                  {{ ruleStatusText((row as MemberPointRule).status) }}
                </span>
              </template>
              <template #cell-actions="{ row }">
                <BaseTableRowActions :actions="ruleRowActions(row as MemberPointRule)" :max-inline="2" />
              </template>
            </BaseTable>
            <div class="flex justify-end">
              <BasePagination
                :page="rulePage"
                :page-size="rulePageSize"
                :total="ruleTotal"
                @update:page="(p) => (rulePage = p)"
                @update:page-size="(s) => (rulePageSize = s)"
              />
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="ruleDlg = false">关闭</BaseButton>
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
  BaseTableRowActions,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn, TableRowAction } from '@/components/base/types'
import {
  assignMemberTags,
  createMemberTag,
  deleteMemberTag,
  listMemberTags,
  listMemberTagsForMember,
  updateMemberTag,
} from '@/api/smsPromotion'
import {
  adjustMemberBalance,
  createMember,
  createMemberPointRule,
  deleteMember,
  deleteMemberPointRule,
  exportMemberConsumptions,
  listMemberConsumptions,
  listMemberGiftRecords,
  listMemberPointRules,
  listMembers,
  updateMember,
  updateMemberPointRule,
} from '@/api/member'
import type { MemberConsumptionRecord, MemberGiftRecord, MemberPointRule, MemberRow, MemberTag, Store } from '@/api/types'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'
import { usePermission } from '@/hooks/usePermission'
import { listAllStores } from '@/api/store'
import { useUserStore } from '@/store/user'

const qc = useQueryClient()
const { hasPerm } = usePermission()
const userStore = useUserStore()
const currentStoreId = computed(() => Number(userStore.currentStoreId || 0))
const canChooseTagStore = computed(() => currentStoreId.value === 0)
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
  { key: 'phone', label: '手机', prop: 'phone', width: '140px' },
  { key: 'name', label: '姓名', prop: 'name', width: '100px' },
  { key: 'balance', label: '余额', width: '100px' },
  { key: 'total_consumption_amount', label: '累计消费', width: '120px' },
  { key: 'unsettled_amount', label: '未结算', width: '110px' },
  { key: 'points', label: '积分', prop: 'points', width: '72px' },
  { key: 'level', label: '等级', prop: 'level', width: '72px' },
  { key: 'actions', label: '操作', width: '240px', align: 'right' },
]
const consColumns: BaseTableColumn[] = [
  { key: 'account_no', label: '记账单号', prop: 'account_no', minWidth: '140px', ellipsis: true },
  { key: 'account_date', label: '日期', prop: 'account_date', width: '120px' },
  { key: 'channel_name', label: '渠道', prop: 'channel_name', width: '120px' },
  { key: 'order_no', label: '订单号', prop: 'order_no', minWidth: '120px', ellipsis: true },
  { key: 'total_amount', label: '销售额', width: '100px' },
  { key: 'consumable_amount', label: '消耗品', width: '100px' },
  { key: 'net_income_amount', label: '净利润', width: '100px' },
]
const giftColumns: BaseTableColumn[] = [
  { key: 'created_at', label: '日期', prop: 'created_at', width: '170px' },
  { key: 'order_no', label: '单号', prop: 'order_no', minWidth: '150px', ellipsis: true },
  { key: 'product_name', label: '商品', prop: 'product_name', minWidth: '160px', ellipsis: true },
  { key: 'unit', label: '规格', prop: 'unit', width: '110px' },
  { key: 'quantity', label: '数量', prop: 'quantity', width: '90px' },
  { key: 'cost_amount', label: '成本金额', width: '110px' },
  { key: 'reason', label: '原因', prop: 'reason', minWidth: '180px', ellipsis: true },
  { key: 'operator_name', label: '操作人', prop: 'operator_name', width: '100px' },
]
const tagColumns: BaseTableColumn[] = [
  { key: 'name', label: '标签', minWidth: '150px' },
  { key: 'store_id', label: '门店', width: '130px' },
  { key: 'member_count', label: '会员数', prop: 'member_count', width: '90px' },
  { key: 'description', label: '说明', prop: 'description', minWidth: '160px', ellipsis: true },
  { key: 'actions', label: '操作', width: '150px', align: 'right' },
]
const ruleColumns: BaseTableColumn[] = [
  { key: 'name', label: '规则名称', prop: 'name', minWidth: '140px', ellipsis: true },
  { key: 'spend_amount', label: '消费金额', width: '100px' },
  { key: 'points', label: '积分', width: '80px' },
  { key: 'status', label: '状态', width: '80px' },
  { key: 'remark', label: '备注', prop: 'remark', minWidth: '140px', ellipsis: true },
  { key: 'actions', label: '操作', width: '120px', align: 'right' },
]
const ruleStatusOptions = [
  { label: '启用', value: 1 },
  { label: '停用', value: 2 },
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

function memberRowActions(row: MemberRow): TableRowAction[] {
  const actions: TableRowAction[] = [
    { label: '消费记录', permission: 'store:member:list', onClick: () => openConsumptions(row) },
    { label: '赠品记录', permission: 'store:member:list', onClick: () => openGifts(row) },
    { label: '编辑', permission: 'store:member:edit', onClick: () => openEdit(row) },
    { label: '调余额', permission: 'store:member:balance', onClick: () => openAdjust(row) },
    { label: '删除', permission: 'store:member:delete', danger: true, onClick: () => void onDelete(row) },
  ]
  if (hasPerm('marketing:sms:list') && hasPerm('marketing:sms:edit')) {
    actions.splice(2, 0, { label: '设置标签', permission: 'marketing:sms:edit', onClick: () => void openMemberTags(row) })
  }
  return actions
}

const { data: tagStoreData } = useQuery({
  queryKey: ['stores', 'member-tags'],
  queryFn: async () => {
    try {
      return await listAllStores()
    } catch {
      return [] as Store[]
    }
  },
})
const tagStores = computed(() => {
  const rows = [...(tagStoreData.value ?? [])]
  const own = userStore.userInfo?.store
  if (currentStoreId.value > 0 && !rows.some((store) => store.id === currentStoreId.value)) {
    rows.push({ id: currentStoreId.value, name: own?.name || `当前门店 #${currentStoreId.value}` })
  }
  return rows
})
const tagStoreOptions = computed(() => tagStores.value.map((store) => ({ label: store.name, value: store.id })))

const tagManagerDlg = ref(false)
const tagListLoading = ref(false)
const tagListError = ref('')
const managedTags = ref<MemberTag[]>([])
const tagEditId = ref(0)
const savingTag = ref(false)
const tagForm = reactive({ store_id: 0, name: '', color: '#4f46e5', description: '' })

function tagStoreName(id: number): string {
  return tagStores.value.find((store) => store.id === id)?.name || `门店 #${id}`
}

function resetTagForm(): void {
  tagEditId.value = 0
  tagForm.store_id = currentStoreId.value || tagStores.value[0]?.id || 0
  tagForm.name = ''
  tagForm.color = '#4f46e5'
  tagForm.description = ''
}

async function loadManagedTags(): Promise<void> {
  tagListLoading.value = true
  tagListError.value = ''
  try {
    managedTags.value = await listMemberTags(currentStoreId.value > 0 ? { store_id: currentStoreId.value } : undefined)
  } catch (error: unknown) {
    tagListError.value = error instanceof Error ? error.message : '加载标签失败，请重试'
  } finally {
    tagListLoading.value = false
  }
}

function openTagManager(): void {
  resetTagForm()
  tagManagerDlg.value = true
  void loadManagedTags()
}

function editTag(tag: MemberTag): void {
  tagEditId.value = tag.id
  tagForm.store_id = tag.store_id
  tagForm.name = tag.name
  tagForm.color = tag.color || '#4f46e5'
  tagForm.description = tag.description || ''
}

async function saveTag(): Promise<void> {
  if (!tagForm.store_id) return toast.warning('请选择标签所属门店')
  if (!tagForm.name.trim()) return toast.warning('请填写标签名称')
  savingTag.value = true
  try {
    const body = { name: tagForm.name.trim(), color: tagForm.color.trim(), description: tagForm.description.trim() }
    if (tagEditId.value) await updateMemberTag(tagEditId.value, body, tagForm.store_id)
    else await createMemberTag({ store_id: tagForm.store_id, ...body })
    toast.success('标签已保存')
    resetTagForm()
    await loadManagedTags()
    await qc.invalidateQueries({ queryKey: ['member-tags'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '保存标签失败')
  } finally {
    savingTag.value = false
  }
}

async function removeTag(tag: MemberTag): Promise<void> {
  const ok = await confirmDialog({ message: `删除标签「${tag.name}」？会员绑定关系也会被清除。` })
  if (!ok) return
  try {
    await deleteMemberTag(tag.id, tag.store_id)
    toast.success('已删除')
    if (tagEditId.value === tag.id) resetTagForm()
    await loadManagedTags()
    await qc.invalidateQueries({ queryKey: ['member-tags'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '删除标签失败')
  }
}

function tagActions(tag: MemberTag): TableRowAction[] {
  return [
    { label: '编辑', permission: 'marketing:sms:edit', onClick: () => editTag(tag), place: 'inline' },
    { label: '删除', permission: 'marketing:sms:delete', danger: true, onClick: () => void removeTag(tag), place: 'inline' },
  ]
}

const tagDlg = ref(false)
const tagLoading = ref(false)
const tagSaving = ref(false)
const tagError = ref('')
const tagMember = ref<MemberRow | null>(null)
const availableTags = ref<MemberTag[]>([])
const selectedTagIds = ref<number[]>([])

async function openMemberTags(row: MemberRow): Promise<void> {
  tagMember.value = row
  availableTags.value = []
  selectedTagIds.value = []
  tagError.value = ''
  tagDlg.value = true
  tagLoading.value = true
  try {
    const [tags, assigned] = await Promise.all([
      listMemberTags(row.store_id ? { store_id: row.store_id } : currentStoreId.value > 0 ? { store_id: currentStoreId.value } : undefined),
      listMemberTagsForMember(row.id),
    ])
    availableTags.value = tags
    const availableIds = new Set(tags.map((tag) => tag.id))
    selectedTagIds.value = assigned.map((tag) => tag.id).filter((id) => availableIds.has(id))
  } catch (error: unknown) {
    tagError.value = error instanceof Error ? error.message : '加载标签失败，请重试'
  } finally {
    tagLoading.value = false
  }
}

async function saveMemberTags(): Promise<void> {
  if (!tagMember.value || tagLoading.value || tagError.value) return
  tagSaving.value = true
  try {
    await assignMemberTags(tagMember.value.id, selectedTagIds.value)
    toast.success('标签已更新')
    tagDlg.value = false
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '保存标签失败')
  } finally {
    tagSaving.value = false
  }
}

const consDlg = ref(false)
const consLoading = ref(false)
const consMember = ref<MemberRow | null>(null)
const consStart = ref('')
const consEnd = ref('')
const consExportDlg = ref(false)
const consExportDate = ref(new Date().toISOString().slice(0, 10))
const consExporting = ref(false)
const consPage = ref(1)
const consPageSize = ref(10)
const consTotal = ref(0)
const consRows = ref<MemberConsumptionRecord[]>([])
const consSummary = reactive({
  count: 0,
  total_amount: 0,
  other_expense_amount: 0,
  consumable_amount: 0,
  net_income_amount: 0,
})

function currentMonthRange(): { start: string; end: string } {
  const t = new Date()
  const y = t.getFullYear()
  const m = String(t.getMonth() + 1).padStart(2, '0')
  const d = String(t.getDate()).padStart(2, '0')
  return { start: `${y}-${m}-01`, end: `${y}-${m}-${d}` }
}

async function loadConsumptions(): Promise<void> {
  if (!consMember.value) return
  consLoading.value = true
  try {
    const data = await listMemberConsumptions(consMember.value.id, {
      start_date: consStart.value || undefined,
      end_date: consEnd.value || undefined,
      page: consPage.value,
      page_size: consPageSize.value,
    })
    consRows.value = data.list ?? []
    consTotal.value = Number(data.total ?? 0)
    consSummary.count = Number(data.summary?.count ?? 0)
    consSummary.total_amount = Number(data.summary?.total_amount ?? 0)
    consSummary.other_expense_amount = Number(data.summary?.other_expense_amount ?? 0)
    consSummary.consumable_amount = Number(data.summary?.consumable_amount ?? 0)
    consSummary.net_income_amount = Number(data.summary?.net_income_amount ?? 0)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载消费记录失败')
  } finally {
    consLoading.value = false
  }
}

function reloadConsumptions(): void {
  consPage.value = 1
  void loadConsumptions()
}

function openConsExportDlg(): void {
  consExportDate.value = consStart.value || consEnd.value || new Date().toISOString().slice(0, 10)
  consExportDlg.value = true
}

async function submitConsExport(): Promise<void> {
  if (!consMember.value) return
  if (!consExportDate.value) {
    toast.warning('请选择导出日期')
    return
  }
  consExporting.value = true
  try {
    await exportMemberConsumptions(consMember.value.id, { date: consExportDate.value })
    consExportDlg.value = false
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '导出失败')
  } finally {
    consExporting.value = false
  }
}

function openConsumptions(row: MemberRow): void {
  consMember.value = row
  const r = currentMonthRange()
  consStart.value = r.start
  consEnd.value = r.end
  consPage.value = 1
  consPageSize.value = 10
  consTotal.value = 0
  consRows.value = []
  consSummary.count = 0
  consSummary.total_amount = 0
  consSummary.other_expense_amount = 0
  consSummary.consumable_amount = 0
  consSummary.net_income_amount = 0
  consDlg.value = true
  void loadConsumptions()
}

watch([consPage, consPageSize], () => {
  if (consDlg.value) {
    void loadConsumptions()
  }
})

const giftDlg = ref(false)
const giftLoading = ref(false)
const giftMember = ref<MemberRow | null>(null)
const giftPage = ref(1)
const giftPageSize = ref(10)
const giftTotal = ref(0)
const giftRows = ref<MemberGiftRecord[]>([])

async function loadGifts(): Promise<void> {
  if (!giftMember.value) return
  giftLoading.value = true
  try {
    const data = await listMemberGiftRecords(giftMember.value.id, {
      page: giftPage.value,
      page_size: giftPageSize.value,
    })
    giftRows.value = data.list ?? []
    giftTotal.value = Number(data.total ?? 0)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载赠品记录失败')
  } finally {
    giftLoading.value = false
  }
}

function openGifts(row: MemberRow): void {
  giftMember.value = row
  giftPage.value = 1
  giftPageSize.value = 10
  giftTotal.value = 0
  giftRows.value = []
  giftDlg.value = true
  void loadGifts()
}

watch([giftPage, giftPageSize], () => {
  if (giftDlg.value) {
    void loadGifts()
  }
})

const ruleDlg = ref(false)
const ruleLoading = ref(false)
const ruleSaving = ref(false)
const ruleKeyword = ref('')
const rulePage = ref(1)
const rulePageSize = ref(10)
const ruleTotal = ref(0)
const ruleRows = ref<MemberPointRule[]>([])
const ruleEditId = ref(0)
const ruleForm = reactive({
  name: '',
  spend_amount: 1,
  points: 1,
  status: 1,
  remark: '',
})

function ruleStatusText(status: number): string {
  return status === 1 ? '启用' : '停用'
}

function resetRuleForm(): void {
  ruleEditId.value = 0
  ruleForm.name = ''
  ruleForm.spend_amount = 1
  ruleForm.points = 1
  ruleForm.status = 1
  ruleForm.remark = ''
}

async function loadRules(): Promise<void> {
  ruleLoading.value = true
  try {
    const data = await listMemberPointRules({
      page: rulePage.value,
      page_size: rulePageSize.value,
      keyword: ruleKeyword.value.trim() || undefined,
    })
    ruleRows.value = data.list ?? []
    ruleTotal.value = Number(data.total ?? 0)
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载会员规则失败')
  } finally {
    ruleLoading.value = false
  }
}

function reloadRules(): void {
  rulePage.value = 1
  void loadRules()
}

function openRuleDialog(): void {
  ruleDlg.value = true
  resetRuleForm()
  rulePage.value = 1
  rulePageSize.value = 10
  ruleKeyword.value = ''
  void loadRules()
}

function editRule(row: MemberPointRule): void {
  ruleEditId.value = row.id
  ruleForm.name = row.name ?? ''
  ruleForm.spend_amount = Number(row.spend_amount ?? 1)
  ruleForm.points = Number(row.points ?? 1)
  ruleForm.status = row.status ?? 1
  ruleForm.remark = row.remark ?? ''
}

async function saveRule(): Promise<void> {
  if (!ruleForm.name.trim()) {
    toast.warning('请填写规则名称')
    return
  }
  if (!ruleForm.spend_amount || ruleForm.spend_amount <= 0) {
    toast.warning('消费金额必须大于0')
    return
  }
  if (!ruleForm.points || ruleForm.points <= 0) {
    toast.warning('积分必须大于0')
    return
  }
  ruleSaving.value = true
  try {
    const body = {
      name: ruleForm.name.trim(),
      spend_amount: Number(ruleForm.spend_amount),
      points: Number(ruleForm.points),
      status: Number(ruleForm.status),
      remark: ruleForm.remark.trim() || undefined,
    }
    if (ruleEditId.value) {
      await updateMemberPointRule(ruleEditId.value, body)
    } else {
      await createMemberPointRule(body)
    }
    toast.success('已保存')
    resetRuleForm()
    await loadRules()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '保存会员规则失败')
  } finally {
    ruleSaving.value = false
  }
}

async function removeRule(row: MemberPointRule): Promise<void> {
  const ok = await confirmDialog({ message: `删除规则「${row.name}」？` })
  if (!ok) return
  try {
    await deleteMemberPointRule(row.id)
    toast.success('已删除')
    if (ruleEditId.value === row.id) resetRuleForm()
    await loadRules()
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除会员规则失败')
  }
}

function ruleRowActions(row: MemberPointRule): TableRowAction[] {
  return [
    { label: '编辑', permission: 'store:member:edit', onClick: () => editRule(row) },
    { label: '删除', permission: 'store:member:edit', danger: true, onClick: () => void removeRule(row) },
  ]
}

watch([rulePage, rulePageSize], () => {
  if (ruleDlg.value) {
    void loadRules()
  }
})
</script>
