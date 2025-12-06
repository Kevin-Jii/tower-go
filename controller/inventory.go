package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
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
// @Tags inventory
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
	roleCode := middleware.GetRoleCode(ctx)

	var req model.ListInventoryReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 非管理员只能查看自己门店
	if roleCode != model.RoleCodeAdmin && roleCode != model.RoleCodeSuperAdmin {
		req.StoreID = storeID
	}

	list, total, err := c.inventoryService.ListInventory(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

// CreateOrder godoc
// @Summary 创建出入库单
// @Tags inventory
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
// @Tags inventory
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
	roleCode := middleware.GetRoleCode(ctx)

	var req model.ListInventoryOrderReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 非管理员只能查看自己门店
	if roleCode != model.RoleCodeAdmin && roleCode != model.RoleCodeSuperAdmin {
		req.StoreID = storeID
	}

	list, total, err := c.inventoryService.ListOrders(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}

// GetOrderByNo godoc
// @Summary 根据单号获取出入库单详情
// @Tags inventory
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

	order, err := c.inventoryService.GetOrderByNo(orderNo)
	if err != nil {
		http.Error(ctx, 500, "未找到该单据")
		return
	}

	http.Success(ctx, order)
}

// GetOrderByID godoc
// @Summary 根据ID获取出入库单详情
// @Tags inventory
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

	order, err := c.inventoryService.GetOrderByID(id)
	if err != nil {
		http.Error(ctx, 500, "未找到该单据")
		return
	}

	http.Success(ctx, order)
}
