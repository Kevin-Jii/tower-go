package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type StatisticsController struct {
	statisticsService *service.StatisticsService
}

func NewStatisticsController(statisticsService *service.StatisticsService) *StatisticsController {
	return &StatisticsController{statisticsService: statisticsService}
}

// Dashboard godoc
// @Summary 统计面板数据
// @Description 获取库存和销售统计数据，支持当天/当周/当月/当季/当年
// @Tags 统计分析
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期: today/week/month/quarter/year" default(today)
// @Param store_id query int false "门店ID（管理员可指定）"
// @Success 200 {object} http.Response{data=model.DashboardStats}
// @Router /statistics/dashboard [get]
func (c *StatisticsController) Dashboard(ctx *gin.Context) {
	queryStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	period := ctx.DefaultQuery("period", "today")

	stats, err := c.statisticsService.GetDashboard(queryStoreID, period)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}

// InventoryStats godoc
// @Summary 库存统计
// @Tags 统计分析
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=model.InventoryStats}
// @Router /statistics/inventory [get]
func (c *StatisticsController) InventoryStats(ctx *gin.Context) {
	queryStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	stats, err := c.statisticsService.GetInventoryStats(queryStoreID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}

// SalesStats godoc
// @Summary 销售统计
// @Tags 统计分析
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期: today/week/month/quarter/year" default(today)
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=model.SalesStats}
// @Router /statistics/sales [get]
func (c *StatisticsController) SalesStats(ctx *gin.Context) {
	queryStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	period := ctx.DefaultQuery("period", "today")

	stats, err := c.statisticsService.GetSalesStats(queryStoreID, period)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}

// SalesTrend godoc
// @Summary 销售趋势（图表数据）
// @Tags 统计分析
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期: week/month/quarter/year" default(month)
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=[]model.SalesTrendItem}
// @Router /statistics/sales-trend [get]
func (c *StatisticsController) SalesTrend(ctx *gin.Context) {
	queryStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	period := ctx.DefaultQuery("period", "month")

	trend, err := c.statisticsService.GetSalesTrend(queryStoreID, period)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, trend)
}

// ChannelStats godoc
// @Summary 渠道销售统计
// @Tags 统计分析
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期: today/week/month/quarter/year" default(month)
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=[]model.ChannelStatsItem}
// @Router /statistics/channel [get]
func (c *StatisticsController) ChannelStats(ctx *gin.Context) {
	queryStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	period := ctx.DefaultQuery("period", "month")

	stats, err := c.statisticsService.GetChannelStats(queryStoreID, period)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}

// BusinessOverview godoc
// @Summary 经营总览统计
// @Tags 统计分析
// @Produce json
// @Security Bearer
// @Param start_date query string true "开始日期 YYYY-MM-DD"
// @Param end_date query string true "结束日期 YYYY-MM-DD"
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=model.BusinessOverviewStats}
// @Router /statistics/business-overview [get]
func (c *StatisticsController) BusinessOverview(ctx *gin.Context) {
	queryStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	stats, err := c.statisticsService.GetBusinessOverview(queryStoreID, startDate, endDate)
	if err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	http.Success(ctx, stats)
}

// HomeCharts godoc
// @Summary 首页图表统计
// @Tags 统计分析
// @Produce json
// @Security Bearer
// @Param start_date query string true "开始日期 YYYY-MM-DD"
// @Param end_date query string true "结束日期 YYYY-MM-DD"
// @Param granularity query string false "折线粒度 day/month" default(day)
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=model.HomeChartsStats}
// @Router /statistics/home-charts [get]
func (c *StatisticsController) HomeCharts(ctx *gin.Context) {
	queryStoreID := middleware.ResolveQueryStoreID(ctx, "store_id")
	startDate := ctx.Query("start_date")
	endDate := ctx.Query("end_date")
	granularity := ctx.DefaultQuery("granularity", "day")
	stats, err := c.statisticsService.GetHomeChartsStats(queryStoreID, startDate, endDate, granularity)
	if err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	http.Success(ctx, stats)
}
