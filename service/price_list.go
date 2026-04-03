package service

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
)

type PriceListService struct {
	priceListModule       *module.PriceListModule
	storeModule           *module.StoreModule
	supplierProductModule *module.SupplierProductModule
}

func NewPriceListService(priceListModule *module.PriceListModule, storeModule *module.StoreModule, supplierProductModule *module.SupplierProductModule) *PriceListService {
	return &PriceListService{
		priceListModule:       priceListModule,
		storeModule:           storeModule,
		supplierProductModule: supplierProductModule,
	}
}

// ===== 价目单相关 =====

// CreatePriceList 创建价目单
func (s *PriceListService) CreatePriceList(req *model.CreatePriceListReq) error {
	// 验证门店存在
	_, err := s.storeModule.GetByID(req.StoreID)
	if err != nil {
		return errors.New("store not found")
	}

	// 如果设置为默认，先清除其他默认标记
	if req.IsDefault == 1 {
		s.priceListModule.ClearDefaultPriceList(req.StoreID)
	}

	priceList := &model.PriceList{
		StoreID:     req.StoreID,
		Name:        req.Name,
		Logo:        req.Logo,
		Description: req.Description,
		Status:      1,
		IsDefault:   req.IsDefault,
	}

	return s.priceListModule.CreatePriceList(priceList)
}

// UpdatePriceList 更新价目单
func (s *PriceListService) UpdatePriceList(id uint, req *model.UpdatePriceListReq) error {
	priceList, err := s.priceListModule.GetPriceListByID(id)
	if err != nil {
		return errors.New("price list not found")
	}

	// 如果设置为默认，先清除其他默认标记
	if req.IsDefault != nil && *req.IsDefault == 1 {
		s.priceListModule.ClearDefaultPriceList(priceList.StoreID)
	}

	updates := updatesPkg.BuildUpdatesFromReq(req)
	if len(updates) == 0 {
		return nil
	}

	return s.priceListModule.UpdatePriceList(id, updates)
}

// DeletePriceList 删除价目单
func (s *PriceListService) DeletePriceList(id uint) error {
	return s.priceListModule.DeletePriceList(id)
}

// GetPriceList 获取价目单详情
func (s *PriceListService) GetPriceList(id uint) (*model.PriceList, error) {
	return s.priceListModule.GetPriceListByID(id)
}

// ListPriceLists 获取门店的价目单列表
func (s *PriceListService) ListPriceLists(storeID uint) ([]*model.PriceList, error) {
	return s.priceListModule.ListPriceListsByStore(storeID)
}

// GetPriceListWithDetails 获取价目单完整结构
func (s *PriceListService) GetPriceListWithDetails(id uint) (*model.PriceListResp, error) {
	priceList, categories, itemsByCategory, err := s.priceListModule.GetPriceListWithDetails(id)
	if err != nil {
		return nil, err
	}

	resp := &model.PriceListResp{
		PriceList:  *priceList,
		Categories: []model.PriceListCategoryResp{},
	}

	// 构建分类和商品结构
	for _, category := range categories {
		categoryResp := model.PriceListCategoryResp{
			PriceListCategory: *category,
			Items:             []model.PriceListItemResp{},
		}

		// 添加商品
		if items, ok := itemsByCategory[category.ID]; ok {
			for _, item := range items {
				itemResp := s.buildItemResp(item)
				categoryResp.Items = append(categoryResp.Items, *itemResp)
			}
		}

		resp.Categories = append(resp.Categories, categoryResp)
	}

	return resp, nil
}

// ===== 价目单分类相关 =====

// CreateCategory 创建价目单分类
func (s *PriceListService) CreateCategory(req *model.CreatePriceListCategoryReq) error {
	// 验证价目单存在
	_, err := s.priceListModule.GetPriceListByID(req.PriceListID)
	if err != nil {
		return errors.New("price list not found")
	}

	category := &model.PriceListCategory{
		PriceListID: req.PriceListID,
		MainTitle:   req.MainTitle,
		SubTitle:    req.SubTitle,
		Sort:        req.Sort,
	}

	return s.priceListModule.CreateCategory(category)
}

