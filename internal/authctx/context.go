// Package authctx 承载「请求级」鉴权与数据权限上下文（P0：从 Gin 注入，供后续 Repository 统一 Scopes 使用）。
package authctx

import (
	"context"
)

// GinKey Gin Context 中存储 *Context 的键（与 middleware 约定一致）。
const GinKey = "authctx"

// GinGetter 由 *gin.Context 实现，用于 FromGin 而不在本包直接依赖 gin。
type GinGetter interface {
	Get(key any) (value any, exists bool)
}

// Context 权限与数据范围快照（JWT + 角色解析后的有效值）。
//
// EffectiveDataScope（D4 单一真源之一）：与 middleware.GetDataScope(c) 同步写入，表示行级范围枚举。
// 「全库」仅当总部未绑店 admin/super_admin（见 HQUnboundAdmin）；绑店超管/管理员一律按门店降级，禁止在 Controller 用 role_code 手写全库分支。
type Context struct {
	UserID             uint
	StoreID            uint
	RoleID             uint
	RoleCode           string
	EffectiveDataScope int8
	HQUnbound          bool

	// 预留：组织 / SaaS / 自定义门店范围
	DeptID         uint
	TenantID       uint
	CustomStoreIDs []uint
}

type ctxKey struct{}

// WithContext 将 Context 写入标准 context（供脱离 Gin 的 service/repository 链路使用）。
func WithContext(parent context.Context, a *Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithValue(parent, ctxKey{}, a)
}

// FromContext 从标准 context 读取；未注入则返回 nil。
func FromContext(ctx context.Context) *Context {
	if ctx == nil {
		return nil
	}
	v, ok := ctx.Value(ctxKey{}).(*Context)
	if !ok || v == nil {
		return nil
	}
	return v
}

// FromGin 从 Gin 风格 Context 读取；未挂载则返回 nil。
func FromGin(c GinGetter) *Context {
	if c == nil {
		return nil
	}
	v, ok := c.Get(GinKey)
	if !ok {
		return nil
	}
	a, ok := v.(*Context)
	if !ok || a == nil {
		return nil
	}
	return a
}

// MustFromGin 与 FromGin 相同，但若未挂载会 panic（仅用于确信已走 AuthMiddleware 的路径）。
func MustFromGin(c GinGetter) *Context {
	a := FromGin(c)
	if a == nil {
		panic("authctx: missing on gin.Context; ensure AuthMiddleware runs before handler")
	}
	return a
}
