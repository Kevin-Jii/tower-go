package middleware

import (
	"net/http"
	"strings"

	"tower-go/model"
	"tower-go/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Error(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Error(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("storeID", claims.StoreID)
		c.Set("roleCode", claims.RoleCode)
		c.Set("roleID", claims.RoleID)
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
