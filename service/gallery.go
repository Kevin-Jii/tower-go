package service

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

type GalleryService struct {
	galleryModule *module.GalleryModule
	rustfsService *RustFSService
}

func NewGalleryService(galleryModule *module.GalleryModule, rustfsService *RustFSService) *GalleryService {
	return &GalleryService{
		galleryModule: galleryModule,
		rustfsService: rustfsService,
	}
}

// Create 创建图库记录
func (s *GalleryService) Create(req *model.CreateGalleryReq, uploadBy uint) (*model.Gallery, error) {
	gallery := &model.Gallery{
		Name:     req.Name,
		Path:     req.Path,
		URL:      req.URL,
		Size:     req.Size,
		MimeType: req.MimeType,
		Category: req.Category,
		StoreID:  req.StoreID,
		UploadBy: uploadBy,
		Remark:   req.Remark,
	}
	if err := s.galleryModule.Create(gallery); err != nil {
		return nil, err
	}
	return gallery, nil
}

// Get 获取图库详情
func (s *GalleryService) Get(id uint) (*model.Gallery, error) {
	return s.galleryModule.GetByID(id)
}

// List 获取图库列表
func (s *GalleryService) List(req *model.GalleryListReq) ([]*model.Gallery, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	return s.galleryModule.List(req)
}

// Update 更新图库
func (s *GalleryService) Update(id uint, req *model.UpdateGalleryReq) error {
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	if len(updates) == 0 {
		return nil
	}
	return s.galleryModule.Update(id, updates)
}

// Delete 删除图库（同时删除存储文件）
func (s *GalleryService) Delete(id uint) error {
	gallery, err := s.galleryModule.GetByID(id)
	if err != nil {
		return err
	}

	// 删除存储文件
	if s.rustfsService != nil && gallery.Path != "" {
		_ = s.rustfsService.Delete(gallery.Path)
	}

	return s.galleryModule.Delete(id)
}

// BatchDelete 批量删除
func (s *GalleryService) BatchDelete(ids []uint) error {
	for _, id := range ids {
		if err := s.Delete(id); err != nil {
			return err
		}
	}
	return nil
}
