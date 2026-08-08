package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	PreOrderStatusPending   int8 = 1
	PreOrderStatusPrepared  int8 = 2
	PreOrderStatusDelivered int8 = 3
	PreOrderStatusCancelled int8 = 4

	PreOrderReminderSending int8 = 1
	PreOrderReminderSent    int8 = 2
	PreOrderReminderFailed  int8 = 3
)

const (
	PreOrderReminderPreviousDay0930 = "previous_day_0930"
	PreOrderReminderPreviousDay1600 = "previous_day_1600"
	PreOrderReminderDueDay0930      = "due_day_0930"
	PreOrderReminderDueDay1600      = "due_day_1600"
)

// PreOrder records a customer's future delivery requirement.
type PreOrder struct {
	ID              uint                  `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo         string                `json:"order_no" gorm:"type:varchar(32);not null;uniqueIndex;comment:预订单号"`
	StoreID         uint                  `json:"store_id" gorm:"not null;index;comment:门店ID"`
	CustomerID      uint                  `json:"customer_id" gorm:"not null;index;comment:B2B客户ID"`
	CustomerName    string                `json:"customer_name" gorm:"type:varchar(100);not null;comment:客户名称快照"`
	ContactPerson   string                `json:"contact_person" gorm:"type:varchar(50);not null;default:'';comment:联系人快照"`
	ContactPhone    string                `json:"contact_phone" gorm:"type:varchar(20);not null;default:'';comment:联系电话快照"`
	DeliveryAddress string                `json:"delivery_address" gorm:"type:varchar(255);not null;default:'';comment:配送地址"`
	ScheduledAt     time.Time             `json:"scheduled_at" gorm:"not null;index;comment:计划配送时间"`
	Status          int8                  `json:"status" gorm:"not null;default:1;index;comment:状态 1=待备货 2=已备货 3=已配送 4=已取消"`
	Remark          string                `json:"remark" gorm:"type:varchar(500);not null;default:'';comment:备注"`
	CreatedBy       uint                  `json:"created_by" gorm:"not null;comment:创建人ID"`
	PreparedAt      *time.Time            `json:"prepared_at,omitempty"`
	DeliveredAt     *time.Time            `json:"delivered_at,omitempty"`
	CancelledAt     *time.Time            `json:"cancelled_at,omitempty"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	DeletedAt       gorm.DeletedAt        `json:"-" gorm:"index"`
	Store           *Store                `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Customer        *B2BCustomer          `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Creator         *User                 `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Items           []PreOrderItem        `json:"items,omitempty" gorm:"foreignKey:PreOrderID"`
	ReminderLogs    []PreOrderReminderLog `json:"reminder_logs,omitempty" gorm:"foreignKey:PreOrderID"`
}

func (PreOrder) TableName() string { return "pre_orders" }

type PreOrderItem struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	PreOrderID  uint           `json:"pre_order_id" gorm:"not null;index;comment:预订单ID"`
	ProductID   uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	ProductName string         `json:"product_name" gorm:"type:varchar(200);not null;comment:商品名称快照"`
	UnitSpecID  uint           `json:"unit_spec_id" gorm:"not null;index;comment:商品规格ID"`
	UnitName    string         `json:"unit_name" gorm:"type:varchar(50);not null;comment:规格名称快照"`
	Quantity    float64        `json:"quantity" gorm:"type:decimal(10,2);not null;comment:预订数量"`
	Remark      string         `json:"remark" gorm:"type:varchar(200);not null;default:'';comment:明细备注"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (PreOrderItem) TableName() string { return "pre_order_items" }

type PreOrderReminderLog struct {
	ID           uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PreOrderID   uint       `json:"pre_order_id" gorm:"not null;uniqueIndex:uk_pre_order_reminder;priority:1;index"`
	ReminderKey  string     `json:"reminder_key" gorm:"type:varchar(40);not null;uniqueIndex:uk_pre_order_reminder;priority:2"`
	Status       int8       `json:"status" gorm:"not null;default:1;index;comment:状态 1=发送中 2=成功 3=失败"`
	SentAt       *time.Time `json:"sent_at,omitempty"`
	ErrorMessage string     `json:"error_message,omitempty" gorm:"type:varchar(500);not null;default:''"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (PreOrderReminderLog) TableName() string { return "pre_order_reminder_logs" }

type CreatePreOrderItemReq struct {
	ProductID  uint    `json:"product_id" binding:"required"`
	UnitSpecID uint    `json:"unit_spec_id" binding:"required"`
	Quantity   float64 `json:"quantity" binding:"required,gt=0"`
	Remark     string  `json:"remark" binding:"max=200"`
}

type CreatePreOrderReq struct {
	StoreID         uint                    `json:"store_id"`
	CustomerID      uint                    `json:"customer_id" binding:"required"`
	ScheduledAt     string                  `json:"scheduled_at" binding:"required"`
	ContactPerson   string                  `json:"contact_person" binding:"max=50"`
	ContactPhone    string                  `json:"contact_phone" binding:"max=20"`
	DeliveryAddress string                  `json:"delivery_address" binding:"max=255"`
	Remark          string                  `json:"remark" binding:"max=500"`
	Items           []CreatePreOrderItemReq `json:"items" binding:"required,min=1,dive"`
}

type UpdatePreOrderReq struct {
	CustomerID      uint                    `json:"customer_id" binding:"required"`
	ScheduledAt     string                  `json:"scheduled_at" binding:"required"`
	ContactPerson   string                  `json:"contact_person" binding:"max=50"`
	ContactPhone    string                  `json:"contact_phone" binding:"max=20"`
	DeliveryAddress string                  `json:"delivery_address" binding:"max=255"`
	Remark          string                  `json:"remark" binding:"max=500"`
	Items           []CreatePreOrderItemReq `json:"items" binding:"required,min=1,dive"`
}

type UpdatePreOrderStatusReq struct {
	Status int8 `json:"status" binding:"required,oneof=1 2 3 4"`
}

type ListPreOrderReq struct {
	StoreID    uint   `form:"store_id"`
	CustomerID uint   `form:"customer_id"`
	Status     *int8  `form:"status" binding:"omitempty,oneof=1 2 3 4"`
	Keyword    string `form:"keyword"`
	StartDate  string `form:"start_date"`
	EndDate    string `form:"end_date"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
}
