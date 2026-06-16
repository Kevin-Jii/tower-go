package controller

import (
	"errors"
	"strconv"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type AuditLogController struct {
	service *service.AuditLogService
}

func NewAuditLogController(service *service.AuditLogService) *AuditLogController {
	return &AuditLogController{service: service}
}

func (c *AuditLogController) List(ctx *gin.Context) {
	req := model.AuditLogListReq{
		Page:      http.GetPage(ctx),
		PageSize:  http.GetPageSize(ctx),
		StartTime: ctx.Query("start_time"),
		EndTime:   ctx.Query("end_time"),
		Module:    ctx.Query("module"),
		Action:    ctx.Query("action"),
		Status:    ctx.Query("status"),
		Keyword:   ctx.Query("keyword"),
	}
	if v, ok := http.ParseUintQuery(ctx, "user_id"); ok {
		req.UserID = v
	}
	if v, ok := http.ParseUintQuery(ctx, "store_id"); ok {
		req.StoreID = v
	}

	rows, total, err := c.service.List(req, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *AuditLogController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil || id == 0 {
		http.Error(ctx, 400, "无效的日志ID")
		return
	}
	row, err := c.service.Get(uint(id), middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		if errors.Is(err, service.ErrForbiddenAuditLog) {
			http.Error(ctx, 403, err.Error())
			return
		}
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, row)
}
