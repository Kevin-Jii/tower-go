package service

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/config"
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/pkg/apicode"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/auth"
)

type UserService struct {
	userModule  *module.UserModule
	storeModule *module.StoreModule
}

type wechatCodeSessionResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func NewUserService(userModule *module.UserModule, storeModule *module.StoreModule) *UserService {
	return &UserService{userModule: userModule, storeModule: storeModule}
}

// --- 用户管理接口 (需要 StoreID 隔离) ---

// CreateUser 在指定 StoreID 下创建用户。
// 这里的 storeID 参数由 Controller 从 Token 中获取并传递。
// 注意：如果这是用户自注册接口，你需要调整逻辑以分配默认的 StoreID。
func (s *UserService) CreateUser(storeID uint, roleCode string, req *model.CreateUserReq) error {
	// 1. 检查手机号在【全局】或【当前门店】是否已存在 (取决于业务需求)
	// 假设我们在【全局】检查手机号唯一性
	exists, err := s.userModule.ExistsByPhone(req.Phone)
	if err != nil {
		return err
	}
	if exists {
		return apicode.New(apicode.PhoneAlreadyExists)
	}

	// 2. 生成唯一工号
	db := s.userModule.GetDB() // 假设 Module 层有 GetDB() 方法
	employeeNo, err := utils.GenerateEmployeeNo(db)
	if err != nil {
		return apicode.Wrap(apicode.InternalError, err)
	}

	// 3. 密码加密
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return err
	}

	targetStoreID := storeID
	if model.HQUnboundAdminRole(roleCode, storeID) && strings.TrimSpace(req.StoreCode) != "" {
		sid, err := s.storeModule.GetIDByStoreCode(req.StoreCode)
		if err != nil {
			return apicode.New(apicode.StoreNotFound)
		}
		targetStoreID = sid
	}

	user := &model.User{
		Phone:      req.Phone,
		Password:   hashedPassword,
		Username:   req.Username,
		Email:      req.Email,
		EmployeeNo: employeeNo,
		Status:     1, // 默认启用
		Gender:     1, // 默认男
		StoreID:    targetStoreID,
		Nickname:   req.Nickname,
		RoleID:     3, // 默认普通员工
	}

	if req.Gender == 2 {
		user.Gender = 2
	}

	if req.RoleID > 0 {
		user.RoleID = req.RoleID
	}

	if user.RoleID > 0 {
		var role model.Role
		if err := s.userModule.GetDB().Select("code").First(&role, user.RoleID).Error; err == nil && model.IsSuperAdminRole(role.Code) {
			user.StoreID = 0
		}
	}

	return s.userModule.Create(user)
}

// GetUserByStoreID 在指定门店下，根据用户ID获取单个用户。
// Module 层将负责使用 StoreID 限制查询。
func (s *UserService) GetUserByStoreID(userID uint, storeID uint) (*model.User, error) {
	// Module 层会使用 userID 和 storeID 进行复合查询
	return s.userModule.GetByUserIDAndStoreID(userID, storeID)
}

// ListUsersByStoreID 获取指定门店下的用户列表。
func (s *UserService) ListUsersByStoreID(storeID uint, page, pageSize int) ([]*model.User, int64, error) {
	// Module 层会使用 storeID 和分页参数进行隔离查询
	return s.userModule.ListByStoreID(storeID, page, pageSize)
}

// ListUsersByStoreIDWithKeyword 支持用户名或手机号模糊匹配
func (s *UserService) ListUsersByStoreIDWithKeyword(storeID uint, keyword string, page, pageSize int) ([]*model.User, int64, error) {
	return s.userModule.ListByStoreIDWithKeyword(storeID, keyword, page, pageSize)
}

// ListAllUsers 获取全部用户（支持分页，用于总部管理员）。storeID>0 时按门店筛选。
func (s *UserService) ListAllUsers(keyword string, storeID uint, page, pageSize int) ([]*model.User, int64, error) {
	return s.userModule.ListAllUsers(keyword, storeID, page, pageSize)
}

