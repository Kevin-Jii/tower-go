package tenant

import "gorm.io/gorm"

// IsolationStrategy 数据隔离策略接口（策略模式）
type IsolationStrategy interface {
	// Apply 应用隔离策略到查询
	Apply(db *gorm.DB, tenant *TenantContext) *gorm.DB
	// GetColumnName 获取隔离字段名
	GetColumnName() string
}

// StoreIsolationStrategy 门店隔离策略
type StoreIsolationStrategy struct {
	columnName string
}

func NewStoreIsolationStrategy() *StoreIsolationStrategy {
	return &StoreIsolationStrategy{columnName: "store_id"}
}

func (s *StoreIsolationStrategy) Apply(db *gorm.DB, tenant *TenantContext) *gorm.DB {
	if tenant == nil || tenant.IsAdmin || tenant.TenantID == 0 {
		return db
	}
	return db.Where(s.columnName+" = ?", tenant.TenantID)
}

func (s *StoreIsolationStrategy) GetColumnName() string {
	return s.columnName
}

// SupplierIsolationStrategy 供应商隔离策略
type SupplierIsolationStrategy struct {
	columnName string
}

func NewSupplierIsolationStrategy() *SupplierIsolationStrategy {
	return &SupplierIsolationStrategy{columnName: "supplier_id"}
}

func (s *SupplierIsolationStrategy) Apply(db *gorm.DB, tenant *TenantContext) *gorm.DB {
	if tenant == nil || tenant.IsAdmin || tenant.TenantID == 0 {
		return db
	}
	return db.Where(s.columnName+" = ?", tenant.TenantID)
}

func (s *SupplierIsolationStrategy) GetColumnName() string {
	return s.columnName
}

// NoIsolationStrategy 无隔离策略（管理员使用）
type NoIsolationStrategy struct{}

func NewNoIsolationStrategy() *NoIsolationStrategy {
	return &NoIsolationStrategy{}
}

func (s *NoIsolationStrategy) Apply(db *gorm.DB, tenant *TenantContext) *gorm.DB {
	return db
}

func (s *NoIsolationStrategy) GetColumnName() string {
	return ""
}

// CompositeIsolationStrategy 组合隔离策略（多字段隔离）
type CompositeIsolationStrategy struct {
	strategies []IsolationStrategy
}

func NewCompositeIsolationStrategy(strategies ...IsolationStrategy) *CompositeIsolationStrategy {
	return &CompositeIsolationStrategy{strategies: strategies}
}

func (s *CompositeIsolationStrategy) Apply(db *gorm.DB, tenant *TenantContext) *gorm.DB {
	for _, strategy := range s.strategies {
		db = strategy.Apply(db, tenant)
	}
	return db
}

func (s *CompositeIsolationStrategy) GetColumnName() string {
	return "composite"
}
