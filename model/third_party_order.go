package model

import (
	"time"

	"gorm.io/gorm"
)

type ThirdPartyOrder struct {
	ID               uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountID        uint           `json:"account_id" gorm:"not null;index;comment:账号池ID"`
	PlatformName     string         `json:"platform_name" gorm:"type:varchar(50);not null;index;comment:平台"`
	OrderNo          string         `json:"order_no" gorm:"type:varchar(100);not null;index:idx_tp_order_no,unique;comment:第三方订单号"`
	PlaceTime        *time.Time     `json:"place_time" gorm:"index;comment:下单时间"`
	PlaceDate        string         `json:"place_date" gorm:"type:varchar(10);index;comment:下单日期"`
	OrderTradeStatus string         `json:"order_trade_status" gorm:"type:varchar(64);comment:交易状态编码"`
	StatusName       string         `json:"status_name" gorm:"type:varchar(100);comment:交易状态名称"`
	PayAmount        float64        `json:"pay_amount" gorm:"type:decimal(12,2);comment:支付金额"`
	TotalAmount      float64        `json:"total_amount" gorm:"type:decimal(12,2);comment:订单金额"`
	TotalItemNum     float64        `json:"total_item_num" gorm:"type:decimal(12,2);comment:总件数"`
	RawJSON          string         `json:"raw_json" gorm:"type:longtext;comment:原始订单JSON"`
	SyncedAt         time.Time      `json:"synced_at" gorm:"index;comment:同步时间"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ThirdPartyOrder) TableName() string {
	return "third_party_orders"
}
