import { http, unwrap } from './http'
import type { ThirdPartyAccount, ThirdPartyOrder, Paginated } from './types'

export async function listThirdPartyAccounts(keyword?: string): Promise<ThirdPartyAccount[]> {
  const res = await http.get<import('./types').ApiEnvelope<ThirdPartyAccount[]>>('/third-party-accounts', {
    params: { keyword: keyword || undefined },
  })
  return unwrap(res)
}

export async function getThirdPartyAccount(id: number): Promise<ThirdPartyAccount> {
  const res = await http.get<import('./types').ApiEnvelope<ThirdPartyAccount>>(`/third-party-accounts/${id}`)
  return unwrap(res)
}

export async function createThirdPartyAccount(body: Record<string, unknown>): Promise<ThirdPartyAccount> {
  const res = await http.post<import('./types').ApiEnvelope<ThirdPartyAccount>>('/third-party-accounts', body)
  return unwrap(res)
}

export async function updateThirdPartyAccount(id: number, body: Record<string, unknown>): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/third-party-accounts/${id}`, body)
}

export async function deleteThirdPartyAccount(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/third-party-accounts/${id}`)
}

export async function testThirdPartyAccountLogin(id: number): Promise<Record<string, unknown>> {
  const res = await http.post<import('./types').ApiEnvelope<Record<string, unknown>>>(`/third-party-accounts/${id}/test-login`)
  return unwrap(res)
}

export async function syncThirdPartyLatestOrders(id: number): Promise<Record<string, unknown>> {
  const res = await http.post<import('./types').ApiEnvelope<Record<string, unknown>>>(`/third-party-accounts/${id}/sync-latest-orders`)
  return unwrap(res)
}

export async function listThirdPartySyncedOrders(id: number, page = 1, pageSize = 20): Promise<Paginated<ThirdPartyOrder>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<ThirdPartyOrder>>>(`/third-party-accounts/${id}/orders`, {
    params: { page, page_size: pageSize },
  })
  return unwrap(res)
}

