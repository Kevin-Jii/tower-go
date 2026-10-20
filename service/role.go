package service

import (
	"errors"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"gorm.io/gorm"
)

// ErrBuiltinRoleNotDeletable 尝试删除系统内置角色（super_admin / admin / store_admin / staff）
var ErrBuiltinRoleNotDeletable = apicode.New(apicode.BuiltinRoleNotDeletable)

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
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apicode.New(apicode.NotFound)
	}
	if err != nil {
		return nil, apicode.Wrap(apicode.InternalError, err)
	}
	InvalidateRolePermissionCache(id)
	return role, nil
}

// DeleteRole 删除角色
func DeleteRole(id uint) error {
	role, err := GetRole(id)
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
	role, err := module.GetRole(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apicode.New(apicode.NotFound)
	}
	if err != nil {
		return nil, apicode.Wrap(apicode.InternalError, err)
	}
	return role, nil
}

// ListRoles 获取角色列表
func ListRoles() ([]model.Role, error) {
	return module.ListRoles()
}

// ListRolesFiltered 过滤查询角色
func ListRolesFiltered(keyword string, status *int8) ([]model.Role, error) {
	return module.ListRolesFiltered(keyword, status)
}
