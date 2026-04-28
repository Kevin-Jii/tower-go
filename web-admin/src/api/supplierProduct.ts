import { http, unwrap } from './http'
import type { Paginated, ProductUnitSpec, StorePurchasableProduct, SupplierCategory } from './types'

export async function listSupplierProducts(params?: {
  supplier_id?: number
  category_id?: number
  keyword?: string
  status?: number
  page?: number
  page_size?: number
}): Promise<Paginated<StorePurchasableProduct>> {
  const res = await http.get<import('./types').ApiEnvelope<Paginated<StorePurchasableProduct>>>('/supplier-products', { params })
  return unwrap(res)
}

export async function createSupplierProduct(body: {
  supplier_id: number
  category_id: number
  name: string
  unit: string
  bottle_price: number
  case_price: number
  bottles_per_case: number
  spec?: string
  remark?: string
}): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/supplier-products', body)
}

export async function updateSupplierProduct(id: number, body: Record<string, unknown>): Promise<void> {
  await http.put<import('./types').ApiEnvelope<unknown>>(`/supplier-products/${id}`, body)
}

export async function getSupplierProduct(id: number): Promise<StorePurchasableProduct> {
  const res = await http.get<import('./types').ApiEnvelope<StorePurchasableProduct>>(`/supplier-products/${id}`)
  return unwrap(res)
}

export async function listSupplierCategories(supplierId: number): Promise<SupplierCategory[]> {
  const res = await http.get<import('./types').ApiEnvelope<SupplierCategory[]>>('/supplier-categories', {
    params: { supplier_id: supplierId },
  })
  return unwrap(res)
}

export async function createSupplierCategory(body: { supplier_id: number; name: string; sort?: number }): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/supplier-categories', body)
}

export async function batchUpsertProductUnitSpecs(body: {
  product_id: number
  units: Array<{
    unit_code: string
    unit_name?: string
    factor_to_base: number
    precision: number
    cost_price: number
    sale_price: number
    is_enabled?: boolean
  }>
}): Promise<void> {
  await http.post<import('./types').ApiEnvelope<unknown>>('/product-unit-specs/batch', body)
}

export async function listProductUnitSpecs(productId: number): Promise<ProductUnitSpec[]> {
  const res = await http.get<import('./types').ApiEnvelope<ProductUnitSpec[]>>('/product-unit-specs', {
    params: { product_id: productId },
  })
  return unwrap(res)
}
