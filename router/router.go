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
	v1 := r.Group("/api/v1")

	// 注册各模块路由
	api.RegisterAuthRoutes(v1, c)
	api.RegisterUserRoutes(v1, c)
	api.RegisterRoleRoutes(v1)
	api.RegisterStoreRoutes(v1, c)
	api.RegisterMenuRoutes(v1, c)
	api.RegisterDingTalkRoutes(v1, c)
	api.RegisterSupplierRoutes(v1, c)
	api.RegisterPurchaseRoutes(v1, c)
	api.RegisterDictRoutes(v1, c)
	api.RegisterFileRoutes(v1, c.File)

	// WebSocket
	r.GET("/ws", controller.WebSocketHandler)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	addr := fmt.Sprintf(":%d", config.GetConfig().App.Port)
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/index.html", addr)
	fmt.Printf("📚 Swagger UI: %s\n\n", swaggerURL)
}
