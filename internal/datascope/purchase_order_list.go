package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// PurchaseOrderListScope 采购单列表数据范围（与 model.ListPurchaseOrderReq 配合）。
func PurchaseOrderListScope(req *model.ListPurchaseOrderReq) func(db *gorm.DB) *gorm.DB {
	if req == nil {
		return func(db *gorm.DB) *gorm.DB { return db.Where("1 = 0") }
	}
	s := listRBACSnap{DataScope: req.DataScope, StoreID: req.StoreID, UserID: req.UserID}
	return listDataScopeScope(s, PolicyPurchaseOrders, false)
}
