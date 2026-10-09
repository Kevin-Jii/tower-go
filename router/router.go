package router

import (
	"fmt"
	"net/http"

	"github.com/Kevin-Jii/tower-go/apidocs"
	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/controller"
	"github.com/Kevin-Jii/tower-go/router/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Setup 初始化路由
func Setup(r *gin.Engine, c *api.Controllers) {
	// 初始化健康检查控制器
	healthController := controller.NewHealthController()

	// 注册健康检查路由（无需认证）
	r.GET("/health", healthController.Check)
	r.GET("/ready", healthController.Ready)
	r.GET("/live", healthController.Live)

	v1 := r.Group("/api/v1")

	// 注册各模块路由
	api.RegisterAuthRoutes(v1, c)
	api.RegisterUserRoutes(v1, c)
	api.RegisterRoleRoutes(v1)
	api.RegisterPermissionRoutes(v1, c)
	api.RegisterStoreRoutes(v1, c)
	api.RegisterMenuRoutes(v1, c)
	api.RegisterDingTalkRoutes(v1, c)
	api.RegisterSupplierRoutes(v1, c)
	api.RegisterPurchaseRoutes(v1, c)
	api.RegisterDictRoutes(v1, c)
	api.RegisterInventoryRoutes(v1, c)
	api.RegisterFileRoutes(v1, c.File)
	api.RegisterGalleryRoutes(v1, c.Gallery)
	api.RegisterStoreAccountRoutes(v1, c)
	api.RegisterStoreExpenseRoutes(v1, c)
	api.RegisterStoreReturnRoutes(v1, c)
	api.RegisterMeituanAIRoutes(v1, c)
	api.RegisterStatisticsRoutes(v1, c)
	api.RegisterAIAssistantRoutes(v1, c)
	api.RegisterMessageTemplateRoutes(v1, c)
	api.RegisterSmsCampaignRoutes(v1, c)
	api.RegisterStoreSmsConfigRoutes(v1, c)
	api.RegisterAliyunSmsTemplateRoutes(v1, c)
	api.RegisterMemberTagRoutes(v1, c)
	api.RegisterMemberRoutes(v1, c)
	api.RegisterPrinterRoutes(v1, c)
	api.RegisterPriceListRoutes(v1, c)
	api.RegisterB2BRoutes(v1, c)
	api.RegisterPreOrderRoutes(v1, c)
	api.RegisterThirdPartyAccountRoutes(v1, c)
	api.RegisterThirdPartyRouteRoutes(v1, c)
	api.RegisterAuditLogRoutes(v1, c)
	api.RegisterInternalRoutes(r, c)

	// WebSocket
	r.GET("/ws", controller.WebSocketHandler)

	registerDocumentationRoutes(r)

	addr := fmt.Sprintf(":%d", config.GetConfig().App.Port)
	fmt.Printf("📚 Swagger UI: http://localhost%s/swagger/index.html\n", addr)
	fmt.Printf("📚 Scalar Docs: http://localhost%s/docs\n\n", addr)
}

func registerDocumentationRoutes(r *gin.Engine) {
	// Serve the frontend brand asset for the documentation workbench.
	r.GET("/docs-assets/tower-logo.svg", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/svg+xml; charset=utf-8", []byte(towerLogoSVG))
	})

	// Serve embedded Swagger UI assets under a path separate from gin-swagger's wildcard route.
	for _, asset := range []string{"swagger-ui.css", "swagger-ui-bundle.js", "swagger-ui-standalone-preset.js"} {
		asset := asset
		r.GET("/docs-assets/"+asset, func(c *gin.Context) {
			file, err := swaggerFiles.HTTP.Open("/" + asset)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
			defer file.Close()
			info, err := file.Stat()
			if err != nil {
				c.Status(http.StatusInternalServerError)
				return
			}
			http.ServeContent(c.Writer, c.Request, asset, info.ModTime(), file)
		})
	}

	// Swagger UI 仍由 gin-swagger 提供，但 doc.json 返回经过转换和校验的 OpenAPI 3 文档。
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger/*any", func(c *gin.Context) {
		if c.Param("any") != "/doc.json" {
			swaggerHandler(c)
			return
		}

		document, err := apidocs.Document()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.Data(http.StatusOK, "application/json; charset=utf-8", document)
	})

	// Custom Swagger UI workbench with the existing OpenAPI 3 endpoint.
	r.GET("/docs", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.Header("Cache-Control", "no-cache")
		c.String(http.StatusOK, documentationPage)
	})
}
