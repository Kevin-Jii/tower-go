<template>
  <div class="sms-page">
    <div class="page-head">
      <div>
        <h2 class="page-title">会员推广</h2>
        <p class="page-subtitle">按会员标签分组发送阿里云短信，支持立即发送与定时排期。</p>
      </div>
      <div class="head-actions">
        <BaseButton v-permission="['marketing:sms:add', 'marketing:sms:edit']" variant="secondary" @click="openTagManager">
          标签管理
        </BaseButton>
        <BaseButton v-permission="'marketing:sms:add'" variant="primary" @click="openCreate">
          新建推广
        </BaseButton>
      </div>
    </div>

    <a-alert v-if="config && (!config.configured || !config.enabled)" type="warning" show-icon>
      <template #title>短信服务尚未可用</template>
      AccessKey {{ config.configured ? '已配置' : '未配置' }}，发送开关{{ config.enabled ? '已开启' : '未开启' }}。
      可先维护标签和活动，但立即发送或排期执行会失败。排期仍按 {{ config.timezone || 'Asia/Shanghai' }}（UTC+8）解释，
      允许窗口为 {{ config.send_window_start }}–{{ config.send_window_end }}（结束时刻不含）。
      <a :href="config.help_url" target="_blank" rel="noreferrer">查看阿里云短信配置指南</a>
    </a-alert>
    <a-alert v-else-if="config" type="info" show-icon>
      <template #title>发送时间以中国标准时间为准</template>
      时区 {{ config.timezone || 'Asia/Shanghai' }}（UTC+8），允许发送窗口为
      {{ config.send_window_start }}–{{ config.send_window_end }}，结束时刻{{ config.send_window_end_exclusive ? '不包含' : '包含' }}。
      阿里云接口不负责排期，系统将在计划时间到达后执行；请使用审核通过的签名和模板 CODE。
    </a-alert>

    <div class="stats-grid">
      <BaseCard v-for="item in stats" :key="item.label" body-padding="16px">
        <div class="stat-label">{{ item.label }}</div>
        <div class="stat-value" :class="item.tone">{{ item.value }}</div>
      </BaseCard>
    </div>

    <BaseCard title="推广活动" body-padding="0">
      <BaseTable
        :columns="campaignColumns"
        :data="(campaigns as unknown) as Record<string, unknown>[]"
        :loading="campaignLoading"
        min-width="1180px"
      >
        <template #cell-name="{ row }">
          <div class="font-medium">{{ (row as SmsCampaign).name }}</div>
          <div class="cell-secondary">{{ campaignTypeText((row as SmsCampaign).campaign_type) }}</div>
        </template>
        <template #cell-owner_store_id="{ row }">
          {{ storeName((row as SmsCampaign).owner_store_id) }}
        </template>
        <template #cell-target_type="{ row }">
          {{ targetText(row as SmsCampaign) }}
        </template>
        <template #cell-segments="{ row }">
          {{ (row as SmsCampaign).segments?.length || ((row as SmsCampaign).template_code ? 1 : 0) }} 组
        </template>
        <template #cell-status="{ row }">
          <BaseTag :variant="statusMeta((row as SmsCampaign).status).variant">
            {{ statusMeta((row as SmsCampaign).status).label }}
          </BaseTag>
        </template>
        <template #cell-counts="{ row }">
          <span class="text-emerald-600">{{ (row as SmsCampaign).success_count || 0 }}</span>
          /
          <span class="text-red-600">{{ (row as SmsCampaign).fail_count || 0 }}</span>
          /
          {{ (row as SmsCampaign).total_count || 0 }}
        </template>
        <template #cell-scheduled_at="{ row }">
          {{ formatChinaTime((row as SmsCampaign).scheduled_at || (row as SmsCampaign).sent_at) }}
        </template>
        <template #cell-actions="{ row }">
          <BaseTableRowActions :actions="campaignActions(row as SmsCampaign)" :max-inline="3" />
        </template>
      </BaseTable>
    </BaseCard>

    <BaseDialog v-model="campaignDlg" :title="editingId ? '编辑推广活动' : '新建推广活动'" max-width="min(980px, 96vw)">
      <div class="campaign-form">
        <section class="form-section">
          <h3>基础信息</h3>
          <div class="form-grid">
            <BaseFormItem label="活动名称" required>
              <BaseInput v-model="form.name" placeholder="如：中秋会员关怀" maxlength="120" />
            </BaseFormItem>
            <BaseFormItem label="活动类型" required>
              <BaseSelect v-model="form.campaign_type" :options="campaignTypeOptions" />
            </BaseFormItem>
            <BaseFormItem v-if="canChooseOwner" label="归属门店" required hint="选择总部全局时可跨门店选受众和标签">
              <BaseSelect v-model="form.owner_store_id" :options="ownerStoreOptions" searchable />
            </BaseFormItem>
            <BaseFormItem label="目标范围" required>
              <BaseSelect v-model="form.target_type" :options="targetOptions" />
            </BaseFormItem>
          </div>

          <BaseFormItem v-if="form.target_type === 'stores'" label="选择门店" required>
            <div class="choice-box">
              <label v-for="store in targetStoreChoices" :key="store.id" class="check-item">
                <input v-model="form.store_ids" type="checkbox" :value="store.id" />
                <span>{{ store.name }}</span>
              </label>
              <span v-if="!targetStoreChoices.length" class="cell-secondary">暂无可选门店</span>
            </div>
          </BaseFormItem>
          <BaseFormItem v-if="form.target_type === 'custom'" label="自定义手机号" required hint="支持换行、逗号、分号或空格分隔；自动去重">
            <BaseTextarea v-model="form.custom_phones" :rows="4" placeholder="13800138000&#10;13900139000" />
          </BaseFormItem>
        </section>

        <section class="form-section">
          <div class="section-head">
            <div>
              <h3>标签分组与短信模板</h3>
              <p>按从上到下顺序匹配，会员命中多个分组时使用第一个；默认分组接收未命中标签的会员。</p>
            </div>
            <BaseButton variant="secondary" size="sm" @click="addSegment(false)">添加标签分组</BaseButton>
          </div>

          <a-alert v-if="form.target_type === 'custom'" type="warning" class="mb-3">
            自定义手机号没有会员标签信息，只会使用默认分组。
          </a-alert>

          <div v-if="form.segments.length" class="segment-list">
            <article v-for="(segment, index) in form.segments" :key="segment.key" class="segment-card">
              <div class="segment-head">
                <div class="segment-title">
                  <span class="segment-index">{{ index + 1 }}</span>
                  <strong>{{ segment.default ? '默认分组' : '标签分组' }}</strong>
                </div>
                <div class="segment-actions">
                  <BaseButton variant="ghost" size="sm" :disabled="index === 0" @click="moveSegment(index, -1)">上移</BaseButton>
                  <BaseButton variant="ghost" size="sm" :disabled="index === form.segments.length - 1" @click="moveSegment(index, 1)">下移</BaseButton>
                  <BaseButton variant="danger" size="sm" :disabled="form.segments.length === 1" @click="removeSegment(index)">删除</BaseButton>
                </div>
              </div>

              <label class="default-switch">
                <a-switch :model-value="segment.default" @update:model-value="setDefault(index, Boolean($event))" />
                <span>设为默认分组</span>
              </label>

              <BaseFormItem v-if="!segment.default" label="匹配标签（任一命中）" required>
                <div class="tag-choice-box">
                  <label v-for="tag in campaignTags" :key="tag.id" class="tag-check" :style="tagCheckStyle(tag)">
                    <input v-model="segment.tag_ids" type="checkbox" :value="tag.id" />
                    <span>{{ tag.name }}（{{ tag.member_count || 0 }}）</span>
                  </label>
                  <span v-if="!campaignTags.length" class="cell-secondary">
                    当前归属门店暂无标签，请先打开“标签管理”创建。
                  </span>
                </div>
              </BaseFormItem>

              <div class="segment-grid">
                <BaseFormItem label="签名" hint="留空使用系统默认签名">
                  <BaseInput v-model="segment.sign" :placeholder="config?.default_sign_name || '系统默认签名'" maxlength="64" />
                </BaseFormItem>
                <BaseFormItem label="模板 CODE" required>
                  <BaseInput v-model="segment.template" placeholder="SMS_123456789" maxlength="64" />
                </BaseFormItem>
              </div>
              <BaseFormItem label="模板变量 JSON" hint='例如 {"activity":"会员日"}；变量须与阿里云模板一致'>
                <BaseTextarea v-model="segment.params" :rows="3" placeholder='{"activity":"会员日"}' />
              </BaseFormItem>
              <label class="default-switch">
                <a-switch v-model="segment.personalize" />
                <span>姓名个性化（自动写入 name 变量；姓名为空时使用“会员”）</span>
              </label>
            </article>
          </div>
        </section>

        <section class="form-section">
          <h3>发送方式</h3>
          <a-radio-group v-model="form.send_mode" type="button">
            <a-radio value="draft">保存草稿</a-radio>
            <a-radio value="now" :disabled="!canSendNow">立即发送</a-radio>
            <a-radio value="scheduled" :disabled="!canSendNow">定时发送</a-radio>
          </a-radio-group>
          <p v-if="!canSendNow" class="permission-hint">当前账号没有短信发送权限，仅可保存草稿。</p>
          <BaseFormItem
            v-if="form.send_mode === 'scheduled'"
            label="计划发送时间（中国标准时间 UTC+8）"
            required
            :hint="`只允许 ${config?.send_window_start || '08:00'}–${config?.send_window_end || '22:00'}，结束时刻不含`"
            class="schedule-field"
          >
            <a-date-picker
              v-model="form.scheduled_at"
              show-time
              value-format="YYYY-MM-DD HH:mm"
              format="YYYY-MM-DD HH:mm"
              :allow-clear="false"
              class="w-full"
            />
          </BaseFormItem>
        </section>
      </div>
      <template #footer>
        <BaseButton variant="ghost" @click="campaignDlg = false">取消</BaseButton>
        <BaseButton variant="primary" :loading="savingCampaign" @click="submitCampaign">
          {{ form.send_mode === 'now' ? '保存并立即发送' : form.send_mode === 'scheduled' ? '保存排期' : '保存草稿' }}
        </BaseButton>
      </template>
    </BaseDialog>

    <BaseDialog v-model="recordsDlg" :title="`发送记录${recordsCampaign ? ` · ${recordsCampaign.name}` : ''}`" max-width="min(980px, 96vw)">
      <div class="records-summary">
        <span>总计 {{ recordsCampaign?.total_count || records.length }}</span>
        <span class="text-emerald-600">成功 {{ recordsCampaign?.success_count || 0 }}</span>
        <span class="text-red-600">失败 {{ recordsCampaign?.fail_count || 0 }}</span>
      </div>
      <BaseTable
        :columns="recordColumns"
        :data="(records as unknown) as Record<string, unknown>[]"
        :loading="recordsLoading"
        min-width="780px"
      >
        <template #cell-segment_id="{ row }">{{ recordSegmentName(row as SmsSendRecord) }}</template>
        <template #cell-status="{ row }">
          <BaseTag :variant="(row as SmsSendRecord).status === 'success' ? 'success' : 'danger'">
            {{ (row as SmsSendRecord).status === 'success' ? '成功' : '失败' }}
          </BaseTag>
        </template>
        <template #cell-created_at="{ row }">{{ formatChinaTime((row as SmsSendRecord).created_at) }}</template>
      </BaseTable>
      <template #footer><BaseButton variant="ghost" @click="recordsDlg = false">关闭</BaseButton></template>
    </BaseDialog>

    <BaseDialog v-model="tagDlg" title="会员标签管理" max-width="min(980px, 96vw)">
      <div class="tag-manager">
        <aside class="tag-editor">
          <h3>{{ tagEditId ? '编辑标签' : '新建标签' }}</h3>
          <BaseFormItem v-if="canChooseOwner" label="所属门店" required>
            <BaseSelect v-model="tagForm.store_id" :options="storeOnlyOptions" searchable :disabled="Boolean(tagEditId)" />
          </BaseFormItem>
          <BaseFormItem label="标签名称" required>
            <BaseInput v-model="tagForm.name" placeholder="如：高频消费" maxlength="80" />
          </BaseFormItem>
          <BaseFormItem label="标识颜色">
            <div class="color-field">
              <input v-model="tagForm.color" type="color" />
              <BaseInput v-model="tagForm.color" placeholder="#4f46e5" />
            </div>
          </BaseFormItem>
          <BaseFormItem label="说明">
            <BaseTextarea v-model="tagForm.description" :rows="3" maxlength="255" />
          </BaseFormItem>
          <div class="editor-actions">
            <BaseButton variant="ghost" @click="resetTagForm">重置</BaseButton>
            <BaseButton v-permission="tagEditId ? 'marketing:sms:edit' : 'marketing:sms:add'" variant="primary" :loading="savingTag" @click="saveTag">保存</BaseButton>
          </div>
        </aside>

        <div class="tag-list-wrap">
          <BaseTable :columns="tagColumns" :data="(visibleTags as unknown) as Record<string, unknown>[]" :loading="tagLoading" min-width="640px">
            <template #cell-name="{ row }">
              <span class="color-dot" :style="{ backgroundColor: (row as MemberTag).color || '#64748b' }"></span>
              {{ (row as MemberTag).name }}
            </template>
            <template #cell-store_id="{ row }">{{ storeName((row as MemberTag).store_id) }}</template>
            <template #cell-actions="{ row }">
              <BaseTableRowActions :actions="tagActions(row as MemberTag)" :max-inline="3" />
            </template>
          </BaseTable>
        </div>
      </div>
      <template #footer><BaseButton variant="ghost" @click="tagDlg = false">关闭</BaseButton></template>
    </BaseDialog>

    <BaseDialog v-model="membersDlg" :title="`标签成员 · ${activeTag?.name || ''}`" max-width="min(900px, 96vw)">
      <div class="member-assignment">
        <div class="member-search">
          <BaseInput v-model="memberKeyword" placeholder="手机号 / UID" clearable @enter="searchMembers" />
          <BaseButton variant="primary" :loading="memberSearchLoading" @click="searchMembers">查找会员</BaseButton>
        </div>
        <div v-if="memberCandidates.length" class="candidate-list">
          <div v-for="member in memberCandidates" :key="member.id" class="member-line">
            <span>{{ member.phone }}<small>{{ member.name || member.uid || '未填写姓名' }}</small></span>
            <BaseButton
              v-permission="'marketing:sms:edit'"
              variant="link"
              size="sm"
              :disabled="taggedMemberIds.has(member.id)"
              @click="addMemberToTag(member)"
            >{{ taggedMemberIds.has(member.id) ? '已在标签中' : '加入' }}</BaseButton>
          </div>
        </div>
        <div class="assigned-head">已分配会员（{{ tagMembers.length }}）</div>
        <BaseTable :columns="memberColumns" :data="(tagMembers as unknown) as Record<string, unknown>[]" :loading="tagMembersLoading" min-width="620px">
          <template #cell-actions="{ row }">
            <BaseButton v-permission="'marketing:sms:edit'" variant="link" size="sm" @click="removeMemberFromTag(row as MemberRow)">移除</BaseButton>
          </template>
        </BaseTable>
      </div>
      <template #footer><BaseButton variant="ghost" @click="membersDlg = false">关闭</BaseButton></template>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseCard,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BaseSelect,
  BaseTable,
  BaseTableRowActions,
  BaseTag,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn, TableRowAction } from '@/components/base/types'
