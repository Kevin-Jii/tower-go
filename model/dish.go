package model

import "time"

// Dish 菜品表
type Dish struct {
	ID          uint          `json:"id" gorm:"primarykey"`
	StoreID     uint          `json:"store_id" gorm:"not null;index"`                      // 所属门店 ID
	Store       *Store        `json:"store,omitempty" gorm:"foreignKey:StoreID"`           // 门店关联
	Name        string        `json:"name" gorm:"not null;type:varchar(100)"`              // 菜品名称
	Price       float64       `json:"price" gorm:"not null;type:decimal(10,2)"`            // 价格
	CategoryID  *uint         `json:"category_id" gorm:"index"`                            // 分类ID
	CategoryRef *DishCategory `json:"category_ref,omitempty" gorm:"foreignKey:CategoryID"` // 分类关联
	Status      int           `json:"status" gorm:"not null;default:1"`                    // 状态：1=上架，2=下架
	Image       string        `json:"image" gorm:"type:varchar(255)"`                      // 图片URL
	Remark      string        `json:"remark" gorm:"type:text"`                             // 备注
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// CreateDishReq 创建菜品请求
type CreateDishReq struct {
	Name       string  `json:"name" binding:"required"`
	Price      float64 `json:"price" binding:"required,gt=0"`
	CategoryID *uint   `json:"category_id"`
	Image      string  `json:"image"`
	Remark     string  `json:"remark"`
}

// UpdateDishReq 更新菜品请求
type UpdateDishReq struct {
	Name       string   `json:"name,omitempty"`
	Price      *float64 `json:"price,omitempty" binding:"omitempty,gt=0"`
	CategoryID *uint    `json:"category_id,omitempty"`
	Status     *int     `json:"status,omitempty"`
	Image      string   `json:"image,omitempty"`
	Remark     string   `json:"remark,omitempty"`
}
