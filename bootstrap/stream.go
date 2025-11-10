package bootstrap

import (
	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/logging"
)

// InitStreamClients 初始化所有 Stream 模式的机器人连接
func InitStreamClients(botModule *module.DingTalkBotModule) {
	streamClient := service.GetStreamClient()

	// 从配置文件读取 Stream 配置
	streamConfig := config.GetDingTalkStreamConfig()

	// 如果配置了 Stream 参数,确保数据库中存在对应机器人
	if streamConfig.ClientID != "" && streamConfig.ClientSecret != "" {
		ensureStreamBotExists(botModule, streamConfig)
	}

	// 查询所有启用的 Stream 类型机器人
	bots, err := botModule.ListEnabledStreamBots()
	if err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to query stream bots", "error", err)
		}
		return
	}

	if len(bots) == 0 {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Info("No stream bots to initialize")
		}
		return
	}

	// 启动所有 Stream 机器人
	successCount := 0
	for _, bot := range bots {
		if err := streamClient.StartBot(bot); err != nil {
			if logging.SugaredLogger != nil {
				logging.SugaredLogger.Errorw("Failed to start stream bot",
					"botID", bot.ID,
					"botName", bot.Name,
					"error", err,
				)
			}
		} else {
			successCount++
		}
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("Stream clients initialized",
			"total", len(bots),
			"success", successCount,
		)
	}
}

// CloseStreamClients 关闭所有 Stream 连接
func CloseStreamClients() {
	streamClient := service.GetStreamClient()
	streamClient.StopAll()

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Info("All stream clients stopped")
	}
}

// ensureStreamBotExists 确保配置文件中的 Stream 机器人在数据库中存在
func ensureStreamBotExists(botModule *module.DingTalkBotModule, streamConfig config.DingTalkStreamConfig) {
	// 先检查是否已存在任意 Stream 类型的机器人
	bots, err := botModule.ListEnabledStreamBots()
	if err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to query stream bots", "error", err)
		}
		return
	}

	// 如果已经存在 Stream 类型机器人,跳过创建
	if len(bots) > 0 {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("⏭️ Stream bot already exists, skipping creation",
				"count", len(bots),
			)
		}
		return
	}

	// 检查是否已存在使用相同 ClientID 的机器人
	existingBot, err := botModule.FindByClientID(streamConfig.ClientID)

	// 如果找到了机器人,直接返回
	if err == nil && existingBot != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Infow("Stream bot already exists with same clientID",
				"botID", existingBot.ID,
				"botName", existingBot.Name,
			)
		}
		return
	}

	// 创建新的 Stream 机器人
	bot := &model.DingTalkBot{
		Name:         "全局Stream机器人(自动创建)",
		BotType:      "stream",
		ClientID:     streamConfig.ClientID,
		ClientSecret: streamConfig.ClientSecret,
		AgentID:      streamConfig.AgentID,
		StoreID:      nil, // 全局机器人
		IsEnabled:    true,
		MsgType:      "markdown",
		Remark:       "由系统根据 config.yaml 自动创建",
	}

	if err := botModule.Create(bot); err != nil {
		if logging.SugaredLogger != nil {
			logging.SugaredLogger.Errorw("Failed to create stream bot from config",
				"error", err,
			)
		}
		return
	}

	if logging.SugaredLogger != nil {
		logging.SugaredLogger.Infow("✅ Stream bot created from config",
			"botID", bot.ID,
			"botName", bot.Name,
			"clientID", bot.ClientID,
		)
	}
}
