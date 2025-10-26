package model

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Username string `json:"username" gorm:"uniqueIndex;not null;type:varchar(191)"`
	Password string `json:"-" gorm:"not null"`

	// --- 权限管理 (RBAC) ---
	// RoleID: 关联到 'roles' 表的主键。这是 RBAC 的核心。
	RoleID uint `json:"role_id" gorm:"not null;default:1"` // default:1 可以是默认的“普通用户”角色ID

	// --- 状态与安全管理 ---
	// Status: 账号状态 (1=正常, 2=禁用, 3=未激活等)
	Status      int       `json:"status" gorm:"not null;default:1"`
	Nickname    string    `json:"nickname" gorm:"type:varchar(100)"` // 推荐限制长度
	Email       string    `json:"email" gorm:"type:varchar(255)"`
	Phone       string    `json:"phone" gorm:"uniqueIndex;type:varchar(20)"`
	LastLoginAt time.Time `json:"last_login_at"` // 记录最后登录时间
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateUserReq struct {
	Phone    string `json:"phone" binding:"required,len=11"`   // 手机号验证
	Password string `json:"password" binding:"required,min=6"` // 密码至少6位
	Username string `json:"username" binding:"required"`       // 强制要求非空
	Email    string `json:"email" binding:"omitempty,email"`   // 可选，但如果提供则必须是有效的邮箱
}

type UpdateUserReq struct {
	Password string `json:"password,omitempty" binding:"omitempty,min=6"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Phone    string `json:"phone,omitempty" binding:"omitempty,len=11"`
}
