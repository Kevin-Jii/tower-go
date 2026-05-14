package middleware

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

// Permission 按权限码进行接口鉴权
func Permission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 超管始终放行；数据范围另见 HQUnboundAdmin / GetDataScope
		if GetRoleCode(c) == model.RoleCodeSuperAdmin || HQUnboundAdmin(c) {
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
		if GetRoleCode(c) == model.RoleCodeSuperAdmin || HQUnboundAdmin(c) {
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

// PermissionStoreBoundSupplierRead 门店端只读本店已绑定供应商下的商品/分类/供应商列表等。
// 具备供应商菜单、库存、采购或记账等相关权限之一即可，避免仅有 inventory:list 时无法加载可采购商品与库存选品。
func PermissionStoreBoundSupplierRead() gin.HandlerFunc {
	return PermissionAny(
		"supplier:list",
		"inventory:list",
		"inventory:in",
		"inventory:out",
		"inventory:record",
		"purchase:list",
		"purchase:add",
		"purchase:edit",
		"store:account:list",
		"store:account:add",
		"store:account:edit",
		"store:list",
		"store:edit",
		"store:menu",
	)
}
