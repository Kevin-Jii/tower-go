import { http, unwrap } from './http'
import type { Paginated, StoreAccount, StoreAccountStats } from './types'

export async function listStoreAccounts(params?: {
  page?: number
  page_size?: number
  store_id?: number
  channel?: string
  order_no?: string
  payment_status?: number
  member_keyword?: string
  tag_code?: string
  start_date?: string
  end_date?: string
}): Promise<Paginated<StoreAccount>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<StoreAccount>>>('/store-accounts', { params })
  return unwrap(res)
}

export async function getStoreAccount(id: number): Promise<StoreAccount> {
  const res = await http.get<import('./types').ApiEnvelope<StoreAccount>>(`/store-accounts/${id}`)
  return unwrap(res)
}

export async function getStoreAccountStats(params?: {
  store_id?: number
  start_date?: string
  end_date?: string
}): Promise<StoreAccountStats> {
  const res = await http.get<import('./types').ApiEnvelope<StoreAccountStats>>('/store-accounts/stats', { params })
  return unwrap(res)
}

export async function createStoreAccount(body: Record<string, unknown>): Promise<StoreAccount> {
  const res = await http.post<import('./types').ApiEnvelope<StoreAccount>>('/store-accounts', body)
  return unwrap(res)
}

export async function bindStoreAccountConsumables(id: number, body: { consumables: Array<Record<string, unknown>> }): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>(`/store-accounts/${id}/consumables`, body)
}

export async function updateStoreAccount(id: number, body: Record<string, unknown>): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/store-accounts/${id}`, body)
}

export async function deleteStoreAccount(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/store-accounts/${id}`)
}
