import { http, unwrap } from './http'
import type { SmsCampaign, SmsSendRecord, SmsServiceConfig } from './types'

export async function getSmsServiceConfig(): Promise<SmsServiceConfig> {
  const res = await http.get<import('./types').ApiEnvelope<SmsServiceConfig>>('/sms-campaigns/config')
  return unwrap(res)
}

export async function listSmsCampaigns(): Promise<SmsCampaign[]> {
  const res = await http.get<import('./types').ApiEnvelope<SmsCampaign[]>>('/sms-campaigns')
  return unwrap(res)
}

export async function getSmsCampaign(id: number): Promise<SmsCampaign> {
  const res = await http.get<import('./types').ApiEnvelope<SmsCampaign>>(`/sms-campaigns/${id}`)
  return unwrap(res)
}

export async function createSmsCampaign(body: Partial<SmsCampaign> & {
  name: string
  campaign_type: string
  template_code: string
  target_type: string
}): Promise<SmsCampaign> {
  const res = await http.post<import('./types').ApiEnvelope<SmsCampaign>>('/sms-campaigns', body)
  return unwrap(res)
}

export async function updateSmsCampaign(id: number, body: Partial<SmsCampaign>): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/sms-campaigns/${id}`, body)
}

export async function deleteSmsCampaign(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/sms-campaigns/${id}`)
}

export async function sendSmsCampaign(id: number): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>(`/sms-campaigns/${id}/send`)
}

export async function cancelSmsCampaign(id: number): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>(`/sms-campaigns/${id}/cancel`)
}

export async function listSmsCampaignRecords(id: number): Promise<SmsSendRecord[]> {
  const res = await http.get<import('./types').ApiEnvelope<SmsSendRecord[]>>(`/sms-campaigns/${id}/records`)
  return unwrap(res)
}
