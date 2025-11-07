package model

import "time"

// DingTalkBot 钉钉机器人配置
type DingTalkBot struct {
	ID           uint      `json:"id" gorm:"primarykey"`
	Name         string    `json:"name" gorm:"size:100;not null"`              // 机器人名称
	BotType      string    `json:"bot_type" gorm:"size:20;default:'webhook'"`  // 机器人类型: webhook, stream
	Webhook      string    `json:"webhook" gorm:"size:500"`                    // Webhook 地址（webhook 模式）
	Secret       string    `json:"secret" gorm:"size:500"`                     // 签名密钥（webhook 模式）
	ClientID     string    `json:"client_id" gorm:"size:200"`                  // AppKey/SuiteKey (stream 模式)
	ClientSecret string    `json:"client_secret" gorm:"size:500"`              // AppSecret/SuiteSecret (stream 模式)
	AgentID      string    `json:"agent_id" gorm:"size:50"`                    // 应用 AgentId (stream 模式推送消息用)
	StoreID      *uint     `json:"store_id" gorm:"index"`                      // 所属门店（null 表示全局）
	Store        *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"`  // 门店关联
	IsEnabled    bool      `json:"is_enabled" gorm:"default:true;index"`       // 是否启用
	MsgType      string    `json:"msg_type" gorm:"size:20;default:'markdown'"` // 消息类型: text, markdown
	Remark       string    `json:"remark" gorm:"type:text"`                    // 备注
	RobotCode    string    `json:"robot_code" gorm:"size:100"`                 // 钉钉机器人编码(robotCode)
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CreateDingTalkBotReq 创建钉钉机器人请求
type CreateDingTalkBotReq struct {
	Name         string `json:"name" binding:"required"`
	BotType      string `json:"bot_type" binding:"omitempty,oneof=webhook stream"` // webhook 或 stream (默认 webhook)
	Webhook      string `json:"webhook"`                                           // webhook 模式必填
	Secret       string `json:"secret"`                                            // webhook 模式可选
	ClientID     string `json:"client_id"`                                         // stream 模式必填
	ClientSecret string `json:"client_secret"`                                     // stream 模式必填
	AgentID      string `json:"agent_id"`                                          // stream 模式推送消息用
	StoreID      *uint  `json:"store_id"`
	IsEnabled    *bool  `json:"is_enabled"`
	MsgType      string `json:"msg_type"`
	Remark       string `json:"remark"`
	RobotCode    string `json:"robot_code"` // 钉钉机器人编码(robotCode)
}

// UpdateDingTalkBotReq 更新钉钉机器人请求
type UpdateDingTalkBotReq struct {
	Name         *string `json:"name" patch:"allowZero"`
	BotType      *string `json:"bot_type" patch:"allowZero"`
	Webhook      *string `json:"webhook" patch:"always"`
	Secret       *string `json:"secret" patch:"always"`
	ClientID     *string `json:"client_id" patch:"always"`
	ClientSecret *string `json:"client_secret" patch:"always"`
	AgentID      *string `json:"agent_id" patch:"always"`
	StoreID      *uint   `json:"store_id" patch:"always"`
	IsEnabled    *bool   `json:"is_enabled" patch:"always"`
	MsgType      *string `json:"msg_type" patch:"allowZero"`
	Remark       *string `json:"remark" patch:"always"`
	RobotCode    *string `json:"robot_code" patch:"always"`
}

// DingTalkTextMessage 钉钉文本消息
type DingTalkTextMessage struct {
	MsgType string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	At *struct {
		AtMobiles []string `json:"atMobiles,omitempty"`
		IsAtAll   bool     `json:"isAtAll,omitempty"`
	} `json:"at,omitempty"`
}

// DingTalkMarkdownMessage 钉钉 Markdown 消息
type DingTalkMarkdownMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
	At *struct {
		AtMobiles []string `json:"atMobiles,omitempty"`
		IsAtAll   bool     `json:"isAtAll,omitempty"`
	} `json:"at,omitempty"`
}
