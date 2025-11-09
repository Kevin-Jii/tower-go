package controller

import (
	"encoding/json"
	"fmt"
	"strconv"
	"tower-go/middleware"
	"tower-go/model"
	"tower-go/service"
	"tower-go/utils/http"

	"github.com/gin-gonic/gin"
)

type MenuController struct {
	menuService *service.MenuService
}

func NewMenuController(menuService *service.MenuService) *MenuController {
	return &MenuController{menuService: menuService}
}

// CreateMenu godoc
// @Summary 创建菜单
// @Description 创建新菜单（仅总部管理员）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param menu body model.CreateMenuReq true "菜单信息"
// @Success 200 {object} utils.Response
// @Router /menus [post]
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅总部管理员可以创建菜单")
		return
	}

	var req model.CreateMenuReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.menuService.CreateMenu(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetMenu godoc
// @Summary 获取菜单详情
// @Description 获取菜单详细信息
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "菜单ID"
// @Success 200 {object} utils.Response{data=model.Menu}
// @Router /menus/{id} [get]
func (c *MenuController) GetMenu(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu ID")
		return
	}

	menu, err := c.menuService.GetMenu(uint(id))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	jsonData, _ := json.MarshalIndent(menu, "", "  ")
	fmt.Printf("[GetMenu] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, menu)
}

// ListMenus godoc
// @Summary 菜单列表
// @Description 获取所有菜单列表（平铺形式）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=[]model.Menu}
// @Router /menus [get]
func (c *MenuController) ListMenus(ctx *gin.Context) {
	menus, err := c.menuService.ListMenus()
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	jsonData, _ := json.MarshalIndent(menus, "", "  ")
	fmt.Printf("[ListMenus] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, menus)
}

// GetMenuTree godoc
// @Summary 菜单树形结构
// @Description 获取菜单树形结构
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=[]model.Menu}
// @Router /menus/tree [get]
func (c *MenuController) GetMenuTree(ctx *gin.Context) {
	tree, err := c.menuService.GetMenuTree()
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Printf("[GetMenuTree] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, tree)
}

// UpdateMenu godoc
// @Summary 更新菜单
// @Description 更新菜单信息（仅总部管理员）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "菜单ID"
// @Param menu body model.UpdateMenuReq true "菜单信息"
// @Success 200 {object} utils.Response
// @Router /menus/{id} [put]
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅总部管理员可以更新菜单")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu ID")
		return
	}

	var req model.UpdateMenuReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.menuService.UpdateMenu(uint(id), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteMenu godoc
// @Summary 删除菜单
// @Description 删除菜单（仅总部管理员）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "菜单ID"
// @Success 200 {object} utils.Response
// @Router /menus/{id} [delete]
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅总部管理员可以删除菜单")
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid menu ID")
		return
	}

	if err := c.menuService.DeleteMenu(uint(id)); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// AssignMenusToRole godoc
