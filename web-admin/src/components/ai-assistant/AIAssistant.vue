<template>
  <div class="pointer-events-none fixed inset-0 z-[900]">
    <Transition name="assistant-launch">
      <a-button
        v-if="!open && canUseAssistant"
        type="primary"
        shape="circle"
        class="assistant-launch pointer-events-auto fixed bottom-6 right-6 max-sm:bottom-3 max-sm:right-3"
        title="打开经营分析助手"
        aria-label="打开经营分析助手"
        @click="openAssistant"
      ><IconRobot /></a-button>
    </Transition>

    <Transition name="assistant-panel">
      <section
        v-if="open"
        ref="panel"
        class="assistant-panel pointer-events-auto fixed flex flex-col overflow-hidden rounded-lg border border-slate-200 bg-white text-slate-800 shadow-2xl"
        :class="{ 'assistant-panel-fullscreen': fullscreen }"
        :style="panelStyle"
        aria-label="经营分析助手"
      >
        <header class="assistant-drag-handle flex h-16 shrink-0 touch-none items-center justify-between border-b border-slate-200 bg-white px-4 sm:px-5" @pointerdown="startDrag">
          <div class="flex min-w-0 items-center gap-3">
            <span class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-teal-50 text-lg text-teal-700"><IconRobot /></span>
            <div class="min-w-0"><strong class="block truncate text-sm font-semibold">经营分析助手</strong><span class="mt-0.5 block truncate text-xs text-slate-500">{{ storeLabel }} · {{ activeConversation?.title || '新对话' }}</span></div>
          </div>
          <div class="ml-3 flex shrink-0 items-center gap-1">
            <a-button v-if="isHQAdmin" type="text" class="assistant-icon-button" :class="{ 'text-teal-700 bg-teal-50': settings }" title="模型设置" aria-label="模型设置" @pointerdown.stop @click="toggleSettings"><IconSettings /></a-button>
            <a-button type="text" class="assistant-icon-button" :title="fullscreen ? '退出全屏' : '全屏显示'" :aria-label="fullscreen ? '退出全屏' : '全屏显示'" @pointerdown.stop @click="toggleFullscreen"><IconFullscreenExit v-if="fullscreen" /><IconFullscreen v-else /></a-button>
            <a-button type="text" class="assistant-icon-button" title="关闭" aria-label="关闭" @pointerdown.stop @click="closeAssistant"><IconClose /></a-button>
          </div>
        </header>

        <div class="flex min-h-0 flex-1">
          <aside class="assistant-sidebar flex w-60 shrink-0 flex-col border-r border-slate-200 bg-slate-50/80 p-3 sm:p-4 max-sm:w-40 max-sm:p-2.5">
            <div class="mb-3 flex items-center justify-between px-1"><span class="text-xs font-semibold tracking-wide text-slate-500">对话</span><span class="text-[11px] tabular-nums text-slate-400">{{ conversations.length }}</span></div>
            <a-button type="primary" class="assistant-new-chat mb-4 shrink-0" long @click="newConversation"><template #icon><IconPlus /></template>新对话</a-button>
            <div class="min-h-0 flex-1 space-y-1 overflow-y-auto">
              <div v-if="!conversations.length" class="px-2 py-4 text-xs leading-5 text-slate-400">还没有对话</div>
              <div v-for="item in conversations" :key="item.id" class="assistant-conversation-row flex items-center gap-1 rounded-md px-1" :class="item.id === activeId ? 'assistant-conversation-active' : ''">
                <button class="assistant-conversation-title min-w-0 flex-1 truncate px-2 py-2.5 text-left text-xs" :title="item.title || '新对话'" @click="selectConversation(item.id)">{{ item.title || '新对话' }}</button>
                <a-button type="text" size="mini" class="assistant-row-action" title="重命名" aria-label="重命名" @click="openRename(item)"><IconEdit /></a-button>
                <a-popconfirm type="warning" :content="`删除“${item.title}”？`" ok-text="删除" @ok="remove(item)"><a-button type="text" size="mini" class="assistant-row-action assistant-delete-action" title="删除" aria-label="删除"><IconDelete /></a-button></a-popconfirm>
              </div>
            </div>
          </aside>

          <div class="relative flex min-w-0 flex-1 flex-col bg-white">
            <Transition name="assistant-settings">
              <section v-if="settings && isHQAdmin" class="absolute inset-0 z-10 overflow-y-auto bg-white p-5 sm:p-7" aria-label="模型设置">
                <div class="mx-auto max-w-xl">
                  <div class="mb-6 flex items-center justify-between"><div><h2 class="m-0 text-base font-semibold">模型设置</h2><p class="mb-0 mt-1 text-xs text-slate-500">配置经营分析助手使用的模型服务</p></div><a-button type="text" class="assistant-icon-button" title="返回对话" aria-label="返回对话" @click="settings = false"><IconClose /></a-button></div>
                  <label class="assistant-field">服务商<a-select v-model="config.provider" @change="providerChanged"><a-option value="openai">OpenAI</a-option><a-option value="deepseek">DeepSeek</a-option><a-option value="qwen">通义千问</a-option><a-option value="custom">OpenAI 兼容服务</a-option></a-select></label>
                  <label class="assistant-field">模型<a-input v-model="config.model" placeholder="模型名称" allow-clear /></label>
                  <label class="assistant-field">API 地址<a-input v-model="config.base_url" placeholder="https://api.openai.com/v1" allow-clear /></label>
                  <label class="assistant-field">API Key<a-input-password v-model="config.api_key" autocomplete="new-password" :placeholder="config.api_key_configured ? config.api_key_masked + '（留空保持不变）' : '填写 API Key'" /></label>
                  <label class="assistant-field">系统指令<a-textarea v-model="config.system_prompt" :auto-size="{ minRows: 4, maxRows: 8 }" placeholder="可选" /></label>
                  <label class="mt-5 flex items-center gap-2 text-sm text-slate-700"><a-switch v-model="config.enabled" />启用助手</label>
                  <p v-if="configError" class="mt-3 text-sm text-rose-700">{{ configError }}</p>
                  <div class="mt-6 flex justify-end gap-2"><a-button :disabled="saving" @click="testConfig">测试连接</a-button><a-button type="primary" :loading="saving" @click="saveConfig">保存设置</a-button></div>
                </div>
              </section>
            </Transition>

            <div v-show="!settings" ref="scrollArea" class="min-h-0 flex-1 overflow-y-auto px-4 py-5 sm:px-7">
              <div v-if="!messages.length" class="flex min-h-full flex-col items-center justify-center text-center">
                <span class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-teal-50 text-2xl text-teal-700"><IconRobot /></span>
                <h2 class="m-0 text-lg font-semibold text-slate-800">你好，需要分析什么？</h2><p class="mt-2 max-w-sm text-sm leading-6 text-slate-500">可以询问经营表现、销售趋势、渠道结构与库存情况。</p>
                <div class="mt-4 flex flex-wrap justify-center gap-2"><a-button size="small" shape="round" @click="usePrompt('分析近期销售趋势和主要变化')">分析近期销售趋势</a-button><a-button size="small" shape="round" @click="usePrompt('总结当前经营情况，并给出优先行动建议')">总结经营情况</a-button></div>
              </div>
              <article v-for="message in messages" :key="message.id" class="assistant-message mb-7" :class="message.role === 'user' ? 'ml-auto max-w-[88%]' : 'max-w-full'">
                <div class="mb-2 flex items-center gap-2 text-xs text-slate-500" :class="message.role === 'user' ? 'justify-end' : ''"><span class="flex h-6 w-6 items-center justify-center rounded-full" :class="message.role === 'user' ? 'bg-slate-200 text-slate-700' : 'bg-teal-100 text-teal-800'"><IconUser v-if="message.role === 'user'" /><IconRobot v-else /></span><span>{{ message.role === 'user' ? '你' : '经营分析助手' }}</span><span v-if="message.model" class="text-slate-400">{{ message.model }}</span></div>
                <div v-if="message.role === 'user'" class="ml-auto w-fit max-w-full whitespace-pre-wrap break-words rounded-2xl rounded-tr-sm bg-slate-100 px-4 py-3 text-sm leading-6 text-slate-800">{{ message.content }}</div>
                <div v-else>
                  <div class="assistant-markdown text-sm leading-7 text-slate-800" v-html="renderMarkdown(message.content)" />
                  <details v-if="message.analysis_context" class="mt-3 rounded-lg border border-slate-200 bg-slate-50 text-xs text-slate-600"><summary class="cursor-pointer px-3 py-2 font-medium">本次分析使用的数据</summary><pre class="m-0 max-h-56 overflow-auto border-t border-slate-200 p-3 font-mono text-[11px] leading-5">{{ formatAnalysisContext(message.analysis_context) }}</pre></details>
                </div>
                <div class="assistant-message-tools" :class="message.role === 'user' ? 'justify-end' : ''"><a-button type="text" size="mini" class="assistant-copy-button" title="复制消息" aria-label="复制消息" @click="copyMessage(message.content)"><IconCopy />复制</a-button></div>
              </article>
              <div v-if="sending" class="mb-5 flex items-center gap-2 text-xs text-slate-500"><span class="flex h-6 w-6 items-center justify-center rounded-full bg-teal-100 text-teal-800"><IconRobot /></span><span>正在分析</span><span class="assistant-thinking-dots"><i></i><i></i><i></i></span></div>
            </div>

            <form v-show="!settings" class="shrink-0 border-t border-slate-200 bg-white px-3 pb-3 pt-2 sm:px-5" @submit.prevent="send">
              <div class="mb-2 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-slate-500"><span>分析区间</span><a-range-picker v-model="dateRange" size="mini" value-format="YYYY-MM-DD" :allow-clear="false" /><span>最多 366 天</span></div>
              <div class="rounded-xl border border-slate-300 bg-white p-2 transition focus-within:border-teal-600 focus-within:shadow-sm">
                <a-textarea v-model="draft" :auto-size="{ minRows: 2, maxRows: 6 }" class="assistant-composer" placeholder="询问经营数据…" :disabled="sending" @keydown.enter.exact.prevent="send" />
                <div class="flex items-center justify-between px-1 pt-1"><span class="text-[11px] text-slate-400">回答基于所选门店和日期范围的数据</span><a-button type="primary" class="assistant-composer-send" shape="circle" html-type="submit" :disabled="!draft.trim() || sending" title="发送" aria-label="发送"><IconSend /></a-button></div>
              </div>
            </form>
          </div>
        </div>
      </section>
    </Transition>

    <a-modal v-model:visible="renameVisible" title="重命名对话" :on-before-ok="renameConversation" unmount-on-close>
      <a-input v-model="renameTitle" placeholder="输入对话名称" :max-length="80" show-word-limit @press-enter="renameConversation" />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, reactive, ref, watch } from 'vue'
