package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterSupplierRoutes 注册供应商相关路由
func RegisterSupplierRoutes(v1 *gin.RouterGroup, c *Controllers) {
	// 公开供应商档案（无需登录）
	publicSuppliers := v1.Group("/public/suppliers")
	{
		publicSuppliers.GET("/:id", c.Supplier.GetSupplierPublic)
	}

	// 供应商管理
	suppliers := v1.Group("/suppliers")
	suppliers.Use(middleware.AuthMiddleware())
	{
		suppliers.POST("", middleware.Permission("supplier:add"), c.Supplier.CreateSupplier)
		suppliers.GET("", middleware.Permission("supplier:list"), c.Supplier.ListSuppliers)
		suppliers.GET("/:id", middleware.Permission("supplier:list"), c.Supplier.GetSupplier)
		suppliers.PUT("/:id", middleware.Permission("supplier:edit"), c.Supplier.UpdateSupplier)
		suppliers.DELETE("/:id", middleware.Permission("supplier:delete"), c.Supplier.DeleteSupplier)
	}

	// 供应商商品管理
	supplierProducts := v1.Group("/supplier-products")
	supplierProducts.Use(middleware.AuthMiddleware())
	{
		supplierProducts.POST("", middleware.Permission("supplier:add"), c.SupplierProduct.CreateProduct)
		supplierProducts.GET("", middleware.Permission("supplier:list"), c.SupplierProduct.ListProducts)
		supplierProducts.GET("/:id", middleware.Permission("supplier:list"), c.SupplierProduct.GetProduct)
		supplierProducts.PUT("/:id", middleware.Permission("supplier:edit"), c.SupplierProduct.UpdateProduct)
		supplierProducts.DELETE("/:id", middleware.Permission("supplier:delete"), c.SupplierProduct.DeleteProduct)
	}

	// 供应商分类管理
	supplierCategories := v1.Group("/supplier-categories")
	supplierCategories.Use(middleware.AuthMiddleware())
	{
		supplierCategories.POST("", middleware.Permission("supplier:add"), c.SupplierProduct.CreateCategory)
		supplierCategories.GET("", middleware.Permission("supplier:list"), c.SupplierProduct.ListCategories)
		supplierCategories.PUT("/:id", middleware.Permission("supplier:edit"), c.SupplierProduct.UpdateCategory)
		supplierCategories.DELETE("/:id", middleware.Permission("supplier:delete"), c.SupplierProduct.DeleteCategory)
	}

	// 商品单位配置管理
	productUnitSpecs := v1.Group("/product-unit-specs")
	productUnitSpecs.Use(middleware.AuthMiddleware())
	{
		productUnitSpecs.POST("", middleware.Permission("supplier:edit"), c.SupplierProduct.CreateProductUnitSpec)
		productUnitSpecs.GET("", middleware.Permission("supplier:list"), c.SupplierProduct.ListProductUnitSpecs)
		productUnitSpecs.POST("/batch", middleware.Permission("supplier:edit"), c.SupplierProduct.BatchUpsertProductUnitSpecs)
		productUnitSpecs.PUT("/:id", middleware.Permission("supplier:edit"), c.SupplierProduct.UpdateProductUnitSpec)
		productUnitSpecs.DELETE("/:id", middleware.Permission("supplier:delete"), c.SupplierProduct.DeleteProductUnitSpec)
	}

	// 门店供应商关联
	storeSuppliers := v1.Group("/store-suppliers")
	storeSuppliers.Use(middleware.StoreAuthMiddleware())
	{
		storeSuppliers.POST("", middleware.Permission("store:menu"), c.StoreSupplier.BindSuppliers)               // 绑定供应商
		storeSuppliers.DELETE("", middleware.Permission("store:menu"), c.StoreSupplier.UnbindSuppliers)           // 解绑供应商
		storeSuppliers.GET("", middleware.Permission("supplier:list"), c.StoreSupplier.ListSuppliers)             // 获取绑定的供应商列表
		storeSuppliers.GET("/products", middleware.Permission("supplier:list"), c.StoreSupplier.ListProducts)     // 获取可采购的商品列表
		storeSuppliers.GET("/categories", middleware.Permission("supplier:list"), c.StoreSupplier.ListCategories) // 获取绑定供应商的分类列表
	}
}
