import { http, unwrap } from './http'
import type { ApiEnvelope, Paginated, StoreExpense, StoreExpenseStats } from './types'
import { downloadBlob, filenameFromDisposition } from '@/utils/download'

export async function listStoreExpenses(params?: Record<string, unknown>): Promise<Paginated<StoreExpense>> {
  const res = await http.get<ApiEnvelope<Paginated<StoreExpense>>>('/store-expenses', { params })
  return unwrap(res)
}

export async function exportStoreExpenses(params: { date: string; store_id?: number }): Promise<void> {
  const res = await http.get<Blob>('/store-expenses/export', { params, responseType: 'blob' })
  const filename = filenameFromDisposition(res.headers['content-disposition'], `store-expenses-${params.date}.xls`)
  downloadBlob(res.data, filename)
}

export async function getStoreExpense(id: number): Promise<StoreExpense> {
  const res = await http.get<ApiEnvelope<StoreExpense>>(`/store-expenses/${id}`)
  return unwrap(res)
}

export async function getStoreExpenseStats(params?: Record<string, unknown>): Promise<StoreExpenseStats> {
  const res = await http.get<ApiEnvelope<StoreExpenseStats>>('/store-expenses/stats', { params })
  return unwrap(res)
}

export async function createStoreExpense(body: Record<string, unknown>): Promise<StoreExpense> {
  const res = await http.post<ApiEnvelope<StoreExpense>>('/store-expenses', body)
  return unwrap(res)
}

export async function updateStoreExpense(id: number, body: Record<string, unknown>): Promise<StoreExpense> {
  const res = await http.put<ApiEnvelope<StoreExpense>>(`/store-expenses/${id}`, body)
  return unwrap(res)
}

export async function deleteStoreExpense(id: number): Promise<void> {
  await http.delete<ApiEnvelope<unknown>>(`/store-expenses/${id}`)
}
