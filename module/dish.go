package module

import (
	"tower-go/model"

	"gorm.io/gorm"
)

type DishModule struct {
	db *gorm.DB
}

func NewDishModule(db *gorm.DB) *DishModule {
	return &DishModule{db: db}
}

// Create 创建菜品
func (m *DishModule) Create(dish *model.Dish) error {
	return m.db.Create(dish).Error
}

// GetByID 根据ID和门店ID获取菜品（数据隔离）
func (m *DishModule) GetByID(id, storeID uint) (*model.Dish, error) {
	var dish model.Dish
	if err := m.db.Where("id = ? AND store_id = ?", id, storeID).First(&dish).Error; err != nil {
		return nil, err
	}
	return &dish, nil
}

// List 获取指定门店的菜品列表（支持分页）
func (m *DishModule) List(storeID uint, page, pageSize int) ([]*model.Dish, int64, error) {
	var dishes []*model.Dish
	var total int64

	query := m.db.Model(&model.Dish{}).Where("store_id = ?", storeID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&dishes).Error; err != nil {
		return nil, 0, err
	}

	return dishes, total, nil
}

// ListByCategory 根据分类获取菜品列表
func (m *DishModule) ListByCategory(storeID uint, category string) ([]*model.Dish, error) {
	var dishes []*model.Dish
	if err := m.db.Where("store_id = ? AND category = ? AND status = 1", storeID, category).Find(&dishes).Error; err != nil {
		return nil, err
	}
	return dishes, nil
}

// Update 更新菜品信息
func (m *DishModule) Update(dish *model.Dish) error {
	return m.db.Save(dish).Error
}

// Delete 删除菜品
func (m *DishModule) Delete(id, storeID uint) error {
	return m.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&model.Dish{}).Error
}
