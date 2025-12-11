package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterMessageTemplateRoutes 注册消息模板路由
func RegisterMessageTemplateRoutes(r *gin.RouterGroup, c *Controllers) {
	group := r.Group("/message-templates")
	group.Use(middleware.AuthMiddleware())
	{
		group.GET("", c.MessageTemplate.List)
		group.GET("/:id", c.MessageTemplate.Get)
		group.POST("", c.MessageTemplate.Create)
		group.PUT("/:id", c.MessageTemplate.Update)
		group.DELETE("/:id", c.MessageTemplate.Delete)
	}
}
