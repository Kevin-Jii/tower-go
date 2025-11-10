package service

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

// CreateRole 创建角色
func CreateRole(req *model.Role) (*model.Role, error) {
	return module.CreateRole(req)
}

// UpdateRole 更新角色
func UpdateRole(id uint, req *model.UpdateRoleReq) (*model.Role, error) {
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

// ListRolesFiltered 过滤查询角色
func ListRolesFiltered(keyword string, status *int8) ([]model.Role, error) {
	return module.ListRolesFiltered(keyword, status)
}
