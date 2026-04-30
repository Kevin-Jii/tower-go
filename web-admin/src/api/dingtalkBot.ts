import { http, unwrap } from './http'
import type { CreateDingTalkBotReq, DingTalkBot, Paginated, UpdateDingTalkBotReq } from './types'

export async function listDingTalkBots(params?: {
  page?: number
  page_size?: number
}): Promise<Paginated<DingTalkBot>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<DingTalkBot>>>('/dingtalk/robots', { params })
  return unwrap(res)
}

export async function getDingTalkBot(id: number): Promise<DingTalkBot> {
  const res = await http.get<import('./types').ApiEnvelope<DingTalkBot>>(`/dingtalk/robots/${id}`)
  return unwrap(res)
}

export async function createDingTalkBot(body: CreateDingTalkBotReq): Promise<DingTalkBot> {
  const res = await http.post<import('./types').ApiEnvelope<DingTalkBot>>('/dingtalk/robots', body)
  return unwrap(res)
}

export async function updateDingTalkBot(id: number, body: UpdateDingTalkBotReq): Promise<DingTalkBot> {
  const res = await http.put<import('./types').ApiEnvelope<DingTalkBot>>(`/dingtalk/robots/${id}`, body)
  return unwrap(res)
}

export async function deleteDingTalkBot(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/dingtalk/robots/${id}`)
}

export async function testDingTalkBot(id: number): Promise<{ message?: string; robot_code?: string }> {
  const res = await http.post<import('./types').ApiEnvelope<{ message?: string; robot_code?: string }>>(
    `/dingtalk/robots/${id}/test`,
  )
  return unwrap(res)
}

export async function testDingTalkStreamCallback(id: number): Promise<Record<string, unknown>> {
  const res = await http.post<import('./types').ApiEnvelope<Record<string, unknown>>>(
    `/dingtalk/robots/${id}/test-callback`,
  )
  return unwrap(res)
}
