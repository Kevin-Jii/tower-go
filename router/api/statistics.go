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
		stats.GET("/dashboard", middleware.Permission("statistics:dashboard"), c.Statistics.Dashboard)
		stats.GET("/inventory", middleware.Permission("statistics:dashboard"), c.Statistics.InventoryStats)
		stats.GET("/sales", middleware.Permission("statistics:dashboard"), c.Statistics.SalesStats)
		stats.GET("/sales-trend", middleware.Permission("statistics:dashboard"), c.Statistics.SalesTrend)
		stats.GET("/channel", middleware.Permission("statistics:dashboard"), c.Statistics.ChannelStats)
		stats.GET("/business-overview", middleware.Permission("statistics:dashboard"), c.Statistics.BusinessOverview)
		stats.GET("/home-charts", middleware.Permission("statistics:dashboard"), c.Statistics.HomeCharts)
	}
}
