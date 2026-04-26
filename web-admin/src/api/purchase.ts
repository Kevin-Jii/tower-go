import { http, unwrap } from './http'
import type { Paginated, PurchaseOrder } from './types'

export { listPurchasableProducts } from './storeSupplier'

export async function listPurchaseOrders(params?: {
  page?: number
  page_size?: number
  store_id?: number
  supplier_id?: number
  status?: number
  date?: string
}): Promise<Paginated<PurchaseOrder>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<PurchaseOrder>>>('/purchase-orders', { params })
  return unwrap(res)
}

export async function getPurchaseOrder(id: number): Promise<PurchaseOrder> {
  const res = await http.get<import('./types').ApiEnvelope<PurchaseOrder>>(`/purchase-orders/${id}`)
  return unwrap(res)
}

export async function createPurchaseOrder(body: Record<string, unknown>): Promise<PurchaseOrder> {
  const res = await http.post<import('./types').ApiEnvelope<PurchaseOrder>>('/purchase-orders', body)
  return unwrap(res)
}

export async function updatePurchaseOrder(id: number, body: Record<string, unknown>): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/purchase-orders/${id}`, body)
}

export async function deletePurchaseOrder(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/purchase-orders/${id}`)
}

export async function getPurchaseOrderActions(id: number): Promise<string[]> {
  const res = await http.get<import('./types').ApiEnvelope<string[]>>(`/purchase-orders/${id}/actions`)
  return unwrap(res)
}

export async function confirmPurchaseOrder(id: number): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>(`/purchase-orders/${id}/confirm`)
}

export async function completePurchaseOrder(id: number): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>(`/purchase-orders/${id}/complete`)
}

export async function cancelPurchaseOrder(id: number): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>(`/purchase-orders/${id}/cancel`)
}
