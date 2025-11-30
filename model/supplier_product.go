package model

import "time"

// SupplierProduct 供应商商品表
type SupplierProduct struct {
	ID         uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierID uint              `json:"supplier_id" gorm:"not null;index;comment:供应商ID"`
	Supplier   *Supplier         `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	CategoryID uint              `json:"category_id" gorm:"not null;index;comment:分类ID"`
	Category   *SupplierCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Name       string            `json:"name" gorm:"type:varchar(200);not null;comment:商品名称"`
	Unit       string            `json:"unit" gorm:"type:varchar(20);not null;default:'斤';comment:单位"`
	Price      float64           `json:"price" gorm:"type:decimal(10,2);not null;default:0;comment:单价"`
	Spec       string            `json:"spec" gorm:"type:varchar(100);comment:规格"`
	Remark     string            `json:"remark" gorm:"type:varchar(500);comment:备注"`
	Status     int8              `json:"status" gorm:"not null;default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}

func (SupplierProduct) TableName() string {
	return "supplier_products"
}

// CreateSupplierProductReq 创建供应商商品请求
type CreateSupplierProductReq struct {
	SupplierID uint    `json:"supplier_id" binding:"required"`
	CategoryID uint    `json:"category_id" binding:"required"`
	Name       string  `json:"name" binding:"required,max=200"`
	Unit       string  `json:"unit" binding:"required,max=20"`
	Price      float64 `json:"price" binding:"gte=0"`
	Spec       string  `json:"spec" binding:"max=100"`
	Remark     string  `json:"remark" binding:"max=500"`
}

// UpdateSupplierProductReq 更新供应商商品请求
type UpdateSupplierProductReq struct {
	CategoryID *uint    `json:"category_id,omitempty"`
	Name       string   `json:"name,omitempty" binding:"max=200"`
	Unit       string   `json:"unit,omitempty" binding:"max=20"`
	Price      *float64 `json:"price,omitempty" binding:"omitempty,gte=0"`
	Spec       string   `json:"spec,omitempty" binding:"max=100"`
	Remark     string   `json:"remark,omitempty" binding:"max=500"`
	Status     *int8    `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
}

// ListSupplierProductReq 供应商商品列表查询
type ListSupplierProductReq struct {
	SupplierID uint   `form:"supplier_id"`
	CategoryID uint   `form:"category_id"`
	Keyword    string `form:"keyword"`
	Status     *int8  `form:"status"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
}
