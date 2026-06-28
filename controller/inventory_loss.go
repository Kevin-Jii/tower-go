package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/excelxml"
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

func (c *InventoryLossController) ExportOrders(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListInventoryLossOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	startDate, endDate, date := exportDateQuery(ctx.Query("date"), req.StartDate, req.EndDate)
	if date == "" {
		http.Error(ctx, 400, "请选择导出日期")
		return
	}
	req.StartDate = startDate
	req.EndDate = endDate
	req.Page = 1
	req.PageSize = exportPageSize
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}

	list, _, err := c.inventoryLossService.ListOrders(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	rows := make([][]interface{}, 0)
	for _, order := range list {
		if len(order.Items) == 0 {
			rows = append(rows, []interface{}{order.CreatedAt, order.OrderNo, inventoryLossTypeLabel(order.Type), memberName(order.Member), "", "", "", "", "", order.TotalCost, order.OperatorName, order.Reason, yesNo(order.IsCanceled)})
			continue
		}
		for _, item := range order.Items {
			rows = append(rows, []interface{}{
				order.CreatedAt,
				order.OrderNo,
				inventoryLossTypeLabel(order.Type),
				memberName(order.Member),
				item.ProductName,
				item.Unit,
				formatAmount(item.Quantity),
				item.BaseUnit,
				formatAmount(item.BaseQuantity),
				formatAmount(item.CostPrice),
				formatAmount(item.CostAmount),
				order.OperatorName,
				order.Reason,
				item.Remark,
				yesNo(order.IsCanceled),
			})
		}
	}
	data := excelxml.Build([]excelxml.Sheet{{
		Name:    "报损自用赠送",
		Headers: []string{"时间", "单据编号", "类型", "会员", "商品", "规格", "数量", "库存单位", "扣减库存", "成本单价", "成本金额", "操作人", "原因", "备注", "是否撤销"},
		Rows:    rows,
	}})
	http.File(ctx, data, excelxml.Filename("inventory-loss-orders-"+date))
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

func (c *InventoryLossController) UpdateOrder(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateInventoryLossOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	order, err := c.inventoryLossService.UpdateOrder(id, middleware.GetStoreID(ctx), &req, middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, err.Error())
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
