package model

import "time"

// StoreSupplier 门店-供应商关联表（门店绑定哪些供应商）
type StoreSupplier struct {
	ID         uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID    uint      `json:"store_id" gorm:"not null;uniqueIndex:idx_store_supplier;comment:门店ID"`
	Store      *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	SupplierID uint      `json:"supplier_id" gorm:"not null;uniqueIndex:idx_store_supplier;comment:供应商ID"`
	Supplier   *Supplier `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	Status     int8      `json:"status" gorm:"not null;default:1;comment:状态 1=启用 0=禁用"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (StoreSupplier) TableName() string {
	return "store_suppliers"
}

// BindStoreSuppliersReq 门店绑定供应商请求
type BindStoreSuppliersReq struct {
	StoreID     uint   `json:"store_id" binding:"required"`           // 门店ID
	SupplierIDs []uint `json:"supplier_ids" binding:"required,min=1"` // 供应商ID列表
}

// UnbindStoreSuppliersReq 门店解绑供应商请求
type UnbindStoreSuppliersReq struct {
	StoreID     uint   `json:"store_id" binding:"required"`           // 门店ID
	SupplierIDs []uint `json:"supplier_ids" binding:"required,min=1"` // 供应商ID列表
}


