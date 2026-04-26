package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterMemberRoutes 注册会员管理路由
func RegisterMemberRoutes(v1 *gin.RouterGroup, c *Controllers) {
	members := v1.Group("/members")
	members.Use(middleware.AuthMiddleware())
	{
		members.POST("", middleware.Permission("store:member:add"), c.Member.CreateMember)
		members.GET("", middleware.Permission("store:member:list"), c.Member.ListMembers)
		members.GET("/phone", middleware.Permission("store:member:list"), c.Member.GetMemberByPhone)
		members.GET("/:id", middleware.Permission("store:member:list"), c.Member.GetMember)
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
	rechargeOrders.Use(middleware.AuthMiddleware())
	{
		rechargeOrders.POST("", middleware.Permission("store:member:add"), c.Member.CreateRechargeOrder)
		rechargeOrders.GET("", middleware.Permission("store:member:list"), c.Member.ListRechargeOrders)
		rechargeOrders.GET("/:id", middleware.Permission("store:member:list"), c.Member.GetRechargeOrder)
		rechargeOrders.POST("/pay", middleware.Permission("store:member:balance"), c.Member.PayRechargeOrder)
		rechargeOrders.POST("/:orderNo/cancel", middleware.Permission("store:member:edit"), c.Member.CancelRechargeOrder)
	}
}
