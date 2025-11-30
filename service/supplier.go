package service

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type SupplierService struct {
	supplierModule *module.SupplierModule
}

func NewSupplierService(supplierModule *module.SupplierModule) *SupplierService {
	return &SupplierService{supplierModule: supplierModule}
}

func (s *SupplierService) CreateSupplier(req *model.CreateSupplierReq) error {
	exists, err := s.supplierModule.ExistsByCode(req.SupplierCode)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("supplier code already exists")
	}

	supplier := &model.Supplier{
		SupplierCode:    req.SupplierCode,
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

	if req.SupplierCode != "" {
		existing, err := s.supplierModule.GetByCode(req.SupplierCode)
		if err == nil && existing.ID != id {
			return errors.New("supplier code already exists")
		}
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
