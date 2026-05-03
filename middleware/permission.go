package middleware

import (
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

// Permission 按权限码进行接口鉴权
func Permission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if HQUnboundAdmin(c) {
			c.Next()
			return
		}

		userID := GetUserID(c)
		storeID := GetStoreID(c)
		roleID := GetRoleID(c)
		roleCode := GetRoleCode(c)
		if userID == 0 || roleID == 0 {
			http.ErrorApp(c, apicode.PermissionDenied)
			c.Abort()
			return
		}

		perms, err := service.GetUserPermissionCodes(userID, storeID, roleID, roleCode)
		if err != nil {
			http.ErrorApp(c, apicode.PermissionLoadFailed)
			c.Abort()
			return
		}

		for _, p := range perms {
			if p == code {
				c.Next()
				return
			}
		}

		http.ErrorApp(c, apicode.PermissionDenied)
		c.Abort()
	}
}

// PermissionAny 满足任意一个权限码即可通过（用于同一接口允许多种操作权限的场景）
func PermissionAny(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if HQUnboundAdmin(c) {
			c.Next()
			return
		}

		userID := GetUserID(c)
		storeID := GetStoreID(c)
		roleID := GetRoleID(c)
		roleCode := GetRoleCode(c)
		if userID == 0 || roleID == 0 {
			http.ErrorApp(c, apicode.PermissionDenied)
			c.Abort()
			return
		}

		perms, err := service.GetUserPermissionCodes(userID, storeID, roleID, roleCode)
		if err != nil {
			http.ErrorApp(c, apicode.PermissionLoadFailed)
			c.Abort()
			return
		}

		allowed := make(map[string]struct{}, len(perms))
		for _, p := range perms {
			allowed[p] = struct{}{}
		}
		for _, code := range codes {
			if _, ok := allowed[code]; ok {
				c.Next()
				return
			}
		}

		http.ErrorApp(c, apicode.PermissionDenied)
		c.Abort()
	}
}
