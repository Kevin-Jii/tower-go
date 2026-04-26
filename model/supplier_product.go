package model

import "time"

// SupplierProduct 供应商商品表
type SupplierProduct struct {
	ID             uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	SupplierID     uint              `json:"supplier_id" gorm:"not null;index;comment:供应商ID"`
	Supplier       *Supplier         `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	CategoryID     uint              `json:"category_id" gorm:"not null;index;comment:分类ID"`
	Category       *SupplierCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Name           string            `json:"name" gorm:"type:varchar(200);not null;comment:商品名称"`
	Unit           string            `json:"unit" gorm:"type:varchar(20);not null;default:'斤';comment:单位"`
	Price          float64           `json:"price" gorm:"type:decimal(10,2);not null;default:0;comment:兼容单价(默认等于单瓶价)"`
	BottlePrice    float64           `json:"bottle_price" gorm:"type:decimal(10,2);not null;default:0;comment:单瓶价格"`
	CasePrice      float64           `json:"case_price" gorm:"type:decimal(10,2);not null;default:0;comment:整箱价格"`
	BottlesPerCase int               `json:"bottles_per_case" gorm:"not null;default:1;comment:每箱瓶数"`
	Spec           string            `json:"spec" gorm:"type:varchar(100);comment:规格"`
	Remark         string            `json:"remark" gorm:"type:varchar(500);comment:备注"`
	Status         int8              `json:"status" gorm:"not null;default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`

	// UnitSpecs 商品单位规格（来自 product_unit_specs，仅接口返回用，不落库）
	UnitSpecs []*ProductUnitSpec `json:"unit_specs,omitempty" gorm:"-"`
}

func (SupplierProduct) TableName() string {
	return "supplier_products"
}

// CreateSupplierProductReq 创建供应商商品请求
type CreateSupplierProductReq struct {
	SupplierID     uint     `json:"supplier_id" binding:"required"`
	CategoryID     uint     `json:"category_id" binding:"required"`
	Name           string   `json:"name" binding:"required,max=200"`
	Unit           string   `json:"unit" binding:"required,max=20"`
	Price          *float64 `json:"price" binding:"omitempty,gte=0"` // 兼容旧字段
	BottlePrice    float64  `json:"bottle_price" binding:"gte=0"`
	CasePrice      float64  `json:"case_price" binding:"gte=0"`
	BottlesPerCase int      `json:"bottles_per_case" binding:"required,gte=1"`
	Spec           string   `json:"spec" binding:"max=100"`
	Remark         string   `json:"remark" binding:"max=500"`
}

// UpdateSupplierProductReq 更新供应商商品请求
type UpdateSupplierProductReq struct {
	CategoryID     *uint    `json:"category_id,omitempty"`
	Name           string   `json:"name,omitempty" binding:"max=200"`
	Unit           string   `json:"unit,omitempty" binding:"max=20"`
	Price          *float64 `json:"price,omitempty" binding:"omitempty,gte=0"`
	BottlePrice    *float64 `json:"bottle_price,omitempty" binding:"omitempty,gte=0"`
	CasePrice      *float64 `json:"case_price,omitempty" binding:"omitempty,gte=0"`
	BottlesPerCase *int     `json:"bottles_per_case,omitempty" binding:"omitempty,gte=1"`
	Spec           string   `json:"spec,omitempty" binding:"max=100"`
	Remark         string   `json:"remark,omitempty" binding:"max=500"`
	Status         *int8    `json:"status,omitempty" binding:"omitempty,oneof=0 1"`
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
