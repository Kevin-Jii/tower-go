import axios, { type AxiosResponse } from 'axios'
import type { ApiEnvelope } from './types'
import { getStoredTenantId, getToken } from '@/utils/storage'

const baseURL = import.meta.env.VITE_API_BASE || '/api/v1'

export const http = axios.create({
  baseURL,
  timeout: 30_000,
})

http.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  const tid = getStoredTenantId()
  if (tid != null && tid > 0) {
    config.headers['x-tenant-id'] = String(tid)
  }
  return config
})

http.interceptors.response.use(
  (res: AxiosResponse<ApiEnvelope>) => {
    const body = res.data
    if (typeof body?.code === 'number' && body.code !== 200) {
      return Promise.reject(new Error(body.message || '请求失败'))
    }
    return res
  },
  (err) => Promise.reject(err),
)

export function unwrap<T>(res: AxiosResponse<ApiEnvelope<T>>): T {
  return res.data.data as T
}
