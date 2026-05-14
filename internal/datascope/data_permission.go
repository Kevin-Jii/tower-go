// Package datascope 提供统一 GORM Scope（P0 起与 pkg/datascope 业务函数并存，逐步收敛）。
package datascope

import (
	"github.com/Kevin-Jii/tower-go/internal/authctx"
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// TablePolicy 描述当前查询主表上用于数据隔离的列（须带表名或别名，如 purchase_orders.store_id）。
type TablePolicy struct {
	StoreColumn   string
	CreatorColumn string // 可选；SELF 时使用
}

// DataPermission 按 AuthContext 的有效数据范围追加 WHERE（不含「总部按 query 筛门店」等特殊逻辑，见各业务 Scope）。
func DataPermission(ac *authctx.Context, p TablePolicy) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if ac == nil {
			return db.Where("1 = 0")
		}
		switch ac.EffectiveDataScope {
		case model.DataScopeAll:
			return db
		case model.DataScopeSelf:
			q := db
			if p.StoreColumn != "" && ac.StoreID > 0 {
				q = q.Where(p.StoreColumn+" = ?", ac.StoreID)
			}
			if p.CreatorColumn != "" {
				return q.Where(p.CreatorColumn+" = ?", ac.UserID)
			}
			if p.StoreColumn != "" {
				return q
			}
			return db.Where("1 = 0")
		case model.DataScopeTenant, model.DataScopeStore:
			fallthrough
		default:
			if len(ac.CustomStoreIDs) > 0 {
				if p.StoreColumn == "" {
					return db.Where("1 = 0")
				}
				return db.Where(p.StoreColumn+" IN ?", ac.CustomStoreIDs)
			}
			if ac.StoreID > 0 && p.StoreColumn != "" {
				return db.Where(p.StoreColumn+" = ?", ac.StoreID)
			}
			return db.Where("1 = 0")
		}
	}
}
