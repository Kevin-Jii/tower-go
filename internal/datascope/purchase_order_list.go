package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// PurchaseOrderListScope 采购单列表数据范围（与 model.ListPurchaseOrderReq 配合）。
// 逻辑与历史 pkg/datascope.ApplyPurchaseOrdersList 一致，总部未绑店 + DataScopeAll 时可用 req.StoreId 筛选门店。
func PurchaseOrderListScope(req *model.ListPurchaseOrderReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req == nil {
			return db.Where("1 = 0")
		}
		if req.DataScope == model.DataScopeAll {
			if req.StoreID > 0 {
				return db.Where("purchase_orders.store_id = ?", req.StoreID)
			}
			return db
		}

		switch req.DataScope {
		case model.DataScopeSelf:
			q := db
			if req.StoreID > 0 {
				q = q.Where("purchase_orders.store_id = ?", req.StoreID)
			}
			return q.Where("purchase_orders.created_by = ?", req.UserID)
		case model.DataScopeTenant, model.DataScopeStore:
			fallthrough
		default:
			if req.StoreID > 0 {
				return db.Where("purchase_orders.store_id = ?", req.StoreID)
			}
			return db
		}
	}
}
