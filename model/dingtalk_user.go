package model

import "time"

// DingTalkUser 钉钉用户缓存表
type DingTalkUser struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Mobile    string    `json:"mobile" gorm:"uniqueIndex;type:varchar(20);not null"` // 手机号
	UserID    string    `json:"user_id" gorm:"type:varchar(100);not null"`           // 钉钉用户ID (staffId)
	Name      string    `json:"name" gorm:"type:varchar(100)"`                       // 用户姓名
	UnionID   string    `json:"union_id" gorm:"type:varchar(100)"`                   // 钉钉unionId
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