import {
  bindMemberTag,
  cancelSmsCampaign,
  createMemberTag,
  createSmsCampaign,
  deleteMemberTag,
  deleteSmsCampaign,
  getSmsCampaign,
  getSmsServiceConfig,
  listMemberTagMembers,
  listMemberTags,
  listSmsCampaignRecords,
  listSmsCampaigns,
  sendSmsCampaign,
  searchSmsMembers,
  unbindMemberTag,
  updateMemberTag,
  updateSmsCampaign,
} from '@/api/smsPromotion'
import { listAllStores } from '@/api/store'
import type {
  MemberRow,
  MemberTag,
  SmsCampaign,
  SmsCampaignPayload,
  SmsCampaignSegment,
  SmsSendRecord,
  Store,
} from '@/api/types'
import { useUserStore } from '@/store/user'
import { usePermission } from '@/hooks/usePermission'
import { toast } from '@/feedback/toast'
import { confirmDialog } from '@/feedback/confirm'

interface SegmentForm {
  key: number
  id?: number
  tag_ids: number[]
  sign: string
  template: string
  params: string
  personalize: boolean
  default: boolean
}

const qc = useQueryClient()
const userStore = useUserStore()
const { hasPerm } = usePermission()
const currentStoreId = computed(() => Number(userStore.currentStoreId || 0))
const canChooseOwner = computed(() => currentStoreId.value === 0)
const canSendNow = computed(() => hasPerm('marketing:sms:send'))

