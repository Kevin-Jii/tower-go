package model

import "time"

// Role 角色表
type Role struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;type:varchar(50)"` // 角色名称
	Code        string    `json:"code" gorm:"uniqueIndex;not null;type:varchar(50)"` // 角色代码：admin(总部管理员), store_admin(门店管理员), staff(普通员工)
	Status      int8      `json:"status" gorm:"type:tinyint(1);not null;default:1"`  // 状态：1=启用 0=禁用
	Description string    `json:"description" gorm:"type:varchar(255)"`              // 角色描述
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateRoleReq 局部更新请求（仅更新出现的字段）
type UpdateRoleReq struct {
	Name        *string `json:"name"`                       // 角色名称（可选）
	Code        *string `json:"code"`                       // 角色代码（可选）
	Status      *int8   `json:"status" patch:"allowZero"`   // 允许更新为0（禁用）
	Description *string `json:"description" patch:"always"` // 允许清空描述
}

const (
	RoleCodeAdmin      = "admin"       // 总部管理员
	RoleCodeStoreAdmin = "store_admin" // 门店管理员
	RoleCodeStaff      = "staff"       // 普通员工
)
