package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterSupplierRoutes 注册供应商相关路由
func RegisterSupplierRoutes(v1 *gin.RouterGroup, c *Controllers) {
	// 供应商管理
	suppliers := v1.Group("/suppliers")
	suppliers.Use(middleware.AuthMiddleware())
	{
		suppliers.POST("", c.Supplier.CreateSupplier)
		suppliers.GET("", c.Supplier.ListSuppliers)
		suppliers.GET("/:id", c.Supplier.GetSupplier)
		suppliers.PUT("/:id", c.Supplier.UpdateSupplier)
		suppliers.DELETE("/:id", c.Supplier.DeleteSupplier)
	}

	// 供应商商品管理
	supplierProducts := v1.Group("/supplier-products")
	supplierProducts.Use(middleware.AuthMiddleware())
	{
		supplierProducts.POST("", c.SupplierProduct.CreateProduct)
		supplierProducts.GET("", c.SupplierProduct.ListProducts)
		supplierProducts.GET("/:id", c.SupplierProduct.GetProduct)
		supplierProducts.PUT("/:id", c.SupplierProduct.UpdateProduct)
		supplierProducts.DELETE("/:id", c.SupplierProduct.DeleteProduct)
	}

	// 供应商分类管理
	supplierCategories := v1.Group("/supplier-categories")
	supplierCategories.Use(middleware.AuthMiddleware())
	{
		supplierCategories.POST("", c.SupplierProduct.CreateCategory)
		supplierCategories.GET("", c.SupplierProduct.ListCategories)
		supplierCategories.PUT("/:id", c.SupplierProduct.UpdateCategory)
		supplierCategories.DELETE("/:id", c.SupplierProduct.DeleteCategory)
	}

	// 门店供应商商品关联
	storeSuppliers := v1.Group("/store-suppliers")
	storeSuppliers.Use(middleware.StoreAuthMiddleware())
	{
		storeSuppliers.POST("/bind", c.StoreSupplier.BindProducts)
		storeSuppliers.DELETE("/unbind", c.StoreSupplier.UnbindProducts)
		storeSuppliers.PUT("/default", c.StoreSupplier.SetDefault)
		storeSuppliers.GET("", c.StoreSupplier.ListByStore)
	}
}
