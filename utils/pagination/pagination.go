package pagination

import (
	"math"

	"gorm.io/gorm"
)

// Pagination 分页参数
type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

// GetOffset 计算偏移量
func (p *Pagination) GetOffset() int {
	return (p.Page - 1) * p.PageSize
}

// GetLimit 获取每页数量
func (p *Pagination) GetLimit() int {
	return p.PageSize
}

// Validate 验证分页参数
func (p *Pagination) Validate() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.PageSize > 100 {
		p.PageSize = 100
	}
}

// Paginate 分页查询
// 返回分页结果和总数
func Paginate(db *gorm.DB, pagination *Pagination, result interface{}) (*gorm.DB, int64) {
	pagination.Validate()

	var total int64
	db.Model(result).Count(&total)

	// 执行分页查询
	db = db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())

	return db, total
}

// PageInfo 分页信息
type PageInfo struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// NewPageInfo 创建分页信息
func NewPageInfo(list interface{}, total int64, page, pageSize int) *PageInfo {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	return &PageInfo{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
