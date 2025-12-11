package module

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type StatisticsModule struct {
	db *gorm.DB
}

func NewStatisticsModule(db *gorm.DB) *StatisticsModule {
	return &StatisticsModule{db: db}
}

// GetInventoryStats 获取库存统计
func (m *StatisticsModule) GetInventoryStats(storeID uint) (*model.InventoryStats, error) {
	stats := &model.InventoryStats{}

	// 商品种类数和总库存
	query := m.db.Model(&model.Inventory{}).Where("deleted_at IS NULL")
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	query.Count(&stats.TotalProducts)
	query.Select("COALESCE(SUM(quantity), 0)").Scan(&stats.TotalQuantity)

	// 出入库单总数
	orderQuery := m.db.Model(&model.InventoryOrder{}).Where("deleted_at IS NULL")
	if storeID > 0 {
		orderQuery = orderQuery.Where("store_id = ?", storeID)
	}
	orderQuery.Count(&stats.TotalRecords)

	// 今日入库/出库（从出入库单统计）
	today := time.Now().Format("2006-01-02")

	// 今日入库
	m.db.Model(&model.InventoryOrder{}).
		Where("deleted_at IS NULL AND DATE(created_at) = ? AND type = ?", today, model.InventoryTypeIn).
		Where(func(db *gorm.DB) *gorm.DB {
			if storeID > 0 {
				return db.Where("store_id = ?", storeID)
			}
			return db
		}(m.db)).
		Select("COALESCE(SUM(total_quantity), 0)").Scan(&stats.TodayIn)

	// 今日出库
	m.db.Model(&model.InventoryOrder{}).
		Where("deleted_at IS NULL AND DATE(created_at) = ? AND type = ?", today, model.InventoryTypeOut).
		Where(func(db *gorm.DB) *gorm.DB {
			if storeID > 0 {
				return db.Where("store_id = ?", storeID)
			}
			return db
		}(m.db)).
		Select("COALESCE(SUM(total_quantity), 0)").Scan(&stats.TodayOut)

	return stats, nil
}

// GetSalesStats 获取销售统计
func (m *StatisticsModule) GetSalesStats(storeID uint, startDate, endDate string) (*model.SalesStats, error) {
	stats := &model.SalesStats{}

	query := m.db.Model(&model.StoreAccount{}).Where("deleted_at IS NULL")
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}

	// 总订单数
	query.Count(&stats.TotalOrders)

	// 总销售额和总数量（使用新字段 total_amount 和 item_count）
	m.db.Model(&model.StoreAccount{}).
		Where("deleted_at IS NULL").
		Where(func(db *gorm.DB) *gorm.DB {
			if storeID > 0 {
				db = db.Where("store_id = ?", storeID)
			}
			if startDate != "" {
				db = db.Where("account_date >= ?", startDate)
			}
			if endDate != "" {
				db = db.Where("account_date <= ?", endDate)
			}
			return db
		}(m.db)).
		Select("COALESCE(SUM(total_amount), 0) as total_amount, COALESCE(SUM(item_count), 0) as total_qty").
		Row().Scan(&stats.TotalAmount, &stats.TotalQty)

	// 平均客单价
	if stats.TotalOrders > 0 {
		stats.AvgAmount = stats.TotalAmount / float64(stats.TotalOrders)
	}

	// 今日销售额
	today := time.Now().Format("2006-01-02")
	fmt.Printf("🔍 [Statistics] 今日日期: %s, storeID: %d\n", today, storeID)
	todayQuery := m.db.Model(&model.StoreAccount{}).Where("deleted_at IS NULL AND DATE(account_date) = ?", today)
	if storeID > 0 {
		todayQuery = todayQuery.Where("store_id = ?", storeID)
	}
	todayQuery.Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.TodayAmount)
	fmt.Printf("🔍 [Statistics] 今日销售额: %.2f\n", stats.TodayAmount)

	// 本月销售额
	monthStart := time.Now().Format("2006-01") + "-01"
	fmt.Printf("🔍 [Statistics] 本月开始: %s\n", monthStart)
	monthQuery := m.db.Model(&model.StoreAccount{}).
		Where("deleted_at IS NULL AND DATE(account_date) >= ?", monthStart)
	if storeID > 0 {
		monthQuery = monthQuery.Where("store_id = ?", storeID)
	}
	monthQuery.Select("COALESCE(SUM(total_amount), 0)").Scan(&stats.MonthAmount)
	fmt.Printf("🔍 [Statistics] 本月销售额: %.2f\n", stats.MonthAmount)

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

	query := m.db.Model(&model.StoreAccount{}).
		Select("DATE_FORMAT(account_date, ?) as date, COALESCE(SUM(total_amount), 0) as amount, COUNT(*) as orders", dateFormat).
		Where("deleted_at IS NULL")

	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}

	query.Group("date").Order("date ASC").Scan(&results)

	return results, nil
}

// GetChannelStats 获取渠道统计
func (m *StatisticsModule) GetChannelStats(storeID uint, startDate, endDate string) ([]model.ChannelStatsItem, error) {
	var results []model.ChannelStatsItem

	query := m.db.Model(&model.StoreAccount{}).
		Select("channel, COALESCE(SUM(total_amount), 0) as amount, COUNT(*) as orders").
		Where("deleted_at IS NULL")

	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}

	query.Group("channel").Order("amount DESC").Scan(&results)

	// 计算总额和占比
	var totalAmount float64
	for _, item := range results {
		totalAmount += item.Amount
	}

	// 获取渠道名称映射
	channelMap := m.getChannelNameMap()

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
func (m *StatisticsModule) getChannelNameMap() map[string]string {
	channelMap := make(map[string]string)

	var dictData []model.DictData
	m.db.Where("type_code = ? AND status = 1", "sales_channel").Find(&dictData)

	for _, d := range dictData {
		channelMap[d.Value] = d.Label
	}

	return channelMap
}
