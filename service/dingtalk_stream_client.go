package service

import (
	"context"
	"fmt"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"sync"
	"time"

	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
)

// DingTalkStreamClient Stream 模式客户端管理器
type DingTalkStreamClient struct {
	clients map[uint]*client.StreamClient // botID -> StreamClient
	mu      sync.RWMutex
	running bool
}

var (
	globalStreamClient     *DingTalkStreamClient
	globalStreamClientOnce sync.Once
)

// GetStreamClient 获取全局 Stream 客户端实例
func GetStreamClient() *DingTalkStreamClient {
	globalStreamClientOnce.Do(func() {
		globalStreamClient = &DingTalkStreamClient{
			clients: make(map[uint]*client.StreamClient),
		}
	})
	return globalStreamClient
}

// StartBot 启动指定机器人的 Stream 连接
func (sc *DingTalkStreamClient) StartBot(bot *model.DingTalkBot) error {
	if bot.BotType != "stream" {
		return fmt.Errorf("bot type is not stream")
	}

	if bot.ClientID == "" || bot.ClientSecret == "" {
		return fmt.Errorf("clientID or clientSecret is empty")
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	// 如果已存在,先停止
	if existingClient, exists := sc.clients[bot.ID]; exists {
		existingClient.Close()
		delete(sc.clients, bot.ID)
	}

	// 配置日志
	logger.SetLogger(logger.NewStdTestLogger())

	// 创建 Stream 客户端
	streamClient := client.NewStreamClient(
		client.WithAppCredential(
			client.NewAppCredentialConfig(bot.ClientID, bot.ClientSecret),
		),
	)

	// 注册机器人消息回调(必须注册,否则连接会失败)
	streamClient.RegisterChatBotCallbackRouter(sc.OnChatBotMessageReceived)

	// 启动客户端
	go func() {
		if err := streamClient.Start(context.Background()); err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Errorw("Stream client start failed",
					"botID", bot.ID,
					"botName", bot.Name,
					"error", err,
				)
			}
		}
	}()

	sc.clients[bot.ID] = streamClient
	sc.running = true
	return nil
}

// StopBot 停止指定机器人的 Stream 连接
func (sc *DingTalkStreamClient) StopBot(botID uint) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if streamClient, exists := sc.clients[botID]; exists {
		streamClient.Close()
		delete(sc.clients, botID)

		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("Stream bot stopped", "botID", botID)
		}
		return nil
	}

	return fmt.Errorf("bot not found")
}

// StopAll 停止所有 Stream 连接
func (sc *DingTalkStreamClient) StopAll() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for botID, streamClient := range sc.clients {
		streamClient.Close()
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("Stream bot stopped", "botID", botID)
		}
	}

	sc.clients = make(map[uint]*client.StreamClient)
	sc.running = false
}

// IsRunning 检查 Stream 客户端是否正在运行
func (sc *DingTalkStreamClient) IsRunning() bool {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.running
}

// GetClient 获取指定机器人的 Stream 客户端
func (sc *DingTalkStreamClient) GetClient(botID uint) (*client.StreamClient, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	streamClient, exists := sc.clients[botID]
	return streamClient, exists
}

// GetBotCount 获取正在运行的机器人数量
func (sc *DingTalkStreamClient) GetBotCount() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return len(sc.clients)
}

// OnChatBotMessageReceived 处理机器人收到的消息回调
// 符合 chatbot.IChatBotMessageHandler 接口：func(context.Context, *BotCallbackDataModel) ([]byte, error)
func (sc *DingTalkStreamClient) OnChatBotMessageReceived(ctx context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	// 记录收到的消息
	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("📨 Received bot message",
			"conversationId", data.ConversationId,
			"senderStaffId", data.SenderStaffId,
			"senderNick", data.SenderNick,
			"text", data.Text.Content,
			"sessionWebhook", data.SessionWebhook != "",
		)
	}

	// TODO: 在这里添加你的业务逻辑
	// 例如：
	// 1. 解析用户消息
	// 2. 调用 AI 接口获取回复
	// 3. 处理业务逻辑

	// 使用 SessionWebhook 回复消息（如果需要）
	if data.SessionWebhook != "" {
		replier := chatbot.NewChatbotReplier()

		// 回复文本消息
		replyMsg := fmt.Sprintf("✅ 消息已收到\n你发送的内容是：%s\n\n时间：%s",
			data.Text.Content,
			time.Now().Format("2006-01-02 15:04:05"))

		if err := replier.SimpleReplyText(ctx, data.SessionWebhook, []byte(replyMsg)); err != nil {
			logging.SugaredLogger.Errorw("Failed to reply text message",
				"error", err,
			)
		}

		// 回复 Markdown 消息
		markdownContent := fmt.Sprintf("### 📨 消息处理完成\n\n**发送者：**@%s\n\n**内容：**\n%s\n\n**处理时间：** %s",
			data.SenderNick,
			data.Text.Content,
			time.Now().Format("2006-01-02 15:04:05"))

		if err := replier.SimpleReplyMarkdown(ctx, data.SessionWebhook,
			[]byte("消息处理结果"), []byte(markdownContent)); err != nil {
			logging.SugaredLogger.Errorw("Failed to reply markdown message",
				"error", err,
			)
		}
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("✅ Message processed successfully")
	}

	// 返回空字节数组（SDK 要求）
	return []byte(""), nil
}
