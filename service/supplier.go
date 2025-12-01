package service

import (
	"errors"
	"fmt"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type SupplierService struct {
	supplierModule *module.SupplierModule
}

func NewSupplierService(supplierModule *module.SupplierModule) *SupplierService {
	return &SupplierService{supplierModule: supplierModule}
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
	return s.supplierModule.Create(supplier)
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

func (s *SupplierService) ListSuppliers(req *model.ListSupplierReq) ([]*model.Supplier, int64, error) {
	return s.supplierModule.List(req)
}

func (s *SupplierService) UpdateSupplier(id uint, req *model.UpdateSupplierReq) error {
	_, err := s.supplierModule.GetByID(id)
	if err != nil {
		return errors.New("supplier not found")
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