// UpdateCategory 更新价目单分类
func (s *PriceListService) UpdateCategory(id uint, req *model.UpdatePriceListCategoryReq) error {
	_, err := s.priceListModule.GetCategoryByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	updates := updatesPkg.BuildUpdatesFromReq(req)
	if len(updates) == 0 {
		return nil
	}

	return s.priceListModule.UpdateCategory(id, updates)
}

// DeleteCategory 删除价目单分类
func (s *PriceListService) DeleteCategory(id uint) error {
	return s.priceListModule.DeleteCategory(id)
}

// ===== 价目单商品相关 =====

// AddItem 添加价目单商品
func (s *PriceListService) AddItem(req *model.AddPriceListItemReq) error {
	// 验证分类存在
	_, err := s.priceListModule.GetCategoryByID(req.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// 验证商品存在
	product, err := s.supplierProductModule.GetByID(req.ProductID)
	if err != nil {
		return errors.New("product not found")
	}

	// 如果没有指定显示名称，使用商品名称
	displayName := req.DisplayName
	if displayName == "" {
		displayName = product.Name
	}

	// 如果没有指定单位，使用商品单位
	unit := req.Unit
	if unit == "" {
		unit = product.Unit
	}

	item := &model.PriceListItem{
		CategoryID:  req.CategoryID,
		ProductID:   req.ProductID,
		DisplayName: displayName,
		Price:       req.Price,
		Unit:        unit,
		Spec:        req.Spec,
		Sort:        req.Sort,
		Status:      1,
	}

	return s.priceListModule.AddItem(item)
}

// UpdateItem 更新价目单商品
func (s *PriceListService) UpdateItem(id uint, req *model.UpdatePriceListItemReq) error {
	_, err := s.priceListModule.GetItemByID(id)
	if err != nil {
		return errors.New("item not found")
	}

	updates := updatesPkg.BuildUpdatesFromReq(req)
	if len(updates) == 0 {
		return nil
	}

	return s.priceListModule.UpdateItem(id, updates)
}

// DeleteItem 删除价目单商品
func (s *PriceListService) DeleteItem(id uint) error {
	return s.priceListModule.DeleteItem(id)
}

// BatchAddItems 批量添加价目单商品到指定分类
func (s *PriceListService) BatchAddItems(req *model.BatchAddPriceListItemsReq) error {
	// 验证分类存在
	_, err := s.priceListModule.GetCategoryByID(req.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// 批量查出商品信息
	productIDs := make([]uint, 0, len(req.Products))
	for _, p := range req.Products {
		productIDs = append(productIDs, p.ProductID)
	}
	products, err := s.supplierProductModule.GetByIDs(productIDs)
	if err != nil {
		return errors.New("failed to query products")
	}

	// 构建商品 ID -> 商品 的 map
	productMap := make(map[uint]*model.SupplierProduct, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	// 构建批量插入列表
	items := make([]*model.PriceListItem, 0, len(req.Products))
	for _, entry := range req.Products {
		product, ok := productMap[entry.ProductID]
		if !ok {
			continue // 跳过不存在的商品
		}

		displayName := entry.DisplayName
		if displayName == "" {
			displayName = product.Name
		}
		unit := entry.Unit
		if unit == "" {
			unit = product.Unit
		}

		items = append(items, &model.PriceListItem{
			CategoryID:  req.CategoryID,
			ProductID:   entry.ProductID,
			DisplayName: displayName,
			Price:       entry.Price,
			Unit:        unit,
			Spec:        entry.Spec,
			Sort:        entry.Sort,
			Status:      1,
		})
	}

	return s.priceListModule.BatchAddItems(items)
}

// buildItemResp 构建商品响应
func (s *PriceListService) buildItemResp(item *model.PriceListItem) *model.PriceListItemResp {
	resp := &model.PriceListItemResp{
		PriceListItem: *item,
	}

	if item.Product != nil {
		resp.ProductName = item.Product.Name
		if item.Product.Supplier != nil {
			resp.SupplierName = item.Product.Supplier.SupplierName
		}
	}

	return resp
}
