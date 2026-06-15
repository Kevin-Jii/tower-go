package api

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterMenuRoutes 注册菜单权限路由
func RegisterMenuRoutes(v1 *gin.RouterGroup, c *Controllers) {
	menus := v1.Group("/menus")
	menus.Use(middleware.AuthMiddleware())
	{
		menus.POST("", middleware.Permission("system:menu:add"), c.Menu.CreateMenu)
		menus.GET("", middleware.Permission("system:menu:list"), c.Menu.ListMenus)
		menus.GET("/tree", middleware.PermissionAny("system:menu:list", "store:menu"), c.Menu.GetMenuTree)
		menus.GET("/:id", middleware.Permission("system:menu:list"), c.Menu.GetMenu)
		menus.PUT("/:id", middleware.Permission("system:menu:edit"), c.Menu.UpdateMenu)
		menus.DELETE("/:id", middleware.Permission("system:menu:delete"), c.Menu.DeleteMenu)
		menus.POST("/assign-role", middleware.Permission("system:role:menu"), c.Menu.AssignMenusToRole)
		menus.GET("/role", middleware.Permission("system:role:list"), c.Menu.GetRoleMenus)
		menus.GET("/role-ids", middleware.Permission("system:role:list"), c.Menu.GetRoleMenuIDs)
		menus.GET("/role-permissions", middleware.Permission("system:role:list"), c.Menu.GetRoleMenuPermissions)
		menus.POST("/assign-store-role", middleware.Permission("store:menu"), c.Menu.AssignMenusToStoreRole)
		menus.GET("/store-role", middleware.Permission("store:menu"), c.Menu.GetStoreRoleMenus)
		menus.GET("/store-role-ids", middleware.Permission("store:menu"), c.Menu.GetStoreRoleMenuIDs)
		menus.GET("/store-role-permissions", middleware.Permission("store:menu"), c.Menu.GetStoreRoleMenuPermissions)
		menus.POST("/copy-store", middleware.Permission("store:menu"), c.Menu.CopyStoreMenus)
		menus.GET("/user-menus", c.Menu.GetUserMenus)
		menus.GET("/user-permissions", c.Menu.GetUserPermissions)
	}

}
