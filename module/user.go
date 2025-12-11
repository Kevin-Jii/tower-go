package module

import (
	"log"

	"github.com/Kevin-Jii/tower-go/model"
	searchPkg "github.com/Kevin-Jii/tower-go/utils/search"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"

	"gorm.io/gorm"
)

type UserModule struct {
	db *gorm.DB
}

func NewUserModule(db *gorm.DB) *UserModule {
	return &UserModule{db: db}
}

// GetDB 返回数据库连接实例
func (m *UserModule) GetDB() *gorm.DB {
	return m.db
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
	// 预加载 Role 和 Store 关联信息
	if err := m.db.Preload("Role").Preload("Store").Where("phone = ?", phone).First(&user).Error; err != nil {
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

// GetByUserIDAndStoreID 根据用户ID和门店ID获取用户
func (m *UserModule) GetByUserIDAndStoreID(userID uint, storeID uint) (*model.User, error) {
	var user model.User
	if err := m.db.Where("id = ? AND store_id = ?", userID, storeID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// ListByStoreID 根据门店ID获取用户列表
func (m *UserModule) ListByStoreID(storeID uint, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	if err := m.db.Model(&model.User{}).Where("store_id = ?", storeID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := m.db.Preload("Role").Where("store_id = ?", storeID).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// ListByStoreIDWithKeyword 根据门店ID与关键字(匹配 username 或 phone)获取用户列表
func (m *UserModule) ListByStoreIDWithKeyword(storeID uint, keyword string, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := m.db.Model(&model.User{}).Where("store_id = ?", storeID)

	if keyword != "" {
		// 使用优化的搜索条件
		conditions := searchPkg.OptimizeSearchKeyword(keyword)
		if len(conditions) > 0 {
			sql, args := searchPkg.BuildSearchSQL(conditions)
			query = query.Where(sql, args...)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Role").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// ListAllUsers 获取所有用户（支持分页，用于总部管理员）
func (m *UserModule) ListAllUsers(keyword string, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	query := m.db.Model(&model.User{})

	if keyword != "" {
		// 使用优化的搜索条件
		conditions := searchPkg.OptimizeSearchKeyword(keyword)
		if len(conditions) > 0 {
			sql, args := searchPkg.BuildSearchSQL(conditions)
			query = query.Where(sql, args...)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("Role").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// UpdateByID 根据ID更新用户信息
func (m *UserModule) UpdateByID(id uint, req *model.UpdateUserReq) error {
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		log.Printf("[UserModule.UpdateByID] id=%d no fields to update", id)
		return nil
	}
	log.Printf("[UserModule.UpdateByID] id=%d updates=%v", id, updateMap)
	return m.db.Model(&model.User{}).Where("id = ?", id).Updates(updateMap).Error
}

// DeleteByUserIDAndStoreID 根据用户ID和门店ID删除用户
func (m *UserModule) DeleteByUserIDAndStoreID(userID uint, storeID uint) error {
	return m.db.Where("id = ? AND store_id = ?", userID, storeID).Delete(&model.User{}).Error
}

// UpdatePasswordByID 仅更新密码字段
func (m *UserModule) UpdatePasswordByID(id uint, hashed string) error {
	return m.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashed).Error
}

// GetByDingTalkID 根据钉钉用户ID获取用户（通过手机号关联）
func (m *UserModule) GetByDingTalkID(dingTalkID string) (*model.User, error) {
	// 先从钉钉用户缓存表查找手机号
	var dingUser model.DingTalkUser
	if err := m.db.Where("user_id = ?", dingTalkID).First(&dingUser).Error; err != nil {
		return nil, err
	}

	// 再通过手机号查找系统用户
	var user model.User
	if err := m.db.Where("phone = ?", dingUser.Mobile).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
