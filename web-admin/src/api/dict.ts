import { http, unwrap } from './http'
import type { DictData, DictType } from './types'

export async function listDictTypes(): Promise<DictType[]> {
  const res = await http.get<import('./types').ApiEnvelope<DictType[]>>('/dict-types')
  return unwrap(res)
}

export async function createDictType(body: Record<string, unknown>): Promise<DictType> {
  const res = await http.post<import('./types').ApiEnvelope<DictType>>('/dict-types', body)
  return unwrap(res)
}

export async function updateDictType(id: number, body: Record<string, unknown>): Promise<DictType> {
  const res = await http.put<import('./types').ApiEnvelope<DictType>>(`/dict-types/${id}`, body)
  return unwrap(res)
}

export async function deleteDictType(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/dict-types/${id}`)
}

export async function listDictDataByTypeCode(typeCode: string): Promise<DictData[]> {
  const res = await http.get<import('./types').ApiEnvelope<DictData[]>>('/dict-data', {
    params: { type_code: typeCode },
  })
  return unwrap(res)
}

export async function createDictData(body: Record<string, unknown>): Promise<DictData> {
  const res = await http.post<import('./types').ApiEnvelope<DictData>>('/dict-data', body)
  return unwrap(res)
}

export async function updateDictData(id: number, body: Record<string, unknown>): Promise<DictData> {
  const res = await http.put<import('./types').ApiEnvelope<DictData>>(`/dict-data/${id}`, body)
  return unwrap(res)
}

export async function deleteDictData(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/dict-data/${id}`)
}

export async function getAllDicts(): Promise<Record<string, DictData[]>> {
  const res = await http.get<import('./types').ApiEnvelope<Record<string, DictData[]>>>('/dicts')
  return unwrap(res)
}
