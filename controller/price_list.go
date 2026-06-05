package controller

import (
	"fmt"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type PriceListController struct {
	priceListService *service.PriceListService
}

func NewPriceListController(priceListService *service.PriceListService) *PriceListController {
	return &PriceListController{priceListService: priceListService}
}

// ===== 价目单相关 =====

// CreatePriceList godoc
// @Summary 创建价目单
// @Tags 价目单
// @Security Bearer
// @Param body body model.CreatePriceListReq true "价目单信息"
// @Success 200 {object} http.Response
// @Router /price-lists [post]
func (c *PriceListController) CreatePriceList(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	var req model.CreatePriceListReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	fmt.Printf("\n========== 创建价目单 ==========\n")
	fmt.Printf("门店ID: %d\n", req.StoreID)
	fmt.Printf("名称: %s\n", req.Name)
	fmt.Printf("LOGO: %s\n", req.Logo)
	fmt.Printf("==============================\n\n")

	if err := c.priceListService.CreatePriceList(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// UpdatePriceList godoc
// @Summary 更新价目单
// @Tags 价目单
// @Security Bearer
// @Param id path int true "价目单ID"
// @Param body body model.UpdatePriceListReq true "价目单信息"
// @Success 200 {object} http.Response
// @Router /price-lists/{id} [put]
func (c *PriceListController) UpdatePriceList(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdatePriceListReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.priceListService.UpdatePriceList(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeletePriceList godoc
// @Summary 删除价目单
// @Tags 价目单
// @Security Bearer
// @Param id path int true "价目单ID"
// @Success 200 {object} http.Response
// @Router /price-lists/{id} [delete]
func (c *PriceListController) DeletePriceList(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.priceListService.DeletePriceList(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetPriceList godoc
// @Summary 获取价目单详情
// @Tags 价目单
// @Param id path int true "价目单ID"
// @Success 200 {object} http.Response{data=model.PriceList}
// @Router /price-lists/{id} [get]
func (c *PriceListController) GetPriceList(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	priceList, err := c.priceListService.GetPriceList(id)
	if err != nil {
		http.Error(ctx, 404, "price list not found")
		return
	}

	http.Success(ctx, priceList)
}

// ListPriceLists godoc
// @Summary 获取价目单列表
// @Tags 价目单
// @Security Bearer
// @Param store_id query int true "门店ID"
// @Success 200 {object} http.Response{data=[]model.PriceList}
// @Router /price-lists [get]
func (c *PriceListController) ListPriceLists(ctx *gin.Context) {
	storeID := middleware.ResolveQueryStoreID(ctx, "store_id")
	if storeID == 0 {
		http.Error(ctx, 400, "store_id is required")
		return
	}

	priceLists, err := c.priceListService.ListPriceLists(storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, priceLists)
}

// GetPriceListWithDetails godoc
// @Summary 获取价目单完整结构（含分类和商品）
// @Tags 价目单
// @Param id path int true "价目单ID"
// @Success 200 {object} http.Response{data=model.PriceListResp}
// @Router /price-lists/{id}/details [get]
func (c *PriceListController) GetPriceListWithDetails(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	resp, err := c.priceListService.GetPriceListWithDetails(id)
	if err != nil {
		http.Error(ctx, 404, "price list not found")
		return
	}

	http.Success(ctx, resp)
}

// ===== 价目单分类相关 =====

// CreateCategory godoc
// @Summary 创建价目单分类
// @Tags 价目单
// @Security Bearer
// @Param body body model.CreatePriceListCategoryReq true "分类信息"
// @Success 200 {object} http.Response
// @Router /price-lists/categories [post]
func (c *PriceListController) CreateCategory(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	var req model.CreatePriceListCategoryReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.priceListService.CreateCategory(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// UpdateCategory godoc
// @Summary 更新价目单分类
// @Tags 价目单
// @Security Bearer
// @Param id path int true "分类ID"
// @Param body body model.UpdatePriceListCategoryReq true "分类信息"
// @Success 200 {object} http.Response
// @Router /price-lists/categories/{id} [put]
func (c *PriceListController) UpdateCategory(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdatePriceListCategoryReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.priceListService.UpdateCategory(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteCategory godoc
// @Summary 删除价目单分类
// @Tags 价目单
// @Security Bearer
// @Param id path int true "分类ID"
// @Success 200 {object} http.Response
// @Router /price-lists/categories/{id} [delete]
func (c *PriceListController) DeleteCategory(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.priceListService.DeleteCategory(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// ===== 价目单商品相关 =====

// AddItem godoc
// @Summary 添加价目单商品
// @Tags 价目单
// @Security Bearer
// @Param body body model.AddPriceListItemReq true "商品信息"
// @Success 200 {object} http.Response
// @Router /price-lists/items [post]
func (c *PriceListController) AddItem(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	var req model.AddPriceListItemReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.priceListService.AddItem(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// UpdateItem godoc
// @Summary 更新价目单商品
// @Tags 价目单
// @Security Bearer
// @Param id path int true "商品ID"
// @Param body body model.UpdatePriceListItemReq true "商品信息"
// @Success 200 {object} http.Response
// @Router /price-lists/items/{id} [put]
func (c *PriceListController) UpdateItem(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdatePriceListItemReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.priceListService.UpdateItem(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteItem godoc
// @Summary 删除价目单商品
// @Tags 价目单
// @Security Bearer
// @Param id path int true "商品ID"
// @Success 200 {object} http.Response
// @Router /price-lists/items/{id} [delete]
func (c *PriceListController) DeleteItem(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.priceListService.DeleteItem(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// BatchAddItems godoc
// @Summary 批量添加商品到价目单分类
// @Tags 价目单
// @Security Bearer
// @Param body body model.BatchAddPriceListItemsReq true "批量商品信息"
// @Success 200 {object} http.Response
// @Router /price-lists/items/batch [post]
func (c *PriceListController) BatchAddItems(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}

	var req model.BatchAddPriceListItemsReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.priceListService.BatchAddItems(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}
