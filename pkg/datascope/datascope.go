package datascope

import (
	idscope "github.com/Kevin-Jii/tower-go/internal/datascope"
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// ApplyPurchaseOrdersList 采购单列表数据范围（与 ListPurchaseOrderReq 配合使用）
func ApplyPurchaseOrdersList(db *gorm.DB, req *model.ListPurchaseOrderReq) *gorm.DB {
	return db.Scopes(idscope.PurchaseOrderListScope(req))
}

// ApplyStoreAccountsList 门店记账列表
func ApplyStoreAccountsList(db *gorm.DB, req *model.ListStoreAccountReq) *gorm.DB {
	return db.Scopes(idscope.StoreAccountListScope(req))
}

// ApplyInventoriesList 库存列表（表别名 i）
func ApplyInventoriesList(db *gorm.DB, req *model.ListInventoryReq) *gorm.DB {
	return db.Scopes(idscope.InventoryListScope(req))
}

// ApplyInventoryOrdersList 出入库单列表
func ApplyInventoryOrdersList(db *gorm.DB, req *model.ListInventoryOrderReq) *gorm.DB {
	return db.Scopes(idscope.InventoryOrderListScope(req))
}
