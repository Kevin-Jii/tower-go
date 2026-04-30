import { http, unwrap } from './http'
import type { RouteImportedProductRow, ThirdPartyLogisticsSheet, ThirdPartyRoute } from './types'

export async function listThirdPartyRoutes(): Promise<ThirdPartyRoute[]> {
  const res = await http.get<import('./types').ApiEnvelope<ThirdPartyRoute[]>>('/third-party-routes')
  return unwrap(res)
}

export async function createThirdPartyRoute(body: Record<string, unknown>): Promise<ThirdPartyRoute> {
  const res = await http.post<import('./types').ApiEnvelope<ThirdPartyRoute>>('/third-party-routes', body)
  return unwrap(res)
}

export async function updateThirdPartyRoute(id: number, body: Record<string, unknown>): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/third-party-routes/${id}`, body)
}

export async function deleteThirdPartyRoute(id: number): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>(`/third-party-routes/${id}`)
}

export async function importThirdPartyRouteByDateRange(
  id: number,
  startDate: string,
  endDate: string,
): Promise<{ start_date: string; end_date: string; count: number; list: RouteImportedProductRow[] }> {
  const res = await http.post<import('./types').ApiEnvelope<{ start_date: string; end_date: string; count: number; list: RouteImportedProductRow[] }>>(
    `/third-party-routes/${id}/import-by-date`,
    { start_date: startDate, end_date: endDate },
  )
  return unwrap(res)
}

export async function saveThirdPartyLogisticsSheet(
  id: number,
  body: {
    start_date: string
    end_date: string
    headers: string[]
    rows: number[][]
    products: string[]
  },
): Promise<ThirdPartyLogisticsSheet> {
  const res = await http.post<import('./types').ApiEnvelope<ThirdPartyLogisticsSheet>>(`/third-party-routes/${id}/logistics-sheets`, body)
  return unwrap(res)
}

export async function listThirdPartyLogisticsSheets(id: number): Promise<ThirdPartyLogisticsSheet[]> {
  const res = await http.get<import('./types').ApiEnvelope<ThirdPartyLogisticsSheet[]>>(`/third-party-routes/${id}/logistics-sheets`)
  return unwrap(res)
}
