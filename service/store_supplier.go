package service

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type StoreSupplierService struct {
	storeSupplierModule *module.StoreSupplierModule
}

func NewStoreSupplierService(storeSupplierModule *module.StoreSupplierModule) *StoreSupplierService {
	return &StoreSupplierService{storeSupplierModule: storeSupplierModule}
}

// BindProducts 门店绑定供应商商品
func (s *StoreSupplierService) BindProducts(storeID uint, productIDs []uint) error {
	return s.storeSupplierModule.BindProducts(storeID, productIDs)
}

// UnbindProducts 门店解绑供应商商品
func (s *StoreSupplierService) UnbindProducts(storeID uint, productIDs []uint) error {
	return s.storeSupplierModule.UnbindProducts(storeID, productIDs)
}

// SetDefault 设置默认供应商
func (s *StoreSupplierService) SetDefault(storeID, productID uint) error {
	return s.storeSupplierModule.SetDefault(storeID, productID)
}

// ListByStoreID 获取门店绑定的所有商品
func (s *StoreSupplierService) ListByStoreID(storeID uint) ([]*model.StoreSupplierProduct, error) {
	return s.storeSupplierModule.ListByStoreID(storeID)
}
