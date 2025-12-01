package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterStoreRoutes 注册门店管理路由
func RegisterStoreRoutes(v1 *gin.RouterGroup, c *Controllers) {
	stores := v1.Group("/stores")
	stores.Use(middleware.AuthMiddleware())
	{
		stores.POST("", c.Store.CreateStore)
		stores.GET("", c.Store.ListStores)
		stores.GET("/all", c.Store.ListAllStores)
		stores.GET("/:id", c.Store.GetStore)
		stores.PUT("/:id", c.Store.UpdateStore)
		stores.DELETE("/:id", c.Store.DeleteStore)
	}
}
