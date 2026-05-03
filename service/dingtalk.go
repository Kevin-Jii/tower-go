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
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
)

type DingTalkService struct {
	botModule    *module.DingTalkBotModule
	userModule   *module.DingTalkUserModule
	streamClient *DingTalkStreamClient
}

func NewDingTalkService(botModule *module.DingTalkBotModule, userModule *module.DingTalkUserModule) *DingTalkService {
	return &DingTalkService{
		botModule:    botModule,
		userModule:   userModule,
		streamClient: GetStreamClient(),
	}
}

// GetStreamClient 获取 Stream 客户端
func (s *DingTalkService) GetStreamClient() *DingTalkStreamClient {
	return s.streamClient
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

// SendTextToBot 发送文本消息到指定机器人配置
func (s *DingTalkService) SendTextToBot(bot *model.DingTalkBot, content string) error {
	msg := model.DingTalkTextMessage{
		MsgType: "text",
	}
	msg.Text.Content = content
	return s.sendMessage(bot, msg)
}

// SendMarkdownToBot 发送 Markdown 消息到指定机器人配置
func (s *DingTalkService) SendMarkdownToBot(bot *model.DingTalkBot, title, text string) error {
	msg := model.DingTalkMarkdownMessage{
		MsgType: "markdown",
	}
	msg.Markdown.Title = title
	msg.Markdown.Text = text
	return s.sendMessage(bot, msg)
}

// SendActionCardToBot Webhook 模式发送钉钉整体跳转卡片
func (s *DingTalkService) SendActionCardToBot(bot *model.DingTalkBot, title, text, singleTitle, singleURL string) error {
	msg := model.DingTalkActionCardMessage{
		MsgType: "actionCard",
	}
	msg.ActionCard.Title = title
	msg.ActionCard.Text = text
	msg.ActionCard.SingleTitle = singleTitle
	msg.ActionCard.SingleURL = singleURL
	msg.ActionCard.BtnOrientation = "0"
	return s.sendMessage(bot, msg)
}

// SendStreamText Stream 模式发送文本消息
func (s *DingTalkService) SendStreamText(bot *model.DingTalkBot, content string) error {
	return s.SendStreamTextToMobile(bot, content, "")
}

// SendStreamTextToMobile Stream 模式发送文本消息到指定手机号用户
func (s *DingTalkService) SendStreamTextToMobile(bot *model.DingTalkBot, content, mobile string) error {
	if bot.RobotCode == "" {
		return errors.New("robotCode is required for stream mode")
	}

	// 获取 access_token
	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// 获取用户ID
	var userIds []string
	if mobile != "" {
		userId, err := s.GetUserIdByMobile(mobile, accessToken)
		if err != nil {
			return fmt.Errorf("failed to get userId by mobile %s: %w", mobile, err)
		}
		userIds = []string{userId}
	}

	// 构造消息体 - 使用 sampleText 模板
	msgBody := map[string]interface{}{
		"msgtype": "sampleText",
		"content": content,
	}

	return s.sendStreamMessageToUsers(bot.RobotCode, accessToken, msgBody, userIds)
}

// SendStreamMarkdown Stream 模式发送 Markdown 消息
// 注意：单聊支持 Markdown
func (s *DingTalkService) SendStreamMarkdown(bot *model.DingTalkBot, title, text string) error {
	return s.SendStreamMarkdownToMobile(bot, title, text, "")
}

// SendStreamMarkdownToMobile Stream 模式发送 Markdown 消息到指定手机号用户
func (s *DingTalkService) SendStreamMarkdownToMobile(bot *model.DingTalkBot, title, text, mobile string) error {
	if bot.RobotCode == "" {
		return errors.New("robotCode is required for stream mode")
	}

	// 获取 access_token
	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// 获取用户ID
	var userIds []string
	if mobile != "" {
		userId, err := s.GetUserIdByMobile(mobile, accessToken)
		if err != nil {
			return fmt.Errorf("failed to get userId by mobile %s: %w", mobile, err)
		}
		userIds = []string{userId}
	}

	// 单聊消息支持 Markdown
	msgBody := map[string]interface{}{
		"msgtype": "sampleMarkdown",
		"title":   title,
		"text":    text,
	}

	return s.sendStreamMessageToUsers(bot.RobotCode, accessToken, msgBody, userIds)
}

// SendStreamCardToMobile Stream 模式发送卡片消息到指定手机号用户
// cardMsgKey 为钉钉卡片模板标识，cardParam 会被序列化为 msgParam
func (s *DingTalkService) SendStreamCardToMobile(bot *model.DingTalkBot, cardMsgKey, mobile string, cardParam map[string]interface{}) error {
	if bot.RobotCode == "" {
		return errors.New("robotCode is required for stream mode")
	}
	if strings.TrimSpace(cardMsgKey) == "" {
		return errors.New("card msg key is required")
	}

	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	var userIds []string
	if mobile != "" {
		userId, err := s.GetUserIdByMobile(mobile, accessToken)
		if err != nil {
			return fmt.Errorf("failed to get userId by mobile %s: %w", mobile, err)
		}
		userIds = []string{userId}
	}

	msgBody := map[string]interface{}{
		"msgtype": strings.TrimSpace(cardMsgKey),
	}
	for k, v := range cardParam {
		msgBody[k] = v
	}
	return s.sendStreamMessageToUsers(bot.RobotCode, accessToken, msgBody, userIds)
}

// SendStreamActionCardToMobile Stream 模式发送钉钉内置整体跳转卡片
func (s *DingTalkService) SendStreamActionCardToMobile(bot *model.DingTalkBot, title, text, singleTitle, singleURL, mobile string) error {
	if bot.RobotCode == "" {
		return errors.New("robotCode is required for stream mode")
	}

	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	var userIds []string
	if mobile != "" {
		userId, err := s.GetUserIdByMobile(mobile, accessToken)
		if err != nil {
			return fmt.Errorf("failed to get userId by mobile %s: %w", mobile, err)
		}
		userIds = []string{userId}
	}

	msgBody := map[string]interface{}{
		"msgtype":        "sampleActionCard",
		"title":          title,
		"text":           text,
		"singleTitle":    singleTitle,
		"singleURL":      singleURL,
		"btnOrientation": "0",
	}
	return s.sendStreamMessageToUsers(bot.RobotCode, accessToken, msgBody, userIds)
}

// SendStreamImageToMobile Stream 模式发送图片消息到指定手机号用户
// imageURL: 图片的公开访问URL
func (s *DingTalkService) SendStreamImageToMobile(bot *model.DingTalkBot, imageURL, mobile string) error {
	if bot.RobotCode == "" {
		return errors.New("robotCode is required for stream mode")
	}

	// 获取 access_token
	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// 获取用户ID
	var userIds []string
	if mobile != "" {
		userId, err := s.GetUserIdByMobile(mobile, accessToken)
		if err != nil {
			return fmt.Errorf("failed to get userId by mobile %s: %w", mobile, err)
		}
		userIds = []string{userId}
	}

	// 图片消息
	msgBody := map[string]interface{}{
		"msgtype":  "sampleImageMsg",
		"photoURL": imageURL,
	}

	return s.sendStreamMessageToUsers(bot.RobotCode, accessToken, msgBody, userIds)
}

// SendStreamMarkdownWithImageToMobile 发送带图片的Markdown消息（先发图片，再发文字）
func (s *DingTalkService) SendStreamMarkdownWithImageToMobile(bot *model.DingTalkBot, title, text, imageURL, mobile string) error {
	// 先发送图片
	if imageURL != "" {
		if err := s.SendStreamImageToMobile(bot, imageURL, mobile); err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Warnw("Failed to send image, will send text only", "error", err)
			}
		}
	}

	// 再发送文字消息
	return s.SendStreamMarkdownToMobile(bot, title, text, mobile)
}

