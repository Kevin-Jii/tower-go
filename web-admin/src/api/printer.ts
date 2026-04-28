import { http, unwrap } from './http'

export interface PrinterRow {
  id: number
  store_id: number
  sn: string
  name: string
  type: number
  status: number
  is_default: number
  online: number
  remark?: string
}

export async function listStorePrinters(store_id: number): Promise<PrinterRow[]> {
  const res = await http.get<import('./types').ApiEnvelope<PrinterRow[]>>('/printers', {
    params: { store_id },
  })
  return unwrap(res)
}

export async function printPurchaseOrder(printerId: number, orderId: number): Promise<{ order_id: string }> {
  const res = await http.post<import('./types').ApiEnvelope<{ order_id: string }>>(
    `/printers/${printerId}/print/purchase-order`,
    { order_id: orderId },
  )
  return unwrap(res)
}
