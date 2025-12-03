package module

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/tenant"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
	"gorm.io/gorm"
)

type PurchaseOrderModule struct {
	db       *gorm.DB
	strategy tenant.IsolationStrategy
}

func NewPurchaseOrderModule(db *gorm.DB) *PurchaseOrderModule {
	return &PurchaseOrderModule{
		db:       db,
		strategy: tenant.NewStoreIsolationStrategy(),
	}
}

// withTenant 应用租户隔离
func (m *PurchaseOrderModule) withTenant(storeID uint) func(*gorm.DB) *gorm.DB {
	if storeID == 0 {
		return tenant.AdminScope()
	}
	return tenant.StoreScope(storeID)
}

func (m *PurchaseOrderModule) GetDB() *gorm.DB {
	return m.db
}

// GenerateOrderNo 生成订单编号
func (m *PurchaseOrderModule) GenerateOrderNo() string {
	return fmt.Sprintf("PO%s%04d", time.Now().Format("20060102150405"), time.Now().Nanosecond()%10000)
}

func (m *PurchaseOrderModule) Create(order *model.PurchaseOrder) error {
	return m.db.Create(order).Error
}

func (m *PurchaseOrderModule) GetByID(id uint) (*model.PurchaseOrder, error) {
	var order model.PurchaseOrder
	if err := m.db.Preload("Store").Preload("Creator").Preload("Items.Supplier").Preload("Items.Product").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (m *PurchaseOrderModule) List(req *model.ListPurchaseOrderReq) ([]*model.PurchaseOrder, int64, error) {
	var orders []*model.PurchaseOrder
	var total int64

	// 使用租户隔离策略
	query := m.db.Model(&model.PurchaseOrder{}).Scopes(m.withTenant(req.StoreID))

	// 其他过滤条件
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StartDate != "" {
		query = query.Where("order_date >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("order_date <= ?", req.EndDate)
	}

	// 如果指定了供应商，需要关联查询
	if req.SupplierID > 0 {
		query = query.Joins("JOIN purchase_order_items ON purchase_order_items.order_id = purchase_orders.id").
			Where("purchase_order_items.supplier_id = ?", req.SupplierID).
			Distinct()
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Store").Preload("Creator").Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (m *PurchaseOrderModule) UpdateByID(id uint, req *model.UpdatePurchaseOrderReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}
	return m.db.Model(&model.PurchaseOrder{}).Where("id = ?", id).Updates(updateMap).Error
}

func (m *PurchaseOrderModule) Delete(id uint) error {
	// 先删除明细
	m.db.Where("order_id = ?", id).Delete(&model.PurchaseOrderItem{})
	// 再删除主单
	return m.db.Delete(&model.PurchaseOrder{}, id).Error
}

// CreateItems 批量创建采购单明细
func (m *PurchaseOrderModule) CreateItems(items []model.PurchaseOrderItem) error {
	return m.db.Create(&items).Error
}

// GetItemsByOrderID 获取采购单明细
func (m *PurchaseOrderModule) GetItemsByOrderID(orderID uint) ([]model.PurchaseOrderItem, error) {
	var items []model.PurchaseOrderItem
	if err := m.db.Preload("Supplier").Preload("Product").Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetOrdersBySupplier 按供应商分组获取采购单明细
func (m *PurchaseOrderModule) GetOrdersBySupplier(orderID uint) (map[uint][]model.PurchaseOrderItem, error) {
	items, err := m.GetItemsByOrderID(orderID)
	if err != nil {
		return nil, err
	}

	result := make(map[uint][]model.PurchaseOrderItem)
	for _, item := range items {
		result[item.SupplierID] = append(result[item.SupplierID], item)
	}
	return result, nil
}
