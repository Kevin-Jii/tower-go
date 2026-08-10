package module

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
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
		var deductItems []model.StoreAccountItem
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
		} else {
			// 无出库单行（例如全部为手写明细）：仅对 product_id>0 的明细扣库存；兼容 outOrder==nil 时按记账明细过滤
			for _, it := range account.Items {
				if it.ProductID == model.StoreAccountItemCustomProductID {
					continue
				}
				deductItems = append(deductItems, it)
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
				return apicode.Newf(apicode.InventoryNotFound, "商品【%s】库存不存在，无法出库", name)
			}
			if inv.Quantity < item.Quantity {
				name := item.ProductName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return apicode.Newf(apicode.InventoryInsufficient, "商品【%s】库存不足，当前库存: %.2f，需出库: %.2f", name, inv.Quantity, item.Quantity)
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
				return apicode.Newf(apicode.InventoryInsufficient, "商品【%s】库存不足，出库失败", name)
			}
		}

		return nil
	})
}

// GetByID 根据ID获取记账（含明细）
func (m *StoreAccountModule) GetByID(id uint) (*model.StoreAccount, error) {
	var account model.StoreAccount
	err := m.db.Preload("Items").Preload("Consumables").Preload("Store").Preload("Operator").Preload("Member").First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (m *StoreAccountModule) GetByIDScoped(id, storeID uint, hqUnbound bool) (*model.StoreAccount, error) {
	var account model.StoreAccount
	query := m.db.Preload("Items").Preload("Consumables").Preload("Store").Preload("Operator").Preload("Member").Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (m *StoreAccountModule) ExistsByStoreChannelDateOrderNo(storeID, excludeID uint, channel, orderNo string, accountDate time.Time) (bool, error) {
	var count int64
	query := m.db.Model(&model.StoreAccount{}).
		Where("store_id = ? AND channel = ? AND order_no = ? AND account_date = ?", storeID, channel, orderNo, accountDate.Format("2006-01-02"))
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
	if req.PaymentStatus > 0 {
		query = query.Where("payment_status = ?", req.PaymentStatus)
	}
	if kw := strings.TrimSpace(req.MemberKeyword); kw != "" {
		like := "%" + kw + "%"
		query = query.Joins("LEFT JOIN t_member AS member_search ON member_search.id = store_accounts.member_id")
		if memberID, err := strconv.ParseUint(kw, 10, 64); err == nil && memberID > 0 {
			query = query.Where("(member_search.phone LIKE ? OR member_search.name LIKE ? OR store_accounts.member_id = ?)", like, like, memberID)
		} else {
			query = query.Where("(member_search.phone LIKE ? OR member_search.name LIKE ?)", like, like)
		}
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
	if err := query.Preload("Items").Preload("Store").Preload("Operator").Preload("Member").
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

// ReplaceItemsWithInventoryAdjustments 原子替换记账明细并应用库存差量。
func (m *StoreAccountModule) ReplaceItemsWithInventoryAdjustments(
	id, storeID uint,
	hqUnbound bool,
	updates map[string]interface{},
	items []model.StoreAccountItem,
	inOrder, outOrder *model.InventoryOrder,
) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		var account model.StoreAccount
		query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id)
		if !hqUnbound {
			query = query.Where("store_id = ?", storeID)
		}
		if err := query.First(&account).Error; err != nil {
			return err
		}
		if account.IsCanceled {
			return apicode.Newf(apicode.OperationDenied, "作废记账单不允许修改")
		}
		if account.IsB2BSupplyOrderAccount() {
			return apicode.Newf(apicode.OperationDenied, "B2B供货生成的记账单仅供查看")
		}

		if outOrder != nil {
			for _, item := range outOrder.Items {
				var inv model.Inventory
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
					Where("store_id = ? AND product_id = ?", account.StoreID, item.ProductID).
					First(&inv).Error; err != nil {
					name := item.ProductName
					if name == "" {
						name = fmt.Sprintf("商品ID:%d", item.ProductID)
					}
					return apicode.Newf(apicode.InventoryNotFound, "商品【%s】库存不存在，无法补扣", name)
				}
				if inv.Quantity < item.Quantity {
					name := item.ProductName
					if name == "" {
						name = fmt.Sprintf("商品ID:%d", item.ProductID)
					}
					return apicode.Newf(apicode.InventoryInsufficient, "商品【%s】库存不足，当前库存: %.2f，需补扣: %.2f", name, inv.Quantity, item.Quantity)
				}
			}
		}

		if inOrder != nil && len(inOrder.Items) > 0 {
			if err := tx.Create(inOrder).Error; err != nil {
				return fmt.Errorf("create account edit stock return order: %w", err)
			}
			for _, item := range inOrder.Items {
				if err := incrementInventoryQuantity(tx, account.StoreID, item.ProductID, item.Quantity, item.Unit); err != nil {
					return fmt.Errorf("return account edit stock for product %d: %w", item.ProductID, err)
				}
			}
		}

		if outOrder != nil && len(outOrder.Items) > 0 {
			if err := tx.Create(outOrder).Error; err != nil {
				return fmt.Errorf("create account edit stock out order: %w", err)
			}
			for _, item := range outOrder.Items {
				res := tx.Model(&model.Inventory{}).
					Where("store_id = ? AND product_id = ? AND quantity >= ?", account.StoreID, item.ProductID, item.Quantity).
					Update("quantity", gorm.Expr("quantity - ?", item.Quantity))
				if res.Error != nil {
					return res.Error
				}
				if res.RowsAffected == 0 {
					return apicode.Newf(apicode.InventoryInsufficient, "商品【%s】库存不足，补扣失败", item.ProductName)
				}
			}
		}

		if err := tx.Where("account_id = ?", account.ID).Delete(&model.StoreAccountItem{}).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].ID = 0
			items[i].AccountID = account.ID
			items[i].DeletedAt = gorm.DeletedAt{}
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}

		return tx.Model(&model.StoreAccount{}).Where("id = ?", account.ID).Updates(updates).Error
	})
}

