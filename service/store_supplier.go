package service

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils"
)

type StoreSupplierService struct {
	storeSupplierModule *module.StoreSupplierModule
}

func NewStoreSupplierService(storeSupplierModule *module.StoreSupplierModule) *StoreSupplierService {
	return &StoreSupplierService{storeSupplierModule: storeSupplierModule}
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
	return s.storeSupplierModule.ListProductsByStoreID(storeID, supplierID, categoryID, keyword)
}