const { data: config } = useQuery({ queryKey: ['sms-config'], queryFn: getSmsServiceConfig })
const { data: campaignData, isLoading: campaignLoading } = useQuery({
  queryKey: computed(() => ['sms-campaigns', currentStoreId.value] as const),
  queryFn: () => listSmsCampaigns(currentStoreId.value > 0 ? { store_id: currentStoreId.value } : undefined),
})
const { data: tagData, isLoading: tagLoading } = useQuery({
  queryKey: computed(() => ['member-tags', currentStoreId.value] as const),
  queryFn: () => listMemberTags(currentStoreId.value > 0 ? { store_id: currentStoreId.value } : undefined),
})
const { data: storeData } = useQuery({
  queryKey: ['stores', 'sms-promotion'],
  queryFn: async () => {
    try {
      return await listAllStores()
    } catch {
      return [] as Store[]
    }
  },
})

const campaigns = computed(() => campaignData.value ?? [])
const allTags = computed(() => tagData.value ?? [])
const stores = computed(() => {
  const rows = [...(storeData.value ?? [])]
  const own = userStore.userInfo?.store
  if (currentStoreId.value > 0 && !rows.some((store) => store.id === currentStoreId.value)) {
    rows.push({ id: currentStoreId.value, name: own?.name || `当前门店 #${currentStoreId.value}` })
  }
  return rows
})
const storeOnlyOptions = computed(() => stores.value.map((s) => ({ label: s.name, value: s.id })))
const ownerStoreOptions = computed(() => [
  { label: '总部全局（跨门店）', value: 0 },
  ...storeOnlyOptions.value,
])

