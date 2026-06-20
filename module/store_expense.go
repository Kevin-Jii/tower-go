package module

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type StoreExpenseModule struct {
	db *gorm.DB
}

func NewStoreExpenseModule(db *gorm.DB) *StoreExpenseModule {
	return &StoreExpenseModule{db: db}
}

func (m *StoreExpenseModule) Create(record *model.StoreExpense) error {
	return m.db.Create(record).Error
}

func (m *StoreExpenseModule) GetByIDScoped(id, storeID uint, hqUnbound bool) (*model.StoreExpense, error) {
	var record model.StoreExpense
	query := m.db.Preload("Store").Preload("Operator").Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.First(&record).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

func (m *StoreExpenseModule) List(req *model.ListStoreExpenseReq) ([]*model.StoreExpense, int64, error) {
	records := make([]*model.StoreExpense, 0)
	var total int64

	query := m.db.Model(&model.StoreExpense{}).Preload("Store").Preload("Operator")
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.CategoryCode != "" {
		query = query.Where("category_code = ?", req.CategoryCode)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("expense_no LIKE ? OR category_name LIKE ? OR remark LIKE ? OR operator_name LIKE ?", like, like, like, like)
	}
	if req.StartDate != "" {
		query = query.Where("expense_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("expense_date <= ?", req.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return records, 0, err
	}
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&records).Error; err != nil {
		return records, 0, err
	}
	return records, total, nil
}

func (m *StoreExpenseModule) Update(id, storeID uint, hqUnbound bool, updates map[string]interface{}) error {
	query := m.db.Model(&model.StoreExpense{}).Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	return query.Updates(updates).Error
}

func (m *StoreExpenseModule) Delete(id, storeID uint, hqUnbound bool) error {
	query := m.db.Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	return query.Delete(&model.StoreExpense{}).Error
}

func (m *StoreExpenseModule) Stats(req *model.ListStoreExpenseReq) (*model.StoreExpenseStats, error) {
	stats := &model.StoreExpenseStats{}
	query := m.db.Model(&model.StoreExpense{})
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.CategoryCode != "" {
		query = query.Where("category_code = ?", req.CategoryCode)
	}
	if req.StartDate != "" {
		query = query.Where("expense_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("expense_date <= ?", req.EndDate)
	}
	if err := query.Select("COALESCE(SUM(amount),0) AS total_amount, COUNT(*) AS count").Scan(stats).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (m *StoreExpenseModule) GenerateExpenseNo() string {
	now := time.Now()
	return fmt.Sprintf("ZC%s%03d", now.Format("20060102150405"), now.UnixNano()%1000)
}
