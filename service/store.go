package service

import (
	"errors"
	"tower-go/model"
	"tower-go/module"
)

type StoreService struct {
	storeModule *module.StoreModule
}

func NewStoreService(storeModule *module.StoreModule) *StoreService {
	return &StoreService{storeModule: storeModule}
}

// CreateStore 创建门店（仅总部管理员）
func (s *StoreService) CreateStore(req *model.CreateStoreReq) error {
	store := &model.Store{
		Name:          req.Name,
		Address:       req.Address,
		Phone:         req.Phone,
		BusinessHours: req.BusinessHours,
		ContactPerson: req.ContactPerson,
		Remark:        req.Remark,
		Status:        1, // 默认正常
	}
	return s.storeModule.Create(store)
}

// GetStore 获取门店详情
func (s *StoreService) GetStore(id uint) (*model.Store, error) {
	return s.storeModule.GetByID(id)
}

// ListStores 获取门店列表
func (s *StoreService) ListStores(page, pageSize int) ([]*model.Store, int64, error) {
	return s.storeModule.List(page, pageSize)
}

// UpdateStore 更新门店信息
func (s *StoreService) UpdateStore(id uint, req *model.UpdateStoreReq) error {
	store, err := s.storeModule.GetByID(id)
	if err != nil {
		return errors.New("store not found")
	}

	if req.Name != "" {
		store.Name = req.Name
	}
	if req.Address != "" {
		store.Address = req.Address
	}
	if req.Phone != "" {
		store.Phone = req.Phone
	}
	if req.BusinessHours != "" {
		store.BusinessHours = req.BusinessHours
	}
	if req.Status != nil {
		store.Status = *req.Status
	}
	if req.ContactPerson != "" {
		store.ContactPerson = req.ContactPerson
	}
	if req.Remark != "" {
		store.Remark = req.Remark
	}

	return s.storeModule.Update(store)
}

// DeleteStore 删除门店
func (s *StoreService) DeleteStore(id uint) error {
	return s.storeModule.Delete(id)
}
