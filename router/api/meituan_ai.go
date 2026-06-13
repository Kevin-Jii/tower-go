package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterMeituanAIRoutes(r *gin.RouterGroup, c *Controllers) {
	group := r.Group("/meituan-ai").Use(middleware.AuthMiddleware(), middleware.StoreBusinessGuard())
	{
		group.GET("/accounts", middleware.Permission("store:account:list"), c.MeituanAI.ListAccounts)
		group.POST("/accounts", middleware.Permission("store:account:add"), c.MeituanAI.CreateAccount)
		group.PUT("/accounts/:id", middleware.Permission("store:account:edit"), c.MeituanAI.UpdateAccount)
		group.GET("/dashboard", middleware.Permission("store:account:list"), c.MeituanAI.Dashboard)
		group.POST("/accounts/:id/orders/import", middleware.Permission("store:account:add"), c.MeituanAI.ImportOrders)
		group.POST("/accounts/:id/reviews/import", middleware.Permission("store:account:add"), c.MeituanAI.ImportReviews)
		group.GET("/orders", middleware.Permission("store:account:list"), c.MeituanAI.ListOrders)
		group.GET("/reviews", middleware.Permission("store:account:list"), c.MeituanAI.ListReviews)
		group.POST("/accounts/:id/suggestions/generate", middleware.Permission("store:account:add"), c.MeituanAI.GenerateSuggestions)
		group.GET("/suggestions", middleware.Permission("store:account:list"), c.MeituanAI.ListSuggestions)
		group.PUT("/suggestions/:id/status", middleware.Permission("store:account:edit"), c.MeituanAI.UpdateSuggestionStatus)
	}
}
