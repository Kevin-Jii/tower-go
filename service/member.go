package service

import (
	"github.com/Kevin-Jii/tower-go/model"
	"github.com/Kevin-Jii/tower-go/module"
)

// MemberService 会员服务
type MemberService struct {
	module *module.MemberModule
}

// NewMemberService 创建会员服务
func NewMemberService(m *module.MemberModule) *MemberService {
	return &MemberService{module: m}
}

// ========== Member 操作 ==========

// CreateMember 创建会员
func (s *MemberService) CreateMember(req *model.CreateMemberReq) (*model.Member, error) {
	return s.module.CreateMember(req)
}

// UpdateMember 更新会员
func (s *MemberService) UpdateMember(id uint, req *model.UpdateMemberReq) (*model.Member, error) {
	return s.module.UpdateMember(id, req)
}

// DeleteMember 删除会员
func (s *MemberService) DeleteMember(id uint) error {
	return s.module.DeleteMember(id)
}

// GetMember 获取会员
func (s *MemberService) GetMember(id uint) (*model.Member, error) {
	return s.module.GetMember(id)
}

// GetMemberByPhone 通过手机号获取会员
func (s *MemberService) GetMemberByPhone(phone string) (*model.Member, error) {
	return s.module.GetMemberByPhone(phone)
}

// GetMemberByUID 通过UID获取会员
func (s *MemberService) GetMemberByUID(uid string) (*model.Member, error) {
	return s.module.GetMemberByUID(uid)
}

// ListMembers 获取会员列表
func (s *MemberService) ListMembers(keyword string, page, pageSize int) ([]model.Member, int64, error) {
	return s.module.ListMembers(keyword, page, pageSize)
}

// AdjustBalance 调整余额
func (s *MemberService) AdjustBalance(id uint, amount model.DecimalType, changeType model.ChangeTypeEnum, remark string, version int) (*model.Member, error) {
	return s.module.AdjustBalanceWithLock(id, amount, changeType, remark, version)
}

// ========== WalletLog 操作 ==========

// CreateWalletLog 创建流水记录
func (s *MemberService) CreateWalletLog(log *model.WalletLog) (*model.WalletLog, error) {
	return s.module.CreateWalletLog(log)
}

// GetWalletLog 获取流水详情
func (s *MemberService) GetWalletLog(id uint) (*model.WalletLog, error) {
	return s.module.GetWalletLog(id)
}

// ListWalletLogs 查询流水列表
func (s *MemberService) ListWalletLogs(req *model.ListWalletLogReq, page, pageSize int) ([]model.WalletLog, int64, error) {
	return s.module.ListWalletLogs(req, page, pageSize)
}

// ========== RechargeOrder 操作 ==========

// CreateRechargeOrder 创建充值单
func (s *MemberService) CreateRechargeOrder(req *model.CreateRechargeOrderReq) (*model.RechargeOrder, error) {
	return s.module.CreateRechargeOrder(req)
}

// GetRechargeOrder 获取充值单
func (s *MemberService) GetRechargeOrder(id uint) (*model.RechargeOrder, error) {
	return s.module.GetRechargeOrder(id)
}

// GetRechargeOrderByNo 通过单号获取充值单
func (s *MemberService) GetRechargeOrderByNo(orderNo string) (*model.RechargeOrder, error) {
	return s.module.GetRechargeOrderByNo(orderNo)
}

// ListRechargeOrders 查询充值单列表
func (s *MemberService) ListRechargeOrders(memberID uint, status *model.PayStatusEnum, page, pageSize int) ([]model.RechargeOrder, int64, error) {
	return s.module.ListRechargeOrders(memberID, status, page, pageSize)
}

// PayRechargeOrder 支付充值单
func (s *MemberService) PayRechargeOrder(orderNo string) (*model.RechargeOrder, error) {
	return s.module.PayRechargeOrder(orderNo)
}

// CancelRechargeOrder 取消充值单
func (s *MemberService) CancelRechargeOrder(orderNo string) (*model.RechargeOrder, error) {
	return s.module.CancelRechargeOrder(orderNo)
}
