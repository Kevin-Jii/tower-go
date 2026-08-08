import { http, unwrap } from './http'
import type { ApiEnvelope, Paginated, PreOrder } from './types'

export interface PreOrderItemPayload {
  product_id: number
  unit_spec_id: number
  quantity: number
  remark?: string
}

export interface PreOrderPayload {
  store_id?: number
  customer_id: number
  scheduled_at: string
  contact_person?: string
  contact_phone?: string
  delivery_address?: string
  remark?: string
  items: PreOrderItemPayload[]
}

export async function listPreOrders(params?: Record<string, unknown>): Promise<Paginated<PreOrder>> {
  const res = await http.get<ApiEnvelope<Paginated<PreOrder>>>('/pre-orders', { params })
  return unwrap(res)
}

export async function getPreOrder(id: number): Promise<PreOrder> {
  const res = await http.get<ApiEnvelope<PreOrder>>(`/pre-orders/${id}`)
  return unwrap(res)
}

export async function createPreOrder(body: PreOrderPayload): Promise<PreOrder> {
  const res = await http.post<ApiEnvelope<PreOrder>>('/pre-orders', body)
  return unwrap(res)
}

export async function updatePreOrder(id: number, body: PreOrderPayload): Promise<PreOrder> {
  const res = await http.put<ApiEnvelope<PreOrder>>(`/pre-orders/${id}`, body)
  return unwrap(res)
}

export async function updatePreOrderStatus(id: number, status: number): Promise<PreOrder> {
  const res = await http.put<ApiEnvelope<PreOrder>>(`/pre-orders/${id}/status`, { status })
  return unwrap(res)
}

export async function deletePreOrder(id: number): Promise<void> {
  await http.delete<ApiEnvelope<unknown>>(`/pre-orders/${id}`)
}
