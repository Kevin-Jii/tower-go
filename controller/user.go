package controller

import (
	"github.com/Kevin-Jii/tower-go/middleware"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/service"
	"github.com/Kevin-Jii/tower-go/utils/auth"
	"github.com/Kevin-Jii/tower-go/utils/http"
	"github.com/Kevin-Jii/tower-go/utils/session"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required,len=11"`   // 手机号必须是11位
	Password string `json:"password" binding:"required,min=6"` // 密码至少6位
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token     string      `json:"token"`
	TokenType string      `json:"token_type"` // 固定 "Bearer"
	ExpiresIn int64       `json:"expires_in"` // 过期时间（秒）
	UserInfo  *model.User `json:"user_info"`
	Strategy  string      `json:"strategy"` // 会话策略 single/multi
}

// ResetPasswordRequest 重置密码请求（可后续扩展指定密码，当前固定无需 body）
type ResetPasswordRequest struct{}

type AssignUserRoleRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	RoleID uint `json:"role_id" binding:"required"`
}

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// CreateUser godoc
// @Summary 创建用户
// @Description 创建新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body model.CreateUserReq true "用户信息"
// @Success 200 {object} http.Response
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	roleCode := middleware.GetRoleCode(ctx)

	var req model.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 使用当前登录用户的门店ID（从Token获取）
	if err := c.userService.CreateUser(storeID, roleCode, &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// GetUser godoc
// @Summary 获取用户
// @Description 获取用户详情
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} http.Response{data=model.User}
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid user ID")
		return
	}

	user, err := c.userService.GetUser(uint(id))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, user)
}

// ListUsers godoc
// @Summary 用户列表
// @Description 获取用户列表。超级管理员或未绑定门店的总部管理员可跨店分页；其余账号仅返回 Token 绑定门店下的用户（分页）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "模糊关键字匹配用户名或手机号任意部分，如手机号135"
// @Param store_id query int false "门店ID，仅未绑定门店的总部管理员或超级管理员有效；>0 时只返回该门店用户"
// @Success 200 {object} http.Response{data=[]model.User} "支持 keyword 模糊匹配用户名或手机号；未绑定门店的总部管理员与超级管理员可跨店；其余账号仅本店"
// @Router /users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	page := http.GetPage(ctx)
	pageSize := http.GetPageSize(ctx)
	storeIDCtx := middleware.GetStoreID(ctx)

	var (
		users []*model.User
		total int64
		err   error
	)

	parseFilterStore := func() uint {
		if sidStr := ctx.Query("store_id"); sidStr != "" {
			if sid, perr := strconv.ParseUint(sidStr, 10, 32); perr == nil && sid > 0 {
				return uint(sid)
			}
		}
		return 0
	}

	// 未绑定门店的总部 admin / super_admin：跨店全量，支持按门店筛选；绑店账号一律走下方本店分支
	if middleware.HQUnboundAdmin(ctx) {
		filterStore := parseFilterStore()
		users, total, err = c.userService.ListAllUsers(keyword, filterStore, page, pageSize)
		if err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		http.SuccessWithPagination(ctx, users, total, page, pageSize)
		return
	}

	// 门店管理员、员工、或已绑定门店的 admin/super_admin：仅能查看本门店用户（忽略 store_id 查询参数）
	if storeIDCtx == 0 {
		http.Error(ctx, 403, "当前账号未绑定门店，无法查看用户列表")
		return
	}

	if keyword != "" {
		users, total, err = c.userService.ListUsersByStoreIDWithKeyword(storeIDCtx, keyword, page, pageSize)
	} else {
		users, total, err = c.userService.ListUsersByStoreID(storeIDCtx, page, pageSize)
	}
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.SuccessWithPagination(ctx, users, total, page, pageSize)
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param user body model.UpdateUserReq true "用户信息"
// @Success 200 {object} http.Response
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid user ID")
		return
	}

	var req model.UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.userService.UpdateUser(uint(id), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 删除用户（管理员可删除任意用户，门店管理员只能删除本门店用户）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} http.Response
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	idStr := ctx.Param("id")
	log.Printf("Attempting to delete user with ID string: %s", idStr)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid user ID")
		return
	}

	// 未绑定门店的总部 admin / super_admin：可删除任意用户
	if middleware.HQUnboundAdmin(ctx) {
		if err := c.userService.DeleteUser(uint(id)); err != nil {
			http.Error(ctx, 500, err.Error())
			return
		}
		http.Success(ctx, nil)
		return
	}

	// 门店侧账号：仅能删除本门店用户
	if storeID == 0 {
		http.Error(ctx, 403, "当前账号未绑定门店，无法删除用户")
		return
	}
	if err := c.userService.DeleteUserByStoreID(uint(id), storeID); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// ResetUserPassword godoc
