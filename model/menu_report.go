package model

import "time"

// MenuReport 报菜单表
type MenuReport struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	StoreID   uint      `json:"store_id" gorm:"not null;index"`            // 所属门店 ID
	Store     *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"` // 门店关联
	DishID    uint      `json:"dish_id" gorm:"not null;index"`             // 菜品 ID
	Dish      *Dish     `json:"dish,omitempty" gorm:"foreignKey:DishID"`   // 菜品关联
	UserID    uint      `json:"user_id" gorm:"not null;index"`             // 操作员 ID
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`   // 操作员关联
	Quantity  int       `json:"quantity" gorm:"not null"`                  // 报菜数量
	Remark    string    `json:"remark" gorm:"type:text"`                   // 备注
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMenuReportReq 创建报菜单请求
type CreateMenuReportReq struct {
	DishID   uint   `json:"dish_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Remark   string `json:"remark"`
}

// UpdateMenuReportReq 更新报菜单请求
type UpdateMenuReportReq struct {
	Quantity *int   `json:"quantity,omitempty" binding:"omitempty,gt=0"`
	Remark   string `json:"remark,omitempty"`
}

// MenuReportStats 报菜单统计
type MenuReportStats struct {
	DishID    uint    `json:"dish_id"`
	DishName  string  `json:"dish_name"`
	Category  string  `json:"category"`
	TotalQty  int     `json:"total_qty"`  // 总数量
	TotalDays int     `json:"total_days"` // 报菜天数
	AvgQty    float64 `json:"avg_qty"`    // 平均数量
}
