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
		accounts.POST("", middleware.Permission("store:account:add"), c.StoreAccount.Create)
		accounts.GET("", middleware.Permission("store:account:list"), c.StoreAccount.List)
		accounts.GET("/stats", middleware.Permission("store:account:list"), c.StoreAccount.Stats)
		accounts.GET("/:id", middleware.Permission("store:account:list"), c.StoreAccount.Get)
		accounts.PUT("/:id", middleware.Permission("store:account:edit"), c.StoreAccount.Update)
		accounts.POST("/:id/consumables", middleware.Permission("store:account:edit"), c.StoreAccount.BindConsumables)
		accounts.DELETE("/:id", middleware.Permission("store:account:delete"), c.StoreAccount.Delete)
	}
}
