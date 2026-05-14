package service

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

// ErrBuiltinRoleNotDeletable 尝试删除系统内置角色（super_admin / admin / store_admin / staff）
var ErrBuiltinRoleNotDeletable = errors.New("builtin role not deletable")

// CreateRole 创建角色
func CreateRole(req *model.Role) (*model.Role, error) {
	role, err := module.CreateRole(req)
	if err == nil {
		InvalidateAllPermissionCache()
	}
	return role, err
}

// UpdateRole 更新角色
func UpdateRole(id uint, req *model.UpdateRoleReq) (*model.Role, error) {
	role, err := module.UpdateRole(id, req)
	if err == nil {
		InvalidateRolePermissionCache(id)
	}
	return role, err
}

// DeleteRole 删除角色
func DeleteRole(id uint) error {
	role, err := module.GetRole(id)
	if err != nil {
		return err
	}
	if model.IsBuiltinRoleNonDeletable(role.Code) {
		return ErrBuiltinRoleNotDeletable
	}
	err = module.DeleteRole(id)
	if err == nil {
		InvalidateRolePermissionCache(id)
	}
	return err
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
