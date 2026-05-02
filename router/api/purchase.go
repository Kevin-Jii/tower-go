package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterPurchaseRoutes 注册采购单管理路由
func RegisterPurchaseRoutes(v1 *gin.RouterGroup, c *Controllers) {
	purchaseOrders := v1.Group("/purchase-orders")
	purchaseOrders.Use(middleware.StoreAuthMiddleware(), middleware.StoreBusinessGuard())
	{
		purchaseOrders.POST("", middleware.Permission("purchase:add"), c.PurchaseOrder.CreateOrder)
		purchaseOrders.GET("", middleware.Permission("purchase:list"), c.PurchaseOrder.ListOrders)
		purchaseOrders.GET("/:id", middleware.Permission("purchase:list"), c.PurchaseOrder.GetOrder)
		purchaseOrders.PUT("/:id", middleware.Permission("purchase:edit"), c.PurchaseOrder.UpdateOrder)
		purchaseOrders.DELETE("/:id", middleware.Permission("purchase:delete"), c.PurchaseOrder.DeleteOrder)
		purchaseOrders.GET("/:id/by-supplier", middleware.Permission("purchase:list"), c.PurchaseOrder.GetOrdersBySupplier)

		// 状态机相关操作
		purchaseOrders.GET("/:id/actions", middleware.Permission("purchase:list"), c.PurchaseOrder.GetAvailableActions)
		purchaseOrders.POST("/:id/confirm", middleware.Permission("purchase:edit"), c.PurchaseOrder.ConfirmOrder)
		purchaseOrders.POST("/:id/complete", middleware.Permission("purchase:edit"), c.PurchaseOrder.CompleteOrder)
		purchaseOrders.POST("/:id/cancel", middleware.Permission("purchase:edit"), c.PurchaseOrder.CancelOrder)
	}
}
