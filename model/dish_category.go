package model

import "time"

// DishCategory 菜品分类表
type DishCategory struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	StoreID   uint      `json:"store_id" gorm:"not null;index"`                                              // 所属门店
	Name      string    `json:"name" gorm:"not null;type:varchar(50);uniqueIndex:idx_store_name,priority:2"` // 分类名称（同门店唯一）
	Code      string    `json:"code" gorm:"type:varchar(50);uniqueIndex"`                                    // 可选的编码
	Sort      int       `json:"sort" gorm:"default:0"`                                                       // 排序（升序）
	Status    int       `json:"status" gorm:"default:1"`                                                     // 1=启用 0=禁用
	Remark    string    `json:"remark" gorm:"type:varchar(255)"`                                             // 备注
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Dishes    []*Dish   `json:"dishes,omitempty" gorm:"foreignKey:CategoryID"` // 关联菜品
}

// CreateDishCategoryReq 创建分类请求
type CreateDishCategoryReq struct {
	Name   string `json:"name" binding:"required"`
	Code   string `json:"code"`
	Sort   int    `json:"sort"`
	Remark string `json:"remark"`
}

// UpdateDishCategoryReq 更新分类请求
type UpdateDishCategoryReq struct {
	Name   string `json:"name,omitempty"`
	Code   string `json:"code,omitempty"`
	Sort   *int   `json:"sort,omitempty"`
	Status *int   `json:"status,omitempty"`
	Remark string `json:"remark,omitempty"`
}

// ReorderDishCategoryItem 排序条目
type ReorderDishCategoryItem struct {
	ID   uint `json:"id" binding:"required"`
	Sort int  `json:"sort"`
}

// ReorderDishCategoriesReq 批量排序请求
type ReorderDishCategoriesReq struct {
	Items []ReorderDishCategoryItem `json:"items" binding:"required,dive"`
}
