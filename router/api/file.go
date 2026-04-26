package api

import (
	"github.com/Kevin-Jii/tower-go/controller"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterFileRoutes 注册文件管理路由
func RegisterFileRoutes(v1 *gin.RouterGroup, fileController *controller.FileController) {
	if fileController == nil {
		return // RustFS未配置时跳过
	}

	files := v1.Group("/files")
	files.Use(middleware.AuthMiddleware())
	{
		files.POST("/upload", middleware.Permission("system:gallery:upload"), fileController.Upload)
		files.POST("/upload-image", middleware.Permission("system:gallery:upload"), fileController.UploadImage)
		files.GET("/list", middleware.Permission("system:gallery:list"), fileController.List)
		files.POST("/delete", middleware.Permission("system:gallery:delete"), fileController.Delete)
		files.GET("/presigned-url", middleware.Permission("system:gallery:list"), fileController.GetPresignedURL)
	}
}
