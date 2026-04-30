package model

import "time"

// Store 门店表
type Store struct {
	ID                  uint      `json:"id" gorm:"primarykey"`
	StoreCode           *string   `json:"store_code,omitempty" gorm:"uniqueIndex;type:varchar(6);default:null"` // 门店编码 JWXXXX
	Name                string    `json:"name" gorm:"not null;type:varchar(100)"`                               // 门店名称
	Address             string    `json:"address" gorm:"type:varchar(255)"`                                     // 门店地址
	AdministrativeUnit  string    `json:"administrative_unit" gorm:"type:varchar(100);comment:归属区"`
	Phone               string    `json:"phone" gorm:"type:varchar(20)"`           // 联系电话
	BusinessHours       string    `json:"business_hours" gorm:"type:varchar(100)"` // 营业时间
	Status              int       `json:"status" gorm:"not null;default:1"`        // 状态：1=正常，2=停业
	ContactPerson       string    `json:"contact_person" gorm:"type:varchar(50)"`  // 联系人
	Remark              string    `json:"remark" gorm:"type:text"`                 // 备注
	ThirdPartyAccountID *uint     `json:"third_party_account_id,omitempty" gorm:"index;comment:绑定第三方账号池ID"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	ThirdPartyAccount *ThirdPartyAccount `json:"third_party_account,omitempty" gorm:"foreignKey:ThirdPartyAccountID"`
}

// CreateStoreReq 创建门店请求
type CreateStoreReq struct {
	Name               string `json:"name" binding:"required"`
	Address            string `json:"address"`
	AdministrativeUnit string `json:"administrative_unit"`
	Phone              string `json:"phone"`
	BusinessHours      string `json:"business_hours"`
	ContactPerson      string `json:"contact_person"`
	Remark             string `json:"remark"`
}

// UpdateStoreReq 更新门店请求（使用指针类型实现部分更新）
type UpdateStoreReq struct {
	Name               *string `json:"name,omitempty"`                // 门店名称
	Address            *string `json:"address,omitempty"`             // 门店地址
	AdministrativeUnit *string `json:"administrative_unit,omitempty"` // 归属区
	Phone              *string `json:"phone,omitempty"`               // 联系电话
	BusinessHours      *string `json:"business_hours,omitempty"`      // 营业时间
	Status             *int    `json:"status,omitempty"`              // 状态：1=正常，2=停业
	ContactPerson      *string `json:"contact_person,omitempty"`      // 联系人
	Remark             *string `json:"remark,omitempty"`              // 备注
}

// BindStoreThirdPartyAccountReq 绑定门店第三方账号（传 null 表示解绑）
type BindStoreThirdPartyAccountReq struct {
	ThirdPartyAccountID *uint `json:"third_party_account_id"`
}
