package module

import (
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type B2BModule struct {
	db *gorm.DB
}

func NewB2BModule(db *gorm.DB) *B2BModule {
	return &B2BModule{db: db}
}

func (m *B2BModule) CreateCustomer(customer *model.B2BCustomer) error {
	return m.db.Create(customer).Error
}

func (m *B2BModule) UpdateCustomer(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.B2BCustomer{}).Where("id = ?", id).Updates(updates).Error
}

func (m *B2BModule) GetCustomer(id uint) (*model.B2BCustomer, error) {
	var customer model.B2BCustomer
	if err := m.db.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (m *B2BModule) ListCustomers(req *model.ListB2BCustomerReq) ([]*model.B2BCustomer, int64, error) {
	var rows []*model.B2BCustomer
	var total int64
	q := m.db.Model(&model.B2BCustomer{})
	if req.StoreID > 0 {
		q = q.Where("store_id = ?", req.StoreID)
	}
	if req.Status > 0 {
		q = q.Where("status = ?", req.Status)
	}
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("(name LIKE ? OR phone LIKE ? OR contact_person LIKE ?)", like, like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (req.Page - 1) * req.PageSize
	if err := q.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (m *B2BModule) UpsertPrice(price *model.B2BCustomerProductPrice) error {
	q := m.db.Model(&model.B2BCustomerProductPrice{}).
		Where("store_id = ? AND product_id = ? AND unit_spec_id = ?", price.StoreID, price.ProductID, price.UnitSpecID)
	if price.CustomerID != nil && *price.CustomerID > 0 {
		q = q.Where("customer_id = ?", *price.CustomerID)
	} else {
		q = q.Where("customer_id IS NULL AND price_level = ?", price.PriceLevel)
	}

	var existing model.B2BCustomerProductPrice
	if err := q.First(&existing).Error; err == nil {
		price.ID = existing.ID
		return m.db.Model(&existing).Updates(map[string]interface{}{
			"customer_id":  price.CustomerID,
			"price_level":  price.PriceLevel,
			"unit_name":    price.UnitName,
			"supply_price": price.SupplyPrice,
			"min_quantity": price.MinQuantity,
			"is_enabled":   price.IsEnabled,
			"remark":       price.Remark,
		}).Error
	} else if err != gorm.ErrRecordNotFound {
		return err
	}
	return m.db.Create(price).Error
}

func (m *B2BModule) UpdatePrice(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.B2BCustomerProductPrice{}).Where("id = ?", id).Updates(updates).Error
}

func (m *B2BModule) DeletePrice(id uint) error {
	return m.db.Delete(&model.B2BCustomerProductPrice{}, id).Error
}

func (m *B2BModule) ListPrices(req *model.ListB2BPriceReq) ([]*model.B2BCustomerProductPrice, int64, error) {
	var rows []*model.B2BCustomerProductPrice
	var total int64
	q := m.db.Model(&model.B2BCustomerProductPrice{}).
		Preload("Customer").
		Preload("Product").
		Preload("UnitSpec")
	if req.StoreID > 0 {
		q = q.Where("store_id = ?", req.StoreID)
	}
	if req.CustomerID > 0 {
		q = q.Where("customer_id = ?", req.CustomerID)
	}
	if strings.TrimSpace(req.PriceLevel) != "" {
		q = q.Where("price_level = ?", strings.TrimSpace(req.PriceLevel))
	}
	if req.ProductID > 0 {
		q = q.Where("product_id = ?", req.ProductID)
	}
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Joins("LEFT JOIN supplier_products sp ON sp.id = b2b_customer_product_prices.product_id").
			Joins("LEFT JOIN b2b_customers bc ON bc.id = b2b_customer_product_prices.customer_id").
			Where("(sp.name LIKE ? OR bc.name LIKE ? OR b2b_customer_product_prices.unit_name LIKE ?)", like, like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (req.Page - 1) * req.PageSize
	if err := q.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (m *B2BModule) ResolvePrice(storeID, customerID, productID, unitSpecID uint, priceLevel string) (*model.B2BCustomerProductPrice, error) {
	var price model.B2BCustomerProductPrice
	if customerID > 0 {
		err := m.db.Where("store_id = ? AND customer_id = ? AND product_id = ? AND unit_spec_id = ? AND is_enabled = 1",
			storeID, customerID, productID, unitSpecID).First(&price).Error
		if err == nil {
			return &price, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	if strings.TrimSpace(priceLevel) != "" {
		err := m.db.Where("store_id = ? AND customer_id IS NULL AND price_level = ? AND product_id = ? AND unit_spec_id = ? AND is_enabled = 1",
			storeID, strings.TrimSpace(priceLevel), productID, unitSpecID).First(&price).Error
		if err == nil {
			return &price, nil
		}
		if err != gorm.ErrRecordNotFound {
			return nil, err
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *B2BModule) CreateSupplyOrderWithInventory(order *model.B2BSupplyOrder) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range order.Items {
			var inv model.Inventory
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("store_id = ? AND product_id = ?", order.StoreID, item.ProductID).
				First(&inv).Error; err != nil {
				return fmt.Errorf("商品【%s】库存不存在，无法供货", item.ProductName)
			}
			if inv.Quantity < item.BaseQuantity {
				return fmt.Errorf("商品【%s】库存不足，当前库存: %.2f，需出库: %.2f", item.ProductName, inv.Quantity, item.BaseQuantity)
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
				return fmt.Errorf("商品【%s】库存不足，出库失败", item.ProductName)
			}
		}

		if order.UnpaidAmount != 0 {
			if err := tx.Model(&model.B2BCustomer{}).
				Where("id = ?", order.CustomerID).
				Update("receivable", gorm.Expr("receivable + ?", order.UnpaidAmount)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *B2BModule) ListSupplyOrders(req *model.ListB2BSupplyOrderReq) ([]*model.B2BSupplyOrder, int64, error) {
	var rows []*model.B2BSupplyOrder
	var total int64
	q := m.db.Model(&model.B2BSupplyOrder{}).Preload("Customer").Preload("Items")
	if req.StoreID > 0 {
		q = q.Where("store_id = ?", req.StoreID)
	}
	if req.CustomerID > 0 {
		q = q.Where("customer_id = ?", req.CustomerID)
	}
	if req.PaymentStatus > 0 {
		q = q.Where("payment_status = ?", req.PaymentStatus)
	}
	if req.DeliveryStatus > 0 {
		q = q.Where("delivery_status = ?", req.DeliveryStatus)
	}
	if req.StartDate != "" {
		q = q.Where("order_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		q = q.Where("order_date <= ?", req.EndDate)
	}
	if kw := strings.TrimSpace(req.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("(order_no LIKE ? OR customer_name LIKE ?)", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (req.Page - 1) * req.PageSize
	if err := q.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (m *B2BModule) GetSupplyOrder(id uint) (*model.B2BSupplyOrder, error) {
	var order model.B2BSupplyOrder
	if err := m.db.Preload("Customer").Preload("Items").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (m *B2BModule) GenerateSupplyOrderNo() string {
	today := time.Now().Format("20060102")
	pattern := "B2B" + today + "%"
	var maxNo string
	m.db.Model(&model.B2BSupplyOrder{}).
		Where("order_no LIKE ?", pattern).
		Order("order_no DESC").
		Limit(1).
		Pluck("order_no", &maxNo)

	seq := 1
	if maxNo != "" && len(maxNo) >= 15 {
		fmt.Sscanf(maxNo[len(maxNo)-4:], "%d", &seq)
		seq++
	}
	return fmt.Sprintf("B2B%s%04d", today, seq)
}
