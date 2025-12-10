package module

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type StoreAccountModule struct {
	db *gorm.DB
}

func NewStoreAccountModule(db *gorm.DB) *StoreAccountModule {
	return &StoreAccountModule{db: db}
}

// Create 创建记账（含明细）
func (m *StoreAccountModule) Create(account *model.StoreAccount) error {
	return m.db.Create(account).Error
}

// GetByID 根据ID获取记账（含明细）
func (m *StoreAccountModule) GetByID(id uint) (*model.StoreAccount, error) {
	var account model.StoreAccount
	err := m.db.Preload("Items").Preload("Store").Preload("Operator").First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// List 记账列表
func (m *StoreAccountModule) List(req *model.ListStoreAccountReq) ([]*model.StoreAccount, int64, error) {
	var accounts []*model.StoreAccount
	var total int64

	query := m.db.Model(&model.StoreAccount{})

	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.Channel != "" {
		query = query.Where("channel = ?", req.Channel)
	}
	if req.OrderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.TagCode != "" {
		query = query.Where("tag_code = ?", req.TagCode)
	}
	if req.StartDate != "" {
		query = query.Where("account_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("account_date <= ?", req.EndDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Items").Preload("Store").Preload("Operator").
		Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&accounts).Error; err != nil {
		return nil, 0, err
	}

	return accounts, total, nil
}

// Update 更新记账
func (m *StoreAccountModule) Update(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.StoreAccount{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除记账（含明细）
func (m *StoreAccountModule) Delete(id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 先删除明细
		if err := tx.Where("account_id = ?", id).Delete(&model.StoreAccountItem{}).Error; err != nil {
			return err
		}
		// 再删除主表
		return tx.Delete(&model.StoreAccount{}, id).Error
	})
}

// GenerateAccountNo 生成记账编号
func (m *StoreAccountModule) GenerateAccountNo() string {
	now := time.Now()
	random := now.UnixNano() % 1000
	return fmt.Sprintf("JZ%s%03d", now.Format("20060102150405"), random)
}

// GetStatsByDateRange 按日期范围统计
func (m *StoreAccountModule) GetStatsByDateRange(storeID uint, startDate, endDate string) (float64, int64, error) {
	var totalAmount float64
	var count int64

	query := m.db.Model(&model.StoreAccount{})
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, 0, err
	}

	query.Select("COALESCE(SUM(total_amount), 0)").Scan(&totalAmount)

	return totalAmount, count, nil
}
