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

// SmsCampaign is a store-owned SMS promotion. Legacy template fields remain as
// the fallback when no segments have been configured.
type SmsCampaign struct {
	ID              uint                 `json:"id" gorm:"primaryKey;autoIncrement"`
	OwnerStoreID    uint                 `json:"owner_store_id" gorm:"not null;default:0;index"`
	Name            string               `json:"name" gorm:"type:varchar(120);not null"`
	CampaignType    string               `json:"campaign_type" gorm:"type:varchar(20);not null;index"`
	SignName        string               `json:"sign_name" gorm:"type:varchar(64)"`
	TemplateCode    string               `json:"template_code" gorm:"type:varchar(64);not null"`
	TemplateParam   string               `json:"template_param" gorm:"type:text"`
	PersonalizeName bool                 `json:"personalize_name" gorm:"default:false"`
	TargetType      string               `json:"target_type" gorm:"type:varchar(32);not null"`
	StoreIDs        UintList             `json:"store_ids" gorm:"type:json"`
	CustomPhones    StringList           `json:"custom_phones" gorm:"type:json"`
	ScheduledAt     *time.Time           `json:"scheduled_at" gorm:"index"`
	Status          string               `json:"status" gorm:"type:varchar(20);not null;index;default:draft"`
	TotalCount      int                  `json:"total_count" gorm:"default:0"`
	SuccessCount    int                  `json:"success_count" gorm:"default:0"`
	FailCount       int                  `json:"fail_count" gorm:"default:0"`
	LastError       string               `json:"last_error" gorm:"type:varchar(500)"`
	SentAt          *time.Time           `json:"sent_at"`
	CreatedBy       uint                 `json:"created_by"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
	Segments        []SmsCampaignSegment `json:"segments" gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE"`
}

func (SmsCampaign) TableName() string { return "sms_campaigns" }

// SmsCampaignSegment is evaluated by Position. The first tag-matching segment
// wins; a default segment is used only when no tagged segment matches.
type SmsCampaignSegment struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CampaignID      uint      `json:"campaign_id" gorm:"not null;index"`
	Position        int       `json:"position" gorm:"not null;default:0;index"`
	TagIDs          UintList  `json:"tag_ids" gorm:"type:json"`
	SignName        string    `json:"sign" gorm:"type:varchar(64)"`
	TemplateCode    string    `json:"template" gorm:"type:varchar(64);not null"`
	TemplateParam   string    `json:"params" gorm:"type:text"`
	PersonalizeName bool      `json:"personalize" gorm:"default:false"`
	IsDefault       bool      `json:"default" gorm:"not null;default:false"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (SmsCampaignSegment) TableName() string { return "sms_campaign_segments" }

type SmsSendRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CampaignID   uint      `json:"campaign_id" gorm:"index;not null"`
	SegmentID    *uint     `json:"segment_id" gorm:"index"`
	Phone        string    `json:"phone" gorm:"type:varchar(20);index;not null"`
	MemberID     *uint     `json:"member_id"`
	BizID        string    `json:"biz_id" gorm:"type:varchar(64)"`
	Status       string    `json:"status" gorm:"type:varchar(20);not null"`
	ErrorMessage string    `json:"error_message" gorm:"type:varchar(500)"`
	CreatedAt    time.Time `json:"created_at"`
}

func (SmsSendRecord) TableName() string { return "sms_send_records" }

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

type SmsCampaignSegmentReq struct {
	TagIDs          UintList `json:"tag_ids"`
	SignName        string   `json:"sign" binding:"max=64"`
	TemplateCode    string   `json:"template" binding:"required,max=64"`
	TemplateParam   string   `json:"params"`
	PersonalizeName bool     `json:"personalize"`
	IsDefault       bool     `json:"default"`
}

type CreateSmsCampaignReq struct {
	OwnerStoreID    uint                    `json:"owner_store_id"`
	Name            string                  `json:"name" binding:"required,max=120"`
	CampaignType    string                  `json:"campaign_type" binding:"required,oneof=activity holiday"`
	SignName        string                  `json:"sign_name" binding:"max=64"`
	TemplateCode    string                  `json:"template_code" binding:"max=64"`
	TemplateParam   string                  `json:"template_param"`
	PersonalizeName *bool                   `json:"personalize_name"`
	TargetType      string                  `json:"target_type" binding:"required,oneof=all_members stores custom"`
	StoreIDs        UintList                `json:"store_ids"`
	CustomPhones    StringList              `json:"custom_phones"`
	ScheduledAt     *time.Time              `json:"scheduled_at"`
	Segments        []SmsCampaignSegmentReq `json:"segments"`
}

type UpdateSmsCampaignReq struct {
	Name             *string                  `json:"name" binding:"omitempty,max=120"`
	CampaignType     *string                  `json:"campaign_type" binding:"omitempty,oneof=activity holiday"`
	SignName         *string                  `json:"sign_name" binding:"omitempty,max=64"`
	TemplateCode     *string                  `json:"template_code" binding:"omitempty,max=64"`
	TemplateParam    *string                  `json:"template_param"`
	PersonalizeName  *bool                    `json:"personalize_name"`
	TargetType       *string                  `json:"target_type" binding:"omitempty,oneof=all_members stores custom"`
	StoreIDs         *UintList                `json:"store_ids"`
	CustomPhones     *StringList              `json:"custom_phones"`
	ScheduledAt      *time.Time               `json:"scheduled_at"`
	ClearScheduledAt bool                     `json:"clear_scheduled_at"`
	Segments         *[]SmsCampaignSegmentReq `json:"segments"`
}

type SmsServiceConfigResp struct {
	Enabled                bool   `json:"enabled"`
	DefaultSign            string `json:"default_sign_name"`
	Region                 string `json:"region"`
	Configured             bool   `json:"configured"`
	HelpURL                string `json:"help_url"`
	TemplateHint           string `json:"template_hint"`
	MaxBatchPhones         int    `json:"max_batch_phones"`
	Timezone               string `json:"timezone"`
	SendWindowStart        string `json:"send_window_start"`
	SendWindowEnd          string `json:"send_window_end"`
	SendWindowEndExclusive bool   `json:"send_window_end_exclusive"`
	QualificationHint      string `json:"qualification_hint,omitempty"`
}
