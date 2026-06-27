import axios, { type AxiosError, type AxiosResponse } from 'axios'
import type { ApiEnvelope } from './types'
import { getStoredTenantId, getToken } from '@/utils/storage'

const baseURL = import.meta.env.VITE_API_BASE || '/api/v1'

/** 与后端 pkg/clientsource 约定；管理后台默认 web-admin，与 Taro 端 weapp/web/app 等区分 */
const clientSource =
  (import.meta.env.VITE_CLIENT_SOURCE && String(import.meta.env.VITE_CLIENT_SOURCE).trim()) ||
  'web-admin'

export const http = axios.create({
  baseURL,
  timeout: 30_000,
})

let authRedirecting = false

function isAuthExpiredCode(code: unknown): boolean {
  return code === 40101 || code === 40102 || code === 40103
}

async function redirectToLogin(): Promise<void> {
  if (authRedirecting) return
  authRedirecting = true
  const [{ useUserStore }, { default: router }] = await Promise.all([
    import('@/store/user'),
    import('@/router'),
  ])
  const userStore = useUserStore()
  await userStore.logout({
    redirect: router.currentRoute.value.fullPath,
    message: '登录已过期，请重新登录',
  })
  authRedirecting = false
}

http.interceptors.request.use((config) => {
  config.headers['X-Client-Source'] = clientSource
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
      if (isAuthExpiredCode(body.code)) {
        void redirectToLogin()
      }
      return Promise.reject(new Error(body.message || '请求失败'))
    }
    return res
  },
  (err: AxiosError<ApiEnvelope>) => {
    if (err.response?.status === 401 || isAuthExpiredCode(err.response?.data?.code)) {
      void redirectToLogin()
    }
    return Promise.reject(err)
  },
)

export function unwrap<T>(res: AxiosResponse<ApiEnvelope<T>>): T {
  return res.data.data as T
}
