package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterPriceListRoutes 注册价目单路由
func RegisterPriceListRoutes(v1 *gin.RouterGroup, c *Controllers) {
	// 价目单路由
	priceLists := v1.Group("/price-lists")
	{
		// 需要认证的管理接口
		priceLists.Use(middleware.AuthMiddleware())
		priceLists.POST("", middleware.Permission("price:add"), c.PriceList.CreatePriceList)
		priceLists.PUT("/:id", middleware.Permission("price:edit"), c.PriceList.UpdatePriceList)
		priceLists.DELETE("/:id", middleware.Permission("price:delete"), c.PriceList.DeletePriceList)
		priceLists.GET("", middleware.Permission("price:list"), c.PriceList.ListPriceLists)

		// 价目单分类管理
		priceLists.POST("/categories", middleware.Permission("price:edit"), c.PriceList.CreateCategory)
		priceLists.PUT("/categories/:id", middleware.Permission("price:edit"), c.PriceList.UpdateCategory)
		priceLists.DELETE("/categories/:id", middleware.Permission("price:delete"), c.PriceList.DeleteCategory)

		// 价目单商品管理
		priceLists.POST("/items", middleware.Permission("price:edit"), c.PriceList.AddItem)
		priceLists.POST("/items/batch", middleware.Permission("price:edit"), c.PriceList.BatchAddItems) // 批量添加商品
		priceLists.PUT("/items/:id", middleware.Permission("price:edit"), c.PriceList.UpdateItem)
		priceLists.DELETE("/items/:id", middleware.Permission("price:delete"), c.PriceList.DeleteItem)
	}

	// 公开接口（不需要认证）
	v1.GET("/price-lists/:id", c.PriceList.GetPriceList)
	v1.GET("/price-lists/:id/details", c.PriceList.GetPriceListWithDetails)
}
