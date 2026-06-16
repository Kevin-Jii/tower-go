import { http, unwrap } from './http'
import type { AuditLog, Paginated } from './types'

export interface AuditLogListParams {
  page?: number
  page_size?: number
  start_time?: string
  end_time?: string
  user_id?: number
  store_id?: number
  module?: string
  action?: string
  status?: string
  keyword?: string
}

export async function listAuditLogs(params: AuditLogListParams): Promise<Paginated<AuditLog>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<AuditLog>>>('/audit-logs', {
    params,
  })
  return unwrap(res)
}

export async function getAuditLog(id: number): Promise<AuditLog> {
  const res = await http.get<import('./types').ApiEnvelope<AuditLog>>(`/audit-logs/${id}`)
  return unwrap(res)
}

