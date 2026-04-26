import { http, unwrap } from './http'
import type { Paginated, User } from './types'

export async function listUsers(params: {
  page?: number
  page_size?: number
  keyword?: string
}): Promise<Paginated<User>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<User>>>('/users', { params })
  return unwrap(res)
}

export async function createUser(body: Record<string, unknown>): Promise<User> {
  const res = await http.post<import('./types').ApiEnvelope<User>>('/users', body)
  return unwrap(res)
}

export async function updateUser(id: number, body: Record<string, unknown>): Promise<User> {
  const res = await http.put<import('./types').ApiEnvelope<User>>(`/users/${id}`, body)
  return unwrap(res)
}

export async function deleteUser(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/users/${id}`)
}

export async function assignUserRole(userId: number, roleId: number): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/permission/assign-user-role', {
    user_id: userId,
    role_id: roleId,
  })
}
