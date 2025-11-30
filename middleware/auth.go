package middleware

import (
	"strings"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/utils/auth"
	"github.com/Kevin-Jii/tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			http.Error(c, 401, "Authorization header is required")
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			http.Error(c, 401, "Invalid authorization header format")
			c.Abort()
			return
		}

		// 解析token
		claims, err := auth.ParseToken(parts[1])
		if err != nil {
			http.Error(c, 401, "Invalid token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文（优先使用 Token 中的值，安全可靠）
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("storeID", claims.StoreID)
		c.Set("roleCode", claims.RoleCode)
		c.Set("roleID", claims.RoleID)

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

// IsAdmin 判断是否是总部管理员
func IsAdmin(c *gin.Context) bool {
	return GetRoleCode(c) == model.RoleCodeAdmin
}
