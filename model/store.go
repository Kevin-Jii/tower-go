package model

import "time"

// Store 门店表
type Store struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	StoreCode     *string   `json:"store_code,omitempty" gorm:"uniqueIndex;type:varchar(6);default:null"` // 门店编码 JWXXXX
	Name          string    `json:"name" gorm:"not null;type:varchar(100)"`                               // 门店名称
	Address       string    `json:"address" gorm:"type:varchar(255)"`                                     // 门店地址
	Phone         string    `json:"phone" gorm:"type:varchar(20)"`                                        // 联系电话
	BusinessHours string    `json:"business_hours" gorm:"type:varchar(100)"`                              // 营业时间
	Status        int       `json:"status" gorm:"not null;default:1"`                                     // 状态：1=正常，2=停业
	ContactPerson string    `json:"contact_person" gorm:"type:varchar(50)"`                               // 联系人
	Remark        string    `json:"remark" gorm:"type:text"`                                              // 备注
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CreateStoreReq 创建门店请求
type CreateStoreReq struct {
	Name          string `json:"name" binding:"required"`
	Address       string `json:"address"`
	Phone         string `json:"phone"`
	BusinessHours string `json:"business_hours"`
	ContactPerson string `json:"contact_person"`
	Remark        string `json:"remark"`
}

// UpdateStoreReq 更新门店请求
type UpdateStoreReq struct {
	Name          string `json:"name,omitempty"`                          // 名称不允许置空，保持非指针
	Address       string `json:"address,omitempty" patch:"always"`        // 允许清空地址
	Phone         string `json:"phone,omitempty" patch:"always"`          // 允许清空电话
	BusinessHours string `json:"business_hours,omitempty" patch:"always"` // 允许清空营业时间
	Status        *int   `json:"status,omitempty" patch:"allowZero"`      // 允许更新为0或2（指针已可）加 tag 标识意图
	ContactPerson string `json:"contact_person,omitempty" patch:"always"` // 允许清空联系人
	Remark        string `json:"remark,omitempty" patch:"always"`         // 允许清空备注
}
