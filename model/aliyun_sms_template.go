package model

import "time"

// AliyunSmsTemplate 阿里云短信模板（本地跟踪记录，用于审核状态轮询与活动引用）。
// 实际审核状态以阿里云为准；本表缓存 TemplateCode、Name、Content、Type、最近一次拉取的审核结果。
type AliyunSmsTemplate struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	TemplateCode    string    `json:"template_code" gorm:"type:varchar(64);uniqueIndex;not null"`
	Name            string    `json:"name" gorm:"type:varchar(120);not null"`
	Content         string    `json:"content" gorm:"type:text;not null"`
	TemplateType    int32     `json:"template_type" gorm:"not null;default:1;comment:0=验证码 1=通知 2=推广 3=国际"`
	RelatedSign     string    `json:"related_sign" gorm:"type:varchar(64);comment:关联签名（可选）"`
	Remark          string    `json:"remark" gorm:"type:varchar(500)"`
	AuditStatus     string    `json:"audit_status" gorm:"type:varchar(20);not null;default:pending;comment:pending|approved|rejected|cancelled|unknown"`
	AuditReason     string    `json:"audit_reason" gorm:"type:varchar(500)"`
	SourceCreatedBy uint      `json:"source_created_by" gorm:"comment:创建人"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (AliyunSmsTemplate) TableName() string { return "aliyun_sms_templates" }

const (
	SmsTemplateAuditPending   = "pending"
	SmsTemplateAuditApproved  = "approved"
	SmsTemplateAuditRejected  = "rejected"
	SmsTemplateAuditCancelled = "cancelled"
	SmsTemplateAuditUnknown   = "unknown"

	SmsTemplateTypeVerification  int32 = 0
	SmsTemplateTypeNotification  int32 = 1
	SmsTemplateTypePromotion     int32 = 2
	SmsTemplateTypeInternational int32 = 3
)

type CreateAliyunSmsTemplateReq struct {
	Name         string `json:"name" binding:"required,max=120"`
	Content      string `json:"content" binding:"required,max=500"`
	TemplateType *int32 `json:"template_type" binding:"required,oneof=0 1 2 3"`
	RelatedSign  string `json:"related_sign" binding:"max=64"`
	Remark       string `json:"remark" binding:"max=500"`
	OwnerStoreID uint   `json:"owner_store_id"`
}

type ListAliyunSmsTemplateReq struct {
	Keyword     string `form:"keyword"`
	AuditStatus string `form:"audit_status"`
}

type RefreshAliyunSmsTemplateReq struct {
	TemplateCode string `json:"template_code" binding:"required,max=64"`
}
