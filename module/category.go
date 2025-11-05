package module

import (
	"tower-go/model"

	"gorm.io/gorm"
)

type DishCategoryModule struct {
	db *gorm.DB
}

func NewDishCategoryModule(db *gorm.DB) *DishCategoryModule {
	return &DishCategoryModule{db: db}
}

func (m *DishCategoryModule) Create(cat *model.DishCategory) error {
	return m.db.Create(cat).Error
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
	return m.db.Transaction(func(tx *gorm.DB) error {
		for _, it := range items {
			if err := tx.Model(&model.DishCategory{}).Where("id = ? AND store_id = ?", it.ID, storeID).Update("sort", it.Sort).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