func (m *StoreAccountModule) CancelWithStockRestore(id, storeID uint, hqUnbound bool, operatorID uint, remark string, restoreOrder *model.InventoryOrder) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		var account model.StoreAccount
		query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id)
		if !hqUnbound {
			query = query.Where("store_id = ?", storeID)
		}
		if err := query.First(&account).Error; err != nil {
			return err
		}
		if account.IsCanceled {
			return apicode.Newf(apicode.DuplicateOperation, "记账单已作废")
		}
		if account.IsB2BSupplyOrderAccount() {
			return apicode.Newf(apicode.OperationDenied, "B2B供货生成的记账单不允许作废")
		}

		if restoreOrder != nil && len(restoreOrder.Items) > 0 {
			if err := tx.Create(restoreOrder).Error; err != nil {
				return fmt.Errorf("create store account cancel inventory order: %w", err)
			}
			for _, item := range restoreOrder.Items {
				if err := incrementInventoryQuantity(tx, account.StoreID, item.ProductID, item.Quantity, item.Unit); err != nil {
					return fmt.Errorf("restore store account inventory for product %d: %w", item.ProductID, err)
				}
			}
		}

		now := time.Now()
		res := tx.Model(&model.StoreAccount{}).
			Where("id = ? AND is_canceled = ?", account.ID, false).
			Updates(map[string]interface{}{
				"is_canceled":    true,
				"canceled_at":    &now,
				"canceled_by_id": operatorID,
				"cancel_remark":  remark,
			})
		if res.Error != nil {
			return fmt.Errorf("mark store account canceled: %w", res.Error)
		}
		if res.RowsAffected == 0 {
			return apicode.Newf(apicode.DuplicateOperation, "记账单已作废")
		}
		return nil
	})
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

// GetStatsByDateRange 按日期范围统计。
func (m *StoreAccountModule) GetStatsByDateRange(storeID uint, startDate, endDate string) (float64, float64, int64, error) {
	return m.GetStatsByDateRangeWithPaymentStatus(storeID, startDate, endDate, 0)
}

