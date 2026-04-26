import type { Directive, DirectiveBinding } from 'vue'
import { useUserStore } from '@/store/user'

function resolveCodes(value: DirectiveBinding['value']): string[] {
  if (value == null) return []
  if (Array.isArray(value)) return value.filter(Boolean).map(String)
  return [String(value)]
}

export const permissionDirective: Directive<HTMLElement, string | string[]> = {
  mounted(el, binding) {
    const store = useUserStore()
    const need = resolveCodes(binding.value)
    if (need.length === 0) return
    const ok = need.some((c) => store.permissions.includes(c))
    if (!ok) el.parentNode?.removeChild(el)
  },
  updated(el, binding) {
    const store = useUserStore()
    const need = resolveCodes(binding.value)
    if (need.length === 0) return
    const ok = need.some((c) => store.permissions.includes(c))
    el.style.display = ok ? '' : 'none'
  },
}
