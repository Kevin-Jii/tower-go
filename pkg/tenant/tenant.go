package tenant

import (
	"context"
	"errors"



)

// 上下文 key
type contextKey string

const (
	TenantIDKey   contextKey = "tenant_id"
	TenantTypeKey contextKey = "tenant_type"
)

// TenantType 租户类型
type TenantType string

const (
	TenantTypeStore    TenantType = "store"    // 门店
	TenantTypeSupplier TenantType = "supplier" // 供应商
	TenantTypeAdmin    TenantType = "admin"    // 总部（无隔离）
)

// TenantContext 租户上下文
type TenantContext struct {
	TenantID   uint
	TenantType TenantType
	IsAdmin    bool // 是否管理员（跳过隔离）
}

// NewTenantContext 创建租户上下文
func NewTenantContext(tenantID uint, tenantType TenantType, isAdmin bool) *TenantContext {
	return &TenantContext{
		TenantID:   tenantID,
		TenantType: tenantType,
		IsAdmin:    isAdmin,
	}
}

// WithTenant 将租户信息注入 context
func WithTenant(ctx context.Context, tenant *TenantContext) context.Context {
	ctx = context.WithValue(ctx, TenantIDKey, tenant.TenantID)
	ctx = context.WithValue(ctx, TenantTypeKey, tenant.TenantType)
	return ctx
}

// GetTenant 从 context 获取租户信息
func GetTenant(ctx context.Context) *TenantContext {
	tenantID, _ := ctx.Value(TenantIDKey).(uint)
	tenantType, _ := ctx.Value(TenantTypeKey).(TenantType)
	return &TenantContext{
		TenantID:   tenantID,
		TenantType: tenantType,
	}
}

// GetTenantID 从 context 获取租户ID
func GetTenantID(ctx context.Context) uint {
	if id, ok := ctx.Value(TenantIDKey).(uint); ok {
		return id
	}
	return 0
}

var ErrNoTenant = errors.New("tenant context not found")
var ErrAccessDenied = errors.New("access denied: data belongs to another tenant")
