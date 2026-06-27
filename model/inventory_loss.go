package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	InventoryLossTypeLoss    = "loss"
	InventoryLossTypeSelfUse = "self_use"
	InventoryLossTypeGift    = "gift"

	InventoryLossReasonDictCode = "PERSONALUSE"
)

// InventoryLossOrder 库存报损/自用/赠送单
type InventoryLossOrder struct {
	ID           uint                     `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo      string                   `json:"order_no" gorm:"type:varchar(50);uniqueIndex;not null;comment:单据编号"`
	StoreID      uint                     `json:"store_id" gorm:"not null;index;comment:门店ID"`
	Type         string                   `json:"type" gorm:"type:varchar(20);not null;index;comment:类型 loss=报损 self_use=自用 gift=赠送"`
	MemberID     *uint                    `json:"member_id,omitempty" gorm:"index;comment:赠送会员ID"`
	Member       *Member                  `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Reason       string                   `json:"reason" gorm:"type:varchar(200);comment:原因"`
	TotalCost    float64                  `json:"total_cost" gorm:"type:decimal(10,2);not null;default:0;comment:总成本"`
	ItemCount    int                      `json:"item_count" gorm:"not null;default:0;comment:明细数量"`
	OperatorID   uint                     `json:"operator_id" gorm:"not null;comment:操作人ID"`
	OperatorName string                   `json:"operator_name" gorm:"type:varchar(50);comment:操作人姓名"`
	IsCanceled   bool                     `json:"is_canceled" gorm:"not null;default:false;index;comment:是否撤销"`
	CanceledAt   *time.Time               `json:"canceled_at,omitempty" gorm:"comment:撤销时间"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	DeletedAt    gorm.DeletedAt           `json:"-" gorm:"index"`
	Items        []InventoryLossOrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (InventoryLossOrder) TableName() string {
	return "inventory_loss_orders"
}

// InventoryLossOrderItem 库存报损/自用/赠送明细
type InventoryLossOrderItem struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID      uint           `json:"order_id" gorm:"not null;index;comment:单据ID"`
	ProductID    uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	ProductName  string         `json:"product_name" gorm:"type:varchar(200);comment:商品名称"`
	Unit         string         `json:"unit" gorm:"type:varchar(50);comment:选择规格单位"`
	Quantity     float64        `json:"quantity" gorm:"type:decimal(10,2);not null;comment:选择规格数量"`
	BaseQuantity float64        `json:"base_quantity" gorm:"type:decimal(10,2);not null;comment:基础库存扣减数量"`
	BaseUnit     string         `json:"base_unit" gorm:"type:varchar(20);comment:基础库存单位"`
	CostPrice    float64        `json:"cost_price" gorm:"type:decimal(10,2);not null;default:0;comment:所选规格成本单价"`
	CostAmount   float64        `json:"cost_amount" gorm:"type:decimal(10,2);not null;default:0;comment:成本金额"`
	Remark       string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (InventoryLossOrderItem) TableName() string {
	return "inventory_loss_order_items"
}

type CreateInventoryLossOrderReq struct {
	StoreID  uint                              `json:"store_id"`
	Type     string                            `json:"type" binding:"required,oneof=loss self_use gift"`
	MemberID *uint                             `json:"member_id"`
	Reason   string                            `json:"reason" binding:"required,max=200"`
	Items    []CreateInventoryLossOrderItemReq `json:"items" binding:"required,min=1,dive"`
}

type UpdateInventoryLossOrderReq struct {
	Reason string `json:"reason" binding:"required,max=200"`
}

type CreateInventoryLossOrderItemReq struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Unit      string  `json:"unit" binding:"required,max=50"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Remark    string  `json:"remark" binding:"max=500"`
}

type ListInventoryLossOrderReq struct {
	StoreID   uint   `form:"store_id"`
	Type      string `form:"type"`
	MemberID  uint   `form:"member_id"`
	Keyword   string `form:"keyword"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`

	DataScope int8   `json:"-"`
	UserID    uint   `json:"-"`
	RoleCode  string `json:"-"`
}

type ListMemberGiftRecordsReq struct {
	StoreID   uint   `form:"store_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

type MemberGiftRecord struct {
	ID           uint      `json:"id"`
	OrderID      uint      `json:"order_id"`
	OrderNo      string    `json:"order_no"`
	ProductID    uint      `json:"product_id"`
	ProductName  string    `json:"product_name"`
	Unit         string    `json:"unit"`
	Quantity     float64   `json:"quantity"`
	CostAmount   float64   `json:"cost_amount"`
	Reason       string    `json:"reason"`
	OperatorName string    `json:"operator_name"`
	CreatedAt    time.Time `json:"created_at"`
}
