package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// InventoryListScope 库存列表数据范围（表别名 i，与 model.ListInventoryReq 配合）。
func InventoryListScope(req *model.ListInventoryReq) func(db *gorm.DB) *gorm.DB {
	if req == nil {
		return func(db *gorm.DB) *gorm.DB { return db.Where("1 = 0") }
	}
	s := listRBACSnap{DataScope: req.DataScope, StoreID: req.StoreID, UserID: req.UserID}
	// 库存行无创建人：SELF 与门店范围一致
	return listDataScopeScope(s, PolicyInventoriesAsI, true)
}
