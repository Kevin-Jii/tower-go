package module

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/businessdate"
	"gorm.io/gorm"
)

type categoryAmountRow struct {
	CategoryID   uint
	CategoryName string
	InAmount     float64
	OutAmount    float64
}

type StatisticsModule struct {
	db *gorm.DB
}

func NewStatisticsModule(db *gorm.DB) *StatisticsModule {
	return &StatisticsModule{db: db}
}

func withStoreID(query *gorm.DB, storeID uint) *gorm.DB {
	if storeID > 0 {
		return query.Where("store_id = ?", storeID)
	}
	return query
}

// GetInventoryStats 获取库存统计
func (m *StatisticsModule) GetInventoryStats(storeID uint) (*model.InventoryStats, error) {
	stats := &model.InventoryStats{}

	var inventorySummary struct {
		TotalProducts int64
		TotalQuantity float64
	}
	inventoryQuery := withStoreID(m.db.Model(&model.Inventory{}).Where("deleted_at IS NULL"), storeID)
	if err := inventoryQuery.Select("COUNT(*) AS total_products, COALESCE(SUM(quantity), 0) AS total_quantity").Scan(&inventorySummary).Error; err != nil {
		return nil, err
	}
	stats.TotalProducts = inventorySummary.TotalProducts
	stats.TotalQuantity = inventorySummary.TotalQuantity

	today := time.Now().Format("2006-01-02")
	var orderSummary struct {
		TotalRecords int64
		TodayIn      float64
		TodayOut     float64
	}
	orderQuery := withStoreID(m.db.Model(&model.InventoryOrder{}).Where("deleted_at IS NULL"), storeID)
	if err := orderQuery.Select(`
		COUNT(*) AS total_records,
		COALESCE(SUM(CASE WHEN type = ? AND created_at >= ? AND created_at < DATE_ADD(?, INTERVAL 1 DAY) THEN total_quantity ELSE 0 END), 0) AS today_in,
		COALESCE(SUM(CASE WHEN type = ? AND created_at >= ? AND created_at < DATE_ADD(?, INTERVAL 1 DAY) THEN total_quantity ELSE 0 END), 0) AS today_out
	`, model.InventoryTypeIn, today, today, model.InventoryTypeOut, today, today).Scan(&orderSummary).Error; err != nil {
		return nil, err
	}
	stats.TotalRecords = orderSummary.TotalRecords
	stats.TodayIn = orderSummary.TodayIn
	stats.TodayOut = orderSummary.TodayOut

	return stats, nil
}

// GetSalesStats 获取销售统计
func (m *StatisticsModule) GetSalesStats(storeID uint, startDate, endDate string) (*model.SalesStats, error) {
	stats := &model.SalesStats{}

	query := withStoreID(m.db.Model(&model.StoreAccount{}).Where("deleted_at IS NULL AND is_canceled = 0"), storeID)
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}

	var summary struct {
		TotalOrders int64
		TotalAmount float64
		TotalQty    float64
	}
	if err := query.Select("COUNT(*) AS total_orders, COALESCE(SUM(total_amount), 0) AS total_amount, COALESCE(SUM(item_count), 0) AS total_qty").Scan(&summary).Error; err != nil {
		return nil, err
	}
	stats.TotalOrders = summary.TotalOrders
	stats.TotalAmount = summary.TotalAmount
	stats.TotalQty = summary.TotalQty

	// 平均客单价
	if stats.TotalOrders > 0 {
		stats.AvgAmount = stats.TotalAmount / float64(stats.TotalOrders)
	}

	// 今日销售额
	today := businessdate.DateString(time.Now())
	todayQuery := withStoreID(m.db.Model(&model.StoreAccount{}).Where("deleted_at IS NULL AND is_canceled = 0 AND account_date >= ? AND account_date < DATE_ADD(?, INTERVAL 1 DAY)", today, today), storeID)
	if err := todayQuery.Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TodayAmount).Error; err != nil {
		return nil, err
	}

	// 本月销售额
	businessToday := businessdate.Date(time.Now())
	monthStart := time.Date(businessToday.Year(), businessToday.Month(), 1, 0, 0, 0, 0, businessToday.Location()).Format("2006-01-02")
	monthQuery := withStoreID(m.db.Model(&model.StoreAccount{}).Where("deleted_at IS NULL AND is_canceled = 0 AND account_date >= ?", monthStart), storeID)
	if err := monthQuery.Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.MonthAmount).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetSalesTrend 获取销售趋势
