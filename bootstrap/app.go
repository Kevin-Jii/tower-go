package bootstrap

import (
	"fmt"
	"tower-go/config"
	"tower-go/middleware"
	"tower-go/utils/logging"
	"tower-go/utils/session"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() {
	closeLogger := InitLogger()
	defer closeLogger()

	// 自动生成 swagger 文档（开发模式）
	GenerateSwaggerDocs()

	LoadAppConfig()
	InitDatabase()
	closeRedis := InitRedisCache()
	defer closeRedis()

	AutoMigrateAndSeeds()

	session.InitSessionManager("single", 3)
	logging.LogInfo("会话管理初始化完成")

	r := gin.Default()
	r.Use(middleware.RequestLoggerMiddleware(4096))

	controllers := BuildControllers()
	RegisterRoutes(r, controllers)

	// 初始化 Stream 客户端
	InitStreamClients(controllers.DingTalkBotModule)
	defer CloseStreamClients()

	addr := fmt.Sprintf(":%d", config.GetConfig().App.Port)
	if err := r.Run(addr); err != nil {
		logging.LogFatal("服务启动失败", zap.Error(err))
	}
}
