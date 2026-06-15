package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	B2BCustomerStatusEnabled  = 1
	B2BCustomerStatusDisabled = 2

	B2BPaymentUnpaid  = 1
	B2BPaymentPartial = 2
	B2BPaymentPaid    = 3

	B2BDeliveryPending = 1
	B2BDeliveryDone    = 2
	B2BDeliveryCancel  = 3
)

type B2BCustomer struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID       uint           `json:"store_id" gorm:"not null;index;comment:所属门店ID"`
	Name          string         `json:"name" gorm:"type:varchar(100);not null;comment:客户名称"`
	CustomerType  string         `json:"customer_type" gorm:"type:varchar(50);comment:客户类型"`
	ContactPerson string         `json:"contact_person" gorm:"type:varchar(50);comment:联系人"`
	Phone         string         `json:"phone" gorm:"type:varchar(20);index;comment:联系电话"`
	Address       string         `json:"address" gorm:"type:varchar(255);comment:地址"`
	Settlement    string         `json:"settlement" gorm:"type:varchar(30);default:'cash';comment:结算方式 cash/week/month"`
	PriceLevel    string         `json:"price_level" gorm:"type:varchar(30);comment:价格等级"`
	CreditLimit   float64        `json:"credit_limit" gorm:"type:decimal(12,2);default:0;comment:信用额度"`
	Receivable    float64        `json:"receivable" gorm:"type:decimal(12,2);default:0;comment:当前应收余额"`
	Status        int            `json:"status" gorm:"not null;default:1;index;comment:状态 1=启用 2=停用"`
	Remark        string         `json:"remark" gorm:"type:text;comment:备注"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (B2BCustomer) TableName() string {
	return "b2b_customers"
}

type B2BCustomerProductPrice struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID     uint           `json:"store_id" gorm:"not null;index;comment:所属门店ID"`
	CustomerID  *uint          `json:"customer_id,omitempty" gorm:"index;comment:客户ID，为空表示价格等级价"`
	PriceLevel  string         `json:"price_level" gorm:"type:varchar(30);index;comment:价格等级"`
	ProductID   uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	UnitSpecID  uint           `json:"unit_spec_id" gorm:"not null;index;comment:规格ID"`
	UnitName    string         `json:"unit_name" gorm:"type:varchar(50);comment:规格名称"`
	SupplyPrice float64        `json:"supply_price" gorm:"type:decimal(10,2);not null;default:0;comment:供货价"`
	MinQuantity float64        `json:"min_quantity" gorm:"type:decimal(10,2);not null;default:1;comment:起订数量"`
	IsEnabled   bool           `json:"is_enabled" gorm:"not null;default:true;index;comment:是否启用"`
	Remark      string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Customer *B2BCustomer     `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Product  *SupplierProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	UnitSpec *ProductUnitSpec `json:"unit_spec,omitempty" gorm:"foreignKey:UnitSpecID"`
}

func (B2BCustomerProductPrice) TableName() string {
	return "b2b_customer_product_prices"
}

