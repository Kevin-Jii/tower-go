package bootstrap

import (
	"fmt"
	"net"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/logging"
	"github.com/Kevin-Jii/tower-go/utils/session"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() {
	const startupSteps = 15
	progress := newStartupProgress(startupSteps)
	progress.begin("初始化日志")
	closeLogger := InitLogger()
	defer closeLogger()
	progress.complete()

	progress.begin("加载应用配置")
	LoadAppConfig()
	progress.complete()

	progress.begin("生成并校验 API 文档")
	// 可通过环境变量 SWAG_AUTO=0 禁用以加快启动速度。
	GenerateSwaggerDocs()
	progress.complete()

	progress.begin("连接数据库")
	InitDatabase()
	progress.complete()

	progress.begin("连接 Redis 缓存")
	closeRedis := InitRedisCache()
	defer closeRedis()
	progress.complete()

	progress.begin("检查并执行数据库迁移")
	AutoMigrateAndSeeds()
	applyUserStorePatches(database.GetDB())
	progress.complete()

	progress.begin("初始化种子数据")
	RunSeedSQL()
	InitDefaultDicts() // 与 RunSeedSQL 共用 SKIP_SEED_DATA=1
	progress.complete()

	progress.begin("初始化事件订阅")
	InitEventSubscribers()
	progress.complete()

	progress.begin("初始化会话管理")
	session.InitSessionManager("single", 3)
	logging.LogInfo("会话管理初始化完成")
	progress.complete()

	progress.begin("初始化 HTTP 中间件")
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.RequestLoggerMiddleware(4096))
	r.Use(middleware.AuditLogMiddleware(4096))
	r.Use(middleware.DuplicateRequestMiddleware())
	progress.complete()

	progress.begin("构建控制器")
	controllers := BuildControllers()
	progress.complete()

	progress.begin("启动定时任务")
	if err := controllers.StartCronJobs(); err != nil {
		logging.LogWarn("定时任务启动失败: " + err.Error())
	}
	progress.complete()

	progress.begin("注册 API 路由")
	RegisterRoutes(r, controllers)
	progress.complete()

	progress.begin("初始化 Stream 客户端")
	InitStreamClients(controllers.DingTalkBotModule)
	defer CloseStreamClients()
	progress.complete()

	addr := fmt.Sprintf(":%d", config.GetConfig().App.Port)
	progress.begin("启动 HTTP 服务")
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logging.LogFatal("服务启动失败", zap.Error(err))
	}
	progress.complete()
	if err := r.RunListener(listener); err != nil {
		logging.LogFatal("服务启动失败", zap.Error(err))
	}
}
