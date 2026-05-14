package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// InventoryOrderListScope 出入库单列表数据范围（与 model.ListInventoryOrderReq 配合）。
// 条件字段与主表 inventory_orders 一致（无别名）。
func InventoryOrderListScope(req *model.ListInventoryOrderReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req == nil {
			return db.Where("1 = 0")
		}
		if req.DataScope == model.DataScopeAll {
			if req.StoreID > 0 {
				return db.Where("store_id = ?", req.StoreID)
			}
			return db
		}

		switch req.DataScope {
		case model.DataScopeSelf:
			q := db
			if req.StoreID > 0 {
				q = q.Where("store_id = ?", req.StoreID)
			}
			return q.Where("operator_id = ?", req.UserID)
		default:
			if req.StoreID > 0 {
				return db.Where("store_id = ?", req.StoreID)
			}
			return db
		}
	}
}
