package api

import (
	"github.com/Kevin-Jii/tower-go/controller"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes 注册用户管理路由
func RegisterUserRoutes(v1 *gin.RouterGroup, c *Controllers) {
	users := v1.Group("/users")
	users.Use(middleware.StoreAuthMiddleware())
	{
		users.GET("/profile", c.User.GetProfile)
		users.PUT("/profile", c.User.UpdateProfile)
		users.POST("", middleware.Permission("system:user:add"), c.User.CreateUser)
		users.GET("", middleware.Permission("system:user:list"), c.User.ListUsers)
		users.GET("/:id", middleware.Permission("system:user:list"), c.User.GetUser)
		users.PUT("/:id", middleware.Permission("system:user:edit"), c.User.UpdateUser)
		users.DELETE("/:id", middleware.Permission("system:user:delete"), c.User.DeleteUser)
		users.POST(":id/reset-password", middleware.Permission("system:user:edit"), c.User.ResetUserPassword)
	}
}

// RegisterRoleRoutes 注册角色管理路由
func RegisterRoleRoutes(v1 *gin.RouterGroup) {
	roles := v1.Group("/roles")
	roles.Use(middleware.AuthMiddleware())
	{
		roles.POST("", middleware.Permission("system:role:add"), controller.CreateRole)
		roles.GET("", middleware.Permission("system:role:list"), controller.ListRoles)
		roles.GET("/:id", middleware.Permission("system:role:list"), controller.GetRole)
		roles.PUT("/:id", middleware.Permission("system:role:edit"), controller.UpdateRole)
		roles.DELETE("/:id", middleware.Permission("system:role:delete"), controller.DeleteRole)
	}
}
