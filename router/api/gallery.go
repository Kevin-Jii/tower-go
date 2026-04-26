package api

import (
	"github.com/Kevin-Jii/tower-go/controller"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterGalleryRoutes(v1 *gin.RouterGroup, galleryController *controller.GalleryController) {
	if galleryController == nil {
		return // RustFS未配置时跳过
	}

	galleries := v1.Group("/galleries")
	galleries.Use(middleware.AuthMiddleware())
	{
		galleries.POST("/upload", middleware.Permission("system:gallery:upload"), galleryController.Upload)
		galleries.GET("", middleware.Permission("system:gallery:list"), galleryController.List)
		galleries.GET("/:id", middleware.Permission("system:gallery:list"), galleryController.Get)
		galleries.PUT("/:id", middleware.Permission("system:gallery:edit"), galleryController.Update)
		galleries.DELETE("/:id", middleware.Permission("system:gallery:delete"), galleryController.Delete)
		galleries.POST("/batch-delete", middleware.Permission("system:gallery:delete"), galleryController.BatchDelete)
	}
}
