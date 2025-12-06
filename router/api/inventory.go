package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterInventoryRoutes 注册库存路由
func RegisterInventoryRoutes(r *gin.RouterGroup, c *Controllers) {
	// 库存管理
	inventories := r.Group("/inventories").Use(middleware.AuthMiddleware())
	{
		inventories.GET("", c.Inventory.ListInventory)
	}

	// 出入库单
	orders := r.Group("/inventory-orders").Use(middleware.AuthMiddleware())
	{
		orders.POST("", c.Inventory.CreateOrder)
		orders.GET("", c.Inventory.ListOrders)
		orders.GET("/no/:order_no", c.Inventory.GetOrderByNo)
		orders.GET("/:id", c.Inventory.GetOrderByID)
	}
}
