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
	TotalAmount float64 `json:"total_amount"` // 总销售额
	TodayAmount float64 `json:"today_amount"` // 今日销售额
	MonthAmount float64 `json:"month_amount"` // 本月销售额
	TotalOrders int64   `json:"total_orders"` // 总订单数
	TotalQty    float64 `json:"total_qty"`    // 总销售数量
	AvgAmount   float64 `json:"avg_amount"`   // 平均客单价
	PeriodLabel string  `json:"period_label"` // 周期标签
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

// BusinessOverviewReq 经营统计查询
type BusinessOverviewReq struct {
	StartDate string `form:"start_date" binding:"required"`
	EndDate   string `form:"end_date" binding:"required"`
	StoreID   uint   `form:"store_id"`
}

// CategoryAmountItem 品类金额统计
type CategoryAmountItem struct {
	CategoryID   uint    `json:"category_id"`
	CategoryName string  `json:"category_name"`
	InAmount     float64 `json:"in_amount"`
	OutAmount    float64 `json:"out_amount"`
	NetAmount    float64 `json:"net_amount"`
}

// MemberConsumptionRankItem 会员消费排行
type MemberConsumptionRankItem struct {
	MemberID    uint    `json:"member_id"`
	MemberName  string  `json:"member_name"`
	MemberPhone string  `json:"member_phone"`
	Amount      float64 `json:"amount"`
	Orders      int64   `json:"orders"`
}

// BusinessOverviewStats 经营汇总统计
type BusinessOverviewStats struct {
	StartDate              string                      `json:"start_date"`
	EndDate                string                      `json:"end_date"`
	StoreID                uint                        `json:"store_id"`
	InboundAmount          float64                     `json:"inbound_amount"`
	OutboundAmount         float64                     `json:"outbound_amount"`
	AllCategoryAmount      float64                     `json:"all_category_amount"`
	SalesAmount            float64                     `json:"sales_amount"`
	ConsumableAmount       float64                     `json:"consumable_amount"`
	B2BSupplyAmount        float64                     `json:"b2b_supply_amount"`
	B2BSupplyOrderCount    int64                       `json:"b2b_supply_order_count"`
	ReturnDepositAmount    float64                     `json:"return_deposit_amount"`
	ReturnLogisticsFee     float64                     `json:"return_logistics_fee"`
	ErrandFeeAmount        float64                     `json:"errand_fee_amount"`
	InventoryLossAmount    float64                     `json:"inventory_loss_amount"`
	InventoryLossCount     int64                       `json:"inventory_loss_count"`
	InventorySelfUseAmount float64                     `json:"inventory_self_use_amount"`
	InventorySelfUseCount  int64                       `json:"inventory_self_use_count"`
	OtherExpenseAmount     float64                     `json:"other_expense_amount"`
	GrossProfitAmount      float64                     `json:"gross_profit_amount"`
	NetProfitAmount        float64                     `json:"net_profit_amount"`
	SalesOrderCount        int64                       `json:"sales_order_count"`
	InventoryInCount       int64                       `json:"inventory_in_count"`
	InventoryOutCount      int64                       `json:"inventory_out_count"`
	Categories             []CategoryAmountItem        `json:"categories"`
	MemberConsumptionRank  []MemberConsumptionRankItem `json:"member_consumption_rank"`
}

// RadarMetricItem 雷达图指标
type RadarMetricItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

// HomeChartsStats 首页图表统计
type HomeChartsStats struct {
	StartDate string                `json:"start_date"`
	EndDate   string                `json:"end_date"`
	Line      []SalesTrendItem      `json:"line"`     // 折线图：销售趋势
	Pie       []ChannelStatsItem    `json:"pie"`      // 扇形图：渠道占比
	Radar     []RadarMetricItem     `json:"radar"`    // 雷达图：经营指标
	Overview  BusinessOverviewStats `json:"overview"` // 汇总卡片
}
