package service

import (
	"errors"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils"
)

type StoreService struct {
	storeModule *module.StoreModule
}

func NewStoreService(storeModule *module.StoreModule) *StoreService {
	return &StoreService{storeModule: storeModule}
}

// CreateStore 创建门店（仅总部管理员）
func (s *StoreService) CreateStore(req *model.CreateStoreReq) error {
	storeCode, err := utils.GenerateStoreCode(s.storeModule.GetDB())
	if err != nil {
		return err
	}

	store := &model.Store{
		StoreCode:     &storeCode,
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

// ListStores 获取门店列表（全部数据）
func (s *StoreService) ListStores() ([]*model.Store, int64, error) {
	return s.storeModule.List()
}

// UpdateStore 更新门店信息
func (s *StoreService) UpdateStore(id uint, req *model.UpdateStoreReq) error {
	// 确认门店存在
	_, err := s.storeModule.GetByID(id)
	if err != nil {
		return errors.New("store not found")
	}

	// 使用动态更新避免整行覆盖
	return s.storeModule.UpdateByID(id, req)
}

// DeleteStore 删除门店
func (s *StoreService) DeleteStore(id uint) error {
	return s.storeModule.Delete(id)
}
