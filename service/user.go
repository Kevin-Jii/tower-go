package service

import (
	"errors"
	"time"
	"tower-go/model"
	"tower-go/module"
	"tower-go/utils"
)

type UserService struct {
	userModule *module.UserModule
}

func NewUserService(userModule *module.UserModule) *UserService {
	return &UserService{userModule: userModule}
}

func (s *UserService) CreateUser(req *model.CreateUserReq) error {
	// 检查手机号是否已存在
	exists, err := s.userModule.ExistsByPhone(req.Phone)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("phone number already exists")
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Phone:    req.Phone,
		Password: hashedPassword,
		Nickname: req.Nickname,
		Email:    req.Email,
		Status:   1, // 默认状态为正常
	}
	return s.userModule.Create(user)
}

func (s *UserService) GetUser(id uint) (*model.User, error) {
	return s.userModule.GetByID(id)
}

func (s *UserService) ListUsers(page, pageSize int) ([]*model.User, int64, error) {
	return s.userModule.List(page, pageSize)
}

func (s *UserService) UpdateUser(id uint, req *model.UpdateUserReq) error {
	user, err := s.userModule.GetByID(id)
	if err != nil {
		return err
	}

	if req.Password != "" {
		// 密码加密
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	return s.userModule.Update(user)
}

func (s *UserService) ValidateUser(phone, password string) (*model.User, error) {
	// 获取用户信息
	user, err := s.userModule.GetByPhone(phone)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 验证密码
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid password")
	}

	// 更新最后登录时间
	user.LastLoginAt = time.Now()
	if err := s.userModule.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	return s.userModule.Delete(id)
}
