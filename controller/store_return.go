package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type StoreReturnController struct {
	storeReturnService *service.StoreReturnService
}

func NewStoreReturnController(storeReturnService *service.StoreReturnService) *StoreReturnController {
	return &StoreReturnController{storeReturnService: storeReturnService}
}

func (c *StoreReturnController) Create(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	var req model.CreateStoreReturnReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	record, err := c.storeReturnService.Create(storeID, userID, &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, record)
}

func (c *StoreReturnController) List(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListStoreReturnReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}

	list, total, err := c.storeReturnService.List(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

func (c *StoreReturnController) Get(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	record, err := c.storeReturnService.Get(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 404, "未找到该返厂记录")
		return
	}
	http.Success(ctx, record)
}

func (c *StoreReturnController) Update(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdateStoreReturnReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	record, err := c.storeReturnService.Update(id, middleware.GetStoreID(ctx), middleware.GetUserID(ctx), &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, record)
}

func (c *StoreReturnController) Delete(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.storeReturnService.Delete(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *StoreReturnController) Stats(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	var req model.ListStoreReturnReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}
	stats, err := c.storeReturnService.Stats(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, stats)
}

func (c *StoreReturnController) CreateProduct(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.CreateStoreReturnProductReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	product, err := c.storeReturnService.CreateProduct(storeID, &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, product)
}

func (c *StoreReturnController) ListProducts(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListStoreReturnProductReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}
	list, total, err := c.storeReturnService.ListProducts(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

func (c *StoreReturnController) UpdateProduct(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdateStoreReturnProductReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	product, err := c.storeReturnService.UpdateProduct(id, middleware.GetStoreID(ctx), &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, product)
}

func (c *StoreReturnController) DeleteProduct(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.storeReturnService.DeleteProduct(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}
