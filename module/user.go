package module

import (
	"log"
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

	if err := m.db.Where("store_id = ?", storeID).
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
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR phone LIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// UpdateByID 根据ID更新用户信息
func (m *UserModule) UpdateByID(id uint, req *model.UpdateUserReq) error {
	updates := make(map[string]interface{})

	if req.Username != "" {
		updates["username"] = req.Username
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Password != "" {
		updates["password"] = req.Password
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Gender != nil { // 支持性别更新 1男2女
		updates["gender"] = *req.Gender
	}

	log.Printf("[UserModule.UpdateByID] id=%d updates=%v", id, updates)
	return m.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteByUserIDAndStoreID 根据用户ID和门店ID删除用户
func (m *UserModule) DeleteByUserIDAndStoreID(userID uint, storeID uint) error {
	return m.db.Where("id = ? AND store_id = ?", userID, storeID).Delete(&model.User{}).Error
}

// UpdatePasswordByID 仅更新密码字段
func (m *UserModule) UpdatePasswordByID(id uint, hashed string) error {
	return m.db.Model(&model.User{}).Where("id = ?", id).Update("password", hashed).Error
}
