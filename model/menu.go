package model

import "time"

// Menu 菜单表（后台管理系统的导航菜单）
type Menu struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	ParentID   uint      `json:"parent_id" gorm:"default:0;index;comment:父菜单ID，0表示顶级菜单"`
	Name       string    `json:"name" gorm:"type:varchar(50);not null;comment:菜单名称"`
	Title      string    `json:"title" gorm:"type:varchar(50);not null;comment:菜单标题（显示用）"`
	Icon       string    `json:"icon" gorm:"type:varchar(100);comment:菜单图标"`
	Path       string    `json:"path" gorm:"type:varchar(200);comment:路由路径"`
	Component  string    `json:"component" gorm:"type:varchar(200);comment:组件路径"`
	Type       int       `json:"type" gorm:"default:1;comment:菜单类型：1=目录，2=菜单，3=按钮"`
	Sort       int       `json:"sort" gorm:"default:0;comment:排序"`
	Permission string    `json:"permission" gorm:"type:varchar(200);comment:权限标识"`
	Visible    int       `json:"visible" gorm:"default:1;comment:是否可见：0=隐藏，1=显示"`
	Status     int       `json:"status" gorm:"default:1;comment:状态：0=禁用，1=启用"`
	Remark     string    `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// 关联
	Children    []*Menu `json:"children,omitempty" gorm:"-"`    // 子菜单（不存数据库）
	Permissions uint8   `json:"permissions,omitempty" gorm:"-"` // 用户对该菜单的权限（运行时填充）
}

// CreateMenuReq 创建菜单请求
type CreateMenuReq struct {
	ParentID   uint   `json:"parent_id"`
	Name       string `json:"name" binding:"required"`
	Title      string `json:"title" binding:"required"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Type       int    `json:"type" binding:"required,oneof=1 2 3"`
	Sort       int    `json:"sort"`
	Permission string `json:"permission"`
	Visible    int    `json:"visible" binding:"oneof=0 1"`
	Status     int    `json:"status" binding:"oneof=0 1"`
	Remark     string `json:"remark"`
}

// UpdateMenuReq 更新菜单请求
type UpdateMenuReq struct {
	ParentID   *uint  `json:"parent_id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Path       string `json:"path"`
	Component  string `json:"component"`
	Type       *int   `json:"type" binding:"omitempty,oneof=1 2 3"`
	Sort       *int   `json:"sort"`
	Permission string `json:"permission"`
	Visible    *int   `json:"visible" binding:"omitempty,oneof=0 1"`
	Status     *int   `json:"status" binding:"omitempty,oneof=0 1"`
	Remark     string `json:"remark"`
}
