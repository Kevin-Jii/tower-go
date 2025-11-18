package bootstrap

import (
	"fmt"
	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/Kevin-Jii/tower-go/utils/session"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() {
	closeLogger := InitLogger()
	defer closeLogger()

	LoadAppConfig()
	
	// 自动生成 swagger 文档（开发模式）- 移到配置加载之后
	// 可通过环境变量 SWAG_AUTO=0 禁用以加快启动速度
	GenerateSwaggerDocs()
	
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
