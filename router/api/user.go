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
		users.POST("", c.User.CreateUser)
		users.GET("", c.User.ListUsers)
		users.GET("/:id", c.User.GetUser)
		users.PUT("/:id", c.User.UpdateUser)
		users.DELETE("/:id", c.User.DeleteUser)
		users.POST(":id/reset-password", c.User.ResetUserPassword)
	}
}

// RegisterRoleRoutes 注册角色管理路由
func RegisterRoleRoutes(v1 *gin.RouterGroup) {
	roles := v1.Group("/roles")
	roles.Use(middleware.AuthMiddleware())
	{
		roles.POST("", controller.CreateRole)
		roles.GET("", controller.ListRoles)
		roles.GET("/:id", controller.GetRole)
		roles.PUT("/:id", controller.UpdateRole)
		roles.DELETE("/:id", controller.DeleteRole)
	}
}