const stats = computed(() => [
  { label: '活动总数', value: campaigns.value.length, tone: 'text-slate-900' },
  { label: '待发送', value: campaigns.value.filter((c) => ['draft', 'scheduled'].includes(c.status)).length, tone: 'text-amber-600' },
  { label: '发送成功', value: campaigns.value.reduce((n, c) => n + Number(c.success_count || 0), 0), tone: 'text-emerald-600' },
  { label: '发送失败', value: campaigns.value.reduce((n, c) => n + Number(c.fail_count || 0), 0), tone: 'text-red-600' },
])

const campaignColumns: BaseTableColumn[] = [
  { key: 'name', label: '活动', minWidth: '180px' },
  { key: 'owner_store_id', label: '归属', width: '120px' },
  { key: 'target_type', label: '受众', minWidth: '150px' },
  { key: 'segments', label: '模板分组', width: '90px' },
  { key: 'status', label: '状态', width: '90px' },
  { key: 'counts', label: '成功/失败/总数', width: '140px' },
  { key: 'scheduled_at', label: '计划/发送时间', width: '180px' },
  { key: 'actions', label: '操作', width: '260px', align: 'right' },
]
const recordColumns: BaseTableColumn[] = [
  { key: 'phone', label: '手机号', prop: 'phone', width: '140px' },
  { key: 'segment_id', label: '分组', minWidth: '130px' },
  { key: 'status', label: '状态', width: '90px' },
  { key: 'biz_id', label: '阿里云 BizID', prop: 'biz_id', minWidth: '150px', ellipsis: true },
  { key: 'error_message', label: '错误信息', prop: 'error_message', minWidth: '200px', ellipsis: true },
  { key: 'created_at', label: '发送时间', width: '180px' },
]
const tagColumns: BaseTableColumn[] = [
  { key: 'name', label: '标签', minWidth: '150px' },
  { key: 'store_id', label: '门店', width: '120px' },
  { key: 'member_count', label: '会员数', prop: 'member_count', width: '90px' },
  { key: 'description', label: '说明', prop: 'description', minWidth: '160px', ellipsis: true },
  { key: 'actions', label: '操作', width: '210px', align: 'right' },
]
const memberColumns: BaseTableColumn[] = [
  { key: 'phone', label: '手机号', prop: 'phone', width: '150px' },
  { key: 'name', label: '姓名', prop: 'name', minWidth: '120px' },
  { key: 'uid', label: 'UID', prop: 'uid', minWidth: '130px' },
  { key: 'actions', label: '操作', width: '90px', align: 'right' },
]

const campaignTypeOptions = [
  { label: '活动推广', value: 'activity' },
  { label: '节日关怀', value: 'holiday' },
]
const targetOptions = [
  { label: '全部会员', value: 'all_members' },
  { label: '指定门店会员', value: 'stores' },
  { label: '自定义手机号', value: 'custom' },
]

