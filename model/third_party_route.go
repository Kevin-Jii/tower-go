package model

import "time"

type ThirdPartyRoute struct {
	ID        uint                   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string                 `json:"name" gorm:"type:varchar(100);not null;comment:路线名称"`
	Remark    string                 `json:"remark" gorm:"type:varchar(500);comment:备注"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Stores    []ThirdPartyRouteStore `json:"stores,omitempty" gorm:"foreignKey:RouteID"`
}

func (ThirdPartyRoute) TableName() string {
	return "third_party_routes"
}

type ThirdPartyRouteStore struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	RouteID   uint      `json:"route_id" gorm:"not null;index;comment:路线ID"`
	StoreID   uint      `json:"store_id" gorm:"not null;index;comment:门店ID"`
	Sort      int       `json:"sort" gorm:"not null;default:0;comment:路线顺序"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Store     *Store    `json:"store,omitempty" gorm:"foreignKey:StoreID"`
}

func (ThirdPartyRouteStore) TableName() string {
	return "third_party_route_stores"
}

type UpsertThirdPartyRouteReq struct {
	Name     string `json:"name" binding:"required,max=100"`
	Remark   string `json:"remark" binding:"max=500"`
	StoreIDs []uint `json:"store_ids"`
}

type ImportRouteOrdersReq struct {
	StartDate string `json:"start_date" binding:"required,len=10"`
	EndDate   string `json:"end_date" binding:"required,len=10"`
}

type RouteStoreQuantity struct {
	StoreID   uint    `json:"store_id"`
	StoreName string  `json:"store_name"`
	Quantity  float64 `json:"quantity"`
}

type RouteImportedProductRow struct {
	ProductName string               `json:"product_name"`
	TotalQty    float64              `json:"total_qty"`
	StoreQty    []RouteStoreQuantity `json:"store_qty"`
}

type SaveRouteSheetReq struct {
	StartDate string      `json:"start_date" binding:"required,len=10"`
	EndDate   string      `json:"end_date" binding:"required,len=10"`
	Headers   []string    `json:"headers"`
	Rows      [][]float64 `json:"rows"`
	Products  []string    `json:"products"`
}

type ThirdPartyLogisticsSheet struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	RouteID      uint      `json:"route_id" gorm:"not null;uniqueIndex:uk_route_logistics_period,priority:1;comment:路线ID"`
	SheetDate    string    `json:"sheet_date" gorm:"type:varchar(10);not null;index;comment:最近一次保存日期(yyyy-mm-dd)"`
	StartDate    string    `json:"start_date" gorm:"type:varchar(10);not null;uniqueIndex:uk_route_logistics_period,priority:2;comment:导入开始日期"`
	EndDate      string    `json:"end_date" gorm:"type:varchar(10);not null;uniqueIndex:uk_route_logistics_period,priority:3;comment:导入结束日期"`
	HeadersJSON  string    `json:"headers_json" gorm:"type:longtext;comment:表头JSON"`
	RowsJSON     string    `json:"rows_json" gorm:"type:longtext;comment:数量矩阵JSON"`
	ProductsJSON string    `json:"products_json" gorm:"type:longtext;comment:商品名JSON"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (ThirdPartyLogisticsSheet) TableName() string {
	return "third_party_logistics_sheets"
}
