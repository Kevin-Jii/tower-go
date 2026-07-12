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

func (m *StoreSupplierModule) activeSupplierIDs(storeID uint) *gorm.DB {
	return m.db.Model(&model.StoreSupplier{}).
		Select("supplier_id").
		Where("store_id = ? AND status = 1", storeID)
}

// BindSuppliers 门店绑定供应商
func (m *StoreSupplierModule) BindSuppliers(storeID uint, supplierIDs []uint) error {
	if len(supplierIDs) == 0 {
		return nil
	}

	uniqueSupplierIDs := make([]uint, 0, len(supplierIDs))
	seen := make(map[uint]struct{}, len(supplierIDs))
	for _, supplierID := range supplierIDs {
		if _, ok := seen[supplierID]; ok {
			continue
		}
		seen[supplierID] = struct{}{}
		uniqueSupplierIDs = append(uniqueSupplierIDs, supplierID)
	}

	return m.db.Transaction(func(tx *gorm.DB) error {
		var bindings []model.StoreSupplier
		if err := tx.
			Where("store_id = ? AND supplier_id IN ?", storeID, uniqueSupplierIDs).
			Find(&bindings).Error; err != nil {
			return err
		}

		boundIDs := make(map[uint]struct{}, len(bindings))
		for _, binding := range bindings {
			boundIDs[binding.SupplierID] = struct{}{}
		}

		if len(bindings) > 0 {
			if err := tx.Model(&model.StoreSupplier{}).
				Where("store_id = ? AND supplier_id IN ?", storeID, uniqueSupplierIDs).
				Updates(map[string]interface{}{
					"status": 1,
				}).Error; err != nil {
				return err
			}
		}

		newBindings := make([]model.StoreSupplier, 0, len(uniqueSupplierIDs)-len(boundIDs))
		for _, supplierID := range uniqueSupplierIDs {
			if _, ok := boundIDs[supplierID]; ok {
				continue
			}
			newBindings = append(newBindings, model.StoreSupplier{
				StoreID:    storeID,
				SupplierID: supplierID,
				Status:     1,
			})
		}
		if len(newBindings) > 0 {
			return tx.Create(&newBindings).Error
		}
		return nil
	})
}

// UnbindSuppliers 门店解绑供应商
func (m *StoreSupplierModule) UnbindSuppliers(storeID uint, supplierIDs []uint) error {
	if len(supplierIDs) == 0 {
		return nil
	}
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

// IsSupplierBound 判断供应商是否已绑定到门店。
func (m *StoreSupplierModule) IsSupplierBound(storeID, supplierID uint) (bool, error) {
	var count int64
	if err := m.db.Model(&model.StoreSupplier{}).
		Where("store_id = ? AND supplier_id = ? AND status = ?", storeID, supplierID, 1).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// ListProductsByStoreID 获取门店可采购的商品列表（绑定供应商的所有商品）
func (m *StoreSupplierModule) ListProductsByStoreID(storeID, supplierID, categoryID uint, keyword string) ([]*model.SupplierProduct, error) {
	query := m.db.Preload("Supplier").Preload("Category").
		Where("supplier_id IN (?) AND status = 1", m.activeSupplierIDs(storeID))

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

	products := make([]*model.SupplierProduct, 0)
	if err := query.Order("supplier_id, category_id, name").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// ListCategoriesByStoreID 获取门店绑定供应商下的所有分类（可按供应商筛选）
func (m *StoreSupplierModule) ListCategoriesByStoreID(storeID, supplierID uint) ([]*model.SupplierCategory, error) {
	query := m.db.Preload("Supplier").
		Where("supplier_id IN (?) AND status = 1", m.activeSupplierIDs(storeID))

	if supplierID > 0 {
		query = query.Where("supplier_id = ?", supplierID)
	}

	categories := make([]*model.SupplierCategory, 0)
	if err := query.Order("supplier_id, sort ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// ValidateStoreProducts 验证商品是否属于门店绑定的供应商
// 返回不可用的商品ID列表
func (m *StoreSupplierModule) ValidateStoreProducts(storeID uint, productIDs []uint) ([]uint, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}

	// 查询这些商品中，属于绑定供应商的商品ID
	var validProductIDs []uint
	if err := m.db.Model(&model.SupplierProduct{}).
		Where("id IN ? AND supplier_id IN (?) AND status = 1", productIDs, m.activeSupplierIDs(storeID)).
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
