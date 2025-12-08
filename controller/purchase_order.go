package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type PurchaseOrderController struct {
	orderService *service.PurchaseOrderService
}

func NewPurchaseOrderController(orderService *service.PurchaseOrderService) *PurchaseOrderController {
	return &PurchaseOrderController{orderService: orderService}
}

// CreateOrder godoc
// @Summary 创建采购单
// @Tags 采购单管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param order body model.CreatePurchaseOrderReq true "采购单信息"
// @Success 200 {object} http.Response{data=model.PurchaseOrder}
// @Router /purchase-orders [post]
func (c *PurchaseOrderController) CreateOrder(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	var req model.CreatePurchaseOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	order, err := c.orderService.CreateOrder(storeID, userID, &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

// GetOrder godoc
// @Summary 获取采购单详情
// @Tags 采购单管理
// @Produce json
// @Security Bearer
// @Param id path int true "采购单ID"
// @Success 200 {object} http.Response{data=model.PurchaseOrder}
// @Router /purchase-orders/{id} [get]
func (c *PurchaseOrderController) GetOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	order, err := c.orderService.GetOrder(id)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, order)
}

// ListOrders godoc
// @Summary 采购单列表
// @Tags 采购单管理
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Param supplier_id query int false "供应商ID"
// @Param status query int false "状态"
// @Param date query string false "日期 格式:2024-01-01"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.PurchaseOrder}
// @Router /purchase-orders [get]
func (c *PurchaseOrderController) ListOrders(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	var req model.ListPurchaseOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 非管理员只能查看自己门店的采购单
	if roleCode != model.RoleCodeAdmin {
		req.StoreID = storeID
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	orders, total, err := c.orderService.ListOrders(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, orders, total, req.Page, req.PageSize)
}

// UpdateOrder godoc
// @Summary 更新采购单
// @Tags 采购单管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "采购单ID"
// @Param order body model.UpdatePurchaseOrderReq true "采购单信息"
// @Success 200 {object} http.Response
// @Router /purchase-orders/{id} [put]
func (c *PurchaseOrderController) UpdateOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdatePurchaseOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.orderService.UpdateOrder(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// DeleteOrder godoc
// @Summary 删除采购单
// @Tags 采购单管理
// @Produce json
// @Security Bearer
// @Param id path int true "采购单ID"
// @Success 200 {object} http.Response
// @Router /purchase-orders/{id} [delete]
func (c *PurchaseOrderController) DeleteOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.orderService.DeleteOrder(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// GetOrdersBySupplier godoc
// @Summary 按供应商分组获取采购单明细
// @Tags 采购单管理
// @Produce json
// @Security Bearer
// @Param id path int true "采购单ID"
// @Success 200 {object} http.Response
// @Router /purchase-orders/{id}/by-supplier [get]
func (c *PurchaseOrderController) GetOrdersBySupplier(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	result, err := c.orderService.GetOrdersBySupplier(id)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, result)
}

// GetAvailableActions godoc
// @Summary 获取采购单可用操作
// @Tags 采购单管理
// @Produce json
// @Security Bearer
// @Param id path int true "采购单ID"
// @Success 200 {object} http.Response{data=[]string}
// @Router /purchase-orders/{id}/actions [get]
func (c *PurchaseOrderController) GetAvailableActions(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	actions, err := c.orderService.GetAvailableActions(id)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, actions)
}

// ConfirmOrder godoc
// @Summary 确认采购单
// @Tags 采购单管理
// @Produce json
// @Security Bearer
// @Param id path int true "采购单ID"
// @Success 200 {object} http.Response
// @Router /purchase-orders/{id}/confirm [post]
func (c *PurchaseOrderController) ConfirmOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.orderService.ConfirmOrder(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// CompleteOrder godoc
// @Summary 完成采购单
// @Tags 采购单管理
// @Produce json
// @Security Bearer
// @Param id path int true "采购单ID"
// @Success 200 {object} http.Response
// @Router /purchase-orders/{id}/complete [post]
func (c *PurchaseOrderController) CompleteOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.orderService.CompleteOrder(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// CancelOrder godoc
// @Summary 取消采购单
// @Tags 采购单管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "采购单ID"
// @Param body body object{reason=string} false "取消原因"
// @Success 200 {object} http.Response
// @Router /purchase-orders/{id}/cancel [post]
func (c *PurchaseOrderController) CancelOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	_ = ctx.ShouldBindJSON(&req)
	if err := c.orderService.CancelOrder(id, req.Reason); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}
