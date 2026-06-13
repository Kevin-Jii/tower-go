package module

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InventoryLossModule struct {
	db *gorm.DB
}

func NewInventoryLossModule(db *gorm.DB) *InventoryLossModule {
	return &InventoryLossModule{db: db}
}

func (m *InventoryLossModule) CreateWithStockDeduct(order *model.InventoryLossOrder) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range order.Items {
			var inv model.Inventory
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("store_id = ? AND product_id = ?", order.StoreID, item.ProductID).
				First(&inv).Error; err != nil {
				name := item.ProductName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return fmt.Errorf("商品【%s】不在库存中，无法扣减", name)
			}
			if inv.Quantity < item.BaseQuantity {
				name := item.ProductName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return fmt.Errorf("商品【%s】库存不足，当前库存: %.2f%s，需要扣减: %.2f%s", name, inv.Quantity, inv.Unit, item.BaseQuantity, item.BaseUnit)
			}
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}

		for _, item := range order.Items {
			res := tx.Model(&model.Inventory{}).
				Where("store_id = ? AND product_id = ? AND quantity >= ?", order.StoreID, item.ProductID, item.BaseQuantity).
				Update("quantity", gorm.Expr("quantity - ?", item.BaseQuantity))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				name := item.ProductName
				if name == "" {
					name = fmt.Sprintf("商品ID:%d", item.ProductID)
				}
				return fmt.Errorf("商品【%s】库存不足，扣减失败", name)
			}
		}

		return nil
	})
}

func (m *InventoryLossModule) GetByIDScoped(id, storeID uint, hqUnbound bool) (*model.InventoryLossOrder, error) {
	var order model.InventoryLossOrder
	query := m.db.Preload("Items").Preload("Member").Where("id = ?", id)
	if !hqUnbound {
		query = query.Where("store_id = ?", storeID)
	}
	if err := query.First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (m *InventoryLossModule) CancelWithStockRestore(id, storeID uint, hqUnbound bool) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		var order model.InventoryLossOrder
		query := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Items").Where("id = ?", id)
		if !hqUnbound {
			query = query.Where("store_id = ?", storeID)
		}
		if err := query.First(&order).Error; err != nil {
			return err
		}
		if order.IsCanceled {
			return fmt.Errorf("单据已撤销")
		}

		for _, item := range order.Items {
			var inv model.Inventory
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("store_id = ? AND product_id = ?", order.StoreID, item.ProductID).
				First(&inv).Error
			if err == gorm.ErrRecordNotFound {
				inv = model.Inventory{
					StoreID:   order.StoreID,
					ProductID: item.ProductID,
					Quantity:  item.BaseQuantity,
					Unit:      item.BaseUnit,
				}
				if err := tx.Create(&inv).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}
			if err := tx.Model(&inv).Update("quantity", gorm.Expr("quantity + ?", item.BaseQuantity)).Error; err != nil {
				return err
			}
		}

		now := time.Now()
		return tx.Model(&order).Updates(map[string]interface{}{
			"is_canceled": true,
			"canceled_at": &now,
		}).Error
	})
}

func (m *InventoryLossModule) List(req *model.ListInventoryLossOrderReq) ([]*model.InventoryLossOrder, int64, error) {
	var orders []*model.InventoryLossOrder
	var total int64

	query := m.db.Model(&model.InventoryLossOrder{}).Preload("Member")
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.MemberID > 0 {
		query = query.Where("member_id = ?", req.MemberID)
	}
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("order_no LIKE ? OR reason LIKE ? OR operator_name LIKE ?", keyword, keyword, keyword)
	}
	if req.StartDate != "" {
		query = query.Where("created_at >= ?", req.StartDate+" 00:00:00")
	}
	if req.EndDate != "" {
		query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func (m *InventoryLossModule) ListMemberGiftRecords(memberID uint, req *model.ListMemberGiftRecordsReq, hqUnbound bool) ([]*model.MemberGiftRecord, int64, error) {
	var records []*model.MemberGiftRecord
	var total int64

	query := m.db.Table("inventory_loss_order_items i").
		Joins("JOIN inventory_loss_orders o ON o.id = i.order_id AND o.deleted_at IS NULL").
		Where("i.deleted_at IS NULL AND o.type = ? AND o.member_id = ? AND o.is_canceled = 0", model.InventoryLossTypeGift, memberID)
	if !hqUnbound || req.StoreID > 0 {
		query = query.Where("o.store_id = ?", req.StoreID)
	}
	if req.StartDate != "" {
		query = query.Where("o.created_at >= ?", req.StartDate+" 00:00:00")
	}
	if req.EndDate != "" {
		query = query.Where("o.created_at <= ?", req.EndDate+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Select(`
			i.id,
			i.order_id,
			o.order_no,
			i.product_id,
			i.product_name,
			i.unit,
			i.quantity,
			i.cost_amount,
			o.reason,
			o.operator_name,
			o.created_at
		`).
		Order("o.id DESC, i.id DESC").
		Offset(offset).
		Limit(req.PageSize).
		Scan(&records).Error; err != nil {
		return nil, 0, err
	}
	return records, total, nil
}

func (m *InventoryLossModule) GenerateOrderNo() string {
	prefix := "BS"
	today := time.Now().Format("20060102")
	pattern := prefix + today + "%"

	var maxNo string
	m.db.Model(&model.InventoryLossOrder{}).
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
