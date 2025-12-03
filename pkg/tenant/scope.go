package tenant

import (
	"gorm.io/gorm"
)

// TenantScope GORM 作用域（装饰器模式）
// 自动为查询添加租户过滤条件
type TenantScope struct {
	tenant   *TenantContext
	strategy IsolationStrategy
}

// NewTenantScope 创建租户作用域
func NewTenantScope(tenant *TenantContext, strategy IsolationStrategy) *TenantScope {
	return &TenantScope{
		tenant:   tenant,
		strategy: strategy,
	}
}

// Apply 应用作用域（实现 gorm 的 Scope 接口）
func (s *TenantScope) Apply(db *gorm.DB) *gorm.DB {
	return s.strategy.Apply(db, s.tenant)
}

// TenantDB 带租户隔离的数据库包装器（代理模式）
type TenantDB struct {
	db       *gorm.DB
	tenant   *TenantContext
	strategy IsolationStrategy
}

// NewTenantDB 创建带租户隔离的数据库
func NewTenantDB(db *gorm.DB, tenant *TenantContext, strategy IsolationStrategy) *TenantDB {
	return &TenantDB{
		db:       db,
		tenant:   tenant,
		strategy: strategy,
	}
}

// DB 获取带隔离的数据库实例
func (t *TenantDB) DB() *gorm.DB {
	return t.strategy.Apply(t.db, t.tenant)
}

// Raw 获取原始数据库（不带隔离，谨慎使用）
func (t *TenantDB) Raw() *gorm.DB {
	return t.db
}

// WithStrategy 切换隔离策略
func (t *TenantDB) WithStrategy(strategy IsolationStrategy) *TenantDB {
	return &TenantDB{
		db:       t.db,
		tenant:   t.tenant,
		strategy: strategy,
	}
}

// SkipIsolation 跳过隔离（管理员操作）
func (t *TenantDB) SkipIsolation() *gorm.DB {
	return t.db
}

// Scopes GORM Scopes 辅助函数
func TenantScopes(tenant *TenantContext, strategy IsolationStrategy) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return strategy.Apply(db, tenant)
	}
}

// StoreScope 门店隔离作用域（便捷函数）
func StoreScope(storeID uint) func(*gorm.DB) *gorm.DB {
	tenant := &TenantContext{TenantID: storeID, TenantType: TenantTypeStore}
	strategy := NewStoreIsolationStrategy()
	return TenantScopes(tenant, strategy)
}

// SupplierScope 供应商隔离作用域（便捷函数）
func SupplierScope(supplierID uint) func(*gorm.DB) *gorm.DB {
	tenant := &TenantContext{TenantID: supplierID, TenantType: TenantTypeSupplier}
	strategy := NewSupplierIsolationStrategy()
	return TenantScopes(tenant, strategy)
}

// AdminScope 管理员作用域（无隔离）
func AdminScope() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db
	}
}
