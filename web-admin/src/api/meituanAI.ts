import { http, unwrap } from './http'
import type {
  ApiEnvelope,
  MeituanAIAccount,
  MeituanAIDashboard,
  MeituanAIOrder,
  MeituanAIReview,
  MeituanAISuggestion,
  Paginated,
} from './types'

export async function listMeituanAIAccounts(): Promise<MeituanAIAccount[]> {
  const res = await http.get<ApiEnvelope<MeituanAIAccount[]>>('/meituan-ai/accounts')
  return unwrap(res)
}

export async function createMeituanAIAccount(body: Record<string, unknown>): Promise<MeituanAIAccount> {
  const res = await http.post<ApiEnvelope<MeituanAIAccount>>('/meituan-ai/accounts', body)
  return unwrap(res)
}

export async function updateMeituanAIAccount(id: number, body: Record<string, unknown>): Promise<void> {
  await http.put<ApiEnvelope<unknown>>(`/meituan-ai/accounts/${id}`, body)
}

export async function getMeituanAIDashboard(params?: Record<string, unknown>): Promise<MeituanAIDashboard> {
  const res = await http.get<ApiEnvelope<MeituanAIDashboard>>('/meituan-ai/dashboard', { params })
  return unwrap(res)
}

export async function importMeituanAIOrders(accountId: number, body: Record<string, unknown>): Promise<{ imported: number }> {
  const res = await http.post<ApiEnvelope<{ imported: number }>>(`/meituan-ai/accounts/${accountId}/orders/import`, body)
  return unwrap(res)
}

export async function syncMeituanAIOrders(accountId: number, body: FormData): Promise<{ imported: number; skipped: number }> {
  const res = await http.post<ApiEnvelope<{ imported: number; skipped: number }>>(`/meituan-ai/accounts/${accountId}/orders/sync`, body, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return unwrap(res)
}

export async function syncMeituanAIOpenAPIOrders(accountId: number, body: { order_id?: string; order_ids?: string[] }): Promise<{ imported: number; skipped: number }> {
  const res = await http.post<ApiEnvelope<{ imported: number; skipped: number }>>(`/meituan-ai/accounts/${accountId}/orders/sync-openapi`, body)
  return unwrap(res)
}

export async function importMeituanAIReviews(accountId: number, body: Record<string, unknown>): Promise<{ imported: number }> {
  const res = await http.post<ApiEnvelope<{ imported: number }>>(`/meituan-ai/accounts/${accountId}/reviews/import`, body)
  return unwrap(res)
}

export async function listMeituanAIOrders(params?: Record<string, unknown>): Promise<Paginated<MeituanAIOrder>> {
  const res = await http.get<ApiEnvelope<Paginated<MeituanAIOrder>>>('/meituan-ai/orders', { params })
  return unwrap(res)
}

export async function listMeituanAIReviews(params?: Record<string, unknown>): Promise<Paginated<MeituanAIReview>> {
  const res = await http.get<ApiEnvelope<Paginated<MeituanAIReview>>>('/meituan-ai/reviews', { params })
  return unwrap(res)
}

export async function generateMeituanAISuggestions(accountId: number, params?: Record<string, unknown>): Promise<{ generated: number; ai_enabled?: boolean; source?: string }> {
  const res = await http.post<ApiEnvelope<{ generated: number; ai_enabled?: boolean; source?: string }>>(`/meituan-ai/accounts/${accountId}/suggestions/generate`, undefined, { params })
  return unwrap(res)
}

export async function listMeituanAISuggestions(params?: Record<string, unknown>): Promise<Paginated<MeituanAISuggestion>> {
  const res = await http.get<ApiEnvelope<Paginated<MeituanAISuggestion>>>('/meituan-ai/suggestions', { params })
  return unwrap(res)
}

export async function updateMeituanAISuggestionStatus(id: number, status: MeituanAISuggestion['status']): Promise<void> {
  await http.put<ApiEnvelope<unknown>>(`/meituan-ai/suggestions/${id}/status`, { status })
}
