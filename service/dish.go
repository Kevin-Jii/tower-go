package service

import (
	"errors"
	"tower-go/model"
	"tower-go/module"
	"tower-go/utils/cache"
)

type DishService struct {
	dishModule *module.DishModule
}

func NewDishService(dishModule *module.DishModule) *DishService {
	return &DishService{dishModule: dishModule}
}

// CreateDish 创建菜品（自动关联当前门店）
func (s *DishService) CreateDish(storeID uint, req *model.CreateDishReq) error {
	// 同分类下名称唯一校验（仅在指定分类时；未分类允许重复不同分类的同名）
	if req.Name != "" {
		exists, err := s.dishModule.ExistsByNameInCategory(storeID, req.CategoryID, req.Name)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("dish name already exists in this category")
		}
	}
	dish := &model.Dish{
		StoreID:    storeID,
		Name:       req.Name,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Image:      req.Image,
		Remark:     req.Remark,
		Status:     1, // 默认上架
	}
	if err := s.dishModule.Create(dish); err != nil {
		return err
	}
	cache.InvalidateDishCategoryCache(storeID)
	return nil
}

// GetDish 获取菜品详情（门店隔离）
func (s *DishService) GetDish(id, storeID uint) (*model.Dish, error) {
	return s.dishModule.GetByID(id, storeID)
}

// ListDishes 获取门店菜品列表
func (s *DishService) ListDishes(storeID uint, page, pageSize int) ([]*model.Dish, int64, error) {
	return s.dishModule.List(storeID, page, pageSize)
}

// ListDishesByCategory 根据分类获取菜品
func (s *DishService) ListDishesByCategory(storeID uint, categoryID uint) ([]*model.Dish, error) {
	return s.dishModule.ListByCategory(storeID, categoryID)
}

// UpdateDish 更新菜品信息
func (s *DishService) UpdateDish(id, storeID uint, req *model.UpdateDishReq) error {
	// 为保持“存在性”语义，先检查是否存在
	if _, err := s.dishModule.GetByID(id, storeID); err != nil {
		return errors.New("dish not found")
	}
	if err := s.dishModule.UpdateByIDAndStoreID(id, storeID, req); err != nil {
		return err
	}
	cache.InvalidateDishCache(id, storeID)
	return nil
}

// DeleteDish 删除菜品
func (s *DishService) DeleteDish(id, storeID uint) error {
	if err := s.dishModule.Delete(id, storeID); err != nil {
		return err
	}
	cache.InvalidateDishCache(id, storeID)
	cache.InvalidateDishCategoryCache(storeID)
	return nil
}
