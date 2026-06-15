package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type B2BController struct {
	service *service.B2BService
}

func NewB2BController(s *service.B2BService) *B2BController {
	return &B2BController{service: s}
}

func (c *B2BController) CreateCustomer(ctx *gin.Context) {
	var req model.CreateB2BCustomerReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	storeID := middleware.GetStoreID(ctx)
	if middleware.HQUnboundAdmin(ctx) && req.StoreID > 0 {
		storeID = req.StoreID
	}
	customer, err := c.service.CreateCustomer(storeID, &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, customer)
}

func (c *B2BController) UpdateCustomer(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateB2BCustomerReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	customer, err := c.service.UpdateCustomer(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, customer)
}

func (c *B2BController) ListCustomers(ctx *gin.Context) {
	var req model.ListB2BCustomerReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	rows, total, err := c.service.ListCustomers(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *B2BController) UpsertPrice(ctx *gin.Context) {
	var req model.UpsertB2BPriceReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.service.UpsertPrice(middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *B2BController) ListPrices(ctx *gin.Context) {
	var req model.ListB2BPriceReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	rows, total, err := c.service.ListPrices(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *B2BController) DeletePrice(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.service.DeletePrice(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *B2BController) CreateSupplyOrder(ctx *gin.Context) {
	var req model.CreateB2BSupplyOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	order, err := c.service.CreateSupplyOrder(middleware.GetStoreID(ctx), middleware.GetUserID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

func (c *B2BController) ListSupplyOrders(ctx *gin.Context) {
	var req model.ListB2BSupplyOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}
	rows, total, err := c.service.ListSupplyOrders(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, rows, total, req.Page, req.PageSize)
}

func (c *B2BController) GetSupplyOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	order, err := c.service.GetSupplyOrder(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}
