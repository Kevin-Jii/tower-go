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

// BindProducts 门店绑定供应商商品
// 首次绑定某商品名称时自动设为默认供应商
func (m *StoreSupplierModule) BindProducts(storeID uint, productIDs []uint) error {
	for _, productID := range productIDs {
		// 检查是否已绑定
		var count int64
		m.db.Model(&model.StoreSupplierProduct{}).Where("store_id = ? AND product_id = ?", storeID, productID).Count(&count)
		if count > 0 {
			continue
		}

		// 获取商品信息，用于检查同名商品
		var product model.SupplierProduct
		if err := m.db.First(&product, productID).Error; err != nil {
			return err
		}

		// 检查门店是否已有同名商品的绑定
		isDefault := false
		var existingCount int64
		m.db.Model(&model.StoreSupplierProduct{}).
			Joins("JOIN supplier_products ON supplier_products.id = store_supplier_products.product_id").
			Where("store_supplier_products.store_id = ? AND supplier_products.name = ?", storeID, product.Name).
			Count(&existingCount)

		// 如果是首次绑定该商品名称，设为默认
		if existingCount == 0 {
			isDefault = true
		}

		// 创建绑定
		binding := &model.StoreSupplierProduct{
			StoreID:   storeID,
			ProductID: productID,
			IsDefault: isDefault,
		}
		if err := m.db.Create(binding).Error; err != nil {
			return err
		}
	}
	return nil
}

// UnbindProducts 门店解绑供应商商品
func (m *StoreSupplierModule) UnbindProducts(storeID uint, productIDs []uint) error {
	return m.db.Where("store_id = ? AND product_id IN ?", storeID, productIDs).Delete(&model.StoreSupplierProduct{}).Error
}

// SetDefault 设置默认供应商
func (m *StoreSupplierModule) SetDefault(storeID, productID uint) error {
	// 先获取商品信息，找到商品名称
	var product model.SupplierProduct
	if err := m.db.First(&product, productID).Error; err != nil {
		return err
	}

	// 找到同名商品的所有ID
	var sameNameProducts []model.SupplierProduct
	m.db.Where("name = ?", product.Name).Find(&sameNameProducts)

	var sameNameIDs []uint
	for _, p := range sameNameProducts {
		sameNameIDs = append(sameNameIDs, p.ID)
	}

	// 将该门店下同名商品的所有绑定设为非默认
	m.db.Model(&model.StoreSupplierProduct{}).Where("store_id = ? AND product_id IN ?", storeID, sameNameIDs).Update("is_default", false)

	// 将指定商品设为默认
	return m.db.Model(&model.StoreSupplierProduct{}).Where("store_id = ? AND product_id = ?", storeID, productID).Update("is_default", true).Error
}

// ListByStoreID 获取门店绑定的所有商品
func (m *StoreSupplierModule) ListByStoreID(storeID uint) ([]*model.StoreSupplierProduct, error) {
	var bindings []*model.StoreSupplierProduct
	if err := m.db.Preload("Product.Supplier").Preload("Product.Category").Where("store_id = ?", storeID).Find(&bindings).Error; err != nil {
		return nil, err
	}
	return bindings, nil
}

// GetStoreProductByName 根据商品名称获取门店的默认供应商商品
func (m *StoreSupplierModule) GetStoreProductByName(storeID uint, productName string) (*model.StoreSupplierProduct, error) {
	var binding model.StoreSupplierProduct
	// 优先找默认的
	err := m.db.Preload("Product.Supplier").
		Joins("JOIN supplier_products ON supplier_products.id = store_supplier_products.product_id").
		Where("store_supplier_products.store_id = ? AND supplier_products.name = ?", storeID, productName).
		Order("store_supplier_products.is_default DESC").
		First(&binding).Error
	if err != nil {
		return nil, err
	}
	return &binding, nil
}

// ValidateStoreProducts 验证商品是否已绑定到门店
// 返回未绑定的商品ID列表
func (m *StoreSupplierModule) ValidateStoreProducts(storeID uint, productIDs []uint) ([]uint, error) {
	if len(productIDs) == 0 {
		return nil, nil
	}

	// 查询已绑定的商品ID
	var boundProductIDs []uint
	err := m.db.Model(&model.StoreSupplierProduct{}).
		Where("store_id = ? AND product_id IN ?", storeID, productIDs).
		Pluck("product_id", &boundProductIDs).Error
	if err != nil {
		return nil, err
	}

	// 构建已绑定商品ID的map
	boundMap := make(map[uint]bool)
	for _, id := range boundProductIDs {
		boundMap[id] = true
	}

	// 找出未绑定的商品ID
	var unboundIDs []uint
	for _, id := range productIDs {
		if !boundMap[id] {
			unboundIDs = append(unboundIDs, id)
		}
	}

	return unboundIDs, nil
}