// convertMarkdownToPlainText 将 Markdown 格式转换为纯文本
func convertMarkdownToPlainText(markdown string) string {
	text := markdown

	// 先移除 Markdown 标记（在处理图片之前）
	text = strings.ReplaceAll(text, "## ", "")
	text = strings.ReplaceAll(text, "**", "")
	text = strings.ReplaceAll(text, "- ", "• ")

	// 处理图片链接：![alt](url) -> 📷 查看报菜图片: url
	// 使用正则表达式或简单字符串替换
	for strings.Contains(text, "![") {
		start := strings.Index(text, "![")
		if start == -1 {
			break
		}

		// 找到对应的 ](
		bracketStart := strings.Index(text[start:], "](")
		if bracketStart == -1 {
			break
		}
		bracketStart += start

		// 找到最后的 )
		urlEnd := strings.Index(text[bracketStart+2:], ")")
		if urlEnd == -1 {
			break
		}
		urlEnd += bracketStart + 2

		// 提取 URL
		url := text[bracketStart+2 : urlEnd]

		// 替换整个图片 Markdown 为纯文本链接
		replacement := "\n\n📷 查看报菜图片:\n" + url + "\n"
		text = text[:start] + replacement + text[urlEnd+1:]

		// 调试日志
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("Converted image markdown to plain text",
				"url", url,
				"replacement", replacement)
		}
	}

	// 移除剩余的星号
	text = strings.ReplaceAll(text, "*", "")

	return text
}

