package module

import (
	"tower-go/model"
	"tower-go/utils/batch"

	"gorm.io/gorm"
)

type DishCategoryModule struct {
	db *gorm.DB
}

func NewDishCategoryModule(db *gorm.DB) *DishCategoryModule {
	return &DishCategoryModule{db: db}
}

// GetDB 公开底层 DB 供编码生成使用
func (m *DishCategoryModule) GetDB() *gorm.DB { return m.db }

func (m *DishCategoryModule) Create(cat *model.DishCategory) error {
	return m.db.Create(cat).Error
}

// ExistsByName 判断某门店下分类名是否已存在
func (m *DishCategoryModule) ExistsByName(storeID uint, name string) (bool, error) {
	var count int64
	err := m.db.Model(&model.DishCategory{}).Where("store_id = ? AND name = ?", storeID, name).Count(&count).Error
	return count > 0, err
}

func (m *DishCategoryModule) Update(cat *model.DishCategory) error {
	return m.db.Save(cat).Error
}

func (m *DishCategoryModule) GetByID(id, storeID uint) (*model.DishCategory, error) {
	var c model.DishCategory
	if err := m.db.Where("id = ? AND store_id = ?", id, storeID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (m *DishCategoryModule) Delete(id, storeID uint) error {
	return m.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&model.DishCategory{}).Error
}

func (m *DishCategoryModule) List(storeID uint) ([]*model.DishCategory, error) {
	var cats []*model.DishCategory
	if err := m.db.Where("store_id = ?", storeID).Order("sort ASC, id ASC").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

// ListWithDishes 预加载该门店分类及分类下启用的菜品
func (m *DishCategoryModule) ListWithDishes(storeID uint) ([]*model.DishCategory, error) {
	var cats []*model.DishCategory
	if err := m.db.Where("store_id = ? AND status = 1", storeID).
		Preload("Dishes", "store_id = ? AND status = 1", storeID).
		Order("sort ASC, id ASC").Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

// BatchUpdateSort 批量更新排序
func (m *DishCategoryModule) BatchUpdateSort(storeID uint, items []model.ReorderDishCategoryItem) error {
	if len(items) == 0 {
		return nil
	}
	// 转换为通用 SortItem
	list := make([]batch.SortItem, 0, len(items))
	for _, it := range items {
		list = append(list, batch.SortItem{ID: it.ID, Sort: it.Sort})
	}
	return batch.BatchUpdateSort(m.db, &model.DishCategory{}, storeID, list)
}
