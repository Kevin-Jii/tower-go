import { reactive } from 'vue'

export type ToastLevel = 'success' | 'error' | 'info' | 'warning'

export interface ToastItem {
  id: number
  level: ToastLevel
  message: string
}

const state = reactive({
  items: [] as ToastItem[],
})

let seq = 0

function push(level: ToastLevel, message: string, duration = 3200): void {
  const id = ++seq
  state.items.push({ id, level, message })
  window.setTimeout(() => {
    state.items = state.items.filter((t) => t.id !== id)
  }, duration)
}

export const toast = {
  state,
  success: (msg: string) => push('success', msg),
  error: (msg: string) => push('error', msg, 4200),
  info: (msg: string) => push('info', msg),
  warning: (msg: string) => push('warning', msg),
}
