package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/encryption"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
)

type AIAssistantService struct {
	module     *module.AIAssistantModule
	statistics *StatisticsService
}

func NewAIAssistantService(m *module.AIAssistantModule, stats *StatisticsService) *AIAssistantService {
	return &AIAssistantService{module: m, statistics: stats}
}

func (s *AIAssistantService) GetConfig() (*model.AIAssistantConfigView, error) {
	row, err := s.module.Config()
	if err != nil {
		return nil, err
	}
	view := &model.AIAssistantConfigView{Provider: "openai", Model: "gpt-4o-mini"}
	if row == nil {
		return view, nil
	}
	view.Provider, view.Model, view.BaseURL, view.Enabled, view.SystemPrompt = row.Provider, row.Model, row.BaseURL, row.Enabled, row.SystemPrompt
	if row.APIKeyCipher != "" {
		secret, e := encryption.Decrypt(row.APIKeyCipher)
		if e == nil && secret != "" {
			view.APIKeyConfigured = true
			view.APIKeyMasked = encryption.MaskSecret(secret)
		}
	}
	return view, nil
}
func (s *AIAssistantService) SaveConfig(req *model.AIAssistantConfigRequest, userID uint) (*model.AIAssistantConfigView, error) {
	row, err := s.module.Config()
	if err != nil {
		return nil, err
	}
	if row == nil {
		row = &model.AIAssistantConfig{ID: 1}
	}
	row.Provider = strings.TrimSpace(req.Provider)
	row.Model = strings.TrimSpace(req.Model)
	row.BaseURL = strings.TrimSpace(req.BaseURL)
	row.Enabled = req.Enabled
	row.SystemPrompt = strings.TrimSpace(req.SystemPrompt)
	row.ConfiguredBy = userID
	if key := strings.TrimSpace(req.APIKey); key != "" {
		row.APIKeyCipher, err = encryption.Encrypt(key)
		if err != nil {
			return nil, fmt.Errorf("加密 API Key 失败: %w", err)
		}
	}
	if row.Enabled && row.APIKeyCipher == "" {
		return nil, errors.New("启用助手前请配置 API Key")
	}
	if err = s.module.SaveConfig(row); err != nil {
		return nil, err
	}
	return s.GetConfig()
}
func (s *AIAssistantService) TestConfig(req *model.AIAssistantConfigRequest) error {
	if strings.TrimSpace(req.APIKey) == "" {
		return errors.New("请填写 API Key 以测试连接")
	}
	client, err := openai.NewChatModel(context.Background(), &openai.ChatModelConfig{APIKey: strings.TrimSpace(req.APIKey), BaseURL: strings.TrimSpace(req.BaseURL), Model: strings.TrimSpace(req.Model)})
	if err != nil {
		return err
	}
	_, err = client.Generate(context.Background(), []*schema.Message{schema.UserMessage("Reply with OK")})
	return err
}
func (s *AIAssistantService) Conversations(userID, storeID uint) ([]model.AIAssistantConversation, error) {
	return s.module.Conversations(userID, storeID)
}
func (s *AIAssistantService) CreateConversation(userID, storeID uint) (*model.AIAssistantConversation, error) {
	row := &model.AIAssistantConversation{UserID: userID, StoreID: storeID, Title: "新对话"}
	err := s.module.CreateConversation(row)
	return row, err
}
func (s *AIAssistantService) Messages(id, userID, storeID uint) ([]model.AIAssistantMessage, error) {
	if _, err := s.module.Conversation(id, userID, storeID); err != nil {
		return nil, err
	}
	return s.module.Messages(id)
}
func (s *AIAssistantService) Rename(id, userID, storeID uint, title string) error {
	title = strings.TrimSpace(title)
	if title == "" || len([]rune(title)) > 160 {
		return errors.New("标题长度必须为 1 到 160 个字符")
	}
	return s.module.RenameConversation(id, userID, storeID, title)
}
func (s *AIAssistantService) Delete(id, userID, storeID uint) error {
	return s.module.DeleteConversation(id, userID, storeID)
}

