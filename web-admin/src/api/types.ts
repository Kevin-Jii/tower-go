export interface ApiEnvelope<T = unknown> {
  code: number
  message: string
  data?: T
  error?: string
}

export interface Paginated<T> {
  list: T[]
  total: number
  page: number
  page_size: number
  page_num: number
}

export interface LoginPayload {
  token: string
  token_type: string
  expires_in: number
  user_info: User
  strategy?: string
}

export interface User {
  id: number
  username: string
  phone: string
  nickname?: string
  email?: string
  role_id?: number
  store_id?: number
  status?: number
  gender?: number
  role?: Role
  store?: Store
}

export interface Role {
  id: number
  name: string
  code: string
  status: number
  description?: string
  data_scope?: number
}

export interface Store {
  id: number
  name: string
  /** 门店编码（与后端 model.Store.store_code 一致） */
  store_code?: string | null
  address?: string
  phone?: string
  business_hours?: string
  contact_person?: string
  remark?: string
  status?: number
  created_at?: string
  updated_at?: string
}

export interface Menu {
  id: number
  parent_id: number
  name: string
  title: string
  icon?: string
  path?: string
  component?: string
  type: number
  sort: number
  permission?: string
  visible: number
  status: number
  remark?: string
  children?: Menu[]
}

export interface DictType {
  id: number
  code: string
  name: string
  remark?: string
  status: number
}

export interface AssignMenusToRoleReq {
  role_id: number
  /** 允许空数组表示清空该角色菜单 */
  menu_ids: number[]
  perms?: Record<number, number>
}

export interface DictData {
  id: number
  type_id: number
  type_code: string
  label: string
  value: string
  sort: number
  css_class?: string
  list_class?: string
  is_default?: boolean
  remark?: string
  status: number
}

/** 供应商 */
export interface Supplier {
  id: number
  supplier_code: string
  supplier_name: string
  contact_person?: string
  contact_phone?: string
  contact_email?: string
  supplier_address?: string
  remark?: string
  status: number
  created_at?: string
  updated_at?: string
}

/** 门店可采购商品（/store-suppliers/products） */
export interface StorePurchasableProduct {
  id: number
  supplier_id: number
  name: string
  unit: string
  supplier?: { supplier_name?: string }
}

export interface PurchaseOrderItem {
  id: number
  order_id?: number
  supplier_id: number
  product_id: number
  quantity: number
  unit_price: number
  amount: number
  remark?: string
  supplier?: { supplier_name?: string }
  product?: { name?: string; unit?: string }
}

export interface PurchaseOrder {
  id: number
  order_no: string
  store_id: number
  total_amount: number
  status: number
  remark?: string
  order_date: string
  created_by?: number
  creator?: { nickname?: string; username?: string }
  store?: Store
  items?: PurchaseOrderItem[]
  created_at?: string
  updated_at?: string
}

export interface InventoryRow {
  id: number
  store_id: number
  store_name?: string
  product_id: number
  product_name: string
  quantity: number
  unit: string
  price?: number
}

export interface InventoryOrder {
  id: number
  order_no: string
  type: number
  store_id?: number
  store_name?: string
  reason?: string
  remark?: string
  total_quantity?: number
  item_count?: number
  operator_name?: string
  created_at?: string
}

export interface StoreAccountItem {
  id: number
  account_id?: number
  product_id: number
  product_name?: string
  spec?: string
  quantity: number
  unit?: string
  price: number
  amount: number
  remark?: string
}

export interface StoreAccount {
  id: number
  account_no: string
  store_id: number
  channel: string
  order_no?: string
  total_amount: number
  other_expense_amount?: number
  net_income_amount?: number
  item_count?: number
  tag_code?: string
  tag_name?: string
  remark?: string
  account_date: string
  created_at?: string
  items?: StoreAccountItem[]
}

export interface MemberRow {
  id: number
  uid: string
  name?: string
  phone: string
  balance: string | number
  points: number
  level: number
  version: number
  createTime?: string
}

export interface InventoryStats {
  total_products: number
  total_quantity: number
  total_records: number
  today_in: number
  today_out: number
}

export interface SalesStats {
  total_amount: number
  today_amount: number
  month_amount: number
  total_orders: number
  total_qty: number
  avg_amount: number
  period_label?: string
}

export interface DashboardStats {
  inventory: InventoryStats
  sales: SalesStats
}

export interface StoreAccountStats {
  total_amount?: number
  count?: number
}

export interface CategoryAmountItem {
  category_id: number
  category_name: string
  in_amount: number
  out_amount: number
  net_amount: number
}

export interface BusinessOverviewStats {
  start_date?: string
  end_date?: string
  store_id?: number
  inbound_amount?: number
  outbound_amount?: number
  sales_amount?: number
  other_expense_amount?: number
  gross_profit_amount?: number
  net_profit_amount?: number
  sales_order_count?: number
  inventory_in_count?: number
  inventory_out_count?: number
  categories?: CategoryAmountItem[]
}

export interface SalesTrendItem {
  date: string
  amount: number
  orders: number
}

export interface ChannelStatsItem {
  channel: string
  channel_name: string
  amount: number
  orders: number
  percent: number
}

export interface HomeChartsStats {
  start_date?: string
  end_date?: string
  line?: SalesTrendItem[]
  pie?: ChannelStatsItem[]
  overview?: BusinessOverviewStats
}
