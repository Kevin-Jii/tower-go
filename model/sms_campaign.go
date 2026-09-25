package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

const (
	SmsCampaignTypeActivity = "activity"
	SmsCampaignTypeHoliday  = "holiday"

	SmsCampaignTargetAllMembers = "all_members"
	SmsCampaignTargetStores     = "stores"
	SmsCampaignTargetCustom     = "custom"

	SmsCampaignStatusDraft     = "draft"
	SmsCampaignStatusScheduled = "scheduled"
	SmsCampaignStatusSending   = "sending"
	SmsCampaignStatusSent      = "sent"
	SmsCampaignStatusFailed    = "failed"
	SmsCampaignStatusCancelled = "cancelled"

	SmsSendRecordSuccess = "success"
	SmsSendRecordFailed  = "failed"
)

// SmsCampaign 短信推广活动（活动通知 / 节日祝福）
type SmsCampaign struct {
	ID              uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name            string     `json:"name" gorm:"type:varchar(120);not null;comment:活动名称"`
	CampaignType    string     `json:"campaign_type" gorm:"type:varchar(20);not null;index;comment:activity|holiday"`
	SignName        string     `json:"sign_name" gorm:"type:varchar(64);comment:短信签名,空则用全局配置"`
	TemplateCode    string     `json:"template_code" gorm:"type:varchar(64);not null;comment:阿里云模板CODE"`
	TemplateParam   string     `json:"template_param" gorm:"type:text;comment:模板变量JSON"`
	PersonalizeName bool       `json:"personalize_name" gorm:"default:false;comment:是否注入会员姓名变量name"`
	TargetType      string     `json:"target_type" gorm:"type:varchar(32);not null;comment:all_members|stores|custom"`
	StoreIDs        UintList   `json:"store_ids" gorm:"type:json;comment:目标门店ID列表"`
	CustomPhones    StringList `json:"custom_phones" gorm:"type:json;comment:自定义手机号列表"`
	ScheduledAt     *time.Time `json:"scheduled_at" gorm:"index;comment:计划发送时间"`
	Status          string     `json:"status" gorm:"type:varchar(20);not null;index;default:draft"`
	TotalCount      int        `json:"total_count" gorm:"default:0"`
	SuccessCount    int        `json:"success_count" gorm:"default:0"`
	FailCount       int        `json:"fail_count" gorm:"default:0"`
	LastError       string     `json:"last_error" gorm:"type:varchar(500)"`
	SentAt          *time.Time `json:"sent_at"`
	CreatedBy       uint       `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (SmsCampaign) TableName() string {
	return "sms_campaigns"
}

// SmsSendRecord 单条短信发送记录
type SmsSendRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CampaignID   uint      `json:"campaign_id" gorm:"index;not null"`
	Phone        string    `json:"phone" gorm:"type:varchar(20);index;not null"`
	MemberID     *uint     `json:"member_id"`
	BizID        string    `json:"biz_id" gorm:"type:varchar(64)"`
	Status       string    `json:"status" gorm:"type:varchar(20);not null"`
	ErrorMessage string    `json:"error_message" gorm:"type:varchar(500)"`
	CreatedAt    time.Time `json:"created_at"`
}

func (SmsSendRecord) TableName() string {
	return "sms_send_records"
}

// UintList 以 JSON 数组存储 uint 列表
type UintList []uint

func (u *UintList) Scan(value interface{}) error {
	if value == nil {
		*u = UintList{}
		return nil
	}
	var data []byte
	switch typed := value.(type) {
	case []byte:
		data = typed
	case string:
		data = []byte(typed)
	default:
		return fmt.Errorf("scan UintList from %T", value)
	}
	if len(data) == 0 {
		*u = UintList{}
		return nil
	}
	if err := json.Unmarshal(data, u); err != nil {
		return fmt.Errorf("decode UintList: %w", err)
	}
	if *u == nil {
		*u = UintList{}
	}
	return nil
}

func (u UintList) Value() (driver.Value, error) {
	if u == nil {
		u = UintList{}
	}
	data, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("encode UintList: %w", err)
	}
	return string(data), nil
}

type CreateSmsCampaignReq struct {
	Name            string     `json:"name" binding:"required,max=120"`
	CampaignType    string     `json:"campaign_type" binding:"required,oneof=activity holiday"`
	SignName        string     `json:"sign_name" binding:"max=64"`
	TemplateCode    string     `json:"template_code" binding:"required,max=64"`
	TemplateParam   string     `json:"template_param"`
	PersonalizeName *bool      `json:"personalize_name"`
	TargetType      string     `json:"target_type" binding:"required,oneof=all_members stores custom"`
	StoreIDs        UintList   `json:"store_ids"`
	CustomPhones    StringList `json:"custom_phones"`
	ScheduledAt     *time.Time `json:"scheduled_at"`
}

type UpdateSmsCampaignReq struct {
	Name            *string     `json:"name" binding:"omitempty,max=120"`
	CampaignType    *string     `json:"campaign_type" binding:"omitempty,oneof=activity holiday"`
	SignName        *string     `json:"sign_name" binding:"omitempty,max=64"`
	TemplateCode    *string     `json:"template_code" binding:"omitempty,max=64"`
	TemplateParam   *string     `json:"template_param"`
	PersonalizeName *bool       `json:"personalize_name"`
	TargetType      *string     `json:"target_type" binding:"omitempty,oneof=all_members stores custom"`
	StoreIDs        *UintList   `json:"store_ids"`
	CustomPhones    *StringList `json:"custom_phones"`
	ScheduledAt     *time.Time  `json:"scheduled_at"`
}

type SmsServiceConfigResp struct {
	Enabled        bool   `json:"enabled"`
	DefaultSign    string `json:"default_sign_name"`
	Region         string `json:"region"`
	Configured     bool   `json:"configured"`
	HelpURL        string `json:"help_url"`
	TemplateHint   string `json:"template_hint"`
	MaxBatchPhones int    `json:"max_batch_phones"`
}