// saveImageToNginx 保存图片到 nginx 托管目录，返回图片访问 URL
func (s *DingTalkService) saveImageToNginx(imageData []byte, filename string) (string, error) {
	imageURL, err := utils.SaveImageFile(filename, imageData)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}
	return imageURL, nil
}

// sendStreamMarkdownWithText 发送 Markdown 消息（Stream 模式，使用已有的 accessToken）
// 注意：钉钉群消息 API 限制较多，直接发送纯文本格式
func (s *DingTalkService) sendStreamMarkdownWithText(bot *model.DingTalkBot, title, markdownText, accessToken string) error {
	// Stream 模式群消息不支持 Markdown，直接发送纯文本
	// 保留图片链接，用户可以点击访问
	plainText := fmt.Sprintf("%s\n\n%s", title, convertMarkdownToPlainText(markdownText))

	// 使用 sampleText 消息类型（钉钉群消息API要求）
	// 注意：msgParam 的格式应该直接是 {"content": "文本"}，而不是嵌套 text 对象
	msgBody := map[string]interface{}{
		"msgtype": "sampleText",
		"content": plainText,
	}

	return s.sendStreamMessage(bot.RobotCode, accessToken, msgBody)
}

// SendStreamImageText Stream 模式发送图文消息
// 新方案：将图片保存到 nginx 托管目录，通过 Markdown 引用图片 URL
func (s *DingTalkService) SendStreamImageText(bot *model.DingTalkBot, title, text string, imageData []byte) error {
	if bot.RobotCode == "" {
		return errors.New("robotCode is required for stream mode")
	}

	// 获取 access_token
	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// 保存图片到 nginx 托管目录并获取 URL
	imageURL, err := s.saveImageToNginx(imageData, "menu_report.png")
	if err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to save image to nginx, sending text only",
				"botID", bot.ID,
				"error", err)
		}
		// 图片保存失败，降级为纯文本消息
		return s.SendStreamMarkdown(bot, title, text)
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Image saved successfully",
			"botID", bot.ID,
			"imageURL", imageURL,
			"imageSize", len(imageData))
	}

	// 在 Markdown 文本中添加图片
	markdownWithImage := fmt.Sprintf("%s\n\n![报菜明细](%s)", text, imageURL)

	// 发送 Markdown 消息（包含图片）
	return s.sendStreamMarkdownWithText(bot, title, markdownWithImage, accessToken)
}

// uploadImage 上传图片到钉钉,返回 mediaId (旧版API，保留向后兼容)
func (s *DingTalkService) uploadImage(accessToken string, imageData []byte) (string, error) {
	apiURL := "https://oapi.dingtalk.com/media/upload?access_token=" + accessToken + "&type=image"

	// 创建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件字段
	part, err := writer.CreateFormFile("media", "image.png")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// 写入图片数据
	if _, err := part.Write(imageData); err != nil {
		return "", fmt.Errorf("failed to write image data: %w", err)
	}

	// 关闭 writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 发送请求
	resp, err := http.Post(apiURL, writer.FormDataContentType(), body)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		return "", fmt.Errorf("dingtalk upload error: code=%v, msg=%v", errCode, result["errmsg"])
	}

	if mediaID, ok := result["media_id"].(string); ok {
		return mediaID, nil
	}

	return "", fmt.Errorf("failed to get media_id from response: %v", result)
}

