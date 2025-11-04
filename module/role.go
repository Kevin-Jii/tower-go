package module

import (
	"gorm.io/gorm"
	"tower-go/model"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

// CreateRole 创建角色
func CreateRole(req *model.Role) (*model.Role, error) {
	if err := db.Create(req).Error; err != nil {
		return nil, err
	}
	return req, nil
}

// UpdateRole 更新角色
func UpdateRole(id uint, req *model.Role) (*model.Role, error) {
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		return nil, err
	}
	role.Name = req.Name
	role.Code = req.Code
	role.Description = req.Description
	if err := db.Save(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// DeleteRole 删除角色
func DeleteRole(id uint) error {
	return db.Delete(&model.Role{}, id).Error
}

// GetRole 获取单个角色
func GetRole(id uint) (*model.Role, error) {
	var role model.Role
	if err := db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRoles 获取角色列表
func ListRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
