package router

import (
	"fmt"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/controller"
	"github.com/Kevin-Jii/tower-go/router/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 初始化路由
func Setup(r *gin.Engine, c *api.Controllers) {
	// 初始化健康检查控制器
	healthController := controller.NewHealthController()

	// 注册健康检查路由（无需认证）
	r.GET("/health", healthController.Check)
	r.GET("/ready", healthController.Ready)
	r.GET("/live", healthController.Live)

	v1 := r.Group("/api/v1")

	// 注册各模块路由
	api.RegisterAuthRoutes(v1, c)
	api.RegisterUserRoutes(v1, c)
	api.RegisterRoleRoutes(v1)
	api.RegisterPermissionRoutes(v1, c)
	api.RegisterStoreRoutes(v1, c)
	api.RegisterMenuRoutes(v1, c)
	api.RegisterDingTalkRoutes(v1, c)
	api.RegisterSupplierRoutes(v1, c)
	api.RegisterPurchaseRoutes(v1, c)
	api.RegisterDictRoutes(v1, c)
	api.RegisterInventoryRoutes(v1, c)
	api.RegisterFileRoutes(v1, c.File)
	api.RegisterGalleryRoutes(v1, c.Gallery)
	api.RegisterStoreAccountRoutes(v1, c)
	api.RegisterStoreExpenseRoutes(v1, c)
	api.RegisterStoreReturnRoutes(v1, c)
	api.RegisterMeituanAIRoutes(v1, c)
	api.RegisterStatisticsRoutes(v1, c)
	api.RegisterMessageTemplateRoutes(v1, c)
	api.RegisterMemberRoutes(v1, c)
	api.RegisterPrinterRoutes(v1, c)
	api.RegisterPriceListRoutes(v1, c)
	api.RegisterB2BRoutes(v1, c)
	api.RegisterThirdPartyAccountRoutes(v1, c)
	api.RegisterThirdPartyRouteRoutes(v1, c)
	api.RegisterAuditLogRoutes(v1, c)

	// WebSocket
	r.GET("/ws", controller.WebSocketHandler)

	// Swagger - 保留原始JSON接口
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Scalar - 美化版API文档
	r.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html")
		c.String(200, `<!DOCTYPE html>
<html>
<head>
    <title>API 文档</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1"/>
</head>
<body>
    <script id="api-reference" data-url="/swagger/doc.json"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`)
	})

	addr := fmt.Sprintf(":%d", config.GetConfig().App.Port)
	fmt.Printf("📚 Swagger UI: http://localhost%s/swagger/index.html\n", addr)
	fmt.Printf("📚 Scalar Docs: http://localhost%s/docs\n\n", addr)
}
