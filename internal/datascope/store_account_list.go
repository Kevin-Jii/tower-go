package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// StoreAccountListScope 门店记账列表数据范围（与 model.ListStoreAccountReq 配合）。
func StoreAccountListScope(req *model.ListStoreAccountReq) func(db *gorm.DB) *gorm.DB {
	if req == nil {
		return func(db *gorm.DB) *gorm.DB { return db.Where("1 = 0") }
	}
	s := listRBACSnap{DataScope: req.DataScope, StoreID: req.StoreID, UserID: req.UserID}
	return listDataScopeScope(s, PolicyStoreAccounts, false)
}
