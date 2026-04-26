import { http, unwrap } from './http'
import type { Role } from './types'

export async function listRoles(params?: { keyword?: string; status?: number }): Promise<Role[]> {
  const res = await http.get<import('./types').ApiEnvelope<Role[]>>('/roles', { params })
  return unwrap(res)
}

export async function createRole(body: Partial<Role>): Promise<Role> {
  const res = await http.post<import('./types').ApiEnvelope<Role>>('/roles', body)
  return unwrap(res)
}

export async function updateRole(id: number, body: Record<string, unknown>): Promise<Role> {
  const res = await http.put<import('./types').ApiEnvelope<Role>>(`/roles/${id}`, body)
  return unwrap(res)
}

export async function deleteRole(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/roles/${id}`)
}
