package datascope

import (
	"github.com/Kevin-Jii/tower-go/model"

	"gorm.io/gorm"
)

// listRBACSnap 列表接口注入的数据范围三要素（与 model.List*Req 上 json:"-" 字段一致）。
type listRBACSnap struct {
	DataScope int8
	StoreID   uint
	UserID    uint
}

// listDataScopeScope 统一「列表类」DataScope 分支（与历史 pkg/datascope 行为对齐）。
//
// selfSameAsStore：为 true 时 DataScopeSelf 仅追加门店条件、不追加本人列（库存行无创建人场景）。
func listDataScopeScope(s listRBACSnap, p TablePolicy, selfSameAsStore bool) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if s.DataScope == model.DataScopeAll {
			if s.StoreID > 0 && p.StoreColumn != "" {
				return db.Where(p.StoreColumn+" = ?", s.StoreID)
			}
			return db
		}

		switch s.DataScope {
		case model.DataScopeSelf:
			q := db
			if s.StoreID > 0 && p.StoreColumn != "" {
				q = q.Where(p.StoreColumn+" = ?", s.StoreID)
			}
			if selfSameAsStore || p.CreatorColumn == "" {
				return q
			}
			return q.Where(p.CreatorColumn+" = ?", s.UserID)
		case model.DataScopeTenant, model.DataScopeStore:
			fallthrough
		default:
			if s.StoreID > 0 && p.StoreColumn != "" {
				return db.Where(p.StoreColumn+" = ?", s.StoreID)
			}
			return db
		}
	}
}
