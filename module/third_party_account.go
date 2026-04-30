package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type ThirdPartyAccountModule struct {
	db *gorm.DB
}

func NewThirdPartyAccountModule(db *gorm.DB) *ThirdPartyAccountModule {
	return &ThirdPartyAccountModule{db: db}
}

func (m *ThirdPartyAccountModule) List(keyword string) ([]*model.ThirdPartyAccount, error) {
	var rows []*model.ThirdPartyAccount
	q := m.db.Model(&model.ThirdPartyAccount{})
	if keyword != "" {
		kw := "%" + keyword + "%"
		q = q.Where("name LIKE ? OR login_name LIKE ? OR phone LIKE ?", kw, kw, kw)
	}
	if err := q.Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (m *ThirdPartyAccountModule) GetByID(id uint) (*model.ThirdPartyAccount, error) {
	var row model.ThirdPartyAccount
	if err := m.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (m *ThirdPartyAccountModule) Create(row *model.ThirdPartyAccount) error {
	return m.db.Create(row).Error
}

func (m *ThirdPartyAccountModule) Update(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.ThirdPartyAccount{}).Where("id = ?", id).Updates(updates).Error
}

func (m *ThirdPartyAccountModule) Delete(id uint) error {
	return m.db.Delete(&model.ThirdPartyAccount{}, id).Error
}
