package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAliyunSmsTemplateRoutes(v1 *gin.RouterGroup, c *Controllers) {
	group := v1.Group("/sms-templates")
	group.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		group.GET("", middleware.Permission("marketing:sms:list"), c.AliyunSmsTemplate.List)
		group.GET("/approved", middleware.Permission("marketing:sms:list"), c.AliyunSmsTemplate.ListApproved)
		group.GET("/signatures", middleware.Permission("marketing:sms:list"), c.AliyunSmsTemplate.ListSignatures)
		group.POST("", middleware.Permission("marketing:sms:add"), c.AliyunSmsTemplate.Create)
		group.POST("/refresh", middleware.Permission("marketing:sms:edit"), c.AliyunSmsTemplate.Refresh)
		group.DELETE("/:code", middleware.Permission("marketing:sms:delete"), c.AliyunSmsTemplate.Delete)
	}
}
