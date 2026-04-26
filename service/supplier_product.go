package service

import (
	"errors"
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

func (s *SupplierProductService) ensureSingleLargeSpec(productID uint, currentUnitCode string, factorToBase float64, excludeID uint) error {
	if factorToBase <= 1 {
		return nil
	}
	specs, err := s.unitSpecModule.ListByProductID(productID)
	if err != nil {
		return err
	}
	current := strings.TrimSpace(currentUnitCode)
	for _, spec := range specs {
		if spec == nil || !spec.IsEnabled || spec.FactorToBase <= 1 {
			continue
		}
		if excludeID > 0 && spec.ID == excludeID {
			continue
		}
		if strings.TrimSpace(spec.UnitCode) != current {
			return errors.New("同一商品只能配置一个大规格单位，请先修改或停用原有大规格")
		}
	}
	return nil
}

func (s *SupplierProductService) validateUnitCodeAndGetName(unitCode string) (string, error) {
	if s.dictModule == nil {
		return "", nil
	}
	dictData, err := s.dictModule.GetDataByTypeAndValue("product_unit", strings.TrimSpace(unitCode))
	if err != nil || dictData == nil {
		return "", errors.New("unit_code 未在字典 product_unit 中定义")
	}
	if dictData.Status != 1 {
		return "", errors.New("unit_code 在字典中已禁用")
	}
	return dictData.Label, nil
}

type SupplierProductService struct {
	productModule  *module.SupplierProductModule
	unitSpecModule *module.ProductUnitSpecModule
	dictModule     *module.DictModule
	categoryModule *module.SupplierCategoryModule
	supplierModule *module.SupplierModule
}

func NewSupplierProductService(
	productModule *module.SupplierProductModule,
	unitSpecModule *module.ProductUnitSpecModule,
	dictModule *module.DictModule,
	categoryModule *module.SupplierCategoryModule,
	supplierModule *module.SupplierModule,
) *SupplierProductService {
	return &SupplierProductService{
		productModule:  productModule,
		unitSpecModule: unitSpecModule,
		dictModule:     dictModule,
		categoryModule: categoryModule,
		supplierModule: supplierModule,
	}
}

// CreateProduct 创建供应商商品
func (s *SupplierProductService) CreateProduct(req *model.CreateSupplierProductReq) error {
	// 验证供应商存在
	_, err := s.supplierModule.GetByID(req.SupplierID)
	if err != nil {
		return errors.New("supplier not found")
	}

	// 验证分类存在
	_, err = s.categoryModule.GetByID(req.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	product := &model.SupplierProduct{
		SupplierID:     req.SupplierID,
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Unit:           req.Unit,
		BottlePrice:    req.BottlePrice,
		CasePrice:      req.CasePrice,
		BottlesPerCase: req.BottlesPerCase,
		Spec:           req.Spec,
		Remark:         req.Remark,
		Status:         1,
	}
	// 兼容旧字段：如果仍传 price，优先使用 price 作为单瓶价
	if req.Price != nil {
		product.BottlePrice = *req.Price
	}
	// 兼容字段 price 默认与单瓶价保持一致
	product.Price = product.BottlePrice
	return s.productModule.Create(product)
}

// GetProduct 获取商品详情
func (s *SupplierProductService) GetProduct(id uint) (*model.SupplierProduct, error) {
	return s.productModule.GetByID(id)
}

// ListProducts 获取商品列表
func (s *SupplierProductService) ListProducts(req *model.ListSupplierProductReq) ([]*model.SupplierProduct, int64, error) {
	return s.productModule.List(req)
}

// UpdateProduct 更新商品
func (s *SupplierProductService) UpdateProduct(id uint, req *model.UpdateSupplierProductReq) error {
	_, err := s.productModule.GetByID(id)
	if err != nil {
		return errors.New("product not found")
	}
	// 兼容旧字段：更新 price 时同步单瓶价；更新单瓶价时同步兼容字段 price
	if req.Price != nil && req.BottlePrice == nil {
		req.BottlePrice = req.Price
	}
	if req.BottlePrice != nil {
		req.Price = req.BottlePrice
	}
	return s.productModule.UpdateByID(id, req)
}

// DeleteProduct 删除商品
func (s *SupplierProductService) DeleteProduct(id uint) error {
	_, err := s.productModule.GetByID(id)
	if err != nil {
		return errors.New("product not found")
	}
	return s.productModule.Delete(id)
}

// CreateCategory 创建供应商分类
func (s *SupplierProductService) CreateCategory(req *model.CreateSupplierCategoryReq) error {
	_, err := s.supplierModule.GetByID(req.SupplierID)
	if err != nil {
		return errors.New("supplier not found")
	}

	category := &model.SupplierCategory{
		SupplierID: req.SupplierID,
		Name:       req.Name,
		Sort:       req.Sort,
		Status:     1,
	}
	return s.categoryModule.Create(category)
}

// GetCategory 获取分类详情
func (s *SupplierProductService) GetCategory(id uint) (*model.SupplierCategory, error) {
	return s.categoryModule.GetByID(id)
}

// ListCategories 获取供应商的分类列表
func (s *SupplierProductService) ListCategories(supplierID uint) ([]*model.SupplierCategory, error) {
	return s.categoryModule.ListBySupplierID(supplierID)
}

// UpdateCategory 更新分类
func (s *SupplierProductService) UpdateCategory(id uint, req *model.UpdateSupplierCategoryReq) error {
	_, err := s.categoryModule.GetByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	return s.categoryModule.UpdateByID(id, req)
}

// DeleteCategory 删除分类
func (s *SupplierProductService) DeleteCategory(id uint) error {
	_, err := s.categoryModule.GetByID(id)
	if err != nil {
		return errors.New("category not found")
	}
	return s.categoryModule.Delete(id)
}

func (s *SupplierProductService) CreateUnitSpec(req *model.CreateProductUnitSpecReq) error {
	if _, err := s.productModule.GetByID(req.ProductID); err != nil {
		return errors.New("product not found")
	}
	dictUnitName, err := s.validateUnitCodeAndGetName(req.UnitCode)
	if err != nil {
		return err
	}
	if err := s.ensureSingleLargeSpec(req.ProductID, req.UnitCode, req.FactorToBase, 0); err != nil {
		return err
	}
	unitName := strings.TrimSpace(req.UnitName)
	if unitName == "" {
		unitName = dictUnitName
	}
	spec := &model.ProductUnitSpec{
		ProductID:    req.ProductID,
		UnitCode:     req.UnitCode,
		UnitName:     unitName,
		FactorToBase: req.FactorToBase,
		Precision:    req.Precision,
		CostPrice:    req.CostPrice,
		SalePrice:    req.SalePrice,
		IsEnabled:    true,
	}
	if req.IsEnabled != nil {
		spec.IsEnabled = *req.IsEnabled
	}
	return s.unitSpecModule.Create(spec)
}

func (s *SupplierProductService) ListUnitSpecs(productID uint) ([]*model.ProductUnitSpec, error) {
	if _, err := s.productModule.GetByID(productID); err != nil {
		return nil, errors.New("product not found")
	}
	return s.unitSpecModule.ListByProductID(productID)
}

func (s *SupplierProductService) UpdateUnitSpec(id uint, req *model.UpdateProductUnitSpecReq) error {
	existing, err := s.unitSpecModule.GetByID(id)
	if err != nil {
		return errors.New("unit spec not found")
	}
	nextUnitCode := existing.UnitCode
	if req.UnitCode != nil {
		nextUnitCode = *req.UnitCode
	}
	nextFactor := existing.FactorToBase
	if req.FactorToBase != nil {
		nextFactor = *req.FactorToBase
	}
	if req.UnitCode != nil {
		dictUnitName, err := s.validateUnitCodeAndGetName(*req.UnitCode)
		if err != nil {
			return err
		}
		if req.UnitName == nil || strings.TrimSpace(*req.UnitName) == "" {
			req.UnitName = &dictUnitName
		}
	}
	if err := s.ensureSingleLargeSpec(existing.ProductID, nextUnitCode, nextFactor, id); err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if req.UnitCode != nil {
		updates["unit_code"] = *req.UnitCode
	}
	if req.UnitName != nil {
		updates["unit_name"] = *req.UnitName
	}
	if req.FactorToBase != nil {
		updates["factor_to_base"] = *req.FactorToBase
	}
	if req.Precision != nil {
		updates["precision"] = *req.Precision
	}
	if req.CostPrice != nil {
		updates["cost_price"] = *req.CostPrice
	}
	if req.SalePrice != nil {
		updates["sale_price"] = *req.SalePrice
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	return s.unitSpecModule.UpdateByID(id, updates)
}

func (s *SupplierProductService) DeleteUnitSpec(id uint) error {
	if _, err := s.unitSpecModule.GetByID(id); err != nil {
		return errors.New("unit spec not found")
	}
	return s.unitSpecModule.DeleteByID(id)
}

func (s *SupplierProductService) BatchUpsertUnitSpecs(req *model.BatchUpsertProductUnitSpecsReq) error {
	if _, err := s.productModule.GetByID(req.ProductID); err != nil {
		return errors.New("product not found")
	}
	largeUnitCode := ""
	for _, unit := range req.Units {
		if unit.FactorToBase <= 1 {
			continue
		}
		if largeUnitCode == "" {
			largeUnitCode = unit.UnitCode
			continue
		}
		if strings.TrimSpace(largeUnitCode) != strings.TrimSpace(unit.UnitCode) {
			return errors.New("同一商品只能配置一个大规格单位")
		}
	}
	if largeUnitCode != "" {
		if err := s.ensureSingleLargeSpec(req.ProductID, largeUnitCode, 2, 0); err != nil {
			return err
		}
	}
	for _, unit := range req.Units {
		dictUnitName, err := s.validateUnitCodeAndGetName(unit.UnitCode)
		if err != nil {
			return err
		}
		unitName := strings.TrimSpace(unit.UnitName)
		if unitName == "" {
			unitName = dictUnitName
		}
		spec := &model.ProductUnitSpec{
			ProductID:    req.ProductID,
			UnitCode:     unit.UnitCode,
			UnitName:     unitName,
			FactorToBase: unit.FactorToBase,
			Precision:    unit.Precision,
			CostPrice:    unit.CostPrice,
			SalePrice:    unit.SalePrice,
			IsEnabled:    true,
		}
		if unit.IsEnabled != nil {
			spec.IsEnabled = *unit.IsEnabled
		}
		if err := s.unitSpecModule.UpsertByProductAndUnit(spec); err != nil {
			return err
		}
	}
	return nil
}
