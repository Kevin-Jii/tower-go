package model

import (
	"time"
)

// Supplier 供应商表
type Supplier struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	SupplierCode   string    `json:"supplier_code" gorm:"uniqueIndex;type:varchar(50);not null;comment:供应商编码"`
	SupplierName   string    `json:"supplier_name" gorm:"type:varchar(200);not null;comment:供应商名称"`
	ContactPerson  string    `json:"contact_person" gorm:"type:varchar(100);comment:联系人"`
	ContactPhone   string    `json:"contact_phone" gorm:"type:varchar(20);comment:联系电话"`
	ContactEmail   string    `json:"contact_email" gorm:"type:varchar(100);comment:联系邮箱"`
	SupplierAddress string   `json:"supplier_address" gorm:"type:varchar(500);comment:供应商地址"`
	Remark         string    `json:"remark" gorm:"type:text;comment:备注"`
	Status         int8      `json:"status" gorm:"not null;default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:datetime(3);comment:创建时间"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:datetime(3);comment:更新时间"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty" gorm:"type:datetime(3);index;comment:软删除时间"`
}

// CreateSupplierReq 创建供应商请求
type CreateSupplierReq struct {
	SupplierCode   string `json:"supplier_code" binding:"required,max=50"`
	SupplierName   string `json:"supplier_name" binding:"required,max=200"`
	ContactPerson  string `json:"contact_person" binding:"max=100"`
	ContactPhone   string `json:"contact_phone" binding:"max=20"`
	ContactEmail   string `json:"contact_email" binding:"max=100,email"`
	SupplierAddress string `json:"supplier_address" binding:"max=500"`
	Remark         string `json:"remark"`
	Status         int8   `json:"status" binding:"oneof=0 1"`
}

// UpdateSupplierReq 更新供应商请求
type UpdateSupplierReq struct {
	ID             uint   `json:"id" binding:"required"`
	SupplierCode   string `json:"supplier_code" binding:"max=50"`
	SupplierName   string `json:"supplier_name" binding:"max=200"`
	ContactPerson  string `json:"contact_person" binding:"max=100"`
	ContactPhone   string `json:"contact_phone" binding:"max=20"`
	ContactEmail   string `json:"contact_email" binding:"max=100,email"`
	SupplierAddress string `json:"supplier_address" binding:"max=500"`
	Remark         string `json:"remark"`
	Status         int8   `json:"status" binding:"oneof=0 1"`
}

// TableName 指定表名
func (Supplier) TableName() string {
	return "suppliers"
}