// uploadImageMedia 上传图片到钉钉媒体库(新版API)，返回 mediaId
func (s *DingTalkService) uploadImageMedia(accessToken string, imageData []byte) (string, error) {
	apiURL := "https://oapi.dingtalk.com/media/upload?access_token=" + accessToken + "&type=image"

	// 创建 multipart form
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建文件字段
	part, err := writer.CreateFormFile("media", "menu_report.png")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// 写入图片数据
	if _, err := part.Write(imageData); err != nil {
		return "", fmt.Errorf("failed to write image data: %w", err)
	}

	// 关闭 writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(apiURL, writer.FormDataContentType(), body)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Upload image response",
			"statusCode", resp.StatusCode,
			"response", result)
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		return "", fmt.Errorf("dingtalk upload error: code=%v, msg=%v", errCode, result["errmsg"])
	}

	if mediaID, ok := result["media_id"].(string); ok {
		return mediaID, nil
	}

	return "", fmt.Errorf("failed to get media_id from response: %v", result)
}

// sendAnnouncement 发送钉钉企业公告
// 参考文档: https://open.dingtalk.com/document/orgapp/create-a-dingtalk-notification
func (s *DingTalkService) sendAnnouncement(accessToken string, agentIDStr, title, content, mediaID string) error {
	apiURL := "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=" + accessToken

	// 将 agentID 从字符串转换为数字
	var agentID int64
	if agentIDStr != "" {
		var err error
		agentID, err = strconv.ParseInt(agentIDStr, 10, 64)
		if err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Warnw("Invalid agentID, using default 0",
					"agentIDStr", agentIDStr,
					"error", err)
			}
			agentID = 0
		}
	}

	// 构建公告消息体
	msgContent := map[string]interface{}{
		"msgtype": "oa",
		"oa": map[string]interface{}{
			"message_url": "dingtalk://dingtalkclient/page/link?url=https://www.dingtalk.com",
			"head": map[string]string{
				"bgcolor": "FFBBBBBB",
				"text":    title,
			},
			"body": map[string]interface{}{
				"title":   title,
				"content": convertMarkdownToPlainText(content),
				"image":   "@lALPDfJ6V_FPDmvNAfTNAfQ", // 这是固定的图片占位符，实际图片通过下面的 form 字段传递
				"form": []map[string]string{
					{
						"key":   "详细信息",
						"value": convertMarkdownToPlainText(content),
					},
				},
			},
		},
	}

	// 如果有图片，添加图片字段
	if mediaID != "" {
		msgContent["oa"].(map[string]interface{})["body"].(map[string]interface{})["image"] = mediaID
	}

	msgJSON, _ := json.Marshal(msgContent)

	reqBody := map[string]interface{}{
		"agent_id":    agentID, // 企业内部应用的AgentId
		"msg":         string(msgJSON),
		"to_all_user": true, // 发送给所有人
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Sending announcement to DingTalk",
			"url", apiURL,
			"agentID", agentID,
			"title", title,
			"mediaID", mediaID,
			"requestBody", string(jsonData))
	}

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send announcement: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Announcement response from DingTalk",
			"statusCode", resp.StatusCode,
			"response", result)
	}

	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		return fmt.Errorf("dingtalk announcement error: code=%v, msg=%v", errCode, result["errmsg"])
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Announcement sent successfully",
			"title", title,
			"taskId", result["task_id"])
	}

	return nil
}

