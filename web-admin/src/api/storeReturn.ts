import { http, unwrap } from './http'
import type { Paginated, StoreReturn, StoreReturnProduct, StoreReturnStats } from './types'

export async function listStoreReturns(params?: {
  page?: number
  page_size?: number
  store_id?: number
  keyword?: string
  start_date?: string
  end_date?: string
}): Promise<Paginated<StoreReturn>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<StoreReturn>>>('/store-returns', { params })
  return unwrap(res)
}

export async function getStoreReturn(id: number): Promise<StoreReturn> {
  const res = await http.get<import('./types').ApiEnvelope<StoreReturn>>(`/store-returns/${id}`)
  return unwrap(res)
}

export async function getStoreReturnStats(params?: {
  store_id?: number
  start_date?: string
  end_date?: string
}): Promise<StoreReturnStats> {
  const res = await http.get<import('./types').ApiEnvelope<StoreReturnStats>>('/store-returns/stats', { params })
  return unwrap(res)
}

export async function createStoreReturn(body: Record<string, unknown>): Promise<StoreReturn> {
  const res = await http.post<import('./types').ApiEnvelope<StoreReturn>>('/store-returns', body)
  return unwrap(res)
}

export async function updateStoreReturn(id: number, body: Record<string, unknown>): Promise<StoreReturn> {
  const res = await http.put<import('./types').ApiEnvelope<StoreReturn>>(`/store-returns/${id}`, body)
  return unwrap(res)
}

export async function deleteStoreReturn(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/store-returns/${id}`)
}

export async function listStoreReturnProducts(params?: {
  page?: number
  page_size?: number
  store_id?: number
  keyword?: string
  status?: number
}): Promise<Paginated<StoreReturnProduct>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<StoreReturnProduct>>>('/store-returns/products', { params })
  return unwrap(res)
}

export async function createStoreReturnProduct(body: Record<string, unknown>): Promise<StoreReturnProduct> {
  const res = await http.post<import('./types').ApiEnvelope<StoreReturnProduct>>('/store-returns/products', body)
  return unwrap(res)
}

export async function updateStoreReturnProduct(id: number, body: Record<string, unknown>): Promise<StoreReturnProduct> {
  const res = await http.put<import('./types').ApiEnvelope<StoreReturnProduct>>(`/store-returns/products/${id}`, body)
  return unwrap(res)
}

export async function deleteStoreReturnProduct(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/store-returns/products/${id}`)
}
