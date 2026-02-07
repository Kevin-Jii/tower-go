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
	MemberName  string          `json:"memberName" gorm:"-"` // 不存数据库，关联查询
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
	CreateTime  time.Time      `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime  time.Time      `json:"updateTime" gorm:"autoUpdateTime"`
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
	Remark     string         `json:"remark"`
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
