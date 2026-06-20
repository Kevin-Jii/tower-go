package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type StoreExpenseController struct {
	storeExpenseService *service.StoreExpenseService
}

func NewStoreExpenseController(storeExpenseService *service.StoreExpenseService) *StoreExpenseController {
	return &StoreExpenseController{storeExpenseService: storeExpenseService}
}

func (c *StoreExpenseController) Create(ctx *gin.Context) {
	var req model.CreateStoreExpenseReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	record, err := c.storeExpenseService.Create(middleware.GetStoreID(ctx), middleware.GetUserID(ctx), &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, record)
}

func (c *StoreExpenseController) List(ctx *gin.Context) {
	var req model.ListStoreExpenseReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	list, total, err := c.storeExpenseService.List(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

func (c *StoreExpenseController) Get(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	record, err := c.storeExpenseService.Get(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 404, "未找到该支出记录")
		return
	}
	http.Success(ctx, record)
}

func (c *StoreExpenseController) Update(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateStoreExpenseReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	record, err := c.storeExpenseService.Update(id, middleware.GetStoreID(ctx), middleware.GetUserID(ctx), &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, record)
}

func (c *StoreExpenseController) Delete(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.storeExpenseService.Delete(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *StoreExpenseController) Stats(ctx *gin.Context) {
	var req model.ListStoreExpenseReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	stats, err := c.storeExpenseService.Stats(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, stats)
}
