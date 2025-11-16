package module

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils"

	"gorm.io/gorm"
)

type MenuReportModule struct {
	db *gorm.DB
}

func NewMenuReportModule(db *gorm.DB) *MenuReportModule {
	return &MenuReportModule{db: db}
}

// CreateOrder 创建报菜记录单（包含详情）
func (m *MenuReportModule) CreateOrder(order *model.MenuReportOrder) error {
	err := m.db.Transaction(func(tx *gorm.DB) error {
		// GORM 会自动保存关联的 Items，无需手动循环插入
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 加载关联信息（用于通知）
		tx.Preload("Items.Dish").Preload("Store").Preload("User").First(order, order.ID)

		return nil
	})

	if err != nil {
		return err
	}

	// 发送通知（在事务成功后异步执行，不影响主流程）
	utils.GlobalEventBus.NotifyAsync(utils.EventMenuReportCreated, order)

	return nil
}

// GetOrderByID 根据ID获取报菜记录单（带门店隔离）
func (m *MenuReportModule) GetOrderByID(id uint, storeID uint) (*model.MenuReportOrder, error) {
	var order model.MenuReportOrder
	if err := m.db.Preload("Items.Dish").Preload("Store").Preload("User").
		Where("id = ? AND store_id = ?", id, storeID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// ListOrders 获取报菜记录单列表（带门店隔离）
func (m *MenuReportModule) ListOrders(storeID uint, page, pageSize int) ([]*model.MenuReportOrder, int64, error) {
	var orders []*model.MenuReportOrder
	var total int64

	query := m.db.Model(&model.MenuReportOrder{}).Where("store_id = ?", storeID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Items.Dish").Preload("Store").Preload("User").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// ListOrdersByDateRange 根据日期范围查询（带门店隔离）
func (m *MenuReportModule) ListOrdersByDateRange(storeID uint, startDate, endDate time.Time) ([]*model.MenuReportOrder, error) {
	var orders []*model.MenuReportOrder
	if err := m.db.Preload("Items.Dish").Preload("Store").Preload("User").
		Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, startDate, endDate).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// UpdateOrder 更新报菜记录单（带门店隔离）
func (m *MenuReportModule) UpdateOrder(id uint, storeID uint, remark *string) error {
	updates := make(map[string]interface{})
	if remark != nil {
		updates["remark"] = *remark
	}

	return m.db.Model(&model.MenuReportOrder{}).
		Where("id = ? AND store_id = ?", id, storeID).
		Updates(updates).Error
}

// DeleteOrder 删除报菜记录单（带门店隔离）
func (m *MenuReportModule) DeleteOrder(id uint, storeID uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("report_order_id = ?", id).Delete(&model.MenuReportItem{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND store_id = ?", id, storeID).Delete(&model.MenuReportOrder{}).Error
	})
}

// AddItem 添加报菜详情项
func (m *MenuReportModule) AddItem(item *model.MenuReportItem) error {
	return m.db.Create(item).Error
}

// UpdateItem 更新报菜详情项
func (m *MenuReportModule) UpdateItem(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.MenuReportItem{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteItem 删除报菜详情项
func (m *MenuReportModule) DeleteItem(id uint) error {
	return m.db.Where("id = ?", id).Delete(&model.MenuReportItem{}).Error
}

// GetStatsByDateRange 获取日期范围内的统计数据（带门店隔离）
func (m *MenuReportModule) GetStatsByDateRange(storeID uint, startDate, endDate time.Time) (*model.MenuReportStats, error) {
	var stats model.MenuReportStats

	err := m.db.Model(&model.MenuReportItem{}).
		Select("COUNT(DISTINCT DATE(menu_report_items.created_at)) as total_days, SUM(quantity) as total_qty").
		Joins("LEFT JOIN menu_report_orders ON menu_report_items.report_order_id = menu_report_orders.id").
		Where("menu_report_orders.store_id = ? AND menu_report_items.created_at BETWEEN ? AND ?", storeID, startDate, endDate).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetStatsByDateRangeAllStores 获取所有门店的统计数据（仅总部可用）
func (m *MenuReportModule) GetStatsByDateRangeAllStores(startDate, endDate time.Time) (*model.MenuReportStats, error) {
	var stats model.MenuReportStats

	err := m.db.Model(&model.MenuReportItem{}).
		Select("COUNT(DISTINCT DATE(menu_report_items.created_at)) as total_days, SUM(quantity) as total_qty").
		Joins("LEFT JOIN menu_report_orders ON menu_report_items.report_order_id = menu_report_orders.id").
		Where("menu_report_items.created_at BETWEEN ? AND ?", startDate, endDate).
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// Create 创建报菜记录
func (m *MenuReportModule) Create(report *model.MenuReport) error {
	return m.db.Create(report).Error
}

// List 获取报菜记录列表（带门店隔离）
func (m *MenuReportModule) List(storeID uint, page, pageSize int) ([]*model.MenuReport, int64, error) {
	var reports []*model.MenuReport
	var total int64

	query := m.db.Model(&model.MenuReport{}).Where("store_id = ?", storeID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&reports).Error; err != nil {
		return nil, 0, err
	}

	return reports, total, nil
}

// ListByDateRange 根据日期范围查询报菜记录（带门店隔离）
func (m *MenuReportModule) ListByDateRange(storeID uint, startDate, endDate time.Time) ([]*model.MenuReport, error) {
	var reports []*model.MenuReport
	if err := m.db.Where("store_id = ? AND created_at BETWEEN ? AND ?", storeID, startDate, endDate).
		Order("created_at DESC").
		Find(&reports).Error; err != nil {
		return nil, err
	}
	return reports, nil
}

// GetByID 根据ID获取报菜记录（带门店隔离）
func (m *MenuReportModule) GetByID(id uint, storeID uint) (*model.MenuReport, error) {
	var report model.MenuReport
	if err := m.db.Where("id = ? AND store_id = ?", id, storeID).First(&report).Error; err != nil {
		return nil, err
	}
	return &report, nil
}

// Update 更新报菜记录（带门店隔离）
func (m *MenuReportModule) Update(id uint, storeID uint, req *model.UpdateMenuReportReq) error {
	updates := make(map[string]interface{})
	if req.DishID != nil {
		updates["dish_id"] = *req.DishID
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	return m.db.Model(&model.MenuReport{}).
		Where("id = ? AND store_id = ?", id, storeID).
		Updates(updates).Error
}

// Delete 删除报菜记录（带门店隔离）
func (m *MenuReportModule) Delete(id uint, storeID uint) error {
	return m.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&model.MenuReport{}).Error
}