// getStreamAccessToken 获取 Stream 模式的 access_token
func (s *DingTalkService) getStreamAccessToken(clientID, clientSecret string) (string, error) {
	apiURL := "https://api.dingtalk.com/v1.0/oauth2/accessToken"

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

// GetUserIdByMobile 通过手机号获取钉钉用户ID（先查缓存，没有再调API）
func (s *DingTalkService) GetUserIdByMobile(mobile, accessToken string) (string, error) {
	if mobile == "" {
		return "", errors.New("mobile is empty")
	}

	// 1. 先从缓存表查询
	if s.userModule != nil {
		user, err := s.userModule.GetByMobile(mobile)
		if err == nil && user != nil && user.UserID != "" {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Infow("Got user from cache", "mobile", mobile, "userId", user.UserID)
			}
			return user.UserID, nil
		}
	}

	// 2. 调用钉钉API获取
	apiURL := "https://oapi.dingtalk.com/topapi/v2/user/getbymobile?access_token=" + accessToken

	reqBody := map[string]string{
		"mobile": mobile,
	}
	jsonData, _ := json.Marshal(reqBody)

	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to call dingtalk api: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("GetUserByMobile response", "mobile", mobile, "response", result)
	}

	// 检查错误
	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		return "", fmt.Errorf("dingtalk api error: code=%v, msg=%v", errCode, result["errmsg"])
	}

	// 提取 userid
	resultData, ok := result["result"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format")
	}
	userId, ok := resultData["userid"].(string)
	if !ok || userId == "" {
		return "", fmt.Errorf("userid not found in response")
	}

	// 3. 保存到缓存表
	if s.userModule != nil {
		name := ""
		if n, ok := resultData["name"].(string); ok {
			name = n
		}
		s.userModule.Upsert(&model.DingTalkUser{
			Mobile: mobile,
			UserID: userId,
			Name:   name,
		})
	}

	return userId, nil
}

// sendStreamMessage 通过钉钉服务端 API 发送单聊消息
// 使用机器人发送单聊消息 API: https://open.dingtalk.com/document/orgapp/chatbots-send-one-on-one-chat-messages-in-batches
func (s *DingTalkService) sendStreamMessage(robotCode, accessToken string, msgBody map[string]interface{}) error {
	return s.sendStreamMessageToUsers(robotCode, accessToken, msgBody, nil)
}

