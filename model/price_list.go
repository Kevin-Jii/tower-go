package model

import "time"

// PriceList 价目单
type PriceList struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	StoreID     uint      `json:"store_id" gorm:"not null;index;comment:门店ID"`
	Store       *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"`
	Name        string    `json:"name" gorm:"type:varchar(200);not null;comment:价目单名称"`
	Logo        string    `json:"logo" gorm:"type:varchar(500);comment:品牌LOGO URL"`
	Description string    `json:"description" gorm:"type:varchar(500);comment:描述"`
	Status      int8      `json:"status" gorm:"not null;default:1;comment:状态 1=启用 0=禁用"`
	IsDefault   int8      `json:"is_default" gorm:"not null;default:0;comment:是否默认 1=是 0=否"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (PriceList) TableName() string {
	return "price_lists"
}

// PriceListCategory 价目单分类（组合分类）
type PriceListCategory struct {
	ID          uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	PriceListID uint       `json:"price_list_id" gorm:"not null;index;comment:价目单ID"`
	PriceList   *PriceList `json:"price_list,omitempty" gorm:"foreignKey:PriceListID"`
	MainTitle   string     `json:"main_title" gorm:"type:varchar(100);not null;comment:主标题"`
	SubTitle    string     `json:"sub_title" gorm:"type:varchar(100);comment:副标题"`
	Sort        int        `json:"sort" gorm:"not null;default:0;comment:排序"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (PriceListCategory) TableName() string {
	return "price_list_categories"
}

// PriceListItem 价目单商品（关联供应商商品）
type PriceListItem struct {
	ID          uint               `json:"id" gorm:"primaryKey;autoIncrement"`
	CategoryID  uint               `json:"category_id" gorm:"not null;index;comment:价目单分类ID"`
	Category    *PriceListCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	ProductID   uint               `json:"product_id" gorm:"not null;index;comment:供应商商品ID"`
	Product     *SupplierProduct   `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	DisplayName string             `json:"display_name" gorm:"type:varchar(200);comment:显示名称（可覆盖商品名称）"`
	Price       float64            `json:"price" gorm:"type:decimal(10,2);not null;comment:价格"`
	Unit        string             `json:"unit" gorm:"type:varchar(20);comment:单位"`
	Spec        string             `json:"spec" gorm:"type:varchar(100);comment:规格说明"`
	Sort        int                `json:"sort" gorm:"not null;default:0;comment:排序"`
	Status      int8               `json:"status" gorm:"not null;default:1;comment:状态 1=显示 0=隐藏"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

func (PriceListItem) TableName() string {
	return "price_list_items"
}

// CreatePriceListReq 创建价目单请求
type CreatePriceListReq struct {
	StoreID     uint   `json:"store_id" binding:"required"`
	Name        string `json:"name" binding:"required,max=200"`
	Logo        string `json:"logo" binding:"max=500"`
	Description string `json:"description" binding:"max=500"`
	IsDefault   int8   `json:"is_default"`
}

// UpdatePriceListReq 更新价目单请求
type UpdatePriceListReq struct {
	Name        string `json:"name" binding:"max=200"`
	Logo        string `json:"logo" binding:"max=500"`
	Description string `json:"description" binding:"max=500"`
	Status      *int8  `json:"status" binding:"omitempty,oneof=0 1"`
	IsDefault   *int8  `json:"is_default" binding:"omitempty,oneof=0 1"`
}

// CreatePriceListCategoryReq 创建价目单分类请求
type CreatePriceListCategoryReq struct {
	PriceListID uint   `json:"price_list_id" binding:"required"`
	MainTitle   string `json:"main_title" binding:"required,max=100"`
	SubTitle    string `json:"sub_title" binding:"max=100"`
	Sort        int    `json:"sort"`
}

// UpdatePriceListCategoryReq 更新价目单分类请求
type UpdatePriceListCategoryReq struct {
	MainTitle string `json:"main_title" binding:"max=100"`
	SubTitle  string `json:"sub_title" binding:"max=100"`
	Sort      *int   `json:"sort"`
}

// AddPriceListItemReq 添加价目单商品请求
type AddPriceListItemReq struct {
	CategoryID  uint    `json:"category_id" binding:"required"`
	ProductID   uint    `json:"product_id" binding:"required"`
	DisplayName string  `json:"display_name" binding:"max=200"`
	Price       float64 `json:"price" binding:"required,gte=0"`
	Unit        string  `json:"unit" binding:"max=20"`
	Spec        string  `json:"spec" binding:"max=100"`
	Sort        int     `json:"sort"`
}

// UpdatePriceListItemReq 更新价目单商品请求
type UpdatePriceListItemReq struct {
	DisplayName string   `json:"display_name" binding:"max=200"`
	Price       *float64 `json:"price" binding:"omitempty,gte=0"`
	Unit        string   `json:"unit" binding:"max=20"`
	Spec        string   `json:"spec" binding:"max=100"`
	Sort        *int     `json:"sort"`
	Status      *int8    `json:"status" binding:"omitempty,oneof=0 1"`
}

// BatchAddItemEntry 批量添加时单个商品条目
type BatchAddItemEntry struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	DisplayName string  `json:"display_name" binding:"max=200"`
	Price       float64 `json:"price" binding:"gte=0"`
	Unit        string  `json:"unit" binding:"max=20"`
	Spec        string  `json:"spec" binding:"max=100"`
	Sort        int     `json:"sort"`
}

// BatchAddPriceListItemsReq 批量添加价目单商品请求
type BatchAddPriceListItemsReq struct {
	CategoryID uint                `json:"category_id" binding:"required"`
	Products   []BatchAddItemEntry `json:"products" binding:"required,min=1"`
}

// PriceListResp 价目单响应（包含完整结构）
type PriceListResp struct {
	PriceList
	Categories []PriceListCategoryResp `json:"categories"`
}

// PriceListCategoryResp 价目单分类响应
type PriceListCategoryResp struct {
	PriceListCategory
	Items []PriceListItemResp `json:"items"`
}

// PriceListItemResp 价目单商品响应
type PriceListItemResp struct {
	PriceListItem
	ProductName  string `json:"product_name"`  // 原商品名称
	SupplierName string `json:"supplier_name"` // 供应商名称
	ProductImage string `json:"product_image"` // 商品图片
}
