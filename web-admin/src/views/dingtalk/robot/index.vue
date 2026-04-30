<template>
  <div class="flex flex-col gap-4">
    <div class="flex flex-col md:flex-row md:items-end gap-3 justify-between">
      <h2 class="page-title">钉钉机器人配置</h2>
      <div class="flex flex-wrap gap-2">
        <BaseButton v-permission="'dingtalk:robot:add'" variant="primary" @click="openCreate">新增机器人</BaseButton>
      </div>
    </div>

    <BaseTable
      :columns="columns"
      :data="(list as unknown) as Record<string, unknown>[]"
      :loading="loading"
      min-width="1100px"
    >
      <template #cell-bot_type="{ row }">
        {{ (row as DingTalkBot).bot_type === 'stream' ? 'Stream' : 'Webhook' }}
      </template>
      <template #cell-webhook="{ row }">
        <span class="font-mono text-xs" :title="(row as DingTalkBot).webhook || ''">{{ previewUrl((row as DingTalkBot).webhook) }}</span>
      </template>
      <template #cell-is_enabled="{ row }">
        <BaseSwitch
          v-permission="'dingtalk:robot:edit'"
          :model-value="(row as DingTalkBot).is_enabled"
          @update:model-value="toggleEnabled(row as DingTalkBot, coerceBool($event))"
        />
      </template>
      <template #cell-actions="{ row }">
        <div class="flex flex-nowrap items-center justify-end gap-2 whitespace-nowrap shrink-0" @click.stop>
          <BaseButton v-permission="'dingtalk:robot:test'" variant="link" size="sm" @click="onTest(row as DingTalkBot)">测试推送</BaseButton>
          <BaseButton
            v-if="(row as DingTalkBot).bot_type === 'stream'"
            v-permission="'dingtalk:robot:test'"
            variant="link"
            size="sm"
            @click="onTestCallback(row as DingTalkBot)"
          >
            回调检测
          </BaseButton>
          <BaseButton v-permission="'dingtalk:robot:edit'" variant="link" size="sm" @click="openEdit(row as DingTalkBot)">编辑</BaseButton>
          <BaseButton v-permission="'dingtalk:robot:delete'" variant="link" size="sm" @click="onDelete(row as DingTalkBot)">删除</BaseButton>
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

    <BaseDialog v-model="dlg" :title="isEdit ? '编辑机器人' : '新增机器人'" max-width="min(640px, 96vw)">
      <div class="space-y-4 max-h-[72vh] overflow-y-auto pr-1">
        <BaseFormItem label="名称">
          <BaseInput v-model="form.name" placeholder="留空则按门店自动生成" />
        </BaseFormItem>
        <BaseFormItem label="类型" required>
          <BaseSelect
            v-model="form.bot_type"
            :options="[
              { label: 'Webhook', value: 'webhook' },
              { label: 'Stream', value: 'stream' },
            ]"
            :disabled="isEdit"
          />
        </BaseFormItem>
        <BaseFormItem label="所属门店">
          <BaseSelect
            v-model="formStoreId"
            :options="storeOptions"
            placeholder="全局（全部门店可用）"
            allow-clear
          />
        </BaseFormItem>
        <BaseFormItem v-if="form.bot_type === 'webhook'" label="Webhook URL" required>
          <BaseInput v-model="form.webhook" placeholder="https://oapi.dingtalk.com/robot/send?access_token=..." />
        </BaseFormItem>
        <BaseFormItem v-if="form.bot_type === 'webhook'" label="加签 Secret">
          <BaseInput v-model="form.secret" placeholder="可选" />
        </BaseFormItem>
        <template v-if="form.bot_type === 'stream'">
          <BaseFormItem label="Client ID" required>
            <BaseInput v-model="form.client_id" />
          </BaseFormItem>
          <BaseFormItem label="Client Secret" required>
            <BaseInput v-model="form.client_secret" type="password" placeholder="编辑时留空表示不修改" />
          </BaseFormItem>
          <BaseFormItem label="Agent ID">
            <BaseInput v-model="form.agent_id" />
          </BaseFormItem>
          <BaseFormItem label="Robot Code">
            <BaseInput v-model="form.robot_code" placeholder="钉钉控制台机器人编码" />
          </BaseFormItem>
        </template>
        <BaseFormItem label="消息类型">
          <BaseSelect
            v-model="form.msg_type"
            :options="[
              { label: 'Markdown', value: 'markdown' },
              { label: 'Text', value: 'text' },
              { label: '卡片', value: 'card' },
            ]"
          />
        </BaseFormItem>
        <BaseFormItem v-if="form.msg_type === 'card'" label="卡片模板Key" required>
          <BaseInput v-model="form.card_msg_key" placeholder="如 sampleActionCard / your.card.template.key" />
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
import { useQuery, useQueryClient } from '@tanstack/vue-query'
import {
  BaseButton,
  BaseDialog,
  BaseFormItem,
  BaseInput,
  BasePagination,
  BaseSelect,
  BaseSwitch,
  BaseTable,
  BaseTextarea,
} from '@/components/base'
import type { BaseTableColumn } from '@/components/base/types'
import {
  createDingTalkBot,
  deleteDingTalkBot,
  getDingTalkBot,
  listDingTalkBots,
  testDingTalkBot,
  testDingTalkStreamCallback,
  updateDingTalkBot,
} from '@/api/dingtalkBot'
import { listAllStores } from '@/api/store'
import type { DingTalkBot, Store, UpdateDingTalkBotReq } from '@/api/types'
import { confirmDialog } from '@/feedback/confirm'
import { toast } from '@/feedback/toast'

