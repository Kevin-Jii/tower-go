package service

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils"
)

type StoreSupplierService struct {
	storeSupplierModule *module.StoreSupplierModule
	unitSpecModule      *module.ProductUnitSpecModule
}

func NewStoreSupplierService(
	storeSupplierModule *module.StoreSupplierModule,
	unitSpecModule *module.ProductUnitSpecModule,
) *StoreSupplierService {
	return &StoreSupplierService{
		storeSupplierModule: storeSupplierModule,
		unitSpecModule:      unitSpecModule,
	}
}

// BindSuppliers 门店绑定供应商
func (s *StoreSupplierService) BindSuppliers(storeID uint, supplierIDs []uint) error {
	if err := s.storeSupplierModule.BindSuppliers(storeID, supplierIDs); err != nil {
		return err
	}

	// 发布供应商绑定事件
	utils.GlobalEventBus.Publish(utils.EventSupplierBound, map[string]interface{}{
		"store_id":     storeID,
		"supplier_ids": supplierIDs,
	})

	return nil
}

// UnbindSuppliers 门店解绑供应商
func (s *StoreSupplierService) UnbindSuppliers(storeID uint, supplierIDs []uint) error {
	return s.storeSupplierModule.UnbindSuppliers(storeID, supplierIDs)
}

// ListSuppliersByStoreID 获取门店绑定的所有供应商
func (s *StoreSupplierService) ListSuppliersByStoreID(storeID uint) ([]*model.StoreSupplier, error) {
	return s.storeSupplierModule.ListSuppliersByStoreID(storeID)
}

// ListProductsByStoreID 获取门店可采购的商品列表（绑定供应商的所有商品）
func (s *StoreSupplierService) ListProductsByStoreID(storeID, supplierID, categoryID uint, keyword string) ([]*model.SupplierProduct, error) {
	products, err := s.storeSupplierModule.ListProductsByStoreID(storeID, supplierID, categoryID, keyword)
	if err != nil {
		return nil, err
	}
	if s.unitSpecModule == nil || len(products) == 0 {
		return products, nil
	}
	ids := make([]uint, 0, len(products))
	for _, p := range products {
		if p != nil {
			ids = append(ids, p.ID)
		}
	}
	specs, err := s.unitSpecModule.ListByProductIDs(ids)
	if err != nil {
		return nil, err
	}
	byProduct := make(map[uint][]*model.ProductUnitSpec)
	for _, sp := range specs {
		if sp == nil {
			continue
		}
		byProduct[sp.ProductID] = append(byProduct[sp.ProductID], sp)
	}
	for _, p := range products {
		if p == nil {
			continue
		}
		if v, ok := byProduct[p.ID]; ok {
			p.UnitSpecs = v
		} else {
			p.UnitSpecs = []*model.ProductUnitSpec{}
		}
	}
	return products, nil
}

// ListCategoriesByStoreID 获取门店绑定供应商下的分类列表
func (s *StoreSupplierService) ListCategoriesByStoreID(storeID, supplierID uint) ([]*model.SupplierCategory, error) {
	return s.storeSupplierModule.ListCategoriesByStoreID(storeID, supplierID)
}

// ValidateStoreProducts 校验商品是否属于门店已绑定的供应商（返回不可用的商品 ID）
func (s *StoreSupplierService) ValidateStoreProducts(storeID uint, productIDs []uint) ([]uint, error) {
	return s.storeSupplierModule.ValidateStoreProducts(storeID, productIDs)
}
