package tenant

import (
	"gorm.io/gorm"
)

// TenantRepository 带租户隔离的通用仓储（代理模式 + 策略模式）
type TenantRepository[T any] struct {
	db       *gorm.DB
	strategy IsolationStrategy
}

// NewTenantRepository 创建租户仓储
func NewTenantRepository[T any](db *gorm.DB, strategy IsolationStrategy) *TenantRepository[T] {
	return &TenantRepository[T]{
		db:       db,
		strategy: strategy,
	}
}

// WithTenant 获取带租户隔离的查询
func (r *TenantRepository[T]) WithTenant(tenant *TenantContext) *gorm.DB {
	return r.strategy.Apply(r.db.Model(new(T)), tenant)
}

// Create 创建记录（自动设置租户ID）
func (r *TenantRepository[T]) Create(tenant *TenantContext, entity *T) error {
	return r.db.Create(entity).Error
}

// FindByID 根据ID查找（带租户隔离）
func (r *TenantRepository[T]) FindByID(tenant *TenantContext, id uint) (*T, error) {
	var entity T
	err := r.WithTenant(tenant).First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll 查找所有（带租户隔离）
func (r *TenantRepository[T]) FindAll(tenant *TenantContext) ([]T, error) {
	var entities []T
	err := r.WithTenant(tenant).Find(&entities).Error
	return entities, err
}

// Update 更新记录（带租户隔离验证）
func (r *TenantRepository[T]) Update(tenant *TenantContext, id uint, updates map[string]interface{}) error {
	result := r.WithTenant(tenant).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return ErrAccessDenied
	}
	return result.Error
}

// Delete 删除记录（带租户隔离验证）
func (r *TenantRepository[T]) Delete(tenant *TenantContext, id uint) error {
	result := r.WithTenant(tenant).Delete(new(T), id)
	if result.RowsAffected == 0 {
		return ErrAccessDenied
	}
	return result.Error
}

// Count 统计数量（带租户隔离）
func (r *TenantRepository[T]) Count(tenant *TenantContext) (int64, error) {
	var count int64
	err := r.WithTenant(tenant).Count(&count).Error
	return count, err
}

// Exists 检查是否存在（带租户隔离）
func (r *TenantRepository[T]) Exists(tenant *TenantContext, id uint) (bool, error) {
	var count int64
	err := r.WithTenant(tenant).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Query 自定义查询（带租户隔离）
func (r *TenantRepository[T]) Query(tenant *TenantContext) *gorm.DB {
	return r.WithTenant(tenant)
}

// Raw 原始查询（不带隔离，谨慎使用）
func (r *TenantRepository[T]) Raw() *gorm.DB {
	return r.db.Model(new(T))
}
