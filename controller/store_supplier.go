package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type StoreSupplierController struct {
	storeSupplierService *service.StoreSupplierService
}

func NewStoreSupplierController(storeSupplierService *service.StoreSupplierService) *StoreSupplierController {
	return &StoreSupplierController{storeSupplierService: storeSupplierService}
}

// BindProducts godoc
// @Summary 门店绑定供应商商品
// @Description 将供应商商品绑定到指定门店，支持批量绑定
// @Tags 门店供应商管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param binding body model.BindStoreSupplierReq true "绑定信息"
// @Success 200 {object} http.Response "绑定成功"
// @Failure 400 {object} http.Response "请求参数错误"
// @Failure 500 {object} http.Response "服务器内部错误"
// @Router /store-suppliers/bind [post]
func (c *StoreSupplierController) BindProducts(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	var req model.BindStoreSupplierReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.storeSupplierService.BindProducts(req.StoreID, req.ProductIDs); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// UnbindProducts godoc
// @Summary 门店解绑供应商商品
// @Description 将供应商商品从指定门店解绑，支持批量解绑
// @Tags 门店供应商管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param binding body model.BindStoreSupplierReq true "解绑信息"
// @Success 200 {object} http.Response "解绑成功"
// @Failure 400 {object} http.Response "请求参数错误"
// @Failure 500 {object} http.Response "服务器内部错误"
// @Router /store-suppliers/unbind [delete]
func (c *StoreSupplierController) UnbindProducts(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	var req model.BindStoreSupplierReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.storeSupplierService.UnbindProducts(req.StoreID, req.ProductIDs); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// SetDefault godoc
// @Summary 设置默认供应商商品
// @Description 为门店设置某个供应商商品为默认选项
// @Tags 门店供应商管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param setting body model.SetDefaultSupplierReq true "设置信息"
// @Success 200 {object} http.Response "设置成功"
// @Failure 400 {object} http.Response "请求参数错误"
// @Failure 500 {object} http.Response "服务器内部错误"
// @Router /store-suppliers/default [put]
func (c *StoreSupplierController) SetDefault(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	var req model.SetDefaultSupplierReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.storeSupplierService.SetDefault(req.StoreID, req.ProductID); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// ListByStore godoc
// @Summary 获取门店绑定的供应商商品列表
// @Description 获取当前门店已绑定的所有供应商商品，管理员可查看指定门店
// @Tags 门店供应商管理
// @Produce json
// @Security Bearer
// @Param store_id query int false "门店ID（管理员可指定，普通用户使用当前门店）"
// @Success 200 {object} http.Response{data=[]model.StoreSupplierProduct} "获取成功"
// @Failure 500 {object} http.Response "服务器内部错误"
// @Router /store-suppliers [get]
func (c *StoreSupplierController) ListByStore(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	// 管理员可以查看指定门店
	if roleCode == model.RoleCodeAdmin {
		if queryStoreID, ok := http.ParseUintQuery(ctx, "store_id"); ok && queryStoreID > 0 {
			storeID = queryStoreID
		}
	}

	bindings, err := c.storeSupplierService.ListByStoreID(storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, bindings)
}
