package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterStoreSmsConfigRoutes(v1 *gin.RouterGroup, c *Controllers) {
	group := v1.Group("/store-sms-configs")
	group.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		group.GET("", middleware.Permission("marketing:sms:list"), c.StoreSmsConfig.Get)
		group.PUT("", middleware.Permission("marketing:sms:edit"), c.StoreSmsConfig.Upsert)
		group.POST("/test", middleware.Permission("marketing:sms:edit"), c.StoreSmsConfig.Test)
	}
}
