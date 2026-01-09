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
		members.POST("", c.Member.CreateMember)
		members.GET("", c.Member.ListMembers)
		members.GET("/phone", c.Member.GetMemberByPhone)
		members.GET("/:id", c.Member.GetMember)
		members.PUT("/:id", c.Member.UpdateMember)
		members.DELETE("/:id", c.Member.DeleteMember)
		members.POST("/:id/adjust-balance", c.Member.AdjustBalance)
	}

	walletLogs := v1.Group("/wallet-logs")
	walletLogs.Use(middleware.AuthMiddleware())
	{
		walletLogs.GET("", c.Member.ListWalletLogs)
	}

	rechargeOrders := v1.Group("/recharge-orders")
	rechargeOrders.Use(middleware.AuthMiddleware())
	{
		rechargeOrders.POST("", c.Member.CreateRechargeOrder)
		rechargeOrders.GET("", c.Member.ListRechargeOrders)
		rechargeOrders.GET("/:id", c.Member.GetRechargeOrder)
		rechargeOrders.POST("/pay", c.Member.PayRechargeOrder)
		rechargeOrders.POST("/:orderNo/cancel", c.Member.CancelRechargeOrder)
	}
}
