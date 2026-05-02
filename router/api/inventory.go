package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterInventoryRoutes 注册库存路由
func RegisterInventoryRoutes(r *gin.RouterGroup, c *Controllers) {
	// 库存管理
	inventories := r.Group("/inventories").Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		inventories.GET("", middleware.Permission("inventory:list"), c.Inventory.ListInventory)
		inventories.PUT("/:id", middleware.PermissionAny("inventory:in", "inventory:out"), c.Inventory.UpdateInventory)
	}

	// 出入库单
	orders := r.Group("/inventory-orders").Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		orders.POST("", middleware.PermissionAny("inventory:in", "inventory:out"), c.Inventory.CreateOrder)
		orders.GET("", middleware.Permission("inventory:record"), c.Inventory.ListOrders)
		orders.GET("/no/:order_no", middleware.Permission("inventory:record"), c.Inventory.GetOrderByNo)
		orders.GET("/:id", middleware.Permission("inventory:record"), c.Inventory.GetOrderByID)
	}
}
