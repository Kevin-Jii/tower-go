package model

import "time"

// Gallery 图库
type Gallery struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(255);not null;comment:文件名"`
	Path        string    `json:"path" gorm:"type:varchar(500);not null;comment:存储路径"`
	URL         string    `json:"url" gorm:"type:varchar(500);not null;comment:访问URL"`
	Size        int64     `json:"size" gorm:"comment:文件大小(字节)"`
	MimeType    string    `json:"mime_type" gorm:"type:varchar(100);comment:MIME类型"`
	Category    string    `json:"category" gorm:"type:varchar(50);index;comment:分类(product/supplier/avatar/purchase/other)"`
	StoreID     uint      `json:"store_id" gorm:"index;comment:所属门店ID"`
	UploadBy    uint      `json:"upload_by" gorm:"comment:上传人ID"`
	UploadByName string   `json:"upload_by_name" gorm:"-"`
	Remark      string    `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Gallery) TableName() string {
	return "galleries"
}

// CreateGalleryReq 创建图库请求
type CreateGalleryReq struct {
	Name     string `json:"name" binding:"required,max=255"`
	Path     string `json:"path" binding:"required,max=500"`
	URL      string `json:"url" binding:"required,max=500"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type" binding:"max=100"`
	Category string `json:"category" binding:"max=50"`
	StoreID  uint   `json:"store_id"`
	Remark   string `json:"remark" binding:"max=500"`
}

// UpdateGalleryReq 更新图库请求
type UpdateGalleryReq struct {
	Name     string `json:"name" binding:"max=255"`
	Category string `json:"category" binding:"max=50"`
	Remark   string `json:"remark" binding:"max=500"`
}

// GalleryListReq 图库列表请求
type GalleryListReq struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
	Category string `form:"category"`
	Keyword  string `form:"keyword"`
	StoreID  uint   `form:"store_id"`
}