// @Summary 重置用户密码
// @Description 将指定用户密码重置为 123456
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} http.Response
// @Router /users/{id}/reset-password [post]
func (c *UserController) ResetUserPassword(ctx *gin.Context) {
	// 未绑店总部 admin/super_admin 可跨店重置；否则仅本店
	requesterStoreID := middleware.GetStoreID(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		http.Error(ctx, 400, "Invalid user ID")
		return
	}

	targetUser, err := c.userService.GetUser(uint(id))
	if err != nil {
		http.Error(ctx, 404, "user not found")
		return
	}

	if !middleware.HQUnboundAdmin(ctx) && targetUser.StoreID != requesterStoreID {
		http.Error(ctx, 403, "无权重置其他门店用户密码")
		return
	}

	// 生成临时密码
	tempPassword, err := auth.GenerateStrongPassword(12)
	if err != nil {
		http.Error(ctx, 500, "生成临时密码失败")
		return
	}

	if err := c.userService.ResetPassword(uint(id), tempPassword); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, gin.H{
		"message":       "密码重置成功",
		"temp_password": tempPassword,
		"warning":       "请立即修改此临时密码",
	})
}

// Register godoc
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param user body model.CreateUserReq true "用户信息"
// @Success 200 {object} http.Response
// @Router /auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req model.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	// 注册时默认不分配门店（由管理员后续分配）
	if err := c.userService.CreateUser(0, "", &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录并获取token
// @Tags 认证
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登录信息"
// @Success 200 {object} http.Response{data=LoginResponse} "登录成功返回 token、token_type=Bearer、expires_in(秒)"
// @Router /auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	user, err := c.userService.ValidateUser(req.Phone, req.Password)
	if err != nil {
		http.Error(ctx, 401, "Invalid phone number or password")
		return
	}

	// 获取角色代码和角色ID
	roleCode := ""
	roleID := uint(0)
	if user.Role != nil {
		roleCode = user.Role.Code
		roleID = user.RoleID
	}

	// 生成token（包含StoreID、RoleCode 和 RoleID）
	token, expiresIn, err := auth.GenerateToken(user.ID, user.Username, user.StoreID, roleCode, roleID)
	if err != nil {
		http.Error(ctx, 500, "Failed to generate token")
		return
	}
	_, _ = service.BuildUserPermissionCache(user.ID, user.StoreID, roleID, roleCode)

	// 如果会话管理器策略为 single，则登录时踢出旧会话（可选增强）
	strategy := ""
	if sm := session.GetSessionManager(); sm != nil {
		strategy = "single"
		if smSessions := sm.ListUserSessions(user.ID); len(smSessions) > 0 && smSessions[0] != nil {
			// 踢出全部旧会话
			sm.KickUser(user.ID, "login_replace")
		}
	}

	http.Success(ctx, LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
		UserInfo:  user,
		Strategy:  strategy,
	})
}

// AssignUserRole 为用户分配角色
func (c *UserController) AssignUserRole(ctx *gin.Context) {
	var req AssignUserRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	updateReq := &model.UpdateUserReq{
		RoleID: &req.RoleID,
	}
	if err := c.userService.UpdateUser(req.UserID, updateReq); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	user, err := c.userService.GetUser(req.UserID)
	if err == nil && user != nil {
		roleCode := ""
		if user.Role != nil {
			roleCode = user.Role.Code
		}
		service.InvalidateUserPermissionCache(user.ID)
		_, _ = service.BuildUserPermissionCache(user.ID, user.StoreID, user.RoleID, roleCode)
	} else {
		service.InvalidateUserPermissionCache(req.UserID)
	}

	http.Success(ctx, nil)
}

// GetProfile godoc
// @Summary 获取用户个人信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} http.Response{data=model.User}
// @Router /users/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	user, err := c.userService.GetUser(userID.(uint))
	if err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, user)
}

// UpdateProfile godoc
// @Summary 更新个人信息
// @Description 更新当前登录用户的个人信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body model.UpdateUserReq true "用户信息"
// @Success 200 {object} http.Response
// @Router /users/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var req model.UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		http.Error(ctx, 400, err.Error())
		return
	}

	if err := c.userService.UpdateUser(userID.(uint), &req); err != nil {
		http.Error(ctx, 500, err.Error())
		return
	}

	http.Success(ctx, nil)
}
