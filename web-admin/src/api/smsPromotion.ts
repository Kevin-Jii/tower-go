import { http, unwrap } from './http'
import type {
  MemberRow,
  MemberTag,
  MemberTagBinding,
  SmsCampaign,
  SmsCampaignPayload,
  SmsSendRecord,
  SmsServiceConfig,
} from './types'

type ApiEnvelope<T> = import('./types').ApiEnvelope<T>

export async function getSmsServiceConfig(): Promise<SmsServiceConfig> {
  const res = await http.get<ApiEnvelope<SmsServiceConfig>>('/sms-campaigns/config')
  return unwrap(res)
}

export async function listSmsCampaigns(params?: { store_id?: number }): Promise<SmsCampaign[]> {
  const res = await http.get<ApiEnvelope<SmsCampaign[]>>('/sms-campaigns', { params })
  return unwrap(res)
}

export async function getSmsCampaign(id: number, storeId?: number): Promise<SmsCampaign> {
  const res = await http.get<ApiEnvelope<SmsCampaign>>(`/sms-campaigns/${id}`, { params: storeId ? { store_id: storeId } : undefined })
  return unwrap(res)
}

export async function createSmsCampaign(body: SmsCampaignPayload): Promise<SmsCampaign> {
  const res = await http.post<ApiEnvelope<SmsCampaign>>('/sms-campaigns', body)
  return unwrap(res)
}

export async function updateSmsCampaign(id: number, body: SmsCampaignPayload, storeId?: number): Promise<void> {
  await http.put<ApiEnvelope<unknown>>(`/sms-campaigns/${id}`, body, { params: storeId ? { store_id: storeId } : undefined })
}

export async function deleteSmsCampaign(id: number, storeId?: number): Promise<void> {
  await http.delete<ApiEnvelope<unknown>>(`/sms-campaigns/${id}`, { params: storeId ? { store_id: storeId } : undefined })
}

export async function sendSmsCampaign(id: number, storeId?: number): Promise<void> {
  await http.post<ApiEnvelope<unknown>>(`/sms-campaigns/${id}/send`, undefined, {
    params: storeId ? { store_id: storeId } : undefined,
    timeout: 10 * 60_000,
  })
}

export async function cancelSmsCampaign(id: number, storeId?: number): Promise<void> {
  await http.post<ApiEnvelope<unknown>>(`/sms-campaigns/${id}/cancel`, undefined, { params: storeId ? { store_id: storeId } : undefined })
}

export async function listSmsCampaignRecords(id: number, storeId?: number): Promise<SmsSendRecord[]> {
  const res = await http.get<ApiEnvelope<SmsSendRecord[]>>(`/sms-campaigns/${id}/records`, { params: storeId ? { store_id: storeId } : undefined })
  return unwrap(res)
}

export async function listMemberTags(params?: { store_id?: number }): Promise<MemberTag[]> {
  const res = await http.get<ApiEnvelope<MemberTag[]>>('/member-tags', { params })
  return unwrap(res)
}

export async function createMemberTag(body: {
  store_id?: number
  name: string
  color?: string
  description?: string
}): Promise<MemberTag> {
  const res = await http.post<ApiEnvelope<MemberTag>>('/member-tags', body)
  return unwrap(res)
}

export async function updateMemberTag(
  id: number,
  body: { name: string; color?: string; description?: string },
  storeId?: number,
): Promise<MemberTag> {
  const res = await http.put<ApiEnvelope<MemberTag>>(`/member-tags/${id}`, body, { params: storeId ? { store_id: storeId } : undefined })
  return unwrap(res)
}

export async function deleteMemberTag(id: number, storeId?: number): Promise<void> {
  await http.delete<ApiEnvelope<unknown>>(`/member-tags/${id}`, { params: storeId ? { store_id: storeId } : undefined })
}

export async function searchSmsMembers(params?: { store_id?: number; keyword?: string; limit?: number }): Promise<MemberRow[]> {
  const res = await http.get<ApiEnvelope<MemberRow[]>>('/member-tags/members/search', { params })
  return unwrap(res)
}

export async function listMemberTagMembers(id: number, storeId?: number): Promise<MemberRow[]> {
  const res = await http.get<ApiEnvelope<MemberRow[]>>(`/member-tags/${id}/members`, { params: storeId ? { store_id: storeId } : undefined })
  return unwrap(res)
}

export async function bindMemberTag(id: number, memberId: number, storeId?: number): Promise<MemberTagBinding> {
  const res = await http.post<ApiEnvelope<MemberTagBinding>>(`/member-tags/${id}/members`, {
    member_id: memberId,
  }, { params: storeId ? { store_id: storeId } : undefined })
  return unwrap(res)
}

export async function unbindMemberTag(id: number, memberId: number, storeId?: number): Promise<void> {
  await http.delete<ApiEnvelope<unknown>>(`/member-tags/${id}/members/${memberId}`, { params: storeId ? { store_id: storeId } : undefined })
}

export async function listMemberTagsForMember(memberId: number): Promise<MemberTag[]> {
  const res = await http.get<ApiEnvelope<MemberTag[]>>(`/members/${memberId}/tags`)
  return unwrap(res)
}

export async function assignMemberTags(memberId: number, tagIds: number[]): Promise<void> {
  await http.put<ApiEnvelope<unknown>>(`/members/${memberId}/tags`, { tag_ids: tagIds })
}