const qc = useQueryClient()
const page = ref(1)
const pageSize = ref(10)

const queryKey = computed(() => ['dingtalk-bots', page.value, pageSize.value] as const)

const { data: pageResult, isLoading: loading } = useQuery({
  queryKey: queryKey,
  queryFn: () => listDingTalkBots({ page: page.value, page_size: pageSize.value }),
})

const list = computed(() => pageResult.value?.list ?? [])
const total = computed(() => pageResult.value?.total ?? 0)

const columns: BaseTableColumn[] = [
  { key: 'name', label: '名称', prop: 'name', minWidth: '160px', ellipsis: true },
  { key: 'bot_type', label: '类型', prop: 'bot_type', width: '100px' },
  { key: 'store_name', label: '门店', prop: 'store_name', width: '120px', ellipsis: true },
  { key: 'webhook', label: 'Webhook', prop: 'webhook', minWidth: '200px' },
  { key: 'robot_code', label: 'RobotCode', prop: 'robot_code', width: '120px', ellipsis: true },
  { key: 'is_enabled', label: '启用', prop: 'is_enabled', width: '88px' },
  { key: 'msg_type', label: '消息', prop: 'msg_type', width: '88px' },
  { key: 'actions', label: '操作', width: '280px', align: 'right' },
]

function previewUrl(u?: string): string {
  if (!u) return '-'
  return u.length > 48 ? `${u.slice(0, 48)}…` : u
}

function coerceBool(v: unknown): boolean {
  if (typeof v === 'boolean') return v
  if (typeof v === 'number') return v !== 0
  return Boolean(v)
}

const stores = ref<Store[]>([])
void listAllStores().then((r) => {
  stores.value = r
})

const storeOptions = computed(() => [
  { label: '全局', value: 0 },
  ...stores.value.map((s) => ({ label: s.name, value: s.id })),
])

const dlg = ref(false)
const saving = ref(false)
const isEdit = ref(false)
const editId = ref(0)

const form = reactive({
  name: '',
  bot_type: 'webhook',
  webhook: '',
  secret: '',
  client_id: '',
  client_secret: '',
  agent_id: '',
  robot_code: '',
  msg_type: 'markdown',
  card_msg_key: '',
  remark: '',
})

const formStoreId = ref<number | undefined>(undefined)

function resetForm(): void {
  form.name = ''
  form.bot_type = 'webhook'
  form.webhook = ''
  form.secret = ''
  form.client_id = ''
  form.client_secret = ''
  form.agent_id = ''
  form.robot_code = ''
  form.msg_type = 'markdown'
  form.card_msg_key = ''
  form.remark = ''
  formStoreId.value = undefined
}

function openCreate(): void {
  isEdit.value = false
  editId.value = 0
  resetForm()
  dlg.value = true
}

async function openEdit(row: DingTalkBot): Promise<void> {
  isEdit.value = true
  editId.value = row.id
  resetForm()
  try {
    const full = await getDingTalkBot(row.id)
    form.name = full.name ?? ''
    form.bot_type = full.bot_type || 'webhook'
    form.webhook = full.webhook ?? ''
    form.secret = full.secret ?? ''
    form.client_id = full.client_id ?? ''
    form.client_secret = ''
    form.agent_id = full.agent_id ?? ''
    form.robot_code = full.robot_code ?? ''
    form.msg_type = full.msg_type || 'markdown'
    form.card_msg_key = full.card_msg_key ?? ''
    form.remark = full.remark ?? ''
    const sid = full.store_id
    formStoreId.value = sid != null && sid > 0 ? Number(sid) : 0
    dlg.value = true
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '加载失败')
  }
}

