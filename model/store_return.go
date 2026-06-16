package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreReturn 门店返厂单
type StoreReturn struct {
	ID           uint              `json:"id" gorm:"primaryKey;autoIncrement"`
	ReturnNo     string            `json:"return_no" gorm:"type:varchar(50);uniqueIndex;not null;comment:返厂单号"`
	ClientReqID  string            `json:"client_request_id" gorm:"type:varchar(64);uniqueIndex;comment:前端提交幂等ID"`
	StoreID      uint              `json:"store_id" gorm:"not null;index;comment:门店ID"`
	ReturnDate   time.Time         `json:"return_date" gorm:"type:date;index;comment:返厂日期"`
	LogisticsFee float64           `json:"logistics_fee" gorm:"type:decimal(10,2);not null;default:0;comment:货拉拉费用"`
	TotalDeposit float64           `json:"total_deposit" gorm:"type:decimal(10,2);not null;default:0;comment:押金总额"`
	ItemCount    int               `json:"item_count" gorm:"not null;default:0;comment:商品明细数量"`
	Remark       string            `json:"remark" gorm:"type:varchar(500);comment:备注"`
	OperatorID   uint              `json:"operator_id" gorm:"not null;comment:操作人ID"`
	OperatorName string            `json:"operator_name" gorm:"type:varchar(50);comment:操作人姓名"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	DeletedAt    gorm.DeletedAt    `json:"-" gorm:"index"`
	Items        []StoreReturnItem `json:"items,omitempty" gorm:"foreignKey:ReturnID"`
	Store        *Store            `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Operator     *User             `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

func (StoreReturn) TableName() string {
	return "store_returns"
}

// StoreReturnItem 门店返厂商品明细
type StoreReturnItem struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ReturnID    uint           `json:"return_id" gorm:"not null;index;comment:返厂单ID"`
	ProductID   uint           `json:"product_id" gorm:"not null;default:0;index;comment:返厂商品档案ID"`
	ProductName string         `json:"product_name" gorm:"type:varchar(200);not null;comment:商品名称"`
	Quantity    float64        `json:"quantity" gorm:"type:decimal(10,2);not null;default:1;comment:返厂数量"`
	Deposit     float64        `json:"deposit" gorm:"type:decimal(10,2);not null;default:0;comment:单件押金"`
	Remark      string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (StoreReturnItem) TableName() string {
	return "store_return_items"
}

// StoreReturnProduct 门店返厂商品档案
type StoreReturnProduct struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID     uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	ProductName string         `json:"product_name" gorm:"type:varchar(200);not null;comment:商品名称"`
	Deposit     float64        `json:"deposit" gorm:"type:decimal(10,2);not null;default:0;comment:默认押金"`
	Remark      string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	Status      int            `json:"status" gorm:"not null;default:1;index;comment:状态 1=启用 0=停用"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Store       *Store         `json:"store,omitempty" gorm:"foreignKey:StoreID"`
}

func (StoreReturnProduct) TableName() string {
	return "store_return_products"
}

type CreateStoreReturnReq struct {
	StoreID      uint                       `json:"store_id"`
	ClientReqID  string                     `json:"client_request_id" binding:"max=64"`
	ReturnDate   string                     `json:"return_date" binding:"required"`
	LogisticsFee float64                    `json:"logistics_fee" binding:"gte=0"`
	Remark       string                     `json:"remark" binding:"max=500"`
	Items        []CreateStoreReturnItemReq `json:"items" binding:"required,min=1,dive"`
}

type CreateStoreReturnItemReq struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name" binding:"max=200"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Deposit     float64 `json:"deposit" binding:"gte=0"`
	Remark      string  `json:"remark" binding:"max=500"`
}

type UpdateStoreReturnReq struct {
	StoreID      uint                       `json:"store_id"`
	ReturnDate   string                     `json:"return_date" binding:"required"`
	LogisticsFee float64                    `json:"logistics_fee" binding:"gte=0"`
	Remark       string                     `json:"remark" binding:"max=500"`
	Items        []CreateStoreReturnItemReq `json:"items" binding:"required,min=1,dive"`
}

type ListStoreReturnReq struct {
	StoreID   uint   `form:"store_id"`
	Keyword   string `form:"keyword"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

type CreateStoreReturnProductReq struct {
	StoreID     uint    `json:"store_id"`
	ProductName string  `json:"product_name" binding:"required,max=200"`
	Deposit     float64 `json:"deposit" binding:"gte=0"`
	Remark      string  `json:"remark" binding:"max=500"`
	Status      int     `json:"status" binding:"omitempty,oneof=0 1"`
}

type UpdateStoreReturnProductReq struct {
	StoreID     uint    `json:"store_id"`
	ProductName string  `json:"product_name" binding:"required,max=200"`
	Deposit     float64 `json:"deposit" binding:"gte=0"`
	Remark      string  `json:"remark" binding:"max=500"`
	Status      int     `json:"status" binding:"omitempty,oneof=0 1"`
}

type ListStoreReturnProductReq struct {
	StoreID  uint   `form:"store_id"`
	Keyword  string `form:"keyword"`
	Status   *int   `form:"status"`
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=500"`
}

type StoreReturnStats struct {
	TotalDeposit float64 `json:"total_deposit"`
	LogisticsFee float64 `json:"logistics_fee"`
	ReturnCount  int64   `json:"return_count"`
	ItemCount    int64   `json:"item_count"`
}
