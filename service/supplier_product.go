package service

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type SupplierProductService struct {
	productModule  *module.SupplierProductModule
	categoryModule *module.SupplierCategoryModule
	supplierModule *module.SupplierModule
}

func NewSupplierProductService(
	productModule *module.SupplierProductModule,
	categoryModule *module.SupplierCategoryModule,
	supplierModule *module.SupplierModule,
) *SupplierProductService {
	return &SupplierProductService{
		productModule:  productModule,
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
		SupplierID: req.SupplierID,
		CategoryID: req.CategoryID,
		Name:       req.Name,
		Unit:       req.Unit,
		Price:      req.Price,
		Spec:       req.Spec,
		Remark:     req.Remark,
		Status:     1,
	}
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
