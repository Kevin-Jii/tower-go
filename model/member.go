package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// DecimalType 金额类型别名
type DecimalType = decimal.Decimal

// DecimalZero 返回零值金额
func DecimalZero() DecimalType {
	return decimal.Zero
}

// Member 会员表
type Member struct {
	ID         uint            `json:"id" gorm:"primaryKey"`
	StoreID    uint            `json:"store_id" gorm:"index;comment:所属门店ID"`
	UID        string          `json:"uid" gorm:"type:varchar(64);uniqueIndex;comment:用户唯一标识"`
	Name       string          `json:"name" gorm:"type:varchar(100);comment:会员姓名"`
	Phone      string          `json:"phone" gorm:"type:varchar(20);uniqueIndex;comment:手机号"`
	Balance    decimal.Decimal `json:"balance" gorm:"type:decimal(10,2);comment:余额"`
	Points     int             `json:"points" gorm:"type:int;default:0;comment:积分"`
	Level      int             `json:"level" gorm:"type:int;default:1;comment:等级"`
	Version    int             `json:"version" gorm:"type:int;default:0;comment:乐观锁版本号"`
	CreateTime time.Time       `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time       `json:"updateTime" gorm:"autoUpdateTime"`
}

// TableName 指定表名为 t_member
func (Member) TableName() string {
	return "t_member"
}

// WalletLog 流水表
type WalletLog struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	MemberID       uint            `json:"memberId" gorm:"index;comment:会员ID"`
	ChangeType     ChangeTypeEnum  `json:"changeType" gorm:"type:int;comment:变动类型"`
	ChangeAmount   decimal.Decimal `json:"changeAmount" gorm:"type:decimal(10,2);comment:变动金额"`
	BalanceAfter   decimal.Decimal `json:"balanceAfter" gorm:"type:decimal(10,2);comment:变动后余额"`
	RelatedOrderNo string          `json:"relatedOrderNo" gorm:"type:varchar(64);index;comment:关联单号"`
	Remark         string          `json:"remark" gorm:"type:varchar(255);comment:备注"`
	CreateTime     time.Time       `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime     time.Time       `json:"updateTime" gorm:"autoUpdateTime"`
}

// TableName 指定表名为 t_member_wallet_log
func (WalletLog) TableName() string {
	return "t_member_wallet_log"
}

// ChangeTypeEnum 流水变动类型枚举
type ChangeTypeEnum int

const (
	ChangeTypeRecharge   ChangeTypeEnum = 1 // 充值
	ChangeTypeConsume    ChangeTypeEnum = 2 // 消费
	ChangeTypeRefund     ChangeTypeEnum = 3 // 退款
	ChangeTypeAdjustAdd  ChangeTypeEnum = 4 // 调增
	ChangeTypeAdjustLess ChangeTypeEnum = 5 // 调减
)

