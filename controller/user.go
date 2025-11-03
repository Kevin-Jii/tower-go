package controller

import (
	"log"
	"net/http"
	"strconv"
	"tower-go/middleware"
	"tower-go/model"
	"tower-go/service"
	"tower-go/utils"

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

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// CreateUser godoc
// @Summary 创建用户
// @Description 创建新用户
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body model.CreateUserReq true "用户信息"
// @Success 200 {object} utils.StandardResponse
// @Router /users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)

	var req model.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.CreateUser(storeID, &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}

// GetUser godoc
// @Summary 获取用户
// @Description 获取用户详情
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} utils.StandardResponse{data=model.User}
// @Router /users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := c.userService.GetUser(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, user)
}

// ListUsers godoc
// @Summary 用户列表
// @Description 获取用户列表。总部管理员返回全部用户（跨门店，支持分页），门店管理员返回其门店用户（分页）
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param keyword query string false "模糊关键字(匹配用户名或手机号任意部分，如手机号后4位)"
// @Success 200 {object} utils.StandardResponse{data=[]model.User} "支持 keyword 模糊匹配用户名或手机号；总部管理员查看全部用户，门店管理员仅查看本门店用户"
// @Router /users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	roleCode := middleware.GetRoleCode(ctx)
	keyword := ctx.Query("keyword")
	page := utils.GetPage(ctx)
	pageSize := utils.GetPageSize(ctx)

	var (
		users []*model.User
		total int64
		err   error
	)

	// 总部管理员返回全部用户（跨门店），支持分页
	if roleCode == model.RoleCodeAdmin {
		users, total, err = c.userService.ListAllUsers(keyword, page, pageSize)
		if err != nil {
			utils.Error(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		utils.SuccessWithPagination(ctx, users, total, page, pageSize)
		return
	}

	// 门店管理员/普通员工：返回门店用户列表，分页
	storeID := middleware.GetStoreID(ctx)

	if keyword != "" {
		users, total, err = c.userService.ListUsersByStoreIDWithKeyword(storeID, keyword, page, pageSize)
	} else {
		users, total, err = c.userService.ListUsersByStoreID(storeID, page, pageSize)
	}
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessWithPagination(ctx, users, total, page, pageSize)
}

// UpdateUser godoc
// @Summary 更新用户
// @Description 更新用户信息
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Param user body model.UpdateUserReq true "用户信息"
// @Success 200 {object} utils.StandardResponse
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req model.UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.UpdateUser(uint(id), &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 删除用户
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} utils.StandardResponse
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	storeID := middleware.GetStoreID(ctx)
	idStr := ctx.Param("id")
	log.Printf("Attempting to delete user with ID string: %s", idStr)

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := c.userService.DeleteUserByStoreID(uint(id), storeID); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}

// ResetUserPassword godoc
// @Summary 重置用户密码
// @Description 将指定用户密码重置为 123456
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "用户ID"
// @Success 200 {object} utils.StandardResponse
// @Router /users/{id}/reset-password [post]
func (c *UserController) ResetUserPassword(ctx *gin.Context) {
	// 仅总部管理员或同门店管理员才能重置（此处：若为 admin 放行；否则必须该用户同门店）
	requesterStoreID := middleware.GetStoreID(ctx)
	requesterRoleCode := middleware.GetRoleCode(ctx)

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	targetUser, err := c.userService.GetUser(uint(id))
	if err != nil {
		utils.Error(ctx, http.StatusNotFound, "user not found")
		return
	}

	if requesterRoleCode != model.RoleCodeAdmin && targetUser.StoreID != requesterStoreID {
		utils.Error(ctx, http.StatusForbidden, "无权重置其他门店用户密码")
		return
	}

	if err := c.userService.ResetPassword(uint(id), "123456"); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}

// Register godoc
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags auth
// @Accept json
// @Produce json
// @Param user body model.CreateUserReq true "用户信息"
// @Success 200 {object} utils.StandardResponse
// @Router /auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req model.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 注册时默认不分配门店（由管理员后续分配）
	if err := c.userService.CreateUser(0, &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录并获取token
// @Tags auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "登录信息"
// @Success 200 {object} utils.StandardResponse{data=LoginResponse} "登录成功返回 token、token_type=Bearer、expires_in(秒)"
// @Router /auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userService.ValidateUser(req.Phone, req.Password)
	if err != nil {
		utils.Error(ctx, http.StatusUnauthorized, "Invalid phone number or password")
		return
	}

	// 获取角色代码和角色ID
	roleCode := ""
	roleID := uint(0)
	if user.Role != nil {
		roleCode = user.Role.Code
		roleID = user.RoleID
	}

	// 生成token（包含 StoreID、RoleCode 和 RoleID）
	token, expiresIn, err := utils.GenerateToken(user.ID, user.Username, user.StoreID, roleCode, roleID)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// 如果会话管理器策略为 single，则登录时踢出旧会话（可选增强）
	strategy := ""
	if sm := utils.GetSessionManager(); sm != nil {
		strategy = "single"
		if smSessions := sm.ListUserSessions(user.ID); len(smSessions) > 0 && smSessions[0] != nil {
			// 踢出全部旧会话
			sm.KickUser(user.ID, "login_replace")
		}
	}

	utils.Success(ctx, LoginResponse{
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: expiresIn,
		UserInfo:  user,
		Strategy:  strategy,
	})
}

// GetProfile godoc
// @Summary 获取用户个人信息
// @Description 获取当前登录用户的详细信息
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.StandardResponse{data=model.User}
// @Router /users/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")
	user, err := c.userService.GetUser(userID.(uint))
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, user)
}

// UpdateProfile godoc
// @Summary 更新个人信息
// @Description 更新当前登录用户的个人信息
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body model.UpdateUserReq true "用户信息"
// @Success 200 {object} utils.StandardResponse
// @Router /users/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, _ := ctx.Get("userID")

	var req model.UpdateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.UpdateUser(userID.(uint), &req); err != nil {
		utils.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, nil)
}