// UpdateUserByStoreID 更新指定门店下的用户数据。
func (s *UserService) UpdateUserByStoreID(userID uint, storeID uint, req *model.UpdateUserReq) error {
	// 1. 先获取用户，并确保用户属于该门店
	user, err := s.userModule.GetByUserIDAndStoreID(userID, storeID)
	if err != nil {
		// 如果用户不存在或不属于该门店，Module 层应该返回 'record not found'
		return apicode.New(apicode.UserNotFound)
	}

	// 2. 更新字段
	if req.Password != "" {
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Status != nil { // 允许设置 0 / 2 等
		log.Printf("[UserService.UpdateUserByStoreID] updating status to %d for user %d", *req.Status, user.ID)
		user.Status = *req.Status
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}

	// StoreID 在这里不需要更新，因为它在数据库中是固定的
	return s.userModule.Update(user)
}

// DeleteUserByStoreID 删除指定门店下的用户。
func (s *UserService) DeleteUserByStoreID(userID uint, storeID uint) error {
	// Module 层将负责在删除前，复合校验 userID 和 storeID
	return s.userModule.DeleteByUserIDAndStoreID(userID, storeID)
}

// DeleteUser 删除用户（管理员使用，不限制门店）
func (s *UserService) DeleteUser(userID uint) error {
	return s.userModule.Delete(userID)
}

// --- 个人档案 / 认证接口 (无需 StoreID 作为查询参数) ---

// GetUser 获取用户详情（用于 Profile 接口）
func (s *UserService) GetUser(id uint) (*model.User, error) {
	// Profile 接口访问的是用户自己的信息，直接使用 ID 即可
	return s.userModule.GetByID(id)
}

// UpdateUser 更新用户（管理员全量更新 / 个人资料等）
func (s *UserService) UpdateUser(id uint, req *model.UpdateUserReq) error {
	if strings.TrimSpace(req.StoreCode) != "" {
		sid, err := s.storeModule.GetIDByStoreCode(req.StoreCode)
		if err != nil {
			return apicode.New(apicode.StoreNotFound)
		}
		req.StoreID = &sid
		req.StoreCode = ""
	}
	if req.RoleID != nil && *req.RoleID > 0 {
		var role model.Role
		if err := s.userModule.GetDB().Select("code").First(&role, *req.RoleID).Error; err == nil && model.IsSuperAdminRole(role.Code) {
			zero := uint(0)
			req.StoreID = &zero
		}
	}
	if req.Password != "" {
		hashed, err := auth.HashPassword(req.Password)
		if err != nil {
			return err
		}
		req.Password = hashed
	}
	return s.userModule.UpdateByID(id, req)
}

// UpdateProfile 当前用户更新个人资料（仅允许昵称/邮箱/用户名/手机号/性别/密码）。
func (s *UserService) UpdateProfile(userID uint, req *model.UpdateUserReq) error {
	safe := &model.UpdateUserReq{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
	}
	if req.Password != "" {
		hashed, err := auth.HashPassword(req.Password)
		if err != nil {
			return err
		}
		safe.Password = hashed
	}
	return s.userModule.UpdateByID(userID, safe)
}

// ValidateUser 登录验证，跨门店查询，用于身份识别。
// **🔑 关键：必须返回包含 StoreID 的 User 对象**
func (s *UserService) ValidateUser(phone, password string) (*model.User, error) {
	// 获取用户信息 (Module 层全局查询)
	user, err := s.userModule.GetByPhone(phone)
	if err != nil {
		return nil, apicode.New(apicode.UserNotFound)
	}
	if user.Status == 2 {
		return nil, apicode.New(apicode.AccountDisabled)
	}

	// 检查门店状态：超级管理员不校验；其他账号若绑店且门店停业则禁止登录
	roleCode := ""
	if user.Role != nil {
		roleCode = user.Role.Code
	}
	if !model.IsSuperAdminRole(roleCode) && user.StoreID > 0 && user.Store != nil && user.Store.Status == 2 {
		return nil, apicode.New(apicode.StoreClosed)
	}

	// 验证密码
	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, apicode.New(apicode.InvalidCredentials)
	}

	// 更新最后登录时间（仅后端维护，避免 Save 全量写回）
	loginTime := time.Now()
	user.LastLoginAt = &loginTime
	if err := s.userModule.UpdateLastLoginAt(user.ID, loginTime); err != nil {
		return nil, err
	}

	// 超级管理员对外始终视为未绑店
	if model.IsSuperAdminRole(roleCode) {
		user.StoreID = 0
		user.Store = nil
	}

	// **🔑 关键：返回的 user 必须包含 StoreID 字段**
	return user, nil
}

