package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type MeituanAIController struct {
	svc *service.MeituanAIService
}

func NewMeituanAIController(svc *service.MeituanAIService) *MeituanAIController {
	return &MeituanAIController{svc: svc}
}

func (c *MeituanAIController) ListAccounts(ctx *gin.Context) {
	rows, err := c.svc.ListAccounts(middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, rows)
}

func (c *MeituanAIController) CreateAccount(ctx *gin.Context) {
	var req model.CreateMeituanAIAccountReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	row, err := c.svc.CreateAccount(middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, row)
}

func (c *MeituanAIController) UpdateAccount(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateMeituanAIAccountReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.UpdateAccount(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *MeituanAIController) Dashboard(ctx *gin.Context) {
	var req model.ListMeituanAIReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	stats, err := c.svc.Dashboard(middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, stats)
}

func (c *MeituanAIController) ImportOrders(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.ImportMeituanAIOrdersReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	res, err := c.svc.ImportOrders(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, res)
}

func (c *MeituanAIController) ImportReviews(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.ImportMeituanAIReviewsReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	res, err := c.svc.ImportReviews(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, res)
}

func (c *MeituanAIController) ListOrders(ctx *gin.Context) {
	var req model.ListMeituanAIReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	rows, total, err := c.svc.ListOrders(middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *MeituanAIController) ListReviews(ctx *gin.Context) {
	var req model.ListMeituanAIReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	rows, total, err := c.svc.ListReviews(middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *MeituanAIController) GenerateSuggestions(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.ListMeituanAIReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	res, err := c.svc.GenerateSuggestions(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, res)
}

func (c *MeituanAIController) ListSuggestions(ctx *gin.Context) {
	var req model.ListMeituanAIReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	rows, total, err := c.svc.ListSuggestions(middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *MeituanAIController) UpdateSuggestionStatus(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateMeituanAISuggestionStatusReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.svc.UpdateSuggestionStatus(id, middleware.GetStoreID(ctx), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}
