package model

import "time"

// RoleMenu 角色菜单关联表（多对多）
type RoleMenu struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	RoleID    uint      `json:"role_id" gorm:"not null;index;comment:角色ID"`
	MenuID    uint      `json:"menu_id" gorm:"not null;index;comment:菜单ID"`
	CreatedAt time.Time `json:"created_at"`

	// 联合唯一索引
	Role *Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Menu *Menu `json:"menu,omitempty" gorm:"foreignKey:MenuID"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "role_menus"
}

// AssignMenusToRoleReq 为角色分配菜单请求
type AssignMenusToRoleReq struct {
	RoleID  uint   `json:"role_id" binding:"required"`
	MenuIDs []uint `json:"menu_ids" binding:"required"`
}
