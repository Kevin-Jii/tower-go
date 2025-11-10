package service

import (
	"errors"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/cache"
)

type DishCategoryService struct {
	module     *module.DishCategoryModule
	dishModule *module.DishModule
}

func NewDishCategoryService(m *module.DishCategoryModule, dishM *module.DishModule) *DishCategoryService {
	return &DishCategoryService{module: m, dishModule: dishM}
}

func (s *DishCategoryService) Create(storeID uint, req *model.CreateDishCategoryReq) (*model.DishCategory, error) {
	// 唯一性校验（门店+名称）
	exists, err := s.module.ExistsByName(storeID, req.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("category name already exists in this store")
	}
	code := req.Code
	if code == "" { // 自动生成编码
		c, err := utils.GenerateDishCategoryCode(s.module.GetDB(), storeID)
		if err != nil {
			return nil, err
		}
		code = c
	}
	cat := &model.DishCategory{StoreID: storeID, Name: req.Name, Code: code, Sort: req.Sort, Status: 1, Remark: req.Remark}
	if err := s.module.Create(cat); err != nil {
		return nil, err
	}
	cache.InvalidateDishCategoryCache(storeID)
	return cat, nil
}

func (s *DishCategoryService) Update(id, storeID uint, req *model.UpdateDishCategoryReq) error {
	cat, err := s.module.GetByID(id, storeID)
	if err != nil {
		return errors.New("category not found")
	}
	// 名称更新需做唯一性校验（同门店不可重复）
	if req.Name != "" && req.Name != cat.Name {
		exists, err := s.module.ExistsByName(storeID, req.Name)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("category name already exists in this store")
		}
		cat.Name = req.Name
	} else if req.Name != "" {
		// 名称相同但仍传值，允许保持（无操作）
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
	cache.InvalidateDishCategoryCache(storeID)
	return nil
}

func (s *DishCategoryService) Delete(id, storeID uint) error {
	// 若分类下仍有菜品则禁止删除
	exists, err := s.dishModule.ExistsDishInCategory(storeID, id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("category has dishes")
	}
	if err := s.module.Delete(id, storeID); err != nil {
		return err
	}
	cache.InvalidateDishCategoryCache(storeID)
	return nil
}

func (s *DishCategoryService) List(storeID uint) ([]*model.DishCategory, error) {
	var cached []*model.DishCategory
	err := cache.DishCategoriesDomain.GetOrSet(&cached, func() (interface{}, error) {
		return s.module.List(storeID)
	}, storeID)
	return cached, err
}

func (s *DishCategoryService) ListWithDishes(storeID uint) ([]*model.DishCategory, error) {
	var cached []*model.DishCategory
	err := cache.DishCategoriesWithDishesDomain.GetOrSet(&cached, func() (interface{}, error) {
		return s.module.ListWithDishes(storeID)
	}, storeID)
	return cached, err
}

func (s *DishCategoryService) Reorder(storeID uint, items []model.ReorderDishCategoryItem) error {
	if err := s.module.BatchUpdateSort(storeID, items); err != nil {
		return err
	}
	cache.DishCategoriesDomain.Invalidate(storeID)
	cache.DishCategoriesWithDishesDomain.Invalidate(storeID)
	return nil
}

// ListDishesForCategory 列出指定门店某分类下的菜品（仅启用菜品）
func (s *DishCategoryService) ListDishesForCategory(storeID, categoryID uint) ([]*model.Dish, error) {
	// 先确认分类存在（避免越权或误请求）
	_, err := s.module.GetByID(categoryID, storeID)
	if err != nil {
		return nil, errors.New("category not found")
	}
	return s.dishModule.ListByCategory(storeID, categoryID)
}

// CreateDishInCategory 在指定门店与分类下创建菜品
// 若分类不存在则返回错误；忽略请求体中的 category_id 以路径参数为准
func (s *DishCategoryService) CreateDishInCategory(storeID, categoryID uint, req *model.CreateDishReq) (*model.Dish, error) {
	// 分类存在性校验
	if _, err := s.module.GetByID(categoryID, storeID); err != nil {
		return nil, errors.New("category not found")
	}
	// 同分类下菜品名称唯一校验
	if req.Name != "" {
		exists, err := s.dishModule.ExistsByNameInCategory(storeID, &categoryID, req.Name)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("dish name already exists in this category")
		}
	}
	// 强制设置分类 ID
	req.CategoryID = &categoryID
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
		return nil, err
	}
	// 缓存失效：分类相关与单菜品（若后续有单菜品缓存）
	cache.InvalidateDishCategoryCache(storeID)
	cache.InvalidateDishCache(dish.ID, storeID)
	return dish, nil
}

// UpdateDishInCategory 在指定门店与分类下更新菜品
// 校验菜品归属、名称唯一性；不允许跨分类更新（忽略请求体中的 category_id）
func (s *DishCategoryService) UpdateDishInCategory(storeID, categoryID, dishID uint, req *model.UpdateDishReq) error {
	// 1. 校验分类存在
	if _, err := s.module.GetByID(categoryID, storeID); err != nil {
		return errors.New("category not found")
	}
	// 2. 校验菜品存在且属于该分类
	dish, err := s.dishModule.GetByID(dishID, storeID)
	if err != nil {
		return errors.New("dish not found")
	}
	if dish.CategoryID == nil || *dish.CategoryID != categoryID {
		return errors.New("dish does not belong to this category")
	}
	// 3. 名称唯一性校验（若改名）
	if req.Name != "" && req.Name != dish.Name {
		exists, err := s.dishModule.ExistsByNameInCategory(storeID, &categoryID, req.Name)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("dish name already exists in this category")
		}
	}
	// 4. 忽略请求中的 category_id 避免跨分类移动（若需移动应另开接口）
	req.CategoryID = nil
	// 5. 执行更新
	if err := s.dishModule.UpdateByIDAndStoreID(dishID, storeID, req); err != nil {
		return err
	}
	// 6. 缓存失效
	cache.InvalidateDishCategoryCache(storeID)
	cache.InvalidateDishCache(dishID, storeID)
	return nil
}

// DeleteDishInCategory 在指定门店与分类下删除菜品
// 校验菜品归属后执行删除
func (s *DishCategoryService) DeleteDishInCategory(storeID, categoryID, dishID uint) error {
	// 1. 校验分类存在
	if _, err := s.module.GetByID(categoryID, storeID); err != nil {
		return errors.New("category not found")
	}
	// 2. 校验菜品存在且属于该分类
	dish, err := s.dishModule.GetByID(dishID, storeID)
	if err != nil {
		return errors.New("dish not found")
	}
	if dish.CategoryID == nil || *dish.CategoryID != categoryID {
		return errors.New("dish does not belong to this category")
	}
	// 3. 执行删除
	if err := s.dishModule.Delete(dishID, storeID); err != nil {
		return err
	}
	// 4. 缓存失效
	cache.InvalidateDishCategoryCache(storeID)
	cache.InvalidateDishCache(dishID, storeID)
	return nil
}
