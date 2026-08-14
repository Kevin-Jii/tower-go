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
		galleries.POST("/multipart/init", middleware.Permission("system:gallery:upload"), galleryController.InitMultipartUpload)
		galleries.GET("/multipart/:session_id", middleware.Permission("system:gallery:upload"), galleryController.GetMultipartUploadStatus)
		galleries.PUT("/multipart/:session_id/parts/:part_number", middleware.Permission("system:gallery:upload"), galleryController.UploadMultipartPart)
		galleries.POST("/multipart/:session_id/complete", middleware.Permission("system:gallery:upload"), galleryController.CompleteMultipartUpload)
		galleries.DELETE("/multipart/:session_id", middleware.Permission("system:gallery:upload"), galleryController.AbortMultipartUpload)
		galleries.GET("", middleware.Permission("system:gallery:list"), galleryController.List)
		galleries.GET("/:id", middleware.Permission("system:gallery:list"), galleryController.Get)
		galleries.PUT("/:id", middleware.Permission("system:gallery:edit"), galleryController.Update)
		galleries.DELETE("/:id", middleware.Permission("system:gallery:delete"), galleryController.Delete)
		galleries.POST("/batch-delete", middleware.Permission("system:gallery:delete"), galleryController.BatchDelete)
	}
}
