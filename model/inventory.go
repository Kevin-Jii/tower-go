package model

import (
	"time"

	"gorm.io/gorm"
)

// Inventory 门店库存表
type Inventory struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID   uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	ProductID uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	Quantity  float64        `json:"quantity" gorm:"type:decimal(10,2);default:0;comment:库存数量"`
	Unit      string         `json:"unit" gorm:"type:varchar(20);comment:单位"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Store   *Store           `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Product *SupplierProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (Inventory) TableName() string {
	return "inventories"
}

// InventoryRecord 库存出入库记录表
type InventoryRecord struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	RecordNo    string         `json:"record_no" gorm:"type:varchar(50);uniqueIndex;not null;comment:单据编号"`
	StoreID     uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	ProductID   uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	Type        int8           `json:"type" gorm:"not null;comment:类型 1=入库 2=出库"`
	Quantity    float64        `json:"quantity" gorm:"type:decimal(10,2);not null;comment:数量"`
	Unit        string         `json:"unit" gorm:"type:varchar(20);comment:单位"`
	Reason      string         `json:"reason" gorm:"type:varchar(100);comment:原因"`
	Remark      string         `json:"remark" gorm:"type:text;comment:备注"`
	OperatorID  uint           `json:"operator_id" gorm:"not null;comment:操作人ID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Store    *Store           `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Product  *SupplierProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Operator *User            `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

func (InventoryRecord) TableName() string {
	return "inventory_records"
}


// 出入库类型常量
const (
	InventoryTypeIn  int8 = 1 // 入库
	InventoryTypeOut int8 = 2 // 出库
)

// 出入库原因
const (
	ReasonPurchase   = "采购入库"
	ReasonReturn     = "退货出库"
	ReasonAdjust     = "库存调整"
	ReasonLoss       = "报损出库"
	ReasonSale       = "销售出库"
	ReasonTransferIn = "调拨入库"
	ReasonTransferOut = "调拨出库"
)

// CreateInventoryRecordReq 创建出入库记录请求
type CreateInventoryRecordReq struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Type      int8    `json:"type" binding:"required,oneof=1 2"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Reason    string  `json:"reason" binding:"required,max=100"`
	Remark    string  `json:"remark" binding:"max=500"`
}

// ListInventoryReq 库存列表查询
type ListInventoryReq struct {
	StoreID   uint   `form:"store_id"`
	ProductID uint   `form:"product_id"`
	Keyword   string `form:"keyword"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListInventoryRecordReq 出入库记录列表查询
type ListInventoryRecordReq struct {
	StoreID   uint   `form:"store_id"`
	ProductID uint   `form:"product_id"`
	Type      *int8  `form:"type"`
	Date      string `form:"date"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// InventoryWithProduct 带商品信息的库存
type InventoryWithProduct struct {
	ID          uint    `json:"id"`
	StoreID     uint    `json:"store_id"`
	StoreName   string  `json:"store_name"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}
