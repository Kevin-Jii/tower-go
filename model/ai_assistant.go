package model

import "time"

type AIAssistantConfig struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Provider     string    `json:"provider" gorm:"type:varchar(64);not null"`
	Model        string    `json:"model" gorm:"type:varchar(128);not null"`
	BaseURL      string    `json:"base_url" gorm:"type:varchar(512)"`
	APIKeyCipher string    `json:"-" gorm:"type:text"`
	Enabled      bool      `json:"enabled" gorm:"not null;default:false"`
	SystemPrompt string    `json:"system_prompt" gorm:"type:text"`
	ConfiguredBy uint      `json:"configured_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (AIAssistantConfig) TableName() string { return "ai_assistant_configs" }

type AIAssistantConversation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index:idx_ai_conversation_owner,priority:1"`
	StoreID   uint      `json:"store_id" gorm:"not null;index:idx_ai_conversation_owner,priority:2"`
	Title     string    `json:"title" gorm:"type:varchar(160);not null;default:'新对话'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (AIAssistantConversation) TableName() string { return "ai_assistant_conversations" }

type AIAssistantMessage struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	ConversationID  uint      `json:"conversation_id" gorm:"not null;index"`
	Role            string    `json:"role" gorm:"type:varchar(16);not null"`
	Content         string    `json:"content" gorm:"type:longtext;not null"`
	AnalysisContext string    `json:"analysis_context,omitempty" gorm:"type:longtext"`
	Provider        string    `json:"provider,omitempty" gorm:"type:varchar(64)"`
	Model           string    `json:"model,omitempty" gorm:"type:varchar(128)"`
	CreatedAt       time.Time `json:"created_at"`
}

func (AIAssistantMessage) TableName() string { return "ai_assistant_messages" }

type AIAssistantConfigView struct {
	Provider         string `json:"provider"`
	Model            string `json:"model"`
	BaseURL          string `json:"base_url"`
	Enabled          bool   `json:"enabled"`
	SystemPrompt     string `json:"system_prompt"`
	APIKeyConfigured bool   `json:"api_key_configured"`
	APIKeyMasked     string `json:"api_key_masked"`
}

type AIAssistantConfigRequest struct {
	Provider     string `json:"provider" binding:"required,max=64"`
	Model        string `json:"model" binding:"required,max=128"`
	BaseURL      string `json:"base_url" binding:"max=512"`
	APIKey       string `json:"api_key" binding:"max=4096"`
	Enabled      bool   `json:"enabled"`
	SystemPrompt string `json:"system_prompt" binding:"max=8000"`
}

type AIAssistantChatRequest struct {
	ConversationID uint   `json:"conversation_id"`
	StoreID        uint   `json:"store_id"`
	Message        string `json:"message" binding:"required,max=8000"`
	StartDate      string `json:"start_date"`
	EndDate        string `json:"end_date"`
}
