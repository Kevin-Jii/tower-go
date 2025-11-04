package service

import (
	"tower-go/model"
	"tower-go/module"
)

// CreateRole 创建角色
func CreateRole(req *model.Role) (*model.Role, error) {
	return module.CreateRole(req)
}

// UpdateRole 更新角色
func UpdateRole(id uint, req *model.Role) (*model.Role, error) {
	return module.UpdateRole(id, req)
}

// DeleteRole 删除角色
func DeleteRole(id uint) error {
	return module.DeleteRole(id)
}

// GetRole 获取单个角色
func GetRole(id uint) (*model.Role, error) {
	return module.GetRole(id)
}

// ListRoles 获取角色列表
func ListRoles() ([]model.Role, error) {
	return module.ListRoles()
}
