package service

import (
	"errors"
	"fmt"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type SupplierService struct {
	supplierModule      *module.SupplierModule
	storeSupplierModule *module.StoreSupplierModule
}

func NewSupplierService(supplierModule *module.SupplierModule, storeSupplierModule *module.StoreSupplierModule) *SupplierService {
	return &SupplierService{
		supplierModule:      supplierModule,
		storeSupplierModule: storeSupplierModule,
	}
}

// CreateSupplier 创建供应商（自动生成编码：门店ID + 4位序号）
func (s *SupplierService) CreateSupplier(storeID uint, req *model.CreateSupplierReq) error {
	// 生成供应商编码：门店ID + 4位序号
	supplierCode, err := s.generateSupplierCode(storeID)
	if err != nil {
		return err
	}

	supplier := &model.Supplier{
		SupplierCode:    supplierCode,
		SupplierName:    req.SupplierName,
		ContactPerson:   req.ContactPerson,
		ContactPhone:    req.ContactPhone,
		ContactEmail:    req.ContactEmail,
		SupplierAddress: req.SupplierAddress,
		Remark:          req.Remark,
		Status:          1,
	}
	if err := s.supplierModule.Create(supplier); err != nil {
		return err
	}
	if s.storeSupplierModule != nil && storeID > 0 {
		return s.storeSupplierModule.BindSuppliers(storeID, []uint{supplier.ID})
	}
	return nil
}

// generateSupplierCode 生成供应商编码：门店ID + 4位序号
func (s *SupplierService) generateSupplierCode(storeID uint) (string, error) {
	// 获取当前最大序号
	maxSeq, err := s.supplierModule.GetMaxSeqByStorePrefix(storeID)
	if err != nil {
		return "", err
	}
	nextSeq := maxSeq + 1
	// 格式：门店ID + 4位序号，如 9990001
	return fmt.Sprintf("%d%04d", storeID, nextSeq), nil
}

func (s *SupplierService) GetSupplier(id uint) (*model.Supplier, error) {
	return s.supplierModule.GetByID(id)
}

func (s *SupplierService) GetSupplierScoped(id, storeID uint, hqUnbound bool) (*model.Supplier, error) {
	supplier, err := s.supplierModule.GetByID(id)
	if err != nil {
		return nil, err
	}
	if hqUnbound {
		return supplier, nil
	}
	if storeID == 0 {
		return nil, errors.New("current user has no store")
	}
	if s.storeSupplierModule == nil {
		return nil, errors.New("store supplier module not configured")
	}
	bound, err := s.storeSupplierModule.IsSupplierBound(storeID, id)
	if err != nil {
		return nil, err
	}
	if !bound {
		return nil, errors.New("supplier not bound to current store")
	}
	return supplier, nil
}

func (s *SupplierService) ListSuppliers(req *model.ListSupplierReq) ([]*model.Supplier, int64, error) {
	return s.supplierModule.List(req)
}

// ListSuppliersByStoreID 仅返回已与门店绑定的供应商（分页、关键词与 ListSuppliers 一致）。
func (s *SupplierService) ListSuppliersByStoreID(storeID uint, req *model.ListSupplierReq) ([]*model.Supplier, int64, error) {
	return s.supplierModule.ListByBoundStore(storeID, req)
}

func (s *SupplierService) UpdateSupplier(id uint, req *model.UpdateSupplierReq) error {
	_, err := s.supplierModule.GetByID(id)
	if err != nil {
		return errors.New("supplier not found")
	}
	return s.supplierModule.UpdateByID(id, req)
}

func (s *SupplierService) UpdateSupplierScoped(id, storeID uint, hqUnbound bool, req *model.UpdateSupplierReq) error {
	if _, err := s.GetSupplierScoped(id, storeID, hqUnbound); err != nil {
		return err
	}
	return s.supplierModule.UpdateByID(id, req)
}

func (s *SupplierService) DeleteSupplier(id uint) error {
	_, err := s.supplierModule.GetByID(id)
	if err != nil {
		return errors.New("supplier not found")
	}
	return s.supplierModule.Delete(id)
}

func (s *SupplierService) DeleteSupplierScoped(id, storeID uint, hqUnbound bool) error {
	if _, err := s.GetSupplierScoped(id, storeID, hqUnbound); err != nil {
		return err
	}
	return s.supplierModule.Delete(id)
}
