package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterAuditLogRoutes(v1 *gin.RouterGroup, c *Controllers) {
	logs := v1.Group("/audit-logs")
	logs.Use(middleware.AuthMiddleware())
	{
		logs.GET("", middleware.Permission("system:audit-log:list"), c.AuditLog.List)
		logs.GET("/:id", middleware.Permission("system:audit-log:detail"), c.AuditLog.Get)
	}
}
