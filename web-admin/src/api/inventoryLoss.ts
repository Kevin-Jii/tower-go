import { http, unwrap } from './http'
import type { InventoryLossOrder, InventoryLossOrderDetail, MemberGiftRecord, Paginated } from './types'

export async function listInventoryLossOrders(params?: {
  page?: number
  page_size?: number
  store_id?: number
  type?: string
  member_id?: number
  start_date?: string
  end_date?: string
  keyword?: string
}): Promise<Paginated<InventoryLossOrder>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<InventoryLossOrder>>>('/inventory-loss-orders', {
    params,
  })
  return unwrap(res)
}

export async function createInventoryLossOrder(body: Record<string, unknown>): Promise<InventoryLossOrderDetail> {
  const res = await http.post<import('./types').ApiEnvelope<InventoryLossOrderDetail>>('/inventory-loss-orders', body)
  return unwrap(res)
}

export async function getInventoryLossOrder(id: number): Promise<InventoryLossOrderDetail> {
  const res = await http.get<import('./types').ApiEnvelope<InventoryLossOrderDetail>>(`/inventory-loss-orders/${id}`)
  return unwrap(res)
}

export async function updateInventoryLossOrder(id: number, body: { reason: string }): Promise<InventoryLossOrderDetail> {
  const res = await http.put<import('./types').ApiEnvelope<InventoryLossOrderDetail>>(
    `/inventory-loss-orders/${id}`,
    body,
  )
  return unwrap(res)
}

export async function cancelInventoryLossOrder(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/inventory-loss-orders/${id}`)
}

export async function listMemberGiftRecords(
  memberId: number,
  params?: { page?: number; page_size?: number; start_date?: string; end_date?: string },
): Promise<Paginated<MemberGiftRecord>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<MemberGiftRecord>>>(
    `/members/${memberId}/gift-records`,
    { params },
  )
  return unwrap(res)
}
