package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	httpPkg "github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type ThirdPartyRouteController struct {
	svc *service.ThirdPartyRouteService
}

func NewThirdPartyRouteController(svc *service.ThirdPartyRouteService) *ThirdPartyRouteController {
	return &ThirdPartyRouteController{svc: svc}
}

func (c *ThirdPartyRouteController) List(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can view")
		return
	}
	rows, err := c.svc.List()
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, rows)
}

func (c *ThirdPartyRouteController) Create(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can create")
		return
	}
	var req model.UpsertThirdPartyRouteReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	row, err := c.svc.Create(&req)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, row)
}

func (c *ThirdPartyRouteController) Update(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can update")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpsertThirdPartyRouteReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.Update(id, &req); err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "updated"})
}

func (c *ThirdPartyRouteController) Delete(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can delete")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.svc.Delete(id); err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, gin.H{"message": "deleted"})
}

func (c *ThirdPartyRouteController) ImportByDate(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can import")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.ImportRouteOrdersReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	items, err := c.svc.ImportByDateRange(id, req.StartDate, req.EndDate)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, gin.H{
		"start_date": req.StartDate,
		"end_date":   req.EndDate,
		"count":      len(items),
		"list":       items,
	})
}

func (c *ThirdPartyRouteController) SaveLogisticsSheet(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can save")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.SaveRouteSheetReq
	if !httpPkg.BindJSON(ctx, &req) {
		return
	}
	row, err := c.svc.SaveLogisticsSheet(id, &req)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, row)
}

func (c *ThirdPartyRouteController) ListLogisticsSheets(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		httpPkg.Error(ctx, 403, "only admin can view")
		return
	}
	id, ok := httpPkg.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	rows, err := c.svc.ListLogisticsSheets(id)
	if err != nil {
		httpPkg.Error(ctx, 500, err.Error())
		return
	}
	httpPkg.Success(ctx, rows)
}
