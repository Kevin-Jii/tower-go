package controller

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type SupplierProductController struct {
	productService *service.SupplierProductService
}

func NewSupplierProductController(productService *service.SupplierProductService) *SupplierProductController {
	return &SupplierProductController{productService: productService}
}

// CreateProduct godoc
// @Summary 创建供应商商品
// @Tags 供应商商品
// @Accept json
// @Produce json
// @Security Bearer
// @Param product body model.CreateSupplierProductReq true "商品信息"
// @Success 200 {object} http.Response
// @Router /supplier-products [post]
func (c *SupplierProductController) CreateProduct(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	var req model.CreateSupplierProductReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.productService.CreateProduct(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// GetProduct godoc
// @Summary 获取供应商商品详情
// @Tags 供应商商品
// @Produce json
// @Security Bearer
// @Param id path int true "商品ID"
// @Success 200 {object} http.Response{data=model.SupplierProduct}
// @Router /supplier-products/{id} [get]
func (c *SupplierProductController) GetProduct(ctx *gin.Context) {
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	product, err := c.productService.GetProduct(id)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, product)
}

// ListProducts godoc
// @Summary 供应商商品列表
// @Tags 供应商商品
// @Produce json
// @Security Bearer
// @Param supplier_id query int false "供应商ID"
// @Param category_id query int false "分类ID"
// @Param keyword query string false "关键词"
// @Param status query int false "状态"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} http.Response{data=[]model.SupplierProduct}
// @Router /supplier-products [get]
func (c *SupplierProductController) ListProducts(ctx *gin.Context) {
	var req model.ListSupplierProductReq
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
	products, total, err := c.productService.ListProducts(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.SuccessWithPagination(ctx, products, total, req.Page, req.PageSize)
}

// UpdateProduct godoc
// @Summary 更新供应商商品
// @Tags 供应商商品
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "商品ID"
// @Param product body model.UpdateSupplierProductReq true "商品信息"
// @Success 200 {object} http.Response
// @Router /supplier-products/{id} [put]
func (c *SupplierProductController) UpdateProduct(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateSupplierProductReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.productService.UpdateProduct(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// DeleteProduct godoc
// @Summary 删除供应商商品
// @Tags 供应商商品
// @Produce json
// @Security Bearer
// @Param id path int true "商品ID"
// @Success 200 {object} http.Response
// @Router /supplier-products/{id} [delete]
func (c *SupplierProductController) DeleteProduct(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.productService.DeleteProduct(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// CreateCategory godoc
// @Summary 创建供应商分类
// @Tags 供应商分类
// @Accept json
// @Produce json
// @Security Bearer
// @Param category body model.CreateSupplierCategoryReq true "分类信息"
// @Success 200 {object} http.Response
// @Router /supplier-categories [post]
func (c *SupplierProductController) CreateCategory(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	var req model.CreateSupplierCategoryReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.productService.CreateCategory(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// ListCategories godoc
// @Summary 供应商分类列表
// @Tags 供应商分类
// @Produce json
// @Security Bearer
// @Param supplier_id query int true "供应商ID"
// @Success 200 {object} http.Response{data=[]model.SupplierCategory}
// @Router /supplier-categories [get]
func (c *SupplierProductController) ListCategories(ctx *gin.Context) {
	supplierID, ok := http.ParseUintQuery(ctx, "supplier_id")
	if !ok {
		http.Error(ctx, 400, "supplier_id is required")
		return
	}
	categories, err := c.productService.ListCategories(supplierID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, categories)
}

// UpdateCategory godoc
// @Summary 更新供应商分类
// @Tags 供应商分类
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "分类ID"
// @Param category body model.UpdateSupplierCategoryReq true "分类信息"
// @Success 200 {object} http.Response
// @Router /supplier-categories/{id} [put]
func (c *SupplierProductController) UpdateCategory(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateSupplierCategoryReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.productService.UpdateCategory(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// DeleteCategory godoc
// @Summary 删除供应商分类
// @Tags 供应商分类
// @Produce json
// @Security Bearer
// @Param id path int true "分类ID"
// @Success 200 {object} http.Response
// @Router /supplier-categories/{id} [delete]
func (c *SupplierProductController) DeleteCategory(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.productService.DeleteCategory(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}
