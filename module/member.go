package module

import (
	"errors"
	"time"

	"github.com/Kevin-Jii/tower-go/model"
	updatesPkg "github.com/Kevin-Jii/tower-go/utils/updates"
	"github.com/google/uuid"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MemberModule 会员模块
type MemberModule struct {
	db *gorm.DB
}

// NewMemberModule 创建会员模块
func NewMemberModule(db *gorm.DB) *MemberModule {
	return &MemberModule{db: db}
}

// GetDB 返回底层数据库实例
func (m *MemberModule) GetDB() *gorm.DB {
	return m.db
}

// ========== Member 操作 ==========

// CreateMember 创建会员
func (m *MemberModule) CreateMember(req *model.CreateMemberReq) (*model.Member, error) {
	// 检查手机号是否已存在
	var count int64
	m.db.Model(&model.Member{}).Where("phone = ?", req.Phone).Count(&count)
	if count > 0 {
		return nil, errors.New("手机号已注册")
	}

	// 如果没有提供 UID，则自动生成
	uid := req.UID
	if uid == "" {
		uid = uuid.NewString()
	}
	member := &model.Member{
		UID:     uid,
		Phone:   req.Phone,
		Balance: model.DecimalZero(),
		Points:  0,
		Level:   1,
		Version: 0,
	}
	if err := m.db.Create(member).Error; err != nil {
		return nil, err
	}
	return member, nil
}

// UpdateMember 更新会员
func (m *MemberModule) UpdateMember(id uint, req *model.UpdateMemberReq) (*model.Member, error) {
	var member model.Member
	if err := m.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	updateMap := updatesPkg.BuildUpdatesFromReq(req)
	if len(updateMap) == 0 {
		return &member, nil
	}
	if err := m.db.Model(&member).Updates(updateMap).Error; err != nil {
		return nil, err
	}
	if err := m.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// DeleteMember 删除会员
func (m *MemberModule) DeleteMember(id uint) error {
	return m.db.Delete(&model.Member{}, id).Error
}

// GetMember 获取会员
func (m *MemberModule) GetMember(id uint) (*model.Member, error) {
	var member model.Member
	if err := m.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMemberByPhone 通过手机号获取会员
func (m *MemberModule) GetMemberByPhone(phone string) (*model.Member, error) {
	var member model.Member
	if err := m.db.Where("phone = ?", phone).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMemberByUID 通过UID获取会员
func (m *MemberModule) GetMemberByUID(uid string) (*model.Member, error) {
	var member model.Member
	if err := m.db.Where("uid = ?", uid).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// ListMembers 获取会员列表
func (m *MemberModule) ListMembers(keyword string) ([]model.Member, error) {
	var members []model.Member
	query := m.db.Model(&model.Member{})
	if keyword != "" {
		query = query.Where("phone LIKE ? OR uid LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	query = query.Order("id DESC")
	if err := query.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// AdjustBalanceWithLock 乐观锁调整余额
func (m *MemberModule) AdjustBalanceWithLock(id uint, amount model.DecimalType, changeType model.ChangeTypeEnum, remark string, version int) (*model.Member, error) {
	var member model.Member
	// 使用行锁查询
	if err := m.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&member, id).Error; err != nil {
		return nil, err
	}

	// 乐观锁校验
	if member.Version != version {
		return nil, errors.New("数据已被修改，请刷新后重试")
	}

	// 计算新余额
	var newBalance model.DecimalType
	if changeType == model.ChangeTypeAdjustAdd {
		newBalance = member.Balance.Add(amount)
	} else if changeType == model.ChangeTypeAdjustLess {
		newBalance = member.Balance.Sub(amount)
		if newBalance.LessThan(model.DecimalZero()) {
			return nil, errors.New("余额不足")
		}
	} else {
		return nil, errors.New("不支持的调整类型")
	}

	// 更新余额和版本号
	if err := m.db.Model(&member).Updates(map[string]interface{}{
		"balance": newBalance,
		"version": member.Version + 1,
	}).Error; err != nil {
		return nil, err
	}

	// 记录流水
	walletLog := &model.WalletLog{
		MemberID:       id,
		ChangeType:     changeType,
		ChangeAmount:   amount,
		BalanceAfter:   newBalance,
		RelatedOrderNo: "",
		Remark:         remark,
	}
	if err := m.db.Create(walletLog).Error; err != nil {
		return nil, err
	}

	// 重新查询
	if err := m.db.First(&member, id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// ========== WalletLog 操作 ==========

// CreateWalletLog 创建流水记录
func (m *MemberModule) CreateWalletLog(log *model.WalletLog) (*model.WalletLog, error) {
	if err := m.db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

// GetWalletLog 获取流水详情
func (m *MemberModule) GetWalletLog(id uint) (*model.WalletLog, error) {
	var log model.WalletLog
	if err := m.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// ListWalletLogs 查询流水列表
func (m *MemberModule) ListWalletLogs(req *model.ListWalletLogReq, page, pageSize int) ([]model.WalletLog, int64, error) {
	var logs []model.WalletLog
	var total int64

	query := m.db.Model(&model.WalletLog{})
	if req.MemberID > 0 {
		query = query.Where("member_id = ?", req.MemberID)
	}
	if req.ChangeType != nil {
		query = query.Where("change_type = ?", *req.ChangeType)
	}
	if req.StartTime != nil {
		query = query.Where("create_time >= ?", *req.StartTime)
	}
	if req.EndTime != nil {
		query = query.Where("create_time <= ?", *req.EndTime)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	query = query.Order("id DESC")
	query = query.Offset((page - 1) * pageSize)
	query = query.Limit(pageSize)

	if err := query.Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ========== RechargeOrder 操作 ==========

// CreateRechargeOrder 创建充值单
func (m *MemberModule) CreateRechargeOrder(req *model.CreateRechargeOrderReq) (*model.RechargeOrder, error) {
	orderNo := generateOrderNo()
	order := &model.RechargeOrder{
		OrderNo:     orderNo,
		MemberID:    req.MemberID,
		PayAmount:   req.PayAmount,
		GiftAmount:  req.GiftAmount,
		TotalAmount: req.PayAmount.Add(req.GiftAmount),
		PayStatus:   model.PayStatusPending,
	}
	if err := m.db.Create(order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

// GetRechargeOrder 获取充值单
func (m *MemberModule) GetRechargeOrder(id uint) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	if err := m.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetRechargeOrderByNo 通过单号获取充值单
func (m *MemberModule) GetRechargeOrderByNo(orderNo string) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	if err := m.db.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// ListRechargeOrders 查询充值单列表
func (m *MemberModule) ListRechargeOrders(memberID uint, status *model.PayStatusEnum, page, pageSize int) ([]model.RechargeOrder, int64, error) {
	var orders []model.RechargeOrder
	var total int64

	query := m.db.Model(&model.RechargeOrder{})
	if memberID > 0 {
		query = query.Where("member_id = ?", memberID)
	}
	if status != nil {
		query = query.Where("pay_status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query = query.Order("id DESC")
	query = query.Offset((page - 1) * pageSize)
	query = query.Limit(pageSize)

	if err := query.Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// PayRechargeOrder 支付充值单（余额充值）
func (m *MemberModule) PayRechargeOrder(orderNo string) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	if err := m.db.Where("order_no = ? AND pay_status = ?", orderNo, model.PayStatusPending).First(&order).Error; err != nil {
		return nil, err
	}

	// 获取会员并加锁
	var member model.Member
	if err := m.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(&member, order.MemberID).Error; err != nil {
		return nil, err
	}

	// 计算新余额
	newBalance := member.Balance.Add(order.TotalAmount)

	// 更新会员余额
	if err := m.db.Model(&member).Updates(map[string]interface{}{
		"balance": newBalance,
		"version": member.Version + 1,
	}).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	// 更新订单状态
	if err := m.db.Model(&order).Updates(map[string]interface{}{
		"pay_status": model.PayStatusPaid,
		"pay_time":   &now,
	}).Error; err != nil {
		return nil, err
	}

	// 记录流水
	walletLog := &model.WalletLog{
		MemberID:       order.MemberID,
		ChangeType:     model.ChangeTypeRecharge,
		ChangeAmount:   order.TotalAmount,
		BalanceAfter:   newBalance,
		RelatedOrderNo: orderNo,
		Remark:         "充值",
	}
	if err := m.db.Create(walletLog).Error; err != nil {
		return nil, err
	}

	// 重新查询订单
	if err := m.db.First(&order, order.ID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// CancelRechargeOrder 取消充值单
func (m *MemberModule) CancelRechargeOrder(orderNo string) (*model.RechargeOrder, error) {
	var order model.RechargeOrder
	if err := m.db.Where("order_no = ? AND pay_status = ?", orderNo, model.PayStatusPending).First(&order).Error; err != nil {
		return nil, err
	}

	if err := m.db.Model(&order).Update("pay_status", model.PayStatusCancelled).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

// generateOrderNo 生成订单号
func generateOrderNo() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

// randomString 生成随机字符串
func randomString(n int) string {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
