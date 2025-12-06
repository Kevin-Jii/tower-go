package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterStatisticsRoutes 注册统计路由
func RegisterStatisticsRoutes(v1 *gin.RouterGroup, c *Controllers) {
	stats := v1.Group("/statistics")
	stats.Use(middleware.AuthMiddleware())
	{
		stats.GET("/dashboard", c.Statistics.Dashboard)
		stats.GET("/inventory", c.Statistics.InventoryStats)
		stats.GET("/sales", c.Statistics.SalesStats)
		stats.GET("/sales-trend", c.Statistics.SalesTrend)
		stats.GET("/channel", c.Statistics.ChannelStats)
	}
}
