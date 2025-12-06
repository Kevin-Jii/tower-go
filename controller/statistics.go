package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
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
// @Tags statistics
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期: today/week/month/quarter/year" default(today)
// @Param store_id query int false "门店ID（管理员可指定）"
// @Success 200 {object} http.Response{data=model.DashboardStats}
// @Router /statistics/dashboard [get]
func (c *StatisticsController) Dashboard(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	// 管理员可以查看指定门店或全部
	queryStoreID := storeID
	if roleCode == model.RoleCodeAdmin || roleCode == model.RoleCodeSuperAdmin {
		if id := ctx.Query("store_id"); id != "" {
			var sid uint
			if _, err := ctx.GetQuery("store_id"); err {
				queryStoreID = 0 // 查看全部
			}
			if n, _ := ctx.GetQuery("store_id"); n != "" && n != "0" {
				http.ParseUint(n, &sid)
				queryStoreID = sid
			}
		} else {
			queryStoreID = 0 // 默认查看全部
		}
	}

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
// @Tags statistics
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=model.InventoryStats}
// @Router /statistics/inventory [get]
func (c *StatisticsController) InventoryStats(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	queryStoreID := storeID
	if roleCode == model.RoleCodeAdmin || roleCode == model.RoleCodeSuperAdmin {
		queryStoreID = 0
		if id := ctx.Query("store_id"); id != "" {
			var sid uint
			http.ParseUint(id, &sid)
			queryStoreID = sid
		}
	}

	stats, err := c.statisticsService.GetInventoryStats(queryStoreID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}

// SalesStats godoc
// @Summary 销售统计
// @Tags statistics
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期: today/week/month/quarter/year" default(today)
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=model.SalesStats}
// @Router /statistics/sales [get]
func (c *StatisticsController) SalesStats(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	queryStoreID := storeID
	if roleCode == model.RoleCodeAdmin || roleCode == model.RoleCodeSuperAdmin {
		queryStoreID = 0
		if id := ctx.Query("store_id"); id != "" {
			var sid uint
			http.ParseUint(id, &sid)
			queryStoreID = sid
		}
	}

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
// @Tags statistics
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期: week/month/quarter/year" default(month)
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=[]model.SalesTrendItem}
// @Router /statistics/sales-trend [get]
func (c *StatisticsController) SalesTrend(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	queryStoreID := storeID
	if roleCode == model.RoleCodeAdmin || roleCode == model.RoleCodeSuperAdmin {
		queryStoreID = 0
		if id := ctx.Query("store_id"); id != "" {
			var sid uint
			http.ParseUint(id, &sid)
			queryStoreID = sid
		}
	}

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
// @Tags statistics
// @Produce json
// @Security Bearer
// @Param period query string false "统计周期: today/week/month/quarter/year" default(month)
// @Param store_id query int false "门店ID"
// @Success 200 {object} http.Response{data=[]model.ChannelStatsItem}
// @Router /statistics/channel [get]
func (c *StatisticsController) ChannelStats(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	queryStoreID := storeID
	if roleCode == model.RoleCodeAdmin || roleCode == model.RoleCodeSuperAdmin {
		queryStoreID = 0
		if id := ctx.Query("store_id"); id != "" {
			var sid uint
			http.ParseUint(id, &sid)
			queryStoreID = sid
		}
	}

	period := ctx.DefaultQuery("period", "month")

	stats, err := c.statisticsService.GetChannelStats(queryStoreID, period)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}