async function submit(): Promise<void> {
  const storePayload =
    formStoreId.value === undefined || formStoreId.value === 0 ? null : Number(formStoreId.value)

  saving.value = true
  try {
    if (form.msg_type === 'card' && form.bot_type !== 'stream') {
      toast.warning('卡片消息仅支持 Stream 机器人')
      return
    }
    if (form.msg_type === 'card' && !form.card_msg_key.trim()) {
      toast.warning('请选择卡片类型并填写卡片模板Key')
      return
    }
    if (!isEdit.value) {
      if (form.bot_type === 'webhook' && !form.webhook.trim()) {
        toast.warning('请填写 Webhook')
        return
      }
      if (form.bot_type === 'stream' && (!form.client_id.trim() || !form.client_secret.trim())) {
        toast.warning('请填写 Client ID 与 Client Secret')
        return
      }
      if (form.bot_type === 'webhook') {
        await createDingTalkBot({
          name: form.name.trim() || undefined,
          bot_type: 'webhook',
          webhook: form.webhook.trim(),
          secret: form.secret.trim(),
          store_id: storePayload,
          msg_type: form.msg_type,
          card_msg_key: form.msg_type === 'card' ? form.card_msg_key.trim() : '',
          remark: form.remark.trim(),
        })
      } else {
        await createDingTalkBot({
          name: form.name.trim() || undefined,
          bot_type: 'stream',
          client_id: form.client_id.trim(),
          client_secret: form.client_secret.trim(),
          agent_id: form.agent_id.trim(),
          robot_code: form.robot_code.trim(),
          store_id: storePayload,
          msg_type: form.msg_type,
          card_msg_key: form.msg_type === 'card' ? form.card_msg_key.trim() : '',
          remark: form.remark.trim(),
        })
      }
      toast.success('已创建')
    } else {
      const body: UpdateDingTalkBotReq = {
        name: form.name.trim(),
        bot_type: form.bot_type,
        msg_type: form.msg_type,
        card_msg_key: form.msg_type === 'card' ? form.card_msg_key.trim() || null : null,
        remark: form.remark.trim() || null,
        store_id: storePayload,
      }
      if (form.bot_type === 'webhook') {
        body.webhook = form.webhook.trim()
        body.secret = form.secret.trim() || null
        body.client_id = null
        body.client_secret = null
        body.agent_id = null
        body.robot_code = null
      } else {
        body.webhook = null
        body.secret = null
        body.client_id = form.client_id.trim() || null
        if (form.client_secret.trim()) body.client_secret = form.client_secret.trim()
        body.agent_id = form.agent_id.trim() || null
        body.robot_code = form.robot_code.trim() || null
      }
      await updateDingTalkBot(editId.value, body)
      toast.success('已保存')
    }
    dlg.value = false
    void qc.invalidateQueries({ queryKey: ['dingtalk-bots'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '失败')
  } finally {
    saving.value = false
  }
}

async function toggleEnabled(row: DingTalkBot, v: boolean): Promise<void> {
  try {
    await updateDingTalkBot(row.id, { is_enabled: v })
    toast.success(v ? '已启用' : '已禁用')
    void qc.invalidateQueries({ queryKey: ['dingtalk-bots'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '更新失败')
    void qc.invalidateQueries({ queryKey: ['dingtalk-bots'] })
  }
}

async function onDelete(row: DingTalkBot): Promise<void> {
  const ok = await confirmDialog({ message: `删除机器人「${row.name}」？` })
  if (!ok) return
  try {
    await deleteDingTalkBot(row.id)
    toast.success('已删除')
    void qc.invalidateQueries({ queryKey: ['dingtalk-bots'] })
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '删除失败')
  }
}

async function onTest(row: DingTalkBot): Promise<void> {
  try {
    await testDingTalkBot(row.id)
    toast.success('测试消息已发送')
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '测试失败')
  }
}

async function onTestCallback(row: DingTalkBot): Promise<void> {
  try {
    const r = await testDingTalkStreamCallback(row.id)
    toast.success(String(r.message ?? '回调正常'))
  } catch (e: unknown) {
    toast.error(e instanceof Error ? e.message : '检测失败')
  }
}
</script>
