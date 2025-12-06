package service

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type StatisticsService struct {
	statisticsModule *module.StatisticsModule
}

func NewStatisticsService(statisticsModule *module.StatisticsModule) *StatisticsService {
	return &StatisticsService{statisticsModule: statisticsModule}
}

// GetDashboard 获取统计面板数据
func (s *StatisticsService) GetDashboard(storeID uint, period string) (*model.DashboardStats, error) {
	inventoryStats, err := s.GetInventoryStats(storeID)
	if err != nil {
		return nil, err
	}

	salesStats, err := s.GetSalesStats(storeID, period)
	if err != nil {
		return nil, err
	}

	return &model.DashboardStats{
		Inventory: *inventoryStats,
		Sales:     *salesStats,
	}, nil
}

// GetInventoryStats 获取库存统计
func (s *StatisticsService) GetInventoryStats(storeID uint) (*model.InventoryStats, error) {
	return s.statisticsModule.GetInventoryStats(storeID)
}

// GetSalesStats 获取销售统计
func (s *StatisticsService) GetSalesStats(storeID uint, period string) (*model.SalesStats, error) {
	startDate, endDate := s.getPeriodRange(period)
	stats, err := s.statisticsModule.GetSalesStats(storeID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	stats.PeriodLabel = s.getPeriodLabel(period)
	return stats, nil
}

// GetSalesTrend 获取销售趋势
func (s *StatisticsService) GetSalesTrend(storeID uint, period string) ([]model.SalesTrendItem, error) {
	startDate, endDate := s.getPeriodRange(period)
	return s.statisticsModule.GetSalesTrend(storeID, startDate, endDate, period)
}

// GetChannelStats 获取渠道统计
func (s *StatisticsService) GetChannelStats(storeID uint, period string) ([]model.ChannelStatsItem, error) {
	startDate, endDate := s.getPeriodRange(period)
	return s.statisticsModule.GetChannelStats(storeID, startDate, endDate)
}

// getPeriodRange 获取周期的日期范围
func (s *StatisticsService) getPeriodRange(period string) (string, string) {
	now := time.Now()
	var startDate, endDate time.Time

	switch period {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 0, 1).Add(-time.Second)
	case "week":
		// 本周一到今天
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		startDate = time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
		endDate = now
	case "month":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	case "quarter":
		quarter := (int(now.Month()) - 1) / 3
		startDate = time.Date(now.Year(), time.Month(quarter*3+1), 1, 0, 0, 0, 0, now.Location())
		endDate = now
	case "year":
		startDate = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endDate = now
	default:
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 0, 1).Add(-time.Second)
	}

	return startDate.Format("2006-01-02"), endDate.Format("2006-01-02")
}

// getPeriodLabel 获取周期标签
func (s *StatisticsService) getPeriodLabel(period string) string {
	switch period {
	case "today":
		return "今日"
	case "week":
		return "本周"
	case "month":
		return "本月"
	case "quarter":
		return "本季度"
	case "year":
		return "本年"
	default:
		return "今日"
	}
}
