package model

import "time"

// StoreSmsConfig 每个门店独立配置阿里云短信（多门店共享 AccessKey 不安全，分开配置）。
// AccessKeySecret 以密文形式落表，密钥由 utils/encryption 包管理（APP_ENCRYPTION_KEY / JWT_SECRET）。
type StoreSmsConfig struct {
	ID                    uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID               uint       `json:"store_id" gorm:"not null;uniqueIndex;comment:门店ID"`
	AccessKeyID           string     `json:"access_key_id" gorm:"type:varchar(64);comment:阿里云 AccessKey ID"`
	AccessKeySecretCipher string     `json:"-" gorm:"type:varchar(1024);comment:AccessKey Secret 密文（仅内部使用）"`
	RegionID              string     `json:"region_id" gorm:"type:varchar(20);not null;default:cn-hangzhou"`
	SignName              string     `json:"sign_name" gorm:"type:varchar(64);comment:默认签名"`
	Enabled               bool       `json:"enabled" gorm:"not null;default:false;comment:是否启用短信服务"`
	SendWindowStart       string     `json:"send_window_start" gorm:"type:varchar(8);comment:HH:MM 起点，留空用全局"`
	SendWindowEnd         string     `json:"send_window_end" gorm:"type:varchar(8);comment:HH:MM 终点（不含）"`
	ConfiguredBy          *uint      `json:"configured_by" gorm:"comment:最后配置人"`
	LastTestedAt          *time.Time `json:"last_tested_at" gorm:"comment:最近一次连通性测试时间"`
	LastTestMessage       string     `json:"last_test_message" gorm:"type:varchar(255);comment:最近一次连通性测试结果"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
}

func (StoreSmsConfig) TableName() string { return "store_sms_configs" }

// StoreSmsConfigResp 返回给前端的视图：Secret 不返回原文，仅返回掩码。
type StoreSmsConfigResp struct {
	ID              uint       `json:"id"`
	StoreID         uint       `json:"store_id"`
	AccessKeyID     string     `json:"access_key_id"`
	AccessKeySecret string     `json:"access_key_secret"` // 掩码形式返回
	RegionID        string     `json:"region_id"`
	SignName        string     `json:"sign_name"`
	Enabled         bool       `json:"enabled"`
	SendWindowStart string     `json:"send_window_start"`
	SendWindowEnd   string     `json:"send_window_end"`
	LastTestedAt    *time.Time `json:"last_tested_at"`
	LastTestMessage string     `json:"last_test_message"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// UpsertStoreSmsConfigReq 保存/更新配置；AccessKeySecret 非必填，留空表示不动现有密钥。
type UpsertStoreSmsConfigReq struct {
	AccessKeyID     string `json:"access_key_id" binding:"max=64"`
	AccessKeySecret string `json:"access_key_secret" binding:"omitempty,max=64"`
	RegionID        string `json:"region_id" binding:"max=20"`
	SignName        string `json:"sign_name" binding:"max=64"`
	Enabled         bool   `json:"enabled"`
	SendWindowStart string `json:"send_window_start" binding:"omitempty,len=5"`
	SendWindowEnd   string `json:"send_window_end" binding:"omitempty,len=5"`
}

// StoreSmsConfigTestReq 连通性测试请求。
type StoreSmsConfigTestReq struct {
	AccessKeyID     string `json:"access_key_id" binding:"required,max=64"`
	AccessKeySecret string `json:"access_key_secret" binding:"required,max=64"`
	RegionID        string `json:"region_id" binding:"required,max=20"`
}