type B2BSupplyOrder struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderNo        string         `json:"order_no" gorm:"type:varchar(50);uniqueIndex;not null;comment:供货单号"`
	StoreID        uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	CustomerID     uint           `json:"customer_id" gorm:"not null;index;comment:客户ID"`
	CustomerName   string         `json:"customer_name" gorm:"type:varchar(100);comment:客户名称快照"`
	OrderDate      time.Time      `json:"order_date" gorm:"type:date;index;comment:供货日期"`
	TotalAmount    float64        `json:"total_amount" gorm:"type:decimal(12,2);default:0;comment:订单金额"`
	PaidAmount     float64        `json:"paid_amount" gorm:"type:decimal(12,2);default:0;comment:已收金额"`
	UnpaidAmount   float64        `json:"unpaid_amount" gorm:"type:decimal(12,2);default:0;comment:未收金额"`
	CostAmount     float64        `json:"cost_amount" gorm:"type:decimal(12,2);default:0;comment:成本金额"`
	ProfitAmount   float64        `json:"profit_amount" gorm:"type:decimal(12,2);default:0;comment:毛利金额"`
	PaymentStatus  int            `json:"payment_status" gorm:"not null;default:1;index;comment:收款状态 1=未收 2=部分 3=已收"`
	DeliveryStatus int            `json:"delivery_status" gorm:"not null;default:1;index;comment:配送状态 1=待配送 2=已配送 3=已取消"`
	Remark         string         `json:"remark" gorm:"type:text;comment:备注"`
	OperatorID     uint           `json:"operator_id" gorm:"not null;comment:操作人ID"`
	OperatorName   string         `json:"operator_name" gorm:"type:varchar(50);comment:操作人名称"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	Customer *B2BCustomer         `json:"customer,omitempty" gorm:"foreignKey:CustomerID"`
	Items    []B2BSupplyOrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
}

func (B2BSupplyOrder) TableName() string {
	return "b2b_supply_orders"
}

type B2BSupplyOrderItem struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	OrderID      uint           `json:"order_id" gorm:"not null;index;comment:供货单ID"`
	ProductID    uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	ProductName  string         `json:"product_name" gorm:"type:varchar(200);comment:商品名称快照"`
	UnitSpecID   uint           `json:"unit_spec_id" gorm:"not null;index;comment:规格ID"`
	UnitName     string         `json:"unit_name" gorm:"type:varchar(50);comment:规格名称快照"`
	FactorToBase float64        `json:"factor_to_base" gorm:"type:decimal(12,6);not null;default:1;comment:换算基础库存系数"`
	Quantity     float64        `json:"quantity" gorm:"type:decimal(10,2);not null;comment:下单数量"`
	BaseQuantity float64        `json:"base_quantity" gorm:"type:decimal(12,2);not null;comment:扣减基础库存数量"`
	SupplyPrice  float64        `json:"supply_price" gorm:"type:decimal(10,2);not null;comment:供货单价"`
	CostPrice    float64        `json:"cost_price" gorm:"type:decimal(10,2);not null;default:0;comment:成本单价"`
	Amount       float64        `json:"amount" gorm:"type:decimal(12,2);not null;comment:行金额"`
	CostAmount   float64        `json:"cost_amount" gorm:"type:decimal(12,2);not null;default:0;comment:行成本"`
	ProfitAmount float64        `json:"profit_amount" gorm:"type:decimal(12,2);not null;default:0;comment:行毛利"`
	Remark       string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (B2BSupplyOrderItem) TableName() string {
	return "b2b_supply_order_items"
}

type CreateB2BCustomerReq struct {
	StoreID       uint    `json:"store_id"`
	Name          string  `json:"name" binding:"required,max=100"`
	CustomerType  string  `json:"customer_type" binding:"max=50"`
	ContactPerson string  `json:"contact_person" binding:"max=50"`
	Phone         string  `json:"phone" binding:"max=20"`
	Address       string  `json:"address" binding:"max=255"`
	Settlement    string  `json:"settlement" binding:"max=30"`
	PriceLevel    string  `json:"price_level" binding:"max=30"`
	CreditLimit   float64 `json:"credit_limit" binding:"gte=0"`
	Remark        string  `json:"remark" binding:"max=500"`
}

type UpdateB2BCustomerReq struct {
	Name          *string  `json:"name,omitempty" binding:"omitempty,max=100"`
	CustomerType  *string  `json:"customer_type,omitempty" binding:"omitempty,max=50"`
	ContactPerson *string  `json:"contact_person,omitempty" binding:"omitempty,max=50"`
	Phone         *string  `json:"phone,omitempty" binding:"omitempty,max=20"`
	Address       *string  `json:"address,omitempty" binding:"omitempty,max=255"`
	Settlement    *string  `json:"settlement,omitempty" binding:"omitempty,max=30"`
	PriceLevel    *string  `json:"price_level,omitempty" binding:"omitempty,max=30"`
	CreditLimit   *float64 `json:"credit_limit,omitempty" binding:"omitempty,gte=0"`
	Status        *int     `json:"status,omitempty" binding:"omitempty,oneof=1 2"`
	Remark        *string  `json:"remark,omitempty" binding:"omitempty,max=500"`
}

type ListB2BCustomerReq struct {
	StoreID  uint   `form:"store_id"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

type UpsertB2BPriceReq struct {
	StoreID     uint    `json:"store_id"`
	CustomerID  *uint   `json:"customer_id"`
	PriceLevel  string  `json:"price_level" binding:"max=30"`
	ProductID   uint    `json:"product_id" binding:"required"`
	UnitSpecID  uint    `json:"unit_spec_id" binding:"required"`
	SupplyPrice float64 `json:"supply_price" binding:"required,gte=0"`
	MinQuantity float64 `json:"min_quantity" binding:"gte=0"`
	IsEnabled   *bool   `json:"is_enabled"`
	Remark      string  `json:"remark" binding:"max=500"`
}

type ListB2BPriceReq struct {
	StoreID    uint   `form:"store_id"`
	CustomerID uint   `form:"customer_id"`
	PriceLevel string `form:"price_level"`
	ProductID  uint   `form:"product_id"`
	Keyword    string `form:"keyword"`
	Page       int    `form:"page,default=1" binding:"min=1"`
	PageSize   int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

type CreateB2BSupplyOrderReq struct {
	StoreID        uint                       `json:"store_id"`
	CustomerID     uint                       `json:"customer_id" binding:"required"`
	OrderDate      string                     `json:"order_date"`
	PaidAmount     float64                    `json:"paid_amount" binding:"gte=0"`
	DeliveryStatus int                        `json:"delivery_status" binding:"omitempty,oneof=1 2"`
	Remark         string                     `json:"remark" binding:"max=500"`
	Items          []CreateB2BSupplyOrderItem `json:"items" binding:"required,min=1,dive"`
}

type CreateB2BSupplyOrderItem struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	UnitSpecID  uint    `json:"unit_spec_id" binding:"required"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	SupplyPrice float64 `json:"supply_price" binding:"gte=0"`
	Remark      string  `json:"remark" binding:"max=500"`
}

type ListB2BSupplyOrderReq struct {
	StoreID        uint   `form:"store_id"`
	CustomerID     uint   `form:"customer_id"`
	Keyword        string `form:"keyword"`
	PaymentStatus  int    `form:"payment_status"`
	DeliveryStatus int    `form:"delivery_status"`
	StartDate      string `form:"start_date"`
	EndDate        string `form:"end_date"`
	Page           int    `form:"page,default=1" binding:"min=1"`
	PageSize       int    `form:"page_size,default=20" binding:"min=1,max=100"`
}
