<template>
  <!-- Toast -->
  <Teleport to="body">
    <div class="pointer-events-none fixed inset-x-0 top-3 z-[100] flex flex-col items-center gap-2 px-3">
      <TransitionGroup name="toast">
        <div
          v-for="t in toast.state.items"
          :key="t.id"
          class="pointer-events-auto max-w-[min(92vw,28rem)] rounded-xl px-4 py-2.5 text-sm font-medium shadow-lg border backdrop-blur-sm"
          :class="toastClass(t.level)"
        >
          {{ t.message }}
        </div>
      </TransitionGroup>
    </div>
  </Teleport>

  <!-- Confirm -->
  <Teleport to="body">
    <div
      v-if="confirmState.open"
      class="fixed inset-0 z-[99] flex items-center justify-center bg-slate-900/40 p-4 backdrop-blur-[2px]"
      role="dialog"
      aria-modal="true"
      @click.self="onCancel"
      @keydown.esc.prevent="onCancel"
    >
      <div
        class="w-full max-w-md rounded-2xl border border-slate-200/80 bg-white p-5 shadow-2xl"
        @click.stop
      >
        <h3 class="m-0 text-base font-semibold text-slate-900">{{ confirmState.title }}</h3>
        <p class="mt-2 mb-5 text-sm text-slate-600 leading-relaxed whitespace-pre-wrap">{{ confirmState.message }}</p>
        <div class="flex justify-end gap-2">
          <BaseButton variant="ghost" @click="onCancel">取消</BaseButton>
          <BaseButton variant="primary" @click="onOk">确定</BaseButton>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import BaseButton from '@/components/base/BaseButton.vue'
import { confirmState, resolveConfirm } from '@/feedback/confirm'
import { toast, type ToastLevel } from '@/feedback/toast'

function toastClass(level: ToastLevel): string {
  const m: Record<ToastLevel, string> = {
    success: 'bg-emerald-50 text-emerald-900 border-emerald-200/80',
    error: 'bg-rose-50 text-rose-900 border-rose-200/80',
    info: 'bg-sky-50 text-sky-900 border-sky-200/80',
    warning: 'bg-amber-50 text-amber-950 border-amber-200/80',
  }
  return m[level]
}

function onOk(): void {
  resolveConfirm(true)
}

function onCancel(): void {
  resolveConfirm(false)
}

function onKey(e: KeyboardEvent): void {
  if (!confirmState.open) return
  if (e.key === 'Escape') onCancel()
}

onMounted(() => window.addEventListener('keydown', onKey))
onUnmounted(() => window.removeEventListener('keydown', onKey))
</script>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.22s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
