package module

import (
	"tower-go/model"
	updatesPkg "tower-go/utils/updates"

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

// ExistsByNameInCategory 判断同门店同分类下是否存在同名菜品
// 若 categoryID 为 nil，则只判断 store_id/category_id IS NULL 场景（允许不同分类同名）
func (m *DishModule) ExistsByNameInCategory(storeID uint, categoryID *uint, name string) (bool, error) {
	var count int64
	q := m.db.Model(&model.Dish{}).Where("store_id = ? AND name = ?", storeID, name)
	if categoryID == nil {
		q = q.Where("category_id IS NULL")
	} else {
		q = q.Where("category_id = ?", *categoryID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
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
func (m *DishModule) ListByCategory(storeID uint, categoryID uint) ([]*model.Dish, error) {
	var dishes []*model.Dish
	if err := m.db.Where("store_id = ? AND category_id = ? AND status = 1", storeID, categoryID).Find(&dishes).Error; err != nil {
		return nil, err
	}
	return dishes, nil
}

// Update 更新菜品信息
func (m *DishModule) Update(dish *model.Dish) error {
	return m.db.Save(dish).Error
}

// UpdateByIDAndStoreID 使用动态更新构造器按需更新菜品字段
func (m *DishModule) UpdateByIDAndStoreID(id, storeID uint, req *model.UpdateDishReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}
	return m.db.Model(&model.Dish{}).Where("id = ? AND store_id = ?", id, storeID).Updates(updateMap).Error
}

// Delete 删除菜品
func (m *DishModule) Delete(id, storeID uint) error {
	return m.db.Where("id = ? AND store_id = ?", id, storeID).Delete(&model.Dish{}).Error
}

// ExistsDishInCategory 判断分类下是否仍存在任意菜品（不限状态）
func (m *DishModule) ExistsDishInCategory(storeID, categoryID uint) (bool, error) {
	var count int64
	if err := m.db.Model(&model.Dish{}).Where("store_id = ? AND category_id = ?", storeID, categoryID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
