package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterPreOrderRoutes(v1 *gin.RouterGroup, c *Controllers) {
	orders := v1.Group("/pre-orders")
	orders.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		orders.POST("", middleware.Permission("preorder:add"), c.PreOrder.Create)
		orders.GET("", middleware.Permission("preorder:list"), c.PreOrder.List)
		orders.GET("/:id", middleware.Permission("preorder:list"), c.PreOrder.Get)
		orders.PUT("/:id", middleware.Permission("preorder:edit"), c.PreOrder.Update)
		orders.PUT("/:id/status", middleware.Permission("preorder:edit"), c.PreOrder.UpdateStatus)
		orders.DELETE("/:id", middleware.Permission("preorder:delete"), c.PreOrder.Delete)
	}
}
