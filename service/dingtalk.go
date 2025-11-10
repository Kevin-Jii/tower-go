package service

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
)

type DingTalkService struct {
	botModule    *module.DingTalkBotModule
	streamClient *DingTalkStreamClient
}

func NewDingTalkService(botModule *module.DingTalkBotModule) *DingTalkService {
	return &DingTalkService{
		botModule:    botModule,
		streamClient: GetStreamClient(),
	}
}

// SendTextMessage 发送文本消息到指定机器人
func (s *DingTalkService) SendTextMessage(botID uint, content string, atMobiles []string, isAtAll bool) error {
	bot, err := s.botModule.GetByID(botID)
	if err != nil {
		return fmt.Errorf("failed to get bot: %w", err)
	}

	if !bot.IsEnabled {
		return errors.New("bot is disabled")
	}

	msg := model.DingTalkTextMessage{
		MsgType: "text",
	}
	msg.Text.Content = content

	if len(atMobiles) > 0 || isAtAll {
		msg.At = &struct {
			AtMobiles []string `json:"atMobiles,omitempty"`
			IsAtAll   bool     `json:"isAtAll,omitempty"`
		}{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		}
	}

	return s.sendMessage(bot, msg)
}

// SendMarkdownMessage 发送 Markdown 消息到指定机器人
func (s *DingTalkService) SendMarkdownMessage(botID uint, title, text string, atMobiles []string, isAtAll bool) error {
	bot, err := s.botModule.GetByID(botID)
	if err != nil {
		return fmt.Errorf("failed to get bot: %w", err)
	}

	if !bot.IsEnabled {
		return errors.New("bot is disabled")
	}

	msg := model.DingTalkMarkdownMessage{
		MsgType: "markdown",
	}
	msg.Markdown.Title = title
	msg.Markdown.Text = text

	if len(atMobiles) > 0 || isAtAll {
		msg.At = &struct {
			AtMobiles []string `json:"atMobiles,omitempty"`
			IsAtAll   bool     `json:"isAtAll,omitempty"`
		}{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		}
	}

	return s.sendMessage(bot, msg)
}

// BroadcastToStore 广播消息到门店的所有机器人（支持 webhook 和 stream 双模式）
func (s *DingTalkService) BroadcastToStore(storeID uint, msgType, title, content string) error {
	bots, err := s.botModule.ListEnabledByStoreID(storeID)
	if err != nil {
		return fmt.Errorf("failed to list bots: %w", err)
	}

	if len(bots) == 0 {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("No enabled bots for store", "storeID", storeID)
		}
		return nil
	}

	var lastErr error
	for _, bot := range bots {
		var err error

		// 根据机器人类型选择发送方式
		if bot.BotType == "stream" {
			// Stream 模式：通过钉钉服务端 API 发送
			if msgType == "markdown" {
				err = s.sendStreamMarkdown(bot, title, content)
			} else {
				err = s.sendStreamText(bot, content)
			}
		} else {
			// Webhook 模式：直接 HTTP POST
			if msgType == "markdown" {
				err = s.sendMarkdownToBot(bot, title, content)
			} else {
				err = s.sendTextToBot(bot, content)
			}
		}

		if err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Errorw("Failed to send to bot",
					"botID", bot.ID,
					"botType", bot.BotType,
					"error", err,
				)
			}
			lastErr = err
		}
	}

	return lastErr
}

// sendTextToBot 发送文本消息到指定机器人配置
func (s *DingTalkService) sendTextToBot(bot *model.DingTalkBot, content string) error {
	msg := model.DingTalkTextMessage{
		MsgType: "text",
	}
	msg.Text.Content = content
	return s.sendMessage(bot, msg)
}

// sendMarkdownToBot 发送 Markdown 消息到指定机器人配置
func (s *DingTalkService) sendMarkdownToBot(bot *model.DingTalkBot, title, text string) error {
	msg := model.DingTalkMarkdownMessage{
		MsgType: "markdown",
	}
	msg.Markdown.Title = title
	msg.Markdown.Text = text
	return s.sendMessage(bot, msg)
}

// sendStreamText Stream 模式发送文本消息
func (s *DingTalkService) sendStreamText(bot *model.DingTalkBot, content string) error {
	if bot.RobotCode == "" {
		return errors.New("robotCode is required for stream mode")
	}

	// 获取 access_token
	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// 构造消息体 - 使用 sampleText 模板
	msgBody := map[string]interface{}{
		"msgtype": "sampleText",
		"text": map[string]string{
			"content": content,
		},
	}

	return s.sendStreamMessage(bot.RobotCode, accessToken, msgBody)
}

