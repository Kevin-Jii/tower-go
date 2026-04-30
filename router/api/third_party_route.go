package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterThirdPartyRouteRoutes(r *gin.RouterGroup, c *Controllers) {
	group := r.Group("/third-party-routes")
	group.Use(middleware.AuthMiddleware())
	{
		group.GET("", middleware.Permission("third:account:list"), c.ThirdPartyRoute.List)
		group.POST("", middleware.Permission("third:account:edit"), c.ThirdPartyRoute.Create)
		group.PUT("/:id", middleware.Permission("third:account:edit"), c.ThirdPartyRoute.Update)
		group.DELETE("/:id", middleware.Permission("third:account:delete"), c.ThirdPartyRoute.Delete)
		group.POST("/:id/import-by-date", middleware.Permission("third:account:edit"), c.ThirdPartyRoute.ImportByDate)
		group.POST("/:id/logistics-sheets", middleware.Permission("third:account:edit"), c.ThirdPartyRoute.SaveLogisticsSheet)
		group.GET("/:id/logistics-sheets", middleware.Permission("third:account:list"), c.ThirdPartyRoute.ListLogisticsSheets)
	}
}
