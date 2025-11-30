package model

import "time"

// SupplierCategory 供应商分类表
type SupplierCategory struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierID uint      `json:"supplier_id" gorm:"not null;index;comment:供应商ID"`
	Supplier   *Supplier `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	Name       string    `json:"name" gorm:"type:varchar(100);not null;comment:分类名称"`
	Sort       int       `json:"sort" gorm:"default:0;comment:排序"`
	Status     int8      `json:"status" gorm:"not null;default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (SupplierCategory) TableName() string {
	return "supplier_categories"
}

// CreateSupplierCategoryReq 创建供应商分类请求
type CreateSupplierCategoryReq struct {
	SupplierID uint   `json:"supplier_id" binding:"required"`
	Name       string `json:"name" binding:"required,max=100"`
	Sort       int    `json:"sort"`
}

// UpdateSupplierCategoryReq 更新供应商分类请求
type UpdateSupplierCategoryReq struct {
	Name   string `json:"name,omitempty" binding:"max=100"`
	Sort   *int   `json:"sort,omitempty"`
	Status *int8  `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
}
