package module

import (
	"fmt"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type InventoryModule struct {
	db *gorm.DB
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
	var inv model.Inventory
	err := m.db.Where("store_id = ? AND product_id = ?", storeID, productID).First(&inv).Error

	if err == gorm.ErrRecordNotFound {
		inv = model.Inventory{
			StoreID:   storeID,
			ProductID: productID,
			Quantity:  quantity,
			Unit:      unit,
		}
		return m.db.Create(&inv).Error
	}

	if err != nil {
		return err
	}

	return m.db.Model(&inv).Update("quantity", gorm.Expr("quantity + ?", quantity)).Error
}

// SubQuantity 减少库存
func (m *InventoryModule) SubQuantity(storeID, productID uint, quantity float64) error {
	var inv model.Inventory
	err := m.db.Where("store_id = ? AND product_id = ?", storeID, productID).First(&inv).Error
	if err != nil {
		return err
	}

	if inv.Quantity < quantity {
		return fmt.Errorf("库存不足，当前库存: %.2f", inv.Quantity)
	}

	return m.db.Model(&inv).Update("quantity", gorm.Expr("quantity - ?", quantity)).Error
}

// List 库存列表
func (m *InventoryModule) List(req *model.ListInventoryReq) ([]*model.InventoryWithProduct, int64, error) {
	var results []*model.InventoryWithProduct
	var total int64

	query := m.db.Table("inventories i").
		Select("i.id, i.store_id, s.name as store_name, i.product_id, sp.name as product_name, i.quantity, i.unit").
		Joins("LEFT JOIN stores s ON s.id = i.store_id").
		Joins("LEFT JOIN supplier_products sp ON sp.id = i.product_id").
		Where("i.deleted_at IS NULL")

	if req.StoreID > 0 {
		query = query.Where("i.store_id = ?", req.StoreID)
	}
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

	query := m.db.Model(&model.InventoryOrder{})

	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.OrderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}
	if req.Date != "" {
		query = query.Where("DATE(created_at) = ?", req.Date)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Items").Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&orders).Error; err != nil {
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
