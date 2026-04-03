package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type PriceListModule struct {
	db *gorm.DB
}

func NewPriceListModule(db *gorm.DB) *PriceListModule {
	return &PriceListModule{db: db}
}

// ===== 价目单相关 =====

// CreatePriceList 创建价目单
func (m *PriceListModule) CreatePriceList(priceList *model.PriceList) error {
	return m.db.Create(priceList).Error
}

// UpdatePriceList 更新价目单
func (m *PriceListModule) UpdatePriceList(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.PriceList{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePriceList 删除价目单
func (m *PriceListModule) DeletePriceList(id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 删除价目单下的所有分类和商品
		if err := tx.Exec("DELETE FROM price_list_items WHERE category_id IN (SELECT id FROM price_list_categories WHERE price_list_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Where("price_list_id = ?", id).Delete(&model.PriceListCategory{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.PriceList{}, id).Error
	})
}

// GetPriceListByID 根据ID获取价目单
func (m *PriceListModule) GetPriceListByID(id uint) (*model.PriceList, error) {
	var priceList model.PriceList
	err := m.db.Preload("Store").Where("id = ?", id).First(&priceList).Error
	return &priceList, err
}

// ListPriceListsByStore 获取门店的价目单列表
func (m *PriceListModule) ListPriceListsByStore(storeID uint) ([]*model.PriceList, error) {
	var priceLists []*model.PriceList
	err := m.db.Where("store_id = ?", storeID).Order("is_default DESC, id DESC").Find(&priceLists).Error
	return priceLists, err
}

// GetDefaultPriceList 获取门店的默认价目单
func (m *PriceListModule) GetDefaultPriceList(storeID uint) (*model.PriceList, error) {
	var priceList model.PriceList
	err := m.db.Where("store_id = ? AND is_default = 1", storeID).First(&priceList).Error
	return &priceList, err
}

// ClearDefaultPriceList 清除门店的默认价目单标记
func (m *PriceListModule) ClearDefaultPriceList(storeID uint) error {
	return m.db.Model(&model.PriceList{}).Where("store_id = ?", storeID).Update("is_default", 0).Error
}

// ===== 价目单分类相关 =====

// CreateCategory 创建价目单分类
func (m *PriceListModule) CreateCategory(category *model.PriceListCategory) error {
	return m.db.Create(category).Error
}

// UpdateCategory 更新价目单分类
func (m *PriceListModule) UpdateCategory(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.PriceListCategory{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCategory 删除价目单分类
func (m *PriceListModule) DeleteCategory(id uint) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		// 删除分类下的所有商品
		if err := tx.Where("category_id = ?", id).Delete(&model.PriceListItem{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.PriceListCategory{}, id).Error
	})
}

// GetCategoryByID 根据ID获取价目单分类
func (m *PriceListModule) GetCategoryByID(id uint) (*model.PriceListCategory, error) {
	var category model.PriceListCategory
	err := m.db.Where("id = ?", id).First(&category).Error
	return &category, err
}

// ListCategoriesByPriceList 获取价目单的分类列表
func (m *PriceListModule) ListCategoriesByPriceList(priceListID uint) ([]*model.PriceListCategory, error) {
	var categories []*model.PriceListCategory
	err := m.db.Where("price_list_id = ?", priceListID).Order("sort ASC, id ASC").Find(&categories).Error
	return categories, err
}

// ===== 价目单商品相关 =====

// AddItem 添加价目单商品
func (m *PriceListModule) AddItem(item *model.PriceListItem) error {
	return m.db.Create(item).Error
}

// BatchAddItems 批量添加价目单商品
func (m *PriceListModule) BatchAddItems(items []*model.PriceListItem) error {
	if len(items) == 0 {
		return nil
	}
	return m.db.Create(&items).Error
}

// UpdateItem 更新价目单商品
func (m *PriceListModule) UpdateItem(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.PriceListItem{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteItem 删除价目单商品
func (m *PriceListModule) DeleteItem(id uint) error {
	return m.db.Delete(&model.PriceListItem{}, id).Error
}

// GetItemByID 根据ID获取价目单商品
func (m *PriceListModule) GetItemByID(id uint) (*model.PriceListItem, error) {
	var item model.PriceListItem
	err := m.db.Preload("Product").Preload("Product.Supplier").Where("id = ?", id).First(&item).Error
	return &item, err
}

// ListItemsByCategory 获取分类下的商品列表
func (m *PriceListModule) ListItemsByCategory(categoryID uint) ([]*model.PriceListItem, error) {
	var items []*model.PriceListItem
	err := m.db.Preload("Product").Preload("Product.Supplier").
		Where("category_id = ?", categoryID).
		Order("sort ASC, id ASC").
		Find(&items).Error
	return items, err
}

// GetPriceListWithDetails 获取价目单完整结构（包含分类和商品）
func (m *PriceListModule) GetPriceListWithDetails(id uint) (*model.PriceList, []*model.PriceListCategory, map[uint][]*model.PriceListItem, error) {
	// 获取价目单
	priceList, err := m.GetPriceListByID(id)
	if err != nil {
		return nil, nil, nil, err
	}

	// 获取分类列表
	categories, err := m.ListCategoriesByPriceList(id)
	if err != nil {
		return nil, nil, nil, err
	}

	// 获取所有商品并按分类分组
	itemsByCategory := make(map[uint][]*model.PriceListItem)
	for _, category := range categories {
		items, err := m.ListItemsByCategory(category.ID)
		if err != nil {
			return nil, nil, nil, err
		}
		itemsByCategory[category.ID] = items
	}

	return priceList, categories, itemsByCategory, nil
}