func (s *UserService) ValidateWechatLogin(code string) (*model.User, error) {
	openID, err := s.exchangeWechatCode(code)
	if err != nil {
		return nil, err
	}
	user, err := s.userModule.GetByWechatOpenID(openID)
	if err != nil {
		return nil, apicode.New(apicode.WechatNotBound)
	}
	return s.prepareLoginUser(user)
}

func (s *UserService) BindWechatCode(userID uint, code string) error {
	openID, err := s.exchangeWechatCode(code)
	if err != nil {
		return err
	}
	if existing, err := s.userModule.GetByWechatOpenID(openID); err == nil && existing != nil && existing.ID != userID {
		return apicode.New(apicode.WechatAlreadyBound)
	}
	user, err := s.userModule.GetByID(userID)
	if err != nil || user == nil {
		return apicode.New(apicode.UserNotFound)
	}
	if user.Status == 2 {
		return apicode.New(apicode.AccountDisabled)
	}
	return s.userModule.BindWechatOpenID(userID, openID)
}

func (s *UserService) prepareLoginUser(user *model.User) (*model.User, error) {
	if user == nil {
		return nil, apicode.New(apicode.UserNotFound)
	}
	if user.Status == 2 {
		return nil, apicode.New(apicode.AccountDisabled)
	}
	roleCode := ""
	if user.Role != nil {
		roleCode = user.Role.Code
	}
	if !model.IsSuperAdminRole(roleCode) && user.StoreID > 0 && user.Store != nil && user.Store.Status == 2 {
		return nil, apicode.New(apicode.StoreClosed)
	}

	loginTime := time.Now()
	user.LastLoginAt = &loginTime
	if err := s.userModule.UpdateLastLoginAt(user.ID, loginTime); err != nil {
		return nil, err
	}
	if model.IsSuperAdminRole(roleCode) {
		user.StoreID = 0
		user.Store = nil
	}
	return user, nil
}

func (s *UserService) exchangeWechatCode(code string) (string, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return "", apicode.New(apicode.MissingParameter)
	}
	cfg := config.GetWechatConfig()
	if strings.TrimSpace(cfg.MiniAppID) == "" || strings.TrimSpace(cfg.MiniAppSecret) == "" {
		return "", apicode.New(apicode.ConfigMissing)
	}
	endpoint := "https://api.weixin.qq.com/sns/jscode2session"
	q := url.Values{}
	q.Set("appid", cfg.MiniAppID)
	q.Set("secret", cfg.MiniAppSecret)
	q.Set("js_code", code)
	q.Set("grant_type", "authorization_code")

	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Get(endpoint + "?" + q.Encode())
	if err != nil {
		return "", apicode.Wrap(apicode.ExternalServiceFailed, err)
	}
	defer resp.Body.Close()

	var body wechatCodeSessionResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", apicode.Wrap(apicode.ExternalServiceFailed, err)
	}
	if body.ErrCode != 0 {
		if body.ErrCode == 40029 {
			return "", apicode.Newf(apicode.InvalidCredentials, "微信登录失败: code无效")
		}
		return "", apicode.Newf(apicode.ExternalServiceFailed, "微信登录失败: %s", body.ErrMsg)
	}
	if strings.TrimSpace(body.OpenID) == "" {
		return "", apicode.New(apicode.ExternalServiceFailed)
	}
	return strings.TrimSpace(body.OpenID), nil
}

// ResetPassword 重置指定用户密码为默认值（已加密）。
func (s *UserService) ResetPassword(userID uint, newPlain string) error {
	// 确认用户存在
	if _, err := s.userModule.GetByID(userID); err != nil {
		return apicode.New(apicode.UserNotFound)
	}
	hashed, err := auth.HashPassword(newPlain)
	if err != nil {
		return err
	}
	return s.userModule.UpdatePasswordByID(userID, hashed)
}
