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

	// 出入库记录
	records := r.Group("/inventory-records").Use(middleware.AuthMiddleware())
	{
		records.GET("", c.Inventory.ListRecords)
		records.POST("", c.Inventory.CreateRecord)
	}
}
