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
		printers.POST("/bind", middleware.Permission("printer:bind"), c.Printer.BindPrinter)
		printers.DELETE("/:id", middleware.Permission("printer:unbind"), c.Printer.UnbindPrinter)

		// 查询
		printers.GET("", middleware.Permission("printer:list"), c.Printer.ListPrinters)
		printers.GET("/all", middleware.Permission("printer:list"), c.Printer.ListAllPrinters)
		printers.GET("/:id", middleware.Permission("printer:list"), c.Printer.GetPrinter)
		printers.GET("/default", middleware.Permission("printer:list"), c.Printer.GetStoreDefaultPrinter)

		// 更新
		printers.PUT("/:id", middleware.Permission("printer:edit"), c.Printer.UpdatePrinter)

		// 测试打印
		printers.POST("/:id/test", middleware.Permission("printer:query"), c.Printer.TestPrint)

		// 打印采购单
		printers.POST("/:id/print/purchase-order", middleware.Permission("printer:query"), c.Printer.PrintPurchaseOrder)

		// 状态查询
		printers.GET("/status", middleware.Permission("printer:query"), c.Printer.QueryPrinterStatus)
		printers.GET("/status/batch", middleware.Permission("printer:query"), c.Printer.BatchQueryStatus)
	}
}
