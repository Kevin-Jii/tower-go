package service

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
	"github.com/Kevin-Jii/tower-go/utils"
	"github.com/Kevin-Jii/tower-go/utils/auth"
)

type UserService struct {
	userModule  *module.UserModule
	storeModule *module.StoreModule
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
		return errors.New("phone number already exists")
	}

	// 2. 生成唯一工号
	db := s.userModule.GetDB() // 假设 Module 层有 GetDB() 方法
	employeeNo, err := utils.GenerateEmployeeNo(db)
	if err != nil {
		return errors.New("生成工号失败: " + err.Error())
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
			return errors.New("无效门店编码: " + err.Error())
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
		return errors.New("user not found or access denied")
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
			return errors.New("无效门店编码: " + err.Error())
		}
		req.StoreID = &sid
		req.StoreCode = ""
	}
	return s.userModule.UpdateByID(id, req)
}

// ValidateUser 登录验证，跨门店查询，用于身份识别。
// **🔑 关键：必须返回包含 StoreID 的 User 对象**
func (s *UserService) ValidateUser(phone, password string) (*model.User, error) {
	// 获取用户信息 (Module 层全局查询)
	user, err := s.userModule.GetByPhone(phone)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user.Status == 2 {
		return nil, errors.New("该账号已经被禁用，请联系管理员")
	}

	// 检查门店状态：如果用户有门店且门店已停业，禁止登录
	if user.Store != nil && user.Store.Status == 2 {
		return nil, errors.New("该门店已停业，暂时无法登录")
	}

	// 验证密码
	if !auth.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid password")
	}

	// 更新最后登录时间（仅后端维护）
	loginTime := time.Now()
	user.LastLoginAt = &loginTime
	if err := s.userModule.Update(user); err != nil {
		return nil, err
	}

	// **🔑 关键：返回的 user 必须包含 StoreID 字段**
	return user, nil
}

// ResetPassword 重置指定用户密码为默认值（已加密）。
func (s *UserService) ResetPassword(userID uint, newPlain string) error {
	// 确认用户存在
	if _, err := s.userModule.GetByID(userID); err != nil {
		return errors.New("user not found")
	}
	hashed, err := auth.HashPassword(newPlain)
	if err != nil {
		return err
	}
	return s.userModule.UpdatePasswordByID(userID, hashed)
}
