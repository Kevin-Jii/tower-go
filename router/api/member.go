package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterMemberRoutes 注册会员管理路由
func RegisterMemberRoutes(v1 *gin.RouterGroup, c *Controllers) {
	members := v1.Group("/members")
	members.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		members.POST("", middleware.Permission("store:member:add"), c.Member.CreateMember)
		members.GET("", middleware.Permission("store:member:list"), c.Member.ListMembers)
		members.GET("/phone", middleware.Permission("store:member:list"), c.Member.GetMemberByPhone)
		members.GET("/point-rules", middleware.Permission("store:member:list"), c.Member.ListPointRules)
		members.POST("/point-rules", middleware.Permission("store:member:edit"), c.Member.CreatePointRule)
		members.PUT("/point-rules/:id", middleware.Permission("store:member:edit"), c.Member.UpdatePointRule)
		members.DELETE("/point-rules/:id", middleware.Permission("store:member:edit"), c.Member.DeletePointRule)
		members.GET("/:id", middleware.Permission("store:member:list"), c.Member.GetMember)
		members.GET("/:id/consumptions", middleware.Permission("store:member:list"), c.Member.ListMemberConsumptions)
		members.GET("/:id/gift-records", middleware.Permission("store:member:list"), c.InventoryLoss.ListMemberGiftRecords)
		members.PUT("/:id", middleware.Permission("store:member:edit"), c.Member.UpdateMember)
		members.DELETE("/:id", middleware.Permission("store:member:delete"), c.Member.DeleteMember)
		members.POST("/:id/adjust-balance", middleware.Permission("store:member:balance"), c.Member.AdjustBalance)
	}

	walletLogs := v1.Group("/wallet-logs")
	walletLogs.Use(middleware.AuthMiddleware())
	{
		walletLogs.GET("", middleware.Permission("store:member:list"), c.Member.ListWalletLogs)
	}

	rechargeOrders := v1.Group("/recharge-orders")
	rechargeOrders.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		rechargeOrders.POST("", middleware.Permission("store:member:add"), c.Member.CreateRechargeOrder)
		rechargeOrders.GET("", middleware.Permission("store:member:list"), c.Member.ListRechargeOrders)
		rechargeOrders.GET("/:id", middleware.Permission("store:member:list"), c.Member.GetRechargeOrder)
		rechargeOrders.POST("/pay", middleware.Permission("store:member:balance"), c.Member.PayRechargeOrder)
		rechargeOrders.POST("/:orderNo/cancel", middleware.Permission("store:member:edit"), c.Member.CancelRechargeOrder)
	}

	memberWines := v1.Group("/member-wines")
	memberWines.Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		memberWines.GET("", middleware.Permission("store:member:list"), c.Member.ListWineStorages)
		memberWines.POST("/deposit", middleware.Permission("store:member:edit"), c.Member.DepositWine)
		memberWines.POST("/withdraw", middleware.Permission("store:member:edit"), c.Member.WithdrawWine)
		memberWines.GET("/transactions", middleware.Permission("store:member:list"), c.Member.ListWineTransactions)
	}
}
