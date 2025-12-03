package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterPurchaseRoutes 注册采购单管理路由
func RegisterPurchaseRoutes(v1 *gin.RouterGroup, c *Controllers) {
	purchaseOrders := v1.Group("/purchase-orders")
	purchaseOrders.Use(middleware.StoreAuthMiddleware())
	{
		purchaseOrders.POST("", c.PurchaseOrder.CreateOrder)
		purchaseOrders.GET("", c.PurchaseOrder.ListOrders)
		purchaseOrders.GET("/:id", c.PurchaseOrder.GetOrder)
		purchaseOrders.PUT("/:id", c.PurchaseOrder.UpdateOrder)
		purchaseOrders.DELETE("/:id", c.PurchaseOrder.DeleteOrder)
		purchaseOrders.GET("/:id/by-supplier", c.PurchaseOrder.GetOrdersBySupplier)

		// 状态机相关操作
		purchaseOrders.GET("/:id/actions", c.PurchaseOrder.GetAvailableActions)
		purchaseOrders.POST("/:id/confirm", c.PurchaseOrder.ConfirmOrder)
		purchaseOrders.POST("/:id/complete", c.PurchaseOrder.CompleteOrder)
		purchaseOrders.POST("/:id/cancel", c.PurchaseOrder.CancelOrder)
	}
}
