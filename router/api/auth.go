package api

import "github.com/gin-gonic/gin"

// RegisterAuthRoutes 注册认证路由
func RegisterAuthRoutes(v1 *gin.RouterGroup, c *Controllers) {
	auth := v1.Group("/auth")
	{
		auth.POST("/register", c.User.Register)
		auth.POST("/login", c.User.Login)
	}
}