function storeName(id?: number): string {
  if (!id) return '总部全局'
  return stores.value.find((s) => s.id === id)?.name || `门店 #${id}`
}
function campaignTypeText(type: string): string {
  return type === 'holiday' ? '节日关怀' : '活动推广'
}
function targetText(row: SmsCampaign): string {
  if (row.target_type === 'custom') return `自定义手机号（${row.custom_phones?.length || 0}）`
  if (row.target_type === 'stores') return `指定门店（${row.store_ids?.length || (row.owner_store_id ? 1 : 0)}）`
  return '全部会员'
}
function statusMeta(status: string): { label: string; variant: 'success' | 'warning' | 'info' | 'danger' | 'neutral' } {
  const map: Record<string, { label: string; variant: 'success' | 'warning' | 'info' | 'danger' | 'neutral' }> = {
    draft: { label: '草稿', variant: 'neutral' },
    scheduled: { label: '已排期', variant: 'info' },
    sending: { label: '发送中', variant: 'warning' },
    sent: { label: '已发送', variant: 'success' },
    failed: { label: '失败', variant: 'danger' },
    cancelled: { label: '已取消', variant: 'neutral' },
  }
  return map[status] || { label: status || '未知', variant: 'neutral' }
}
function formatChinaTime(value?: string | null): string {
  if (!value) return '-'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('zh-CN', {
    timeZone: 'Asia/Shanghai',
    year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false,
  }).format(date).replaceAll('/', '-')
}

let segmentKey = 0
function emptySegment(isDefault: boolean): SegmentForm {
  return { key: ++segmentKey, tag_ids: [], sign: '', template: '', params: '{}', personalize: false, default: isDefault }
}

const campaignDlg = ref(false)
const savingCampaign = ref(false)
const editingId = ref(0)
const form = reactive({
  owner_store_id: 0,
  name: '',
  campaign_type: 'activity' as 'activity' | 'holiday',
  target_type: 'all_members' as 'all_members' | 'stores' | 'custom',
  store_ids: [] as number[],
  custom_phones: '',
  send_mode: 'now' as 'draft' | 'now' | 'scheduled',
  scheduled_at: '',
  segments: [] as SegmentForm[],
})

const campaignTags = computed(() => {
  if (!form.owner_store_id) return allTags.value
  return allTags.value.filter((t) => t.store_id === form.owner_store_id)
})
const targetStoreChoices = computed(() => {
  if (form.owner_store_id > 0) return stores.value.filter((s) => s.id === form.owner_store_id)
  return stores.value
})

watch(() => form.owner_store_id, (id) => {
  if (id > 0) form.store_ids = [id]
  else form.store_ids = form.store_ids.filter((storeId) => stores.value.some((s) => s.id === storeId))
  const valid = new Set(campaignTags.value.map((t) => t.id))
  form.segments.forEach((segment) => { segment.tag_ids = segment.tag_ids.filter((tagId) => valid.has(tagId)) })
})

