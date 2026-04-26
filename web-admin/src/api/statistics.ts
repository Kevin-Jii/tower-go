import { http, unwrap } from './http'
import type {
  BusinessOverviewStats,
  ChannelStatsItem,
  DashboardStats,
  HomeChartsStats,
  InventoryStats,
  SalesStats,
  SalesTrendItem,
} from './types'

export async function getStatisticsDashboard(params?: {
  period?: string
  store_id?: number
}): Promise<DashboardStats> {
  const res = await http.get<import('./types').ApiEnvelope<DashboardStats>>('/statistics/dashboard', { params })
  return unwrap(res)
}

export async function getStatisticsInventory(params?: { store_id?: number }): Promise<InventoryStats> {
  const res = await http.get<import('./types').ApiEnvelope<InventoryStats>>('/statistics/inventory', { params })
  return unwrap(res)
}

export async function getStatisticsSales(params?: { period?: string; store_id?: number }): Promise<SalesStats> {
  const res = await http.get<import('./types').ApiEnvelope<SalesStats>>('/statistics/sales', { params })
  return unwrap(res)
}

export async function getStatisticsSalesTrend(params?: {
  period?: string
  store_id?: number
}): Promise<SalesTrendItem[]> {
  const res = await http.get<import('./types').ApiEnvelope<SalesTrendItem[]>>('/statistics/sales-trend', { params })
  return unwrap(res)
}

export async function getStatisticsChannel(params?: {
  period?: string
  store_id?: number
}): Promise<ChannelStatsItem[]> {
  const res = await http.get<import('./types').ApiEnvelope<ChannelStatsItem[]>>('/statistics/channel', { params })
  return unwrap(res)
}

export async function getBusinessOverview(params: {
  start_date: string
  end_date: string
  store_id?: number
}): Promise<BusinessOverviewStats> {
  const res = await http.get<import('./types').ApiEnvelope<BusinessOverviewStats>>('/statistics/business-overview', {
    params,
  })
  return unwrap(res)
}

export async function getHomeCharts(params: {
  start_date: string
  end_date: string
  granularity?: string
  store_id?: number
}): Promise<HomeChartsStats> {
  const res = await http.get<import('./types').ApiEnvelope<HomeChartsStats>>('/statistics/home-charts', { params })
  return unwrap(res)
}
