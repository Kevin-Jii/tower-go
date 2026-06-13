package model

import "time"

// ProductUnitSpec 商品多单位配置（基础单位：L）
type ProductUnitSpec struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProductID    uint      `json:"product_id" gorm:"not null;index;uniqueIndex:uk_product_unit_specs_product_unit_name,priority:1;comment:商品ID"`
	UnitCode     string    `json:"unit_code" gorm:"type:varchar(50);not null;uniqueIndex:uk_product_unit_specs_product_unit_name,priority:2;comment:单位编码，如 bottle/case/barrel/liter"`
	UnitName     string    `json:"unit_name" gorm:"type:varchar(50);not null;uniqueIndex:uk_product_unit_specs_product_unit_name,priority:3;comment:单位名称，如 瓶/箱/桶/L"`
	FactorToBase float64   `json:"factor_to_base" gorm:"type:decimal(12,6);not null;default:1;comment:换算到基础单位L的系数"`
	Precision    int       `json:"precision" gorm:"not null;default:0;comment:数量精度(小数位)"`
	CostPrice    float64   `json:"cost_price" gorm:"type:decimal(10,2);not null;default:0;comment:单位成本价"`
	SalePrice    float64   `json:"sale_price" gorm:"type:decimal(10,2);not null;default:0;comment:单位售价"`
	IsEnabled    bool      `json:"is_enabled" gorm:"not null;default:true;comment:是否启用"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ProductUnitSpec) TableName() string {
	return "product_unit_specs"
}

type CreateProductUnitSpecReq struct {
	ProductID    uint    `json:"product_id" binding:"required"`
	UnitCode     string  `json:"unit_code" binding:"required,max=50"`
	UnitName     string  `json:"unit_name" binding:"omitempty,max=50"`
	FactorToBase float64 `json:"factor_to_base" binding:"required,gt=0"`
	Precision    int     `json:"precision" binding:"gte=0,lte=6"`
	CostPrice    float64 `json:"cost_price" binding:"gte=0"`
	SalePrice    float64 `json:"sale_price" binding:"gte=0"`
	IsEnabled    *bool   `json:"is_enabled"`
}

type UpdateProductUnitSpecReq struct {
	UnitCode     *string  `json:"unit_code,omitempty" binding:"omitempty,max=50"`
	UnitName     *string  `json:"unit_name,omitempty" binding:"omitempty,max=50"`
	FactorToBase *float64 `json:"factor_to_base,omitempty" binding:"omitempty,gt=0"`
	Precision    *int     `json:"precision,omitempty" binding:"omitempty,gte=0,lte=6"`
	CostPrice    *float64 `json:"cost_price,omitempty" binding:"omitempty,gte=0"`
	SalePrice    *float64 `json:"sale_price,omitempty" binding:"omitempty,gte=0"`
	IsEnabled    *bool    `json:"is_enabled,omitempty"`
}

type BatchUpsertProductUnitSpecsReq struct {
	ProductID uint                            `json:"product_id" binding:"required"`
	Units     []CreateProductUnitSpecUnitItem `json:"units" binding:"required,min=1,dive"`
}

type CreateProductUnitSpecUnitItem struct {
	UnitCode     string  `json:"unit_code" binding:"required,max=50"`
	UnitName     string  `json:"unit_name" binding:"omitempty,max=50"`
	FactorToBase float64 `json:"factor_to_base" binding:"required,gt=0"`
	Precision    int     `json:"precision" binding:"gte=0,lte=6"`
	CostPrice    float64 `json:"cost_price" binding:"gte=0"`
	SalePrice    float64 `json:"sale_price" binding:"gte=0"`
	IsEnabled    *bool   `json:"is_enabled"`
}
