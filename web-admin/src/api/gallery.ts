import { http, unwrap } from './http'
import type { Gallery } from './types'

export async function listGalleries(params?: {
  page?: number
  page_size?: number
  category?: string
  keyword?: string
  store_id?: number
}): Promise<Gallery[]> {
  const res = await http.get<import('./types').ApiEnvelope<Gallery[]>>('/galleries', { params })
  return unwrap(res)
}

export async function getGallery(id: number): Promise<Gallery> {
  const res = await http.get<import('./types').ApiEnvelope<Gallery>>(`/galleries/${id}`)
  return unwrap(res)
}

export async function uploadGallery(body: FormData): Promise<Gallery> {
  const res = await http.post<import('./types').ApiEnvelope<Gallery>>('/galleries/upload', body, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return unwrap(res)
}

export async function updateGallery(id: number, body: { name?: string; category?: string; remark?: string }): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/galleries/${id}`, body)
}

export async function deleteGallery(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/galleries/${id}`)
}

export async function batchDeleteGallery(ids: number[]): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/galleries/batch-delete', { ids })
}