function resetCampaignForm(): void {
  editingId.value = 0
  form.owner_store_id = currentStoreId.value
  form.name = ''
  form.campaign_type = 'activity'
  form.target_type = 'all_members'
  form.store_ids = currentStoreId.value ? [currentStoreId.value] : []
  form.custom_phones = ''
  form.send_mode = canSendNow.value ? 'now' : 'draft'
  form.scheduled_at = defaultScheduleTime()
  form.segments = [emptySegment(true)]
}
function defaultScheduleTime(): string {
  const now = new Date(Date.now() + 60 * 60 * 1000)
  const parts = new Intl.DateTimeFormat('en-CA', {
    timeZone: 'Asia/Shanghai', year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hourCycle: 'h23',
  }).formatToParts(now)
  const get = (type: string) => parts.find((p) => p.type === type)?.value || ''
  return `${get('year')}-${get('month')}-${get('day')} ${get('hour')}:${get('minute')}`
}
function chinaPickerValue(value?: string | null): string {
  if (!value) return defaultScheduleTime()
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return defaultScheduleTime()
  const formatted = new Intl.DateTimeFormat('sv-SE', {
    timeZone: 'Asia/Shanghai', year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hourCycle: 'h23',
  }).format(date)
  return formatted.slice(0, 16)
}
function openCreate(): void {
  resetCampaignForm()
  campaignDlg.value = true
}
async function openEdit(row: SmsCampaign): Promise<void> {
  try {
    const full = await getSmsCampaign(row.id, currentStoreId.value)
    editingId.value = full.id
    form.owner_store_id = Number(full.owner_store_id || currentStoreId.value || 0)
    form.name = full.name
    form.campaign_type = full.campaign_type === 'holiday' ? 'holiday' : 'activity'
    form.target_type = ['all_members', 'stores', 'custom'].includes(full.target_type)
      ? full.target_type as typeof form.target_type
      : 'all_members'
    form.store_ids = [...(full.store_ids || [])]
    if (form.owner_store_id && !form.store_ids.length) form.store_ids = [form.owner_store_id]
    form.custom_phones = (full.custom_phones || []).join('\n')
    form.send_mode = full.scheduled_at && full.status === 'scheduled' ? 'scheduled' : 'draft'
    form.scheduled_at = chinaPickerValue(full.scheduled_at)
    const sourceSegments = [...(full.segments || [])].sort((a, b) => Number(a.position || 0) - Number(b.position || 0))
    form.segments = sourceSegments.length
      ? sourceSegments.map(toSegmentForm)
      : [{
          ...emptySegment(true),
          sign: full.sign_name || '',
          template: full.template_code || '',
          params: full.template_param || '{}',
          personalize: Boolean(full.personalize_name),
        }]
    campaignDlg.value = true
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '加载活动失败')
  }
}
function toSegmentForm(segment: SmsCampaignSegment): SegmentForm {
  return {
    key: ++segmentKey,
    id: segment.id,
    tag_ids: [...(segment.tag_ids || [])],
    sign: segment.sign || '',
    template: segment.template || '',
    params: segment.params || '{}',
    personalize: Boolean(segment.personalize),
    default: Boolean(segment.default),
  }
}
function addSegment(isDefault: boolean): void {
  if (isDefault && form.segments.some((s) => s.default)) return
  form.segments.push(emptySegment(isDefault))
}
function removeSegment(index: number): void {
  form.segments.splice(index, 1)
  if (!form.segments.some((s) => s.default) && form.segments.length) form.segments[form.segments.length - 1].default = true
}
function moveSegment(index: number, direction: -1 | 1): void {
  const next = index + direction
  if (next < 0 || next >= form.segments.length) return
  const [item] = form.segments.splice(index, 1)
  form.segments.splice(next, 0, item)
}
function setDefault(index: number, enabled: boolean): void {
  if (!enabled) {
    toast.warning('活动必须保留一个默认分组')
    return
  }
  form.segments.forEach((segment, i) => {
    segment.default = i === index
    if (i === index) segment.tag_ids = []
  })
}
function parsePhones(raw: string): string[] {
  return [...new Set(raw.split(/[\s,;，；]+/).map((v) => v.trim()).filter(Boolean))]
}
function scheduledISO(raw: string): string | undefined {
  const value = raw.trim().replace(' ', 'T')
  if (!/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/.test(value)) return undefined
  return `${value}:00+08:00`
}
function validateSegments(): string | null {
  if (!form.segments.length) return '请至少配置一个短信分组'
  if (form.segments.filter((s) => s.default).length !== 1) return '请且只能配置一个默认分组'
  for (let i = 0; i < form.segments.length; i += 1) {
    const segment = form.segments[i]
    if (!segment.default && !segment.tag_ids.length) return `第 ${i + 1} 个分组请选择至少一个标签`
    if (!segment.template.trim()) return `第 ${i + 1} 个分组请填写模板 CODE`
    if (segment.params.trim()) {
      try {
        const parsed: unknown = JSON.parse(segment.params)
        if (!parsed || Array.isArray(parsed) || typeof parsed !== 'object') return `第 ${i + 1} 个分组的模板变量必须是 JSON 对象`
      } catch {
        return `第 ${i + 1} 个分组的模板变量不是合法 JSON`
      }
    }
  }
  return null
}
async function submitCampaign(): Promise<void> {
  if (!form.name.trim()) return toast.warning('请填写活动名称')
  if (canChooseOwner.value && form.owner_store_id < 0) return toast.warning('请选择归属门店')
  if (form.target_type === 'stores' && !form.store_ids.length) return toast.warning('请选择至少一个门店')
  const phones = parsePhones(form.custom_phones)
  if (form.target_type === 'custom' && !phones.length) return toast.warning('请填写至少一个手机号')
  const segmentError = validateSegments()
  if (segmentError) return toast.warning(segmentError)
  const schedule = form.send_mode === 'scheduled' ? scheduledISO(form.scheduled_at) : undefined
  if (form.send_mode === 'scheduled' && !schedule) return toast.warning('请选择计划发送时间')
  if (form.send_mode === 'now' && !canSendNow.value) return toast.warning('当前账号没有立即发送权限')
  if (form.send_mode === 'now') {
    const ok = await confirmDialog({ message: '保存后将立即调用阿里云短信发送，确认继续？' })
    if (!ok) return
  }

  const defaultSegment = form.segments.find((s) => s.default) || form.segments[0]
  const payload: SmsCampaignPayload = {
    owner_store_id: form.owner_store_id,
    name: form.name.trim(),
    campaign_type: form.campaign_type,
    sign_name: defaultSegment.sign.trim(),
    template_code: defaultSegment.template.trim(),
    template_param: defaultSegment.params.trim() || '{}',
    personalize_name: defaultSegment.personalize,
    target_type: form.target_type,
    store_ids: form.target_type === 'stores' ? [...form.store_ids] : [],
    custom_phones: form.target_type === 'custom' ? phones : [],
    scheduled_at: schedule,
    clear_scheduled_at: form.send_mode !== 'scheduled',
    segments: form.segments.map((segment) => ({
      tag_ids: segment.default ? [] : [...segment.tag_ids],
      sign: segment.sign.trim(),
      template: segment.template.trim(),
      params: segment.params.trim() || '{}',
      personalize: segment.personalize,
      default: segment.default,
    })),
  }

  savingCampaign.value = true
  try {
    let id = editingId.value
    if (id) await updateSmsCampaign(id, payload, currentStoreId.value)
    else {
      id = (await createSmsCampaign(payload)).id
      // If the subsequent immediate send fails, retries must update this draft instead of creating duplicates.
      editingId.value = id
    }
    if (form.send_mode === 'now') {
      await sendSmsCampaign(id, currentStoreId.value)
      toast.success('活动已发送')
    } else {
      toast.success(form.send_mode === 'scheduled' ? '活动已排期' : '草稿已保存')
    }
    campaignDlg.value = false
    await qc.invalidateQueries({ queryKey: ['sms-campaigns'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '保存活动失败')
    await qc.invalidateQueries({ queryKey: ['sms-campaigns'] })
  } finally {
    savingCampaign.value = false
  }
}

async function sendNow(row: SmsCampaign): Promise<void> {
  const ok = await confirmDialog({ message: `立即发送活动「${row.name}」？此操作会产生短信费用。` })
  if (!ok) return
  try {
    await sendSmsCampaign(row.id, currentStoreId.value)
    toast.success('发送完成')
    await qc.invalidateQueries({ queryKey: ['sms-campaigns'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '发送失败')
    await qc.invalidateQueries({ queryKey: ['sms-campaigns'] })
  }
}
async function cancelCampaign(row: SmsCampaign): Promise<void> {
  const ok = await confirmDialog({ message: `取消活动「${row.name}」？` })
  if (!ok) return
  try {
    await cancelSmsCampaign(row.id, currentStoreId.value)
    toast.success('已取消')
    await qc.invalidateQueries({ queryKey: ['sms-campaigns'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '取消失败')
  }
}
async function removeCampaign(row: SmsCampaign): Promise<void> {
  const ok = await confirmDialog({ message: `删除活动「${row.name}」及其发送记录？` })
  if (!ok) return
  try {
    await deleteSmsCampaign(row.id, currentStoreId.value)
    toast.success('已删除')
    await qc.invalidateQueries({ queryKey: ['sms-campaigns'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '删除失败')
  }
}
function campaignActions(row: SmsCampaign): TableRowAction[] {
  const editable = ['draft', 'scheduled'].includes(row.status)
  const sendable = ['draft', 'scheduled'].includes(row.status)
  const actions: TableRowAction[] = [
    { label: '记录', permission: 'marketing:sms:list', onClick: () => void openRecords(row), place: 'inline' },
  ]
  if (editable) actions.push({ label: '编辑', permission: 'marketing:sms:edit', onClick: () => void openEdit(row), place: 'inline' })
  if (sendable) actions.push({ label: '立即发送', permission: 'marketing:sms:send', onClick: () => void sendNow(row), place: 'inline' })
  if (['draft', 'scheduled'].includes(row.status)) actions.push({ label: '取消', permission: 'marketing:sms:edit', onClick: () => void cancelCampaign(row), place: 'more' })
  if (row.status !== 'sending') actions.push({ label: '删除', permission: 'marketing:sms:delete', danger: true, onClick: () => void removeCampaign(row), place: 'more' })
  return actions
}

const recordsDlg = ref(false)
const recordsLoading = ref(false)
const records = ref<SmsSendRecord[]>([])
const recordsCampaign = ref<SmsCampaign | null>(null)
async function openRecords(row: SmsCampaign): Promise<void> {
  recordsCampaign.value = row
  records.value = []
  recordsDlg.value = true
  recordsLoading.value = true
  try {
    records.value = await listSmsCampaignRecords(row.id, currentStoreId.value)
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '加载发送记录失败')
  } finally {
    recordsLoading.value = false
  }
}
function recordSegmentName(record: SmsSendRecord): string {
  if (!record.segment_id) return '活动默认模板'
  const segment = recordsCampaign.value?.segments?.find((item) => item.id === record.segment_id)
  if (!segment) return `分组 #${record.segment_id}`
  if (segment.default) return '默认分组'
  const names = (segment.tag_ids || []).map((id) => allTags.value.find((tag) => tag.id === id)?.name).filter(Boolean)
  return names.join(' / ') || `分组 #${record.segment_id}`
}

const tagDlg = ref(false)
const tagEditId = ref(0)
const savingTag = ref(false)
const tagForm = reactive({ store_id: 0, name: '', color: '#4f46e5', description: '' })
const visibleTags = computed(() => {
  if (!tagForm.store_id) return allTags.value
  return allTags.value.filter((tag) => tag.store_id === tagForm.store_id)
})
function resetTagForm(): void {
  tagEditId.value = 0
  tagForm.store_id = currentStoreId.value || stores.value[0]?.id || 0
  tagForm.name = ''
  tagForm.color = '#4f46e5'
  tagForm.description = ''
}
function openTagManager(): void {
  resetTagForm()
  tagDlg.value = true
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
    await qc.invalidateQueries({ queryKey: ['member-tags'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '删除标签失败')
  }
}
function tagActions(tag: MemberTag): TableRowAction[] {
  return [
    { label: '成员', permission: 'marketing:sms:list', onClick: () => void openTagMembers(tag), place: 'inline' },
    { label: '编辑', permission: 'marketing:sms:edit', onClick: () => editTag(tag), place: 'inline' },
    { label: '删除', permission: 'marketing:sms:delete', danger: true, onClick: () => void removeTag(tag), place: 'inline' },
  ]
}
function tagCheckStyle(tag: MemberTag): Record<string, string> {
  return { '--tag-color': tag.color || '#64748b' }
}

const membersDlg = ref(false)
const activeTag = ref<MemberTag | null>(null)
const tagMembers = ref<MemberRow[]>([])
const tagMembersLoading = ref(false)
const memberKeyword = ref('')
const memberCandidates = ref<MemberRow[]>([])
const memberSearchLoading = ref(false)
const taggedMemberIds = computed(() => new Set(tagMembers.value.map((member) => member.id)))
async function loadTagMembers(): Promise<void> {
  if (!activeTag.value) return
  tagMembersLoading.value = true
  try {
    tagMembers.value = await listMemberTagMembers(activeTag.value.id, activeTag.value.store_id)
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '加载标签成员失败')
  } finally {
    tagMembersLoading.value = false
  }
}
async function openTagMembers(tag: MemberTag): Promise<void> {
  activeTag.value = tag
  memberKeyword.value = ''
  memberCandidates.value = []
  tagMembers.value = []
  membersDlg.value = true
  await loadTagMembers()
}
async function searchMembers(): Promise<void> {
  if (!activeTag.value) return
  memberSearchLoading.value = true
  try {
    memberCandidates.value = await searchSmsMembers({
      store_id: activeTag.value.store_id,
      keyword: memberKeyword.value.trim() || undefined,
      limit: 20,
    })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '查找会员失败')
  } finally {
    memberSearchLoading.value = false
  }
}
async function addMemberToTag(member: MemberRow): Promise<void> {
  if (!activeTag.value) return
  try {
    await bindMemberTag(activeTag.value.id, member.id, activeTag.value.store_id)
    toast.success('已加入标签')
    await loadTagMembers()
    await qc.invalidateQueries({ queryKey: ['member-tags'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '分配标签失败')
  }
}
async function removeMemberFromTag(member: MemberRow): Promise<void> {
  if (!activeTag.value) return
  try {
    await unbindMemberTag(activeTag.value.id, member.id, activeTag.value.store_id)
    toast.success('已移除')
    await loadTagMembers()
    await qc.invalidateQueries({ queryKey: ['member-tags'] })
  } catch (error: unknown) {
    toast.error(error instanceof Error ? error.message : '移除失败')
  }
}
</script>

<style scoped>
.sms-page { display: flex; flex-direction: column; gap: 16px; min-width: 0; }
.page-head { display: flex; align-items: flex-end; justify-content: space-between; gap: 16px; }
.page-title { margin: 0; }
.page-subtitle { margin: 6px 0 0; color: var(--color-text-3); font-size: 14px; }
.head-actions { display: flex; flex-wrap: wrap; gap: 8px; flex-shrink: 0; }
.stats-grid { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); gap: 12px; }
.stat-label { color: var(--color-text-3); font-size: 13px; }
.stat-value { margin-top: 4px; font-size: 27px; font-weight: 700; line-height: 1.2; }
.cell-secondary { color: var(--color-text-3); font-size: 12px; margin-top: 2px; }
.campaign-form { display: flex; flex-direction: column; gap: 16px; max-height: 72vh; overflow-y: auto; padding-right: 4px; }
.form-section { padding: 16px; border: 1px solid var(--color-border-2); border-radius: 10px; background: var(--color-fill-1); }
.form-section h3, .tag-editor h3 { margin: 0 0 14px; font-size: 16px; }
.form-grid, .segment-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 0 16px; }
.section-head { display: flex; justify-content: space-between; align-items: flex-start; gap: 12px; }
.section-head h3 { margin-bottom: 5px; }
.section-head p { margin: 0 0 14px; color: var(--color-text-3); font-size: 13px; }
.choice-box, .tag-choice-box { display: flex; flex-wrap: wrap; gap: 8px 16px; padding: 10px; border: 1px solid var(--color-border-2); border-radius: 6px; background: var(--color-bg-2); }
.check-item, .tag-check { display: inline-flex; align-items: center; gap: 6px; cursor: pointer; }
.tag-check { padding: 5px 9px; border: 1px solid color-mix(in srgb, var(--tag-color) 38%, var(--color-border-2)); border-radius: 999px; }
.tag-check input { accent-color: var(--tag-color); }
.segment-list { display: flex; flex-direction: column; gap: 12px; }
.segment-card { padding: 14px; border: 1px solid var(--color-border-2); border-radius: 9px; background: var(--color-bg-2); }
.segment-head, .segment-title, .segment-actions, .default-switch { display: flex; align-items: center; }
.segment-head { justify-content: space-between; gap: 12px; margin-bottom: 10px; }
.segment-title, .segment-actions, .default-switch { gap: 8px; }
.segment-index { display: inline-grid; place-items: center; width: 24px; height: 24px; border-radius: 50%; color: white; background: rgb(var(--primary-6)); font-size: 12px; font-weight: 700; }
.default-switch { margin-bottom: 12px; font-size: 13px; color: var(--color-text-2); }
.schedule-field { max-width: 480px; margin-top: 14px; }
.permission-hint { margin: 9px 0 0; color: var(--color-text-3); font-size: 12px; }
.records-summary { display: flex; flex-wrap: wrap; gap: 18px; margin-bottom: 12px; font-weight: 600; }
.tag-manager { display: grid; grid-template-columns: 280px minmax(0, 1fr); gap: 18px; max-height: 68vh; overflow-y: auto; }
.tag-editor { padding: 14px; border: 1px solid var(--color-border-2); border-radius: 9px; align-self: start; }
.tag-list-wrap { min-width: 0; }
.color-field { display: grid; grid-template-columns: 42px 1fr; gap: 8px; align-items: center; }
.color-field input[type='color'] { width: 42px; height: 32px; padding: 2px; border: 1px solid var(--color-border-2); border-radius: 5px; background: transparent; }
.editor-actions { display: flex; justify-content: flex-end; gap: 8px; }
.color-dot { display: inline-block; width: 10px; height: 10px; margin-right: 7px; border-radius: 50%; }
.member-assignment { display: flex; flex-direction: column; gap: 12px; max-height: 68vh; overflow-y: auto; }
.member-search { display: grid; grid-template-columns: minmax(0, 1fr) auto; gap: 8px; }
.candidate-list { max-height: 220px; overflow-y: auto; border: 1px solid var(--color-border-2); border-radius: 7px; }
.member-line { display: flex; align-items: center; justify-content: space-between; gap: 12px; padding: 8px 12px; border-bottom: 1px solid var(--color-border-2); }
.member-line:last-child { border-bottom: 0; }
.member-line small { margin-left: 10px; color: var(--color-text-3); }
.assigned-head { margin-top: 4px; font-weight: 600; }
@media (max-width: 760px) {
  .page-head { align-items: stretch; flex-direction: column; }
  .head-actions > * { flex: 1; }
  .stats-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .form-grid, .segment-grid, .tag-manager { grid-template-columns: 1fr; }
  .section-head, .segment-head { flex-direction: column; align-items: stretch; }
  .segment-actions { flex-wrap: wrap; }
  .member-search { grid-template-columns: 1fr; }
}
</style>
