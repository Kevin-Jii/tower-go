package module

import (
	"tower-go/model"

	"gorm.io/gorm"
)

type StoreModule struct {
	db *gorm.DB
}

func NewStoreModule(db *gorm.DB) *StoreModule {
	return &StoreModule{db: db}
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

// List 获取门店列表（支持分页）
func (m *StoreModule) List(page, pageSize int) ([]*model.Store, int64, error) {
	var stores []*model.Store
	var total int64

	if err := m.db.Model(&model.Store{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := m.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&stores).Error; err != nil {
		return nil, 0, err
	}

	return stores, total, nil
}

// Update 更新门店信息
func (m *StoreModule) Update(store *model.Store) error {
	return m.db.Save(store).Error
}

// Delete 删除门店
func (m *StoreModule) Delete(id uint) error {
	return m.db.Delete(&model.Store{}, id).Error
}
