package model

// DashboardStats 统计面板数据
type DashboardStats struct {
	Inventory InventoryStats `json:"inventory"` // 库存统计
	Sales     SalesStats     `json:"sales"`     // 销售统计
}

// InventoryStats 库存统计
type InventoryStats struct {
	TotalProducts int64   `json:"total_products"` // 商品种类数
	TotalQuantity float64 `json:"total_quantity"` // 总库存数量
	TotalRecords  int64   `json:"total_records"`  // 出入库记录数
	TodayIn       float64 `json:"today_in"`       // 今日入库
	TodayOut      float64 `json:"today_out"`      // 今日出库
}

// SalesStats 销售统计
type SalesStats struct {
	TotalAmount  float64 `json:"total_amount"`  // 总销售额
	TotalOrders  int64   `json:"total_orders"`  // 总订单数
	TotalQty     float64 `json:"total_qty"`     // 总销售数量
	AvgAmount    float64 `json:"avg_amount"`    // 平均客单价
	PeriodLabel  string  `json:"period_label"`  // 周期标签
}

// SalesTrendItem 销售趋势项
type SalesTrendItem struct {
	Date   string  `json:"date"`   // 日期
	Amount float64 `json:"amount"` // 销售额
	Orders int64   `json:"orders"` // 订单数
}

// ChannelStatsItem 渠道统计项
type ChannelStatsItem struct {
	Channel     string  `json:"channel"`      // 渠道编码
	ChannelName string  `json:"channel_name"` // 渠道名称
	Amount      float64 `json:"amount"`       // 销售额
	Orders      int64   `json:"orders"`       // 订单数
	Percent     float64 `json:"percent"`      // 占比
}

// TopProductItem 热销商品项
type TopProductItem struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Amount      float64 `json:"amount"`
}
