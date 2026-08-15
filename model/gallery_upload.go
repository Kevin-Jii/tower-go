package model

import "time"

const (
	GalleryUploadStatusUploading  = "uploading"
	GalleryUploadStatusCompleting = "completing"
	GalleryUploadStatusCompleted  = "completed"
	GalleryUploadStatusAborted    = "aborted"
	GalleryUploadStatusExpired    = "expired"
	GalleryUploadStatusFailed     = "failed"
)

// GalleryUploadSession 保存可恢复分片上传的服务端状态。文件分片由 RustFS 管理，数据库不保存二进制内容。
type GalleryUploadSession struct {
	ID          string     `json:"id" gorm:"type:char(36);primaryKey"`
	UploadID    string     `json:"-" gorm:"type:varchar(512);not null"`
	ObjectPath  string     `json:"-" gorm:"type:varchar(500);not null"`
	Fingerprint string     `json:"fingerprint" gorm:"type:char(64);not null;index:idx_gallery_upload_resume,priority:3"`
	FileName    string     `json:"file_name" gorm:"type:varchar(255);not null"`
	FileSize    int64      `json:"file_size" gorm:"not null"`
	ContentType string     `json:"content_type" gorm:"type:varchar(100)"`
	Category    string     `json:"category" gorm:"type:varchar(50);not null"`
	Remark      string     `json:"remark" gorm:"type:varchar(500)"`
	ChunkSize   int64      `json:"chunk_size" gorm:"not null"`
	TotalParts  int        `json:"total_parts" gorm:"not null"`
	StoreID     uint       `json:"store_id" gorm:"not null;index:idx_gallery_upload_resume,priority:2"`
	UserID      uint       `json:"user_id" gorm:"not null;index:idx_gallery_upload_resume,priority:1"`
	Status      string     `json:"status" gorm:"type:varchar(20);not null;index"`
	GalleryID   uint       `json:"gallery_id" gorm:"index"`
	ExpiresAt   time.Time  `json:"expires_at" gorm:"not null;index"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (GalleryUploadSession) TableName() string {
	return "gallery_upload_sessions"
}

type InitGalleryMultipartUploadReq struct {
	FileName    string `json:"file_name" binding:"required,max=255"`
	FileSize    int64  `json:"file_size" binding:"required,gt=0"`
	ContentType string `json:"content_type" binding:"max=100"`
	Category    string `json:"category" binding:"omitempty,oneof=product supplier avatar purchase store-return other"`
	Remark      string `json:"remark" binding:"max=500"`
	Fingerprint string `json:"fingerprint" binding:"required,len=64,hexadecimal"`
}

type GalleryUploadedPart struct {
	PartNumber int   `json:"part_number"`
	Size       int64 `json:"size"`
}

type GalleryMultipartUploadStatus struct {
	SessionID     string                `json:"session_id"`
	ChunkSize     int64                 `json:"chunk_size"`
	TotalParts    int                   `json:"total_parts"`
	UploadedParts []GalleryUploadedPart `json:"uploaded_parts"`
	UploadedBytes int64                 `json:"uploaded_bytes"`
	FileSize      int64                 `json:"file_size"`
	Status        string                `json:"status"`
	ExpiresAt     time.Time             `json:"expires_at"`
	Gallery       *Gallery              `json:"gallery,omitempty"`
}

type GalleryUploadPartResult struct {
	PartNumber int    `json:"part_number"`
	Size       int64  `json:"size"`
	ETag       string `json:"etag"`
}
