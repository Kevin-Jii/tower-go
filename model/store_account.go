package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreAccount 门店记账
type StoreAccount struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountNo    string         `json:"account_no" gorm:"type:varchar(50);uniqueIndex;not null;comment:记账编号"`
	StoreID      uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	ProductID    uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	Spec         string         `json:"spec" gorm:"type:varchar(100);comment:规格"`
	Quantity     float64        `json:"quantity" gorm:"type:decimal(10,2);default:1;comment:数量"`
	Unit         string         `json:"unit" gorm:"type:varchar(20);comment:单位"`
	Price        float64        `json:"price" gorm:"type:decimal(10,2);comment:单价"`
	Amount       float64        `json:"amount" gorm:"type:decimal(10,2);comment:金额"`
	Channel      string         `json:"channel" gorm:"type:varchar(50);index;comment:销售渠道(字典:sales_channel)"`
	OrderSource  string         `json:"order_source" gorm:"type:varchar(50);index;comment:订单来源(字典:order_source)"`
	OrderNo      string         `json:"order_no" gorm:"type:varchar(100);index;comment:订单编号"`
	Remark       string         `json:"remark" gorm:"type:text;comment:备注"`
	OperatorID   uint           `json:"operator_id" gorm:"not null;comment:操作人ID"`
	AccountDate  time.Time      `json:"account_date" gorm:"type:date;index;comment:记账日期"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联
	Store    *Store           `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Product  *SupplierProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Operator *User            `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

func (StoreAccount) TableName() string {
	return "store_accounts"
}

// CreateStoreAccountReq 创建记账请求
type CreateStoreAccountReq struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	Spec        string  `json:"spec" binding:"max=100"`
	Quantity    float64 `json:"quantity" binding:"required,gt=0"`
	Unit        string  `json:"unit" binding:"max=20"`
	Price       float64 `json:"price" binding:"gte=0"`
	Amount      float64 `json:"amount" binding:"gte=0"`
	Channel     string  `json:"channel" binding:"required,max=50"`
	OrderSource string  `json:"order_source" binding:"required,max=50"`
	OrderNo     string  `json:"order_no" binding:"max=100"`
	Remark      string  `json:"remark" binding:"max=500"`
	AccountDate string  `json:"account_date"` // 格式: 2024-12-07
}

// UpdateStoreAccountReq 更新记账请求
type UpdateStoreAccountReq struct {
	ProductID   *uint    `json:"product_id"`
	Spec        string   `json:"spec" binding:"max=100"`
	Quantity    *float64 `json:"quantity" binding:"omitempty,gt=0"`
	Unit        string   `json:"unit" binding:"max=20"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0"`
	Amount      *float64 `json:"amount" binding:"omitempty,gte=0"`
	Channel     string   `json:"channel" binding:"max=50"`
	OrderSource string   `json:"order_source" binding:"max=50"`
	OrderNo     string   `json:"order_no" binding:"max=100"`
	Remark      string   `json:"remark" binding:"max=500"`
	AccountDate string   `json:"account_date"`
}

// ListStoreAccountReq 记账列表查询
type ListStoreAccountReq struct {
	StoreID     uint   `form:"store_id"`
	ProductID   uint   `form:"product_id"`
	Channel     string `form:"channel"`
	OrderSource string `form:"order_source"`
	OrderNo     string `form:"order_no"`
	StartDate   string `form:"start_date"`
	EndDate     string `form:"end_date"`
	Page        int    `form:"page,default=1" binding:"min=1"`
	PageSize    int    `form:"page_size,default=20" binding:"min=1,max=100"`
}
