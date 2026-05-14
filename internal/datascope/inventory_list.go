package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// InventoryListScope 库存列表数据范围（表别名 i，与 model.ListInventoryReq 配合）。
func InventoryListScope(req *model.ListInventoryReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req == nil {
			return db.Where("1 = 0")
		}
		if req.DataScope == model.DataScopeAll {
			if req.StoreID > 0 {
				return db.Where("i.store_id = ?", req.StoreID)
			}
			return db
		}

		switch req.DataScope {
		case model.DataScopeSelf:
			// 库存行无创建人字段，个人范围与门店范围一致
			fallthrough
		default:
			if req.StoreID > 0 {
				return db.Where("i.store_id = ?", req.StoreID)
			}
			return db
		}
	}
}
