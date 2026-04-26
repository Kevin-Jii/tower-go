const TENANT_KEY = 'tower_tenant_id'
const TOKEN_KEY = 'tower_token'

export function getToken(): string {
  return localStorage.getItem(TOKEN_KEY) ?? ''
}

export function setToken(token: string): void {
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}

export function getStoredTenantId(): number | null {
  const v = localStorage.getItem(TENANT_KEY)
  if (v === null || v === '') return null
  const n = Number(v)
  return Number.isFinite(n) ? n : null
}

export function setStoredTenantId(id: number): void {
  localStorage.setItem(TENANT_KEY, String(id))
}

export function clearStoredTenantId(): void {
  localStorage.removeItem(TENANT_KEY)
}

/**
 * 切换租户：默认整页刷新；也可在业务里改为仅 invalidateQueries。
 */
export function applyTenantSwitch(
  tenantId: number,
  mode: 'reload' | 'soft' = 'reload',
): void {
  setStoredTenantId(tenantId)
  if (mode === 'reload') {
    window.location.reload()
  }
}