func (m *StatisticsModule) GetSalesTrend(storeID uint, startDate, endDate, period string) ([]model.SalesTrendItem, error) {
	var results []model.SalesTrendItem

	dateFormat := "%Y-%m-%d"
	if period == "month" || period == "quarter" {
		dateFormat = "%Y-%m-%d"
	} else if period == "year" {
		dateFormat = "%Y-%m"
	}

	query := withStoreID(m.db.Model(&model.StoreAccount{}).
		Select("DATE_FORMAT(account_date, ?) as date, COALESCE(SUM(total_amount), 0) as amount, COUNT(*) as orders", dateFormat).
		Where("deleted_at IS NULL AND is_canceled = 0"), storeID)
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}

	if err := query.Group("date").Order("date ASC").Scan(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}

// GetSalesTrendByGranularity 按粒度获取销售趋势
func (m *StatisticsModule) GetSalesTrendByGranularity(storeID uint, startDate, endDate, granularity string) ([]model.SalesTrendItem, error) {
	var results []model.SalesTrendItem
	dateFormat := "%Y-%m-%d"
	if granularity == "month" {
		dateFormat = "%Y-%m"
	}

	query := withStoreID(m.db.Model(&model.StoreAccount{}).
		Select("DATE_FORMAT(account_date, ?) as date, COALESCE(SUM(total_amount), 0) as amount, COUNT(*) as orders", dateFormat).
		Where("deleted_at IS NULL AND is_canceled = 0"), storeID)
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}
	if err := query.Group("date").Order("date ASC").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

// GetChannelStats 获取渠道统计
func (m *StatisticsModule) GetChannelStats(storeID uint, startDate, endDate string) ([]model.ChannelStatsItem, error) {
	var results []model.ChannelStatsItem

	query := withStoreID(m.db.Model(&model.StoreAccount{}).
		Select("channel, COALESCE(SUM(total_amount), 0) as amount, COUNT(*) as orders").
		Where("deleted_at IS NULL AND is_canceled = 0"), storeID)
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}

	if err := query.Group("channel").Order("amount DESC").Scan(&results).Error; err != nil {
		return nil, err
	}

	// 计算总额和占比
	var totalAmount float64
	for _, item := range results {
		totalAmount += item.Amount
	}

	// 获取渠道名称映射
	channelMap, err := m.getChannelNameMap()
	if err != nil {
		return nil, err
	}

	for i := range results {
		if totalAmount > 0 {
			results[i].Percent = results[i].Amount / totalAmount * 100
		}
		if name, ok := channelMap[results[i].Channel]; ok {
			results[i].ChannelName = name
		} else {
			results[i].ChannelName = results[i].Channel
		}
	}

	return results, nil
}

// getChannelNameMap 获取渠道名称映射
func (m *StatisticsModule) getChannelNameMap() (map[string]string, error) {
	channelMap := make(map[string]string)

	var dictData []model.DictData
	if err := m.db.Where("type_code = ? AND status = 1", "sales_channel").Find(&dictData).Error; err != nil {
		return nil, err
	}

	for _, d := range dictData {
		channelMap[d.Value] = d.Label
	}

	return channelMap, nil
}

