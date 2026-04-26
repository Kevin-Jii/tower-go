import { computed } from 'vue'
import { useUserStore } from '@/store/user'

export function usePermission() {
  const userStore = useUserStore()

  const codes = computed(() => new Set(userStore.permissions))

  function hasPerm(code: string): boolean {
    if (!code) return true
    return codes.value.has(code)
  }

  return { hasPerm, codes }
}
