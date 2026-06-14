import { http, unwrap } from './http'
import type { ApiEnvelope, Paginated } from './types'

export interface PrinterRow {
  id: number
  store_id: number
  sn: string
  name: string
  type: number
  status: number
  is_default: number
  online: number
  last_heartbeat?: string
  remark?: string
  created_at?: string
  updated_at?: string
}

export interface PrinterPayload {
  store_id: number
  sn: string
  name?: string
  type?: number
  is_default?: number
  remark?: string
}

export interface UpdatePrinterPayload {
  name?: string
  type?: number
  status?: number
  is_default?: number
  remark?: string
}

export interface PrinterStatusRow {
  sn: string
  online: number
  status?: number
}

export async function listStorePrinters(store_id: number): Promise<PrinterRow[]> {
  const res = await http.get<ApiEnvelope<PrinterRow[]>>('/printers', {
    params: { store_id },
  })
  return unwrap(res)
}

export async function listAllPrinters(): Promise<Paginated<PrinterRow>> {
  const res = await http.get<ApiEnvelope<Paginated<PrinterRow>>>('/printers/all')
  return unwrap(res)
}

export async function bindPrinter(body: PrinterPayload): Promise<void> {
  await http.post<ApiEnvelope<unknown>>('/printers/bind', body)
}

export async function updatePrinter(id: number, body: UpdatePrinterPayload): Promise<void> {
  await http.put<ApiEnvelope<unknown>>(`/printers/${id}`, body)
}

export async function unbindPrinter(id: number): Promise<void> {
  await http.delete<ApiEnvelope<unknown>>(`/printers/${id}`)
}

export async function batchQueryPrinterStatus(store_id: number): Promise<PrinterStatusRow[]> {
  const res = await http.get<ApiEnvelope<PrinterStatusRow[]>>('/printers/status/batch', {
    params: { store_id },
  })
  return unwrap(res)
}

export async function testPrint(
  printerId: number,
  body: { content?: string; copies?: number },
): Promise<{ order_id: string }> {
  const res = await http.post<ApiEnvelope<{ order_id: string }>>(`/printers/${printerId}/test`, body)
  return unwrap(res)
}

export async function printPurchaseOrder(printerId: number, orderId: number): Promise<{ order_id: string }> {
  const res = await http.post<ApiEnvelope<{ order_id: string }>>(`/printers/${printerId}/print/purchase-order`, {
    order_id: orderId,
  })
  return unwrap(res)
}
