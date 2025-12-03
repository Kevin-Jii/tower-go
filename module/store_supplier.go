package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type StoreSupplierModule struct {
	db *gorm.DB
}

func NewStoreSupplierModule(db *gorm.DB) *StoreSupplierModule {
	return &StoreSupplierModule{db: db}
}

// BindSuppliers 门店绑定供应商
func (m *StoreSupplierModule) BindSuppliers(storeID uint, supplierIDs []uint) error {
	for _, supplierID := range supplierIDs {
		// 检查是否已绑定
		var count int64
		m.db.Model(&model.StoreSupplier{}).Where("store_id = ? AND supplier_id = ?", storeID, supplierID).Count(&count)
		if count > 0 {
			continue
		}

		// 创建绑定
		binding := &model.StoreSupplier{
			StoreID:    storeID,
			SupplierID: supplierID,
			Status:     1,
		}
		if err := m.db.Create(binding).Error; err != nil {
			return err
		}
	}
	return nil
}

// UnbindSuppliers 门店解绑供应商
func (m *StoreSupplierModule) UnbindSuppliers(storeID uint, supplierIDs []uint) error {
	return m.db.Where("store_id = ? AND supplier_id IN ?", storeID, supplierIDs).Delete(&model.StoreSupplier{}).Error
}

// ListSuppliersByStoreID 获取门店绑定的所有供应商
func (m *StoreSupplierModule) ListSuppliersByStoreID(storeID uint) ([]*model.StoreSupplier, error) {
	var bindings []*model.StoreSupplier
	if err := m.db.Preload("Supplier").Where("store_id = ? AND status = 1", storeID).Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

// ListProductsByStoreID 获取门店可采购的商品列表（绑定供应商的所有商品）
func (m *StoreSupplierModule) ListProductsByStoreID(storeID, supplierID, categoryID uint, keyword string) ([]*model.SupplierProduct, error) {
	// 先获取门店绑定的供应商ID列表
	var supplierIDs []uint
	if err := m.db.Model(&model.StoreSupplier{}).
		Where("store_id = ? AND status = 1", storeID).
		Pluck("supplier_id", &supplierIDs).Error; err != nil {
		return nil, err
	}

	if len(supplierIDs) == 0 {
		return []*model.SupplierProduct{}, nil
	}

	// 查询这些供应商的商品
	query := m.db.Preload("Supplier").Preload("Category").
		Where("supplier_id IN ? AND status = 1", supplierIDs)

	// 可选筛选条件
	if supplierID > 0 {
		query = query.Where("supplier_id = ?", supplierID)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	var products []*model.SupplierProduct
	if err := query.Order("supplier_id, category_id, name").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// ValidateStoreProducts 验证商品是否属于门店绑定的供应商
// 返回不可用的商品ID列表
func (m *StoreSupplierModule) ValidateStoreProducts(storeID uint, productIDs []uint) ([]uint, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}

	// 获取门店绑定的供应商ID列表
	var supplierIDs []uint
	if err := m.db.Model(&model.StoreSupplier{}).
		Where("store_id = ? AND status = 1", storeID).
		Pluck("supplier_id", &supplierIDs).Error; err != nil {
		return nil, err
	}

	if len(supplierIDs) == 0 {
		return productIDs, nil // 没有绑定供应商，所有商品都不可用
	}

	// 查询这些商品中，属于绑定供应商的商品ID
	var validProductIDs []uint
	if err := m.db.Model(&model.SupplierProduct{}).
		Where("id IN ? AND supplier_id IN ? AND status = 1", productIDs, supplierIDs).
		Pluck("id", &validProductIDs).Error; err != nil {
		return nil, err
	}

	// 构建有效商品ID的map
	validMap := make(map[uint]bool)
	for _, id := range validProductIDs {
		validMap[id] = true
	}

	// 找出无效的商品ID
	var invalidIDs []uint
	for _, id := range productIDs {
		if !validMap[id] {
			invalidIDs = append(invalidIDs, id)
		}
	}

	return invalidIDs, nil
}
