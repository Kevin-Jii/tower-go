package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// InventoryOrderListScope 出入库单列表数据范围（与 model.ListInventoryOrderReq 配合）。
func InventoryOrderListScope(req *model.ListInventoryOrderReq) func(db *gorm.DB) *gorm.DB {
	if req == nil {
		return func(db *gorm.DB) *gorm.DB { return db.Where("1 = 0") }
	}
	s := listRBACSnap{DataScope: req.DataScope, StoreID: req.StoreID, UserID: req.UserID}
	return listDataScopeScope(s, PolicyInventoryOrders, false)
}
