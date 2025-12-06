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

// InventoryOrder 出入库单（主表）
type InventoryOrder struct {
	ID            uint                  `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo       string                `json:"order_no" gorm:"type:varchar(50);uniqueIndex;not null;comment:单据编号"`
	Type          int8                  `json:"type" gorm:"not null;comment:类型 1=入库 2=出库"`
	StoreID       uint                  `json:"store_id" gorm:"not null;index;comment:门店ID"`
	StoreName     string                `json:"store_name" gorm:"type:varchar(100);comment:门店名称"`
	Reason        string                `json:"reason" gorm:"type:varchar(100);comment:原因"`
	Remark        string                `json:"remark" gorm:"type:text;comment:备注"`
	TotalQuantity float64               `json:"total_quantity" gorm:"type:decimal(10,2);comment:总数量"`
	ItemCount     int                   `json:"item_count" gorm:"comment:商品种类数"`
	OperatorID    uint                  `json:"operator_id" gorm:"not null;comment:操作人ID"`
	OperatorName  string                `json:"operator_name" gorm:"type:varchar(50);comment:操作人姓名"`
	OperatorPhone string                `json:"operator_phone" gorm:"type:varchar(20);comment:操作人手机号"`
	CreatedAt     time.Time             `json:"created_at"`
	UpdatedAt     time.Time             `json:"updated_at"`
	DeletedAt     gorm.DeletedAt        `json:"-" gorm:"index"`
	Items         []InventoryOrderItem  `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (InventoryOrder) TableName() string {
	return "inventory_orders"
}

// InventoryOrderItem 出入库单明细（子表）
type InventoryOrderItem struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID        uint           `json:"order_id" gorm:"not null;index;comment:出入库单ID"`
	ProductID      uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	ProductName    string         `json:"product_name" gorm:"type:varchar(200);comment:商品名称"`
	Quantity       float64        `json:"quantity" gorm:"type:decimal(10,2);not null;comment:数量"`
	Unit           string         `json:"unit" gorm:"type:varchar(20);comment:单位"`
	ProductionDate *time.Time     `json:"production_date" gorm:"type:date;comment:生产日期"`
	ExpiryDate     *time.Time     `json:"expiry_date" gorm:"type:date;comment:截止日期/保质期"`
	Remark         string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt      time.Time      `json:"created_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (InventoryOrderItem) TableName() string {
	return "inventory_order_items"
}

// 出入库类型常量
const (
	InventoryTypeIn  int8 = 1 // 入库
	InventoryTypeOut int8 = 2 // 出库
)

// 出入库原因
const (
	ReasonPurchase    = "采购入库"
	ReasonReturn      = "退货出库"
	ReasonAdjust      = "库存调整"
	ReasonLoss        = "报损出库"
	ReasonSale        = "销售出库"
	ReasonTransferIn  = "调拨入库"
	ReasonTransferOut = "调拨出库"
)

// CreateInventoryOrderReq 创建出入库单请求
type CreateInventoryOrderReq struct {
	Type   int8                          `json:"type" binding:"required,oneof=1 2"`
	Reason string                        `json:"reason" binding:"required,max=100"`
	Remark string                        `json:"remark" binding:"max=500"`
	Items  []CreateInventoryOrderItemReq `json:"items" binding:"required,min=1,dive"`
}

// CreateInventoryOrderItemReq 出入库单明细请求
type CreateInventoryOrderItemReq struct {
	ProductID      uint    `json:"product_id" binding:"required"`
	Quantity       float64 `json:"quantity" binding:"required,gt=0"`
	ProductionDate string  `json:"production_date"` // 生产日期
	ExpiryDate     string  `json:"expiry_date"`     // 截止日期
	Remark         string  `json:"remark" binding:"max=500"`
}

// ListInventoryReq 库存列表查询
type ListInventoryReq struct {
	StoreID   uint   `form:"store_id"`
	ProductID uint   `form:"product_id"`
	Keyword   string `form:"keyword"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ListInventoryOrderReq 出入库单列表查询
type ListInventoryOrderReq struct {
	StoreID  uint   `form:"store_id"`
	Type     *int8  `form:"type"`
	OrderNo  string `form:"order_no"`
	Date     string `form:"date"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
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
