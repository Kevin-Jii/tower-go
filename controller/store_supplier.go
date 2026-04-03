package controller

import (
	"encoding/json"
	"fmt"

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

// BindSuppliers 门店绑定供应商
func (c *StoreSupplierController) BindSuppliers(ctx *gin.Context) {
	var req model.BindStoreSuppliersReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.storeSupplierService.BindSuppliers(req.StoreID, req.SupplierIDs); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// UnbindSuppliers 门店解绑供应商
func (c *StoreSupplierController) UnbindSuppliers(ctx *gin.Context) {
	var req model.UnbindStoreSuppliersReq
	if !http.BindJSON(ctx, &req) {
		return
	}
	if err := c.storeSupplierService.UnbindSuppliers(req.StoreID, req.SupplierIDs); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// ListSuppliers 获取门店绑定的供应商列表
func (c *StoreSupplierController) ListSuppliers(ctx *gin.Context) {
	// 从 token 获取 storeID
	storeID := middleware.GetStoreID(ctx)

	// 管理员可以通过 query 参数查看其他门店
	if middleware.IsAdmin(ctx) {
		if queryStoreID, ok := http.ParseUintQuery(ctx, "store_id"); ok && queryStoreID > 0 {
			storeID = queryStoreID
		}
	}

	suppliers, err := c.storeSupplierService.ListSuppliersByStoreID(storeID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, suppliers)
}

// ListCategories godoc
// @Summary 获取门店绑定供应商的分类列表
// @Tags 门店供应商
// @Security Bearer
// @Param store_id query int false "门店ID（管理员可选）"
// @Param supplier_id query int false "供应商ID（可选，不传则返回所有绑定供应商的分类）"
// @Success 200 {object} http.Response{data=[]model.SupplierCategory}
// @Router /store-suppliers/categories [get]
func (c *StoreSupplierController) ListCategories(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	// 管理员可以通过 query 参数查看其他门店
	if middleware.IsAdmin(ctx) {
		if queryStoreID, ok := http.ParseUintQuery(ctx, "store_id"); ok && queryStoreID > 0 {
			storeID = queryStoreID
		}
	}

	supplierID, _ := http.ParseUintQuery(ctx, "supplier_id")

	categories, err := c.storeSupplierService.ListCategoriesByStoreID(storeID, supplierID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, categories)
}

// ListProducts 获取门店可采购的商品列表
func (c *StoreSupplierController) ListProducts(ctx *gin.Context) {
	// 从 token 获取 storeID
	storeID := middleware.GetStoreID(ctx)

	// 管理员可以通过 query 参数查看其他门店
	if middleware.IsAdmin(ctx) {
		if queryStoreID, ok := http.ParseUintQuery(ctx, "store_id"); ok && queryStoreID > 0 {
			storeID = queryStoreID
		}
	}

	fmt.Printf("🔍 ListProducts: storeID=%d, isAdmin=%v\n", storeID, middleware.IsAdmin(ctx))

	supplierID, _ := http.ParseUintQuery(ctx, "supplier_id")
	categoryID, _ := http.ParseUintQuery(ctx, "category_id")
	keyword := ctx.Query("keyword")
	products, err := c.storeSupplierService.ListProductsByStoreID(storeID, supplierID, categoryID, keyword)
	if err != nil {
		fmt.Printf("❌ ListProducts error: %v\n", err)
		http.Error(ctx, 500, err.Error())
		return
	}

	// 打印响应数据
	jsonData, _ := json.MarshalIndent(products, "", "  ")
	fmt.Printf("✅ ListProducts response (%d items):\n%s\n", len(products), string(jsonData))

	http.Success(ctx, products)
}
