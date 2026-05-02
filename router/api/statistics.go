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
		// 统计页对门店内登录用户开放，不再做单独权限码校验
		stats.GET("/dashboard", c.Statistics.Dashboard)
		stats.GET("/inventory", c.Statistics.InventoryStats)
		stats.GET("/sales", c.Statistics.SalesStats)
		stats.GET("/sales-trend", c.Statistics.SalesTrend)
		stats.GET("/channel", c.Statistics.ChannelStats)
		stats.GET("/business-overview", c.Statistics.BusinessOverview)
		stats.GET("/home-charts", c.Statistics.HomeCharts)
	}
}
