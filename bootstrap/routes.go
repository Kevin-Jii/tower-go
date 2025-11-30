package bootstrap

import (
	"github.com/Kevin-Jii/tower-go/router"
	"github.com/Kevin-Jii/tower-go/router/api"
	"github.com/gin-gonic/gin"
)

// BuildControllers 构建控制器（委托给 router/api 包）
func BuildControllers() *api.Controllers {
	return api.BuildControllers()
}

// RegisterRoutes 注册路由（委托给 router 包）
func RegisterRoutes(r *gin.Engine, c *api.Controllers) {
	router.Setup(r, c)
}
