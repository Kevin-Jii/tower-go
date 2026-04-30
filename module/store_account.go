package module

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/datascope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// CreateWithInventoryOut 创建记账并自动出库（同事务）
func (m *StoreAccountModule) CreateWithInventoryOut(account *model.StoreAccount, outOrder *model.InventoryOrder) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		deductItems := account.Items
		if outOrder != nil && len(outOrder.Items) > 0 {
			deductItems = make([]model.StoreAccountItem, 0, len(outOrder.Items))
			for _, item := range outOrder.Items {
				deductItems = append(deductItems, model.StoreAccountItem{
					ProductID:   item.ProductID,
					ProductName: item.ProductName,
					Quantity:    item.Quantity,
					Unit:        item.Unit,
				})
			}
		}

		// 先锁库存并做充足性校验，避免并发下出现负库存
		for _, item := range deductItems {
			var inv model.Inventory
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("store_id = ? AND product_id = ?", account.StoreID, item.ProductID).
				First(&inv).Error; err != nil {
				name := item.ProductName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return fmt.Errorf("商品【%s】库存不存在，无法出库", name)
			}
			if inv.Quantity < item.Quantity {
				name := item.ProductName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return fmt.Errorf("商品【%s】库存不足，当前库存: %.2f，需出库: %.2f", name, inv.Quantity, item.Quantity)
			}
		}

		if err := tx.Create(account).Error; err != nil {
			return err
		}

		if outOrder != nil {
			if err := tx.Create(outOrder).Error; err != nil {
				return err
			}
		}

		for _, item := range deductItems {
			res := tx.Model(&model.Inventory{}).
				Where("store_id = ? AND product_id = ? AND quantity >= ?", account.StoreID, item.ProductID, item.Quantity).
				Update("quantity", gorm.Expr("quantity - ?", item.Quantity))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				name := item.ProductName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return fmt.Errorf("商品【%s】库存不足，出库失败", name)
			}
		}

		return nil
	})
}

// GetByID 根据ID获取记账（含明细）
func (m *StoreAccountModule) GetByID(id uint) (*model.StoreAccount, error) {
	var account model.StoreAccount
	err := m.db.Preload("Items").Preload("Consumables").Preload("Store").Preload("Operator").First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// List 记账列表
func (m *StoreAccountModule) List(req *model.ListStoreAccountReq) ([]*model.StoreAccount, int64, error) {
	accounts := make([]*model.StoreAccount, 0) // 初始化为空数组，避免返回null
	var total int64

	query := datascope.ApplyStoreAccountsList(m.db.Model(&model.StoreAccount{}), req)
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
		return accounts, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Items").Preload("Store").Preload("Operator").
		Preload("Consumables").
		Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&accounts).Error; err != nil {
		return accounts, 0, err
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
func (m *StoreAccountModule) GetStatsByDateRange(storeID uint, startDate, endDate string) (float64, float64, int64, error) {
	var totalAmount float64
	var netIncomeAmount float64
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
		return 0, 0, 0, err
	}

	query.Select("COALESCE(SUM(total_amount), 0)").Scan(&totalAmount)
	query.Select("COALESCE(SUM(net_income_amount), 0)").Scan(&netIncomeAmount)

	return totalAmount, netIncomeAmount, count, nil
}

func (m *StoreAccountModule) ReplaceConsumables(accountID uint, consumables []model.StoreAccountConsumable) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id = ?", accountID).Delete(&model.StoreAccountConsumable{}).Error; err != nil {
			return err
		}
		if len(consumables) > 0 {
			if err := tx.Create(&consumables).Error; err != nil {
				return err
			}
		}

		var account model.StoreAccount
		if err := tx.First(&account, accountID).Error; err != nil {
			return err
		}

		var consumableTotal float64
		if err := tx.Model(&model.StoreAccountConsumable{}).
			Where("account_id = ?", accountID).
			Select("COALESCE(SUM(amount),0)").
			Scan(&consumableTotal).Error; err != nil {
			return err
		}

		netIncome := account.TotalAmount - account.OtherExpenseAmount - consumableTotal
		return tx.Model(&model.StoreAccount{}).Where("id = ?", accountID).Update("net_income_amount", netIncome).Error
	})
}