// @Summary 为角色分配菜单
// @Description 为角色分配默认菜单权限（仅总部管理员）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param assignment body model.AssignMenusToRoleReq true "分配信息"
// @Success 200 {object} utils.Response
// @Router /menus/assign-role [post]
func (c *MenuController) AssignMenusToRole(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅总部管理员可以分配角色菜单")
		return
	}

	var req model.AssignMenusToRoleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.menuService.AssignMenusToRole(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetRoleMenus godoc
// @Summary 获取角色菜单
// @Description 获取角色的默认菜单权限（树形结构）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param role_id query int true "角色ID"
// @Success 200 {object} utils.Response{data=[]model.Menu}
// @Router /menus/role [get]
func (c *MenuController) GetRoleMenus(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Query("role_id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid role ID")
		return
	}

	tree, err := c.menuService.GetRoleMenuTree(uint(roleID))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Printf("[GetRoleMenus] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, tree)
}

// GetRoleMenuIDs godoc
// @Summary 获取角色菜单ID列表
// @Description 获取角色的所有菜单ID（用于权限回显）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param role_id query int true "角色ID"
// @Success 200 {object} utils.Response{data=[]uint}
// @Router /menus/role-ids [get]
func (c *MenuController) GetRoleMenuIDs(ctx *gin.Context) {
	roleID, err := strconv.ParseUint(ctx.Query("role_id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid role ID")
		return
	}

	menuIDs, err := c.menuService.GetRoleMenuIDs(uint(roleID))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	jsonData, _ := json.MarshalIndent(menuIDs, "", "  ")
	fmt.Printf("[GetRoleMenuIDs] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, menuIDs)
}

// AssignMenusToStoreRole godoc
// @Summary 为门店角色分配菜单
// @Description 为特定门店的角色定制菜单权限（总部管理员或门店管理员）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param assignment body model.AssignStoreMenusReq true "分配信息"
// @Success 200 {object} utils.Response
// @Router /menus/assign-store-role [post]
func (c *MenuController) AssignMenusToStoreRole(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	isAdmin := middleware.IsAdmin(ctx)

	var req model.AssignStoreMenusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 门店管理员只能配置自己门店的权限
	if !isAdmin && req.StoreID != storeID {
		http.Error(ctx, 403, "无权配置其他门店的菜单权限")
		return
	}

	if err := c.menuService.AssignMenusToStoreRole(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetStoreRoleMenus godoc
// @Summary 获取门店角色菜单
// @Description 获取特定门店角色的菜单权限（树形结构）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param store_id query int true "门店ID"
// @Param role_id query int true "角色ID"
// @Success 200 {object} utils.Response{data=[]model.Menu}
// @Router /menus/store-role [get]
func (c *MenuController) GetStoreRoleMenus(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Query("store_id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid store ID")
		return
	}

	roleID, err := strconv.ParseUint(ctx.Query("role_id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid role ID")
		return
	}

	tree, err := c.menuService.GetStoreRoleMenuTree(uint(storeID), uint(roleID))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Printf("[GetStoreRoleMenus] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, tree)
}

// GetStoreRoleMenuIDs godoc
// @Summary 获取门店角色菜单ID列表
// @Description 获取门店角色的所有菜单ID（用于权限回显）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param store_id query int true "门店ID"
// @Param role_id query int true "角色ID"
// @Success 200 {object} utils.Response{data=[]uint}
// @Router /menus/store-role-ids [get]
func (c *MenuController) GetStoreRoleMenuIDs(ctx *gin.Context) {
	storeID, err := strconv.ParseUint(ctx.Query("store_id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid store ID")
		return
	}

	roleID, err := strconv.ParseUint(ctx.Query("role_id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid role ID")
		return
	}

	menuIDs, err := c.menuService.GetStoreRoleMenuIDs(uint(storeID), uint(roleID))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	jsonData, _ := json.MarshalIndent(menuIDs, "", "  ")
	fmt.Printf("[GetStoreRoleMenuIDs] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, menuIDs)
}

// CopyStoreMenus godoc
// @Summary 复制门店菜单权限
// @Description 将一个门店的菜单配置复制到另一个门店（仅总部管理员）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Param copy body model.CopyStoreMenusReq true "复制信息"
// @Success 200 {object} utils.Response
// @Router /menus/copy-store [post]
func (c *MenuController) CopyStoreMenus(ctx *gin.Context) {
	if !middleware.IsAdmin(ctx) {
		http.Error(ctx, 403, "仅总部管理员可以复制菜单权限")
		return
	}

	var req model.CopyStoreMenusReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.menuService.CopyStoreMenus(&req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetUserMenus godoc
// @Summary 获取当前用户菜单
// @Description 获取当前登录用户的菜单树（用于渲染后台侧边栏）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=[]model.Menu}
// @Router /menus/user-menus [get]
func (c *MenuController) GetUserMenus(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)
	roleID := middleware.GetRoleID(ctx)

	// 调试日志：打印当前用户的权限信息
	fmt.Printf("[GetUserMenus] 用户权限信息 - storeID: %d, roleCode: %s, roleID: %d\n", storeID, roleCode, roleID)

	// 总部管理员获取所有菜单
	if roleCode == model.RoleCodeAdmin {
		tree, err := c.menuService.GetMenuTree()
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		jsonData, _ := json.MarshalIndent(tree, "", "  ")
		fmt.Printf("[GetUserMenus-Admin] 查询结果:\n%s\n", string(jsonData))
		http.Success(ctx, tree)
		return
	}

	// 其他用户根据门店和角色获取菜单
	tree, err := c.menuService.GetUserMenuTree(storeID, roleID)
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}
	jsonData, _ := json.MarshalIndent(tree, "", "  ")
	fmt.Printf("[GetUserMenus] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, tree)
}

// GetUserPermissions godoc
// @Summary 获取当前用户权限标识列表
// @Description 获取当前用户所有可访问的权限标识（用于前端按钮级权限控制）
// @Tags menus
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=[]string}
// @Router /menus/user-permissions [get]
func (c *MenuController) GetUserPermissions(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)
	roleID := middleware.GetRoleID(ctx)

	var permissions []string
	var err error

	// 总部管理员获取所有权限
	if roleCode == model.RoleCodeAdmin {
		permissions, err = c.menuService.GetAllPermissions()
	} else {
		// 其他用户根据门店和角色获取权限
		permissions, err = c.menuService.GetUserPermissions(storeID, roleID)
	}

	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	jsonData, _ := json.MarshalIndent(permissions, "", "  ")
	fmt.Printf("[GetUserPermissions] 查询结果:\n%s\n", string(jsonData))
	http.Success(ctx, permissions)
}
