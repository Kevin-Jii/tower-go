package model

import "time"

// MenuReportOrder 报菜记录单（主表）
type MenuReportOrder struct {
	ID        uint              `json:"id" gorm:"primaryKey"`
	StoreID   uint              `json:"store_id" gorm:"not null;index:idx_store_created"` // 所属门店 ID
	Store     *Store            `json:"store,omitempty" gorm:"foreignKey:StoreID"`        // 门店关联
	UserID    uint              `json:"user_id" gorm:"not null;index"`                    // 操作员 ID
	User      *User             `json:"user,omitempty" gorm:"foreignKey:UserID"`          // 操作员关联
	Remark    string            `json:"remark" gorm:"type:text"`                          // 备注
	Items     []*MenuReportItem `json:"items" gorm:"foreignKey:ReportOrderID"`            // 报菜详情列表
	CreatedAt time.Time         `json:"created_at" gorm:"index:idx_store_created"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// MenuReportItem 报菜详情（从表）
type MenuReportItem struct {
	ID            uint             `json:"id" gorm:"primaryKey;autoIncrement"`
	ReportOrderID uint             `json:"report_order_id" gorm:"not null;index:idx_report_order_dish"`  // 报菜记录单 ID（复合索引）
	ReportOrder   *MenuReportOrder `json:"report_order,omitempty" gorm:"foreignKey:ReportOrderID"`       // 报菜记录单关联
	DishID        uint             `json:"dish_id" gorm:"not null;index:idx_report_order_dish;index"`    // 菜品 ID（复合索引+单独索引）
	Dish          *Dish            `json:"dish,omitempty" gorm:"foreignKey:DishID"`                      // 菜品关联
	Quantity      int              `json:"quantity" gorm:"not null"`                                     // 报菜数量
	Remark        string           `json:"remark" gorm:"type:text"`                                      // 备注
	CreatedAt     time.Time        `json:"created_at" gorm:"index"`                                      // 添加时间索引用于统计查询
	UpdatedAt     time.Time        `json:"updated_at"`
}

// MenuReportItemReq 创建报菜详情项请求
type MenuReportItemReq struct {
	DishID   uint   `json:"dish_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Remark   string `json:"remark"`
}

// CreateMenuReportOrderReq 创建报菜记录单请求
type CreateMenuReportOrderReq struct {
	Remark string               `json:"remark"`
	Items  []*MenuReportItemReq `json:"items" binding:"required,min=1"`
}

// UpdateMenuReportOrderReq 更新报菜记录单请求
type UpdateMenuReportOrderReq struct {
	Remark *string `json:"remark,omitempty"`
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

// MenuReport 报菜记录（简化版，用于单个报菜）
type MenuReport struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	StoreID   uint      `json:"store_id" gorm:"not null;index"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	DishID    uint      `json:"dish_id" gorm:"not null;index"`
	Quantity  int       `json:"quantity" gorm:"not null"`
	Remark    string    `json:"remark" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateMenuReportReq 创建报菜记录请求
type CreateMenuReportReq struct {
	DishID   uint   `json:"dish_id" binding:"required"`
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Remark   string `json:"remark"`
}

// UpdateMenuReportReq 更新报菜记录请求
type UpdateMenuReportReq struct {
	DishID   *uint   `json:"dish_id,omitempty"`
	Quantity *int    `json:"quantity,omitempty"`
	Remark   *string `json:"remark,omitempty"`
}
