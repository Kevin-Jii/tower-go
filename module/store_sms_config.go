package module

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type StoreSmsConfigModule struct{ db *gorm.DB }

func NewStoreSmsConfigModule(db *gorm.DB) *StoreSmsConfigModule {
	return &StoreSmsConfigModule{db: db}
}

func (m *StoreSmsConfigModule) GetByStoreID(storeID uint) (*model.StoreSmsConfig, error) {
	if storeID == 0 {
		return nil, errors.New("门店 ID 不能为空")
	}
	var row model.StoreSmsConfig
	if err := m.db.Where("store_id = ?", storeID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *StoreSmsConfigModule) Upsert(row *model.StoreSmsConfig) error {
	return m.db.Save(row).Error
}

func (m *StoreSmsConfigModule) DeleteByStoreID(storeID uint) error {
	return m.db.Where("store_id = ?", storeID).Delete(&model.StoreSmsConfig{}).Error
}
