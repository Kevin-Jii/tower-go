import { reactive } from 'vue'

interface ConfirmState {
  open: boolean
  title: string
  message: string
  resolve: ((ok: boolean) => void) | null
}

const state = reactive<ConfirmState>({
  open: false,
  title: '',
  message: '',
  resolve: null,
})

export const confirmState = state

export function confirmDialog(opts: { title?: string; message: string }): Promise<boolean> {
  return new Promise((resolve) => {
    state.title = opts.title ?? '确认'
    state.message = opts.message
    state.resolve = resolve
    state.open = true
  })
}

export function resolveConfirm(ok: boolean): void {
  state.open = false
  const r = state.resolve
  state.resolve = null
  r?.(ok)
}
