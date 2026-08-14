import {
  completeGalleryMultipartUpload,
  getGalleryMultipartUploadStatus,
  initGalleryMultipartUpload,
  uploadGalleryMultipartPart,
} from '@/api/gallery'
import type { GalleryMultipartUploadStatus } from '@/api/gallery'
import type { Gallery } from '@/api/types'

const FINGERPRINT_SAMPLE_SIZE = 64 * 1024
const UPLOAD_CONCURRENCY = 3
const PART_MAX_ATTEMPTS = 3

export interface GalleryUploadProgress {
  percent: number
  uploadedBytes: number
  totalBytes: number
  uploadedParts: number
  totalParts: number
  resumedParts: number
  phase: 'uploading' | 'merging'
}

export interface GalleryUploadOptions {
  category: string
  remark: string
  signal?: AbortSignal
  onProgress?: (progress: GalleryUploadProgress) => void
}

export async function uploadGalleryResumable(
  file: File,
  options: GalleryUploadOptions,
): Promise<Gallery> {
  const fingerprint = await createGalleryFileFingerprint(file, options.category)
  throwIfAborted(options.signal)

  let session = await initGalleryMultipartUpload({
    file_name: file.name,
    file_size: file.size,
    content_type: file.type,
    category: options.category,
    remark: options.remark,
    fingerprint,
  })
  session = await waitForCompletingSession(session, options)
  if (session.gallery) return session.gallery

  const uploadedPartNumbers = new Set(session.uploaded_parts.map((part) => part.part_number))
  const missingParts = Array.from({ length: session.total_parts }, (_, index) => index + 1)
    .filter((partNumber) => !uploadedPartNumbers.has(partNumber))
  const resumedParts = session.uploaded_parts.length
  let completedBytes = session.uploaded_bytes
  let completedParts = resumedParts
  let nextPartIndex = 0
  const inFlightBytes = new Map<number, number>()

  const reportProgress = (phase: GalleryUploadProgress['phase']): void => {
    const inFlight = Array.from(inFlightBytes.values()).reduce((sum, value) => sum + value, 0)
    const uploadedBytes = Math.min(file.size, completedBytes + inFlight)
    options.onProgress?.({
      percent: file.size > 0 ? Math.round((uploadedBytes / file.size) * 100) : 0,
      uploadedBytes,
      totalBytes: file.size,
      uploadedParts: completedParts,
      totalParts: session.total_parts,
      resumedParts,
      phase,
    })
  }
  reportProgress('uploading')

  const workerController = new AbortController()
  const abortWorkers = (): void => workerController.abort()
  options.signal?.addEventListener('abort', abortWorkers, { once: true })
  if (options.signal?.aborted) workerController.abort()

  const worker = async (): Promise<void> => {
    while (!workerController.signal.aborted) {
      const queueIndex = nextPartIndex
      nextPartIndex += 1
      if (queueIndex >= missingParts.length) return

      const partNumber = missingParts[queueIndex]
      const start = (partNumber - 1) * session.chunk_size
      const end = Math.min(file.size, start + session.chunk_size)
      const chunk = file.slice(start, end)
      await uploadPartWithRetry(session.session_id, partNumber, chunk, workerController.signal, (loaded) => {
        inFlightBytes.set(partNumber, Math.min(chunk.size, loaded))
        reportProgress('uploading')
      })
      inFlightBytes.delete(partNumber)
      completedBytes += chunk.size
      completedParts += 1
      reportProgress('uploading')
    }
    throw abortError()
  }

  try {
    const workerCount = Math.min(UPLOAD_CONCURRENCY, Math.max(1, missingParts.length))
    await Promise.all(Array.from({ length: workerCount }, () => worker()))
    throwIfAborted(options.signal)
    completedBytes = file.size
    inFlightBytes.clear()
    reportProgress('merging')
    return await completeGalleryMultipartUpload(session.session_id, workerController.signal)
  } catch (error) {
    workerController.abort()
    throw error
  } finally {
    options.signal?.removeEventListener('abort', abortWorkers)
  }
}

async function waitForCompletingSession(
  initial: GalleryMultipartUploadStatus,
  options: GalleryUploadOptions,
): Promise<GalleryMultipartUploadStatus> {
  let session = initial
  while (session.status === 'completing') {
    options.onProgress?.({
      percent: 100,
      uploadedBytes: session.file_size,
      totalBytes: session.file_size,
      uploadedParts: session.total_parts,
      totalParts: session.total_parts,
      resumedParts: session.total_parts,
      phase: 'merging',
    })
    await abortableDelay(2_000, options.signal ?? new AbortController().signal)
    throwIfAborted(options.signal)
    session = await getGalleryMultipartUploadStatus(session.session_id)
  }
  return session
}

export function isGalleryUploadAbort(error: unknown): boolean {
  if (error instanceof DOMException && error.name === 'AbortError') return true
  if (typeof error !== 'object' || error === null) return false
  const candidate = error as { code?: unknown; name?: unknown }
  return candidate.code === 'ERR_CANCELED' || candidate.name === 'CanceledError' || candidate.name === 'AbortError'
}

async function uploadPartWithRetry(
  sessionId: string,
  partNumber: number,
  chunk: Blob,
  signal: AbortSignal,
  onProgress: (loaded: number) => void,
): Promise<void> {
  let lastError: unknown
  for (let attempt = 1; attempt <= PART_MAX_ATTEMPTS; attempt += 1) {
    throwIfAborted(signal)
    onProgress(0)
    try {
      await uploadGalleryMultipartPart(sessionId, partNumber, chunk, {
        signal,
        onProgress: (event) => onProgress(event.loaded),
      })
      return
    } catch (error) {
      if (isGalleryUploadAbort(error)) throw error
      lastError = error
      if (attempt < PART_MAX_ATTEMPTS) {
        await abortableDelay(400 * 2 ** (attempt - 1), signal)
      }
    }
  }
  throw lastError
}

async function createGalleryFileFingerprint(file: File, category: string): Promise<string> {
  const first = new Uint8Array(await file.slice(0, FINGERPRINT_SAMPLE_SIZE).arrayBuffer())
  const lastStart = Math.max(0, file.size - FINGERPRINT_SAMPLE_SIZE)
  const last = new Uint8Array(await file.slice(lastStart).arrayBuffer())
  const metadata = new TextEncoder().encode(
    `${file.name}\0${file.size}\0${file.lastModified}\0${file.type}\0${category}\0`,
  )
  const fingerprintData = new Uint8Array(metadata.length + first.length + last.length)
  fingerprintData.set(metadata, 0)
  fingerprintData.set(first, metadata.length)
  fingerprintData.set(last, metadata.length + first.length)
  const digest = new Uint8Array(await crypto.subtle.digest('SHA-256', fingerprintData))
  return Array.from(digest, (byte) => byte.toString(16).padStart(2, '0')).join('')
}

function throwIfAborted(signal?: AbortSignal): void {
  if (signal?.aborted) throw abortError()
}

function abortError(): DOMException {
  return new DOMException('Upload aborted', 'AbortError')
}

function abortableDelay(milliseconds: number, signal: AbortSignal): Promise<void> {
  return new Promise((resolve, reject) => {
    if (signal.aborted) {
      reject(abortError())
      return
    }
    const timer = window.setTimeout(() => {
      signal.removeEventListener('abort', onAbort)
      resolve()
    }, milliseconds)
    const onAbort = (): void => {
      window.clearTimeout(timer)
      reject(abortError())
    }
    signal.addEventListener('abort', onAbort, { once: true })
  })
}
