package model

import "time"

// RoleMenu 角色菜单关联表（多对多）
type RoleMenu struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	RoleID      uint      `json:"role_id" gorm:"not null;index;comment:角色ID"`
	MenuID      uint      `json:"menu_id" gorm:"not null;index;comment:菜单ID"`
	Permissions uint8     `json:"permissions" gorm:"not null;default:0;comment:权限位：bit0=查看,bit1=新增,bit2=修改,bit3=删除"`
	CreatedAt   time.Time `json:"created_at"`

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
	RoleID  uint                `json:"role_id" binding:"required"`
	MenuIDs []uint              `json:"menu_ids" binding:"required"`
	Perms   map[uint]uint8      `json:"perms"` // 菜单ID -> 权限位映射
}

// MenuPermission 菜单权限位常量
const (
	PermView   uint8 = 1 << 3 // 1000 = 8  查看权限
	PermCreate uint8 = 1 << 2 // 0100 = 4  新增权限
	PermUpdate uint8 = 1 << 1 // 0010 = 2  修改权限
	PermDelete uint8 = 1 << 0 // 0001 = 1  删除权限
	PermAll    uint8 = 15      // 1111 = 15 所有权限
)

// HasPermission 检查是否有指定权限
func HasPermission(perms uint8, perm uint8) bool {
	return perms&perm == perm
}
