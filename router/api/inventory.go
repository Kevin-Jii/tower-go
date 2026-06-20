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

	// 库存报损/自用/赠送
	lossOrders := r.Group("/inventory-loss-orders").Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		lossOrders.POST("", middleware.PermissionAny("inventory:out", "inventory:in"), c.InventoryLoss.CreateOrder)
		lossOrders.GET("", middleware.Permission("inventory:record"), c.InventoryLoss.ListOrders)
		lossOrders.GET("/:id", middleware.Permission("inventory:record"), c.InventoryLoss.GetOrderByID)
		lossOrders.PUT("/:id", middleware.PermissionAny("inventory:out", "inventory:in"), c.InventoryLoss.UpdateOrder)
		lossOrders.DELETE("/:id", middleware.PermissionAny("inventory:out", "inventory:in"), c.InventoryLoss.CancelOrder)
	}
}
