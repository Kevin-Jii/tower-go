package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type PreOrderController struct {
	service *service.PreOrderService
}

func NewPreOrderController(service *service.PreOrderService) *PreOrderController {
	return &PreOrderController{service: service}
}

func (c *PreOrderController) Create(ctx *gin.Context) {
	var req model.CreatePreOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	storeID := middleware.GetStoreID(ctx)
	if middleware.HQUnboundAdmin(ctx) && req.StoreID > 0 {
		storeID = req.StoreID
	}
	order, err := c.service.Create(storeID, middleware.GetUserID(ctx), &req)
	if err != nil {
		http.ErrorFrom(ctx, err)
		return
	}
	http.Success(ctx, order)
}

func (c *PreOrderController) Update(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdatePreOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	order, err := c.service.Update(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.ErrorFrom(ctx, err)
		return
	}
	http.Success(ctx, order)
}

func (c *PreOrderController) Get(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	order, err := c.service.Get(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.ErrorFrom(ctx, err)
		return
	}
	http.Success(ctx, order)
}

func (c *PreOrderController) List(ctx *gin.Context) {
	var req model.ListPreOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}
	req.StoreID = middleware.ResolveQueryStoreID(ctx, "store_id")
	rows, total, err := c.service.List(&req)
	if err != nil {
		http.ErrorFrom(ctx, err)
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *PreOrderController) UpdateStatus(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdatePreOrderStatusReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	order, err := c.service.UpdateStatus(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), req.Status)
	if err != nil {
		http.ErrorFrom(ctx, err)
		return
	}
	http.Success(ctx, order)
}

func (c *PreOrderController) Delete(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.service.Delete(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx)); err != nil {
		http.ErrorFrom(ctx, err)
		return
	}
	http.Success(ctx, nil)
}
