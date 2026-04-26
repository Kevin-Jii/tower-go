import { http, unwrap } from './http'
import type { Paginated, Supplier } from './types'

export async function listSuppliers(params?: {
  page?: number
  page_size?: number
  keyword?: string
  status?: number
}): Promise<Paginated<Supplier>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<Supplier>>>('/suppliers', { params })
  return unwrap(res)
}

export async function getSupplier(id: number): Promise<Supplier> {
  const res = await http.get<import('./types').ApiEnvelope<Supplier>>(`/suppliers/${id}`)
  return unwrap(res)
}

export async function createSupplier(body: Record<string, unknown>): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/suppliers', body)
}

export async function updateSupplier(id: number, body: Record<string, unknown>): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/suppliers/${id}`, body)
}

export async function deleteSupplier(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/suppliers/${id}`)
}
