<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col gap-3">
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <h2 class="page-title">操作日志</h2>
        <div class="flex gap-2">
          <BaseButton variant="ghost" @click="resetFilters">重置</BaseButton>
          <BaseButton variant="primary" :loading="loading" @click="reload">查询</BaseButton>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-6 gap-2">
        <BaseInput v-model="keyword" placeholder="关键词 / 单号 / 路径" clearable @enter="reload" />
        <BaseSelect v-model="moduleFilter" :options="moduleOptions" placeholder="模块" />
        <BaseSelect v-model="actionFilter" :options="actionOptions" placeholder="操作类型" />
        <BaseSelect v-model="statusFilter" :options="statusOptions" placeholder="结果" />
        <BaseInput v-model="startTime" placeholder="开始时间 2026-06-15 00:00:00" clearable @enter="reload" />
        <BaseInput v-model="endTime" placeholder="结束时间 2026-06-15 23:59:59" clearable @enter="reload" />
      </div>
    </div>

    <BaseTable
      :columns="columns"
      :data="(list as unknown) as Record<string, unknown>[]"
      :loading="loading"
      min-width="1180px"
      row-clickable
      @row-dblclick="openDetailFromRow"
    >
      <template #cell-created_at="{ row }">
        {{ formatTime((row as AuditLog).created_at) }}
      </template>
      <template #cell-operator="{ row }">
        <div class="leading-5">
          <div class="font-medium text-slate-800">{{ operatorName(row as AuditLog) }}</div>
          <div class="text-xs text-slate-500">{{ (row as AuditLog).phone || (row as AuditLog).username || '-' }}</div>
        </div>
      </template>
      <template #cell-action="{ row }">
        <div class="flex items-center gap-2">
          <BaseTag variant="info">{{ (row as AuditLog).module_name || (row as AuditLog).module }}</BaseTag>
          <span>{{ (row as AuditLog).action_name || (row as AuditLog).action }}</span>
        </div>
      </template>
      <template #cell-resource="{ row }">
        {{ resourceText(row as AuditLog) }}
      </template>
      <template #cell-status="{ row }">
        <BaseTag :variant="(row as AuditLog).status === 'success' ? 'success' : 'danger'">
          {{ (row as AuditLog).status === 'success' ? '成功' : '失败' }}
        </BaseTag>
      </template>
      <template #cell-source="{ row }">
        <div class="leading-5">
          <div>{{ (row as AuditLog).client_source || '-' }}</div>
          <div class="text-xs text-slate-500">{{ (row as AuditLog).client_ip || '-' }}</div>
        </div>
      </template>
      <template #cell-device="{ row }">
        <div class="leading-5">
          <div>{{ deviceText(row as AuditLog) }}</div>
          <div class="text-xs text-slate-500">{{ (row as AuditLog).browser || '-' }}</div>
        </div>
      </template>
      <template #cell-actions="{ row }">
        <BaseTableRowActions :actions="[{ label: '详情', permission: 'system:audit-log:detail', onClick: () => openDetail(row as AuditLog) }]" />
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

    <BaseDialog v-model="detailVisible" title="操作日志详情" max-width="min(920px, 96vw)">
      <div v-if="detail" class="max-h-[72vh] overflow-y-auto space-y-4 pr-1">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
          <InfoItem label="操作时间" :value="formatTime(detail.created_at)" />
          <InfoItem label="操作结果" :value="detail.status === 'success' ? '成功' : '失败'" />
          <InfoItem label="操作人" :value="operatorName(detail)" />
          <InfoItem label="所属门店" :value="detail.store_name || (detail.store_id ? String(detail.store_id) : '总部/未知')" />
          <InfoItem label="模块" :value="detail.module_name || detail.module" />
          <InfoItem label="操作" :value="detail.action_name || detail.action" />
          <InfoItem label="操作对象" :value="resourceText(detail)" />
          <InfoItem label="来源" :value="`${detail.client_source || '-'} / ${detail.client_ip || '-'}`" />
          <InfoItem label="设备" :value="deviceDetailText(detail)" />
          <InfoItem label="接口" :value="`${detail.method || '-'} ${detail.path || '-'}`" />
          <InfoItem label="耗时" :value="`${detail.latency_ms ?? 0} ms`" />
        </div>

        <div v-if="detail.error_message" class="rounded border border-red-200 bg-red-50 p-3 text-sm text-red-700">
          {{ detail.error_message }}
        </div>

        <DetailBlock title="请求参数" :content="detail.request_body" />
        <DetailBlock title="修改前" :content="detail.before_data" />
        <DetailBlock title="修改后" :content="detail.after_data" />
        <DetailBlock title="字段差异" :content="detail.diff_data" />
        <DetailBlock title="User-Agent" :content="detail.user_agent" plain />
      </div>
    </BaseDialog>
  </div>
</template>

<script setup lang="ts">
import { defineComponent, h, ref, watch } from 'vue'
import {
  BaseButton,
  BaseDialog,
  BaseInput,
  BasePagination,
  BaseSelect,
  BaseTable,
  BaseTableRowActions,
  BaseTag,
} from '@/components/base'
import type { BaseSelectOption, BaseTableColumn } from '@/components/base/types'
import { getAuditLog, listAuditLogs } from '@/api/auditLog'
import type { AuditLog } from '@/api/types'
import { toast } from '@/feedback/toast'

