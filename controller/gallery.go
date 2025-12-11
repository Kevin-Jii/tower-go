package controller

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GalleryController struct {
	galleryService *service.GalleryService
	rustfsService  *service.RustFSService
}

func NewGalleryController(galleryService *service.GalleryService, rustfsService *service.RustFSService) *GalleryController {
	return &GalleryController{
		galleryService: galleryService,
		rustfsService:  rustfsService,
	}
}

// Upload godoc
// @Summary 上传图片到图库
// @Tags 图库管理
// @Accept multipart/form-data
// @Produce json
// @Security Bearer
// @Param file formData file true "图片文件"
// @Param category formData string false "分类(product/supplier/avatar/other)"
// @Param remark formData string false "备注"
// @Success 200 {object} http.Response{data=model.Gallery}
// @Router /galleries/upload [post]
func (c *GalleryController) Upload(ctx *gin.Context) {
	if c.rustfsService == nil {
		http.Error(ctx, 500, "文件服务未启用")
		return
	}

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

	// 限制文件大小（20MB）
	if header.Size > 20*1024*1024 {
		http.Error(ctx, 400, "图片大小不能超过10MB")
		return
	}

	category := ctx.PostForm("category")
	if category == "" {
		category = "other"
	}

	// 上传到RustFS
	result, err := c.uploadToRustFS(file, header, category)
	if err != nil {
		http.Error(ctx, 500, "上传失败: "+err.Error())
		return
	}

	// 保存到数据库
	userID := middleware.GetUserID(ctx)
	storeID := middleware.GetStoreID(ctx)

	gallery, err := c.galleryService.Create(&model.CreateGalleryReq{
		Name:     header.Filename,
		Path:     result.Path,
		URL:      result.URL,
		Size:     result.Size,
		MimeType: header.Header.Get("Content-Type"),
		Category: category,
		StoreID:  storeID,
		Remark:   ctx.PostForm("remark"),
	}, userID)

	if err != nil {
		http.Error(ctx, 500, "保存图库记录失败: "+err.Error())
		return
	}

	http.Success(ctx, gallery)
}

func (c *GalleryController) uploadToRustFS(file multipart.File, header *multipart.FileHeader, category string) (*service.RustFSUploadResult, error) {
	now := time.Now()
	dateFolder := now.Format("2006/01/02")
	folder := fmt.Sprintf("gallery/%s/%s", category, dateFolder)

	ext := filepath.Ext(header.Filename)
	uniqueName := fmt.Sprintf("%s_%s%s", now.Format("150405"), uuid.New().String()[:8], ext)

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/" + strings.TrimPrefix(ext, ".")
	}

	return c.rustfsService.Upload(folder, uniqueName, file, header.Size, contentType)
}

// List godoc
// @Summary 获取图库列表
// @Tags 图库管理
// @Produce json
// @Security Bearer
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param category query string false "分类"
// @Param keyword query string false "关键词"
// @Success 200 {object} http.Response{data=[]model.Gallery}
// @Router /galleries [get]
func (c *GalleryController) List(ctx *gin.Context) {
	var req model.GalleryListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 20
	}

	// 非管理员只能看自己门店的图片
	if !middleware.IsAdmin(ctx) {
		req.StoreID = middleware.GetStoreID(ctx)
	}

	galleries, total, err := c.galleryService.List(&req)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, gin.H{
		"list":  galleries,
		"total": total,
		"page":  req.Page,
		"size":  req.PageSize,
	})
}

// Get godoc
// @Summary 获取图库详情
// @Tags 图库管理
// @Produce json
// @Security Bearer
// @Param id path int true "图库ID"
// @Success 200 {object} http.Response{data=model.Gallery}
// @Router /galleries/{id} [get]
func (c *GalleryController) Get(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	gallery, err := c.galleryService.Get(uint(id))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, gallery)
}

// Update godoc
// @Summary 更新图库信息
// @Tags 图库管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "图库ID"
// @Param body body model.UpdateGalleryReq true "更新信息"
// @Success 200 {object} http.Response
// @Router /galleries/{id} [put]
func (c *GalleryController) Update(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	var req model.UpdateGalleryReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.galleryService.Update(uint(id), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// Delete godoc
// @Summary 删除图库
// @Tags 图库管理
// @Produce json
// @Security Bearer
// @Param id path int true "图库ID"
// @Success 200 {object} http.Response
// @Router /galleries/{id} [delete]
func (c *GalleryController) Delete(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可删除图片")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "无效的ID")
		return
	}

	if err := c.galleryService.Delete(uint(id)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// BatchDelete godoc
// @Summary 批量删除图库
// @Tags 图库管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body BatchDeleteReq true "删除请求"
// @Success 200 {object} http.Response
// @Router /galleries/batch-delete [post]
func (c *GalleryController) BatchDelete(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅管理员可删除图片")
		return
	}

	var req BatchDeleteReq
	if !http.BindJSON(ctx, &req) {
		return
	}

	if err := c.galleryService.BatchDelete(req.IDs); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	http.Success(ctx, nil)
}

// BatchDeleteReq 批量删除请求
type BatchDeleteReq struct {
	IDs []uint `json:"ids" binding:"required"`
}