func (s *AIAssistantService) Chat(ctx context.Context, req *model.AIAssistantChatRequest, userID, storeID uint) (*model.AIAssistantMessage, error) {
	cfg, err := s.module.Config()
	if err != nil {
		return nil, err
	}
	if cfg == nil || !cfg.Enabled || cfg.APIKeyCipher == "" {
		return nil, errors.New("AI 助手尚未配置或未启用")
	}
	key, err := encryption.Decrypt(cfg.APIKeyCipher)
	if err != nil || key == "" {
		return nil, errors.New("AI API Key 无法解密，请重新配置")
	}
	var conversation *model.AIAssistantConversation
	if req.ConversationID == 0 {
		conversation, err = s.CreateConversation(userID, storeID)
	} else {
		conversation, err = s.module.Conversation(req.ConversationID, userID, storeID)
	}
	if err != nil {
		return nil, err
	}
	text := strings.TrimSpace(req.Message)
	start, end := req.StartDate, req.EndDate
	if start == "" || end == "" {
		end = time.Now().Format("2006-01-02")
		start = time.Now().AddDate(0, 0, -30).Format("2006-01-02")
	}
	st, e1 := time.Parse("2006-01-02", start)
	et, e2 := time.Parse("2006-01-02", end)
	if e1 != nil || e2 != nil || et.Before(st) || et.Sub(st) > 366*24*time.Hour {
		return nil, errors.New("日期范围无效，最多支持 366 天")
	}
	contextData, err := s.businessContext(storeID, start, end)
	if err != nil {
		return nil, err
	}
	raw, _ := json.Marshal(contextData)
	messages := []*schema.Message{schema.SystemMessage(defaultSystemPrompt(cfg.SystemPrompt)), schema.SystemMessage(fmt.Sprintf("当前分析范围：%s 至 %s。只读业务数据快照：\n%s", start, end, string(raw)))}
	history, err := s.module.Messages(conversation.ID)
	if err != nil {
		return nil, err
	}
	if len(history) > 12 {
		history = history[len(history)-12:]
	}
	for _, item := range history {
		if item.Role == "user" {
			messages = append(messages, schema.UserMessage(item.Content))
		} else if item.Role == "assistant" {
			messages = append(messages, schema.AssistantMessage(item.Content, nil))
		}
	}
	messages = append(messages, schema.UserMessage(text))
	client, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{APIKey: key, BaseURL: cfg.BaseURL, Model: cfg.Model})
	if err != nil {
		return nil, err
	}
	answer, err := client.Generate(ctx, messages)
	if err != nil {
		return nil, err
	}
	userMsg := &model.AIAssistantMessage{ConversationID: conversation.ID, Role: "user", Content: text}
	if err = s.module.AddMessage(userMsg); err != nil {
		return nil, err
	}
	response := &model.AIAssistantMessage{ConversationID: conversation.ID, Role: "assistant", Content: answer.Content, AnalysisContext: string(raw), Provider: cfg.Provider, Model: cfg.Model}
	if err = s.module.AddMessage(response); err != nil {
		return nil, err
	}
	_ = s.module.TouchConversation(conversation.ID)
	_ = s.module.RenameIfUntitled(conversation.ID, truncateTitle(text, 60))
	return response, nil
}
func defaultSystemPrompt(prompt string) string {
	if strings.TrimSpace(prompt) != "" {
		return prompt
	}
	return "你是 Tower Go 的经营分析助手。只使用提供的数据，不得编造数字；指出数据口径、变化和可执行建议。涉及缺失数据时明确说明。回答请使用清晰的 Markdown：先给简短结论，再按‘分析过程’列出可核验的计算口径、数据对比与依据，复杂内容使用标题和列表；有行列数据时输出标准 Markdown 表格（表头行后必须紧跟分隔行，如 | --- | ---: |）；代码放入带语言标识的 fenced code block。不要输出隐藏的逐字内心推理或思维链，只提供简洁、可验证的分析依据和结论。"
}
func truncateTitle(v string, n int) string {
	r := []rune(strings.TrimSpace(v))
	if len(r) > n {
		r = r[:n]
	}
	if len(r) == 0 {
		return "新对话"
	}
	return string(r)
}
func (s *AIAssistantService) businessContext(storeID uint, start, end string) (map[string]any, error) {
	overview, err := s.statistics.GetBusinessOverview(storeID, start, end)
	if err != nil {
		return nil, err
	}
	inventory, err := s.statistics.GetInventoryStats(storeID)
	if err != nil {
		return nil, err
	}
	trend, err := s.statistics.BusinessModule().GetSalesTrend(storeID, start, end, "month")
	if err != nil {
		return nil, err
	}
	channels, err := s.statistics.BusinessModule().GetChannelStats(storeID, start, end)
	if err != nil {
		return nil, err
	}
	payments, err := s.module.PaymentSummary(storeID, start, end)
	if err != nil {
		return nil, err
	}
	return map[string]any{"overview": overview, "inventory": inventory, "sales_trend": trend, "channels": channels, "payments": payments}, nil
}
