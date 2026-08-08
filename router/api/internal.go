package api

import (
	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterInternalRoutes registers endpoints intended only for trusted backend services.
func RegisterInternalRoutes(r *gin.Engine, c *Controllers) {
	internal := r.Group("/api/internal/v1")
	internal.Use(middleware.InternalServiceAuthMiddleware(config.GetConfig().InternalService.Token))
	{
		internal.GET("/reports/daily-turnover", c.DailyTurnover.List)
	}
}