// GetStatsByDateRangeWithPaymentStatus 按日期范围和支付状态统计。
func (m *StoreAccountModule) GetStatsByDateRangeWithPaymentStatus(storeID uint, startDate, endDate string, paymentStatus int) (float64, float64, int64, error) {
	var netIncomeAmount float64
	var summary struct {
		TotalAmount float64
		Count       int64
	}

	query := m.db.Model(&model.StoreAccount{}).Where("is_canceled = ?", false)
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}
	if paymentStatus > 0 {
		query = query.Where("payment_status = ?", paymentStatus)
	}

	if err := query.Select("COALESCE(SUM(total_amount - COALESCE(round_amount, 0)), 0) AS total_amount, COUNT(*) AS count").Scan(&summary).Error; err != nil {
		return 0, 0, 0, err
	}

	// 实时净利润：销售额 - 其他支出 - 跑腿费 - 消耗品金额 - 商品成本 - 赠酒成本 - 抹零金额（不依赖历史 net_income_amount 存量值）
	costSub := m.db.Table("store_account_items AS sai").
		Select("sai.account_id, COALESCE(SUM(sai.quantity * COALESCE(ps.cost_price,0)),0) AS cost_amount").
		Joins("JOIN store_accounts AS sa_cost ON sa_cost.id = sai.account_id AND sa_cost.deleted_at IS NULL AND sa_cost.is_canceled = 0").
		Joins("LEFT JOIN product_unit_specs AS ps ON ps.product_id = sai.product_id AND ps.is_enabled = 1 AND (ps.unit_code = sai.unit OR ps.unit_name = sai.unit)").
		Group("sai.account_id")
	netQuery := m.db.Model(&model.StoreAccount{}).
		Select("COALESCE(SUM(store_accounts.total_amount - store_accounts.other_expense_amount - store_accounts.errand_fee - COALESCE(cons.sum_amount, 0) - COALESCE(costs.cost_amount,0) - store_accounts.gift_wine_cost_amount - store_accounts.round_amount), 0)").
		Joins("LEFT JOIN (SELECT sac.account_id, COALESCE(SUM(sac.amount),0) AS sum_amount FROM store_account_consumables AS sac JOIN store_accounts AS sa_cons ON sa_cons.id = sac.account_id AND sa_cons.deleted_at IS NULL AND sa_cons.is_canceled = 0 GROUP BY sac.account_id) AS cons ON cons.account_id = store_accounts.id").
		Joins("LEFT JOIN (?) AS costs ON costs.account_id = store_accounts.id", costSub).
		Where("store_accounts.is_canceled = ?", false)
	if storeID > 0 {
		netQuery = netQuery.Where("store_accounts.store_id = ?", storeID)
	}
	if startDate != "" {
		netQuery = netQuery.Where("store_accounts.account_date >= ?", startDate)
	}
	if endDate != "" {
		netQuery = netQuery.Where("store_accounts.account_date <= ?", endDate)
	}
	if paymentStatus > 0 {
		netQuery = netQuery.Where("store_accounts.payment_status = ?", paymentStatus)
	}
	if err := netQuery.Scan(&netIncomeAmount).Error; err != nil {
		return 0, 0, 0, err
	}

	return summary.TotalAmount, netIncomeAmount, summary.Count, nil
}

