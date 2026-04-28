import { http, unwrap } from './http'
import type { MessageTemplate } from './types'

export async function listMessageTemplates(): Promise<MessageTemplate[]> {
  const res = await http.get<import('./types').ApiEnvelope<MessageTemplate[]>>('/message-templates')
  return unwrap(res)
}

export async function getMessageTemplate(id: number): Promise<MessageTemplate> {
  const res = await http.get<import('./types').ApiEnvelope<MessageTemplate>>(`/message-templates/${id}`)
  return unwrap(res)
}

export async function createMessageTemplate(body: {
  code: string
  name: string
  title?: string
  content: string
  description?: string
  variables?: string
  is_enabled?: boolean
}): Promise<MessageTemplate> {
  const res = await http.post<import('./types').ApiEnvelope<MessageTemplate>>('/message-templates', body)
  return unwrap(res)
}

export async function updateMessageTemplate(
  id: number,
  body: {
    name?: string
    title?: string
    content?: string
    description?: string
    variables?: string
    is_enabled?: boolean
  },
): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/message-templates/${id}`, body)
}

export async function deleteMessageTemplate(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/message-templates/${id}`)
}
