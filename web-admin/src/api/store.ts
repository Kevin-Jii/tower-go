import { http, unwrap } from './http'
import type { Paginated, Store } from './types'

export async function listAllStores(): Promise<Store[]> {
  const res = await http.get<import('./types').ApiEnvelope<Store[]>>('/stores/all')
  return unwrap(res)
}

/** 门店列表（分页结构，当前后端一次返回全量） */
export async function listStores(): Promise<Paginated<Store>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<Store>>>('/stores')
  return unwrap(res)
}

export async function getStore(id: number): Promise<Store> {
  const res = await http.get<import('./types').ApiEnvelope<Store>>(`/stores/${id}`)
  return unwrap(res)
}

export async function createStore(body: Record<string, unknown>): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/stores', body)
}

export async function updateStore(id: number, body: Record<string, unknown>): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/stores/${id}`, body)
}

export async function deleteStore(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/stores/${id}`)
}

export async function bindStoreThirdPartyAccount(id: number, thirdPartyAccountId: number | null): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/stores/${id}/third-party-account`, {
    third_party_account_id: thirdPartyAccountId,
  })
}