// GetChannelStatsByDateRange 按渠道统计销售额和订单数。
func (m *StoreAccountModule) GetChannelStatsByDateRange(storeID uint, startDate, endDate string, paymentStatus int) ([]model.ChannelStatsItem, error) {
	results := make([]model.ChannelStatsItem, 0)
	query := m.db.Model(&model.StoreAccount{}).
		Select("channel, COALESCE(SUM(total_amount), 0) AS amount, COUNT(*) AS orders").
		Where("deleted_at IS NULL AND is_canceled = 0")
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	if startDate != "" {
		query = query.Where("account_date >= ?", startDate)
	}
	if endDate != "" {
		query = query.Where("account_date <= ?", endDate)
	}
	if paymentStatus > 0 {
		query = query.Where("payment_status = ?", paymentStatus)
	}

	if err := query.Group("channel").Order("amount DESC").Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
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

		var itemCostTotal float64
		if err := tx.Table("store_account_items AS sai").
			Select("COALESCE(SUM(sai.quantity * COALESCE(ps.cost_price,0)),0)").
			Joins("LEFT JOIN product_unit_specs AS ps ON ps.product_id = sai.product_id AND ps.is_enabled = 1 AND (ps.unit_code = sai.unit OR ps.unit_name = sai.unit)").
			Where("sai.account_id = ?", accountID).
			Scan(&itemCostTotal).Error; err != nil {
			return err
		}

		netIncome := account.TotalAmount - account.OtherExpenseAmount - account.ErrandFee - consumableTotal - itemCostTotal - account.GiftWineCostAmount - account.RoundAmount
		return tx.Model(&model.StoreAccount{}).Where("id = ?", accountID).Update("net_income_amount", netIncome).Error
	})
}

func (m *StoreAccountModule) CreateConsumableProduct(product *model.StoreAccountConsumableProduct) error {
	return m.db.Create(product).Error
}

func (m *StoreAccountModule) GetConsumableProductByIDScoped(id, storeID uint, hqUnbound bool) (*model.StoreAccountConsumableProduct, error) {
	var product model.StoreAccountConsumableProduct
	query := m.db.Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (m *StoreAccountModule) GetConsumableProductMap(ids []uint, storeID uint, hqUnbound bool) (map[uint]*model.StoreAccountConsumableProduct, error) {
	result := make(map[uint]*model.StoreAccountConsumableProduct)
	if len(ids) == 0 {
		return result, nil
	}
	var products []*model.StoreAccountConsumableProduct
	query := m.db.Where("id IN ?", ids)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}
	for _, product := range products {
		result[product.ID] = product
	}
	return result, nil
}

func (m *StoreAccountModule) ListConsumableProducts(req *model.ListStoreAccountConsumableProductReq) ([]*model.StoreAccountConsumableProduct, int64, error) {
	products := make([]*model.StoreAccountConsumableProduct, 0)
	var total int64

	query := m.db.Model(&model.StoreAccountConsumableProduct{}).Preload("Store")
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if keyword := strings.TrimSpace(req.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR remark LIKE ?", like, like)
	}
	if err := query.Count(&total).Error; err != nil {
		return products, 0, err
	}
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&products).Error; err != nil {
		return products, 0, err
	}
	return products, total, nil
}

func (m *StoreAccountModule) ListAllConsumableProducts(storeID uint) ([]*model.StoreAccountConsumableProduct, error) {
	products := make([]*model.StoreAccountConsumableProduct, 0)
	query := m.db.Model(&model.StoreAccountConsumableProduct{}).Preload("Store")
	if storeID > 0 {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.Order("name ASC, id ASC").Find(&products).Error; err != nil {
		return products, err
	}
	return products, nil
}

func (m *StoreAccountModule) UpdateConsumableProduct(product *model.StoreAccountConsumableProduct) error {
	return m.db.Model(&model.StoreAccountConsumableProduct{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
		"store_id":   product.StoreID,
		"name":       product.Name,
		"cost_price": product.CostPrice,
		"remark":     product.Remark,
	}).Error
}

func (m *StoreAccountModule) DeleteConsumableProduct(id, storeID uint, hqUnbound bool) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		query := tx.Where("id = ?", id)
		if !hqUnbound {
			query = query.Where("store_id = ?", storeID)
		}
		var product model.StoreAccountConsumableProduct
		if err := query.First(&product).Error; err != nil {
			return err
		}
		if err := tx.Where("consumable_product_id = ? AND store_id = ?", product.ID, product.StoreID).
			Delete(&model.ProductUnitSpecConsumable{}).Error; err != nil {
			return err
		}
		return tx.Delete(&product).Error
	})
}
