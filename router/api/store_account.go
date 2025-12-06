package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterStoreAccountRoutes 注册门店记账路由
func RegisterStoreAccountRoutes(v1 *gin.RouterGroup, c *Controllers) {
	accounts := v1.Group("/store-accounts")
	accounts.Use(middleware.AuthMiddleware())
	{
		accounts.POST("", c.StoreAccount.Create)
		accounts.GET("", c.StoreAccount.List)
		accounts.GET("/stats", c.StoreAccount.Stats)
		accounts.GET("/:id", c.StoreAccount.Get)
		accounts.PUT("/:id", c.StoreAccount.Update)
		accounts.DELETE("/:id", c.StoreAccount.Delete)
	}
}
