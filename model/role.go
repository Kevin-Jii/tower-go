package model

import "time"

// Role 角色表
type Role struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;type:varchar(50)"` // 角色名称
	Code        string    `json:"code" gorm:"uniqueIndex;not null;type:varchar(50)"` // 角色代码：admin(总部管理员), store_admin(门店管理员), staff(普通员工)
	DataScope   int8      `json:"data_scope" gorm:"type:tinyint;not null;default:3;comment:数据范围 1=全部 2=租户 3=门店 4=仅本人"`
	Status      int8      `json:"status" gorm:"type:tinyint(1);not null;default:1"` // 状态：1=启用 0=禁用
	Description string    `json:"description" gorm:"type:varchar(255)"`             // 角色描述
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateRoleReq 局部更新请求（仅更新出现的字段）
type UpdateRoleReq struct {
	Name        *string `json:"name"`                       // 角色名称（可选）
	Code        *string `json:"code"`                       // 角色代码（可选）
	DataScope   *int8   `json:"data_scope"`                 // 数据范围（可选）
	Status      *int8   `json:"status" patch:"allowZero"`   // 允许更新为0（禁用）
	Description *string `json:"description" patch:"always"` // 允许清空描述
}

const (
	RoleCodeSuperAdmin = "super_admin" // 超级管理员
	RoleCodeAdmin      = "admin"       // 总部管理员
	RoleCodeStoreAdmin = "store_admin" // 门店管理员
	RoleCodeStaff      = "staff"       // 普通员工
)

// IsSuperAdminRole 超级管理员：不绑定门店，始终拥有全库数据与跨店操作能力。
func IsSuperAdminRole(roleCode string) bool {
	return roleCode == RoleCodeSuperAdmin
}

// HQUnboundAdminRole 判定「可跨店操作的总部身份」。
//
// - super_admin：始终视为未绑店总部，不校验 store_id。
// - admin：须 store_id==0 才可跨店；store_id>0 时与门店管理员一致，仅本店。
func HQUnboundAdminRole(roleCode string, storeID uint) bool {
	if IsSuperAdminRole(roleCode) {
		return true
	}
	if storeID > 0 {
		return false
	}
	return roleCode == RoleCodeAdmin
}

// IsBuiltinRoleNonDeletable 系统内置角色，禁止删除
func IsBuiltinRoleNonDeletable(code string) bool {
	switch code {
	case RoleCodeSuperAdmin, RoleCodeAdmin, RoleCodeStoreAdmin, RoleCodeStaff:
		return true
	default:
		return false
	}
}

// 数据范围（与角色 data_scope 一致）
const (
	DataScopeAll    int8 = 1 // 全部（总部）
	DataScopeTenant int8 = 2 // 租户/公司（预留）
	DataScopeStore  int8 = 3 // 本门店
	DataScopeSelf   int8 = 4 // 仅本人
)
