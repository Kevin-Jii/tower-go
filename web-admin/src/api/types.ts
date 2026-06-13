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
  administrative_unit?: string
  phone?: string
  business_hours?: string
  contact_person?: string
  remark?: string
  third_party_account_id?: number | null
  third_party_account?: ThirdPartyAccount
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

export interface Gallery {
  id: number
  name: string
  path: string
  url: string
  size: number
  mime_type?: string
  category?: string
  store_id?: number
  upload_by?: number
  upload_by_name?: string
  remark?: string
  created_at?: string
  updated_at?: string
}

export interface MessageTemplate {
  id: number
  code: string
  name: string
  title?: string
  content: string
  description?: string
  variables?: string
  is_enabled: boolean
  created_at?: string
  updated_at?: string
}

export interface ThirdPartyAccount {
  id: number
  platform_name: string
  name: string
  login_name: string
  phone?: string
  password: string
  application_key: string
  login_type?: string
  channel?: string
  shop_id?: string
  customer_id?: string
  is_enabled: boolean
  last_test_ok?: boolean
  last_test_msg?: string
  last_token?: string
  token_valid_time?: number
  last_test_at?: string
  last_sync_at?: string
  last_sync_msg?: string
  last_sync_count?: number
  remark?: string
  created_at?: string
  updated_at?: string
}

export interface ThirdPartyOrder {
  id: number
  account_id: number
  platform_name: string
  order_no: string
  place_time?: string
  place_date?: string
  order_trade_status?: string
  status_name?: string
  pay_amount?: number
  total_amount?: number
  total_item_num?: number
  raw_json?: string
  synced_at?: string
  created_at?: string
}

export interface ThirdPartyRouteStore {
  id: number
  route_id: number
  store_id: number
  sort: number
  store?: Store
}

export interface ThirdPartyRoute {
  id: number
  name: string
  remark?: string
  stores?: ThirdPartyRouteStore[]
  created_at?: string
  updated_at?: string
}

export interface RouteStoreQuantity {
  store_id: number
  store_name: string
  quantity: number
}

export interface RouteImportedProductRow {
  product_name: string
  total_qty: number
  store_qty: RouteStoreQuantity[]
}

export interface ThirdPartyLogisticsSheet {
  id: number
  route_id: number
  sheet_date: string
  start_date: string
  end_date: string
  headers: string[]
  rows: number[][]
  products: string[]
  created_at?: string
  updated_at?: string
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

/** 门店-供应商绑定记录 */
export interface StoreSupplierBinding {
  id: number
  store_id: number
  supplier_id: number
  status: number
  supplier?: Supplier
}

export interface SupplierCategory {
  id: number
  supplier_id: number
  name: string
  sort?: number
  status?: number
}

/** 门店可采购商品（/store-suppliers/products，结构与 SupplierProduct 对齐） */
export interface StorePurchasableProduct {
  id: number
  supplier_id: number
  category_id?: number
  name: string
  unit: string
  price?: number
  bottle_price?: number
  case_price?: number
  bottles_per_case?: number
  spec?: string
  status?: number
  supplier?: { supplier_name?: string; supplier_code?: string }
  category?: { name?: string }
}

export interface ProductUnitSpec {
  id: number
  product_id: number
  unit_code: string
  unit_name: string
  factor_to_base: number
  precision: number
  cost_price: number
  sale_price: number
  is_enabled: boolean
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

export type InventoryLossType = 'loss' | 'self_use' | 'gift'

export interface InventoryLossOrderItem {
  id: number
  order_id?: number
  product_id: number
  product_name?: string
  unit: string
  quantity: number
  base_quantity?: number
  base_unit?: string
  cost_price: number
  cost_amount: number
  remark?: string
  created_at?: string
}

export interface InventoryLossOrder {
  id: number
  order_no: string
  store_id: number
  type: InventoryLossType
  member_id?: number
  member?: {
    id: number
    name?: string
    phone?: string
  }
  reason?: string
  total_cost: number
  item_count: number
  operator_id?: number
  operator_name?: string
  created_at?: string
  updated_at?: string
  items?: InventoryLossOrderItem[]
}

export interface InventoryLossOrderDetail extends InventoryLossOrder {
  items: InventoryLossOrderItem[]
}

export interface MemberGiftRecord {
  id: number
  order_id?: number
  order_no: string
  product_id: number
  product_name?: string
  unit: string
  quantity: number
  cost_amount: number
  reason?: string
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

export interface StoreAccountConsumable {
  id: number
  account_id?: number
  product_id: number
  product_name?: string
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
  member_id?: number
  payment_status?: number
  member?: {
    id: number
    name?: string
    phone?: string
  }
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
  consumables?: StoreAccountConsumable[]
}

export interface MemberRow {
  id: number
  store_id?: number
  uid: string
  name?: string
  phone: string
  balance: string | number
  points: number
  level: number
  version: number
  createTime?: string
}

export interface MemberConsumptionRecord {
  account_id: number
  account_no: string
  account_date: string
  channel: string
  channel_name?: string
  order_no?: string
  total_amount: number
  other_expense_amount: number
  consumable_amount: number
  net_income_amount: number
  created_at?: string
}

export interface MemberConsumptionSummary {
  count: number
  total_amount: number
  other_expense_amount: number
  consumable_amount: number
  net_income_amount: number
}

export interface MemberConsumptionPage {
  list: MemberConsumptionRecord[]
  total: number
  page: number
  page_size: number
  summary?: MemberConsumptionSummary
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
  net_income_amount?: number
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

export interface RadarMetricItem {
  name: string
  value: number
}

export interface HomeChartsStats {
  start_date?: string
  end_date?: string
  line?: SalesTrendItem[]
  pie?: ChannelStatsItem[]
  radar?: RadarMetricItem[]
  overview?: BusinessOverviewStats
}

/** 钉钉机器人（列表/详情，与后端 model.DingTalkBot JSON 对齐） */
export interface DingTalkBot {
  id: number
  name: string
  bot_type: string
  webhook?: string
  secret?: string
  client_id?: string
  client_secret?: string
  agent_id?: string
  robot_code?: string
  store_id?: number | null
  store_code?: string
  store_name?: string
  is_enabled: boolean
  msg_type?: string
  card_msg_key?: string
  remark?: string
  created_at?: string
  updated_at?: string
}

export interface CreateDingTalkBotReq {
  name?: string
  bot_type?: string
  webhook?: string
  secret?: string
  client_id?: string
  client_secret?: string
  agent_id?: string
  robot_code?: string
  store_id?: number | null
  is_enabled?: boolean
  msg_type?: string
  card_msg_key?: string
  remark?: string
}

export type UpdateDingTalkBotReq = Partial<{
  name: string
  bot_type: string
  webhook: string | null
  secret: string | null
  client_id: string | null
  client_secret: string | null
  agent_id: string | null
  robot_code: string | null
  store_id: number | null
  is_enabled: boolean
  msg_type: string
  card_msg_key: string | null
  remark: string | null
}>
