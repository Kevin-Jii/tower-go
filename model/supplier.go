package model

import (
	"time"

	"gorm.io/gorm"
)

// Supplier 供应商表
type Supplier struct {
	ID              uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierCode    string         `json:"supplier_code" gorm:"uniqueIndex;type:varchar(50);not null;comment:供应商编码"`
	SupplierName    string         `json:"supplier_name" gorm:"type:varchar(200);not null;comment:供应商名称"`
	ContactPerson   string         `json:"contact_person" gorm:"type:varchar(100);comment:联系人"`
	ContactPhone    string         `json:"contact_phone" gorm:"type:varchar(20);comment:联系电话"`
	ContactEmail    string         `json:"contact_email" gorm:"type:varchar(100);comment:联系邮箱"`
	SupplierAddress string         `json:"supplier_address" gorm:"type:varchar(500);comment:供应商地址"`
	Remark          string         `json:"remark" gorm:"type:text;comment:备注"`
	Status          int8           `json:"status" gorm:"not null;default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (Supplier) TableName() string {
	return "suppliers"
}

// CreateSupplierReq 创建供应商请求
type CreateSupplierReq struct {
	SupplierName    string `json:"supplier_name" binding:"required,max=200"`         // 供应商名称
	ContactPerson   string `json:"contact_person" binding:"max=100"`                 // 联系人
	ContactPhone    string `json:"contact_phone" binding:"max=20"`                   // 联系电话
	ContactEmail    string `json:"contact_email" binding:"omitempty,max=100,email"`  // 联系邮箱
	SupplierAddress string `json:"supplier_address" binding:"max=500"`               // 地址
	Remark          string `json:"remark"`                                           // 备注
}

// UpdateSupplierReq 更新供应商请求
type UpdateSupplierReq struct {
	SupplierName    *string `json:"supplier_name,omitempty" binding:"omitempty,max=200"`
	ContactPerson   *string `json:"contact_person,omitempty" binding:"omitempty,max=100"`
	ContactPhone    *string `json:"contact_phone,omitempty" binding:"omitempty,max=20"`
	ContactEmail    *string `json:"contact_email,omitempty" binding:"omitempty,max=100,email"`
	SupplierAddress *string `json:"supplier_address,omitempty" binding:"omitempty,max=500"`
	Remark          *string `json:"remark,omitempty"`
	Status          *int8   `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
}

// ListSupplierReq 供应商列表查询请求
type ListSupplierReq struct {
	Keyword  string `form:"keyword"`
	Status   *int8  `form:"status"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
}
