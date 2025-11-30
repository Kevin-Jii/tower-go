package model

import "time"

// PurchaseOrder 采购单（报菜单）
type PurchaseOrder struct {
	ID          uint                 `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo     string               `json:"order_no" gorm:"uniqueIndex;type:varchar(32);not null;comment:订单编号"`
	StoreID     uint                 `json:"store_id" gorm:"not null;index;comment:门店ID"`
	Store       *Store               `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	TotalAmount float64              `json:"total_amount" gorm:"type:decimal(12,2);not null;default:0;comment:总金额"`
	Status      int8                 `json:"status" gorm:"not null;default:1;comment:状态 1=待确认 2=已确认 3=已完成 4=已取消"`
	Remark      string               `json:"remark" gorm:"type:varchar(500);comment:备注"`
	OrderDate   time.Time            `json:"order_date" gorm:"type:date;not null;comment:报菜日期"`
	CreatedBy   uint                 `json:"created_by" gorm:"not null;comment:创建人ID"`
	Creator     *User                `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Items       []PurchaseOrderItem  `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}

// PurchaseOrderItem 采购单明细
type PurchaseOrderItem struct {
	ID         uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID    uint             `json:"order_id" gorm:"not null;index;comment:采购单ID"`
	SupplierID uint             `json:"supplier_id" gorm:"not null;index;comment:供应商ID"`
	Supplier   *Supplier        `json:"supplier,omitempty" gorm:"foreignKey:SupplierID"`
	ProductID  uint             `json:"product_id" gorm:"not null;comment:商品ID"`
	Product    *SupplierProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity   float64          `json:"quantity" gorm:"type:decimal(10,2);not null;comment:数量"`
	UnitPrice  float64          `json:"unit_price" gorm:"type:decimal(10,2);not null;comment:单价"`
	Amount     float64          `json:"amount" gorm:"type:decimal(12,2);not null;comment:金额"`
	Remark     string           `json:"remark" gorm:"type:varchar(200);comment:备注"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}

// CreatePurchaseOrderReq 创建采购单请求
type CreatePurchaseOrderReq struct {
	OrderDate string                       `json:"order_date" binding:"required"` // 格式: 2024-01-01
	Remark    string                       `json:"remark" binding:"max=500"`
	Items     []CreatePurchaseOrderItemReq `json:"items" binding:"required,min=1,dive"`
}

// CreatePurchaseOrderItemReq 采购单明细请求
type CreatePurchaseOrderItemReq struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Remark    string  `json:"remark" binding:"max=200"`
}

// UpdatePurchaseOrderReq 更新采购单请求
type UpdatePurchaseOrderReq struct {
	Status *int8  `json:"status,omitempty" binding:"omitempty,oneof=1 2 3 4"`
	Remark string `json:"remark,omitempty" binding:"max=500"`
}

// ListPurchaseOrderReq 采购单列表查询
type ListPurchaseOrderReq struct {
	StoreID    uint   `form:"store_id"`
	SupplierID uint   `form:"supplier_id"`
	Status     *int8  `form:"status"`
	StartDate  string `form:"start_date"` // 格式: 2024-01-01
	EndDate    string `form:"end_date"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// PurchaseOrderStatus 采购单状态
const (
	PurchaseStatusPending   int8 = 1 // 待确认
	PurchaseStatusConfirmed int8 = 2 // 已确认
	PurchaseStatusCompleted int8 = 3 // 已完成
	PurchaseStatusCancelled int8 = 4 // 已取消
)

// SupplierGroupedItems 按供应商分组的采购单明细
type SupplierGroupedItems struct {
	SupplierID   uint                        `json:"supplier_id"`
	SupplierName string                      `json:"supplier_name"`
	Items        []SupplierGroupedItemDetail `json:"items"`
	SubTotal     float64                     `json:"sub_total"`
}

// SupplierGroupedItemDetail 分组明细详情
type SupplierGroupedItemDetail struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Unit        string  `json:"unit"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
	Remark      string  `json:"remark"`
}
