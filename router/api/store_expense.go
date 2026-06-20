package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterStoreExpenseRoutes 注册门店支出路由
func RegisterStoreExpenseRoutes(v1 *gin.RouterGroup, c *Controllers) {
	expenses := v1.Group("/store-expenses")
	expenses.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		expenses.POST("", middleware.Permission("store:expenses:add"), c.StoreExpense.Create)
		expenses.GET("", middleware.Permission("store:expenses:list"), c.StoreExpense.List)
		expenses.GET("/stats", middleware.Permission("store:expenses:list"), c.StoreExpense.Stats)
		expenses.GET("/:id", middleware.Permission("store:expenses:list"), c.StoreExpense.Get)
		expenses.PUT("/:id", middleware.Permission("store:expenses:edit"), c.StoreExpense.Update)
		expenses.DELETE("/:id", middleware.Permission("store:expenses:delete"), c.StoreExpense.Delete)
	}
}
