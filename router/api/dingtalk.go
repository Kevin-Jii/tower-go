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
			robots.POST("", c.DingTalkBot.CreateBot)
			robots.GET("", c.DingTalkBot.ListBots)
			robots.GET("/:id", c.DingTalkBot.GetBot)
			robots.PUT("/:id", c.DingTalkBot.UpdateBot)
			robots.DELETE("/:id", c.DingTalkBot.DeleteBot)
			robots.POST("/:id/test", c.DingTalkBot.TestBot)
			robots.POST("/:id/test-callback", c.DingTalkBot.TestStreamBotCallback)
		}
	}

}
