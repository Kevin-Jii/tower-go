package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// StoreAccountListScope 门店记账列表数据范围（与 model.ListStoreAccountReq 配合）。
func StoreAccountListScope(req *model.ListStoreAccountReq) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if req == nil {
			return db.Where("1 = 0")
		}
		if req.DataScope == model.DataScopeAll {
			if req.StoreID > 0 {
				return db.Where("store_accounts.store_id = ?", req.StoreID)
			}
			return db
		}

		switch req.DataScope {
		case model.DataScopeSelf:
			q := db
			if req.StoreID > 0 {
				q = q.Where("store_accounts.store_id = ?", req.StoreID)
			}
			return q.Where("store_accounts.operator_id = ?", req.UserID)
		default:
			if req.StoreID > 0 {
				return db.Where("store_accounts.store_id = ?", req.StoreID)
			}
			return db
		}
	}
}
