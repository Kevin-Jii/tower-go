package model

import "time"

type User struct {
	ID       uint   `json:"id" gorm:"primarykey"`
	Password string `json:"-" gorm:"not null"`
	Phone    string `json:"phone" gorm:"uniqueIndex;type:varchar(20)"`
	Username string `json:"username" gorm:"not null;uniqueIndex:idx_store_username;type:varchar(191)"`
	Nickname string `json:"nickname" gorm:"type:varchar(100)"`
	Email    string `json:"email" gorm:"type:varchar(255)"`

	// --- 门店关联 ---
	StoreID uint   `json:"store_id" gorm:"not null;uniqueIndex:idx_store_username"`
	Store   *Store `json:"store,omitempty" gorm:"foreignKey:StoreID"` // 门店关联

	// --- 权限管理 (RBAC) ---
	RoleID uint  `json:"role_id" gorm:"not null;default:3"`       // 角色ID：1=总部管理员，2=门店管理员，3=普通员工
	Role   *Role `json:"role,omitempty" gorm:"foreignKey:RoleID"` // 角色关联

	// --- 状态与安全管理 ---
	Status      int       `json:"status" gorm:"not null;default:1"` // 账号状态：1=正常，2=禁用
	LastLoginAt time.Time `json:"last_login_at"`                    // 记录最后登录时间
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
