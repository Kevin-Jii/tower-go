import { http, unwrap } from './http'
import type { AssignMenusToRoleReq, Menu } from './types'

export async function fetchUserMenus(): Promise<Menu[]> {
  const res = await http.get<import('./types').ApiEnvelope<Menu[]>>('/menus/user-menus')
  return unwrap(res)
}

export async function fetchUserPermissions(): Promise<string[]> {
  const res = await http.get<import('./types').ApiEnvelope<string[]>>('/menus/user-permissions')
  return unwrap(res)
}

export async function fetchMenuTree(): Promise<Menu[]> {
  const res = await http.get<import('./types').ApiEnvelope<Menu[]>>('/menus/tree')
  return unwrap(res)
}

export async function createMenu(body: Record<string, unknown>): Promise<Menu> {
  const res = await http.post<import('./types').ApiEnvelope<Menu>>('/menus', body)
  return unwrap(res)
}

export async function updateMenu(id: number, body: Record<string, unknown>): Promise<Menu> {
  const res = await http.put<import('./types').ApiEnvelope<Menu>>(`/menus/${id}`, body)
  return unwrap(res)
}

export async function deleteMenu(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/menus/${id}`)
}

export async function fetchRoleMenuIds(roleId: number): Promise<number[]> {
  const res = await http.get<import('./types').ApiEnvelope<number[]>>('/menus/role-ids', {
    params: { role_id: roleId },
  })
  return unwrap(res)
}

export async function assignRoleMenus(body: AssignMenusToRoleReq): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/menus/assign-role', body)
}

/** 角色在各菜单上的权限位（用于分配菜单回显） */
export async function fetchRoleMenuPermissions(roleId: number): Promise<Record<number, number>> {
  const res = await http.get<import('./types').ApiEnvelope<Record<string, number>>>(
    '/menus/role-permissions',
    { params: { role_id: roleId } },
  )
  const raw = unwrap(res) ?? {}
  const out: Record<number, number> = {}
  for (const [k, v] of Object.entries(raw)) {
    out[Number(k)] = v
  }
  return out
}
