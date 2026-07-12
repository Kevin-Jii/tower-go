package module

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/pkg/datascope"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryModule struct {
	db *gorm.DB
}

func incrementInventoryQuantity(db *gorm.DB, storeID, productID uint, quantity float64, unit string) error {
	inv := &model.Inventory{
		StoreID:   storeID,
		ProductID: productID,
		Quantity:  quantity,
		Unit:      unit,
	}
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "store_id"}, {Name: "product_id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"quantity": gorm.Expr("quantity + VALUES(quantity)"),
		}),
	}).Create(inv).Error
}

func NewInventoryModule(db *gorm.DB) *InventoryModule {
	return &InventoryModule{db: db}
}

// GetByStoreAndProduct 获取门店商品库存
func (m *InventoryModule) GetByStoreAndProduct(storeID, productID uint) (*model.Inventory, error) {
	var inv model.Inventory
	err := m.db.Where("store_id = ? AND product_id = ?", storeID, productID).First(&inv).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// AddQuantity 增加库存
func (m *InventoryModule) AddQuantity(storeID, productID uint, quantity float64, unit string) error {
	return incrementInventoryQuantity(m.db, storeID, productID, quantity, unit)
}

// SubQuantity 减少库存
func (m *InventoryModule) SubQuantity(storeID, productID uint, quantity float64) error {
	res := m.db.Model(&model.Inventory{}).
		Where("store_id = ? AND product_id = ? AND quantity >= ?", storeID, productID, quantity).
		Update("quantity", gorm.Expr("quantity - ?", quantity))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return nil
	}

	var inv model.Inventory
	if err := m.db.Where("store_id = ? AND product_id = ?", storeID, productID).First(&inv).Error; err != nil {
		return err
	}
	return apicode.Newf(apicode.InventoryInsufficient, "库存不足，当前库存: %.2f", inv.Quantity)
}

// UpdateQuantity 直接修改库存数量
func (m *InventoryModule) UpdateQuantity(id uint, quantity float64) error {
	return m.db.Model(&model.Inventory{}).Where("id = ?", id).Update("quantity", quantity).Error
}

// UpdateQuantityAndUnit 同时更新库存数量和单位
func (m *InventoryModule) UpdateQuantityAndUnit(id uint, quantity float64, unit string) error {
	return m.db.Model(&model.Inventory{}).Where("id = ?", id).Updates(map[string]interface{}{
		"quantity": quantity,
		"unit":     unit,
	}).Error
}

// GetByID 根据ID获取库存
func (m *InventoryModule) GetByID(id uint) (*model.Inventory, error) {
	var inv model.Inventory
	err := m.db.First(&inv, id).Error
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

// List 库存列表
func (m *InventoryModule) List(req *model.ListInventoryReq) ([]*model.InventoryWithProduct, int64, error) {
	var results []*model.InventoryWithProduct
	var total int64

	query := datascope.ApplyInventoriesList(m.db.Table("inventories i").
		Select("i.id, i.store_id, s.name as store_name, i.product_id, sp.name as product_name, COALESCE(sp.price, 0) as price, i.quantity, i.unit").
		Joins("LEFT JOIN stores s ON s.id = i.store_id").
		Joins("LEFT JOIN supplier_products sp ON sp.id = i.product_id").
		Where("i.deleted_at IS NULL"), req)
	if req.ProductID > 0 {
		query = query.Where("i.product_id = ?", req.ProductID)
	}
	if req.Keyword != "" {
		query = query.Where("sp.name LIKE ?", "%"+req.Keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("i.id DESC").Offset(offset).Limit(req.PageSize).Scan(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// CreateOrder 创建出入库单
func (m *InventoryModule) CreateOrder(order *model.InventoryOrder) error {
	return m.db.Create(order).Error
}

// CreateOrderWithStockApply 创建出入库单并更新库存（同事务）
func (m *InventoryModule) CreateOrderWithStockApply(order *model.InventoryOrder) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if order.Type == model.InventoryTypeOut {
			for _, item := range order.Items {
				var inv model.Inventory
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
					Where("store_id = ? AND product_id = ?", order.StoreID, item.ProductID).
					First(&inv).Error; err != nil {
					name := item.ProductName
					if name == "" {
						name = fmt.Sprintf("商品ID:%d", item.ProductID)
					}
					return apicode.Newf(apicode.InventoryNotFound, "商品【%s】不在库存中，无法出库", name)
				}
				if inv.Quantity < item.Quantity {
					name := item.ProductName
					if name == "" {
						name = fmt.Sprintf("商品ID:%d", item.ProductID)
					}
					return apicode.Newf(apicode.InventoryInsufficient, "商品【%s】库存不足，当前库存: %.2f，出库数量: %.2f", name, inv.Quantity, item.Quantity)
				}
			}
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for _, item := range order.Items {
			if order.Type == model.InventoryTypeIn {
				if err := incrementInventoryQuantity(tx, order.StoreID, item.ProductID, item.Quantity, item.Unit); err != nil {
					return err
				}
				continue
			}

			res := tx.Model(&model.Inventory{}).
				Where("store_id = ? AND product_id = ? AND quantity >= ?", order.StoreID, item.ProductID, item.Quantity).
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

// GetOrderByNo 根据单号获取出入库单
func (m *InventoryModule) GetOrderByNo(orderNo string) (*model.InventoryOrder, error) {
	var order model.InventoryOrder
	err := m.db.Preload("Items").Where("order_no = ?", orderNo).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrderByID 根据ID获取出入库单
func (m *InventoryModule) GetOrderByID(id uint) (*model.InventoryOrder, error) {
	var order model.InventoryOrder
	err := m.db.Preload("Items").First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// ListOrders 出入库单列表
func (m *InventoryModule) ListOrders(req *model.ListInventoryOrderReq) ([]*model.InventoryOrder, int64, error) {
	var orders []*model.InventoryOrder
	var total int64

	query := datascope.ApplyInventoryOrdersList(m.db.Model(&model.InventoryOrder{}), req)
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.OrderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.Date != "" {
		query = query.Where("created_at >= ? AND created_at < DATE_ADD(?, INTERVAL 1 DAY)", req.Date, req.Date)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	// 使用Preload避免N+1查询问题
	if err := query.
		Preload("Items").
		Order("id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GenerateOrderNo 生成单据编号
// 入库: RK + 日期 + 序号，如 RK202412070001
// 出库: CK + 日期 + 序号，如 CK202412070001
func (m *InventoryModule) GenerateOrderNo(orderType int8) string {
	prefix := "RK" // 入库
	if orderType == model.InventoryTypeOut {
		prefix = "CK" // 出库
	}

	today := time.Now().Format("20060102")
	pattern := prefix + today + "%"

	var maxNo string
	m.db.Model(&model.InventoryOrder{}).
		Where("order_no LIKE ?", pattern).
		Order("order_no DESC").
		Limit(1).
		Pluck("order_no", &maxNo)

	seq := 1
	if maxNo != "" && len(maxNo) >= 14 {
		fmt.Sscanf(maxNo[len(maxNo)-4:], "%d", &seq)
		seq++
	}

	return fmt.Sprintf("%s%s%04d", prefix, today, seq)
}
