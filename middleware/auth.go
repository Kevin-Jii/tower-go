package middleware

import (
	"fmt"
	"strings"

	"github.com/Kevin-Jii/tower-go/internal/authctx"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/utils/auth"
	"github.com/Kevin-Jii/tower-go/utils/database"
	"github.com/Kevin-Jii/tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			http.ErrorApp(c, apicode.AuthHeaderRequired)
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			http.ErrorApp(c, apicode.AuthHeaderFormat)
			c.Abort()
			return
		}

		// 解析token
		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			http.ErrorApp(c, apicode.TokenInvalid)
			c.Abort()
			return
		}

		// 将用户信息存储到上下文（优先使用 Token 中的值，安全可靠）
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("storeID", claims.StoreID)
		c.Set("roleCode", claims.RoleCode)
		c.Set("roleID", claims.RoleID)

		// 加载角色（数据权限 data_scope）
		if claims.RoleID > 0 && database.DB != nil {
			var role model.Role
			if err := database.DB.First(&role, claims.RoleID).Error; err == nil {
				c.Set("roleModel", &role)
			}
		}

		// 统一 AuthContext（数据权限 Pipeline P0，供 internal/datascope 与后续 Repository 使用）
		c.Set(authctx.GinKey, buildAuthContext(c, claims))

		// 同时保存请求头中的值（供日志或特殊场景使用）
		c.Set("headerUserID", c.GetHeader("X-User-Id"))
		c.Set("headerStoreID", c.GetHeader("X-Store-Id"))

		c.Next()
	}
}

func buildAuthContext(c *gin.Context, claims *auth.Claims) *authctx.Context {
	return &authctx.Context{
		UserID:             claims.UserID,
		StoreID:            claims.StoreID,
		RoleID:             claims.RoleID,
		RoleCode:           claims.RoleCode,
		EffectiveDataScope: GetDataScope(c),
		HQUnbound:          HQUnboundAdmin(c),
	}
}

// AttachAuthContextToHTTPRequest 将 Gin 上的 AuthContext 写入标准 http.Request.Context，
// 供 service / repository 使用 authctx.FromContext(r.Context())。
func AttachAuthContextToHTTPRequest(c *gin.Context) {
	if ac := authctx.FromGin(c); ac != nil {
		r := c.Request
		c.Request = r.WithContext(authctx.WithContext(r.Context(), ac))
	}
}

// StoreAuthMiddleware 门店鉴权中间件 (别名，向后兼容)
func StoreAuthMiddleware() gin.HandlerFunc {
	return AuthMiddleware()
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("userID"); exists {
		return userID.(uint)
	}
	return 0
}

// GetStoreID 从上下文获取门店 ID
func GetStoreID(c *gin.Context) uint {
	if storeID, exists := c.Get("storeID"); exists {
		return storeID.(uint)
	}
	return 0
}

// GetRoleCode 从上下文获取角色代码
func GetRoleCode(c *gin.Context) string {
	if roleCode, exists := c.Get("roleCode"); exists {
		return roleCode.(string)
	}
	return ""
}

// GetRoleID 从上下文获取角色 ID
func GetRoleID(c *gin.Context) uint {
	if roleID, exists := c.Get("roleID"); exists {
		return roleID.(uint)
	}
	return 0
}

// IsAdmin 判断是否是管理员角色。注意：绑定门店的 admin 仍应按门店隔离数据，
// 是否可跨店请使用 HQUnboundAdmin。
func IsAdmin(c *gin.Context) bool {
	roleCode := GetRoleCode(c)
	return roleCode == model.RoleCodeAdmin || roleCode == model.RoleCodeSuperAdmin
}

// HQUnboundAdmin 未绑定门店的总部 admin / super_admin：可跨店查看/操作数据。
// Token 中 store_id>0 时一律视为已绑店，与角色码无关（绑店 super_admin 也仅本店）。
func HQUnboundAdmin(c *gin.Context) bool {
	return model.HQUnboundAdminRole(GetRoleCode(c), GetStoreID(c))
}

// ResolveQueryStoreID 解析列表/统计接口可访问的门店范围。
// 总部未绑定 admin / super_admin 可以按 query store_id 查看指定门店，缺省或 0 表示全部；
// 绑定门店的 admin、门店管理员、员工一律只能看 Token 里的本店，忽略 query store_id。
func ResolveQueryStoreID(c *gin.Context, queryKey string) uint {
	if !HQUnboundAdmin(c) {
		return GetStoreID(c)
	}
	if queryKey == "" {
		queryKey = "store_id"
	}
	raw := strings.TrimSpace(c.Query(queryKey))
	if raw == "" || raw == "0" {
		return 0
	}
	var sid uint
	if _, err := fmt.Sscanf(raw, "%d", &sid); err != nil {
		return 0
	}
	return sid
}

// GetRoleModel 已加载的角色（可能为 nil）
func GetRoleModel(c *gin.Context) *model.Role {
	if v, ok := c.Get("roleModel"); ok {
		if r, ok := v.(*model.Role); ok {
			return r
		}
	}
	return nil
}

// GetDataScope 返回本请求用于「行级数据范围」的 DataScope 枚举（D4：单一真源，勿在 Controller 用角色码推断全库）。
//
// 规则摘要：
//  1) HQUnboundAdmin(c) 为真（admin 或 super_admin 且 token 中 store_id==0）→ DataScopeAll；列表可再按 query store_id 筛店。
//  2) token 中 store_id>0 且角色在库中为「全部」→ 强制降为 DataScopeStore。
//  3) 其余使用 roles.data_scope；无角色模型时默认 DataScopeStore。
//
// 同步写入 authctx.Context.EffectiveDataScope；HTTP 列表请求的 DataScope 字段由 service 从 authctx 注入（见 list_authctx.go）。
func GetDataScope(c *gin.Context) int8 {
	if HQUnboundAdmin(c) {
		return model.DataScopeAll
	}
	if r := GetRoleModel(c); r != nil {
		if r.DataScope == model.DataScopeAll && GetStoreID(c) > 0 {
			return model.DataScopeStore
		}
		return r.DataScope
	}
	return model.DataScopeStore
}

// ListRBAC 返回 dataScope、userID、roleCode（dataScope 同 GetDataScope）。
// 注意：HTTP 列表路径已改为 authctx + service 注入 req，新代码勿在 Controller 调用本函数写 DataScope；保留供尚未迁移的 handler / 脚本使用。
func ListRBAC(c *gin.Context) (dataScope int8, userID uint, roleCode string) {
	userID = GetUserID(c)
	roleCode = GetRoleCode(c)
	dataScope = GetDataScope(c)
	return dataScope, userID, roleCode
}
