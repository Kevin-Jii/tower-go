package controller

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FileController 文件控制器
type FileController struct {
	rustfsService *service.RustFSService
}

// NewFileController 创建文件控制器
func NewFileController(rustfsService *service.RustFSService) *FileController {
	return &FileController{
		rustfsService: rustfsService,
	}
}

// Upload godoc
// @Summary 上传文件
// @Description 上传文件到RustFS对象存储
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "文件"
// @Param folder formData string false "目标文件夹，如 documents"
// @Success 200 {object} http.Response{data=service.RustFSUploadResult}
// @Router /files/upload [post]
func (c *FileController) Upload(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		http.Error(ctx, 400, "请选择要上传的文件")
		return
	}
	defer file.Close()

	// 获取目标文件夹
	baseFolder := ctx.PostForm("folder")
	if baseFolder == "" {
		baseFolder = "uploads"
	}

	// 按日期创建子文件夹: uploads/2025/12/02/
	now := time.Now()
	dateFolder := now.Format("2006/01/02")
	folder := fmt.Sprintf("%s/%s", baseFolder, dateFolder)

	// 生成唯一文件名
	ext := filepath.Ext(header.Filename)
	uniqueName := fmt.Sprintf("%s_%s%s", now.Format("150405"), uuid.New().String()[:8], ext)

	// 获取Content-Type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// 上传文件
	result, err := c.rustfsService.Upload(folder, uniqueName, file, header.Size, contentType)
	if err != nil {
		http.Error(ctx, 500, "上传失败: "+err.Error())
		return
	}

	// 返回原始文件名
	result.Name = header.Filename
	http.Success(ctx, result)
}

// UploadImage godoc
// @Summary 上传图片
// @Description 上传图片文件，仅支持jpg/png/gif/webp格式
// @Tags 文件管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "图片文件"
// @Param type formData string false "图片类型：product/supplier/avatar/purchase"
// @Param store_id formData int false "门店ID（type=purchase时必填）"
// @Success 200 {object} http.Response{data=service.RustFSUploadResult}
// @Router /files/upload-image [post]
func (c *FileController) UploadImage(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		http.Error(ctx, 400, "请选择要上传的图片")
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowedExts[ext] {
		http.Error(ctx, 400, "仅支持jpg/png/gif/webp格式的图片")
		return
	}

	// 限制文件大小（10MB）
	if header.Size > 10*1024*1024 {
		http.Error(ctx, 400, "图片大小不能超过10MB")
		return
	}

	// 根据类型确定文件夹
	imageType := ctx.PostForm("type")
	storeID := ctx.PostForm("store_id")
	baseFolder := ""
	
	switch imageType {
	case "product":
		baseFolder = "products"
	case "supplier":
		baseFolder = "suppliers"
	case "avatar":
		baseFolder = "avatars"
	case "purchase":
		// 采购单按门店分文件夹: purchases/store_1/2025/12/02/
		if storeID == "" {
			// 如果没传store_id，从token中获取
			storeID = fmt.Sprintf("%d", middleware.GetStoreID(ctx))
		}
		baseFolder = fmt.Sprintf("purchases/store_%s", storeID)
	default:
		baseFolder = "images"
	}

	// 按日期创建子文件夹: products/2025/12/02/
	now := time.Now()
	dateFolder := now.Format("2006/01/02")
	folder := fmt.Sprintf("%s/%s", baseFolder, dateFolder)

	// 生成唯一文件名
	uniqueName := fmt.Sprintf("%s_%s%s", now.Format("150405"), uuid.New().String()[:8], ext)

	// 获取Content-Type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/" + strings.TrimPrefix(ext, ".")
	}

	result, err := c.rustfsService.Upload(folder, uniqueName, file, header.Size, contentType)
	if err != nil {
		http.Error(ctx, 500, "上传失败: "+err.Error())
		return
	}

	result.Name = header.Filename
	http.Success(ctx, result)
}

// List godoc
// @Summary 文件列表
// @Description 获取指定目录的文件列表
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param prefix query string false "路径前缀"
// @Success 200 {object} http.Response
// @Router /files/list [get]
func (c *FileController) List(ctx *gin.Context) {
	prefix := ctx.Query("prefix")

	objects, err := c.rustfsService.List(prefix, true)
	if err != nil {
		http.Error(ctx, 500, "获取文件列表失败: "+err.Error())
		return
	}

	// 转换为简单格式
	var files []map[string]interface{}
	for _, obj := range objects {
		files = append(files, map[string]interface{}{
			"name":          obj.Key,
			"size":          obj.Size,
			"last_modified": obj.LastModified,
			"url":           c.rustfsService.GetPublicURL(obj.Key),
		})
	}

	http.Success(ctx, files)
}

// Delete godoc
// @Summary 删除文件
// @Description 删除指定的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body DeleteFileReq true "删除请求"
// @Success 200 {object} http.Response
// @Router /files/delete [post]
func (c *FileController) Delete(ctx *gin.Context) {
	// 仅管理员可删除文件
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可删除文件")
		return
	}

	var req DeleteFileReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	for _, name := range req.Names {
		path := name
		if req.Dir != "" {
			path = req.Dir + "/" + name
		}
		if err := c.rustfsService.Delete(path); err != nil {
			http.Error(ctx, 500, "删除失败: "+err.Error())
			return
		}
	}

	http.Success(ctx, nil)
}

// GetPresignedURL godoc
// @Summary 获取预签名URL
// @Description 获取文件的临时访问URL
// @Tags 文件管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param path query string true "文件路径"
// @Param expires query int false "过期时间(分钟)" default(60)
// @Success 200 {object} http.Response
// @Router /files/presigned-url [get]
func (c *FileController) GetPresignedURL(ctx *gin.Context) {
	path := ctx.Query("path")
	if path == "" {
		http.Error(ctx, 400, "请指定文件路径")
		return
	}

	expires := 60 // 默认60分钟
	if e := ctx.Query("expires"); e != "" {
		fmt.Sscanf(e, "%d", &expires)
	}

	url, err := c.rustfsService.GetPresignedURL(path, time.Duration(expires)*time.Minute)
	if err != nil {
		http.Error(ctx, 500, "获取URL失败: "+err.Error())
		return
	}

	http.Success(ctx, gin.H{"url": url})
}

// DeleteFileReq 删除文件请求
type DeleteFileReq struct {
	Dir   string   `json:"dir"`                      // 目录路径
	Names []string `json:"names" binding:"required"` // 文件名列表
}
