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


// CreateRecord godoc
// @Summary 创建出入库记录
// @Tags inventory
// @Accept json
// @Produce json
// @Security Bearer
// @Param record body model.CreateInventoryRecordReq true "出入库信息"
// @Success 200 {object} http.Response
// @Router /inventory-records [post]
func (c *InventoryController) CreateRecord(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	userID := middleware.GetUserID(ctx)

	var req model.CreateInventoryRecordReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.inventoryService.CreateRecord(storeID, userID, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// ListRecords godoc
// @Summary 出入库记录列表
// @Tags inventory
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID"
// @Param product_id query int false "商品ID"
// @Param type query int false "类型 1=入库 2=出库"
// @Param date query string false "日期"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.InventoryRecord}
// @Router /inventory-records [get]
func (c *InventoryController) ListRecords(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	var req model.ListInventoryRecordReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 非管理员只能查看自己门店
	if roleCode != model.RoleCodeAdmin && roleCode != model.RoleCodeSuperAdmin {
		req.StoreID = storeID
	}

	list, total, err := c.inventoryService.ListRecords(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, list, total, req.Page, req.PageSize)
}