import { IconClose, IconCopy, IconDelete, IconEdit, IconFullscreen, IconFullscreenExit, IconPlus, IconRobot, IconSend, IconSettings, IconUser } from '@arco-design/web-vue/es/icon'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useUserStore } from '@/store/user'
import { usePermission } from '@/hooks/usePermission'
import { toast } from '@/feedback/toast'
import { createAIConversation, deleteAIConversation, getAIConfig, getAIMessages, listAIConversations, renameAIConversation, saveAIConfig, sendAIMessage, testAIConfig, type AIConfigInput, type AIConversation, type AIMessage } from '@/api/aiAssistant'

marked.setOptions({ gfm: true, breaks: true })
const user = useUserStore()
const { hasPerm } = usePermission()
const open = ref(false), settings = ref(false), fullscreen = ref(false), sending = ref(false), saving = ref(false)
const panel = ref<HTMLElement>(), scrollArea = ref<HTMLElement>()
const position = reactive({ left: 0, top: 0 })
const conversations = ref<AIConversation[]>([]), messages = ref<AIMessage[]>([]), activeId = ref(0), draft = ref(''), configError = ref('')
const renameVisible = ref(false), renameTitle = ref(''), renameTarget = ref<AIConversation>()
const today = new Date().toISOString().slice(0, 10)
const dateRange = ref([new Date(Date.now() - 30 * 86400000).toISOString().slice(0, 10), today])
const isHQAdmin = computed(() => {
  const roleCode = user.userInfo?.role?.code ?? ''
  return roleCode === 'super_admin' || (roleCode === 'admin' && Number(user.userInfo?.store_id ?? 0) === 0)
})
const canUseAssistant = computed(() => user.userInfo?.role?.code === 'store_admin' || isHQAdmin.value || hasPerm('ai:assistant:use'))
const storeLabel = computed(() => user.currentStoreId ? `门店 ${user.currentStoreId}` : '全部门店')
const effectiveStoreId = computed(() => isHQAdmin.value ? user.currentStoreId : Number(user.userInfo?.store_id ?? 0))
const activeConversation = computed(() => conversations.value.find(item => item.id === activeId.value))
const panelStyle = computed(() => fullscreen.value ? undefined : { left: `${position.left}px`, top: `${position.top}px` })
const config = reactive<AIConfigInput & { api_key_configured: boolean; api_key_masked: string }>({ provider: 'openai', model: 'gpt-4o-mini', base_url: '', enabled: false, system_prompt: '', api_key: '', api_key_configured: false, api_key_masked: '' })

