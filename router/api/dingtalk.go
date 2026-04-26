package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterDingTalkRoutes 注册钉钉相关路由
func RegisterDingTalkRoutes(v1 *gin.RouterGroup, c *Controllers) {
	dingtalk := v1.Group("/dingtalk")
	dingtalk.Use(middleware.AuthMiddleware())
	{
		robots := dingtalk.Group("/robots")
		{
			robots.POST("", middleware.Permission("dingtalk:robot:add"), c.DingTalkBot.CreateBot)
			robots.GET("", middleware.Permission("dingtalk:robot:list"), c.DingTalkBot.ListBots)
			robots.GET("/:id", middleware.Permission("dingtalk:robot:list"), c.DingTalkBot.GetBot)
			robots.PUT("/:id", middleware.Permission("dingtalk:robot:edit"), c.DingTalkBot.UpdateBot)
			robots.DELETE("/:id", middleware.Permission("dingtalk:robot:delete"), c.DingTalkBot.DeleteBot)
			robots.POST("/:id/test", middleware.Permission("dingtalk:robot:test"), c.DingTalkBot.TestBot)
			robots.POST("/:id/test-callback", middleware.Permission("dingtalk:robot:test"), c.DingTalkBot.TestStreamBotCallback)
		}
	}

}
