package tenant

import (
	"github.com/gin-gonic/gin"
)

// GinTenantMiddleware Gin 租户中间件
// 从请求上下文提取租户信息并注入到 context
func GinTenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 gin context 获取已认证的用户信息
		storeID, _ := c.Get("storeID")
		roleCode, _ := c.Get("roleCode")

		var tenantCtx *TenantContext

		// 根据角色判断租户类型
		if roleCode == "admin" || roleCode == "super_admin" {
			tenantCtx = &TenantContext{
				TenantID:   0,
				TenantType: TenantTypeAdmin,
				IsAdmin:    true,
			}
		} else if sid, ok := storeID.(uint); ok && sid > 0 {
			tenantCtx = &TenantContext{
				TenantID:   sid,
				TenantType: TenantTypeStore,
				IsAdmin:    false,
			}
		}

		if tenantCtx != nil {
			// 注入到 gin context
			c.Set("tenant", tenantCtx)
			// 注入到 request context
			ctx := WithTenant(c.Request.Context(), tenantCtx)
			c.Request = c.Request.WithContext(ctx)
		}

		c.Next()
	}
}

// GetTenantFromGin 从 Gin context 获取租户信息
func GetTenantFromGin(c *gin.Context) *TenantContext {
	if tenant, exists := c.Get("tenant"); exists {
		if t, ok := tenant.(*TenantContext); ok {
			return t
		}
	}
	return nil
}

// MustGetTenant 必须获取租户信息，否则返回错误
func MustGetTenant(c *gin.Context) (*TenantContext, error) {
	tenant := GetTenantFromGin(c)
	if tenant == nil {
		return nil, ErrNoTenant
	}
	return tenant, nil
}

// RequireTenant 要求必须有租户信息的中间件
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenant := GetTenantFromGin(c)
		if tenant == nil || (tenant.TenantID == 0 && !tenant.IsAdmin) {
			c.JSON(403, gin.H{"error": "tenant context required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
