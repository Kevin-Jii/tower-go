import { http, unwrap } from './http'
import type { LoginPayload } from './types'

export async function login(phone: string, password: string): Promise<LoginPayload> {
  const res = await http.post<import('./types').ApiEnvelope<LoginPayload>>('/auth/login', {
    phone,
    password,
  })
  return unwrap(res)
}
