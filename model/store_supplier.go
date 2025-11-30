package model

import "time"

// StoreSupplierProduct 门店-供应商商品关联表（门店从哪个供应商进哪些货）
type StoreSupplierProduct struct {
	ID        uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID   uint             `json:"store_id" gorm:"not null;uniqueIndex:idx_store_product;comment:门店ID"`
	Store     *Store           `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	ProductID uint             `json:"product_id" gorm:"not null;uniqueIndex:idx_store_product;comment:供应商商品ID"`
	Product   *SupplierProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	IsDefault bool             `json:"is_default" gorm:"default:false;comment:是否默认供应商（同商品多供应商时）"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (StoreSupplierProduct) TableName() string {
	return "store_supplier_products"
}

// BindStoreSupplierReq 门店绑定供应商商品请求
type BindStoreSupplierReq struct {
	StoreID    uint   `json:"store_id" binding:"required"`
	ProductIDs []uint `json:"product_ids" binding:"required,min=1"`
}

// SetDefaultSupplierReq 设置默认供应商请求
type SetDefaultSupplierReq struct {
	StoreID   uint `json:"store_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
}
