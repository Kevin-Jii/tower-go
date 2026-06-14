package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	MeituanSuggestionStatusPending  = "pending"
	MeituanSuggestionStatusApproved = "approved"
	MeituanSuggestionStatusDone     = "done"
	MeituanSuggestionStatusIgnored  = "ignored"
)

type MeituanAIOperatorAccount struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID        uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	ShopName       string         `json:"shop_name" gorm:"type:varchar(100);not null;comment:美团店铺名称"`
	ShopID         string         `json:"shop_id" gorm:"type:varchar(100);index;comment:美团店铺ID"`
	LoginName      string         `json:"login_name" gorm:"type:varchar(100);comment:登录账号"`
	DeveloperID    string         `json:"developer_id" gorm:"type:varchar(100);comment:美团开放平台DeveloperId"`
	SignKey        string         `json:"sign_key,omitempty" gorm:"type:varchar(255);comment:美团开放平台SignKey"`
	AppAuthToken   string         `json:"app_auth_token,omitempty" gorm:"type:varchar(500);comment:美团门店授权Token"`
	BusinessID     int            `json:"business_id" gorm:"not null;default:2;comment:美团业务ID 2外卖"`
	APIVersion     string         `json:"api_version" gorm:"type:varchar(20);not null;default:'2';comment:美团开放平台版本"`
	APIBaseURL     string         `json:"api_base_url" gorm:"type:varchar(255);not null;default:'https://api-open-cater.meituan.com';comment:美团开放平台地址"`
	QueryOrderPath string         `json:"query_order_path" gorm:"type:varchar(120);not null;default:'/api/order/queryById';comment:订单详情接口路径"`
	AuthStatus     string         `json:"auth_status" gorm:"type:varchar(20);not null;default:'manual';comment:授权状态 manual/oauth/expired"`
	IsEnabled      bool           `json:"is_enabled" gorm:"not null;default:true;comment:是否启用"`
	LastImportedAt *time.Time     `json:"last_imported_at,omitempty" gorm:"comment:最后导入时间"`
	Remark         string         `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (MeituanAIOperatorAccount) TableName() string {
	return "meituan_ai_accounts"
}

type MeituanAIOrder struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID        uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	AccountID      uint           `json:"account_id" gorm:"not null;index;uniqueIndex:uk_meituan_ai_orders_account_order;comment:美团账号ID"`
	OrderNo        string         `json:"order_no" gorm:"type:varchar(100);not null;index;uniqueIndex:uk_meituan_ai_orders_account_order;comment:美团订单号"`
	OrderTime      time.Time      `json:"order_time" gorm:"not null;index;comment:下单时间"`
	CustomerName   string         `json:"customer_name" gorm:"type:varchar(100);comment:顾客昵称"`
	ProductSummary string         `json:"product_summary" gorm:"type:varchar(500);comment:商品摘要"`
	OriginalAmount float64        `json:"original_amount" gorm:"type:decimal(10,2);not null;default:0;comment:原价"`
	ActualAmount   float64        `json:"actual_amount" gorm:"type:decimal(10,2);not null;default:0;comment:实收"`
	DeliveryFee    float64        `json:"delivery_fee" gorm:"type:decimal(10,2);not null;default:0;comment:配送费"`
	PlatformFee    float64        `json:"platform_fee" gorm:"type:decimal(10,2);not null;default:0;comment:平台服务费"`
	RefundAmount   float64        `json:"refund_amount" gorm:"type:decimal(10,2);not null;default:0;comment:退款金额"`
	Status         string         `json:"status" gorm:"type:varchar(50);comment:订单状态"`
	StoreAccountID *uint          `json:"store_account_id,omitempty" gorm:"index;comment:关联门店记账ID"`
	ImportedAt     time.Time      `json:"imported_at"`
	RawJSON        string         `json:"raw_json" gorm:"type:longtext;comment:原始数据"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (MeituanAIOrder) TableName() string {
	return "meituan_ai_orders"
}

