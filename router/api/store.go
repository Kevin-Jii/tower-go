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
		stores.POST("", middleware.Permission("store:add"), c.Store.CreateStore)
		stores.GET("", middleware.Permission("store:list"), c.Store.ListStores)
		stores.GET("/all", middleware.Permission("store:list"), c.Store.ListAllStores)
		stores.GET("/:id", middleware.Permission("store:list"), c.Store.GetStore)
		stores.PUT("/:id", middleware.Permission("store:edit"), c.Store.UpdateStore)
		stores.PUT("/:id/third-party-account", middleware.Permission("store:menu"), c.Store.BindThirdPartyAccount)
		stores.DELETE("/:id", middleware.Permission("store:delete"), c.Store.DeleteStore)
	}
}
