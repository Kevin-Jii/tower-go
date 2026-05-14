package api

import (
	"errors"
	"strconv"

	"github.com/Kevin-Jii/tower-go/controller"
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/gin-gonic/gin"
)

type updateRoleRequest struct {
	ID uint `json:"id" binding:"required"`
	model.UpdateRoleReq
}

// RegisterPermissionRoutes 注册 RBAC 标准化权限路由
func RegisterPermissionRoutes(v1 *gin.RouterGroup, c *Controllers) {
	permission := v1.Group("/permission")
	permission.Use(middleware.AuthMiddleware())

	permission.POST("/role/create", middleware.Permission("system:role:add"), controller.CreateRole)
	permission.PUT("/role/update", middleware.Permission("system:role:edit"), func(ctx *gin.Context) {
		var req updateRoleRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			http.Error(ctx, 400, err.Error())
			return
		}
		role, err := service.UpdateRole(req.ID, &req.UpdateRoleReq)
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		http.Success(ctx, role)
	})
	permission.DELETE("/role/delete", middleware.Permission("system:role:delete"), func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Query("id"), 10, 32)
		if err != nil || id == 0 {
			http.ErrorApp(ctx, apicode.InvalidID)
			return
		}
		if err := service.DeleteRole(uint(id)); err != nil {
			if errors.Is(err, service.ErrBuiltinRoleNotDeletable) {
				http.ErrorApp(ctx, apicode.BuiltinRoleNotDeletable)
				return
			}
			http.Error(ctx, 500, err.Error())
			return
		}
		http.Success(ctx, nil)
	})
	permission.GET("/role/page", middleware.Permission("system:role:list"), controller.ListRoles)

	permission.POST("/menu/create", middleware.Permission("system:menu:add"), c.Menu.CreateMenu)
	permission.GET("/menu/list", middleware.Permission("system:menu:list"), c.Menu.GetMenuTree)

	permission.POST("/assign-role-menu", middleware.Permission("system:role:menu"), c.Menu.AssignMenusToRole)
	permission.POST("/assign-user-role", middleware.Permission("system:user:edit"), c.User.AssignUserRole)
}
