package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	dingtalkMod "github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ReportBotController struct {
	module *dingtalkMod.MenuReportModule
}

func NewReportBotController(module *dingtalkMod.MenuReportModule) *ReportBotController {
	return &ReportBotController{module: module}
}

// CreateBot godoc
// @Summary 创建报菜机器人（创建报菜记录）
// @Tags report-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param bot body model.CreateMenuReportReq true "报菜信息"
// @Success 200 {object} http.Response
// @Router /report/robots [post]
func (c *ReportBotController) CreateBot(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	var req model.CreateMenuReportReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	report := &model.MenuReport{
		StoreID:  storeID,
		UserID:   userID,
		DishID:   req.DishID,
		Quantity: req.Quantity,
		Remark:   req.Remark,
	}

	if err := c.module.Create(report); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, report)
}

// ListBots godoc
// @Summary 获取报菜记录列表
// @Tags report-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param start_date query string false "开始日期 (YYYY-MM-DD)"
// @Param end_date query string false "结束日期 (YYYY-MM-DD)"
// @Success 200 {object} http.Response{data=[]model.MenuReport}
// @Router /report/robots [get]
func (c *ReportBotController) ListBots(ctx *gin.Context) {
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

		reports, err := c.module.ListByDateRange(storeID, startDate, endDate)
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		http.Success(ctx, reports)
		return
	}

	// 否则返回分页列表
	page := http.GetPage(ctx)
	pageSize := http.GetPageSize(ctx)

	reports, total, err := c.module.List(storeID, page, pageSize)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, reports, total, page, pageSize)
}

// GetBot godoc
// @Summary 获取报菜记录详情
// @Tags report-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "报菜记录ID"
// @Success 200 {object} http.Response{data=model.MenuReport}
// @Router /report/robots/{id} [get]
func (c *ReportBotController) GetBot(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu report ID")
		return
	}

	report, err := c.module.GetByID(uint(id), storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, report)
}

// UpdateBot godoc
// @Summary 更新报菜记录
// @Tags report-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "报菜记录ID"
// @Param bot body model.UpdateMenuReportReq true "报菜信息"
// @Success 200 {object} http.Response
// @Router /report/robots/{id} [put]
func (c *ReportBotController) UpdateBot(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu report ID")
		return
	}

	var req model.UpdateMenuReportReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.module.Update(uint(id), storeID, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteBot godoc
// @Summary 删除报菜记录
// @Tags report-bots
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "报菜记录ID"
// @Success 200 {object} http.Response
// @Router /report/robots/{id} [delete]
func (c *ReportBotController) DeleteBot(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu report ID")
		return
	}

	if err := c.module.Delete(uint(id), storeID); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}
