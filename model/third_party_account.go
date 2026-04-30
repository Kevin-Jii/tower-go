package model

import (
	"time"

	"gorm.io/gorm"
)

// ThirdPartyAccount 第三方账号池
type ThirdPartyAccount struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	PlatformName   string         `json:"platform_name" gorm:"type:varchar(50);not null;default:'tsbeer';comment:平台标识"`
	Name           string         `json:"name" gorm:"type:varchar(100);not null;comment:账号名称"`
	LoginName      string         `json:"login_name" gorm:"type:varchar(100);not null;index;comment:登录名"`
	Phone          string         `json:"phone" gorm:"type:varchar(30);comment:手机号"`
	Password       string         `json:"password" gorm:"type:varchar(255);not null;comment:密码(可为加密串)"`
	ApplicationKey string         `json:"application_key" gorm:"type:varchar(128);not null;comment:第三方application-key"`
	LoginType      string         `json:"login_type" gorm:"type:varchar(10);not null;default:'2';comment:登录类型"`
	Channel        string         `json:"channel" gorm:"type:varchar(20);not null;default:'WEB';comment:渠道"`
	ShopID         string         `json:"shop_id" gorm:"type:varchar(64);comment:第三方店铺ID"`
	CustomerID     string         `json:"customer_id" gorm:"type:varchar(64);comment:第三方客户ID"`
	IsEnabled      bool           `json:"is_enabled" gorm:"not null;default:true;comment:是否启用"`
	LastTestOK     bool           `json:"last_test_ok" gorm:"not null;default:false;comment:最后一次测试是否成功"`
	LastTestMsg    string         `json:"last_test_msg" gorm:"type:varchar(500);comment:最后一次测试消息"`
	LastToken      string         `json:"last_token" gorm:"type:text;comment:最后一次token"`
	TokenValidTime int64          `json:"token_valid_time" gorm:"comment:token有效期秒"`
	LastTestAt     *time.Time     `json:"last_test_at" gorm:"comment:最后测试时间"`
	LastSyncAt     *time.Time     `json:"last_sync_at" gorm:"comment:最后同步时间"`
	LastSyncMsg    string         `json:"last_sync_msg" gorm:"type:varchar(500);comment:最后同步消息"`
	LastSyncCount  int            `json:"last_sync_count" gorm:"comment:最后同步订单数"`
	Remark         string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (ThirdPartyAccount) TableName() string {
	return "third_party_accounts"
}

type CreateThirdPartyAccountReq struct {
	PlatformName   string `json:"platform_name" binding:"required,max=50"`
	Name           string `json:"name" binding:"required,max=100"`
	LoginName      string `json:"login_name" binding:"required,max=100"`
	Phone          string `json:"phone" binding:"max=30"`
	Password       string `json:"password" binding:"required,max=255"`
	ApplicationKey string `json:"application_key" binding:"required,max=128"`
	LoginType      string `json:"login_type" binding:"max=10"`
	Channel        string `json:"channel" binding:"max=20"`
	ShopID         string `json:"shop_id" binding:"max=64"`
	CustomerID     string `json:"customer_id" binding:"max=64"`
	IsEnabled      *bool  `json:"is_enabled"`
	Remark         string `json:"remark" binding:"max=500"`
}

type UpdateThirdPartyAccountReq struct {
	PlatformName   *string `json:"platform_name" binding:"omitempty,max=50"`
	Name           *string `json:"name" binding:"omitempty,max=100"`
	LoginName      *string `json:"login_name" binding:"omitempty,max=100"`
	Phone          *string `json:"phone" binding:"omitempty,max=30"`
	Password       *string `json:"password" binding:"omitempty,max=255"`
	ApplicationKey *string `json:"application_key" binding:"omitempty,max=128"`
	LoginType      *string `json:"login_type" binding:"omitempty,max=10"`
	Channel        *string `json:"channel" binding:"omitempty,max=20"`
	ShopID         *string `json:"shop_id" binding:"omitempty,max=64"`
	CustomerID     *string `json:"customer_id" binding:"omitempty,max=64"`
	IsEnabled      *bool   `json:"is_enabled"`
	Remark         *string `json:"remark" binding:"omitempty,max=500"`
}
