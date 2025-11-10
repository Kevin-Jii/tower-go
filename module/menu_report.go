package module

import (
	"time"
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

type MenuReportModule struct {
	db *gorm.DB
}

func NewMenuReportModule(db *gorm.DB) *MenuReportModule {
	return &MenuReportModule{db: db}
}

// Create 创建报菜记录
func (m *MenuReportModule) Create(report *model.MenuReport) error {
	return m.db.Create(report).Error
}

// GetByID 根据ID获取报菜记录（带门店隔离）
func (m *MenuReportModule) GetByID(id uint, storeID uint) (*model.MenuReport, error) {
	var report model.MenuReport
	if err := m.db.Preload("Dish").Where("id = ? AND store_id = ?", id, storeID).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// List 获取报菜记录列表（带门店隔离）
func (m *MenuReportModule) List(storeID uint, page, pageSize int) ([]*model.MenuReport, int64, error) {
	var reports []*model.MenuReport
	var total int64

	query := m.db.Model(&model.MenuReport{}).Where("store_id = ?", storeID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Dish").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// ListByDateRange 根据日期范围查询（带门店隔离）
func (m *MenuReportModule) ListByDateRange(storeID uint, startDate, endDate time.Time) ([]*model.MenuReport, error) {
	var reports []*model.MenuReport
	if err := m.db.Preload("Dish").
		Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, startDate, endDate).
		Order("created_at DESC").
		Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

// Update 更新报菜记录（带门店隔离）
func (m *MenuReportModule) Update(id uint, storeID uint, req *model.UpdateMenuReportReq) error {
	updates := make(map[string]interface{})

	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	return m.db.Model(&model.MenuReport{}).
		Where("id = ? AND store_id = ?", id, storeID).
		Updates(updates).Error
} // Delete 删除报菜记录（带门店隔离）
func (m *MenuReportModule) Delete(id uint, storeID uint) error {
	return m.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&model.MenuReport{}).Error
}

// GetStatsByDateRange 获取日期范围内的统计数据（带门店隔离）
func (m *MenuReportModule) GetStatsByDateRange(storeID uint, startDate, endDate time.Time) (*model.MenuReportStats, error) {
	var stats model.MenuReportStats

	err := m.db.Model(&model.MenuReport{}).
		Select("COUNT(*) as total_days, SUM(quantity) as total_qty").
		Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, startDate, endDate).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetStatsByDateRangeAllStores 获取所有门店的统计数据（仅总部可用）
func (m *MenuReportModule) GetStatsByDateRangeAllStores(startDate, endDate time.Time) (*model.MenuReportStats, error) {
	var stats model.MenuReportStats

	err := m.db.Model(&model.MenuReport{}).
		Select("COUNT(*) as total_days, SUM(quantity) as total_qty").
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}