// RechargeOrder 充值单表
type RechargeOrder struct {
	ID          uint            `json:"id" gorm:"primaryKey"`
	OrderNo     string          `json:"orderNo" gorm:"type:varchar(32);uniqueIndex;comment:单号"`
	MemberID    uint            `json:"memberId" gorm:"index;comment:会员ID"`
	MemberName  string          `json:"memberName" gorm:"-"`  // 不存数据库，关联查询
	MemberPhone string          `json:"memberPhone" gorm:"-"` // 不存数据库，关联查询
	PayAmount   decimal.Decimal `json:"payAmount" gorm:"type:decimal(10,2);comment:实付金额"`
	GiftAmount  decimal.Decimal `json:"giftAmount" gorm:"type:decimal(10,2);comment:赠送金额"`
	TotalAmount decimal.Decimal `json:"totalAmount" gorm:"type:decimal(10,2);comment:总金额"`
	PayType     int             `json:"payType" gorm:"type:int;default:0;comment:支付方式"`
	PayTypeName string          `json:"payTypeName" gorm:"-"` // 不存数据库
	PayStatus   PayStatusEnum   `json:"payStatus" gorm:"type:int;default:0;comment:支付状态"`
	StatusName  string          `json:"statusName" gorm:"-"` // 不存数据库
	PayTime     *time.Time      `json:"payTime" gorm:"comment:支付时间"`
	Remark      string          `json:"remark" gorm:"type:varchar(255);comment:备注"`
	CreateTime  time.Time       `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime  time.Time       `json:"updateTime" gorm:"autoUpdateTime"`
}

// TableName 指定表名为 t_recharge_order
func (RechargeOrder) TableName() string {
	return "t_recharge_order"
}

// PayStatusEnum 支付状态枚举
type PayStatusEnum int

const (
	PayStatusPending   PayStatusEnum = 0 // 待支付
	PayStatusPaid      PayStatusEnum = 1 // 已支付
	PayStatusCancelled PayStatusEnum = 2 // 已取消
	PayStatusRefunded  PayStatusEnum = 3 // 已退款
)

// String 获取状态名称
func (s PayStatusEnum) String() string {
	switch s {
	case PayStatusPending:
		return "待支付"
	case PayStatusPaid:
		return "已支付"
	case PayStatusCancelled:
		return "已取消"
	case PayStatusRefunded:
		return "已退款"
	default:
		return "未知"
	}
}

// ========== 请求结构体 ==========

// CreateMemberReq 创建会员请求
type CreateMemberReq struct {
	UID    string  `json:"uid"`  // 可选，不传则自动生成
	Name   string  `json:"name"` // 会员姓名
	Phone  string  `json:"phone" binding:"required"`
	Level  *int    `json:"level_id"` // 等级（可选）
	Remark *string `json:"remark"`   // 备注（可选，暂不存储）
}

// UpdateMemberReq 更新会员请求
type UpdateMemberReq struct {
	Name   *string `json:"name"`
	Phone  *string `json:"phone"`
	Points *int    `json:"points"`
	Level  *int    `json:"level"`
}

// AdjustBalanceReq 调整余额请求
type AdjustBalanceReq struct {
	Amount  decimal.Decimal `json:"amount" binding:"required"`
	Type    ChangeTypeEnum  `json:"type" binding:"required,oneof=4 5"` // 4=调增 5=调减
	Remark  string          `json:"remark"`
	Version int             `json:"version"` // 乐观锁版本号
}

// CreateRechargeOrderReq 创建充值单请求
type CreateRechargeOrderReq struct {
	MemberID   uint            `json:"memberId" binding:"required"`
	PayAmount  decimal.Decimal `json:"payAmount" binding:"required"`
	GiftAmount decimal.Decimal `json:"giftAmount"`
	PayType    int             `json:"payType" binding:"required"`
	Remark     string          `json:"remark"`
}

// PayRechargeOrderReq 支付充值单请求
type PayRechargeOrderReq struct {
	OrderNo string `json:"orderNo" binding:"required"`
}

// ListWalletLogReq 查询流水请求
type ListWalletLogReq struct {
	MemberID   uint            `form:"memberId"`
	ChangeType *ChangeTypeEnum `form:"changeType"`
	StartTime  *time.Time      `form:"startTime"`
	EndTime    *time.Time      `form:"endTime"`
}

// MemberConsumptionRecord 会员消费记录（来自门店记账）
type MemberConsumptionRecord struct {
	AccountID          uint      `json:"account_id"`
	AccountNo          string    `json:"account_no"`
	AccountDate        time.Time `json:"account_date"`
	Channel            string    `json:"channel"`
	ChannelName        string    `json:"channel_name"`
	OrderNo            string    `json:"order_no"`
	TotalAmount        float64   `json:"total_amount"`
	OtherExpenseAmount float64   `json:"other_expense_amount"`
	RoundAmount        float64   `json:"round_amount"`
	GiftWineCostAmount float64   `json:"gift_wine_cost_amount"`
	ConsumableAmount   float64   `json:"consumable_amount"`
	NetIncomeAmount    float64   `json:"net_income_amount"`
	CreatedAt          time.Time `json:"created_at"`
}

// MemberConsumptionSummary 会员消费汇总
type MemberConsumptionSummary struct {
	Count              int64   `json:"count"`
	TotalAmount        float64 `json:"total_amount"`
	OtherExpenseAmount float64 `json:"other_expense_amount"`
	RoundAmount        float64 `json:"round_amount"`
	GiftWineCostAmount float64 `json:"gift_wine_cost_amount"`
	ConsumableAmount   float64 `json:"consumable_amount"`
	NetIncomeAmount    float64 `json:"net_income_amount"`
}

const (
	MemberWineTxnDeposit  = 1 // 存入
	MemberWineTxnWithdraw = 2 // 取出
)

// MemberWineStorage 会员存酒当前存量
type MemberWineStorage struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StoreID   uint      `json:"store_id" gorm:"not null;index;uniqueIndex:idx_member_wine_storage_unique,priority:1;comment:门店ID"`
	MemberID  uint      `json:"member_id" gorm:"not null;index;uniqueIndex:idx_member_wine_storage_unique,priority:2;comment:会员ID"`
	Member    *Member   `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	WineName  string    `json:"wine_name" gorm:"type:varchar(120);not null;uniqueIndex:idx_member_wine_storage_unique,priority:3;comment:酒品名称"`
	Unit      string    `json:"unit" gorm:"type:varchar(20);not null;default:'瓶';uniqueIndex:idx_member_wine_storage_unique,priority:4;comment:单位"`
	Quantity  float64   `json:"quantity" gorm:"type:decimal(12,2);not null;default:0;comment:当前数量"`
	Remark    string    `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MemberWineStorage) TableName() string {
	return "member_wine_storages"
}

// MemberWineTransaction 会员存取酒流水
type MemberWineTransaction struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	StoreID      uint      `json:"store_id" gorm:"not null;index;comment:门店ID"`
	StorageID    uint      `json:"storage_id" gorm:"not null;index;comment:存酒记录ID"`
	MemberID     uint      `json:"member_id" gorm:"not null;index;comment:会员ID"`
	Member       *Member   `json:"member,omitempty" gorm:"foreignKey:MemberID"`
	Type         int       `json:"type" gorm:"not null;index;comment:类型 1=存入 2=取出"`
	WineName     string    `json:"wine_name" gorm:"type:varchar(120);not null;comment:酒品名称"`
	Unit         string    `json:"unit" gorm:"type:varchar(20);not null;default:'瓶';comment:单位"`
	Quantity     float64   `json:"quantity" gorm:"type:decimal(12,2);not null;comment:本次数量"`
	BalanceAfter float64   `json:"balance_after" gorm:"type:decimal(12,2);not null;default:0;comment:操作后数量"`
	Remark       string    `json:"remark" gorm:"type:varchar(500);comment:备注"`
	OperatorID   uint      `json:"operator_id" gorm:"index;comment:操作人ID"`
	OperatorName string    `json:"operator_name" gorm:"type:varchar(100);comment:操作人"`
	CreatedAt    time.Time `json:"created_at"`
}

func (MemberWineTransaction) TableName() string {
	return "member_wine_transactions"
}

type ListMemberWineStorageReq struct {
	StoreID   uint   `form:"store_id"`
	MemberID  uint   `form:"member_id"`
	Keyword   string `form:"keyword"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	OnlyStock int    `form:"only_stock"`
}

type MemberWineAdjustReq struct {
	MemberID uint    `json:"member_id" binding:"required"`
	WineName string  `json:"wine_name" binding:"required,max=120"`
	Unit     string  `json:"unit" binding:"max=20"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	Remark   string  `json:"remark" binding:"max=500"`
}

type ListMemberWineTransactionReq struct {
	StoreID   uint   `form:"store_id"`
	StorageID uint   `form:"storage_id"`
	MemberID  uint   `form:"member_id"`
	Type      int    `form:"type"`
	Keyword   string `form:"keyword"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
}
