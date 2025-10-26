package module

import (
	"tower-go/model"

	"gorm.io/gorm"
)

type UserModule struct {
	db *gorm.DB
}

func NewUserModule(db *gorm.DB) *UserModule {
	return &UserModule{db: db}
}

func (m *UserModule) Create(user *model.User) error {
	return m.db.Create(user).Error
}

func (m *UserModule) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := m.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModule) GetByPhone(phone string) (*model.User, error) {
	var user model.User
	if err := m.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModule) ExistsByPhone(phone string) (bool, error) {
	var count int64
	if err := m.db.Model(&model.User{}).Where("phone = ?", phone).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (m *UserModule) List(page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := m.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := m.db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (m *UserModule) Update(user *model.User) error {
	return m.db.Save(user).Error
}

func (m *UserModule) Delete(id uint) error {
	return m.db.Delete(&model.User{}, id).Error
}
