import { http, unwrap } from './http'
import type { ApiEnvelope, B2BCustomer, B2BCustomerProductPrice, B2BSupplyOrder, Paginated } from './types'
import { downloadBlob, filenameFromDisposition } from '@/utils/download'

export async function listB2BCustomers(params?: Record<string, unknown>): Promise<Paginated<B2BCustomer>> {
  const res = await http.get<ApiEnvelope<Paginated<B2BCustomer>>>('/b2b/customers', { params })
  return unwrap(res)
}

export async function createB2BCustomer(body: Record<string, unknown>): Promise<B2BCustomer> {
  const res = await http.post<ApiEnvelope<B2BCustomer>>('/b2b/customers', body)
  return unwrap(res)
}

export async function updateB2BCustomer(id: number, body: Record<string, unknown>): Promise<B2BCustomer> {
  const res = await http.put<ApiEnvelope<B2BCustomer>>(`/b2b/customers/${id}`, body)
  return unwrap(res)
}

export async function listB2BPrices(params?: Record<string, unknown>): Promise<Paginated<B2BCustomerProductPrice>> {
  const res = await http.get<ApiEnvelope<Paginated<B2BCustomerProductPrice>>>('/b2b/prices', { params })
  return unwrap(res)
}

export async function upsertB2BPrice(body: Record<string, unknown>): Promise<void> {
  await http.post<ApiEnvelope<unknown>>('/b2b/prices', body)
}

export async function deleteB2BPrice(id: number): Promise<void> {
  await http.delete<ApiEnvelope<unknown>>(`/b2b/prices/${id}`)
}

export async function listB2BSupplyOrders(params?: Record<string, unknown>): Promise<Paginated<B2BSupplyOrder>> {
  const res = await http.get<ApiEnvelope<Paginated<B2BSupplyOrder>>>('/b2b/supply-orders', { params })
  return unwrap(res)
}

export async function exportB2BSupplyOrders(params: { date: string; store_id?: number }): Promise<void> {
  const res = await http.get<Blob>('/b2b/supply-orders/export', { params, responseType: 'blob' })
  const filename = filenameFromDisposition(res.headers['content-disposition'], `b2b-supply-orders-${params.date}.xls`)
  downloadBlob(res.data, filename)
}

export async function createB2BSupplyOrder(body: Record<string, unknown>): Promise<B2BSupplyOrder> {
  const res = await http.post<ApiEnvelope<B2BSupplyOrder>>('/b2b/supply-orders', body)
  return unwrap(res)
}

export async function getB2BSupplyOrder(id: number): Promise<B2BSupplyOrder> {
  const res = await http.get<ApiEnvelope<B2BSupplyOrder>>(`/b2b/supply-orders/${id}`)
  return unwrap(res)
}

export async function updateB2BSupplyOrderDeliveryStatus(
  id: number,
  body: { delivery_status: number },
): Promise<B2BSupplyOrder> {
  const res = await http.put<ApiEnvelope<B2BSupplyOrder>>(`/b2b/supply-orders/${id}/delivery-status`, body)
  return unwrap(res)
}

export async function updateB2BSupplyOrderPaymentStatus(
  id: number,
  body: { payment_status: number; paid_amount?: number },
): Promise<B2BSupplyOrder> {
  const res = await http.put<ApiEnvelope<B2BSupplyOrder>>(`/b2b/supply-orders/${id}/payment-status`, body)
  return unwrap(res)
}
