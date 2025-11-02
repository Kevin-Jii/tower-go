package model

import "time"

// Role 角色表
type Role struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null;type:varchar(50)"` // 角色名称
	Code        string    `json:"code" gorm:"uniqueIndex;not null;type:varchar(50)"` // 角色代码：admin(总部管理员), store_admin(门店管理员), staff(普通员工)
	Description string    `json:"description" gorm:"type:varchar(255)"`              // 角色描述
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const (
	RoleCodeAdmin      = "admin"       // 总部管理员
	RoleCodeStoreAdmin = "store_admin" // 门店管理员
	RoleCodeStaff      = "staff"       // 普通员工
)
