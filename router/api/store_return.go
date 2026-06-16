package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterStoreReturnRoutes 注册门店返厂管理路由
func RegisterStoreReturnRoutes(v1 *gin.RouterGroup, c *Controllers) {
	returns := v1.Group("/store-returns")
	returns.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		returns.POST("/products", middleware.PermissionAny("store:return:product", "store:return:add", "store:return:edit"), c.StoreReturn.CreateProduct)
		returns.GET("/products", middleware.Permission("store:return:list"), c.StoreReturn.ListProducts)
		returns.PUT("/products/:id", middleware.PermissionAny("store:return:product", "store:return:add", "store:return:edit"), c.StoreReturn.UpdateProduct)
		returns.DELETE("/products/:id", middleware.PermissionAny("store:return:product", "store:return:delete", "store:return:edit"), c.StoreReturn.DeleteProduct)
		returns.POST("", middleware.Permission("store:return:add"), c.StoreReturn.Create)
		returns.GET("", middleware.Permission("store:return:list"), c.StoreReturn.List)
		returns.GET("/stats", middleware.Permission("store:return:list"), c.StoreReturn.Stats)
		returns.GET("/:id", middleware.Permission("store:return:list"), c.StoreReturn.Get)
		returns.PUT("/:id", middleware.Permission("store:return:edit"), c.StoreReturn.Update)
		returns.DELETE("/:id", middleware.Permission("store:return:delete"), c.StoreReturn.Delete)
	}
}
