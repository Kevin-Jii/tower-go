package middleware

import (
	"strings"

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

		// 同时保存请求头中的值（供日志或特殊场景使用）
		c.Set("headerUserID", c.GetHeader("X-User-Id"))
		c.Set("headerStoreID", c.GetHeader("X-Store-Id"))

		c.Next()
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

// IsAdmin 判断是否是总部管理员或超级管理员
func IsAdmin(c *gin.Context) bool {
	roleCode := GetRoleCode(c)
	return roleCode == model.RoleCodeAdmin || roleCode == model.RoleCodeSuperAdmin
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

// GetDataScope 当前请求的数据范围（管理员视为全部）
func GetDataScope(c *gin.Context) int8 {
	if IsAdmin(c) {
		return model.DataScopeAll
	}
	if r := GetRoleModel(c); r != nil {
		return r.DataScope
	}
	return model.DataScopeStore
}

// ListRBAC 列表接口注入：数据范围、当前用户、角色码（门店 ID 由调用方按是否管理员自行处理）
func ListRBAC(c *gin.Context) (dataScope int8, userID uint, roleCode string) {
	userID = GetUserID(c)
	roleCode = GetRoleCode(c)
	dataScope = GetDataScope(c)
	return dataScope, userID, roleCode
}
