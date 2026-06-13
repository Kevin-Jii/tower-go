package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type InventoryLossController struct {
	inventoryLossService *service.InventoryLossService
}

func NewInventoryLossController(inventoryLossService *service.InventoryLossService) *InventoryLossController {
	return &InventoryLossController{inventoryLossService: inventoryLossService}
}

func (c *InventoryLossController) CreateOrder(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	var req model.CreateInventoryLossOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	order, err := c.inventoryLossService.CreateOrder(storeID, userID, &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

func (c *InventoryLossController) ListOrders(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListInventoryLossOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}

	list, total, err := c.inventoryLossService.ListOrders(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

func (c *InventoryLossController) GetOrderByID(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	order, err := c.inventoryLossService.GetOrderByIDScoped(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 404, "未找到该单据")
		return
	}
	http.Success(ctx, order)
}

func (c *InventoryLossController) CancelOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.inventoryLossService.CancelOrder(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

func (c *InventoryLossController) ListMemberGiftRecords(ctx *gin.Context) {
	memberID, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.ListMemberGiftRecordsReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	list, total, err := c.inventoryLossService.ListMemberGiftRecords(memberID, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}
