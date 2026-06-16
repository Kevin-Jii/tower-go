package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterB2BRoutes(v1 *gin.RouterGroup, c *Controllers) {
	b2b := v1.Group("/b2b")
	b2b.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		b2b.POST("/customers", middleware.Permission("b2b:customer:add"), c.B2B.CreateCustomer)
		b2b.GET("/customers", middleware.Permission("b2b:customer:list"), c.B2B.ListCustomers)
		b2b.PUT("/customers/:id", middleware.Permission("b2b:customer:edit"), c.B2B.UpdateCustomer)

		b2b.POST("/prices", middleware.Permission("b2b:price:edit"), c.B2B.UpsertPrice)
		b2b.GET("/prices", middleware.Permission("b2b:price:list"), c.B2B.ListPrices)
		b2b.DELETE("/prices/:id", middleware.Permission("b2b:price:delete"), c.B2B.DeletePrice)

		b2b.POST("/supply-orders", middleware.Permission("b2b:order:add"), c.B2B.CreateSupplyOrder)
		b2b.GET("/supply-orders", middleware.Permission("b2b:order:list"), c.B2B.ListSupplyOrders)
		b2b.GET("/supply-orders/:id", middleware.Permission("b2b:order:list"), c.B2B.GetSupplyOrder)
		b2b.PUT("/supply-orders/:id/delivery-status", middleware.Permission("b2b:order:edit"), c.B2B.UpdateSupplyOrderDelivery)
		b2b.PUT("/supply-orders/:id/payment-status", middleware.Permission("b2b:order:edit"), c.B2B.UpdateSupplyOrderPayment)
	}
}