// sendStreamMessageToUsers 发送单聊消息给指定用户
func (s *DingTalkService) sendStreamMessageToUsers(robotCode, accessToken string, msgBody map[string]interface{}, userIds []string) error {
	if len(userIds) == 0 {
		return fmt.Errorf("no userIds specified for single chat message")
	}

	// 使用单聊消息 API
	apiURL := "https://api.dingtalk.com/v1.0/robot/oToMessages/batchSend"

	// 构造完整请求体
	msgType := msgBody["msgtype"].(string)
	delete(msgBody, "msgtype")

	msgParamJSON, _ := json.Marshal(msgBody)

	reqBody := map[string]interface{}{
		"msgKey":    msgType,
		"msgParam":  string(msgParamJSON),
		"robotCode": robotCode,
		"userIds":   userIds,
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

	// 记录请求详情（调试用）
	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Sending stream message to DingTalk API",
			"url", apiURL,
			"robotCode", robotCode,
			"msgKey", msgType,
			"accessTokenPreview", accessToken[:10]+"...",
			"requestBody", string(jsonData),
		)
	}

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

	// 记录完整响应（调试用）
	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Received response from DingTalk API",
			"statusCode", resp.StatusCode,
			"responseBody", string(body),
			"responseMap", result,
		)
	}

	// 检查是否有错误
	if errCode, ok := result["code"].(string); ok && errCode != "" {
		return fmt.Errorf("dingtalk api error: code=%v, msg=%v", errCode, result["message"])
	}

	// 如果 code 不存在，但是 response 是其他格式，视为成功
	// 实际 DingTalk API 可能返回成功但不包含 code 字段
	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Stream message sent successfully",
			"robotCode", robotCode,
			"response", result,
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

	// 创建机器人（暂时使用临时名称）
	bot := &model.DingTalkBot{
		Name:         "临时名称", // 先用临时名称，创建后再更新
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
		CardMsgKey:   strings.TrimSpace(req.CardMsgKey),
		Remark:       req.Remark,
	}

	if req.IsEnabled != nil {
		bot.IsEnabled = *req.IsEnabled
	}
	if req.MsgType != "" {
		bot.MsgType = req.MsgType
	}
	if bot.MsgType == "card" && bot.BotType != "stream" {
		return nil, errors.New("card msg type only supports stream bot")
	}
	if bot.MsgType == "card" && strings.TrimSpace(bot.CardMsgKey) == "" {
		return nil, errors.New("card_msg_key is required when msg_type is card")
	}

	// 创建机器人
	if err := s.botModule.Create(bot); err != nil {
		return nil, err
	}

	// 生成机器人名称：门店名称_机器人ID
	botName := ""
	if req.Name != "" {
		// 如果前端提供了名称，使用前端提供的
		botName = req.Name
	} else {
		// 否则自动生成：门店名称_机器人ID
		if bot.StoreID != nil {
			store, err := s.botModule.GetStoreByID(*bot.StoreID)
			if err == nil && store != nil {
				botName = fmt.Sprintf("%s_机器人%d", store.Name, bot.ID)
			} else {
				botName = fmt.Sprintf("机器人%d", bot.ID)
			}
		} else {
			botName = fmt.Sprintf("全局机器人%d", bot.ID)
		}
	}

	// 更新机器人名称
	if err := s.botModule.UpdateName(bot.ID, botName); err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("Failed to update bot name",
				"botID", bot.ID,
				"error", err)
		}
	}
	bot.Name = botName

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

	targetBotType := bot.BotType
	if req.BotType != nil && *req.BotType != "" {
		targetBotType = *req.BotType
	}
	targetMsgType := bot.MsgType
	if req.MsgType != nil && *req.MsgType != "" {
		targetMsgType = *req.MsgType
	}
	targetCardMsgKey := strings.TrimSpace(bot.CardMsgKey)
	if req.CardMsgKey != nil {
		targetCardMsgKey = strings.TrimSpace(*req.CardMsgKey)
	}
	if targetMsgType == "card" && targetBotType != "stream" {
		return errors.New("card msg type only supports stream bot")
	}
	if targetMsgType == "card" && targetCardMsgKey == "" {
		return errors.New("card_msg_key is required when msg_type is card")
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

// GetOpenConversationIdByBotId 通过机器人ID和群号获取群会话ID
// 这个方法会尝试多种方式获取 openConversationId
func (s *DingTalkService) GetOpenConversationIdByBotId(botID uint, chatId string) (map[string]interface{}, error) {
	// 获取机器人配置
	bot, err := s.botModule.GetByID(botID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bot: %w", err)
	}

	if bot.BotType != "stream" {
		return nil, errors.New("only stream bots support this operation")
	}

	// 获取 access_token
	accessToken, err := s.getStreamAccessToken(bot.ClientID, bot.ClientSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	// 方法1: 尝试使用新版API (v1.0)
	apiURL := fmt.Sprintf("https://api.dingtalk.com/v1.0/im/conversations/%s", chatId)

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)
	req.Header.Set("Content-Type", "application/json")

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Requesting group conversation info",
			"botID", botID,
			"chatId", chatId,
			"apiURL", apiURL,
			"accessTokenPreview", accessToken[:10]+"...")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Get group info response",
			"botID", botID,
			"chatId", chatId,
			"statusCode", resp.StatusCode,
			"response", result)
	}

	// 检查新版API的错误格式
	if code, ok := result["code"].(string); ok && code != "" {
		errMsg := "unknown error"
		if msg, ok := result["message"].(string); ok {
			errMsg = msg
		}

		// 如果新版API失败，尝试旧版API
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Warnw("New API failed, trying old API",
				"code", code,
				"message", errMsg)
		}

		return s.getOpenConversationIdByOldAPI(bot, chatId, accessToken, botID)
	}

	// 尝试提取 openConversationId
	var openConversationId string
	if id, ok := result["openConversationId"].(string); ok {
		openConversationId = id
	}

	// 返回结果
	response := map[string]interface{}{
		"success":              openConversationId != "",
		"chat_id":              chatId,
		"bot_id":               botID,
		"open_conversation_id": openConversationId,
		"api_version":          "v1.0",
		"raw_response":         result,
	}

	if openConversationId == "" {
		response["suggestion"] = "API返回成功但未包含openConversationId，请通过Stream事件回调获取：\n" +
			"1. 在群聊中@机器人发送消息\n" +
			"2. 查看应用日志找到conversationId\n" +
			"3. 使用该ID更新机器人配置"
	} else {
		response["message"] = "成功获取openConversationId，可以使用此ID更新机器人配置"
	}

	return response, nil
}

