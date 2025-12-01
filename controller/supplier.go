package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

type SupplierController struct {
	supplierService *service.SupplierService
}

func NewSupplierController(supplierService *service.SupplierService) *SupplierController {
	return &SupplierController{supplierService: supplierService}
}

// CreateSupplier godoc
// @Summary 创建供应商
// @Description 创建新供应商，编码自动生成（门店ID+4位序号）
// @Tags 供应商管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param supplier body model.CreateSupplierReq true "供应商信息"
// @Success 200 {object} http.Response "创建成功"
// @Failure 400 {object} http.Response "请求参数错误"
// @Failure 500 {object} http.Response "服务器内部错误"
// @Router /suppliers [post]
func (c *SupplierController) CreateSupplier(ctx *gin.Context) {

	var req model.CreateSupplierReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	// 获取当前用户门店ID
	storeID := middleware.GetStoreID(ctx)

	if err := c.supplierService.CreateSupplier(storeID, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetSupplier godoc
// @Summary 获取供应商详情
// @Description 获取供应商详细信息
// @Tags suppliers
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "供应商ID"
// @Success 200 {object} http.Response{data=model.Supplier}
// @Router /suppliers/{id} [get]
func (c *SupplierController) GetSupplier(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	supplier, err := c.supplierService.GetSupplier(id)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, supplier)
}

// ListSuppliers godoc
// @Summary 供应商列表
// @Description 获取供应商列表（支持分页和搜索）
// @Tags suppliers
// @Accept json
// @Produce json
// @Security Bearer
// @Param keyword query string false "搜索关键词（名称或编码）"
// @Param status query int false "状态筛选（0=禁用，1=启用）"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} http.Response{data=[]model.Supplier}
// @Router /suppliers [get]
func (c *SupplierController) ListSuppliers(ctx *gin.Context) {
	var req model.ListSupplierReq
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

	suppliers, total, err := c.supplierService.ListSuppliers(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, suppliers, total, req.Page, req.PageSize)
}

// UpdateSupplier godoc
// @Summary 更新供应商信息
// @Description 更新供应商信息（仅管理员）
// @Tags suppliers
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "供应商ID"
// @Param supplier body model.UpdateSupplierReq true "供应商信息"
// @Success 200 {object} http.Response
// @Router /suppliers/{id} [put]
func (c *SupplierController) UpdateSupplier(ctx *gin.Context) {

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	var req model.UpdateSupplierReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.supplierService.UpdateSupplier(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteSupplier godoc
// @Summary 删除供应商
// @Description 删除供应商（仅管理员）
// @Tags suppliers
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "供应商ID"
// @Success 200 {object} http.Response
// @Router /suppliers/{id} [delete]
func (c *SupplierController) DeleteSupplier(ctx *gin.Context) {

	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}

	if err := c.supplierService.DeleteSupplier(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}
