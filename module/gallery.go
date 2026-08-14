package module

import (
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"gorm.io/gorm"
)

type GalleryModule struct {
	db *gorm.DB
}

func NewGalleryModule(db *gorm.DB) *GalleryModule {
	return &GalleryModule{db: db}
}

// Create 创建图库记录
func (m *GalleryModule) Create(gallery *model.Gallery) error {
	return m.db.Create(gallery).Error
}

// GetByID 根据ID获取图库
func (m *GalleryModule) GetByID(id uint) (*model.Gallery, error) {
	var gallery model.Gallery
	if err := m.db.First(&gallery, id).Error; err != nil {
		return nil, err
	}
	return &gallery, nil
}

// GetByPath 根据路径获取图库
func (m *GalleryModule) GetByPath(path string) (*model.Gallery, error) {
	var gallery model.Gallery
	if err := m.db.Where("path = ?", path).First(&gallery).Error; err != nil {
		return nil, err
	}
	return &gallery, nil
}

// List 获取图库列表
func (m *GalleryModule) List(req *model.GalleryListReq) ([]*model.Gallery, int64, error) {
	var galleries []*model.Gallery
	var total int64

	query := m.db.Model(&model.Gallery{})

	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("id DESC").Offset(offset).Limit(req.PageSize).Find(&galleries).Error; err != nil {
		return nil, 0, err
	}

	return galleries, total, nil
}

// Update 更新图库
func (m *GalleryModule) Update(id uint, updates map[string]interface{}) error {
	return m.db.Model(&model.Gallery{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除图库记录
func (m *GalleryModule) Delete(id uint) error {
	return m.db.Delete(&model.Gallery{}, id).Error
}

// DeleteByPath 根据路径删除图库记录
func (m *GalleryModule) DeleteByPath(path string) error {
	return m.db.Where("path = ?", path).Delete(&model.Gallery{}).Error
}

// BatchDelete 批量删除
func (m *GalleryModule) BatchDelete(ids []uint) error {
	return m.db.Delete(&model.Gallery{}, ids).Error
}

func (m *GalleryModule) CreateUploadSession(session *model.GalleryUploadSession) error {
	return m.db.Create(session).Error
}

func (m *GalleryModule) FindResumableUploadSession(userID, storeID uint, fingerprint string) (*model.GalleryUploadSession, error) {
	var session model.GalleryUploadSession
	err := m.db.Where(
		"user_id = ? AND store_id = ? AND fingerprint = ? AND status IN ?",
		userID,
		storeID,
		fingerprint,
		[]string{model.GalleryUploadStatusUploading, model.GalleryUploadStatusCompleting, model.GalleryUploadStatusCompleted},
	).Order("created_at DESC").First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (m *GalleryModule) GetUploadSession(sessionID string, userID, storeID uint) (*model.GalleryUploadSession, error) {
	var session model.GalleryUploadSession
	if err := m.db.Where("id = ? AND user_id = ? AND store_id = ?", sessionID, userID, storeID).First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func (m *GalleryModule) UpdateUploadSessionStatus(sessionID, status string) error {
	return m.db.Model(&model.GalleryUploadSession{}).
		Where("id = ?", sessionID).
		Update("status", status).Error
}

func (m *GalleryModule) MarkUploadSessionCompleting(sessionID string) (bool, error) {
	result := m.db.Model(&model.GalleryUploadSession{}).
		Where("id = ? AND status = ?", sessionID, model.GalleryUploadStatusUploading).
		Update("status", model.GalleryUploadStatusCompleting)
	return result.RowsAffected == 1, result.Error
}

func (m *GalleryModule) CompleteUploadSession(session *model.GalleryUploadSession, gallery *model.Gallery) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(gallery).Error; err != nil {
			return err
		}
		now := time.Now()
		result := tx.Model(&model.GalleryUploadSession{}).
			Where("id = ? AND status = ?", session.ID, model.GalleryUploadStatusCompleting).
			Updates(map[string]interface{}{
				"status":       model.GalleryUploadStatusCompleted,
				"gallery_id":   gallery.ID,
				"completed_at": &now,
			})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

func (m *GalleryModule) ListExpiredUploadSessions(now time.Time, limit int) ([]*model.GalleryUploadSession, error) {
	if limit <= 0 {
		limit = 100
	}
	sessions := make([]*model.GalleryUploadSession, 0)
	err := m.db.Where(
		"expires_at <= ? AND status IN ?",
		now,
		[]string{model.GalleryUploadStatusUploading, model.GalleryUploadStatusCompleting},
	).Order("expires_at ASC").Limit(limit).Find(&sessions).Error
	return sessions, err
}
