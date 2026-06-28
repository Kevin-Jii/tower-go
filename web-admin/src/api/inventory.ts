import { http, unwrap } from './http'
import type { InventoryOrder, InventoryRow, Paginated } from './types'
import { downloadBlob, filenameFromDisposition } from '@/utils/download'

export async function listInventories(params?: {
  page?: number
  page_size?: number
  store_id?: number
  product_id?: number
  keyword?: string
}): Promise<Paginated<InventoryRow>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<InventoryRow>>>('/inventories', { params })
  return unwrap(res)
}

export async function updateInventoryQuantity(
  id: number,
  body: { quantity: number; remark?: string },
): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/inventories/${id}`, body)
}

export async function listInventoryOrders(params?: {
  page?: number
  page_size?: number
  store_id?: number
  type?: number
  order_no?: string
  date?: string
}): Promise<Paginated<InventoryOrder>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<InventoryOrder>>>('/inventory-orders', {
    params,
  })
  return unwrap(res)
}

export async function exportInventoryOrders(params: { date: string; store_id?: number }): Promise<void> {
  const res = await http.get<Blob>('/inventory-orders/export', { params, responseType: 'blob' })
  const filename = filenameFromDisposition(res.headers['content-disposition'], `inventory-orders-${params.date}.xls`)
  downloadBlob(res.data, filename)
}

export async function getInventoryOrder(id: number): Promise<InventoryOrder> {
  const res = await http.get<import('./types').ApiEnvelope<InventoryOrder>>(`/inventory-orders/${id}`)
  return unwrap(res)
}

export async function createInventoryOrder(body: Record<string, unknown>): Promise<unknown> {
  const res = await http.post<import('./types').ApiEnvelope<unknown>>('/inventory-orders', body)
  return unwrap(res)
}
