package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/excelxml"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type InventoryController struct {
	inventoryService *service.InventoryService
}

func NewInventoryController(inventoryService *service.InventoryService) *InventoryController {
	return &InventoryController{inventoryService: inventoryService}
}

// ListInventory godoc
// @Summary 库存列表
// @Tags 库存管理
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Param product_id query int false "商品ID"
// @Param keyword query string false "商品名称关键词"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.InventoryWithProduct}
// @Router /inventories [get]
func (c *InventoryController) ListInventory(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListInventoryReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}

	middleware.AttachAuthContextToHTTPRequest(ctx)

	list, total, err := c.inventoryService.ListInventory(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

// CreateOrder godoc
// @Summary 创建出入库单
// @Tags 库存管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param order body model.CreateInventoryOrderReq true "出入库单信息"
// @Success 200 {object} http.Response{data=model.InventoryOrder}
// @Router /inventory-orders [post]
func (c *InventoryController) CreateOrder(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	var req model.CreateInventoryOrderReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	order, err := c.inventoryService.CreateOrder(storeID, userID, &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, order)
}

// ListOrders godoc
// @Summary 出入库单列表
// @Tags 库存管理
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Param type query int false "类型 1=入库 2=出库"
// @Param order_no query string false "单号"
// @Param date query string false "日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.InventoryOrder}
// @Router /inventory-orders [get]
func (c *InventoryController) ListOrders(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListInventoryOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}

	middleware.AttachAuthContextToHTTPRequest(ctx)

	list, total, err := c.inventoryService.ListOrders(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

func (c *InventoryController) ExportOrders(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.ListInventoryOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	date := ctx.Query("date")
	if date == "" {
		date = req.Date
	}
	if date == "" {
		http.Error(ctx, 400, "请选择导出日期")
		return
	}
	req.Date = date
	req.Page = 1
	req.PageSize = exportPageSize
	if !middleware.HQUnboundAdmin(ctx) {
		req.StoreID = storeID
	}
	middleware.AttachAuthContextToHTTPRequest(ctx)

	list, _, err := c.inventoryService.ListOrders(ctx.Request.Context(), &req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	rows := make([][]interface{}, 0)
	for _, order := range list {
		if len(order.Items) == 0 {
			rows = append(rows, []interface{}{order.CreatedAt, order.OrderNo, inventoryTypeLabel(order.Type), order.StoreName, "", "", "", order.TotalQuantity, order.ItemCount, order.OperatorName, order.Reason, order.Remark})
			continue
		}
		for _, item := range order.Items {
			rows = append(rows, []interface{}{
				order.CreatedAt,
				order.OrderNo,
				inventoryTypeLabel(order.Type),
				order.StoreName,
				item.ProductName,
				formatAmount(item.Quantity),
				item.Unit,
				order.TotalQuantity,
				order.ItemCount,
				order.OperatorName,
				order.Reason,
				item.Remark,
			})
		}
	}
	data := excelxml.Build([]excelxml.Sheet{{
		Name:    "库存流水",
		Headers: []string{"时间", "单据编号", "类型", "门店", "商品", "明细数量", "单位", "单据总数量", "商品种类", "操作人", "原因", "备注"},
		Rows:    rows,
	}})
	http.File(ctx, data, excelxml.Filename("inventory-orders-"+date))
}

// GetOrderByNo godoc
// @Summary 根据单号获取出入库单详情
// @Tags 库存管理
// @Produce json
// @Security Bearer
// @Param order_no path string true "单据编号"
// @Success 200 {object} http.Response{data=model.InventoryOrder}
// @Router /inventory-orders/no/{order_no} [get]
func (c *InventoryController) GetOrderByNo(ctx *gin.Context) {
	orderNo := ctx.Param("order_no")
	if orderNo == "" {
		http.Error(ctx, 400, "单据编号不能为空")
		return
	}

	order, err := c.inventoryService.GetOrderByNoScoped(orderNo, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, "未找到该单据")
		return
	}

	http.Success(ctx, order)
}

// GetOrderByID godoc
// @Summary 根据ID获取出入库单详情
// @Tags 库存管理
// @Produce json
// @Security Bearer
// @Param id path int true "出入库单ID"
// @Success 200 {object} http.Response{data=model.InventoryOrder}
// @Router /inventory-orders/{id} [get]
func (c *InventoryController) GetOrderByID(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	order, err := c.inventoryService.GetOrderByIDScoped(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx))
	if err != nil {
		http.Error(ctx, 500, "未找到该单据")
		return
	}

	http.Success(ctx, order)
}

// UpdateInventory godoc
// @Summary 修改库存数量
// @Description 直接修改库存数量（仅管理员）
// @Tags 库存管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "库存ID"
// @Param body body model.UpdateInventoryReq true "库存信息"
// @Success 200 {object} http.Response
// @Router /inventories/{id} [put]
func (c *InventoryController) UpdateInventory(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可修改库存")
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdateInventoryReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	// 检查库存是否存在
	if _, err := c.inventoryService.GetInventoryByIDScoped(id, middleware.GetStoreID(ctx), middleware.HQUnboundAdmin(ctx)); err != nil {
		http.Error(ctx, 404, "库存记录不存在")
		return
	}

	if err := c.inventoryService.UpdateInventory(id, req.Quantity); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}
