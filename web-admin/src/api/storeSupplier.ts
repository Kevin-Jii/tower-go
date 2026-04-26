import { http, unwrap } from './http'
import type { StorePurchasableProduct, StoreSupplierBinding } from './types'

/** 当前门店（或管理员指定 store_id）已绑定的供应商 */
export async function listStoreBoundSuppliers(params?: { store_id?: number }): Promise<StoreSupplierBinding[]> {
  const res = await http.get<import('./types').ApiEnvelope<StoreSupplierBinding[]>>('/store-suppliers', { params })
  return unwrap(res)
}

/** 绑定供应商下、门店可采购的商品（含 supplier / category 等关联） */
export async function listPurchasableProducts(params?: {
  keyword?: string
  supplier_id?: number
  category_id?: number
  store_id?: number
}): Promise<StorePurchasableProduct[]> {
  const res = await http.get<import('./types').ApiEnvelope<StorePurchasableProduct[]>>('/store-suppliers/products', {
    params,
  })
  return unwrap(res)
}

export async function bindStoreSuppliers(body: { store_id: number; supplier_ids: number[] }): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/store-suppliers', body)
}

export async function unbindStoreSuppliers(body: { store_id: number; supplier_ids: number[] }): Promise<void> {
  await http.delete<import('./types').ApiEnvelope<unknown>>('/store-suppliers', { data: body })
}
