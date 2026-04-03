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
		priceLists.POST("", c.PriceList.CreatePriceList)
		priceLists.PUT("/:id", c.PriceList.UpdatePriceList)
		priceLists.DELETE("/:id", c.PriceList.DeletePriceList)
		priceLists.GET("", c.PriceList.ListPriceLists)

		// 价目单分类管理
		priceLists.POST("/categories", c.PriceList.CreateCategory)
		priceLists.PUT("/categories/:id", c.PriceList.UpdateCategory)
		priceLists.DELETE("/categories/:id", c.PriceList.DeleteCategory)

		// 价目单商品管理
		priceLists.POST("/items", c.PriceList.AddItem)
		priceLists.POST("/items/batch", c.PriceList.BatchAddItems) // 批量添加商品
		priceLists.PUT("/items/:id", c.PriceList.UpdateItem)
		priceLists.DELETE("/items/:id", c.PriceList.DeleteItem)
	}

	// 公开接口（不需要认证）
	v1.GET("/price-lists/:id", c.PriceList.GetPriceList)
	v1.GET("/price-lists/:id/details", c.PriceList.GetPriceListWithDetails)
}
