package module

import (
	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type DingTalkUserModule struct {
	db *gorm.DB
}

func NewDingTalkUserModule(db *gorm.DB) *DingTalkUserModule {
	return &DingTalkUserModule{db: db}
}

// GetByMobile 通过手机号获取钉钉用户
func (m *DingTalkUserModule) GetByMobile(mobile string) (*model.DingTalkUser, error) {
	var user model.DingTalkUser
	if err := m.db.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建钉钉用户缓存
func (m *DingTalkUserModule) Create(user *model.DingTalkUser) error {
	return m.db.Create(user).Error
}

// Upsert 创建或更新钉钉用户（根据手机号）
func (m *DingTalkUserModule) Upsert(user *model.DingTalkUser) error {
	return m.db.Where("mobile = ?", user.Mobile).
		Assign(model.DingTalkUser{
			UserID:  user.UserID,
			Name:    user.Name,
			UnionID: user.UnionID,
		}).
		FirstOrCreate(user).Error
}