let drag: { pointerId: number; x: number; y: number; left: number; top: number } | undefined
function placePanel() {
  if (!panel.value || fullscreen.value) return
  const rect = panel.value.getBoundingClientRect()
  position.left = Math.max(10, window.innerWidth - rect.width - 24)
  position.top = Math.max(10, window.innerHeight - rect.height - 24)
}
function openAssistant() { open.value = true; void nextTick(placePanel) }
function closeAssistant() { open.value = false; settings.value = false }
function toggleSettings() { settings.value = !settings.value }
function toggleFullscreen() { fullscreen.value = !fullscreen.value; if (!fullscreen.value) void nextTick(placePanel) }
function startDrag(event: PointerEvent) {
  if (fullscreen.value || event.button !== 0 || (event.target as HTMLElement).closest('button')) return
  drag = { pointerId: event.pointerId, x: event.clientX, y: event.clientY, left: position.left, top: position.top }
  ;(event.currentTarget as HTMLElement).setPointerCapture(event.pointerId)
}
function moveDrag(event: PointerEvent) {
  if (!drag || event.pointerId !== drag.pointerId || !panel.value) return
  const rect = panel.value.getBoundingClientRect()
  position.left = Math.max(0, Math.min(window.innerWidth - rect.width, drag.left + event.clientX - drag.x))
  position.top = Math.max(0, Math.min(window.innerHeight - rect.height, drag.top + event.clientY - drag.y))
}
function stopDrag(event: PointerEvent) { if (drag?.pointerId === event.pointerId) drag = undefined }
window.addEventListener('pointermove', moveDrag)
window.addEventListener('pointerup', stopDrag)
window.addEventListener('pointercancel', stopDrag)
window.addEventListener('resize', placePanel)
onBeforeUnmount(() => {
  window.removeEventListener('pointermove', moveDrag)
  window.removeEventListener('pointerup', stopDrag)
  window.removeEventListener('pointercancel', stopDrag)
  window.removeEventListener('resize', placePanel)
})

