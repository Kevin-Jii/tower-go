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
		galleries.POST("/upload", galleryController.Upload)
		galleries.GET("", galleryController.List)
		galleries.GET("/:id", galleryController.Get)
		galleries.PUT("/:id", galleryController.Update)
		galleries.DELETE("/:id", galleryController.Delete)
		galleries.POST("/batch-delete", galleryController.BatchDelete)
	}
}
