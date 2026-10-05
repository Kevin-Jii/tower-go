package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterSmsCampaignRoutes(r *gin.RouterGroup, c *Controllers) {
	group := r.Group("/sms-campaigns")
	group.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		group.GET("/config", middleware.Permission("marketing:sms:list"), c.SmsCampaign.Config)
		group.GET("", middleware.Permission("marketing:sms:list"), c.SmsCampaign.List)
		group.GET("/:id", middleware.Permission("marketing:sms:list"), c.SmsCampaign.Get)
		group.GET("/:id/records", middleware.Permission("marketing:sms:list"), c.SmsCampaign.Records)
		group.POST("/:id/records/:record_id/retry", middleware.Permission("marketing:sms:send"), c.SmsCampaign.RetryRecord)
		group.POST("", middleware.Permission("marketing:sms:add"), c.SmsCampaign.Create)
		group.PUT("/:id", middleware.Permission("marketing:sms:edit"), c.SmsCampaign.Update)
		group.DELETE("/:id", middleware.Permission("marketing:sms:delete"), c.SmsCampaign.Delete)
		group.POST("/:id/send", middleware.Permission("marketing:sms:send"), c.SmsCampaign.Send)
		group.POST("/:id/cancel", middleware.Permission("marketing:sms:edit"), c.SmsCampaign.Cancel)
	}
}
