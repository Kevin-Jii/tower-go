package controller

import (
	"strconv"
	"time"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

type MenuReportController struct {
	menuReportService *service.MenuReportService
}

func NewMenuReportController(menuReportService *service.MenuReportService) *MenuReportController {
	return &MenuReportController{menuReportService: menuReportService}
}

// CreateMenuReportOrder godoc
// @Summary 创建报菜记录单
// @Description 创建新的报菜记录单（包含报菜详情列表，自动关联当前门店和操作员）
// @Tags menu-reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param order body model.CreateMenuReportOrderReq true "报菜记录单信息"
// @Success 200 {object} http.Response{data=model.MenuReportOrder}
// @Router /menu-reports [post]
func (c *MenuReportController) CreateMenuReportOrder(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	var req model.CreateMenuReportOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	order, err := c.menuReportService.CreateMenuReportOrder(storeID, userID, &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, order)
}

// GetMenuReportOrder godoc
// @Summary 获取报菜记录单详情
// @Description 获取报菜记录单详细信息（包含报菜详情列表）
// @Tags menu-reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "报菜记录单ID"
// @Success 200 {object} http.Response{data=model.MenuReportOrder}
// @Router /menu-reports/{id} [get]
func (c *MenuReportController) GetMenuReportOrder(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu report order ID")
		return
	}

	order, err := c.menuReportService.GetMenuReportOrder(uint(id), storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, order)
}

// ListMenuReportOrders godoc
// @Summary 报菜记录单列表
// @Description 获取当前门店的报菜记录单列表，支持分页和日期范围筛选
// @Tags menu-reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} http.Response{data=[]model.MenuReportOrder} "分页 meta: total,page,page_size,page_count,has_more"
// @Router /menu-reports [get]
func (c *MenuReportController) ListMenuReportOrders(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	// 如果提供了日期范围，则按日期查询
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			http.Error(ctx, 400, "Invalid start_date format")
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			http.Error(ctx, 400, "Invalid end_date format")
			return
		}

		orders, err := c.menuReportService.ListMenuReportOrdersByDateRange(storeID, startDate, endDate)
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		http.Success(ctx, orders)
		return
	}

	// 否则返回分页列表
	page := http.GetPage(ctx)
	pageSize := http.GetPageSize(ctx)

	orders, total, err := c.menuReportService.ListMenuReportOrders(storeID, page, pageSize)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, orders, total, page, pageSize)
}

// UpdateMenuReportOrder godoc
// @Summary 更新报菜记录单
// @Description 更新报菜记录单信息
// @Tags menu-reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "报菜记录单ID"
// @Param order body model.UpdateMenuReportOrderReq true "报菜记录单信息"
// @Success 200 {object} http.Response
// @Router /menu-reports/{id} [put]
func (c *MenuReportController) UpdateMenuReportOrder(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu report order ID")
		return
	}

	var req model.UpdateMenuReportOrderReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.menuReportService.UpdateMenuReportOrder(uint(id), storeID, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteMenuReportOrder godoc
// @Summary 删除报菜记录单
// @Description 删除报菜记录单（会同时删除关联的报菜详情）
// @Tags menu-reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "报菜记录单ID"
// @Success 200 {object} http.Response
// @Router /menu-reports/{id} [delete]
func (c *MenuReportController) DeleteMenuReportOrder(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu report order ID")
		return
	}

	if err := c.menuReportService.DeleteMenuReportOrder(uint(id), storeID); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetStatistics godoc
// @Summary 获取统计数据
// @Description 获取指定日期范围的报菜统计数量
// @Tags menu-reports
// @Accept json
// @Produce json
// @Security Bearer
// @Param start_date query string true "开始日期 (YYYY-MM-DD)"
// @Param end_date query string true "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} http.Response{data=model.MenuReportStats}
// @Router /menu-reports/statistics [get]
func (c *MenuReportController) GetStatistics(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	isAdmin := middleware.IsAdmin(ctx)

	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		http.Error(ctx, 400, "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		http.Error(ctx, 400, "Invalid start_date format")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		http.Error(ctx, 400, "Invalid end_date format")
		return
	}

	var stats *model.MenuReportStats

	// 如果是总部管理员，可以查看所有门店的统计
	if isAdmin {
		stats, err = c.menuReportService.GetStatsByDateRangeAllStores(startDate, endDate)
	} else {
		stats, err = c.menuReportService.GetStatsByDateRange(storeID, startDate, endDate)
	}

	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, stats)
}