const columns: BaseTableColumn[] = [
  { key: 'created_at', label: '操作时间', width: '170px' },
  { key: 'operator', label: '操作人', minWidth: '150px' },
  { key: 'store_name', label: '门店', prop: 'store_name', minWidth: '130px', ellipsis: true },
  { key: 'action', label: '模块 / 操作', minWidth: '190px' },
  { key: 'resource', label: '操作对象', minWidth: '180px', ellipsis: true },
  { key: 'status', label: '结果', width: '86px' },
  { key: 'source', label: '来源', minWidth: '150px' },
  { key: 'device', label: '设备', minWidth: '150px' },
  { key: 'path', label: '接口路径', prop: 'path', minWidth: '220px', ellipsis: true },
  { key: 'actions', label: '操作', width: '100px', align: 'right' },
]

const moduleOptions: BaseSelectOption[] = [
  { label: '全部模块', value: '' },
  { label: '认证', value: 'auth' },
  { label: '用户管理', value: 'user' },
  { label: '角色管理', value: 'role' },
  { label: '菜单管理', value: 'menu' },
  { label: '门店管理', value: 'store' },
  { label: '供应商管理', value: 'supplier' },
  { label: '采购管理', value: 'purchase_order' },
  { label: '库存管理', value: 'inventory' },
  { label: '门店记账', value: 'store_account' },
  { label: '门店退货', value: 'store_return' },
  { label: '会员管理', value: 'member' },
  { label: 'B2B', value: 'b2b' },
  { label: '第三方账号', value: 'third_party_account' },
  { label: '第三方路线', value: 'third_party_route' },
  { label: '打印机', value: 'printer' },
  { label: '图库管理', value: 'gallery' },
]

const actionOptions: BaseSelectOption[] = [
  { label: '全部操作', value: '' },
  { label: '登录', value: 'login' },
  { label: '新增', value: 'create' },
  { label: '修改', value: 'update' },
  { label: '删除', value: 'delete' },
  { label: '导入', value: 'import' },
  { label: '导出', value: 'export' },
  { label: '审核', value: 'approve' },
  { label: '打印', value: 'print' },
  { label: '同步', value: 'sync' },
]

const statusOptions: BaseSelectOption[] = [
  { label: '全部结果', value: '' },
  { label: '成功', value: 'success' },
  { label: '失败', value: 'fail' },
]

const keyword = ref('')
const moduleFilter = ref<string | number | undefined>('')
const actionFilter = ref<string | number | undefined>('')
const statusFilter = ref<string | number | undefined>('')
const startTime = ref('')
const endTime = ref('')
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const list = ref<AuditLog[]>([])
const loading = ref(false)
const detailVisible = ref(false)
const detail = ref<AuditLog | null>(null)

async function load(): Promise<void> {
  loading.value = true
  try {
    const res = await listAuditLogs({
      page: page.value,
      page_size: pageSize.value,
      keyword: keyword.value || undefined,
      module: String(moduleFilter.value || '') || undefined,
      action: String(actionFilter.value || '') || undefined,
      status: String(statusFilter.value || '') || undefined,
      start_time: startTime.value || undefined,
      end_time: endTime.value || undefined,
    })
    list.value = res.list ?? []
    total.value = res.total ?? 0
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载操作日志失败')
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

function resetFilters(): void {
  keyword.value = ''
  moduleFilter.value = ''
  actionFilter.value = ''
  statusFilter.value = ''
  startTime.value = ''
  endTime.value = ''
  reload()
}

watch([page, pageSize], () => {
  void load()
})

void load()

function openDetailFromRow(row: Record<string, unknown>): void {
  void openDetail(row as unknown as AuditLog)
}

async function openDetail(row: AuditLog): Promise<void> {
  try {
    detail.value = await getAuditLog(row.id)
    detailVisible.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载详情失败')
  }
}

function operatorName(row: AuditLog): string {
  return row.nickname || row.username || row.phone || (row.user_id ? `用户 #${row.user_id}` : '未知用户')
}

function resourceText(row: AuditLog): string {
  return row.resource_name || row.resource_no || row.resource_id || '-'
}

function deviceTypeLabel(value?: string): string {
  const m: Record<string, string> = {
    desktop: '桌面端',
    mobile: '移动端',
    tablet: '平板',
    bot: '机器人',
    unknown: '未知设备',
  }
  return m[value || ''] || value || '未知设备'
}

function deviceText(row: AuditLog): string {
  const parts = [deviceTypeLabel(row.device_type), row.os].filter(Boolean)
  return parts.join(' / ') || '-'
}

function deviceDetailText(row: AuditLog): string {
  return [deviceTypeLabel(row.device_type), row.os, row.browser].filter(Boolean).join(' / ') || '-'
}

function formatTime(value?: string): string {
  if (!value) return '-'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return value
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function formatJSON(content?: string): string {
  if (!content) return ''
  try {
    return JSON.stringify(JSON.parse(content), null, 2)
  } catch {
    return content
  }
}

const InfoItem = defineComponent({
  props: { label: { type: String, required: true }, value: { type: String, default: '-' } },
  setup(props) {
    return () =>
      h('div', { class: 'rounded border border-slate-200 bg-slate-50 px-3 py-2 min-w-0' }, [
        h('div', { class: 'text-xs text-slate-500 mb-1' }, props.label),
        h('div', { class: 'text-slate-800 break-words' }, props.value || '-'),
      ])
  },
})

const DetailBlock = defineComponent({
  props: {
    title: { type: String, required: true },
    content: { type: String, default: '' },
    plain: { type: Boolean, default: false },
  },
  setup(props) {
    return () => {
      if (!props.content) return null
      return h('div', { class: 'space-y-2' }, [
        h('div', { class: 'font-medium text-slate-700' }, props.title),
        h(
          'pre',
          {
            class:
              'max-h-64 overflow-auto rounded border border-slate-200 bg-slate-950 p-3 text-xs leading-5 text-slate-100 whitespace-pre-wrap break-words',
          },
          props.plain ? props.content : formatJSON(props.content),
        ),
      ])
    }
  },
})
</script>
