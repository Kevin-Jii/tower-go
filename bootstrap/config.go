package bootstrap

import (
	"os"
	"tower-go/config"
	"tower-go/utils"

	"go.uber.org/zap"
)

func LoadAppConfig() {
	// 初始化环境变量配置
	config.InitConfig()

	// 调试：检查环境变量加载
	config.DebugConfig()

	// 环境变量端口覆盖
	if portEnv := os.Getenv("APP_PORT"); portEnv != "" {
		cfg := config.GetConfig()
		if cfg.App.Port > 0 {
			utils.LogInfo("使用环境变量端口", zap.Int("port", cfg.App.Port))
		}
	}
	utils.LogInfo("配置加载完成")
}