// GetBusinessOverview 获取经营总览统计（按日期）
func (m *StatisticsModule) GetBusinessOverview(storeID uint, startDate, endDate string) (*model.BusinessOverviewStats, error) {
	stats := &model.BusinessOverviewStats{
		StartDate:              startDate,
		EndDate:                endDate,
		StoreID:                storeID,
		StoreExpenseCategories: make([]model.StoreExpenseCategoryAmountItem, 0),
	}

	var categoryRows []categoryAmountRow
	categorySQL := `
SELECT
	COALESCE(sp.category_id, 0) AS category_id,
	COALESCE(sc.name, '未分类') AS category_name,
	COALESCE(SUM(CASE WHEN io.type = 1 THEN ioi.quantity * COALESCE(sp.price, 0) ELSE 0 END), 0) AS in_amount,
	COALESCE(SUM(CASE WHEN io.type = 2 THEN ioi.quantity * COALESCE(sp.price, 0) ELSE 0 END), 0) AS out_amount
FROM inventory_order_items ioi
JOIN inventory_orders io ON io.id = ioi.order_id AND io.deleted_at IS NULL
LEFT JOIN supplier_products sp ON sp.id = ioi.product_id
LEFT JOIN supplier_categories sc ON sc.id = sp.category_id
WHERE io.created_at >= ? AND io.created_at < DATE_ADD(?, INTERVAL 1 DAY)
`
	args := []interface{}{startDate, endDate}
	if storeID > 0 {
		categorySQL += " AND io.store_id = ?"
		args = append(args, storeID)
	}
	categorySQL += " GROUP BY COALESCE(sp.category_id, 0), COALESCE(sc.name, '未分类') ORDER BY in_amount DESC, out_amount DESC"
	if err := m.db.Raw(categorySQL, args...).Scan(&categoryRows).Error; err != nil {
		return nil, err
	}

	stats.Categories = make([]model.CategoryAmountItem, 0, len(categoryRows))
	for _, row := range categoryRows {
		item := model.CategoryAmountItem{
			CategoryID:   row.CategoryID,
			CategoryName: row.CategoryName,
			InAmount:     row.InAmount,
			OutAmount:    row.OutAmount,
			NetAmount:    row.OutAmount - row.InAmount,
		}
		stats.Categories = append(stats.Categories, item)
		stats.InboundAmount += row.InAmount
		stats.OutboundAmount += row.OutAmount
	}
	stats.AllCategoryAmount = stats.InboundAmount

	salesQuery := m.db.Model(&model.StoreAccount{}).
		Where("deleted_at IS NULL AND is_canceled = 0 AND account_date >= ? AND account_date <= ?", startDate, endDate)
	if storeID > 0 {
		salesQuery = salesQuery.Where("store_id = ?", storeID)
	}
	var salesSummary struct {
		SalesOrderCount    int64
		SalesAmount        float64
		OtherExpenseAmount float64
		ErrandFeeAmount    float64
		RoundAmount        float64
		GiftWineCostAmount float64
	}
	if err := salesQuery.Select("COUNT(*) AS sales_order_count, COALESCE(SUM(total_amount), 0) AS sales_amount, COALESCE(SUM(other_expense_amount), 0) AS other_expense_amount, COALESCE(SUM(errand_fee), 0) AS errand_fee_amount, COALESCE(SUM(round_amount), 0) AS round_amount, COALESCE(SUM(gift_wine_cost_amount), 0) AS gift_wine_cost_amount").Scan(&salesSummary).Error; err != nil {
		return nil, err
	}
	stats.SalesOrderCount = salesSummary.SalesOrderCount
	stats.SalesAmount = salesSummary.SalesAmount
	stats.OtherExpenseAmount = salesSummary.OtherExpenseAmount
	stats.ErrandFeeAmount = salesSummary.ErrandFeeAmount
	stats.RoundAmount = salesSummary.RoundAmount
	stats.GiftWineCostAmount = salesSummary.GiftWineCostAmount
	consumableQuery := m.db.Table("store_account_consumables AS sac").
		Joins("JOIN store_accounts AS sa ON sa.id = sac.account_id AND sa.deleted_at IS NULL AND sa.is_canceled = 0").
		Where("sa.account_date >= ? AND sa.account_date <= ?", startDate, endDate)
	if storeID > 0 {
		consumableQuery = consumableQuery.Where("sa.store_id = ?", storeID)
	}
	if err := consumableQuery.Select("COALESCE(SUM(sac.amount), 0)").Scan(&stats.ConsumableAmount).Error; err != nil {
		return nil, err
	}

	expenseQuery := m.db.Model(&model.StoreExpense{}).
		Where("deleted_at IS NULL AND expense_date >= ? AND expense_date <= ?", startDate, endDate)
	if storeID > 0 {
		expenseQuery = expenseQuery.Where("store_id = ?", storeID)
	}
	var expenseSummary struct {
		StoreExpenseAmount     float64
		TakeoutPromotionAmount float64
	}
	if err := expenseQuery.Select("COALESCE(SUM(amount), 0) AS store_expense_amount, COALESCE(SUM(CASE WHEN category_code = ? THEN amount ELSE 0 END), 0) AS takeout_promotion_amount", "takeout_promotion").Scan(&expenseSummary).Error; err != nil {
		return nil, err
	}
	stats.StoreExpenseAmount = expenseSummary.StoreExpenseAmount
	stats.TakeoutPromotionAmount = expenseSummary.TakeoutPromotionAmount

	expenseCategorySQL := `
SELECT
	dd.value AS category_code,
	dd.label AS category_name,
	COALESCE(SUM(se.amount), 0) AS amount,
	COUNT(se.id) AS count
FROM dict_data dd
LEFT JOIN store_expenses se
	ON se.category_code = dd.value
	AND se.deleted_at IS NULL
	AND se.expense_date >= ?
	AND se.expense_date <= ?
`
	expenseCategoryArgs := []interface{}{startDate, endDate}
	if storeID > 0 {
		expenseCategorySQL += " AND se.store_id = ?"
		expenseCategoryArgs = append(expenseCategoryArgs, storeID)
	}
	expenseCategorySQL += `
WHERE dd.type_code = ? AND dd.status = 1
GROUP BY dd.id, dd.value, dd.label, dd.sort
ORDER BY dd.sort ASC, dd.id ASC
`
	expenseCategoryArgs = append(expenseCategoryArgs, model.StoreExpenseCategoryDictCode)
	if err := m.db.Raw(expenseCategorySQL, expenseCategoryArgs...).Scan(&stats.StoreExpenseCategories).Error; err != nil {
		return nil, err
	}

	takeoutSalesQuery := m.db.Model(&model.StoreAccount{}).
		Where(`deleted_at IS NULL AND is_canceled = 0 AND account_date >= ? AND account_date <= ? AND (
			LOWER(channel) REGEXP ? OR
			channel LIKE ? OR channel LIKE ? OR channel LIKE ? OR channel LIKE ? OR channel LIKE ? OR
			channel LIKE ? OR channel LIKE ? OR channel LIKE ? OR channel LIKE ? OR channel LIKE ?
		)`, startDate, endDate,
			`(^|[^a-z0-9])(takeout|waimai|meituan|eleme|elm|taobao|tb|flash|shangou|jd|jingdong|douyin|tiktok|tuangou|groupbuy|group_buy|groupon|wechat_mini|mini_program|miniprogram|mall)([^a-z0-9]|$)`,
			"%外卖%", "%美团%", "%饿了么%", "%闪购%", "%淘宝%",
			"%抖音%", "%团购%", "%微信小程序%", "%小程序%", "%商城%")
	if storeID > 0 {
		takeoutSalesQuery = takeoutSalesQuery.Where("store_id = ?", storeID)
	}
	if err := takeoutSalesQuery.Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TakeoutSalesAmount).Error; err != nil {
		return nil, err
	}
	if stats.TakeoutPromotionAmount > 0 {
		stats.TakeoutPromotionROI = stats.TakeoutSalesAmount / stats.TakeoutPromotionAmount
	}

	b2bQuery := m.db.Model(&model.B2BSupplyOrder{}).
		Where("deleted_at IS NULL AND delivery_status <> ? AND order_date >= ? AND order_date <= ?", model.B2BDeliveryCancel, startDate, endDate)
	if storeID > 0 {
		b2bQuery = b2bQuery.Where("store_id = ?", storeID)
	}
	var b2bSummary struct {
		B2BSupplyOrderCount int64
		B2BSupplyAmount     float64
	}
	if err := b2bQuery.Select("COUNT(*) AS b2b_supply_order_count, COALESCE(SUM(total_amount), 0) AS b2b_supply_amount").Scan(&b2bSummary).Error; err != nil {
		return nil, err
	}
	stats.B2BSupplyOrderCount = b2bSummary.B2BSupplyOrderCount
	stats.B2BSupplyAmount = b2bSummary.B2BSupplyAmount

	returnQuery := m.db.Model(&model.StoreReturn{}).
		Where("deleted_at IS NULL AND return_date >= ? AND return_date <= ?", startDate, endDate)
	if storeID > 0 {
		returnQuery = returnQuery.Where("store_id = ?", storeID)
	}
	if err := returnQuery.Select("COALESCE(SUM(total_deposit), 0), COALESCE(SUM(logistics_fee), 0)").
		Row().Scan(&stats.ReturnDepositAmount, &stats.ReturnLogisticsFee); err != nil {
		return nil, err
	}

	lossBaseQuery := m.db.Model(&model.InventoryLossOrder{}).
		Where("deleted_at IS NULL AND is_canceled = 0 AND created_at >= ? AND created_at < DATE_ADD(?, INTERVAL 1 DAY)", startDate, endDate)
	if storeID > 0 {
		lossBaseQuery = lossBaseQuery.Where("store_id = ?", storeID)
	}
	var lossSummary struct {
		InventoryLossCount     int64
		InventoryLossAmount    float64
		InventorySelfUseCount  int64
		InventorySelfUseAmount float64
	}
	if err := lossBaseQuery.Select(`
		COUNT(CASE WHEN type = ? THEN 1 END) AS inventory_loss_count,
		COALESCE(SUM(CASE WHEN type = ? THEN total_cost ELSE 0 END), 0) AS inventory_loss_amount,
		COUNT(CASE WHEN type = ? THEN 1 END) AS inventory_self_use_count,
		COALESCE(SUM(CASE WHEN type = ? THEN total_cost ELSE 0 END), 0) AS inventory_self_use_amount
	`, model.InventoryLossTypeLoss, model.InventoryLossTypeLoss, model.InventoryLossTypeSelfUse, model.InventoryLossTypeSelfUse).Scan(&lossSummary).Error; err != nil {
		return nil, err
	}
	stats.InventoryLossCount = lossSummary.InventoryLossCount
	stats.InventoryLossAmount = lossSummary.InventoryLossAmount
	stats.InventorySelfUseCount = lossSummary.InventorySelfUseCount
	stats.InventorySelfUseAmount = lossSummary.InventorySelfUseAmount

	memberRankQuery := m.db.Table("store_accounts AS sa").
		Select(`
			COALESCE(tm.id, 0) AS member_id,
			COALESCE(NULLIF(tm.name, ''), '未知会员') AS member_name,
			COALESCE(tm.phone, '') AS member_phone,
			COALESCE(SUM(sa.total_amount), 0) AS amount,
			COUNT(sa.id) AS orders
		`).
		Joins("LEFT JOIN t_member AS tm ON tm.id = sa.member_id").
		Where("sa.deleted_at IS NULL AND sa.is_canceled = 0 AND sa.member_id IS NOT NULL AND sa.account_date >= ? AND sa.account_date <= ?", startDate, endDate)
	if storeID > 0 {
		memberRankQuery = memberRankQuery.Where("sa.store_id = ?", storeID)
	}
	if err := memberRankQuery.
		Group("tm.id, tm.name, tm.phone").
		Order("amount DESC").
		Limit(10).
		Scan(&stats.MemberConsumptionRank).Error; err != nil {
		return nil, err
	}

	var itemCostAmount float64
	itemCostQuery := m.db.Table("store_account_items AS sai").
		Joins("JOIN store_accounts AS sa ON sa.id = sai.account_id AND sa.deleted_at IS NULL AND sa.is_canceled = 0").
		Joins("LEFT JOIN product_unit_specs AS ps ON ps.product_id = sai.product_id AND ps.is_enabled = 1 AND (ps.unit_code = sai.unit OR ps.unit_name = sai.unit)").
		Where("sa.account_date >= ? AND sa.account_date <= ?", startDate, endDate)
	if storeID > 0 {
		itemCostQuery = itemCostQuery.Where("sa.store_id = ?", storeID)
	}
	if err := itemCostQuery.Select("COALESCE(SUM(sai.quantity * COALESCE(ps.cost_price, 0)), 0)").Scan(&itemCostAmount).Error; err != nil {
		return nil, err
	}

	inOutQuery := m.db.Model(&model.InventoryOrder{}).
		Where("deleted_at IS NULL AND created_at >= ? AND created_at < DATE_ADD(?, INTERVAL 1 DAY)", startDate, endDate)
	if storeID > 0 {
		inOutQuery = inOutQuery.Where("store_id = ?", storeID)
	}
	var inOutSummary struct {
		InventoryInCount  int64
		InventoryOutCount int64
	}
	if err := inOutQuery.Select(`
		COUNT(CASE WHEN type = ? THEN 1 END) AS inventory_in_count,
		COUNT(CASE WHEN type = ? THEN 1 END) AS inventory_out_count
	`, model.InventoryTypeIn, model.InventoryTypeOut).Scan(&inOutSummary).Error; err != nil {
		return nil, err
	}
	stats.InventoryInCount = inOutSummary.InventoryInCount
	stats.InventoryOutCount = inOutSummary.InventoryOutCount

	stats.GrossProfitAmount = stats.SalesAmount - itemCostAmount
	stats.NetProfitAmount = calculateBusinessOverviewNetProfit(stats, itemCostAmount)

	return stats, nil
}

func calculateBusinessOverviewNetProfit(stats *model.BusinessOverviewStats, itemCostAmount float64) float64 {
	if stats == nil {
		return 0
	}
	// 门店支出在大屏单独展示；记账净利保持与有效记账单的净利润口径一致。
	return stats.SalesAmount - stats.OtherExpenseAmount - stats.ErrandFeeAmount - stats.ConsumableAmount - itemCostAmount - stats.GiftWineCostAmount - stats.RoundAmount
}
