package model

import (
	"time"

	"gorm.io/gorm"
)

// StoreAccount 门店记账（主表）
type StoreAccount struct {
	ID          uint               `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountNo   string             `json:"account_no" gorm:"type:varchar(50);uniqueIndex;not null;comment:记账编号"`
	StoreID     uint               `json:"store_id" gorm:"not null;index;comment:门店ID"`
	Channel     string             `json:"channel" gorm:"type:varchar(50);index;comment:渠道来源(字典:sales_channel)"`
	OrderNo     string             `json:"order_no" gorm:"type:varchar(100);index;comment:订单编号"`
	TotalAmount float64            `json:"total_amount" gorm:"type:decimal(10,2);comment:总金额"`
	ItemCount   int                `json:"item_count" gorm:"comment:商品数量"`
	TagCode     string             `json:"tag_code" gorm:"type:varchar(50);index;comment:标签编码"`
	TagName     string             `json:"tag_name" gorm:"type:varchar(100);comment:标签名称"`
	Remark      string             `json:"remark" gorm:"type:text;comment:备注"`
	OperatorID  uint               `json:"operator_id" gorm:"not null;comment:操作人ID"`
	AccountDate time.Time          `json:"account_date" gorm:"type:date;index;comment:记账日期"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	DeletedAt   gorm.DeletedAt     `json:"-" gorm:"index"`
	Items       []StoreAccountItem `json:"items,omitempty" gorm:"foreignKey:AccountID"`

	// 关联
	Store    *Store `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Operator *User  `json:"operator,omitempty" gorm:"foreignKey:OperatorID"`
}

func (StoreAccount) TableName() string {
	return "store_accounts"
}

// StoreAccountItem 门店记账明细（子表）
type StoreAccountItem struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID   uint           `json:"account_id" gorm:"not null;index;comment:记账ID"`
	ProductID   uint           `json:"product_id" gorm:"not null;index;comment:商品ID"`
	ProductName string         `json:"product_name" gorm:"type:varchar(200);comment:商品名称"`
	Spec        string         `json:"spec" gorm:"type:varchar(100);comment:规格"`
	Quantity    float64        `json:"quantity" gorm:"type:decimal(10,2);default:1;comment:数量"`
	Unit        string         `json:"unit" gorm:"type:varchar(20);comment:单位"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);comment:单价"`
	Amount      float64        `json:"amount" gorm:"type:decimal(10,2);comment:金额"`
	Remark      string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (StoreAccountItem) TableName() string {
	return "store_account_items"
}

// CreateStoreAccountReq 创建记账请求
type CreateStoreAccountReq struct {
	Channel     string                      `json:"channel" binding:"required,max=50"`
	OrderNo     string                      `json:"order_no" binding:"max=100"`
	TagCode     string                      `json:"tag_code" binding:"max=50"`
	TagName     string                      `json:"tag_name" binding:"max=100"`
	Remark      string                      `json:"remark" binding:"max=500"`
	AccountDate string                      `json:"account_date"`
	Items       []CreateStoreAccountItemReq `json:"items" binding:"required,min=1,dive"`
	NotifyImage string                      `json:"notify_image"` // 通知图片URL（前端生成）
}

// CreateStoreAccountItemReq 创建记账明细请求
type CreateStoreAccountItemReq struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Spec      string  `json:"spec" binding:"max=100"`
	Quantity  float64 `json:"quantity" binding:"required,gt=0"`
	Unit      string  `json:"unit" binding:"max=20"`
	Price     float64 `json:"price" binding:"gte=0"`
	Amount    float64 `json:"amount" binding:"gte=0"`
	Remark    string  `json:"remark" binding:"max=500"`
}

// UpdateStoreAccountReq 更新记账请求
type UpdateStoreAccountReq struct {
	Channel     string `json:"channel" binding:"max=50"`
	OrderNo     string `json:"order_no" binding:"max=100"`
	TagCode     string `json:"tag_code" binding:"max=50"`
	TagName     string `json:"tag_name" binding:"max=100"`
	Remark      string `json:"remark" binding:"max=500"`
	AccountDate string `json:"account_date"`
}

// ListStoreAccountReq 记账列表查询
type ListStoreAccountReq struct {
	StoreID   uint   `form:"store_id"`
	Channel   string `form:"channel"`
	OrderNo   string `form:"order_no"`
	TagCode   string `form:"tag_code"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}
