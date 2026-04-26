package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// ApplyPurchaseOrdersList 采购单列表数据范围（与 ListPurchaseOrderReq 配合使用）
func ApplyPurchaseOrdersList(db *gorm.DB, req *model.ListPurchaseOrderReq) *gorm.DB {
	if req.RoleCode == model.RoleCodeAdmin || req.RoleCode == model.RoleCodeSuperAdmin || req.DataScope == model.DataScopeAll {
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

// ApplyStoreAccountsList 门店记账列表
func ApplyStoreAccountsList(db *gorm.DB, req *model.ListStoreAccountReq) *gorm.DB {
	if req.RoleCode == model.RoleCodeAdmin || req.RoleCode == model.RoleCodeSuperAdmin || req.DataScope == model.DataScopeAll {
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

// ApplyInventoriesList 库存列表（表别名 i）
func ApplyInventoriesList(db *gorm.DB, req *model.ListInventoryReq) *gorm.DB {
	if req.RoleCode == model.RoleCodeAdmin || req.RoleCode == model.RoleCodeSuperAdmin || req.DataScope == model.DataScopeAll {
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

// ApplyInventoryOrdersList 出入库单列表
func ApplyInventoryOrdersList(db *gorm.DB, req *model.ListInventoryOrderReq) *gorm.DB {
	if req.RoleCode == model.RoleCodeAdmin || req.RoleCode == model.RoleCodeSuperAdmin || req.DataScope == model.DataScopeAll {
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
