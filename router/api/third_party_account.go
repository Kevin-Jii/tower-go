package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterThirdPartyAccountRoutes(r *gin.RouterGroup, c *Controllers) {
	group := r.Group("/third-party-accounts")
	group.Use(middleware.AuthMiddleware())
	{
		group.GET("", middleware.Permission("third:account:list"), c.ThirdPartyAccount.List)
		group.GET("/:id", middleware.Permission("third:account:list"), c.ThirdPartyAccount.Get)
		group.POST("", middleware.Permission("third:account:add"), c.ThirdPartyAccount.Create)
		group.PUT("/:id", middleware.Permission("third:account:edit"), c.ThirdPartyAccount.Update)
		group.DELETE("/:id", middleware.Permission("third:account:delete"), c.ThirdPartyAccount.Delete)
		group.GET("/:id/orders", middleware.Permission("third:account:list"), c.ThirdPartyAccount.ListSyncedOrders)
		group.POST("/:id/test-login", middleware.Permission("third:account:edit"), c.ThirdPartyAccount.TestLogin)
		group.POST("/:id/sync-latest-orders", middleware.Permission("third:account:edit"), c.ThirdPartyAccount.SyncLatestOrders)
	}
}
