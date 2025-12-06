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

// CreateOrUpdate 创建或更新库存
func (m *InventoryModule) CreateOrUpdate(storeID, productID uint, quantity float64, unit string) error {
	var inv model.Inventory
	err := m.db.Where("store_id = ? AND product_id = ?", storeID, productID).First(&inv).Error
	
	if err == gorm.ErrRecordNotFound {
		// 创建新库存
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
	
	// 更新库存数量
	return m.db.Model(&inv).Update("quantity", quantity).Error
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

// CreateRecord 创建出入库记录
func (m *InventoryModule) CreateRecord(record *model.InventoryRecord) error {
	return m.db.Create(record).Error
}

// ListRecords 出入库记录列表
func (m *InventoryModule) ListRecords(req *model.ListInventoryRecordReq) ([]*model.InventoryRecord, int64, error) {
	var records []*model.InventoryRecord
	var total int64

	query := m.db.Model(&model.InventoryRecord{})

	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.ProductID > 0 {
		query = query.Where("product_id = ?", req.ProductID)
	}
	if req.Type != nil {
		query = query.Where("type = ?", *req.Type)
	}
	if req.Date != "" {
		query = query.Where("DATE(created_at) = ?", req.Date)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Preload("Store").Preload("Product").Preload("Operator").
		Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GenerateRecordNo 生成单据编号
// 入库: RK + 日期 + 序号，如 RK202412070001
// 出库: CK + 日期 + 序号，如 CK202412070001
func (m *InventoryModule) GenerateRecordNo(recordType int8) string {
	prefix := "RK" // 入库
	if recordType == model.InventoryTypeOut {
		prefix = "CK" // 出库
	}

	today := time.Now().Format("20060102")
	pattern := prefix + today + "%"

	// 查询今天该类型的最大编号
	var maxNo string
	m.db.Model(&model.InventoryRecord{}).
		Where("record_no LIKE ?", pattern).
		Order("record_no DESC").
		Limit(1).
		Pluck("record_no", &maxNo)

	seq := 1
	if maxNo != "" && len(maxNo) >= 14 {
		// 提取序号部分（最后4位）
		fmt.Sscanf(maxNo[len(maxNo)-4:], "%d", &seq)
		seq++
	}

	return fmt.Sprintf("%s%s%04d", prefix, today, seq)
}
