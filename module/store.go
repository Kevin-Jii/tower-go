package module

import (
	"tower-go/model"
	updatesPkg "tower-go/utils/updates"

	"gorm.io/gorm"
)

type StoreModule struct {
	db *gorm.DB
}

func NewStoreModule(db *gorm.DB) *StoreModule {
	return &StoreModule{db: db}
}

// GetDB 返回底层数据库实例
func (m *StoreModule) GetDB() *gorm.DB {
	return m.db
}

// Create 创建门店
func (m *StoreModule) Create(store *model.Store) error {
	return m.db.Create(store).Error
}

// GetByID 根据ID获取门店
func (m *StoreModule) GetByID(id uint) (*model.Store, error) {
	var store model.Store
	if err := m.db.First(&store, id).Error; err != nil {
		return nil, err
	}
	return &store, nil
}

// List 获取门店列表（全部数据，不分页）
func (m *StoreModule) List() ([]*model.Store, int64, error) {
	var stores []*model.Store
	var total int64

	if err := m.db.Model(&model.Store{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := m.db.Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

// Update 更新门店信息
func (m *StoreModule) Update(store *model.Store) error {
	return m.db.Save(store).Error
}

// UpdateByID 根据ID更新门店信息（动态更新，避免整行覆盖）
func (m *StoreModule) UpdateByID(id uint, req *model.UpdateStoreReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return nil
	}
	return m.db.Model(&model.Store{}).Where("id = ?", id).Updates(updateMap).Error
}

// Delete 删除门店
func (m *StoreModule) Delete(id uint) error {
	return m.db.Delete(&model.Store{}, id).Error
}
