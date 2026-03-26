package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterPrinterRoutes 注册打印机管理路由
func RegisterPrinterRoutes(v1 *gin.RouterGroup, c *Controllers) {
	printers := v1.Group("/printers")
	printers.Use(middleware.AuthMiddleware())
	{
		// 绑定/解绑
		printers.POST("/bind", c.Printer.BindPrinter)
		printers.DELETE("/:id", c.Printer.UnbindPrinter)

		// 查询
		printers.GET("", c.Printer.ListPrinters)
		printers.GET("/all", c.Printer.ListAllPrinters)
		printers.GET("/:id", c.Printer.GetPrinter)
		printers.GET("/default", c.Printer.GetStoreDefaultPrinter)

		// 更新
		printers.PUT("/:id", c.Printer.UpdatePrinter)

		// 状态查询
		printers.GET("/status", c.Printer.QueryPrinterStatus)
		printers.GET("/status/batch", c.Printer.BatchQueryStatus)
	}
}