type MeituanAIReview struct {
	ID             uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID        uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	AccountID      uint           `json:"account_id" gorm:"not null;index;uniqueIndex:uk_meituan_ai_reviews_account_review;comment:美团账号ID"`
	ReviewID       string         `json:"review_id" gorm:"type:varchar(100);index;uniqueIndex:uk_meituan_ai_reviews_account_review;comment:评价ID"`
	OrderNo        string         `json:"order_no" gorm:"type:varchar(100);index;comment:订单号"`
	Rating         int            `json:"rating" gorm:"not null;default:0;comment:评分"`
	Content        string         `json:"content" gorm:"type:text;comment:评价内容"`
	Sentiment      string         `json:"sentiment" gorm:"type:varchar(20);comment:情绪 positive/neutral/negative"`
	Tags           string         `json:"tags" gorm:"type:varchar(255);comment:标签逗号分隔"`
	SuggestedReply string         `json:"suggested_reply" gorm:"type:text;comment:建议回复"`
	ReviewTime     time.Time      `json:"review_time" gorm:"not null;index;comment:评价时间"`
	ReplyStatus    string         `json:"reply_status" gorm:"type:varchar(20);default:'pending';comment:回复状态"`
	ImportedAt     time.Time      `json:"imported_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (MeituanAIReview) TableName() string {
	return "meituan_ai_reviews"
}

type MeituanAISuggestion struct {
	ID            uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID       uint           `json:"store_id" gorm:"not null;index;comment:门店ID"`
	AccountID     uint           `json:"account_id" gorm:"not null;index;comment:美团账号ID"`
	Type          string         `json:"type" gorm:"type:varchar(30);not null;index;comment:建议类型"`
	Title         string         `json:"title" gorm:"type:varchar(120);not null;comment:标题"`
	Content       string         `json:"content" gorm:"type:text;comment:建议内容"`
	Reason        string         `json:"reason" gorm:"type:text;comment:原因"`
	ImpactScore   int            `json:"impact_score" gorm:"not null;default:0;comment:影响分"`
	Status        string         `json:"status" gorm:"type:varchar(20);not null;default:'pending';index;comment:状态"`
	ActionPayload string         `json:"action_payload" gorm:"type:text;comment:执行参数JSON"`
	GeneratedAt   time.Time      `json:"generated_at"`
	ApprovedAt    *time.Time     `json:"approved_at,omitempty"`
	DoneAt        *time.Time     `json:"done_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}

func (MeituanAISuggestion) TableName() string {
	return "meituan_ai_suggestions"
}

type CreateMeituanAIAccountReq struct {
	StoreID        uint   `json:"store_id"`
	ShopName       string `json:"shop_name" binding:"required,max=100"`
	ShopID         string `json:"shop_id" binding:"max=100"`
	LoginName      string `json:"login_name" binding:"max=100"`
	DeveloperID    string `json:"developer_id" binding:"max=100"`
	SignKey        string `json:"sign_key" binding:"max=255"`
	AppAuthToken   string `json:"app_auth_token" binding:"max=500"`
	BusinessID     int    `json:"business_id"`
	APIVersion     string `json:"api_version" binding:"max=20"`
	APIBaseURL     string `json:"api_base_url" binding:"max=255"`
	QueryOrderPath string `json:"query_order_path" binding:"max=120"`
	IsEnabled      *bool  `json:"is_enabled"`
	Remark         string `json:"remark" binding:"max=500"`
}

type UpdateMeituanAIAccountReq struct {
	ShopName       string `json:"shop_name" binding:"max=100"`
	ShopID         string `json:"shop_id" binding:"max=100"`
	LoginName      string `json:"login_name" binding:"max=100"`
	DeveloperID    string `json:"developer_id" binding:"max=100"`
	SignKey        string `json:"sign_key" binding:"max=255"`
	AppAuthToken   string `json:"app_auth_token" binding:"max=500"`
	BusinessID     int    `json:"business_id"`
	APIVersion     string `json:"api_version" binding:"max=20"`
	APIBaseURL     string `json:"api_base_url" binding:"max=255"`
	QueryOrderPath string `json:"query_order_path" binding:"max=120"`
	IsEnabled      *bool  `json:"is_enabled"`
	Remark         string `json:"remark" binding:"max=500"`
}

type ImportMeituanAIOrdersReq struct {
	Orders []ImportMeituanAIOrderItem `json:"orders" binding:"required,min=1,dive"`
}

type SyncMeituanAIOrdersResp struct {
	Imported int `json:"imported"`
	Skipped  int `json:"skipped"`
}

type SyncMeituanAIOpenAPIOrdersReq struct {
	OrderID  string   `json:"order_id"`
	OrderIDs []string `json:"order_ids"`
}

type ImportMeituanAIOrderItem struct {
	OrderNo        string  `json:"order_no" binding:"required,max=100"`
	OrderTime      string  `json:"order_time" binding:"required"`
	CustomerName   string  `json:"customer_name" binding:"max=100"`
	ProductSummary string  `json:"product_summary" binding:"max=500"`
	OriginalAmount float64 `json:"original_amount" binding:"gte=0"`
	ActualAmount   float64 `json:"actual_amount" binding:"gte=0"`
	DeliveryFee    float64 `json:"delivery_fee" binding:"gte=0"`
	PlatformFee    float64 `json:"platform_fee" binding:"gte=0"`
	RefundAmount   float64 `json:"refund_amount" binding:"gte=0"`
	Status         string  `json:"status" binding:"max=50"`
	RawJSON        string  `json:"raw_json"`
}

type ImportMeituanAIReviewsReq struct {
	Reviews []ImportMeituanAIReviewItem `json:"reviews" binding:"required,min=1,dive"`
}

type ImportMeituanAIReviewItem struct {
	ReviewID   string `json:"review_id" binding:"max=100"`
	OrderNo    string `json:"order_no" binding:"max=100"`
	Rating     int    `json:"rating" binding:"required,min=1,max=5"`
	Content    string `json:"content" binding:"max=2000"`
	ReviewTime string `json:"review_time" binding:"required"`
}

type ListMeituanAIReq struct {
	StoreID   uint   `form:"store_id"`
	AccountID uint   `form:"account_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Keyword   string `form:"keyword"`
	Page      int    `form:"page,default=1" binding:"min=1"`
	PageSize  int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

type UpdateMeituanAISuggestionStatusReq struct {
	Status string `json:"status" binding:"required,oneof=pending approved done ignored"`
}

type MeituanAIDashboard struct {
	OrderCount         int64    `json:"order_count"`
	SalesAmount        float64  `json:"sales_amount"`
	RefundAmount       float64  `json:"refund_amount"`
	PlatformFee        float64  `json:"platform_fee"`
	AvgOrderAmount     float64  `json:"avg_order_amount"`
	ReviewCount        int64    `json:"review_count"`
	NegativeCount      int64    `json:"negative_count"`
	NegativeRate       float64  `json:"negative_rate"`
	PendingSuggestions int64    `json:"pending_suggestions"`
	TopProducts        []string `json:"top_products"`
}