// sendStreamMarkdown Stream 模式发送 Markdown 消息
func (s *DingTalkService) sendStreamMarkdown(bot *model.DingTalkBot, title, text string) error {
	if bot.RobotCode == "" {
		return errors.New("robotCode is required for stream mode")
	}

	// 获取 access_token
	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// 构造消息体 - 使用 sampleMarkdown 模板
	msgBody := map[string]interface{}{
		"msgtype": "sampleMarkdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
	}

	return s.sendStreamMessage(bot.RobotCode, accessToken, msgBody)
}

// getStreamAccessToken 获取 Stream 模式的 access_token
func (s *DingTalkService) getStreamAccessToken(clientID, clientSecret string) (string, error) {
	apiURL := "https://api.dingtalk.com/v1.0/oauth2/accessToken"

	// Trim any accidental whitespace from config values
	clientID = strings.TrimSpace(clientID)
	clientSecret = strings.TrimSpace(clientSecret)

	reqBody := map[string]string{
		"appKey":    clientID,
		"appSecret": clientSecret,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if accessToken, ok := result["accessToken"].(string); ok {
		return accessToken, nil
	}

	// Log DingTalk response for debugging (do not log clientSecret)
	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Errorw("Failed to get access token from DingTalk",
			"clientID", clientID,
			"response", result,
		)
	}

	return "", fmt.Errorf("failed to get access token: %v", result)
}

var defaultStreamUserIds = []string{"010903622624-181076934"}

// sendStreamMessage 通过钉钉服务端 API 发送消息到群聊
// 使用机器人批量发送消息 API: https://open.dingtalk.com/document/orgapp/robot-send-group-messages
func (s *DingTalkService) sendStreamMessage(robotCode, accessToken string, msgBody map[string]interface{}) error {
	apiURL := "https://api.dingtalk.com/v1.0/robot/oToMessages/batchSend"

	// 构造完整请求体
	reqBody := map[string]interface{}{
		"msgKey": msgBody["msgtype"].(string),
		"msgParam": func() string {
			content := msgBody
			delete(content, "msgtype")
			jsonStr, _ := json.Marshal(content)
			return string(jsonStr)
		}(),
		"robotCode": robotCode,
		"userIds":   defaultStreamUserIds,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// 创建请求
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)

	// 发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// 解析响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查是否有错误
	if errCode, ok := result["code"].(string); ok && errCode != "0" {
		return fmt.Errorf("dingtalk api error: code=%v, msg=%v", errCode, result["message"])
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Stream message sent successfully",
			"robotCode", robotCode,
			"processQueryKey", result["processQueryKey"],
		)
	}

	return nil
}

// sendMessage 发送消息到钉钉（Webhook 模式通用方法）
func (s *DingTalkService) sendMessage(bot *model.DingTalkBot, message interface{}) error {
	webhook := bot.Webhook

	// 如果配置了签名密钥，则添加签名
	if bot.Secret != "" {
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		sign := s.generateSign(timestamp, bot.Secret)
		webhook = fmt.Sprintf("%s&timestamp=%s&sign=%s", webhook, timestamp, sign)
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	resp, err := http.Post(webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		return fmt.Errorf("dingtalk api error: code=%v, msg=%v", errCode, result["errmsg"])
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("DingTalk message sent successfully", "botID", bot.ID)
	}
	return nil
}

// generateSign 生成钉钉签名
func (s *DingTalkService) generateSign(timestamp, secret string) string {
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return url.QueryEscape(signature)
}

// CreateBot 创建机器人配置
func (s *DingTalkService) CreateBot(req *model.CreateDingTalkBotReq) (*model.DingTalkBot, error) {
	// 默认值处理
	if req.BotType == "" {
		req.BotType = "webhook" // 默认使用 webhook 模式
	}

	// 验证字段
	if req.BotType == "webhook" {
		if req.Webhook == "" {
			return nil, errors.New("webhook is required for webhook type")
		}
		// 检查 webhook 是否已存在
		exists, err := s.botModule.ExistsByWebhook(req.Webhook, 0)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("webhook already exists")
		}
	} else if req.BotType == "stream" {
		if req.ClientID == "" || req.ClientSecret == "" {
			return nil, errors.New("clientID and clientSecret are required for stream type")
		}
	}

	bot := &model.DingTalkBot{
		Name:         req.Name,
		BotType:      req.BotType,
		Webhook:      req.Webhook,
		Secret:       req.Secret,
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		AgentID:      req.AgentID,
		RobotCode:    req.RobotCode,
		StoreID:      req.StoreID,
		IsEnabled:    true,
		MsgType:      "markdown",
		Remark:       req.Remark,
	}

	if req.IsEnabled != nil {
		bot.IsEnabled = *req.IsEnabled
	}
	if req.MsgType != "" {
		bot.MsgType = req.MsgType
	}

	if err := s.botModule.Create(bot); err != nil {
		return nil, err
	}

	// 如果是 Stream 类型且启用,自动启动连接
	if bot.BotType == "stream" && bot.IsEnabled {
		if err := s.streamClient.StartBot(bot); err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Errorw("Failed to start stream bot after creation",
					"botID", bot.ID,
					"error", err,
				)
			}
		}
	}

	return bot, nil
}

