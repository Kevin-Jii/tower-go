import { http, unwrap } from './http'
import type { Gallery } from './types'
import type { AxiosProgressEvent } from 'axios'

export interface GalleryMultipartUploadInit {
  file_name: string
  file_size: number
  content_type?: string
  category?: string
  remark?: string
  fingerprint: string
}

export interface GalleryUploadedPart {
  part_number: number
  size: number
}

export interface GalleryMultipartUploadStatus {
  session_id: string
  chunk_size: number
  total_parts: number
  uploaded_parts: GalleryUploadedPart[]
  uploaded_bytes: number
  file_size: number
  status: string
  expires_at: string
  gallery?: Gallery
}

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
  const res = await http.post<import('./types').ApiEnvelope<Gallery>>('/galleries/upload', body)
  return unwrap(res)
}

export async function initGalleryMultipartUpload(body: GalleryMultipartUploadInit): Promise<GalleryMultipartUploadStatus> {
  const res = await http.post<import('./types').ApiEnvelope<GalleryMultipartUploadStatus>>('/galleries/multipart/init', body)
  return unwrap(res)
}

export async function getGalleryMultipartUploadStatus(sessionId: string): Promise<GalleryMultipartUploadStatus> {
  const res = await http.get<import('./types').ApiEnvelope<GalleryMultipartUploadStatus>>(`/galleries/multipart/${sessionId}`)
  return unwrap(res)
}

export async function uploadGalleryMultipartPart(
  sessionId: string,
  partNumber: number,
  body: Blob,
  options?: { signal?: AbortSignal; onProgress?: (event: AxiosProgressEvent) => void },
): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(
    `/galleries/multipart/${sessionId}/parts/${partNumber}`,
    body,
    {
      headers: { 'Content-Type': 'application/octet-stream' },
      timeout: 120_000,
      signal: options?.signal,
      onUploadProgress: options?.onProgress,
    },
  )
}

export async function completeGalleryMultipartUpload(sessionId: string, signal?: AbortSignal): Promise<Gallery> {
  const res = await http.post<import('./types').ApiEnvelope<Gallery>>(
    `/galleries/multipart/${sessionId}/complete`,
    undefined,
    { signal, timeout: 120_000 },
  )
  return unwrap(res)
}

export async function abortGalleryMultipartUpload(sessionId: string): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/galleries/multipart/${sessionId}`)
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
