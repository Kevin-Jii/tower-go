import { http, unwrap } from './http'
import type { MemberConsumptionPage, MemberGiftRecord, MemberRow, MemberWineStorage, MemberWineTransaction, Paginated } from './types'

export async function listMembers(params?: {
  page?: number
  page_size?: number
  keyword?: string
}): Promise<Paginated<MemberRow>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<MemberRow>>>('/members', { params })
  return unwrap(res)
}

export async function getMember(id: number): Promise<MemberRow> {
  const res = await http.get<import('./types').ApiEnvelope<MemberRow>>(`/members/${id}`)
  return unwrap(res)
}

export async function createMember(body: Record<string, unknown>): Promise<MemberRow> {
  const res = await http.post<import('./types').ApiEnvelope<MemberRow>>('/members', body)
  return unwrap(res)
}

export async function updateMember(id: number, body: Record<string, unknown>): Promise<MemberRow> {
  const res = await http.put<import('./types').ApiEnvelope<MemberRow>>(`/members/${id}`, body)
  return unwrap(res)
}

export async function deleteMember(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/members/${id}`)
}

export async function adjustMemberBalance(
  id: number,
  body: { amount: string; type: number; remark?: string; version: number },
): Promise<MemberRow> {
  const res = await http.post<import('./types').ApiEnvelope<MemberRow>>(`/members/${id}/adjust-balance`, body)
  return unwrap(res)
}

export async function listMemberConsumptions(
  id: number,
  params?: { start_date?: string; end_date?: string; page?: number; page_size?: number },
): Promise<MemberConsumptionPage> {
  const res = await http.get<import('./types').ApiEnvelope<MemberConsumptionPage>>(`/members/${id}/consumptions`, { params })
  return unwrap(res)
}

export async function listMemberGiftRecords(
  id: number,
  params?: { start_date?: string; end_date?: string; page?: number; page_size?: number },
): Promise<Paginated<MemberGiftRecord>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<MemberGiftRecord>>>(`/members/${id}/gift-records`, {
    params,
  })
  return unwrap(res)
}

export async function listMemberWineStorages(params?: {
  page?: number
  page_size?: number
  keyword?: string
  member_id?: number
  only_stock?: number
}): Promise<Paginated<MemberWineStorage>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<MemberWineStorage>>>('/member-wines', { params })
  return unwrap(res)
}

export async function depositMemberWine(body: {
  member_id: number
  wine_name: string
  unit?: string
  quantity: number
  remark?: string
}): Promise<MemberWineStorage> {
  const res = await http.post<import('./types').ApiEnvelope<MemberWineStorage>>('/member-wines/deposit', body)
  return unwrap(res)
}

export async function withdrawMemberWine(body: {
  member_id: number
  wine_name: string
  unit?: string
  quantity: number
  remark?: string
}): Promise<MemberWineStorage> {
  const res = await http.post<import('./types').ApiEnvelope<MemberWineStorage>>('/member-wines/withdraw', body)
  return unwrap(res)
}

export async function listMemberWineTransactions(params?: {
  page?: number
  page_size?: number
  storage_id?: number
  member_id?: number
  type?: number
  keyword?: string
  start_date?: string
  end_date?: string
}): Promise<Paginated<MemberWineTransaction>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<MemberWineTransaction>>>('/member-wines/transactions', {
    params,
  })
  return unwrap(res)
}