function renderMarkdown(content: string): string {
  return DOMPurify.sanitize(marked.parse(content) as string, { USE_PROFILES: { html: true } })
}
function formatAnalysisContext(value: string): string {
  try { return JSON.stringify(JSON.parse(value), null, 2) }
  catch { return value }
}
function configPayload(): AIConfigInput { return { provider: config.provider, model: config.model, base_url: config.base_url, enabled: config.enabled, system_prompt: config.system_prompt, api_key: config.api_key } }
async function loadConversations() {
  if (!isHQAdmin.value && !effectiveStoreId.value) { conversations.value = []; return }
  conversations.value = await listAIConversations(effectiveStoreId.value)
  if (!activeId.value && conversations.value.length) await selectConversation(conversations.value[0].id)
}
async function loadConfig() { if (!isHQAdmin.value) return; try { Object.assign(config, await getAIConfig()) } catch (error) { configError.value = errorMessage(error) } }
async function selectConversation(id: number) { activeId.value = id; messages.value = await getAIMessages(id, effectiveStoreId.value); scrollBottom() }
async function newConversation() {
  try { const row = await createAIConversation(effectiveStoreId.value); conversations.value.unshift(row); activeId.value = row.id; messages.value = []; settings.value = false }
  catch (error) { toast.error(errorMessage(error)) }
}
async function send() {
  const text = draft.value.trim()
  if (!text || sending.value || (!isHQAdmin.value && !effectiveStoreId.value)) return
  sending.value = true; draft.value = ''
  try {
    if (!activeId.value) { const row = await createAIConversation(effectiveStoreId.value); conversations.value.unshift(row); activeId.value = row.id }
    messages.value.push({ id: Date.now(), conversation_id: activeId.value, role: 'user', content: text, created_at: new Date().toISOString() }); scrollBottom()
    const answer = await sendAIMessage({ conversation_id: activeId.value, store_id: effectiveStoreId.value, message: text, start_date: dateRange.value[0], end_date: dateRange.value[1] })
    messages.value.push(answer); await loadConversations(); scrollBottom()
  } catch (error) { draft.value = text; toast.error(errorMessage(error)) }
  finally { sending.value = false }
}
function openRename(item: AIConversation) {
  renameTarget.value = item
  renameTitle.value = item.title
  renameVisible.value = true
}
async function renameConversation(): Promise<boolean> {
  const item = renameTarget.value
  const title = renameTitle.value.trim()
  if (!item || !title) return false
  try {
    await renameAIConversation(item.id, title, effectiveStoreId.value)
    item.title = title
    return true
  } catch (error) {
    toast.error(errorMessage(error))
    return false
  }
}
async function remove(item: AIConversation) {
  try {
    await deleteAIConversation(item.id, effectiveStoreId.value)
    conversations.value = conversations.value.filter(row => row.id !== item.id)
    if (activeId.value === item.id) { activeId.value = 0; messages.value = []; if (conversations.value.length) await selectConversation(conversations.value[0].id) }
  } catch (error) { toast.error(errorMessage(error)) }
}
function providerChanged() {
  const bases: Record<string, string> = { openai: 'https://api.openai.com/v1', deepseek: 'https://api.deepseek.com/v1', qwen: 'https://dashscope.aliyuncs.com/compatible-mode/v1', custom: '' }
  config.base_url = bases[config.provider] ?? ''
  if (config.provider === 'deepseek') config.model = 'deepseek-chat'
  if (config.provider === 'qwen') config.model = 'qwen-plus'
  if (config.provider === 'openai') config.model = 'gpt-4o-mini'
}
async function saveConfig() {
  saving.value = true; configError.value = ''
  try { Object.assign(config, await saveAIConfig(configPayload())); config.api_key = ''; toast.success('设置已保存') }
  catch (error) { configError.value = errorMessage(error) }
  finally { saving.value = false }
}
async function testConfig() {
  saving.value = true; configError.value = ''
  try { const result = await testAIConfig(configPayload()); toast.success(result.message) }
  catch (error) { configError.value = errorMessage(error) }
  finally { saving.value = false }
}
async function copyMessage(content: string) {
  try {
    await navigator.clipboard.writeText(content)
    toast.success('消息已复制')
  } catch {
    toast.error('复制失败，请检查浏览器剪贴板权限')
  }
}
function usePrompt(text: string) { draft.value = text }
function scrollBottom() { void nextTick(() => { if (scrollArea.value) scrollArea.value.scrollTop = scrollArea.value.scrollHeight }) }
function errorMessage(error: unknown) { return error instanceof Error ? error.message : '请求失败' }
watch(effectiveStoreId, () => { activeId.value = 0; messages.value = []; if (open.value) void loadConversations() })
watch(open, value => { if (value) { void loadConversations(); void loadConfig(); void nextTick(placePanel) } })
</script>

