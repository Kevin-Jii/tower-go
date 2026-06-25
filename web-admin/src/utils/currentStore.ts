import type { User } from '@/api/types'

const SWITCHABLE_ROLE_CODES = new Set(['admin', 'super_admin'])

function positiveId(value: unknown): number {
  const id = Number(value ?? 0)
  return Number.isFinite(id) && id > 0 ? id : 0
}

export function isTenantSwitchableUser(user: User | null | undefined): boolean {
  const roleCode = user?.role?.code ?? ''
  return SWITCHABLE_ROLE_CODES.has(roleCode) || positiveId(user?.store_id) === 0
}

export function resolveCurrentStoreId(user: User | null | undefined, tenantId: unknown): number {
  const storeId = positiveId(user?.store_id)
  if (storeId > 0 && !isTenantSwitchableUser(user)) return storeId

  const selectedTenantId = positiveId(tenantId)
  return selectedTenantId || storeId
}

export function currentStoreIdOrUndefined(user: User | null | undefined, tenantId: unknown): number | undefined {
  return resolveCurrentStoreId(user, tenantId) || undefined
}
