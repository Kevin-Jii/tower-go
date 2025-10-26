package controller

import (
	"net/http"
	"strconv"
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
	Token    string      `json:"token"`
	UserInfo *model.User `json:"user_info"`
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
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.CreateUserReq true "用户信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req model.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.CreateUser(&req); err != nil {
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
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response{data=model.User}
// @Router /api/v1/users/{id} [get]
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
// @Description 获取用户列表，支持分页
// @Tags users
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} utils.Response{data=[]model.User}
// @Router /api/v1/users [get]
func (c *UserController) ListUsers(ctx *gin.Context) {
	page := utils.GetPage(ctx)
	pageSize := utils.GetPageSize(ctx)

	users, total, err := c.userService.ListUsers(page, pageSize)
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
// @Param id path int true "用户ID"
// @Param user body model.UpdateUserReq true "用户信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/users/{id} [put]
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
// @Param id path int true "用户ID"
// @Success 200 {object} utils.Response
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		utils.Error(ctx, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := c.userService.DeleteUser(uint(id)); err != nil {
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
// @Success 200 {object} utils.Response
// @Router /api/v1/auth/register [post]
func (c *UserController) Register(ctx *gin.Context) {
	var req model.CreateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.userService.CreateUser(&req); err != nil {
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
// @Success 200 {object} utils.Response{data=LoginResponse}
// @Router /api/v1/auth/login [post]
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

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.Error(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.Success(ctx, LoginResponse{
		Token:    token,
		UserInfo: user,
	})
}

// GetProfile godoc
// @Summary 获取用户个人信息
// @Description 获取当前登录用户的详细信息
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} utils.Response{data=model.User}
// @Router /api/v1/users/profile [get]
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
// @Description 更新当前登录用户的信息
// @Tags users
// @Accept json
// @Produce json
// @Security Bearer
// @Param user body model.UpdateUserReq true "用户信息"
// @Success 200 {object} utils.Response
// @Router /api/v1/users/profile [put]
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