<style scoped>
.assistant-panel{width:min(860px,calc(100vw - 32px));height:min(720px,calc(100dvh - 32px));min-width:0;z-index:1}
.assistant-panel-fullscreen{inset:0;width:100vw;height:100dvh;max-width:none;border-radius:0}
.assistant-drag-handle{cursor:grab}.assistant-drag-handle:active{cursor:grabbing}
.assistant-launch{width:58px!important;height:58px!important;border:2px solid #fff!important;border-radius:50%!important;background:#0f766e!important;color:#fff!important;font-size:24px!important;box-shadow:0 8px 24px #0f172a55!important;transition:transform .2s,background-color .2s!important}.assistant-launch:hover{background:#115e59!important;transform:translateY(-3px)}.assistant-launch :deep(.arco-icon){color:#fff!important}
.assistant-icon-button{display:flex!important;width:34px!important;height:34px!important;align-items:center;justify-content:center;color:#536273!important;font-size:18px!important}.assistant-icon-button:hover{background:#f1f5f9!important;color:#0f766e!important}
.assistant-conversation-row{min-width:0;min-height:38px;border:1px solid transparent;transition:background-color .15s,border-color .15s}.assistant-conversation-row:hover{background:#eef2f5}.assistant-conversation-active{border-color:#b9ded8;background:#e7f4f1;color:#124e47}.assistant-conversation-active:hover{background:#e0f0ec}.assistant-conversation-title{min-width:0;border:0;background:transparent;color:inherit;font-size:12px;line-height:1.4}.assistant-row-action{width:28px!important;height:28px!important;min-width:28px!important;flex:none;opacity:1!important;color:#536273!important;font-size:14px!important}.assistant-row-action:hover{background:#dce6e8!important;color:#0f766e!important}.assistant-delete-action:hover{color:#be123c!important}.assistant-new-chat{height:38px!important;border-radius:6px!important;background:#0f766e!important;color:white!important;font-weight:600!important;box-shadow:0 1px 2px #0f172a1a!important}.assistant-new-chat:hover{background:#115e59!important}.assistant-new-chat :deep(.arco-icon){color:white!important}
.assistant-field{display:block;margin-top:16px;color:#475569;font-size:13px}.assistant-field :deep(.arco-select),.assistant-field :deep(.arco-input-wrapper),.assistant-field :deep(.arco-textarea-wrapper){display:flex;width:100%;margin-top:6px}
.assistant-composer:deep(.arco-textarea-wrapper),.assistant-composer:deep(.arco-textarea){border:0!important;background:transparent!important;box-shadow:none!important;resize:none!important}.assistant-composer :deep(textarea){padding:3px 1px!important}.assistant-composer-send{width:36px!important;height:36px!important;border:0!important;border-radius:8px!important;background:#0f766e!important;color:#fff!important;font-size:17px!important;box-shadow:0 2px 5px #0f766e44!important}.assistant-composer-send:hover{background:#115e59!important}.assistant-composer-send:disabled{background:#94a3b8!important;box-shadow:none!important}.assistant-composer-send :deep(.arco-icon){color:white!important}.assistant-message-tools{display:flex;margin-top:5px;min-height:28px}.assistant-copy-button{height:26px!important;padding:0 8px!important;color:#536273!important;font-size:11px!important}.assistant-copy-button:hover{background:#edf5f3!important;color:#0f766e!important}.assistant-copy-button :deep(.arco-icon){font-size:13px!important}
.assistant-thinking-dots{display:flex;gap:3px}.assistant-thinking-dots i{width:4px;height:4px;border-radius:50%;background:#0f766e;animation:assistant-pulse .8s infinite alternate}.assistant-thinking-dots i:nth-child(2){animation-delay:.2s}.assistant-thinking-dots i:nth-child(3){animation-delay:.4s}@keyframes assistant-pulse{to{opacity:.25;transform:translateY(-3px)}}
.assistant-markdown :deep(:first-child){margin-top:0}.assistant-markdown :deep(:last-child){margin-bottom:0}.assistant-markdown :deep(p){margin:0 0 14px;line-height:1.8}.assistant-markdown :deep(h1),.assistant-markdown :deep(h2),.assistant-markdown :deep(h3){margin:22px 0 9px;color:#172b36;font-weight:650;line-height:1.4}.assistant-markdown :deep(h1){font-size:20px}.assistant-markdown :deep(h2){font-size:17px}.assistant-markdown :deep(h3){font-size:15px}.assistant-markdown :deep(ul),.assistant-markdown :deep(ol){margin:8px 0 16px;padding-left:24px}.assistant-markdown :deep(li){padding-left:3px}.assistant-markdown :deep(li+li){margin-top:4px}.assistant-markdown :deep(strong){font-weight:650;color:#172b36}.assistant-markdown :deep(blockquote){margin:14px 0;border-left:3px solid #14b8a6;background:#f0fdfa;padding:8px 14px;color:#475569}.assistant-markdown :deep(code){border-radius:4px;background:#f1f5f9;padding:2px 5px;font-family:ui-monospace,SFMono-Regular,Menlo,monospace;font-size:.88em}.assistant-markdown :deep(pre){margin:12px 0 16px;overflow-x:auto;border:1px solid #263746;border-radius:8px;background:#14212b;padding:15px 16px;color:#e2e8f0;line-height:1.6}.assistant-markdown :deep(pre code){background:transparent;padding:0;color:inherit;font-size:12px}.assistant-markdown :deep(table){display:block;width:100%;overflow-x:auto;border-collapse:collapse;font-size:12px;line-height:1.55}.assistant-markdown :deep(.table-wrap){max-width:100%;margin:14px 0;overflow-x:auto;border:1px solid #dbe4e7;border-radius:7px}.assistant-markdown :deep(th){background:#f1f5f6;color:#334155;font-weight:600;text-align:left}.assistant-markdown :deep(th),.assistant-markdown :deep(td){border-bottom:1px solid #e2e8f0;padding:8px 10px;white-space:nowrap}.assistant-markdown :deep(tr:last-child td){border-bottom:0}.assistant-markdown :deep(a){color:#0f766e;text-decoration:underline;text-underline-offset:2px}.assistant-markdown :deep(hr){margin:20px 0;border:0;border-top:1px solid #e2e8f0}
.assistant-launch-enter-active,.assistant-launch-leave-active{transition:opacity .2s,transform .2s}.assistant-launch-enter-from,.assistant-launch-leave-to{opacity:0;transform:scale(.75) translateY(8px)}
.assistant-panel-enter-active,.assistant-panel-leave-active{transition:opacity .22s ease,transform .22s ease}.assistant-panel-enter-from,.assistant-panel-leave-to{opacity:0;transform:translateY(14px) scale(.98)}
.assistant-settings-enter-active,.assistant-settings-leave-active{transition:opacity .18s,transform .18s}.assistant-settings-enter-from,.assistant-settings-leave-to{opacity:0;transform:translateY(6px)}
@media(max-width:640px){.assistant-panel{width:calc(100vw - 16px);height:calc(100dvh - 16px);max-height:760px}.assistant-panel-fullscreen{width:100vw;height:100dvh}.assistant-sidebar{width:152px}.assistant-row-action{width:24px!important;min-width:24px!important;height:26px!important}.assistant-conversation-row{gap:0;padding-left:0;padding-right:0}.assistant-conversation-title{padding-left:6px;padding-right:3px;font-size:11px}.assistant-markdown{font-size:13px}.assistant-markdown :deep(table){font-size:11px}}
@media(prefers-reduced-motion:reduce){.assistant-launch-enter-active,.assistant-launch-leave-active,.assistant-panel-enter-active,.assistant-panel-leave-active,.assistant-settings-enter-active,.assistant-settings-leave-active{transition:none}}
</style>
