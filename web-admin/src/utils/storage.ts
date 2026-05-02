const TENANT_KEY = 'tower_tenant_id'
const TOKEN_KEY = 'tower_token'
const LOGIN_REMEMBER_KEY = 'tower_login_remember'

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

export interface RememberedLogin {
  phone: string
  password: string
}

export function getRememberedLogin(): RememberedLogin | null {
  try {
    const raw = localStorage.getItem(LOGIN_REMEMBER_KEY)
    if (!raw) return null
    const o = JSON.parse(raw) as { phone?: unknown; password?: unknown }
    if (typeof o.phone !== 'string' || typeof o.password !== 'string') return null
    return { phone: o.phone, password: o.password }
  } catch {
    return null
  }
}

export function setRememberedLogin(phone: string, password: string): void {
  localStorage.setItem(LOGIN_REMEMBER_KEY, JSON.stringify({ phone, password }))
}

export function clearRememberedLogin(): void {
  localStorage.removeItem(LOGIN_REMEMBER_KEY)
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
