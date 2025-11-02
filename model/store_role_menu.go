package model

import "time"

// StoreRoleMenu 门店角色菜单权限表（实现门店级别的菜单权限定制）
type StoreRoleMenu struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	StoreID   uint      `json:"store_id" gorm:"not null;index;comment:门店ID"`
	RoleID    uint      `json:"role_id" gorm:"not null;index;comment:角色ID"`
	MenuID    uint      `json:"menu_id" gorm:"not null;index;comment:菜单ID"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	Store *Store `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Role  *Role  `json:"role,omitempty" gorm:"foreignKey:RoleID"`
	Menu  *Menu  `json:"menu,omitempty" gorm:"foreignKey:MenuID"`
}

// TableName 指定表名
func (StoreRoleMenu) TableName() string {
	return "store_role_menus"
}

// AssignStoreMenusReq 为门店角色分配菜单权限请求
type AssignStoreMenusReq struct {
	StoreID uint   `json:"store_id" binding:"required"`
	RoleID  uint   `json:"role_id" binding:"required"`
	MenuIDs []uint `json:"menu_ids" binding:"required"`
}

// CopyStoreMenusReq 复制门店菜单权限请求
type CopyStoreMenusReq struct {
	FromStoreID uint `json:"from_store_id" binding:"required"`
	ToStoreID   uint `json:"to_store_id" binding:"required"`
	RoleID      uint `json:"role_id" binding:"required"`
}