// GetBot 获取机器人配置
func (s *DingTalkService) GetBot(id uint) (*model.DingTalkBot, error) {
	return s.botModule.GetByID(id)
}

// ListBots 获取机器人列表
func (s *DingTalkService) ListBots(page, pageSize int) ([]*model.DingTalkBot, int64, error) {
	return s.botModule.List(page, pageSize)
}

// UpdateBot 更新机器人配置
func (s *DingTalkService) UpdateBot(id uint, req *model.UpdateDingTalkBotReq) error {
	// 验证机器人是否存在
	bot, err := s.botModule.GetByID(id)
	if err != nil {
		return err
	}

	// 如果更新 webhook，检查是否重复
	if req.Webhook != nil && *req.Webhook != "" {
		exists, err := s.botModule.ExistsByWebhook(*req.Webhook, id)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("webhook already exists")
		}
	}

	updates := updatesPkg.BuildUpdatesFromReq(req)
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	// 记录旧的状态
	wasEnabled := bot.IsEnabled
	oldBotType := bot.BotType

	if err := s.botModule.Update(id, updates); err != nil {
		return err
	}

	// 重新获取更新后的机器人信息
	updatedBot, err := s.botModule.GetByID(id)
	if err != nil {
		return err
	}

	// Stream 机器人连接管理
	if updatedBot.BotType == "stream" {
		// 如果从禁用变为启用,启动连接
		if !wasEnabled && updatedBot.IsEnabled {
			if err := s.streamClient.StartBot(updatedBot); err != nil {
				if logging.SugaredLogger != nil {
					logging.SugaredLogger.Errorw("Failed to start stream bot",
						"botID", id,
						"error", err,
					)
				}
			}
		}
		// 如果从启用变为禁用,停止连接
		if wasEnabled && !updatedBot.IsEnabled {
			if err := s.streamClient.StopBot(id); err != nil {
				if logging.SugaredLogger != nil {
					logging.SugaredLogger.Warnw("Failed to stop stream bot",
						"botID", id,
						"error", err,
					)
				}
			}
		}
		// 如果类型改变或配置更新,重启连接
		if oldBotType == "stream" && updatedBot.IsEnabled {
			s.streamClient.StopBot(id) // 忽略错误
			if err := s.streamClient.StartBot(updatedBot); err != nil {
				if logging.SugaredLogger != nil {
					logging.SugaredLogger.Errorw("Failed to restart stream bot",
						"botID", id,
						"error", err,
					)
				}
			}
		}
	} else if oldBotType == "stream" {
		// 如果从 stream 改为 webhook,停止连接
		s.streamClient.StopBot(id) // 忽略错误
	}

	return nil
}

// DeleteBot 删除机器人配置
func (s *DingTalkService) DeleteBot(id uint) error {
	// 获取机器人信息
	bot, err := s.botModule.GetByID(id)
	if err != nil {
		return err
	}

	// 禁止删除 Stream 类型的机器人
	if bot.BotType == "stream" {
		return errors.New("不允许删除 Stream 类型的机器人,如需停用请修改启用状态")
	}

	return s.botModule.Delete(id)
}

// TestBot 测试机器人连接
func (s *DingTalkService) TestBot(id uint) error {
	bot, err := s.botModule.GetByID(id)
	if err != nil {
		return err
	}

	testMsg := fmt.Sprintf("🔔 机器人测试消息\n\n发送时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	if bot.BotType == "stream" {
		if bot.MsgType == "markdown" {
			return s.sendStreamMarkdown(bot, "机器人测试", testMsg)
		}
		return s.sendStreamText(bot, testMsg)
	}

	// webhook 模式
	if bot.MsgType == "markdown" {
		return s.sendMarkdownToBot(bot, "机器人测试", testMsg)
	}
	return s.sendTextToBot(bot, testMsg)
}
