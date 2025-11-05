package bootstrap

import (
	"fmt"
	"os"
	"tower-go/config"
	"tower-go/utils"

	"go.uber.org/zap"
)

func LoadAppConfig() {
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		utils.LogFatal("配置文件加载失败", zap.Error(err))
	}
	// 环境变量端口覆盖
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		var p int
		if _, err := fmt.Sscanf(portEnv, "%d", &p); err == nil && p > 0 {
			cfg := config.GetConfig()
			cfg.App.Port = p
			utils.LogInfo("使用环境变量端口", zap.Int("port", p))
		}
	}
	utils.LogInfo("配置加载完成")
}
