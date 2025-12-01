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
		menus.POST("", c.Menu.CreateMenu)
		menus.GET("", c.Menu.ListMenus)
		menus.GET("/tree", c.Menu.GetMenuTree)
		menus.GET("/:id", c.Menu.GetMenu)
		menus.PUT("/:id", c.Menu.UpdateMenu)
		menus.DELETE("/:id", c.Menu.DeleteMenu)
		menus.POST("/assign-role", c.Menu.AssignMenusToRole)
		menus.GET("/role", c.Menu.GetRoleMenus)
		menus.GET("/role-ids", c.Menu.GetRoleMenuIDs)
		menus.GET("/role-permissions", c.Menu.GetRoleMenuPermissions)
		menus.POST("/assign-store-role", c.Menu.AssignMenusToStoreRole)
		menus.GET("/store-role", c.Menu.GetStoreRoleMenus)
		menus.GET("/store-role-ids", c.Menu.GetStoreRoleMenuIDs)
		menus.GET("/store-role-permissions", c.Menu.GetStoreRoleMenuPermissions)
		menus.POST("/copy-store", c.Menu.CopyStoreMenus)
		menus.GET("/user-menus", c.Menu.GetUserMenus)
		menus.GET("/user-permissions", c.Menu.GetUserPermissions)
	}

}
