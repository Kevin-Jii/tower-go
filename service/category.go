package service

import (
	"errors"
	"fmt"
	"time"
	"tower-go/model"
	"tower-go/module"
	"tower-go/utils"
)

type DishCategoryService struct {
	module *module.DishCategoryModule
}

func NewDishCategoryService(m *module.DishCategoryModule) *DishCategoryService {
	return &DishCategoryService{module: m}
}

func (s *DishCategoryService) Create(storeID uint, req *model.CreateDishCategoryReq) error {
	cat := &model.DishCategory{StoreID: storeID, Name: req.Name, Code: req.Code, Sort: req.Sort, Status: 1, Remark: req.Remark}
	if err := s.module.Create(cat); err != nil {
		return err
	}
	utils.InvalidateDishCategoryCache(storeID)
	return nil
}

func (s *DishCategoryService) Update(id, storeID uint, req *model.UpdateDishCategoryReq) error {
	cat, err := s.module.GetByID(id, storeID)
	if err != nil {
		return errors.New("category not found")
	}
	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.Code != "" {
		cat.Code = req.Code
	}
	if req.Sort != nil {
		cat.Sort = *req.Sort
	}
	if req.Status != nil {
		cat.Status = *req.Status
	}
	if req.Remark != "" {
		cat.Remark = req.Remark
	}
	if err := s.module.Update(cat); err != nil {
		return err
	}
	utils.InvalidateDishCategoryCache(storeID)
	return nil
}

func (s *DishCategoryService) Delete(id, storeID uint) error {
	if err := s.module.Delete(id, storeID); err != nil {
		return err
	}
	utils.InvalidateDishCategoryCache(storeID)
	return nil
}

func (s *DishCategoryService) List(storeID uint) ([]*model.DishCategory, error) {
	key := fmt.Sprintf(utils.CacheKeyDishCategories, storeID)
	var cached []*model.DishCategory
	err := utils.CacheGetOrSet(key, &cached, time.Minute*10, func() (interface{}, error) {
		return s.module.List(storeID)
	})
	return cached, err
}

func (s *DishCategoryService) ListWithDishes(storeID uint) ([]*model.DishCategory, error) {
	key := fmt.Sprintf(utils.CacheKeyDishCategoriesWithDishes, storeID)
	var cached []*model.DishCategory
	err := utils.CacheGetOrSet(key, &cached, time.Minute*5, func() (interface{}, error) {
		return s.module.ListWithDishes(storeID)
	})
	return cached, err
}

func (s *DishCategoryService) Reorder(storeID uint, items []model.ReorderDishCategoryItem) error {
	if err := s.module.BatchUpdateSort(storeID, items); err != nil {
		return err
	}
	utils.InvalidateDishCategoryCache(storeID)
	return nil
}
