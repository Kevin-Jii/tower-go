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
		group.GET("", middleware.Permission("message:template:list"), c.MessageTemplate.List)
		group.GET("/:id", middleware.Permission("message:template:list"), c.MessageTemplate.Get)
		group.POST("", middleware.Permission("message:template:add"), c.MessageTemplate.Create)
		group.PUT("/:id", middleware.Permission("message:template:edit"), c.MessageTemplate.Update)
		group.DELETE("/:id", middleware.Permission("message:template:delete"), c.MessageTemplate.Delete)
	}
}
