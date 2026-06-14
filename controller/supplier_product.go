package controller

import (
	"strconv"
	"strings"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type SupplierProductController struct {
	productService       *service.SupplierProductService
	storeSupplierService *service.StoreSupplierService
}

func NewSupplierProductController(
	productService *service.SupplierProductService,
	storeSupplierService *service.StoreSupplierService,
) *SupplierProductController {
	return &SupplierProductController{
		productService:       productService,
		storeSupplierService: storeSupplierService,
	}
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
// @Success 200 {object} http.Response{data=[]model.SupplierProduct}
// @Router /supplier-products [get]
func (c *SupplierProductController) ListProducts(ctx *gin.Context) {
	var req model.ListSupplierProductReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}
	products, err := c.productService.ListProducts(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, products)
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

// CreateProductUnitSpec godoc
// @Summary 创建商品单位配置
// @Tags 商品单位配置
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body model.CreateProductUnitSpecReq true "单位配置"
// @Success 200 {object} http.Response
// @Router /product-unit-specs [post]
func (c *SupplierProductController) CreateProductUnitSpec(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	var req model.CreateProductUnitSpecReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.productService.CreateUnitSpec(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// ListProductUnitSpecs godoc
// @Summary 商品单位配置列表
// @Tags 商品单位配置
// @Produce json
// @Security Bearer
// @Param product_id query int true "商品ID"
// @Success 200 {object} http.Response{data=[]model.ProductUnitSpec}
// @Router /product-unit-specs [get]
func (c *SupplierProductController) ListProductUnitSpecs(ctx *gin.Context) {
	productID, ok := http.ParseUintQuery(ctx, "product_id")
	if !ok {
		http.Error(ctx, 400, "product_id is required")
		return
	}
	// 门店账号：仅允许查询本店已绑定供应商下的商品单位，防止凭 product_id 枚举其它供应商商品
	storeID := middleware.GetStoreID(ctx)
	if storeID > 0 && !middleware.HQUnboundAdmin(ctx) && c.storeSupplierService != nil {
		invalid, err := c.storeSupplierService.ValidateStoreProducts(storeID, []uint{productID})
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		if len(invalid) > 0 {
			http.ErrorApp(ctx, apicode.PermissionDenied)
			return
		}
	}
	specs, err := c.productService.ListUnitSpecs(productID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, specs)
}

// BatchListProductUnitSpecs godoc
// @Summary 批量查询商品单位配置
// @Tags 商品单位配置
// @Produce json
// @Security Bearer
// @Param product_ids query string true "商品ID列表，逗号分隔，如 1,2,3"
// @Success 200 {object} http.Response{data=[]model.ProductUnitSpec}
// @Router /product-unit-specs/batch [get]
func (c *SupplierProductController) BatchListProductUnitSpecs(ctx *gin.Context) {
	productIDs := parseProductIDs(ctx)
	if len(productIDs) == 0 {
		http.Error(ctx, 400, "product_ids is required")
		return
	}

	storeID := middleware.GetStoreID(ctx)
	if storeID > 0 && !middleware.HQUnboundAdmin(ctx) && c.storeSupplierService != nil {
		invalid, err := c.storeSupplierService.ValidateStoreProducts(storeID, productIDs)
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		if len(invalid) > 0 {
			http.ErrorApp(ctx, apicode.PermissionDenied)
			return
		}
	}

	specs, err := c.productService.ListUnitSpecsByProductIDs(productIDs)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, specs)
}

func parseProductIDs(ctx *gin.Context) []uint {
	rawValues := ctx.QueryArray("product_ids")
	if single := strings.TrimSpace(ctx.Query("product_ids")); single != "" && len(rawValues) == 0 {
		rawValues = []string{single}
	}

	seen := make(map[uint]struct{})
	ids := make([]uint, 0)
	for _, raw := range rawValues {
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			v, err := strconv.ParseUint(part, 10, 32)
			if err != nil || v == 0 {
				continue
			}
			id := uint(v)
			if _, ok := seen[id]; ok {
				continue
			}
			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}
	return ids
}

// UpdateProductUnitSpec godoc
// @Summary 更新商品单位配置
// @Tags 商品单位配置
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "单位配置ID"
// @Param body body model.UpdateProductUnitSpecReq true "单位配置"
// @Success 200 {object} http.Response
// @Router /product-unit-specs/{id} [put]
func (c *SupplierProductController) UpdateProductUnitSpec(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	var req model.UpdateProductUnitSpecReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.productService.UpdateUnitSpec(id, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// DeleteProductUnitSpec godoc
// @Summary 删除商品单位配置
// @Tags 商品单位配置
// @Produce json
// @Security Bearer
// @Param id path int true "单位配置ID"
// @Success 200 {object} http.Response
// @Router /product-unit-specs/{id} [delete]
func (c *SupplierProductController) DeleteProductUnitSpec(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	id, ok := http.ParseUintParam(ctx, "id")
	if !ok {
		return
	}
	if err := c.productService.DeleteUnitSpec(id); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// BatchUpsertProductUnitSpecs godoc
// @Summary 批量保存商品单位配置
// @Tags 商品单位配置
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body model.BatchUpsertProductUnitSpecsReq true "批量单位配置"
// @Success 200 {object} http.Response
// @Router /product-unit-specs/batch [post]
func (c *SupplierProductController) BatchUpsertProductUnitSpecs(ctx *gin.Context) {
	if !http.RequireAdmin(ctx) {
		return
	}
	var req model.BatchUpsertProductUnitSpecsReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.productService.BatchUpsertUnitSpecs(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}