// getOpenConversationIdByOldAPI 使用旧版API获取群会话ID
func (s *DingTalkService) getOpenConversationIdByOldAPI(bot *model.DingTalkBot, chatId, accessToken string, botID uint) (map[string]interface{}, error) {
	// 使用旧版API
	apiURL := fmt.Sprintf("https://oapi.dingtalk.com/chat/get?access_token=%s&chatid=%s", accessToken, chatId)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get group info (old API): %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Get group info response (old API)",
			"botID", botID,
			"chatId", chatId,
			"response", result)
	}

	// 检查错误码
	if errCode, ok := result["errcode"].(float64); ok && errCode != 0 {
		errMsg := "unknown error"
		if msg, ok := result["errmsg"].(string); ok {
			errMsg = msg
		}

		// 返回详细的错误信息和建议
		return map[string]interface{}{
			"success": false,
			"error":   fmt.Sprintf("钉钉API错误: code=%v, msg=%v", errCode, errMsg),
			"suggestion": "无法通过API直接获取openConversationId，请使用以下方法：\n" +
				"1. 在群聊中@机器人发送消息\n" +
				"2. 查看应用日志中的Stream事件回调\n" +
				"3. 从事件中提取conversationId字段\n" +
				"4. 使用该conversationId更新机器人配置",
			"chat_id":     chatId,
			"bot_id":      botID,
			"api_version": "oapi (old)",
		}, nil
	}

	// 尝试提取 openConversationId
	var openConversationId string
	if chatInfo, ok := result["chat_info"].(map[string]interface{}); ok {
		if id, ok := chatInfo["openConversationId"].(string); ok {
			openConversationId = id
		}
	}

	// 返回结果
	response := map[string]interface{}{
		"success":              openConversationId != "",
		"chat_id":              chatId,
		"bot_id":               botID,
		"open_conversation_id": openConversationId,
		"api_version":          "oapi (old)",
		"raw_response":         result,
	}

	if openConversationId == "" {
		response["suggestion"] = "API返回成功但未包含openConversationId，请通过Stream事件回调获取：\n" +
			"1. 在群聊中@机器人发送消息\n" +
			"2. 查看应用日志找到conversationId\n" +
			"3. 使用该ID更新机器人配置"
	} else {
		response["message"] = "成功获取openConversationId，可以使用此ID更新机器人配置"
	}

	return response, nil
}

// TestBot 测试机器人连接
func (s *DingTalkService) TestBot(id uint) error {
	bot, err := s.botModule.GetByID(id)
	if err != nil {
		return err
	}

	testMsg := fmt.Sprintf("🔔 机器人测试消息\n\n发送时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	if bot.BotType == "stream" {
		// 获取门店负责人手机号
		var mobile string
		if bot.StoreID != nil {
			store, err := s.botModule.GetStoreByID(*bot.StoreID)
			if err == nil && store != nil && store.Phone != "" {
				mobile = store.Phone
			}
		}

		if mobile == "" {
			return errors.New("robot_code is empty for stream bot: please set the store phone number (门店联系电话) for testing")
		}

		if bot.MsgType == "card" && strings.TrimSpace(bot.CardMsgKey) != "" {
			cardParam := map[string]interface{}{
				"title":      "机器人测试",
				"content":    testMsg,
				"createTime": time.Now().Format("2006-01-02 15:04:05"),
			}
			return s.SendStreamCardToMobile(bot, bot.CardMsgKey, mobile, cardParam)
		}
		if bot.MsgType == "markdown" {
			return s.SendStreamMarkdownToMobile(bot, "机器人测试", testMsg, mobile)
		}
		return s.SendStreamTextToMobile(bot, testMsg, mobile)
	}

	// webhook 模式
	if bot.MsgType == "markdown" {
		return s.SendMarkdownToBot(bot, "机器人测试", testMsg)
	}
	return s.SendTextToBot(bot, testMsg)
}
