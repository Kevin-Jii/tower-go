package model

import "time"

// MessageTemplate 消息模板
type MessageTemplate struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string    `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:模板编码"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null;comment:模板名称"`
	Title       string    `json:"title" gorm:"type:varchar(200);comment:消息标题模板"`
	Content     string    `json:"content" gorm:"type:text;not null;comment:消息内容模板"`
	Description string    `json:"description" gorm:"type:varchar(500);comment:模板说明"`
	Variables   string    `json:"variables" gorm:"type:text;comment:可用变量说明(JSON)"`
	IsEnabled   bool      `json:"is_enabled" gorm:"default:true;comment:是否启用"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MessageTemplate) TableName() string {
	return "message_templates"
}

// 预定义模板编码
const (
	TemplateStoreAccountCreated = "store_account_created" // 记账通知
	TemplateInventoryCreated    = "inventory_created"     // 入库通知
	TemplatePurchaseCreated     = "purchase_created"      // 采购通知

	// 钉钉机器人命令回复模板
	TemplateBotHelp           = "bot_help"            // 帮助菜单
	TemplateBotInventoryQuery = "bot_inventory_query" // 库存查询
	TemplateBotTodayAccount   = "bot_today_account"   // 今日记账
	TemplateBotTodayInventory = "bot_today_inventory" // 今日入库
	TemplateBotSearchResult   = "bot_search_result"   // 搜索结果
	TemplateBotUnknown        = "bot_unknown"         // 未知命令
)

// CreateMessageTemplateReq 创建消息模板请求
type CreateMessageTemplateReq struct {
	Code        string `json:"code" binding:"required,max=50"`
	Name        string `json:"name" binding:"required,max=100"`
	Title       string `json:"title" binding:"max=200"`
	Content     string `json:"content" binding:"required"`
	Description string `json:"description" binding:"max=500"`
	Variables   string `json:"variables"`
	IsEnabled   *bool  `json:"is_enabled"`
}

// UpdateMessageTemplateReq 更新消息模板请求
type UpdateMessageTemplateReq struct {
	Name        *string `json:"name" binding:"omitempty,max=100"`
	Title       *string `json:"title" binding:"omitempty,max=200"`
	Content     *string `json:"content"`
	Description *string `json:"description" binding:"omitempty,max=500"`
	Variables   *string `json:"variables"`
	IsEnabled   *bool   `json:"is_enabled"`
}
